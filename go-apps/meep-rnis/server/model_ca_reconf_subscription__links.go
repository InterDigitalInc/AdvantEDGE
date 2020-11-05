/*
 * ETSI GS MEC 012 - Radio Network Information API
 *
 * The ETSI MEC ISG MEC012 Radio Network Information API described using OpenAPI.
 *
 * API version: 2.1.1
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package server

// Hyperlink related to the resource. This shall be only included in the HTTP responses and in HTTP PUT requests.
type CaReconfSubscriptionLinks struct {
	Self *LinkType `json:"self"`
}
