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
 * AdvantEDGE WLAN Access Information API
 *
 * WLAN Access Information Service is AdvantEDGE's implementation of [ETSI MEC ISG MEC028 WAI API](http://www.etsi.org/deliver/etsi_gs/MEC/001_099/028/02.01.01_60/gs_MEC028v020101p.pdf) <p>[Copyright (c) ETSI 2020](https://forge.etsi.org/etsi-forge-copyright-notice.txt) <p>**Micro-service**<br>[meep-wais](https://github.com/InterDigitalInc/AdvantEDGE/tree/master/go-apps/meep-wais) <p>**Type & Usage**<br>Edge Service used by edge applications that want to get information about WLAN access information in the network <p>**Details**<br>API details available at _your-AdvantEDGE-ip-address/api_ <p>AdvantEDGE supports a selected subset of WAI API subscription types. <p>Supported subscriptions: <p> - AssocStaSubscription
 *
 * API version: 2.1.1
 * Contact: AdvantEDGE@InterDigital.com
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package server

type StatisticsGroupData struct {
	Dot11AMPDUDelimiterCRCErrorCount int32 `json:"dot11AMPDUDelimiterCRCErrorCount,omitempty"`

	Dot11AMPDUReceivedCount int32 `json:"dot11AMPDUReceivedCount,omitempty"`

	Dot11AMSDUAckFailureCount int32 `json:"dot11AMSDUAckFailureCount,omitempty"`

	Dot11AckFailureCount int32 `json:"dot11AckFailureCount,omitempty"`

	Dot11BeamformingFrameCount int32 `json:"dot11BeamformingFrameCount,omitempty"`

	Dot11ChannelWidthSwitchCount int32 `json:"dot11ChannelWidthSwitchCount,omitempty"`

	Dot11DualCTSFailureCount int32 `json:"dot11DualCTSFailureCount,omitempty"`

	Dot11DualCTSSuccessCount int32 `json:"dot11DualCTSSuccessCount,omitempty"`

	Dot11ExplicitBARFailureCount int32 `json:"dot11ExplicitBARFailureCount,omitempty"`

	Dot11FCSErrorCount int32 `json:"dot11FCSErrorCount,omitempty"`

	Dot11FailedAMSDUCount int32 `json:"dot11FailedAMSDUCount,omitempty"`

	Dot11FailedCount int32 `json:"dot11FailedCount,omitempty"`

	Dot11FortyMHzFrameReceivedCount int32 `json:"dot11FortyMHzFrameReceivedCount,omitempty"`

	Dot11FortyMHzFrameTransmittedCount int32 `json:"dot11FortyMHzFrameTransmittedCount,omitempty"`

	Dot11FrameDuplicateCount int32 `json:"dot11FrameDuplicateCount,omitempty"`

	Dot11GrantedRDGUnusedCount int32 `json:"dot11GrantedRDGUnusedCount,omitempty"`

	Dot11GrantedRDGUsedCount int32 `json:"dot11GrantedRDGUsedCount,omitempty"`

	Dot11GroupReceivedFrameCount int32 `json:"dot11GroupReceivedFrameCount,omitempty"`

	Dot11GroupTransmittedFrameCount int32 `json:"dot11GroupTransmittedFrameCount,omitempty"`

	Dot11ImplicitBARFailureCount int32 `json:"dot11ImplicitBARFailureCount,omitempty"`

	Dot11MPDUInReceivedAMPDUCount int32 `json:"dot11MPDUInReceivedAMPDUCount,omitempty"`

	Dot11MultipleRetryAMSDUCount int32 `json:"dot11MultipleRetryAMSDUCount,omitempty"`

	Dot11MultipleRetryCount int32 `json:"dot11MultipleRetryCount,omitempty"`

	Dot11PSMPUTTGrantDuration int32 `json:"dot11PSMPUTTGrantDuration,omitempty"`

	Dot11PSMPUTTUsedDuration int32 `json:"dot11PSMPUTTUsedDuration,omitempty"`

	Dot11QosAckFailureCount int32 `json:"dot11QosAckFailureCount,omitempty"`

	Dot11QosDiscardedFrameCount int32 `json:"dot11QosDiscardedFrameCount,omitempty"`

	Dot11QosFailedCount int32 `json:"dot11QosFailedCount,omitempty"`

	Dot11QosFrameDuplicateCount int32 `json:"dot11QosFrameDuplicateCount,omitempty"`

	Dot11QosMPDUsReceivedCount int32 `json:"dot11QosMPDUsReceivedCount,omitempty"`

	Dot11QosMultipleRetryCount int32 `json:"dot11QosMultipleRetryCount,omitempty"`

	Dot11QosRTSFailureCount int32 `json:"dot11QosRTSFailureCount,omitempty"`

	Dot11QosRTSSuccessCount int32 `json:"dot11QosRTSSuccessCount,omitempty"`

	Dot11QosReceivedFragmentCount int32 `json:"dot11QosReceivedFragmentCount,omitempty"`

	Dot11QosRetriesReceivedCount int32 `json:"dot11QosRetriesReceivedCount,omitempty"`

	Dot11QosRetryCount int32 `json:"dot11QosRetryCount,omitempty"`

	Dot11QosTransmittedFragmentCount int32 `json:"dot11QosTransmittedFragmentCount,omitempty"`

	Dot11QosTransmittedFrameCount int32 `json:"dot11QosTransmittedFrameCount,omitempty"`

	Dot11RSNAStatsBIPMICErrors int32 `json:"dot11RSNAStatsBIPMICErrors,omitempty"`

	Dot11RSNAStatsCCMPDecryptErrors int32 `json:"dot11RSNAStatsCCMPDecryptErrors,omitempty"`

	Dot11RSNAStatsCCMPReplays int32 `json:"dot11RSNAStatsCCMPReplays,omitempty"`

	Dot11RSNAStatsCMACReplays int32 `json:"dot11RSNAStatsCMACReplays,omitempty"`

	Dot11RSNAStatsRobustMgmtCCMPReplays int32 `json:"dot11RSNAStatsRobustMgmtCCMPReplays,omitempty"`

	Dot11RSNAStatsTKIPICVErrors int32 `json:"dot11RSNAStatsTKIPICVErrors,omitempty"`

	Dot11RSNAStatsTKIPReplays int32 `json:"dot11RSNAStatsTKIPReplays,omitempty"`

	Dot11RTSFailureCount int32 `json:"dot11RTSFailureCount,omitempty"`

	Dot11RTSLSIGFailureCount int32 `json:"dot11RTSLSIGFailureCount,omitempty"`

	Dot11RTSLSIGSuccessCount int32 `json:"dot11RTSLSIGSuccessCount,omitempty"`

	Dot11RTSSuccessCount int32 `json:"dot11RTSSuccessCount,omitempty"`

	Dot11ReceivedAMSDUCount int32 `json:"dot11ReceivedAMSDUCount,omitempty"`

	Dot11ReceivedFragmentCount int32 `json:"dot11ReceivedFragmentCount,omitempty"`

	Dot11ReceivedOctetsInAMPDUCount int64 `json:"dot11ReceivedOctetsInAMPDUCount,omitempty"`

	Dot11ReceivedOctetsInAMSDUCount int64 `json:"dot11ReceivedOctetsInAMSDUCount,omitempty"`

	Dot11RetryAMSDUCount int32 `json:"dot11RetryAMSDUCount,omitempty"`

	Dot11RetryCount int32 `json:"dot11RetryCount,omitempty"`

	Dot11STAStatisticsAPAverageAccessDelay int32 `json:"dot11STAStatisticsAPAverageAccessDelay,omitempty"`

	Dot11STAStatisticsAverageAccessDelayBackGround int32 `json:"dot11STAStatisticsAverageAccessDelayBackGround,omitempty"`

	Dot11STAStatisticsAverageAccessDelayBestEffort int32 `json:"dot11STAStatisticsAverageAccessDelayBestEffort,omitempty"`

	Dot11STAStatisticsAverageAccessDelayVideo int32 `json:"dot11STAStatisticsAverageAccessDelayVideo,omitempty"`

	Dot11STAStatisticsAverageAccessDelayVoice int32 `json:"dot11STAStatisticsAverageAccessDelayVoice,omitempty"`

	Dot11STAStatisticsChannelUtilization int32 `json:"dot11STAStatisticsChannelUtilization,omitempty"`

	Dot11STAStatisticsStationCount int32 `json:"dot11STAStatisticsStationCount,omitempty"`

	Dot11STBCCTSFailureCount int32 `json:"dot11STBCCTSFailureCount,omitempty"`

	Dot11STBCCTSSuccessCount int32 `json:"dot11STBCCTSSuccessCount,omitempty"`

	Dot11TransmittedAMPDUCount int32 `json:"dot11TransmittedAMPDUCount,omitempty"`

	Dot11TransmittedAMSDUCount int32 `json:"dot11TransmittedAMSDUCount,omitempty"`

	Dot11TransmittedFragmentCount int32 `json:"dot11TransmittedFragmentCount,omitempty"`

	Dot11TransmittedFrameCount int32 `json:"dot11TransmittedFrameCount,omitempty"`

	Dot11TransmittedFramesInGrantedRDGCount int32 `json:"dot11TransmittedFramesInGrantedRDGCount,omitempty"`

	Dot11TransmittedMPDUsInAMPDUCount int32 `json:"dot11TransmittedMPDUsInAMPDUCount,omitempty"`

	Dot11TransmittedOctetsInAMPDUCount int64 `json:"dot11TransmittedOctetsInAMPDUCount,omitempty"`

	Dot11TransmittedOctetsInAMSDUCount int64 `json:"dot11TransmittedOctetsInAMSDUCount,omitempty"`

	Dot11TransmittedOctetsInGrantedRDGCount int32 `json:"dot11TransmittedOctetsInGrantedRDGCount,omitempty"`

	Dot11TwentyMHzFrameReceivedCount int32 `json:"dot11TwentyMHzFrameReceivedCount,omitempty"`

	Dot11TwentyMHzFrameTransmittedCount int32 `json:"dot11TwentyMHzFrameTransmittedCount,omitempty"`

	Dot11nonSTBCCTSFailureCount int32 `json:"dot11nonSTBCCTSFailureCount,omitempty"`

	Dot11nonSTBCCTSSuccessCount int32 `json:"dot11nonSTBCCTSSuccessCount,omitempty"`
}
