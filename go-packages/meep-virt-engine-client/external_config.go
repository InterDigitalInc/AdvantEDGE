/*
 * Copyright (c) 2019
 * InterDigital Communications, Inc.
 * All rights reserved.
 *
 * The information provided herein is the proprietary and confidential
 * information of InterDigital Communications, Inc.
 */

package client

// External Process configuration. NOTE: Only valid if 'isExternal' is set.
type ExternalConfig struct {
	IngressServiceMap []ServiceMap `json:"ingressServiceMap,omitempty"`

	EgressServiceMap []ServiceMap `json:"egressServiceMap,omitempty"`
}
