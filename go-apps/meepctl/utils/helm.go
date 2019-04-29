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
	cmd := exec.Command("helm", "ls", name, "--short")
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
		exist = strings.HasPrefix(s, name)

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
	cmd := exec.Command("helm", "delete", name, "--purge")
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
	cmd := exec.Command("helm", "install", "--name", name, chart, "--replace")
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
		err = errors.New("Error intalling component [" + name + "]")
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
