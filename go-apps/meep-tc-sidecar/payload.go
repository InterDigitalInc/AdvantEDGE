/*
 * Copyright (c) 2019
 * InterDigital Communications, Inc.
 * All rights reserved.
 *
 * The information provided herein is the proprietary and confidential
 * information of InterDigital Communications, Inc.
 */
package main

import (
	"math/rand"
	"time"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Payload represents additional data appended to outgoing ICMP Echo
// Requests.
type Payload []byte

// Resize will assign a new payload of the given size to p.
func (p *Payload) Resize(size uint16) {
	buf := make([]byte, size)
	if _, err := rand.Read(buf); err != nil {
		log.Error("error resizing payload: ", err)
		return
	}
	*p = Payload(buf)
}
