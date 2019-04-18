package helm

import (
	"bufio"
	"errors"
	"os/exec"
	"strings"

	log "github.com/InterDigitalInc/AdvantEDGE/go-apps/meep-virt-engine/log"
)

const NAMESPACE string = "NAMESPACE:"
const STATUS string = "STATUS:"
const RESOURCE string = "==>"

// Returns the status of a release
func GetReleaseStatus(name string) (*Status, error) {
	out, err := getStatus(name)
	if err != nil {
		return nil, err
	}

	status, err := parseStatus(out)
	if err != nil {
		return nil, err
	}
	return status, nil
}

func getStatus(name string) ([]byte, error) {
	var cmd = exec.Command("helm", "status", name)
	out, err := cmd.Output()
	if err != nil {
		err = errors.New("Error getting status for Release [" + name + "]")
		log.Error(err)
		return nil, err
	}
	return out, nil
}

func parseStatus(buf []byte) (*Status, error) {
	var status Status

	s := string(buf)
	scanLines := bufio.NewScanner(strings.NewReader(s))
	scanLines.Split(bufio.ScanLines)
	for i := 0; scanLines.Scan(); i++ {
		scanWords := bufio.NewScanner(strings.NewReader(scanLines.Text()))
		scanWords.Split(bufio.ScanWords)
		scanWords.Scan()
		word := scanWords.Text()

		if word == NAMESPACE {
			scanWords.Scan()
			status.Namespace = scanWords.Text()
		} else if word == STATUS {
			scanWords.Scan()
			status.State = scanWords.Text()
		} else if word == "==>" {
			var r Resource
			// Scan Type
			scanWords.Scan()
			t := strings.Split(scanWords.Text(), "/")
			r.Type = t[1]

			// Skip a line
			scanLines.Scan()

			// Scan Name
			scanLines.Scan()
			scanRes := bufio.NewScanner(strings.NewReader(scanLines.Text()))
			scanRes.Split(bufio.ScanWords)
			scanRes.Scan()
			r.Name = scanRes.Text()
			for scanRes.Scan() {
				r.Age = scanRes.Text()
			}
			status.Resources = append(status.Resources, r)
		}
	}
	return &status, nil
}
