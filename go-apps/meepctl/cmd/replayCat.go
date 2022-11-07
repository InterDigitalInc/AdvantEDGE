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

	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
)

// replayCatCmd represents the replay cat command
var replayCatCmd = &cobra.Command{
	Use:     "cat <filename>",
	Short:   "Prints a replay file content",
	Long:    "Prints a replay file content.",
	Args:    cobra.ExactValidArgs(1),
	Example: "meepctl replay cat myfilename",
	Run: func(cmd *cobra.Command, args []string) {

		v, _ := cmd.Flags().GetBool("verbose")
		if v {
			fmt.Println("Replay cat called")
			fmt.Println("[flag] verbose:", v)
		}
		replayGet(cmd, args[0])
	},
}

func init() {
	setSandboxFlag(replayCatCmd)
	replayCmd.AddCommand(replayCatCmd)
}

func replayGet(cobraCmd *cobra.Command, filename string) {
	verbose, _ := cobraCmd.Flags().GetBool("verbose")

	client, err := createClient(getBasePath(cobraCmd))
	if err != nil {
		printError("Error creating client: ", err, verbose)
		return
	}

	replay, _, err := client.EventReplayApi.GetReplayFile(context.TODO(), filename)
	if err != nil {
		printError("Error: ", err, verbose)
	} else {
		json, err := json.Marshal(replay)
		if err != nil {
			printError("Error: ", err, verbose)
		}

		jsonToYaml, err := yaml.JSONToYAML(json)
		if err != nil {
			printError("Error converting JSON to YAML: ", err, verbose)
			return
		}

		fmt.Println(string(jsonToYaml))
	}
}
