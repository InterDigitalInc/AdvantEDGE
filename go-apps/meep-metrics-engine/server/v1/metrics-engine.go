/*
 * Copyright (c) 2019  InterDigital Communications, Inc
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"

	"github.com/olivere/elastic"
)

type ElasticFormatedLogResponse struct {
	Msg       string `json:"msg"`
	MsgType   string `json:"meep.log.msgType"`
	Src       string `json:"meep.log.src"`
	Dest      string `json:"meep.log.dest"`
	Timestamp string `json:"@timestamp"`

	/*** specific fields for all message types

	/*** ingressPacketStats ***/
	Rx         int32   `json:"meep.log.rx"`
	RxBytes    int32   `json:"meep.log.rxBytes"`
	PacketLoss string  `json:"meep.log.packet-loss"`
	Throughput float32 `json:"meep.log.throughput"`

	/*** latency ***/
	Latency int32 `json:"meep.log.latency-latest"`

	/*** mobilityEvent ***/
	NewPoa string `json:"meep.log.newPoa"`
	OldPoa string `json:"meep.log.oldPoa"`
}

func metricsGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	client, err := elastic.NewClient(elastic.SetURL("http://meep-elasticsearch-client:9200"))

	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Search with a term query
	bq := elastic.NewBoolQuery()
	bq = bq.Must(elastic.NewTermQuery("msg", "Measurements log"))

	u, _ := url.Parse(r.URL.String())
	q := u.Query()

	msgType := q.Get("dataType")
	if msgType != "" {
		bq = bq.Must(elastic.NewTermQuery("meep.log.msgType", msgType))
	}

	dst := q.Get("dest")
	if dst != "" {
		bq = bq.Must(elastic.NewTermQuery("meep.log.dest", dst))
	}

	src := q.Get("src")
	if src != "" {
		bq = bq.Must(elastic.NewTermQuery("meep.log.src", src))
	}

	timeBegin := q.Get("startTime")
	timeEnd := q.Get("stopTime")

	//default values
	if timeBegin == "" {
		timeBegin = "now-1m"
	}
	if timeEnd == "" {
		timeEnd = "now"
	}
	bq = bq.Must(elastic.NewRangeQuery("@timestamp").Gte(timeBegin).Lte(timeEnd))

	log.Info("Search query: ", "Measurements log", " + ", msgType, " + ", dst, " + ", src, " + ", timeBegin, " + ", timeEnd)

	searchQuery := client.Scroll("filebeat*").
		Query(bq). // specify the query
		Size(1000) // take documents 0-9
		//		Pretty(true) // pretty print request and response JSON

	docs := 0
	pages := 0
	print := 0
	var logResponseList LogResponseList
	for {
		res, err := searchQuery.Do(context.Background())
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Info("Error while querying ES: ", err)
			break
		}
		if res == nil {
			log.Info("Null result from ES")
			break
		}
		if res.Hits == nil {
			log.Info("Not even a single hit in ES")
			break
		}

		pages++

		for _, hit := range res.Hits.Hits {
			//item := make(map[string]interface{})
			var t ElasticFormatedLogResponse
			err := json.Unmarshal(*hit.Source, &t)
			if err != nil {
				log.Info("Deserialization failed")
				//                                continue
			}
			logResponse := convertToLogResponse(&t)
			logResponseList.LogResponse = append(logResponseList.LogResponse, *logResponse)
			print++
			docs++
		}
	}
	log.Info("Total number of results: ", docs, " in ", pages, " different queries")
	if docs > 0 {
		jsonResponse, err := json.Marshal(logResponseList)

		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, string(jsonResponse))
	}
	w.WriteHeader(http.StatusOK)
}
