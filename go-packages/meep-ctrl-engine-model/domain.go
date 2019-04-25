/*
 * Copyright (c) 2019
 * InterDigital Communications, Inc.
 * All rights reserved.
 *
 * The information provided herein is the proprietary and confidential
 * information of InterDigital Communications, Inc.
 */

package model

// Operator domain object
type Domain struct {

	// Unique domain ID
	Id string `json:"id,omitempty"`

	// Domain name
	Name string `json:"name,omitempty"`

	// Domain type
	Type_ string `json:"type,omitempty"`

	// Latency in ms between zones within domain
	InterZoneLatency int32 `json:"interZoneLatency,omitempty"`

	// Latency variation in ms between zones within domain
	InterZoneLatencyVariation int32 `json:"interZoneLatencyVariation,omitempty"`

	// The limit of the traffic supported between zones within the domain
	InterZoneThroughput int32 `json:"interZoneThroughput,omitempty"`

	// Packet lost (in terms of percentage) between zones within the domain
	InterZonePacketLoss float64 `json:"interZonePacketLoss,omitempty"`

	Zones []Zone `json:"zones,omitempty"`
}
