/*
 * Location API
 *
 * The ETSI MEC ISG MEC012 Location API described using OpenAPI. The API is based on the Open Mobile Alliance's specification RESTful Network API for Zonal Presence
 *
 * API version: 1.1.1
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package server

// UserEventType : User event
type UserEventType string

// List of UserEventType
const (
	ENTERING     UserEventType = "Entering"
	LEAVING      UserEventType = "Leaving"
	TRANSFERRING UserEventType = "Transferring"
)