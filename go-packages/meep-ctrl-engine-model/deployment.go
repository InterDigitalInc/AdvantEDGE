/*
 * Copyright (c) 2019
 * InterDigital Communications, Inc.
 * All rights reserved.
 *
 * The information provided herein is the proprietary and confidential
 * information of InterDigital Communications, Inc.
 */

package model

// Network deployment object
type Deployment struct {

	// Latency in ms between domains
	InterDomainLatency int32 `json:"interDomainLatency,omitempty"`

	// Latency variation in ms between domains
	InterDomainLatencyVariation int32 `json:"interDomainLatencyVariation,omitempty"`

	// The limit of the traffic supported between domains
	InterDomainThroughput int32 `json:"interDomainThroughput,omitempty"`

	// Packet lost (in terms of percentage) between domains
	InterDomainPacketLoss float64 `json:"interDomainPacketLoss,omitempty"`

	Domains []Domain `json:"domains,omitempty"`
}
