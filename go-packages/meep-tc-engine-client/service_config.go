/*
 * Copyright (c) 2019
 * InterDigital Communications, Inc.
 * All rights reserved.
 *
 * The information provided herein is the proprietary and confidential
 * information of InterDigital Communications, Inc.
 */

package client

// Service object
type ServiceConfig struct {

	// Unique service name
	Name string `json:"name,omitempty"`

	// Multi-Edge service name, if any
	MeSvcName string `json:"meSvcName,omitempty"`

	Ports []ServicePort `json:"ports,omitempty"`
}
