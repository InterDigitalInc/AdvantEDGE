/*
 * Copyright (c) 2020  InterDigital Communications, Inc
 *
 * Licensed under the Apache License, Version 2.0 (the \"License\");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an \"AS IS\" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * MEC Demo 3 API
 *
 * Demo 3 is an edge application that can be used with AdvantEDGE or ETSI MEC Sandbox to demonstrate MEC011 and MEC021 usage
 *
 * API version: 0.0.1
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package server

import (
	"net/http"
	"strings"

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
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./static/"))))

	return router
}

var routes = Routes{

	Route{
		"DeleteAmsDevice",
		strings.ToUpper("Delete"),
		"/service/ams/delete/{device}",
		DeleteAmsDevice,
	},

	Route{
		"Deregister",
		strings.ToUpper("Delete"),
		"/info/application/delete",
		Deregister,
	},

	Route{
		"GetActivityLogs",
		strings.ToUpper("Get"),
		"/info/logs",
		GetActivityLogs,
	},

	Route{
		"GetAmsDevices",
		strings.ToUpper("Get"),
		"/info/ams",
		GetAmsDevices,
	},

	Route{
		"GetPlatformInfo",
		strings.ToUpper("Get"),
		"/info/application",
		GetPlatformInfo,
	},

	Route{
		"Register",
		strings.ToUpper("Post"),
		"/register/app",
		Register,
	},

	Route{
		"UpdateAmsDevices",
		strings.ToUpper("Put"),
		"/service/ams/update/{device}",
		UpdateAmsDevices,
	},

	Route{
		"AmsNotificationCallback",
		strings.ToUpper("Post"),
		"/services/callback/amsevent",
		AmsNotificationCallback,
	},

	Route{
		"AppTerminationNotificationCallback",
		strings.ToUpper("Post"),
		"/application/termination",
		AppTerminationNotificationCallback,
	},

	Route{
		"ContextTransferNotificationCallback",
		strings.ToUpper("Post"),
		"/application/transfer",
		ContextTransferNotificationCallback,
	},

	Route{
		"ServiceAvailNotificationCallback",
		strings.ToUpper("Post"),
		"/services/callback/service-availability",
		ServiceAvailNotificationCallback,
	},
}
