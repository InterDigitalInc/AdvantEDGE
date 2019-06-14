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
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/roymx/viper"
	yaml "gopkg.in/yaml.v2"
)

var RepoCfg *viper.Viper

// Config version needs to be bumped only when new elements are added
var defaultConfig = `
version: 1.0.0

node:
  ip: ""

meep:
  gitdir: ""
  workdir: "<DEFAULT>/.meep"
  registry: "meep-docker-registry:30001"
`

// Node parameters node
type Node struct {
	// Node IP Address
	IP string `json:"ip,omitempty"`
}

// Meep parameters node
type Meep struct {
	// GIT directory
	Gitdir string `json:"gitdir,omitempty"`
	// MEEP work directory
	Workdir string `json:"workdir,omitempty"`
	// MEEP docker registry
	Registry string `json:"registry,omitempty"`
}

// Config structure
type Config struct {
	Version string `json:"version,omitempty"`
	// Node parameters
	Node Node `json:"ip,omitempty"`
	// Meep parameters
	Meep Meep `json:"meep,omitempty"`
}

// ConfigGetDefaultPath get default config file path
func ConfigGetDefaultPath() (path string) {
	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		return path
	}
	return home + "/.meepctl.yaml"
}

// ConfigCreateIfNotExist test for config file existence & create default if needed
func ConfigCreateIfNotExist() (exist bool, err error) {
	path := ConfigGetDefaultPath()
	if path == "" {
		os.Exit(1)
	}

	// meepctl uses this config file located in $(HOME)/.meepctl.yaml
	// If it does not exist, create one.
	_, err = os.Stat(path)
	if err == nil {
		// we're good
		return true, nil
	} else if os.IsNotExist(err) {
		// file does not exist so create default
		cfg := ConfigReadFile("") // returns default config
		err = ConfigWriteFile(cfg, path)
		fmt.Println("Creating default config file @ " + path)
		return true, err
	} // else file may exist ... just return the error
	return true, err
}

// ConfigReadFile read the configuration file
func ConfigReadFile(filePath string) (cfg *Config) {
	var err error
	var data []byte
	cfg = new(Config)

	// Read from config file
	if filePath != "" {
		data, err = ioutil.ReadFile(filePath)
		if err != nil {
			fmt.Println("Error reading config file [" + filePath + "]")
			fmt.Println(err)
			return nil
		}
	}
	// Revert to default if readfile failed
	if len(data) == 0 {
		data = []byte(defaultConfig)
		str := fmt.Sprintf("%s", data)
		home, _ := homedir.Dir()
		newStr := strings.Replace(str, "<DEFAULT>", home, -1)
		data = []byte(newStr)
	}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		fmt.Println("Error unmarshalling config file [" + filePath + "]")
		fmt.Println(err)
		// Revert to default if unmarshall failed
		_ = yaml.Unmarshal([]byte(defaultConfig), cfg)
	}

	return cfg
}

// ConfigValidate validates config file
func ConfigValidate(filePath string) (valid bool) {
	if filePath == "" {
		filePath = ConfigGetDefaultPath()
	}
	cfg := ConfigReadFile(filePath)
	configValid := true

	// Validate IPV4
	valid, reason := ConfigIPValid(cfg.Node.IP)
	if !valid {
		fmt.Println("")
		fmt.Println("  WARNING    invalid meepctl config: node.ip")
		fmt.Println("             Reason: " + reason)
		fmt.Println("             Fix with:  meepctl config ip <node-ip-address>")
		fmt.Println("")
		configValid = false
	}

	// Validate Gitdir
	valid, reason = ConfigGitdirValid(cfg.Meep.Gitdir)
	if !valid {
		fmt.Println("")
		fmt.Println("  WARNING    invalid meepctl config: meep.gitdir")
		fmt.Println("             Reason: " + reason)
		fmt.Println("             Fix with:  meepctl config gitdir <path-to-gitdir>")
		fmt.Println("")
		configValid = false
	}
	return configValid
}

// ConfigGitdirValid validates IP address
func ConfigGitdirValid(gitdir string) (valid bool, reason string) {
	valid = true
	fi, err := os.Stat(gitdir)

	if err != nil {
		reason = "Path error  [" + gitdir + "]"
		valid = false
	} else {
		if !fi.IsDir() {
			reason = "Not a directory [" + gitdir + "]"
			valid = false
		}
	}
	return valid, reason
}

// ConfigIPValid validates IP address
func ConfigIPValid(ipAddr string) (valid bool, reason string) {
	valid = true
	// only ipv4 address
	if ConfigIsIpv4(ipAddr) {
		// not localhost
		ip := net.ParseIP(ipAddr)
		if ip.IsLoopback() {
			reason = "Invalid local IP address [" + ipAddr + "] (loopback)"
			valid = false
		}
		// only local address
		addrs, _ := net.InterfaceAddrs()
		var local = false
		// var localV4 []string
		for _, a := range addrs {
			if strings.Contains(a.String(), ipAddr) {
				local = true
			}
			// aIP := strings.Split(a.String(), "/")[0]
			// if ConfigIsIpv4(aIP) {
			// 	localV4 = append(localV4, aIP)
			// }
		}
		if !local {
			reason = "Not a local IP address [" + ipAddr + "]"
			valid = false
		}
	} else {
		reason = "Not an IPV4 address [" + ipAddr + "]"
		valid = false
	}
	return valid, reason
}

// ConfigWriteFile writes the configuration file
func ConfigWriteFile(cfg *Config, filePath string) (err error) {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		err = errors.New("Error marshalling config")
		return err
	}

	err = ioutil.WriteFile(filePath, data, os.ModePerm)
	if err != nil {
		err = errors.New("Error writing config file [" + filePath + "]")
		return err
	}

	return nil
}

// ConfigIsIpv4 checks if IP address is IPV4
func ConfigIsIpv4(host string) bool {
	parts := strings.Split(host, ".")
	if len(parts) < 4 {
		return false
	}
	for _, x := range parts {
		if i, err := strconv.Atoi(x); err == nil {
			if i < 0 || i > 255 {
				return false
			}
		} else {
			return false
		}
	}
	return true
}

// InitRepoConfig initializes & returns the repo config
func InitRepoConfig() *viper.Viper {
	repodir := viper.GetString("meep.gitdir")
	RepoCfg = viper.New()
	RepoCfg.SetConfigFile(repodir + "/.meepctl-repocfg.yaml")
	if err := RepoCfg.ReadInConfig(); err == nil {
		fmt.Println("Using repo config file:", RepoCfg.ConfigFileUsed())
	} else {
		RepoCfg = nil
		fmt.Println(err)
	}
	return RepoCfg
}
