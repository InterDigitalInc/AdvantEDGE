/*
 * ETSI GS MEC 012 - Radio Network Information API
 *
 * The ETSI MEC ISG MEC012 Radio Network Information API described using OpenAPI.
 *
 * API version: 2.1.1
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package server

// The temporary identifier allocated for the specific UE as defined below.
type RabEstNotificationTempUeId struct {
	// MMEC as defined in ETSI TS 136 413 [i.3].
	Mmec string `json:"mmec"`
	// M-TMSI as defined in ETSI TS 136 413 [i.3].
	Mtmsi string `json:"mtmsi"`
}
