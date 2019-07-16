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
	"net/url"
	"fmt"
	"net/http"
	"io"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	"github.com/gorilla/mux"
	"github.com/olivere/elastic"
)

type LogDataResponse struct {
	Msg       string `json:"msg"`
	MsgType   string `json:"meep.sidecar.msgType"`
	Src       string `json:"meep.sidecar.src"`
	Dest      string `json:"meep.sidecar.dest"`
	Timestamp string `json:"@timestamp"`

	Rx         int32   `json:"meep.sidecar.rx"`
	RxBytes    int32   `json:"meep.sidecar.rxBytes"`
	PacketLoss string  `json:"meep.sidecar.packet-loss"`
	Latency    int32   `json:"meep.sidecar.latency-latest"`
	Throughput float32 `json:"meep.sidecar.throughput"`
}

// Init - Location Service initialization
func Init() (err error) {
	return nil
}

func metricsGet(w http.ResponseWriter, r *http.Request) {
	getMetrics(w, r, "*", "*", "*")
}

func metricsGetByMsgType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	getMetrics(w, r, vars["msgType"], "*", "*")
}

func metricsGetByMsgTypeByDst(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	getMetrics(w, r, vars["msgType"], vars["dst"], "*")
}

func metricsGetByMsgTypeByDstBySrc(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	getMetrics(w, r, vars["msgType"], vars["dst"], vars["src"])
}

func getMetrics(w http.ResponseWriter, r *http.Request, msgType string, dst string, src string) {
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
	if msgType != "*" {
		bq = bq.Must(elastic.NewTermQuery("meep.sidecar.msgType", msgType))
	}
	if dst != "*" {
		bq = bq.Must(elastic.NewTermQuery("meep.sidecar.dest", dst))
	}
	if src != "*" {
		bq = bq.Must(elastic.NewTermQuery("meep.sidecar.src", src))
	}
        u, _ := url.Parse(r.URL.String())
        q := u.Query()
        timeBegin := q.Get("timeBegin")
        timeEnd := q.Get("timeEnd")

	//default values
	if timeBegin == "" {
		timeBegin= "now-1m"
	}
	if timeEnd == "" {
		timeEnd = "now"
	}
	bq = bq.Must(elastic.NewRangeQuery("@timestamp").Gte(timeBegin).Lte(timeEnd))

	log.Info("Search query: ", "Measurements log", " + ", msgType, " + ", dst, " + ", src, " + ", timeBegin, " + ", timeEnd)

	searchQuery := client.Scroll("filebeat*").
		Query(bq).   // specify the query
		Size(1000)    // take documents 0-9
//		Pretty(true) // pretty print request and response JSON

	docs := 0
	pages := 0
	print := 0
	var dataResponseList DataResponseList
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
			var t LogDataResponse
			err := json.Unmarshal(*hit.Source, &t)
			if err != nil {
				log.Info("Deserialization failed")
				//                                continue
			}
			dataResponse := convertToDataResponse(&t)
			dataResponseList.DataResponse = append(dataResponseList.DataResponse, *dataResponse)
			print++
			docs++
		}
	}
	log.Info("Total number of results: ", docs, " in ", pages, " different queries")
	jsonResponse, err := json.Marshal(dataResponseList)

        if err != nil {
                log.Error(err.Error())
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
        }
        w.WriteHeader(http.StatusOK)
        fmt.Fprintf(w, string(jsonResponse))
}

func convertToDataResponse(logDataResponse *LogDataResponse) (*DataResponse) {
	
	if logDataResponse == nil {
		return nil
	}

	msgType := logDataResponse.MsgType

        var resp DataResponse
        resp.DataType = msgType
        resp.Src = logDataResponse.Src
        resp.Dest = logDataResponse.Dest
        resp.Timestamp = logDataResponse.Timestamp

	switch(msgType) {
	case "latency":
	        var data DataResponseData
        	data.Latency = logDataResponse.Latency
		resp.Data = &data
	case "packetStats":
		var data DataResponseData
		data.Rx = logDataResponse.Rx
		data.RxBytes = logDataResponse.RxBytes
		data.Throughput = logDataResponse.Throughput
		data.PacketLoss = logDataResponse.PacketLoss
		resp.Data = &data
	default:
	}
	return &resp
}
