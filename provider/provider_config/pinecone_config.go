// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// Copyright 2024 FeatureForm Inc.
//

package provider_config

import (
	"encoding/json"

	"github.com/featureform/fferr"

	ss "github.com/featureform/helpers/stringset"
)

// https://docs.pinecone.io/docs/projects
type PineconeConfig struct {
	// NOTE: it appears the only place to fetch the project ID is
	// from the URL of the project page in the Pinecone dashboard.
	// For example:
	// https://app.pinecone.io/organizations/<ORG ID>/projects/us-west4-gcp-free:<PROJECT ID>/indexes
	ProjectID   string
	Environment string
	ApiKey      string
}

func (pc PineconeConfig) Serialize() SerializedConfig {
	config, err := json.Marshal(pc)
	if err != nil {
		panic(err)
	}
	return config
}

func (pc *PineconeConfig) Deserialize(config SerializedConfig) error {
	err := json.Unmarshal(config, pc)
	if err != nil {
		return fferr.NewInternalError(err)
	}
	return nil
}

func (pc PineconeConfig) MutableFields() ss.StringSet {
	return ss.StringSet{
		// Based on this note in the documentation:
		// "The environment cannot be changed after the project is created."
		// - https://docs.pinecone.io/docs/projects
		// it seems that only the API key should be mutable.
		"ApiKey": true,
	}
}

func (a PineconeConfig) DifferingFields(b PineconeConfig) (ss.StringSet, error) {
	return differingFields(a, b)
}
