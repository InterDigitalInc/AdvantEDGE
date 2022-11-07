/*
 * Copyright (c) 2022  The AdvantEDGE Authors
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

package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"

	sandbox "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-sandbox-ctrl-client"
)

// replayImportCmd represents the replay import command
var replayImportCmd = &cobra.Command{
	Use:     "import <yaml-file-name.yaml> <replay-file-name>",
	Short:   "Copies local yaml file to the replay store",
	Long:    "Copies local yaml file to the replay store.",
	Args:    cobra.ExactValidArgs(2),
	Example: "meepctl replay import yamlfilename.yaml replayfilename",
	Run: func(cmd *cobra.Command, args []string) {

		v, _ := cmd.Flags().GetBool("verbose")
		if v {
			fmt.Println("Replay import called")
			fmt.Println("[flag] verbose:", v)
		}
		replayAdd(cmd, args[0], args[1])
	},
}

func init() {
	setSandboxFlag(replayImportCmd)
	replayCmd.AddCommand(replayImportCmd)
}

func replayAdd(cobraCmd *cobra.Command, yamlFilename string, replayFilename string) {
	verbose, _ := cobraCmd.Flags().GetBool("verbose")
	client, err := createClient(getBasePath(cobraCmd))
	if err != nil {
		printError("Error creating client: ", err, verbose)
		return
	}

	file, err := os.Open(yamlFilename)
	if err != nil {
		printError("Error opening file: ", err, verbose)
		return
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		printError("Error reading file: ", err, verbose)
		return
	}

	yamltojson, err := yaml.YAMLToJSON(b)
	if err != nil {
		printError("Error converting YAML to JSON: ", err, verbose)
		return
	}

	var replay sandbox.Replay
	err = json.Unmarshal(yamltojson, &replay)
	if err != nil {
		printError("Error decoding JSON: ", err, verbose)
		return
	}

	_, err = client.EventReplayApi.CreateReplayFile(context.TODO(), replayFilename, replay)

	if err != nil {
		printError("Error: ", err, verbose)
	} else {
		if verbose {
			fmt.Println("Command successful")
		}
	}
}
