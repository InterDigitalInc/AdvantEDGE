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

// IsDockerImage  Returns true if a Docker image exists
func IsDockerImage(name string, cobraCmd *cobra.Command) (exist bool, err error) {
	exist = false
	err = nil
	verbose, _ := cobraCmd.Flags().GetBool("verbose")

	start := time.Now()
	cmd := exec.Command("docker", "images", name, "-q")
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
		exist = s != ""
	}
	if verbose {
		r := FormatResult("Result: "+string(out), elapsed, cobraCmd)
		fmt.Println(r)
	}

	return exist, err

}

// GetDockerTag  Returns image tag
func GetDockerTag(name string, cobraCmd *cobra.Command) (tag string, err error) {
	err = nil
	verbose, _ := cobraCmd.Flags().GetBool("verbose")

	start := time.Now()
	cmd := exec.Command("docker", "images", name)
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
		sl := strings.Fields(s)
		tag = ""
		for i, f := range sl {
			if f == name {
				tag = sl[i+1]
			}
		}
	}
	if verbose {
		r := FormatResult("Result: "+string(out), elapsed, cobraCmd)
		fmt.Println(r)
	}

	return tag, err

}

// TagDockerImage  Tag a docker image
func TagDockerImage(localRepo string, newRepo string, cobraCmd *cobra.Command) (err error) {
	err = nil
	verbose, _ := cobraCmd.Flags().GetBool("verbose")

	start := time.Now()
	cmd := exec.Command("docker", "tag", localRepo, newRepo)
	if verbose {
		fmt.Println("Cmd:", cmd.Args)
	}
	out, err := cmd.CombinedOutput()
	elapsed := time.Since(start)
	if err != nil {
		err = errors.New("Error tagging Docker image [" + localRepo + "] to [" + newRepo + "]")
		fmt.Println(err)
	}
	if verbose {
		r := FormatResult("Result: "+string(out), elapsed, cobraCmd)
		fmt.Println(r)
	}

	return err

}

// PushDockerImage  Push a docker image
func PushDockerImage(newRepo string, cobraCmd *cobra.Command) (err error) {
	err = nil
	verbose, _ := cobraCmd.Flags().GetBool("verbose")

	start := time.Now()
	cmd := exec.Command("docker", "push", newRepo)
	if verbose {
		fmt.Println("Cmd:", cmd.Args)
	}
	out, err := cmd.CombinedOutput()
	elapsed := time.Since(start)
	if err != nil {
		err = errors.New("Error pushing Docker image [" + newRepo + "]")
		fmt.Println(err)
	}
	if verbose {
		r := FormatResult("Result: "+string(out), elapsed, cobraCmd)
		fmt.Println(r)
	}

	return err

}

// PullDockerImage  Pulls a docker image
func PullDockerImage(repo string, cobraCmd *cobra.Command) (err error) {
	err = nil
	verbose, _ := cobraCmd.Flags().GetBool("verbose")

	start := time.Now()
	cmd := exec.Command("docker", "pull", repo)
	if verbose {
		fmt.Println("Cmd:", cmd.Args)
	}
	out, err := cmd.CombinedOutput()
	elapsed := time.Since(start)
	if err != nil {
		err = errors.New("Error pulling Docker image [" + repo + "]")
		fmt.Println(err)
	}
	if verbose {
		r := FormatResult("Result: "+string(out), elapsed, cobraCmd)
		fmt.Println(r)
	}

	return err
}
