package helm

import (
	"os/exec"
	log "github.com/InterDigitalInc/AdvantEDGE/go-apps/meep-virt-engine/log"
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
