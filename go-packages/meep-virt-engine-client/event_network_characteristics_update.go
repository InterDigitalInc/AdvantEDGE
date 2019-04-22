/*
 * Copyright (c) 2019
 * InterDigital Communications, Inc.
 * All rights reserved.
 *
 * The information provided herein is the proprietary and confidential
 * information of InterDigital Communications, Inc.
 */

package client

// Network Characteristics update Event object
type EventNetworkCharacteristicsUpdate struct {

	// Name of the network element to be updated
	ElementName string `json:"elementName,omitempty"`

	// Type of the network element to be updated
	ElementType string `json:"elementType,omitempty"`

	// Latency in ms
	Latency int32 `json:"latency,omitempty"`

	// Latency variation in ms
	LatencyVariation int32 `json:"latencyVariation,omitempty"`

	// Throughput limit
	Throughput int32 `json:"throughput,omitempty"`

	// Packet loss percentage
	PacketLoss float64 `json:"packetLoss,omitempty"`
}
