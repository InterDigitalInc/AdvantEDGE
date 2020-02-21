// Copyright © 2020 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
)

// replayExportCmd represents the replay export command
var replayExportCmd = &cobra.Command{
	Use:     "export <replay-file-name> <yaml-file-name.yaml>",
	Short:   "Copies a replay store file content into a local yaml file",
	Long:    "Copies a replay store file content into a local yaml file.",
	Args:    cobra.ExactValidArgs(2),
	Example: "meepctl replay export replayfilename yamlfilename.yaml",
	Run: func(cmd *cobra.Command, args []string) {

		v, _ := cmd.Flags().GetBool("verbose")
		if v {
			fmt.Println("Replay export called")
			fmt.Println("[flag] verbose:", v)
		}
		replayExport(cmd, args[0], args[1])
	},
}

func init() {
	replayCmd.AddCommand(replayExportCmd)
}

func replayExport(cobraCmd *cobra.Command, replayFilename string, yamlFilename string) {
	verbose, _ := cobraCmd.Flags().GetBool("verbose")
	client, err := createClient(getBasePath())
	if err != nil {
		printError("Error creating client: ", err, verbose)
		return
	}

	replay, _, err := client.EventReplayApi.GetReplayFile(context.TODO(), replayFilename)
	if err != nil {
		printError("Error getting replay file: ", err, verbose)
		return
	}

	json, err := json.Marshal(replay)
	if err != nil {
		printError("Error creating JSON: ", err, verbose)
		return
	}

	jsonToYaml, err := yaml.JSONToYAML(json)
	if err != nil {
		printError("Error converting JSON to YAML: ", err, verbose)
		return
	}

	err = ioutil.WriteFile(yamlFilename, jsonToYaml, 0644)
	if err != nil {
		printError("Error creating yaml file: ", err, verbose)
		return
	}

	if verbose {
		fmt.Println("Command successful")
	}
}
