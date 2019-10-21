/*
 * Copyright (c) 2019  InterDigital Communications, Inc
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package bws

import (
	"fmt"
	"testing"
	"time"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
)

const redisAddr string = "localhost:30379"

func TestBwsharingBasic(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	bwSharing, err := NewBwSharing("test", redisAddr, nil, nil)
	if err != nil {
		t.Errorf("Failed to create a bwSharing object.")
	} else {
		bwSharing.UpdateControls()
		_ = bwSharing.Start()

		time.Sleep(1000 * time.Millisecond)
		bwSharing.Stop()
	}
}
