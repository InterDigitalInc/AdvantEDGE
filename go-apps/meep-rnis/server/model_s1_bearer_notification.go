/*
 * ETSI GS MEC 012 - Radio Network Information API
 *
 * The ETSI MEC ISG MEC012 Radio Network Information API described using OpenAPI.
 *
 * API version: 2.1.1
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package server

type S1BearerNotification struct {
	// Shall be set to \"S1BearerNotification\".
	NotificationType string `json:"notificationType"`
	// The subscribed event that triggered this notification in S1BearerSubscription.
	S1Event int32 `json:"s1Event"`
	// Information on specific UE that matches the criteria in S1BearerSubscription as defined below.
	S1UeInfo []S1BearerNotificationS1UeInfo `json:"s1UeInfo"`

	TimeStamp *TimeStamp `json:"timeStamp,omitempty"`
}
