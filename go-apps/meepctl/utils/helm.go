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

package utils

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// IsHelmRelease  Returns true if a Helm release exists
func IsHelmRelease(name string, cobraCmd *cobra.Command) (exist bool, err error) {
	exist = false
	err = nil
	verbose, _ := cobraCmd.Flags().GetBool("verbose")

	start := time.Now()
	cmd := exec.Command("helm", "ls", "--filter", name, "--short")
	if verbose {
		fmt.Println("Cmd:", cmd.Args)
	}
	out, err := cmd.CombinedOutput()
	elapsed := time.Since(start)
	if err != nil {
		err = errors.New("Error listing component [" + name + "]")
		fmt.Println(err)
	} else {
		s := string(out)
		lines := strings.Split(s, "\n")
		for _, line := range lines {
			if line == name {
				exist = true
				break
			}
		}
	}
	if verbose {
		r := FormatResult("Result: "+string(out), elapsed, cobraCmd)
		fmt.Println(r)
	}

	return exist, err
}

// HelmDelete  Deletes specified release
func HelmDelete(name string, cobraCmd *cobra.Command) (err error) {
	err = nil
	verbose, _ := cobraCmd.Flags().GetBool("verbose")

	start := time.Now()
	cmd := exec.Command("helm", "uninstall", name)
	if verbose {
		fmt.Println("Cmd:", cmd.Args)
	}
	out, err := cmd.CombinedOutput()
	elapsed := time.Since(start)

	if err != nil {
		err = errors.New("Error deleting component [" + name + "]")
		fmt.Println(err)
	} else {
		r := FormatResult("Deleted "+name, elapsed, cobraCmd)
		fmt.Println(r)
	}
	if verbose {
		fmt.Println("Result: " + string(out))
	}

	return err
}

// HelmInstall Install specified releases
func HelmInstall(name string, chart string, flags [][]string, cobraCmd *cobra.Command) (err error) {
	err = nil
	verbose, _ := cobraCmd.Flags().GetBool("verbose")

	start := time.Now()
	cmd := exec.Command("helm", "install", name, "--set", "fullnameOverride="+name, chart, "--replace")
	for _, f := range flags {
		cmd.Args = append(cmd.Args, f[0])
		cmd.Args = append(cmd.Args, f[1])
	}
	if verbose {
		fmt.Println("Cmd:", cmd.Args)
	}
	out, err := cmd.CombinedOutput()
	elapsed := time.Since(start)
	if err != nil {
		err = errors.New("Error installing component [" + name + "]")
		fmt.Println(err)
	} else {
		r := FormatResult("Deployed "+name, elapsed, cobraCmd)
		fmt.Println(r)
	}
	if verbose {
		fmt.Println("Result: " + string(out))
	}
	return err
}

// HelmFlags Takes helm flag & value pair and formats it into an array of flag value pair
func HelmFlags(flagsIn [][]string, flag string, value string) (flagsOut [][]string) {
	if flagsIn == nil {
		flagsOut = make([][]string, 0)
	} else {
		flagsOut = flagsIn
	}
	if flag != "" {
		f := make([]string, 0)
		f = append(f, flag)
		f = append(f, value)
		flagsOut = append(flagsOut, f)
	}
	return flagsOut

}
