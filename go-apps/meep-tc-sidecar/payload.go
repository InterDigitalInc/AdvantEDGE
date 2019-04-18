package main

import (
	"math/rand"
	log "github.com/InterDigitalInc/AdvantEDGE/go-apps/meep-tc-sidecar/log"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Payload represents additional data appended to outgoing ICMP Echo
// Requests.
type Payload []byte

// Resize will assign a new payload of the given size to p.
func (p *Payload) Resize(size uint16) {
	buf := make([]byte, size, size)
	if _, err := rand.Read(buf); err != nil {
		log.Error("error resizing payload: ", err)
		return
	}
	*p = Payload(buf)
}
