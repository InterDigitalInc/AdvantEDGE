package utils

import (
	"errors"
	"fmt"
	"os/exec"
	"time"

	"github.com/spf13/cobra"
)

// SendWget command
func SendWget(url string, fileLocation string, cobraCmd *cobra.Command) (err error) {
	err = nil
	verbose, _ := cobraCmd.Flags().GetBool("verbose")

	start := time.Now()
	cmd := exec.Command("wget", "-O", fileLocation, url)
	if verbose {
		fmt.Println("Cmd:", cmd.Args)
	}
	out, err := cmd.CombinedOutput()
	elapsed := time.Since(start)
	if err != nil {
		err = errors.New("Error sending a wget command")
		fmt.Println(err)
	}
	if verbose {
		r := FormatResult("Result: "+string(out), elapsed, cobraCmd)
		fmt.Println(r)
	}

	return err

}
