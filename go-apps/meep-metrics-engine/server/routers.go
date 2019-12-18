/*
 * Copyright (c) 2019  InterDigital Communications, Inc
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
 * AdvantEDGE Metrics Service REST API
 *
 * Metrics Service provides metrics about the active scenario <p>**Micro-service**<br>[meep-metrics-engine](https://github.com/InterDigitalInc/AdvantEDGE/tree/master/go-apps/meep-metrics-engine) <p>**Type & Usage**<br>Platform Service used by control/monitoring software and possibly by edge applications that require metrics <p>**Details**<br>API details available at _your-AdvantEDGE-ip-address:30000/api_ <p>**Default Port**<br>`30005`
 *
 * API version: 1.0.0
 * Contact: AdvantEDGE@InterDigital.com
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package server

import (
	"fmt"
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

	return router
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World! on v1")
}

func IndexV2(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World! on v2")
}

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/v1/",
		Index,
	},

	Route{
		"MetricsGet",
		strings.ToUpper("Get"),
		"/v1/metrics",
		MetricsGet,
	},

	Route{
		"IndexV2",
		"GET",
		"/v2/",
		Index,
	},

	Route{
		"PostEventQuery",
		strings.ToUpper("Post"),
		"/v2/metrics/query/event",
		PostEventQuery,
	},

	Route{
		"PostNetworkQuery",
		strings.ToUpper("Post"),
		"/v2/metrics/query/network",
		PostNetworkQuery,
	},

	Route{
		"CreateEventSubscription",
		strings.ToUpper("Post"),
		"/v2/metrics/subscriptions/event",
		CreateEventSubscription,
	},

	Route{
		"CreateNetworkSubscription",
		strings.ToUpper("Post"),
		"/v2/metrics/subscriptions/network",
		CreateNetworkSubscription,
	},

	Route{
		"DeleteEventSubscriptionById",
		strings.ToUpper("Delete"),
		"/v2/metrics/subscriptions/event/{subscriptionId}",
		DeleteEventSubscriptionById,
	},

	Route{
		"DeleteNetworkSubscriptionById",
		strings.ToUpper("Delete"),
		"/v2/metrics/subscriptions/network/{subscriptionId}",
		DeleteNetworkSubscriptionById,
	},

	Route{
		"GetEventSubscription",
		strings.ToUpper("Get"),
		"/v2/metrics/subscriptions/event",
		GetEventSubscription,
	},

	Route{
		"GetEventSubscriptionById",
		strings.ToUpper("Get"),
		"/v2/metrics/subscriptions/event/{subscriptionId}",
		GetEventSubscriptionById,
	},

	Route{
		"GetNetworkSubscription",
		strings.ToUpper("Get"),
		"/v2/metrics/subscriptions/network",
		GetNetworkSubscription,
	},

	Route{
		"GetNetworkSubscriptionById",
		strings.ToUpper("Get"),
		"/v2/metrics/subscriptions/network/{subscriptionId}",
		GetNetworkSubscriptionById,
	},
}
