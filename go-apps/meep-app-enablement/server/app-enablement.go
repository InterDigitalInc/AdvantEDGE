/*
 * Copyright (c) 2020  InterDigital Communications, Inc
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

package server

import (
	"errors"
	"os"
	"strings"
	"sync"

	appInfo "github.com/InterDigitalInc/AdvantEDGE/go-apps/meep-app-enablement/server/app-info"
	appSupport "github.com/InterDigitalInc/AdvantEDGE/go-apps/meep-app-enablement/server/app-support"
	servMgmt "github.com/InterDigitalInc/AdvantEDGE/go-apps/meep-app-enablement/server/service-mgmt"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
)

const serviceName = "Edge Platform Application Enablement Service"

var mutex sync.Mutex
var sandboxName string

// Init - EPAE Service initialization
func Init() (err error) {

	// Retrieve Sandbox name from environment variable
	sandboxNameEnv := strings.TrimSpace(os.Getenv("MEEP_SANDBOX_NAME"))
	if sandboxNameEnv != "" {
		sandboxName = sandboxNameEnv
	}
	if sandboxName == "" {
		err = errors.New("MEEP_SANDBOX_NAME env variable not set")
		log.Error(err.Error())
		return err
	}
	log.Info("MEEP_SANDBOX_NAME: ", sandboxName)

	err = servMgmt.Init(&mutex)
	if err != nil {
		return err
	}

	err = appSupport.Init(&mutex)
	if err != nil {
		return err
	}

	err = appInfo.Init(&mutex)
	if err != nil {
		return err
	}

	return nil
}

// Run - Start
func Run() (err error) {

	err = servMgmt.Run()
	if err != nil {
		return err
	}

	err = appSupport.Run()
	if err != nil {
		return err
	}

	err = appInfo.Run()
	if err != nil {
		return err
	}

	return nil
}
