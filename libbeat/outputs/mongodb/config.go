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
	"strconv"
	"time"
)

type mongodbConfig struct {
	Host            string  `config:"host"`
	Port            int     `config:"port"`
	Username        string  `config:"username"`
	Password        string  `config:"password"`
	DatabasePrefix  string  `config:"database_prefix"`
	BulkMaxSize     int     `config:"bulk_max_size"`
	MaxRetries      int     `config:"max_retries"`
	ConnectionCount int     `config:"conn_count"`
	MessageKey 		string  `config:"message_key"`
	Backoff         Backoff  `config:"backoff"`
}

const (
	defaultBulkSize = 50
)

type Backoff struct {
	Init time.Duration
	Max  time.Duration
}

var (
	defaultConfig = mongodbConfig{
		Host:           "localhost",
		Port:           27017,
		Username:       "",
		Password:       "",
		DatabasePrefix: "filebeat",
		BulkMaxSize:    defaultBulkSize,
		MaxRetries:     3,
		ConnectionCount: 10,
		Backoff: Backoff{
			Init: 1 * time.Second,
			Max:  60 * time.Second,
		},
	}
)

func (c *mongodbConfig) Validate() error {
	if c.Username != "" && c.Password == "" {
		return fmt.Errorf("need a password for username")
	}
	return nil
}


func (c *mongodbConfig) composeURI() string {
	if c.Username != "" {
		return "mongodb://" + c.Username + ":" + c.Password + "@" + c.Host + ":" + strconv.Itoa(c.Port)
	} else {
		return "mongodb://"+ c.Host + ":" + strconv.Itoa(c.Port)
	}
}
