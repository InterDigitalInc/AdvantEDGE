/*
 * Copyright (c) 2019
 * InterDigital Communications, Inc.
 * All rights reserved.
 *
 * The information provided herein is the proprietary and confidential
 * information of InterDigital Communications, Inc.
 */
package helm

import (
	"os/exec"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
)

func deleteReleases(charts []Chart) error {
	for _, c := range charts {
		var cmd = exec.Command("helm", "delete", c.ReleaseName, "--purge")
		_, err := cmd.CombinedOutput()
		if err != nil {
			log.Error(err)
			return err
		}
	}

	return nil
}
