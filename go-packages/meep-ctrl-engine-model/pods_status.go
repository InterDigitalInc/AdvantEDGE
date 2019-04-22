/*
 * Copyright (c) 2019
 * InterDigital Communications, Inc.
 * All rights reserved.
 *
 * The information provided herein is the proprietary and confidential
 * information of InterDigital Communications, Inc.
 */

package model

// List of all pods status
type PodsStatus struct {
	PodStatus []PodStatus `json:"podStatus,omitempty"`
}
