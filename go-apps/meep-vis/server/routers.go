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
 * AdvantEDGE V2X Information Service REST API
 *
 * V2X Information Service is AdvantEDGE's implementation of [ETSI MEC ISG MEC030 V2XI API](http://www.etsi.org/deliver/etsi_gs/MEC/001_099/030/02.02.01_60/gs_MEC030v020201p.pdf) <p>[Copyright (c) ETSI 2017](https://forge.etsi.org/etsi-forge-copyright-notice.txt) <p>**Micro-service**<br>[meep-vis](https://github.com/InterDigitalInc/AdvantEDGE/tree/master/go-apps/meep-vis) <p>**Type & Usage**<br>Edge Service used by edge applications that want to get information about radio conditions in the network <p>**Note**<br>AdvantEDGE supports a selected subset of RNI API endpoints (see below) and a subset of subscription types.
 *
 * API version: 2.2.1
 * Contact: AdvantEDGE@InterDigital.com
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package server

import (
	"fmt"
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
	handler = http.StripPrefix("/vis/v2/api/", http.FileServer(http.Dir("./api/")))
	router.
		PathPrefix("/vis/v2/api/").
		Name("Api").
		Handler(handler)
	// User supplied service API files
	handler = http.StripPrefix("/vis/v2/user-api/", http.FileServer(http.Dir("./user-api/")))
	router.
		PathPrefix("/vis/v2/user-api/").
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
		"/vis/v2/",
		Index,
	},

	Route{
		"IndividualSubscriptionDELETE",
		strings.ToUpper("Delete"),
		"/vis/v2/subscriptions/{subscriptionId}",
		IndividualSubscriptionDELETE,
	},

	Route{
		"IndividualSubscriptionGET",
		strings.ToUpper("Get"),
		"/vis/v2/subscriptions/{subscriptionId}",
		IndividualSubscriptionGET,
	},

	Route{
		"IndividualSubscriptionPUT",
		strings.ToUpper("Put"),
		"/vis/v2/subscriptions/{subscriptionId}",
		IndividualSubscriptionPUT,
	},

	Route{
		"ProvInfoGET",
		strings.ToUpper("Get"),
		"/vis/v2/queries/pc5_provisioning_info",
		ProvInfoGET,
	},

	Route{
		"ProvInfoUuMbmsGET",
		strings.ToUpper("Get"),
		"/vis/v2/queries/uu_mbms_provisioning_info",
		ProvInfoUuMbmsGET,
	},

	Route{
		"ProvInfoUuUnicastGET",
		strings.ToUpper("Get"),
		"/vis/v2/queries/uu_unicast_provisioning_info",
		ProvInfoUuUnicastGET,
	},

	Route{
		"SubGET",
		strings.ToUpper("Get"),
		"/vis/v2/subscriptions",
		SubGET,
	},

	Route{
		"SubPOST",
		strings.ToUpper("Post"),
		"/vis/v2/subscriptions",
		SubPOST,
	},

	Route{
		"V2xMessagePOST",
		strings.ToUpper("Post"),
		"/vis/v2/publish_v2x_message",
		V2xMessagePOST,
	},

	Route{
		"Mec011AppTerminationPOST",
		strings.ToUpper("Post"),
		"/vis/v2/notifications/mec011/appTermination",
		Mec011AppTerminationPOST,
	},

	Route{
		"PredictedQosPOST",
		strings.ToUpper("Post"),
		"/vis/v2/provide_predicted_qos",
		PredictedQosPOST,
	},
}
