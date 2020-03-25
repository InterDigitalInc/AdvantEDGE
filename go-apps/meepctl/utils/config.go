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
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/roymx/viper"
	yaml "gopkg.in/yaml.v2"
)

const configVersion string = "1.4.1"

const defaultNotSet string = "not set"
const defaultIP string = ""
const defaultGitDir string = ""
const defaultWorkDir string = ".meep"
const defaultRegistry string = "meep-docker-registry:30001"

var Cfg *Config
var RepoCfg *viper.Viper

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
	// Node parameters
	Node Node `json:"ip,omitempty"`
	// Meep parameters
	Meep Meep `json:"meep,omitempty"`
}

// ConfigInit initializes the meep configuration
func ConfigInit() bool {

	// Initialize Config variable
	Cfg = &Config{
		Node: Node{
			IP: defaultNotSet,
		},
		Meep: Meep{
			Gitdir:   defaultNotSet,
			Workdir:  defaultNotSet,
			Registry: defaultNotSet,
		},
	}

	// Locate configuration file or create a new one if it does not exist
	// NOTE: meepctl uses config file located in $(HOME)/.meepctl.yaml
	path := ConfigGetDefaultPath()
	if path == "" {
		fmt.Println("Error accessing config file at $(HOME)/.meepctl.yaml")
		os.Exit(1)
	}
	_, err := os.Stat(path)
	if err == nil {
		// Update configuration object from config file
		_ = ConfigReadFile(Cfg, path)
	} else if !os.IsNotExist(err) {
		fmt.Println("Error accessing config file at $(HOME)/.meepctl.yaml")
		fmt.Println(err)
		return false
	}

	// Create default entries if they don't exist
	valuesUpdated := ConfigSetDefaultValues(Cfg)
	if valuesUpdated {
		err = ConfigWriteFile(Cfg, path)
		if err != nil {
			fmt.Println("Failed to update config file with error: " + err.Error())
			return false
		}
		fmt.Println("Updated config file @ " + path)
	}

	// Set Repo config if gitdir is set
	repoCfgFile := Cfg.Meep.Gitdir + "/.meepctl-repocfg.yaml"
	if Cfg.Meep.Gitdir != "" {
		RepoCfg = viper.New()
		RepoCfg.SetConfigFile(repoCfgFile)
		if err = RepoCfg.ReadInConfig(); err == nil {
			fmt.Println("Using repo config file:", RepoCfg.ConfigFileUsed())
		} else {
			RepoCfg = nil
		}
	}

	return true
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

func ConfigSetDefaultValues(cfg *Config) bool {
	updated := false
	if cfg.Node.IP == defaultNotSet {
		cfg.Node.IP = defaultIP
		updated = true
	}
	if cfg.Meep.Gitdir == defaultNotSet {
		cfg.Meep.Gitdir = defaultGitDir
		updated = true
	}
	if cfg.Meep.Workdir == defaultNotSet {
		home, _ := homedir.Dir()
		cfg.Meep.Workdir = home + "/" + defaultWorkDir
		updated = true
	}
	if cfg.Meep.Registry == defaultNotSet {
		cfg.Meep.Registry = defaultRegistry
		updated = true
	}
	return updated
}

// ConfigReadFile read the configuration file
func ConfigReadFile(cfg *Config, filePath string) (err error) {
	if filePath == "" {
		return nil
	}

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading config file [" + filePath + "]")
		fmt.Println(err)
		return err
	}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		fmt.Println("Error unmarshalling config file [" + filePath + "]")
		fmt.Println(err)
		return err
	}

	return nil
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

// ConfigValidate validates config file
func ConfigValidate(filePath string) (valid bool) {
	configValid := true

	// Validate IPV4
	valid, reason := ConfigIPValid(Cfg.Node.IP)
	if !valid {
		fmt.Println("")
		fmt.Println("  WARNING    invalid meepctl config: node.ip")
		fmt.Println("             Reason: " + reason)
		fmt.Println("             Fix:  meepctl config ip <node-ip-address>")
		fmt.Println("")
		configValid = false
	}

	// Validate Gitdir & repo version
	valid, reason = ConfigPathValid(Cfg.Meep.Gitdir)
	if !valid {
		fmt.Println("")
		fmt.Println("  WARNING    invalid meepctl config: meep.gitdir")
		fmt.Println("             Reason: " + reason)
		fmt.Println("             Fix:  meepctl config gitdir <path-to-gitdir>")
		fmt.Println("")
		configValid = false
	} else if RepoCfg == nil {
		fmt.Println("")
		fmt.Println("  WARNING    repocfg file not found")
		fmt.Println("             Fix: set gitdir to point to a valid repo")
		fmt.Println("")
		configValid = false
	} else {
		repoVer := RepoCfg.GetString("version")
		if repoVer != configVersion {
			fmt.Println("")
			fmt.Println("  WARNING    meepctl version[" + configVersion + "] != repocfg version[" + repoVer + "]")
			fmt.Println("             repocfg file: " + RepoCfg.ConfigFileUsed())
			fmt.Println("             Fix: upgrade meepctl binary to matching version or set gitdir to repo with matching version")
			fmt.Println("")
			configValid = false
		}
	}

	return configValid
}

// ConfigPathValid validates IP address
func ConfigPathValid(path string) (valid bool, reason string) {
	valid = true
	fi, err := os.Stat(path)

	if err != nil {
		reason = "Path error  [" + path + "]"
		valid = false
	} else {
		if !fi.IsDir() {
			reason = "Not a directory [" + path + "]"
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
