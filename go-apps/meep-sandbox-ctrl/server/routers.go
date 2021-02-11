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
 * AdvantEDGE Sandbox Controller REST API
 *
 * This API is the main Sandbox Controller API for scenario deployment & event injection <p>**Micro-service**<br>[meep-sandbox-ctrl](https://github.com/InterDigitalInc/AdvantEDGE/tree/master/go-apps/meep-sandbox-ctrl) <p>**Type & Usage**<br>Platform runtime interface to manage active scenarios and inject events in AdvantEDGE platform <p>**Details**<br>API details available at _your-AdvantEDGE-ip-address/api_
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

func NewRouter(priSw string, altSw string) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	for _, route := range routes {
		var handler http.Handler = Logger(route.HandlerFunc, route.Name)
		handler = met.MetricsHandler(handler, sbxCtrl.sandboxName, moduleName)
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	// Path prefix router order is important
	if altSw != "" {
		var handler http.Handler = http.StripPrefix("/alt/api/", http.FileServer(http.Dir(altSw)))
		router.
			PathPrefix("/alt/api/").
			Name("AltSw").
			Handler(handler)
	}
	if priSw != "" {
		var handler http.Handler = http.StripPrefix("/api/", http.FileServer(http.Dir(priSw)))
		router.
			PathPrefix("/api/").
			Name("PriSw").
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
		"/sandbox-ctrl/v1/",
		Index,
	},

	Route{
		"ActivateScenario",
		strings.ToUpper("Post"),
		"/sandbox-ctrl/v1/active/{name}",
		ActivateScenario,
	},

	Route{
		"GetActiveNodeServiceMaps",
		strings.ToUpper("Get"),
		"/sandbox-ctrl/v1/active/serviceMaps",
		GetActiveNodeServiceMaps,
	},

	Route{
		"GetActiveScenario",
		strings.ToUpper("Get"),
		"/sandbox-ctrl/v1/active",
		GetActiveScenario,
	},

	Route{
		"TerminateScenario",
		strings.ToUpper("Delete"),
		"/sandbox-ctrl/v1/active",
		TerminateScenario,
	},

	Route{
		"CreateReplayFile",
		strings.ToUpper("Post"),
		"/sandbox-ctrl/v1/replay/{name}",
		CreateReplayFile,
	},

	Route{
		"CreateReplayFileFromScenarioExec",
		strings.ToUpper("Post"),
		"/sandbox-ctrl/v1/replay/{name}/generate",
		CreateReplayFileFromScenarioExec,
	},

	Route{
		"DeleteReplayFile",
		strings.ToUpper("Delete"),
		"/sandbox-ctrl/v1/replay/{name}",
		DeleteReplayFile,
	},

	Route{
		"DeleteReplayFileList",
		strings.ToUpper("Delete"),
		"/sandbox-ctrl/v1/replay",
		DeleteReplayFileList,
	},

	Route{
		"GetReplayFile",
		strings.ToUpper("Get"),
		"/sandbox-ctrl/v1/replay/{name}",
		GetReplayFile,
	},

	Route{
		"GetReplayFileList",
		strings.ToUpper("Get"),
		"/sandbox-ctrl/v1/replay",
		GetReplayFileList,
	},

	Route{
		"GetReplayStatus",
		strings.ToUpper("Get"),
		"/sandbox-ctrl/v1/replaystatus",
		GetReplayStatus,
	},

	Route{
		"LoopReplay",
		strings.ToUpper("Post"),
		"/sandbox-ctrl/v1/replay/{name}/loop",
		LoopReplay,
	},

	Route{
		"PlayReplayFile",
		strings.ToUpper("Post"),
		"/sandbox-ctrl/v1/replay/{name}/play",
		PlayReplayFile,
	},

	Route{
		"StopReplayFile",
		strings.ToUpper("Post"),
		"/sandbox-ctrl/v1/replay/{name}/stop",
		StopReplayFile,
	},

	Route{
		"SendEvent",
		strings.ToUpper("Post"),
		"/sandbox-ctrl/v1/events/{type}",
		SendEvent,
	},
}
