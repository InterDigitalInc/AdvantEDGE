/*
 * ETSI GS MEC 012 - Radio Network Information API
 *
 * The ETSI MEC ISG MEC012 Radio Network Information API described using OpenAPI.
 *
 * API version: 2.1.1
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package server

type ResultsPerCsiRsIndex struct {
	CsiRsIndex int32 `json:"csiRsIndex"`

	CsiRsResults *MeasQuantityResultsNr `json:"csiRsResults,omitempty"`
}
