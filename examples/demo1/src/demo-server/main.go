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

package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	sw "github.com/InterDigitalInc/AdvantEDGE/demoserver/go"
	"github.com/gorilla/handlers"

	locServClient "github.com/InterDigitalInc/AdvantEDGE/locservapi"
)

func init() {
	// Initialize App
	sw.Init()
}

func main() {
	log.Printf("DemoSvc App API Server started")

	run := true

	go func() {
		sigchan := make(chan os.Signal, 10)
		signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
		<-sigchan
		log.Printf("Program killed !")
		// do last actions and wait for all write operations to end
		run = false
	}()

	go func() {
		router := sw.NewRouter()

		methods := handlers.AllowedMethods([]string{"OPTIONS", "DELETE", "GET", "HEAD", "POST", "PUT"})
		header := handlers.AllowedHeaders([]string{"content-type"})

		registerLocServ("ue2-ext")
		registerLocServ("ue1")

		log.Fatal(http.ListenAndServe(":80", handlers.CORS(methods, header)(router)))

		run = false
	}()

	count := 0
	for {
		if !run {
			log.Printf("Ran for %d seconds", count)
			break
		}
		time.Sleep(time.Second)
		count++
	}
}

func registerLocServ(ue string) {
	locServCfg := locServClient.NewConfiguration()
	locServCfg.BasePath = "http://meep-loc-serv/location/v2"

	locServ := locServClient.NewAPIClient(locServCfg)
	log.Printf("Created Location Service client before")

	if locServ == nil {
		log.Printf("Cannot find the Location Service API")
		return
	}
	log.Printf("Created Location Service client")

	var subscription locServClient.UserTrackingSubscription
	subscription.ClientCorrelator = "001" //don't care
	subscription.Address = ue
	var userCriteria []locServClient.UserEventType
	userCriteria = append(userCriteria, "Entering")
	userCriteria = append(userCriteria, "Transferring")
	subscription.UserEventCriteria = userCriteria

	serviceName := os.Getenv("MGM_APP_ID")
	newString := strings.ToUpper(serviceName) + "_SERVICE_HOST"
	newString = strings.Replace(newString, "-", "_", -1)

	myPodIp := os.Getenv(newString)
	var cb locServClient.CallbackReference
	//using the current server implementation, but with new loc-serv api which returns exaclty the notifyURL
	cb.NotifyURL = "http://" + myPodIp + "/v1/location_notifications/1"
	subscription.CallbackReference = &cb

	var subscriptionBody locServClient.InlineUserTrackingSubscription
	subscriptionBody.UserTrackingSubscription = &subscription
	_, resp, err := locServ.LocationApi.UserTrackingSubPOST(nil, subscriptionBody)
	if err != nil {
		log.Printf(err.Error())
	}
	defer resp.Body.Close()
	if resp != nil {
		if resp.StatusCode != http.StatusOK {
			log.Printf("Not OK status response")
		}
		responseData, err := ioutil.ReadAll(resp.Body)

		if err == nil {
			responseString := string(responseData)
			log.Printf(responseString)
		} else {
			log.Printf("response decoding error, %s", err)
		}
	}
}
