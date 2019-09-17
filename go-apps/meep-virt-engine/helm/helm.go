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

package helm

import (
	"fmt"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
)

const (
	StateIdle       = "IDLE"
	StateInstalling = "INSTALLING"
	StateDeleting   = "DELETING"
)

var state string = StateIdle

func GetReleasesName() ([]Release, error) {
	return getReleasesName()
}

func GetReleases() ([]Release, error) {
	return getReleases()
}

func InstallCharts(charts []Chart) error {
	if state == StateIdle {
		state = StateInstalling
		go func() {
			log.Debug("Installing ", len(charts), " Charts...")
			_ = installCharts(charts)
			log.Debug("Charts installed (", len(charts), ")")
			state = StateIdle
		}()
		return nil
	}
	err := fmt.Errorf("Service busy [%s]", state)
	return err
}

func DeleteReleases(charts []Chart) error {
	if state == StateIdle {
		state = StateDeleting
		go func() {
			log.Debug("Deleting ", len(charts), " Releases...")
			_ = deleteReleases(charts)
			log.Debug("Releases deleted (", len(charts), ")")
			state = StateIdle
		}()
		return nil
	}
	err := fmt.Errorf("Service busy [%s]", state)
	return err
}
