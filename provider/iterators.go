// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// Copyright 2024 FeatureForm Inc.
//

package provider

import (
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/parquet-go/parquet-go"

	_ "gocloud.dev/blob/fileblob"
	_ "gocloud.dev/blob/memblob"

	"github.com/featureform/fferr"
	filestore "github.com/featureform/filestore"
)

// PARQUET
type parquetIterator struct {
	reader        *parquet.Reader
	currentValues GenericRecord
	err           error
	// Using fields instead of column names gives us access to
	// parquet field metadata, which we need to determine whether
	// a field is a timestamp or not, for example
	fields []parquet.Field
	limit  int64
	idx    int64
}

func (p *parquetIterator) Next() bool {
	if p.idx >= p.limit {
		return false
	}
	row := make(map[string]interface{}, 0)
	err := p.reader.Read(&row)
	if err != nil {
		if err == io.EOF {
			return false
		} else {
			p.err = err
			return false
		}
	}
	records := make(GenericRecord, 0)
	for _, f := range p.fields {
		var recordVal interface{}
		switch assertedVal := row[f.Name()].(type) {
		// We're currently converting int32 to int to decrease/simplify the number of
		// data types we support; however, this process should involve a larger discussion
		// with the team to assess impact.
		case int32:
			recordVal = int(assertedVal)
		case int64:
			// Given we're instructing Spark to output timestamps as int64 (microseconds),
			// we need to rely on the parquet schema's field metadata to determine whether
			// the field is a timestamp or not. If it is, we need to convert it to its
			// corresponding Go type (time.Time).
			if reflect.DeepEqual(f.Type(), parquet.Timestamp(parquet.Millisecond).Type()) {
				recordVal = time.UnixMilli(assertedVal).UTC()
			} else {
				recordVal = int(assertedVal)
			}
		// This is the type that a []float32 is returned, so we have to parse it.
		case map[string]interface{}:
			vec, err := parseFloatVec(assertedVal)
			if err != nil {
				p.err = err
				return false
			}
			recordVal = vec
		default:
			recordVal = assertedVal
		}
		records = append(records, recordVal)
	}
	p.currentValues = records
	p.idx += 1
	return true
}

// parseFloatVec parses a generic float array that is received via a parquet file. It shows up in the form:
// map[list:[map[element: 1] map[element:2] map[element:3]]]
// Though, sometimes when people are using Databricks they tend to get VectorUDT types.
// In that situation its:
// map[indices: , type: , values: list:[map[element: 1] map[element:2] map[element:3]]]
func parseFloatVec(val map[string]interface{}) ([]float32, error) {
	if _, ok := val["indices"]; ok {
		unwrapped, err := unwrapVectorUDT(val)
		if err != nil {
			return nil, err
		}
		val = unwrapped
	}
	list, ok := val["list"]
	if !ok {
		return nil, fferr.NewDataTypeNotFoundErrorf(val, "expected to find field 'list' when parsing float vector")
	}
	// To iterate over the list, we need to cast it to []interface{}
	elementsSlice, ok := list.([]interface{})
	if !ok {
		return nil, fferr.NewDataTypeNotFoundErrorf(list, "failed to cast to []interface{} when parsing float vector")
	}
	vec := make([]float32, len(elementsSlice))
	for i, e := range elementsSlice {
		// To access the 'element' field, which holds the float value,
		// we need to cast it to map[string]interface{}
		m, ok := e.(map[string]interface{})
		if !ok {
			return nil, fferr.NewDataTypeNotFoundErrorf(e, "failed to cast to map[string]interface{} when parsing float vector")
		}

		switch casted := m["element"].(type) {
		case float32:
			vec[i] = casted
		case float64:
			vec[i] = float32(casted)
		case int:
			vec[i] = float32(casted)
		case int32:
			vec[i] = float32(casted)
		case int64:
			vec[i] = float32(casted)
		case string:
			parsedVec, err := strconv.ParseFloat(casted, 32)
			if err != nil {
				return nil, fferr.NewDataTypeNotFoundErrorf(casted, "unexpected type in parquet vector list when parsing float vector")
			}
			vec[i] = float32(parsedVec)
		default:
			return nil, fferr.NewDataTypeNotFoundErrorf(casted, "unexpected type in parquet vector list when parsing float vector")
		}
	}
	return vec, nil
}

func unwrapVectorUDT(val map[string]interface{}) (map[string]interface{}, error) {
	expKeys := []string{"indices", "type", "size", "values"}
	for _, key := range expKeys {
		if _, has := val[key]; !has {
			return val, nil
		}
	}
	if val["indices"] != nil {
		// VectorUDT supports spare vectors. We don't have a parser for them yet.
		return nil, fferr.NewInternalErrorf("SpareVectors not supported")
	}
	vec := val["values"]
	if vec == nil {
		// Unsure what this type is, lets let it ride
		return val, nil
	}
	if dict, ok := vec.(map[string]interface{}); ok {
		return dict, nil
	} else {
		// Unsure what this type is, lets let it ride
		return val, nil
	}
}

func (p *parquetIterator) Values() GenericRecord {
	return p.currentValues
}

