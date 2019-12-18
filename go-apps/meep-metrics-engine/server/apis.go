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
	"net/http"

	v1 "github.com/InterDigitalInc/AdvantEDGE/go-apps/meep-metrics-engine/server/v1"
	v2 "github.com/InterDigitalInc/AdvantEDGE/go-apps/meep-metrics-engine/server/v2"
)

// Init - Metrics engine initialization
func Init() (err error) {
	//no v1.Init()
	err = v2.Init()
	if err != nil {
		return err
	}
	return nil
}

func MetricsGet(w http.ResponseWriter, r *http.Request) {
	v1.MetricsGet(w, r)
}

func PostEventQuery(w http.ResponseWriter, r *http.Request) {
	v2.PostEventQuery(w, r)
}

func PostNetworkQuery(w http.ResponseWriter, r *http.Request) {
	v2.PostNetworkQuery(w, r)
}

func CreateEventSubscription(w http.ResponseWriter, r *http.Request) {
	v2.CreateEventSubscription(w, r)
}

func CreateNetworkSubscription(w http.ResponseWriter, r *http.Request) {
	v2.CreateNetworkSubscription(w, r)
}

func DeleteEventSubscriptionById(w http.ResponseWriter, r *http.Request) {
	v2.DeleteEventSubscriptionById(w, r)
}

func DeleteNetworkSubscriptionById(w http.ResponseWriter, r *http.Request) {
	v2.DeleteNetworkSubscriptionById(w, r)
}

func GetEventSubscription(w http.ResponseWriter, r *http.Request) {
	v2.GetEventSubscription(w, r)
}

func GetEventSubscriptionById(w http.ResponseWriter, r *http.Request) {
	v2.GetEventSubscriptionById(w, r)
}

func GetNetworkSubscription(w http.ResponseWriter, r *http.Request) {
	v2.GetNetworkSubscription(w, r)
}

func GetNetworkSubscriptionById(w http.ResponseWriter, r *http.Request) {
	v2.GetNetworkSubscriptionById(w, r)
}
