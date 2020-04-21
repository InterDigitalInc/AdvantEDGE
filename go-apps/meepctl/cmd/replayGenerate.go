/*
 * Copyright (c) 2020  InterDigital Communications, Inc
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
	"fmt"

	sandbox "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-sandbox-ctrl-client"
	"github.com/spf13/cobra"
)

// replayGenerateCmd represents the replay generate command
var replayGenerateCmd = &cobra.Command{
	Use:     `generate <filename> <scenarioname> ["description"]`,
	Short:   "Creates a new replay file from scenario events",
	Long:    "Creates a new replay file from scenario events.",
	Args:    cobra.RangeArgs(2, 3),
	Example: `meepctl replay generate <filename> <scenarioname> ["description-of-the-content (string)"]`,
	Run: func(cmd *cobra.Command, args []string) {

		v, _ := cmd.Flags().GetBool("verbose")
		if v {
			fmt.Println("Replay generate called")
			fmt.Println("[flag] verbose:", v)
		}

		desc := ""
		if len(args) == 3 {
			desc = args[2]
		}
		replayAddFromScenario(cmd, args[0], args[1], desc)
	},
}

func init() {
	setSandboxFlag(replayGenerateCmd)
	replayCmd.AddCommand(replayGenerateCmd)
}

func replayAddFromScenario(cobraCmd *cobra.Command, filename string, scenarioName string, description string) {
	verbose, _ := cobraCmd.Flags().GetBool("verbose")
	client, err := createClient(getBasePath(cobraCmd))
	if err != nil {
		printError("Error creating client: ", err, verbose)
		return
	}

	var replayInfo sandbox.ReplayInfo
	replayInfo.ScenarioName = scenarioName
	replayInfo.Description = description

	_, err = client.EventReplayApi.CreateReplayFileFromScenarioExec(context.TODO(), filename, replayInfo)

	if err != nil {
		printError("Error: ", err, verbose)
	} else {
		if verbose {
			fmt.Println("Command successful")
		}
	}
}
