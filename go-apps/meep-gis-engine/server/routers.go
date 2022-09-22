/*
 * Copyright (c) 2022  InterDigital Communications, Inc
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
 * AdvantEDGE GIS Engine REST API
 *
 * This API allows to control geo-spatial behavior and simulation. <p>**Micro-service**<br>[meep-gis-engine](https://github.com/InterDigitalInc/AdvantEDGE/tree/master/go-apps/meep-gis-engine) <p>**Type & Usage**<br>Platform runtime interface to control geo-spatial behavior and simulation <p>**Details**<br>API details available at _your-AdvantEDGE-ip-address/api_
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
		handler = met.MetricsHandler(handler, ge.sandboxName, serviceName)
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	// Path prefix router order is important
	// Service Api files
	handler = http.StripPrefix("/gis/v1/api/", http.FileServer(http.Dir("./api/")))
	router.
		PathPrefix("/gis/v1/api/").
		Name("Api").
		Handler(handler)
	// User supplied service API files
	handler = http.StripPrefix("/gis/v1/user-api/", http.FileServer(http.Dir("./user-api/")))
	router.
		PathPrefix("/gis/v1/user-api/").
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
		"/gis/v1/",
		Index,
	},

	Route{
		"GetAutomationState",
		strings.ToUpper("Get"),
		"/gis/v1/automation",
		GetAutomationState,
	},

	Route{
		"GetAutomationStateByName",
		strings.ToUpper("Get"),
		"/gis/v1/automation/{type}",
		GetAutomationStateByName,
	},

	Route{
		"SetAutomationStateByName",
		strings.ToUpper("Post"),
		"/gis/v1/automation/{type}",
		SetAutomationStateByName,
	},

	Route{
		"DeleteGeoDataByName",
		strings.ToUpper("Delete"),
		"/gis/v1/geodata/{assetName}",
		DeleteGeoDataByName,
	},

	Route{
		"GetAssetData",
		strings.ToUpper("Get"),
		"/gis/v1/geodata",
		GetAssetData,
	},

	Route{
		"GetDistanceGeoDataByName",
		strings.ToUpper("Post"),
		"/gis/v1/geodata/{assetName}/distanceTo",
		GetDistanceGeoDataByName,
	},

	Route{
		"GetGeoDataByName",
		strings.ToUpper("Get"),
		"/gis/v1/geodata/{assetName}",
		GetGeoDataByName,
	},

	Route{
		"GetGeoDataPowerValues",
		strings.ToUpper("Post"),
		"/gis/v1/geodata/cellularPower",
		GetGeoDataPowerValues,
	},

	Route{
		"GetWithinRangeByName",
		strings.ToUpper("Post"),
		"/gis/v1/geodata/{assetName}/withinRange",
		GetWithinRangeByName,
	},

	Route{
		"UpdateGeoDataByName",
		strings.ToUpper("Post"),
		"/gis/v1/geodata/{assetName}",
		UpdateGeoDataByName,
	},
}
