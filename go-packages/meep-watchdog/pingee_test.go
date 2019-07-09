package watchdog

import (
	"fmt"
	"testing"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
)

const pingeeRedisAddr string = "localhost:30379"
const pingeeName string = "pingee-tester"

func TestNewPingee(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())
	// Keep this one first...
	fmt.Println("Invalid Redis DB address")
	_, err := NewPingee("ExpectedFailure-InvalidDbLocation", pingeeName)
	if err == nil {
		t.Errorf("Should report error on invalid Redis db")
	}

	fmt.Println("Create normal")
	_, err = NewPingee(pingeeRedisAddr, pingeeName)
	if err != nil {
		t.Errorf("Unable to create pingee")
	}

	fmt.Println("Create no name")
	_, err = NewPingee(pingeeRedisAddr, "")
	if err == nil {
		t.Errorf("Should not allow creating pingee without a name")
	}
}
