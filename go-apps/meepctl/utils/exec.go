/*
 * Copyright (c) 2019
 * InterDigital Communications, Inc.
 * All rights reserved.
 *
 * The information provided herein is the proprietary and confidential
 * information of InterDigital Communications, Inc.
 */

package utils

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// ExecuteCmd wraps exec.Command.CombinedOutput with verbosity
func ExecuteCmd(cmd *exec.Cmd, cobraCmd *cobra.Command) (output string, err error) {
	err = nil
	verbose, _ := cobraCmd.Flags().GetBool("verbose")

	if verbose {
		fmt.Println("Cmd:", cmd.Args)
	}
	start := time.Now()
	out, err := cmd.CombinedOutput()
	elapsed := time.Since(start)
	output = string(out)
	if verbose {
		var r string
		if len(output) > 0 {
			r = "Result: " + output
		} else {
			r = "Result: OK"
		}
		r = FormatResult(r, elapsed, cobraCmd)
		fmt.Println(r)
	}
	return output, err
}

// KillProcess kill a process by name
func KillProcess(name string, cobraCmd *cobra.Command) {
	cmd := exec.Command("pidof", name)
	pid, err := ExecuteCmd(cmd, cobraCmd)
	if err == nil {
		pid = strings.TrimSuffix(pid, "\n")
		cmd = exec.Command("kill", "-9", pid)
		_, err = ExecuteCmd(cmd, cobraCmd)
		if err != nil {
			fmt.Println("Error terminating " + name)
			fmt.Println(err)
		}
	}
}

// InterruptProcess SIGINT a process by name
func InterruptProcess(name string, cobraCmd *cobra.Command) {
	cmd := exec.Command("pidof", name)
	pid, err := ExecuteCmd(cmd, cobraCmd)
	if err == nil {
		pid = strings.TrimSuffix(pid, "\n")
		cmd = exec.Command("kill", "-2", pid)
		_, err = ExecuteCmd(cmd, cobraCmd)
		if err != nil {
			fmt.Println("Error interrupting " + name)
			fmt.Println(err)
		}
	}
}
