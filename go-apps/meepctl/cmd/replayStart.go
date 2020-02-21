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
	"fmt"

	"github.com/spf13/cobra"
)

// replayStartCmd represents the replay start command
var replayStartCmd = &cobra.Command{
	Use:     "start <filename>",
	Short:   "Executes auto-replay file",
	Long:    "Executes auto-replay file.",
	Args:    cobra.ExactValidArgs(1),
	Example: "meepctl replay start myfilename",
	Run: func(cmd *cobra.Command, args []string) {

		v, _ := cmd.Flags().GetBool("verbose")
		l, _ := cmd.Flags().GetBool("loop")
		if v {
			fmt.Println("Replay start called")
			fmt.Println("[flag] verbose:", v)
			fmt.Println("[flag] loop:", l)
		}

		replayPlay(cmd, args[0])
	},
}

func init() {
	replayCmd.AddCommand(replayStartCmd)

	replayStartCmd.Flags().BoolP("loop", "l", false, "Enables replay files to loop indefinitely")
}

func replayPlay(cobraCmd *cobra.Command, filename string) {
	loop, _ := cobraCmd.Flags().GetBool("loop")
	verbose, _ := cobraCmd.Flags().GetBool("verbose")

	client, err := createClient(getBasePath())
	if err != nil {
		printError("Error creating client: ", err, verbose)
		return
	}

	if loop {
		_, err = client.EventReplayApi.LoopReplay(context.TODO(), filename)
	} else {
		_, err = client.EventReplayApi.PlayReplayFile(context.TODO(), filename)
	}

	if err != nil {
		printError("Error: ", err, verbose)
	} else {
		if verbose {
			fmt.Println("Command successful")
		}
	}
}
