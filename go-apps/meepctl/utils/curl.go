package utils

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// SendCurl command
func SendCurl(cmdtype string, url string, argument string, cobraCmd *cobra.Command) (err error) {
	err = nil
	verbose, _ := cobraCmd.Flags().GetBool("verbose")

	start := time.Now()
	var cmd *exec.Cmd
	cmd = exec.Command("curl", "-vX", cmdtype, url, "-H", "Content-Type: application/json", "-H", "kbn-xsrf: true", "-d", argument)

	if verbose {
		fmt.Println("Cmd:", cmd.Args)
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		err = errors.New("Error sending a curl command")
	} else {
		str := string(out)
		isServiceUnavailable := strings.Contains(str, "Service Unavailable")

		if isServiceUnavailable == true {
			err = errors.New("Error: Service Unavailable")
		}
	}

	elapsed := time.Since(start)
	if verbose {
		if err != nil {
			fmt.Println(err)
		}
		r := FormatResult("Result: "+string(out), elapsed, cobraCmd)
		fmt.Println(r)
	}
	return err

}
