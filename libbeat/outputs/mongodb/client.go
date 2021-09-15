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
	"context"
	"encoding/json"
	"regexp"
	"time"

	"github.com/elastic/beats/v7/libbeat/logp"
	"github.com/elastic/beats/v7/libbeat/outputs"
	"github.com/elastic/beats/v7/libbeat/publisher"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Client struct {
	mc 		 *mongo.Client
	config 	 mongodbConfig
	log      *logp.Logger
	ctx 	 context.Context
	observer outputs.Observer
}

func NewClient(s mongodbConfig, observer outputs.Observer) (*Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log := logp.NewLogger("mongodb")

	uri := s.composeURI()
	mc, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Error("Mongodb fail to new client to %s %+v", uri, err)
	}

	client := &Client{
		mc:     	mc,
		config: 	s,
		ctx: 		ctx,
		log: 		log,
		observer: 	observer,
	}

	return client, nil
}

func (client *Client) Publish(ctx context.Context, batch publisher.Batch) error {
	events := batch.Events()
	rest, err := client.publishEvents(ctx, events)
	if len(rest) == 0 {
		batch.ACK()
	} else {
		batch.RetryEvents(rest)
	}
	return err
}

func (client *Client) publishEvents(ctx context.Context, data []publisher.Event) ([]publisher.Event, error) {
	st := client.observer

	if st != nil {
		st.NewBatch(len(data))
	}

	if len(data) == 0 {
		return nil, nil
	}

	var failedEvents []publisher.Event
	origCount := len(data)
	dropped := 0
	for i := 0; i < origCount; i++ {
		message, err := data[i].Content.GetValue(client.config.MessageKey)
		if err != nil {
			dropped++
			msgBytes, _ := json.Marshal(message)
			client.log.Info("Mongodb output drop event %s", string(msgBytes))
			continue
		}
		client.filterMap(message.(map[string]interface{}))
		collection := client.mc.Database(client.getDbName(message)).Collection(client.getCollectionName(message))

		_, err = collection.InsertOne(ctx, message)
		if err != nil {
			failedEvents = append(failedEvents, data[i])
		}
	}

	failed := len(failedEvents)

	st.Acked(origCount - failed - dropped)
	st.Failed(len(failedEvents))
	st.Dropped(dropped)
	st.Duplicate(0)
	st.ErrTooMany(0)

	if failed > 0 {
		return failedEvents, nil
	}
	return nil, nil
}

func (client *Client) Connect() error {
	return client.mc.Connect(client.ctx)
}

func (client *Client) Close() error {
	return client.mc.Disconnect(client.ctx)
}

func (client *Client) String() string {
	return "mongodb-connection"
}

//从文档信息中获取数据库名称
func (client *Client) getDbName(json interface{}) string {
	for k, v := range json.(map[string]interface{}) {
		if k == "logTime" {
			tm := time.Unix(v.(int64), 0)
			return tm.Format("2006-01-02") + "_Log"
		}
	}
	return "logTimeMissing"
}

//从文档信息中获取集合名称
func (client *Client) getCollectionName(doc interface{}) string {
	formalColName := ""
	for k, v := range doc.(map[string]interface{}) {
		if k == "logType" {
			reg, _ := regexp.Compile("system\\.|/|\\\\|\\$") //system. \ / $
			formalColName = reg.ReplaceAllString(v.(string), "");
		}
	}
	if formalColName == "" {
		return "actionTypeMissing"
	}
	return formalColName
}

func (client *Client)checkIllegalKey(key string, replace bool ) (string, bool) {
	reg, _ := regexp.Compile("^[$_]+|\\.") //system. \ / $
	matched := reg.MatchString(key)
	if replace {
		return reg.ReplaceAllString(key, ""), matched
	} else {
		return "", matched
	}
}

//过滤key中的非法字符
func (client *Client)filterMap(md map[string]interface{}) {
	for key, val := range md {
		keyRep, matched := client.checkIllegalKey(key, true)
		if matched {
			delete(md, key)
			md[keyRep] = val
		}
		switch val.(type) {
		case map[string]interface{}:
			client.filterMap(val.(map[string]interface{}))
		case []interface{}:
			client.filterArray(val.([]interface{}))
		default:

		}
	}
}

func (client *Client) filterArray(ad []interface{}) {
	for _, val := range ad {
		switch val.(type) {
		case map[string]interface{}:
			client.filterMap(val.(map[string]interface{}))
		case []interface{}:
			client.filterArray(val.([]interface{}))
		}
	}
}

