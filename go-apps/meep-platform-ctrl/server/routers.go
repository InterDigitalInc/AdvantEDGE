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
 * AdvantEDGE Platform Controller REST API
 *
 * This API is the main Platform Controller API for scenario configuration & sandbox management <p>**Micro-service**<br>[meep-pfm-ctrl](https://github.com/InterDigitalInc/AdvantEDGE/tree/master/go-apps/meep-platform-ctrl) <p>**Type & Usage**<br>Platform main interface used by controller software to configure scenarios and manage sandboxes in the AdvantEDGE platform <p>**Details**<br>API details available at _your-AdvantEDGE-ip-address/api_
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

func NewRouter(priFe string, priSw string, altFe string, altSw string) *mux.Router {
	var handler http.Handler
	router := mux.NewRouter().StrictSlash(true)

	for _, route := range routes {
		handler = Logger(route.HandlerFunc, route.Name)
		handler = met.MetricsHandler(handler, "", serviceName)
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	// Path prefix router order is important
	// Service Api files
	handler = http.StripPrefix("/platform-ctrl/v1/api/", http.FileServer(http.Dir("./api/")))
	router.
		PathPrefix("/platform-ctrl/v1/api/").
		Name("Api").
		Handler(handler)
	// User supplied service API files
	handler = http.StripPrefix("/platform-ctrl/v1/user-api/", http.FileServer(http.Dir("./user-api/")))
	router.
		PathPrefix("/platform-ctrl/v1/user-api/").
		Name("UserApi").
		Handler(handler)
	// Alternate Swagger UI
	if altSw != "" {
		handler = http.StripPrefix("/alt/api/", http.FileServer(http.Dir(altSw)))
		router.
			PathPrefix("/alt/api/").
			Name("AltSw").
			Handler(handler)
	}
	// Primary Swagger UI
	if priSw != "" {
		handler = http.StripPrefix("/api/", http.FileServer(http.Dir(priSw)))
		router.
			PathPrefix("/api/").
			Name("PriSw").
			Handler(handler)
	}
	// Alternate Frontend
	if altFe != "" {
		handler = http.StripPrefix("/alt/", http.FileServer(http.Dir(altFe)))
		router.
			PathPrefix("/alt/").
			Name("AltFe").
			Handler(handler)
	}
	// Primary Frontend
	if priFe != "" {
		handler = http.StripPrefix("/", http.FileServer(http.Dir(priFe)))
		router.
			PathPrefix("/").
			Name("PriFe").
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
		"/platform-ctrl/v1/",
		Index,
	},

	Route{
		"CreateSandbox",
		strings.ToUpper("Post"),
		"/platform-ctrl/v1/sandboxes",
		CreateSandbox,
	},

	Route{
		"CreateSandboxWithName",
		strings.ToUpper("Post"),
		"/platform-ctrl/v1/sandboxes/{name}",
		CreateSandboxWithName,
	},

	Route{
		"DeleteSandbox",
		strings.ToUpper("Delete"),
		"/platform-ctrl/v1/sandboxes/{name}",
		DeleteSandbox,
	},

	Route{
		"DeleteSandboxList",
		strings.ToUpper("Delete"),
		"/platform-ctrl/v1/sandboxes",
		DeleteSandboxList,
	},

	Route{
		"GetSandbox",
		strings.ToUpper("Get"),
		"/platform-ctrl/v1/sandboxes/{name}",
		GetSandbox,
	},

	Route{
		"GetSandboxList",
		strings.ToUpper("Get"),
		"/platform-ctrl/v1/sandboxes",
		GetSandboxList,
	},

	Route{
		"CreateScenario",
		strings.ToUpper("Post"),
		"/platform-ctrl/v1/scenarios/{name}",
		CreateScenario,
	},

	Route{
		"DeleteScenario",
		strings.ToUpper("Delete"),
		"/platform-ctrl/v1/scenarios/{name}",
		DeleteScenario,
	},

	Route{
		"DeleteScenarioList",
		strings.ToUpper("Delete"),
		"/platform-ctrl/v1/scenarios",
		DeleteScenarioList,
	},

	Route{
		"GetScenario",
		strings.ToUpper("Get"),
		"/platform-ctrl/v1/scenarios/{name}",
		GetScenario,
	},

	Route{
		"GetScenarioList",
		strings.ToUpper("Get"),
		"/platform-ctrl/v1/scenarios",
		GetScenarioList,
	},

	Route{
		"SetScenario",
		strings.ToUpper("Put"),
		"/platform-ctrl/v1/scenarios/{name}",
		SetScenario,
	},
}
