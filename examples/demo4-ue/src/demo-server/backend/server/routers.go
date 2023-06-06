/*
 * Copyright (c) 2022  The AdvantEDGE Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * MEEP Demo 4 API
 * Demo 4 is an edge application that can be used with AdvantEDGE or ETSI MEC Sandbox to demonstrate MEC016 usage
 *
 * No description provided (generated by Swagger Codegen https://github.com/swagger-api/swagger-codegen)
 *
 * API version: 0.0.1
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package server

import (
	"net/http"
	"strings"

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
		handler = met.MetricsHandler(handler, "", moduleName)
		handler = httpLog.LogRx(handler)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	// Path prefix router order is important
	// Service Api files
	handler = http.StripPrefix("/demo4-ue/api/", http.FileServer(http.Dir("./api/")))
	router.
		PathPrefix("/demo4-ue/api/").
		Name("Api").
		Handler(handler)
	// User supplied service API files
	handler = http.StripPrefix("/demo4-ue/user-api/", http.FileServer(http.Dir("./user-api/")))
	router.
		PathPrefix("/demo4-ue/user-api/").
		Name("UserApi").
		Handler(handler)

	return router
}

var routes = Routes{

	Route{
		"GetActivityLogs",
		strings.ToUpper("Get"),
		"/info/logs",
		GetActivityLogs,
	},

	Route{
		"GetPlatformInfo",
		strings.ToUpper("Get"),
		"/info/application",
		GetPlatformInfo,
	},

	Route{
		"DaiAppListGET",
		strings.ToUpper("Get"),
		"/dai/apps",
		DaiAppListGET,
	},

	Route{
		"DaiDoPingDELETE",
		strings.ToUpper("Delete"),
		"/dai/delete/{appcontextid}",
		DaiDoPingDELETE,
	},

	Route{
		"DaiDoPingGET",
		strings.ToUpper("Get"),
		"/dai/doping/{appcontextid}",
		DaiDoPingGET,
	},

	Route{
		"DaiDoPingPOST",
		strings.ToUpper("Post"),
		"/dai/instantiate",
		DaiDoPingPOST,
	},

	Route{
		"DaiAppLocationAvailabilityPOST",
		strings.ToUpper("Post"),
		"/dai/availability/{appcontextid}",
		DaiAppLocationAvailabilityPOST,
	},

	Route{
		"ServiceAvailNotificationCallback",
		strings.ToUpper("Post"),
		"/services/callback/service-availability",
		ServiceAvailNotificationCallback,
	},

	Route{
		"AppTerminationNotificationCallback",
		strings.ToUpper("Post"),
		"/application/termination",
		AppTerminationNotificationCallback,
	},

	Route{
		"ApplicationContextDeleteNotificationCallback",
		strings.ToUpper("Post"),
		"/dai/callback/ApplicationContextDeleteNotification",
		ApplicationContextDeleteNotificationCallback,
	},
}
