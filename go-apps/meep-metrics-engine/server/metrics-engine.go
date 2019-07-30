/*
 * Copyright (c) 2019
 * InterDigital Communications, Inc.
 * All rights reserved.
 *
 * The information provided herein is the proprietary and confidential
 * information of InterDigital Communications, Inc.
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

// Init - Location Service initialization
func Init() (err error) {
	return nil
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

func convertToLogResponse(esLogResponse *ElasticFormatedLogResponse) *LogResponse {

	if esLogResponse == nil {
		return nil
	}

	msgType := esLogResponse.MsgType

	var resp LogResponse
	resp.DataType = msgType
	resp.Src = esLogResponse.Src
	resp.Dest = esLogResponse.Dest
	resp.Timestamp = esLogResponse.Timestamp

	switch msgType {
	case "latency":
		var data LogResponseData
		data.Latency = esLogResponse.Latency
		resp.Data = &data
	case "ingressPacketStats":
		var data LogResponseData
		data.Rx = esLogResponse.Rx
		data.RxBytes = esLogResponse.RxBytes
		data.Throughput = esLogResponse.Throughput
		data.PacketLoss = esLogResponse.PacketLoss
		resp.Data = &data
	case "mobilityEvent":
		var data LogResponseData
		data.NewPoa = esLogResponse.NewPoa
		data.OldPoa = esLogResponse.OldPoa
		resp.Data = &data
	default:
	}
	return &resp
}
