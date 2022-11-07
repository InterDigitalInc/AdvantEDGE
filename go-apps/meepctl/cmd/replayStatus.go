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

// replayStatusCmd represents the replay status command
var replayStatusCmd = &cobra.Command{
	Use:     "status",
	Short:   "Retrieve replay status",
	Long:    "Retrieve replay execution status from backend.",
	Example: "meepctl replay status",
	Run: func(cmd *cobra.Command, args []string) {
		v, _ := cmd.Flags().GetBool("verbose")
		if v {
			fmt.Println("Replay status called")
			fmt.Println("[flag] verbose:", v)
		}
		replayStatus(cmd)
	},
}

func init() {
	setSandboxFlag(replayStatusCmd)
	replayCmd.AddCommand(replayStatusCmd)
}

func replayStatus(cobraCmd *cobra.Command) {
	verbose, _ := cobraCmd.Flags().GetBool("verbose")

	client, err := createClient(getBasePath(cobraCmd))
	if err != nil {
		printError("Error creating client: ", err, verbose)
		return
	}

	status, _, err := client.EventReplayApi.GetReplayStatus(context.TODO())
	if err != nil {
		if err.Error() == "404 Not Found" {
			fmt.Println("Replay file not running...")
		} else {
			printError("Error: ", err, verbose)
		}
	} else {
		json, err := json.Marshal(status)
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
