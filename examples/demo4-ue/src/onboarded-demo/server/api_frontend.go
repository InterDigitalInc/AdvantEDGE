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
)

// A FrontendApiController binds http requests to an api service and writes the service results to the http response
type FrontendApiController struct {
	service FrontendApiServicer
}

// NewFrontendApiController creates a default api controller
func NewFrontendApiController(s FrontendApiServicer) Router {
	return &FrontendApiController{service: s}
}

// Routes returns all of the api route for the FrontendApiController
func (c *FrontendApiController) Routes() Routes {
	return Routes{
		{
			"Ping",
			strings.ToUpper("Get"),
			"/ping",
			c.Ping,
		},
		{
			"Terminate",
			strings.ToUpper("Delete"),
			"/",
			c.Terminate,
		},
	}
}

// Ping - Await for ping request and reply winth pong text body
func (c *FrontendApiController) Ping(w http.ResponseWriter, r *http.Request) {
	result, err := c.service.Ping(r.Context())
	if err != nil {
		w.WriteHeader(500)
		return
	}

	EncodeJSONResponse(result, nil, w)
}

// Terminate - Terminate gracefully the application
func (c *FrontendApiController) Terminate(w http.ResponseWriter, r *http.Request) {
	result, err := c.service.Terminate(r.Context())
	if err != nil {
		w.WriteHeader(500)
		return
	}

	EncodeJSONResponse(result, nil, w)
}
