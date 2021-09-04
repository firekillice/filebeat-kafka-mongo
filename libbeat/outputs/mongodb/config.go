// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package mongodb

import (
	"fmt"
	"time"
)

type mongodbConfig struct {
	Host           string  `config:"host"`
	Port           int     `config:"port"`
	Username       string  `config:"username"`
	Password       string  `config:"password"`
	DatabasePrefix string  `config:"database_prefix"`
	BulkMaxSize    int     `config:"bulk_max_size"`
	MaxRetries     int     `config:"max_retries"`
	Backoff        Backoff `config:"backoff"`
}

type Backoff struct {
	Init time.Duration
	Max  time.Duration
}

const (
	defaultBulkSize = 50
)

var (
	defaultConfig = mongodbConfig{
		Host:           "localhost",
		Port:           27017,
		Username:       "",
		Password:       "",
		DatabasePrefix: "filebeat",
		BulkMaxSize:    defaultBulkSize,
		MaxRetries:     3,
		Backoff: Backoff{
			Init: 1 * time.Second,
			Max:  60 * time.Second,
		},
	}
)

func (c *mongodbConfig) Validate() error {
	if c.APIKey != "" && (c.Username != "" || c.Password != "") {
		return fmt.Errorf("cannot set both api_key and username/password")
	}

	return nil
}