func (p *parquetIterator) Columns() []string {
	columns := make([]string, len(p.fields))
	for i, field := range p.fields {
		columns[i] = field.Name()
	}
	return columns
}

func (p *parquetIterator) Err() error {
	return p.err
}

func (p *parquetIterator) Close() error {
	err := p.reader.Close()
	if err != nil {
		return fferr.NewInternalError(err)
	}
	return nil
}

func newParquetIterator(src io.ReaderAt, limit int64) (GenericTableIterator, error) {
	reader := parquet.NewReader(src)
	if limit == -1 {
		limit = math.MaxInt64
	}
	return &parquetIterator{
		reader: reader,
		fields: reader.Schema().Fields(),
		limit:  limit,
		idx:    0,
	}, nil
}

type multipleFileParquetIterator struct {
	iterator      *parquetIterator
	store         FileStore
	currentRecord GenericRecord
	err           error
	columns       []string
	limit         int64
	idx           int64
	files         []filestore.Filepath
	fileIdx       int64
}

func (p *multipleFileParquetIterator) Next() bool {
	hasNext := p.iterator.Next()
	if hasNext {
		p.currentRecord = p.iterator.Values()
		p.idx += 1
		return true
	}
	err := p.iterator.Err()
	if err != nil {
		p.err = err
		return false
	}
	if p.fileIdx >= int64(len(p.files)) {
		return false
	}
	if p.idx >= p.limit {
		return false
	}
	src, err := p.store.ReaderAt(p.files[p.fileIdx])
	if err != nil {
		p.err = err
		return false
	}
	updatedLimit := p.limit - p.idx
	iterator, err := newParquetIterator(src, updatedLimit)
	if err != nil {
		p.err = err
		return false
	}
	parquetIterator, isParquetIterator := iterator.(*parquetIterator)
	if !isParquetIterator {
		p.err = fferr.NewInternalError(fmt.Errorf("iterator is not a parquet iterator"))
		return false
	}
	p.iterator = parquetIterator
	p.fileIdx += 1
	return p.Next()
}

func (p *multipleFileParquetIterator) Values() GenericRecord {
	return p.currentRecord
}

func (p *multipleFileParquetIterator) Columns() []string {
	return p.columns
}

func (p *multipleFileParquetIterator) Err() error {
	return p.err
}

func (p *multipleFileParquetIterator) Close() error {
	err := p.iterator.reader.Close()
	if err != nil {
		return fferr.NewInternalError(err)
	}
	return nil
}

func newMultipleFileParquetIterator(files []filestore.Filepath, store FileStore, limit int64) (GenericTableIterator, error) {
	if len(files) == 0 {
		return nil, fferr.NewInvalidArgumentError(fmt.Errorf("no files to read"))
	}
	for _, file := range files {
		if file.Ext() != filestore.Parquet {
			return nil, fferr.NewInvalidArgumentError(fmt.Errorf("one or more files have an extension that is not .parquet: %s", file.Ext()))
		}
	}
	if limit == -1 {
		limit = math.MaxInt64
	}
	src, err := store.ReaderAt(files[0])
	if err != nil {
		return nil, err
	}
	iterator, err := newParquetIterator(src, limit)
	if err != nil {
		return nil, err
	}
	parquetIterator, isParquetIterator := iterator.(*parquetIterator)
	if !isParquetIterator {
		return nil, fferr.NewInternalError(fmt.Errorf("iterator is not a parquet iterator"))
	}
	return &multipleFileParquetIterator{
		iterator:      parquetIterator,
		store:         store,
		currentRecord: nil,
		columns:       parquetIterator.Columns(),
		limit:         limit,
		idx:           0,
		files:         files,
		fileIdx:       1,
	}, nil
}

type ParquetIteratorMultipleFiles struct {
	fileList       []filestore.Filepath
	currentIndex   int64
	fileIterator   Iterator
	featureColumns []string
	labelColumn    string
	store          FileStore
}

func parquetIteratorOverMultipleFiles(fileParts []filestore.Filepath, store FileStore) (Iterator, error) {
	src, err := store.ReaderAt(fileParts[0])
	if err != nil {
		return nil, err
	}
	iterator, err := parquetIteratorFromBytes(src)
	if err != nil {
		return nil, err
	}
	return &ParquetIteratorMultipleFiles{
		fileList:       fileParts,
		currentIndex:   int64(0),
		fileIterator:   iterator,
		store:          store,
		featureColumns: iterator.FeatureColumns(),
		labelColumn:    iterator.LabelColumn(),
	}, nil
}

func (p *ParquetIteratorMultipleFiles) FeatureColumns() []string {
	return p.featureColumns
}

func (p *ParquetIteratorMultipleFiles) LabelColumn() string {
	return p.labelColumn
}

