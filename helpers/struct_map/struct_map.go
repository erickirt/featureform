// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// Copyright 2024 FeatureForm Inc.
//

package struct_map

import (
	"fmt"
	"reflect"

	"github.com/featureform/fferr"
)

type StructMap struct {
	keys []string
	size int
	val  reflect.Value
}

func (sm StructMap) Get(key string) (interface{}, bool) {
	if idx := sm.indexOf(key); idx == -1 {
		return nil, false
	}
	return sm.val.FieldByName(key).Interface(), true
}

func (sm StructMap) Size() int {
	return sm.size
}

func (sm StructMap) Keys() []string {
	return sm.keys
}

func (sm StructMap) Has(key string, val interface{}) bool {
	v, ok := sm.Get(key)
	if !ok {
		return false
	}
	return reflect.DeepEqual(v, val)
}

func (sm StructMap) indexOf(key string) int {
	for i, k := range sm.keys {
		if k == key {
			return i
		}
	}
	return -1
}

func NewStructMap(s interface{}) (*StructMap, error) {
	t := reflect.TypeOf(s).Kind()
	if t != reflect.Struct {
		return nil, fferr.NewInternalError(fmt.Errorf("cannot create map from type %T", s))
	}

	v := reflect.ValueOf(s)
	keys := []string{}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if !field.CanInterface() {
			continue
		}
		keys = append(keys, v.Type().Field(i).Name)
	}

	return &StructMap{
		keys: keys,
		size: len(keys),
		val:  v,
	}, nil
}
