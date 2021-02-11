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
 * AdvantEDGE WLAN Access Information API
 *
 * WLAN Access Information Service is AdvantEDGE's implementation of [ETSI MEC ISG MEC028 WAI API](http://www.etsi.org/deliver/etsi_gs/MEC/001_099/028/02.01.01_60/gs_MEC028v020101p.pdf) <p>[Copyright (c) ETSI 2020](https://forge.etsi.org/etsi-forge-copyright-notice.txt) <p>**Micro-service**<br>[meep-wais](https://github.com/InterDigitalInc/AdvantEDGE/tree/master/go-apps/meep-wais) <p>**Type & Usage**<br>Edge Service used by edge applications that want to get information about WLAN access information in the network <p>**Details**<br>API details available at _your-AdvantEDGE-ip-address/api_ <p>AdvantEDGE supports a selected subset of WAI API subscription types. <p>Supported subscriptions: <p> - AssocStaSubscription
 *
 * API version: 2.1.1
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
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler = route.HandlerFunc
		handler = Logger(handler, route.Name)
		handler = met.MetricsHandler(handler, sandboxName, logModuleWAIS)
		handler = httpLog.LogRx(handler, "")
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/wai/v2/",
		Index,
	},

	Route{
		"ApInfoGET",
		strings.ToUpper("Get"),
		"/wai/v2/queries/ap/ap_information",
		ApInfoGET,
	},

	Route{
		"StaInfoGET",
		strings.ToUpper("Get"),
		"/wai/v2/queries/sta/sta_information",
		StaInfoGET,
	},

	Route{
		"SubscriptionLinkListSubscriptionsGET",
		strings.ToUpper("Get"),
		"/wai/v2/subscriptions",
		SubscriptionLinkListSubscriptionsGET,
	},

	Route{
		"SubscriptionsDELETE",
		strings.ToUpper("Delete"),
		"/wai/v2/subscriptions/{subscriptionId}",
		SubscriptionsDELETE,
	},

	Route{
		"SubscriptionsGET",
		strings.ToUpper("Get"),
		"/wai/v2/subscriptions/{subscriptionId}",
		SubscriptionsGET,
	},

	Route{
		"SubscriptionsPOST",
		strings.ToUpper("Post"),
		"/wai/v2/subscriptions",
		SubscriptionsPOST,
	},

	Route{
		"SubscriptionsPUT",
		strings.ToUpper("Put"),
		"/wai/v2/subscriptions/{subscriptionId}",
		SubscriptionsPUT,
	},
}