func (p *ParquetIteratorMultipleFiles) Next() (map[string]interface{}, error) {
	nextRow, err := p.fileIterator.Next()
	if err != nil {
		return nil, err
	}
	if nextRow == nil {
		if p.currentIndex+1 == int64(len(p.fileList)) {
			return nil, nil
		}
		p.currentIndex += 1
		src, err := p.store.ReaderAt(p.fileList[p.currentIndex])
		if err != nil {
			return nil, err
		}
		iterator, err := parquetIteratorFromBytes(src)
		if err != nil {
			return nil, err
		}
		p.fileIterator = iterator
		return p.fileIterator.Next()
	}
	return nextRow, nil
}

type ParquetIterator struct {
	reader         *parquet.Reader
	index          int64
	featureColumns []string
	labelColumn    string
	fields         []parquet.Field
}

func (p *ParquetIterator) Next() (map[string]interface{}, error) {
	row := make(map[string]interface{})
	err := p.reader.Read(&row)
	if err != nil {
		if err == io.EOF {
			return nil, nil
		} else {
			return nil, fferr.NewInternalError(err)
		}
	}
	for _, f := range p.fields {
		switch assertedVal := row[f.Name()].(type) {
		case int32:
			row[f.Name()] = int(assertedVal)
		case int64:
			if reflect.DeepEqual(f.Type(), parquet.Timestamp(parquet.Millisecond).Type()) {
				// This check for a negative value is necessary because Spark uses a different
				// calendar than Go. For example. a 0 value in Spark starts at the year 0001; however,
				// in Go, a 0 value starts at 1970. This means that if we don't check for negative
				// values, we'll end up with a time.Time value that's prior to epoch start.
				if assertedVal < 0 {
					row[f.Name()] = time.UnixMilli(0).UTC()
				} else {
					row[f.Name()] = time.UnixMilli(assertedVal).UTC()
				}
			} else {
				row[f.Name()] = int(assertedVal)
			}
		case map[string]interface{}:
			vec, err := parseFloatVec(assertedVal)
			if err != nil {
				return nil, err
			}
			row[f.Name()] = vec
		default:
			row[f.Name()] = assertedVal
		}
	}
	return row, nil
}

func (p *ParquetIterator) FeatureColumns() []string {
	return p.featureColumns
}

func (p *ParquetIterator) LabelColumn() string {
	return p.labelColumn
}

func getParquetNumRows(src io.ReaderAt) (int64, error) {
	r := parquet.NewReader(src)
	return r.NumRows(), nil
}

type columnType string

const (
	labelType   columnType = "Label"
	featureType columnType = "Feature"
)

type parquetSchema struct {
	featureColumns []string
	labelColumn    string
	fields         []parquet.Field
}

func (s *parquetSchema) parseParquetColumnName(r *parquet.Reader) {
	s.fields = r.Schema().Fields()
	columnList := r.Schema().Columns()
	for _, column := range columnList {
		columnName := column[0]
		colType := s.getColumnType(columnName)
		s.setColumn(colType, columnName)
	}
}

func (s *parquetSchema) getColumnType(name string) columnType {
	columnSections := strings.Split(name, "__")
	return columnType(columnSections[0])
}

func (s *parquetSchema) setColumn(colType columnType, name string) {
	if colType == labelType {
		s.labelColumn = name
	} else if colType == featureType {
		s.featureColumns = append(s.featureColumns, name)
	}
}

func parquetIteratorFromBytes(src io.ReaderAt) (Iterator, error) {
	r := parquet.NewReader(src)
	schema := parquetSchema{}
	schema.parseParquetColumnName(r)
	return &ParquetIterator{
		reader:         r,
		index:          int64(0),
		featureColumns: schema.featureColumns,
		labelColumn:    schema.labelColumn,
		fields:         schema.fields,
	}, nil
}

// / CSV
type csvIterator struct {
	reader        *csv.Reader
	currentValues GenericRecord
	err           error
	columnNames   []string
	idx           int64
	limit         int64
}

func (c *csvIterator) Next() bool {
	if c.idx >= c.limit {
		return false
	}
	row, err := c.reader.Read()
	if err != nil {
		if err == io.EOF {
			return false
		} else {
			c.err = err
			return false
		}
	}
	c.currentValues = c.ParseRow(row)
	c.idx += 1
	return true
}

func (c *csvIterator) Values() GenericRecord {
	return c.currentValues
}

func (c *csvIterator) Columns() []string {
	return c.columnNames
}

func (c *csvIterator) Err() error {
	return c.err
}

func (c *csvIterator) Close() error {
	return nil
}

func (c *csvIterator) ParseRow(row []string) GenericRecord {
	records := make(GenericRecord, len(row))
	for i, value := range row {
		if integer, err := strconv.Atoi(value); err == nil {
			records[i] = integer
			continue
		}
		if float, err := strconv.ParseFloat(value, 64); err == nil {
			records[i] = float
			continue
		}
		records[i] = value
	}
	return records
}

func newCSVIterator(src io.Reader, limit int64) (GenericTableIterator, error) {
	reader := csv.NewReader(src)
	headers, err := reader.Read()
	if err != nil {
		return nil, fferr.NewInternalError(err)
	}
	if limit == -1 {
		limit = math.MaxInt64
	}
	return &csvIterator{
		reader:      reader,
		columnNames: headers,
		limit:       limit,
		idx:         0,
	}, nil
}
