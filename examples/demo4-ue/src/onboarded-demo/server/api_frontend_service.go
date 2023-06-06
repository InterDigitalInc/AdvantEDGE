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
	"context"
	"syscall"
)

// FrontendApiService is a service that implents the logic for the FrontendApiServicer
// This service should implement the business logic for every endpoint for the FrontendApi API.
// Include any external packages or services that will be required by this service.
type FrontendApiService struct {
}

// NewFrontendApiService creates a default api service
func NewFrontendApiService() FrontendApiServicer {
	return &FrontendApiService{}
}

// Ping - Await for ping request and reply winth pong text body
func (s *FrontendApiService) Ping(ctx context.Context) (interface{}, error) {
	// TODO - update Ping with the required logic for this service method.
	// Add api_frontend_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.
	//return nil, errors.New("service method 'Ping' not implemented")
	return "pong", nil
}

// Terminate - Terminate gracefully the application
func (s *FrontendApiService) Terminate(ctx context.Context) (interface{}, error) {
	// TODO - update Terminate with the required logic for this service method.
	// Add api_frontend_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.
	//return nil, errors.New("service method 'Terminate' not implemented")

	syscall.Kill(syscall.Getpid(), syscall.SIGINT)

	return nil, nil
}