/*
 * Copyright (c) 2022  InterDigital Communications, Inc
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
	"fmt"
	"net/http"
	"strings"

	bwm "github.com/InterDigitalInc/AdvantEDGE/go-apps/meep-tm/server/bwm"
	mts "github.com/InterDigitalInc/AdvantEDGE/go-apps/meep-tm/server/mts"
	httpLog "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-http-logger"
	met "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-metrics"

	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {
	var handler http.Handler
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)
		handler = met.MetricsHandler(handler, sandboxName, serviceName)
		handler = httpLog.LogRx(handler)
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	// Path prefix router order is important
	// Service Api files
	handler = http.StripPrefix("/bwm/v1/api/", http.FileServer(http.Dir("./api/")))
	router.
		PathPrefix("/bwm/v1/api/").
		Name("Api").
		Handler(handler)
	// User supplied service API files
	handler = http.StripPrefix("/bwm/v1/user-api/", http.FileServer(http.Dir("./user-api/")))
	router.
		PathPrefix("/bwm/v1/user-api/").
		Name("UserApi").
		Handler(handler)

	return router
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/bwm/v1/",
		Index,
	},

	Route{
		"Mec011AppTerminationPOST",
		strings.ToUpper("Post"),
		"/bwm/v1/notifications/mec011/appTermination",
		bwm.Mec011AppTerminationPOST,
	},

	Route{
		"BandwidthAllocationDELETE",
		strings.ToUpper("Delete"),
		"/bwm/v1/bw_allocations/{allocationId}",
		bwm.BandwidthAllocationDELETE,
	},

	Route{
		"BandwidthAllocationGET",
		strings.ToUpper("Get"),
		"/bwm/v1/bw_allocations/{allocationId}",
		bwm.BandwidthAllocationGET,
	},

	Route{
		"BandwidthAllocationListGET",
		strings.ToUpper("Get"),
		"/bwm/v1/bw_allocations",
		bwm.BandwidthAllocationListGET,
	},

	Route{
		"BandwidthAllocationPATCH",
		strings.ToUpper("Patch"),
		"/bwm/v1/bw_allocations/{allocationId}",
		bwm.BandwidthAllocationPATCH,
	},

	Route{
		"BandwidthAllocationPOST",
		strings.ToUpper("Post"),
		"/bwm/v1/bw_allocations",
		bwm.BandwidthAllocationPOST,
	},

	Route{
		"BandwidthAllocationPUT",
		strings.ToUpper("Put"),
		"/bwm/v1/bw_allocations/{allocationId}",
		bwm.BandwidthAllocationPUT,
	},

	Route{
		"Index",
		"GET",
		"/mts/v1/",
		Index,
	},

	Route{
		"Mec011AppTerminationPOST",
		strings.ToUpper("Post"),
		"/mts/v1/notifications/mec011/appTermination",
		mts.Mec011AppTerminationPOST,
	},

	Route{
		"MtsCapabilityInfoGET",
		strings.ToUpper("Get"),
		"/mts/v1/mts_capability_info",
		mts.MtsCapabilityInfoGET,
	},

	Route{
		"MtsSessionDELETE",
		strings.ToUpper("Delete"),
		"/mts/v1/mts_sessions/{sessionId}",
		mts.MtsSessionDELETE,
	},

	Route{
		"MtsSessionGET",
		strings.ToUpper("Get"),
		"/mts/v1/mts_sessions/{sessionId}",
		mts.MtsSessionGET,
	},

	Route{
		"MtsSessionPOST",
		strings.ToUpper("Post"),
		"/mts/v1/mts_sessions",
		mts.MtsSessionPOST,
	},

	Route{
		"MtsSessionPUT",
		strings.ToUpper("Put"),
		"/mts/v1/mts_sessions/{sessionId}",
		mts.MtsSessionPUT,
	},

	Route{
		"MtsSessionsListGET",
		strings.ToUpper("Get"),
		"/mts/v1/mts_sessions",
		mts.MtsSessionsListGET,
	},
}
