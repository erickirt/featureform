// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// Copyright 2024 FeatureForm Inc.
//

package provider

import (
	"os"
	"testing"

	pc "github.com/featureform/provider/provider_config"
	pt "github.com/featureform/provider/provider_type"
	"github.com/joho/godotenv"
)

func TestOnlineStoreCassandra(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration tests")
	}
	err := godotenv.Load("../.env")
	if err != nil {
		t.Logf("could not open .env file... Checking environment: %s", err)
	}
	cassandraUsername, ok := os.LookupEnv("CASSANDRA_USER")
	if !ok {
		t.Fatalf("missing CASSANDRA_USER variable")
	}
	cassandraPassword, ok := os.LookupEnv("CASSANDRA_PASSWORD")
	if !ok {
		t.Fatalf("missing CASSANDRA_PASSWORD variable")
	}
	cassandraAddr := "localhost:9042"
	cassandraConfig := &pc.CassandraConfig{
		Addr:        cassandraAddr,
		Username:    cassandraUsername,
		Consistency: "ONE",
		Password:    cassandraPassword,
		Replication: 3,
	}

	store, err := GetOnlineStore(pt.CassandraOnline, cassandraConfig.Serialized())
	if err != nil {
		t.Fatalf("could not initialize store: %s\n", err)
	}

	test := OnlineStoreTest{
		t:     t,
		store: store,
	}
	test.Run()
}
