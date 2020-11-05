/*
 * ETSI GS MEC 012 - Radio Network Information API
 *
 * The ETSI MEC ISG MEC012 Radio Network Information API described using OpenAPI.
 *
 * API version: 2.1.1
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package server

type L2MeasCellInfo struct {
	// It indicates the packet discard rate in percentage of the downlink GBR traffic in a cell, as defined in ETSI TS 136 314 [i.11].
	DlGbrPdrCell int32 `json:"dl_gbr_pdr_cell,omitempty"`
	// It indicates the PRB usage for downlink GBR traffic, as defined in ETSI TS 136 314 [i.11] and ETSI TS 136 423 [i.12].
	DlGbrPrbUsageCell int32 `json:"dl_gbr_prb_usage_cell,omitempty"`
	// It indicates the packet discard rate in percentage of the downlink non-GBR traffic in a cell, as defined in ETSI TS 136 314 [i.11].
	DlNongbrPdrCell int32 `json:"dl_nongbr_pdr_cell,omitempty"`
	// It indicates (in percentage) the PRB usage for downlink non-GBR traffic, as defined in ETSI TS 136 314 [i.11] and ETSI TS 136 423 [i.12].
	DlNongbrPrbUsageCell int32 `json:"dl_nongbr_prb_usage_cell,omitempty"`
	// It indicates (in percentage) the PRB usage for total downlink traffic, as defined in ETSI TS 136 314 [i.11] and ETSI TS 136 423 [i.12].
	DlTotalPrbUsageCell int32 `json:"dl_total_prb_usage_cell,omitempty"`

	Ecgi *Ecgi `json:"ecgi,omitempty"`
	// It indicates the number of active UEs with downlink GBR traffic, as defined in ETSI TS 136 314 [i.11].
	NumberOfActiveUeDlGbrCell int32 `json:"number_of_active_ue_dl_gbr_cell,omitempty"`
	// It indicates the number of active UEs with downlink non-GBR traffic, as defined in ETSI TS 136 314 [i.11].
	NumberOfActiveUeDlNongbrCell int32 `json:"number_of_active_ue_dl_nongbr_cell,omitempty"`
	// It indicates the number of active UEs with uplink GBR traffic, as defined in ETSI TS 136 314 [i.11].
	NumberOfActiveUeUlGbrCell int32 `json:"number_of_active_ue_ul_gbr_cell,omitempty"`
	// It indicates the number of active UEs with uplink non-GBR traffic, as defined in ETSI TS 136 314 [i.11].
	NumberOfActiveUeUlNongbrCell int32 `json:"number_of_active_ue_ul_nongbr_cell,omitempty"`
	// It indicates (in percentage) the received dedicated preamples, as defined in ETSI TS 136 314 [i.11].
	ReceivedDedicatedPreamblesCell int32 `json:"received_dedicated_preambles_cell,omitempty"`
	// It indicates (in percentage) the received randomly selected preambles in the high range, as defined in ETSI TS 136 314 [i.11].
	ReceivedRandomlySelectedPreamblesHighRangeCell int32 `json:"received_randomly_selected_preambles_high_range_cell,omitempty"`
	// It indicates (in percentage) the received randomly selected preambles in the low range, as defined in ETSI TS 136 314 [i.11].
	ReceivedRandomlySelectedPreamblesLowRangeCell int32 `json:"received_randomly_selected_preambles_low_range_cell,omitempty"`
	// It indicates the packet discard rate in percentage of the uplink GBR traffic in a cell, as defined in ETSI TS 136 314 [i.11].
	UlGbrPdrCell int32 `json:"ul_gbr_pdr_cell,omitempty"`
	// It indicates (in percentage) the PRB usage for uplink GBR traffic, as defined in ETSI TS 136 314 [i.11] and ETSI TS 136 423 [i.12].
	UlGbrPrbUsageCell int32 `json:"ul_gbr_prb_usage_cell,omitempty"`
	// It indicates the packet discard rate in percentage of the uplink non-GBR traffic in a cell, as defined in ETSI TS 136 314 [i.11].
	UlNongbrPdrCell int32 `json:"ul_nongbr_pdr_cell,omitempty"`
	// It indicates (in percentage) the PRB usage for uplink non-GBR traffic, as defined in ETSI TS 136 314 [i.11] and ETSI TS 136 423 [i.12].
	UlNongbrPrbUsageCell int32 `json:"ul_nongbr_prb_usage_cell,omitempty"`
	// It indicates (in percentage) the PRB usage for total uplink traffic, as defined in ETSI TS 136 314 [i.11] and ETSI TS 136 423 [i.12].
	UlTotalPrbUsageCell int32 `json:"ul_total_prb_usage_cell,omitempty"`
}
