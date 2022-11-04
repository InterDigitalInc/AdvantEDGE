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
	"fmt"

	"github.com/spf13/cobra"
)

// replayStopCmd represents the replay stop command
var replayStopCmd = &cobra.Command{
	Use:     "stop <filename>",
	Short:   "Stops execution of an auto-replay",
	Long:    "Stops execution of an auto-replay.",
	Args:    cobra.ExactValidArgs(1),
	Example: "meepctl replay stop myfilename",
	Run: func(cmd *cobra.Command, args []string) {

		v, _ := cmd.Flags().GetBool("verbose")
		if v {
			fmt.Println("Replay stop called")
			fmt.Println("[flag] verbose:", v)
		}

		replayStop(cmd, args[0])
	},
}

func init() {
	setSandboxFlag(replayStopCmd)
	replayCmd.AddCommand(replayStopCmd)
}

func replayStop(cobraCmd *cobra.Command, filename string) {
	verbose, _ := cobraCmd.Flags().GetBool("verbose")

	client, err := createClient(getBasePath(cobraCmd))
	if err != nil {
		printError("Error creating client: ", err, verbose)
		return
	}

	_, err = client.EventReplayApi.StopReplayFile(context.TODO(), filename)
	if err != nil {
		printError("Error: ", err, verbose)
	} else {
		if verbose {
			fmt.Println("Command successful")
		}
	}
}
