// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// Copyright 2024 FeatureForm Inc.
//

package metadata

import pc "github.com/featureform/provider/provider_config"

func isValidBigQueryConfigUpdate(sa, sb pc.SerializedConfig) (bool, error) {
	a := pc.BigQueryConfig{}
	b := pc.BigQueryConfig{}
	if err := a.Deserialize(sa); err != nil {
		return false, err
	}
	if err := b.Deserialize(sb); err != nil {
		return false, err
	}
	diff, err := a.DifferingFields(b)
	if err != nil {
		return false, err
	}
	return a.MutableFields().Contains(diff), nil
}

func isValidCassandraConfigUpdate(sa, sb pc.SerializedConfig) (bool, error) {
	a := pc.CassandraConfig{}
	b := pc.CassandraConfig{}
	if err := a.Deserialize(sa); err != nil {
		return false, err
	}
	if err := b.Deserialize(sb); err != nil {
		return false, err
	}
	diff, err := a.DifferingFields(b)
	if err != nil {
		return false, err
	}
	return a.MutableFields().Contains(diff), nil
}

func isValidDynamoConfigUpdate(sa, sb pc.SerializedConfig) (bool, error) {
	a := pc.DynamodbConfig{}
	b := pc.DynamodbConfig{}
	if err := a.Deserialize(sa); err != nil {
		return false, err
	}
	if err := b.Deserialize(sb); err != nil {
		return false, err
	}
	diff, err := a.DifferingFields(b)
	if err != nil {
		return false, err
	}
	return a.MutableFields().Contains(diff), nil
}

func isValidFirestoreConfigUpdate(sa, sb pc.SerializedConfig) (bool, error) {
	a := pc.FirestoreConfig{}
	b := pc.FirestoreConfig{}
	if err := a.Deserialize(sa); err != nil {
		return false, err
	}
	if err := b.Deserialize(sb); err != nil {
		return false, err
	}
	diff, err := a.DifferingFields(b)
	if err != nil {
		return false, err
	}
	return a.MutableFields().Contains(diff), nil
}

func isValidMongoConfigUpdate(sa, sb pc.SerializedConfig) (bool, error) {
	a := pc.MongoDBConfig{}
	b := pc.MongoDBConfig{}
	if err := a.Deserialize(sa); err != nil {
		return false, err
	}
	if err := b.Deserialize(sb); err != nil {
		return false, err
	}
	diff, err := a.DifferingFields(b)
	if err != nil {
		return false, err
	}
	return a.MutableFields().Contains(diff), nil
}

func isValidMySqlConfigUpdate(sa, sb pc.SerializedConfig) (bool, error) {
	a := pc.MySqlConfig{}
	b := pc.MySqlConfig{}
	if err := a.Deserialize(sa); err != nil {
		return false, err
	}
	if err := b.Deserialize(sb); err != nil {
		return false, err
	}
	diff, err := a.DifferingFields(b)
	if err != nil {
		return false, err
	}
	return a.MutableFields().Contains(diff), nil
}

func isValidPostgresConfigUpdate(sa, sb pc.SerializedConfig) (bool, error) {
	a := pc.PostgresConfig{}
	b := pc.PostgresConfig{}
	if err := a.Deserialize(sa); err != nil {
		return false, err
	}
	if err := b.Deserialize(sb); err != nil {
		return false, err
	}
	diff, err := a.DifferingFields(b)
	if err != nil {
		return false, err
	}
	return a.MutableFields().Contains(diff), nil
}

func isValidClickHouseConfigUpdate(sa, sb pc.SerializedConfig) (bool, error) {
	a := pc.ClickHouseConfig{}
	b := pc.ClickHouseConfig{}
	if err := a.Deserialize(sa); err != nil {
		return false, err
	}
	if err := b.Deserialize(sb); err != nil {
		return false, err
	}
	diff, err := a.DifferingFields(b)
	if err != nil {
		return false, err
	}
	return a.MutableFields().Contains(diff), nil
}

func isValidRedisConfigUpdate(sa, sb pc.SerializedConfig) (bool, error) {
	a := pc.RedisConfig{}
	b := pc.RedisConfig{}
	if err := a.Deserialize(sa); err != nil {
		return false, err
	}
	if err := b.Deserialize(sb); err != nil {
		return false, err
	}
	diff, err := a.DifferingFields(b)
	if err != nil {
		return false, err
	}
	return a.MutableFields().Contains(diff), nil
}

func isValidSnowflakeConfigUpdate(sa, sb pc.SerializedConfig) (bool, error) {
	a := pc.SnowflakeConfig{}
	b := pc.SnowflakeConfig{}
	if err := a.Deserialize(sa); err != nil {
		return false, err
	}
	if err := b.Deserialize(sb); err != nil {
		return false, err
	}
	diff, err := a.DifferingFields(b)
	if err != nil {
		return false, err
	}
	return a.MutableFields().Contains(diff), nil
}

func isValidRedshiftConfigUpdate(sa, sb pc.SerializedConfig) (bool, error) {
	a := pc.RedshiftConfig{}
	b := pc.RedshiftConfig{}
	if err := a.Deserialize(sa); err != nil {
		return false, err
	}
	if err := b.Deserialize(sb); err != nil {
		return false, err
	}
	diff, err := a.DifferingFields(b)
	if err != nil {
		return false, err
	}
	return a.MutableFields().Contains(diff), nil
}

func isValidK8sConfigUpdate(sa, sb pc.SerializedConfig) (bool, error) {
	a := pc.K8sConfig{}
	b := pc.K8sConfig{}
	if err := a.Deserialize(sa); err != nil {
		return false, err
	}
	if err := b.Deserialize(sb); err != nil {
		return false, err
	}
	diff, err := a.DifferingFields(b)
	if err != nil {
		return false, err
	}
	return a.MutableFields().Contains(diff), nil
}

func isValidSparkConfigUpdate(sa, sb pc.SerializedConfig) (bool, error) {
	a := pc.SparkConfig{}
	b := pc.SparkConfig{}
	if err := a.Deserialize(sa); err != nil {
		return false, err
	}
	if err := b.Deserialize(sb); err != nil {
		return false, err
	}
	diff, err := a.DifferingFields(b)
	if err != nil {
		return false, err
	}
	return a.MutableFields().Contains(diff), nil
}
