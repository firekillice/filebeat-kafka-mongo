// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

// Code generated by beats/dev-tools/cmd/asset/asset.go - DO NOT EDIT.

package suricata

import (
	"github.com/elastic/beats/v7/libbeat/asset"
)

func init() {
	if err := asset.SetFields("filebeat", "suricata", asset.ModuleFieldsPri, AssetSuricata); err != nil {
		panic(err)
	}
}

// AssetSuricata returns asset data.
// This is the base64 encoded gzipped contents of module/suricata.
func AssetSuricata() string {
	return "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"
}
