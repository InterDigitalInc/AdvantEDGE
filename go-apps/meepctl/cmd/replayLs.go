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

// replayLsCmd represents the replay ls command
var replayLsCmd = &cobra.Command{
	Use:     "ls",
	Short:   "Gets a list of replay files name",
	Long:    "Gets a list of replay files name.",
	Example: "meepctl replay ls",
	Run: func(cmd *cobra.Command, args []string) {

		v, _ := cmd.Flags().GetBool("verbose")
		l, _ := cmd.Flags().GetBool("long")
		if v {
			fmt.Println("Replay ls called")
			fmt.Println("[flag] verbose:", v)
			fmt.Println("[flag] long:", l)
		}
		replayGetAll(cmd)

	},
}

func init() {
	replayCmd.AddCommand(replayLsCmd)

	replayLsCmd.Flags().BoolP("long", "l", false, "Displays description of each file")
}

func replayGetAll(cobraCmd *cobra.Command) {
	verbose, _ := cobraCmd.Flags().GetBool("verbose")
	long, _ := cobraCmd.Flags().GetBool("long")

	client, err := createClient(getBasePath())
	if err != nil {
		printError("Error creating client: ", err, verbose)
		return
	}

	replayFileNameList, _, err := client.EventReplayApi.GetReplayFileList(context.TODO())
	if err != nil {
		printError("Error: ", err, verbose)
	} else {
		for _, replayFilename := range replayFileNameList.ReplayFiles {
			replay, _, err := client.EventReplayApi.GetReplayFile(context.TODO(), replayFilename)
			if err != nil {
				printError("Error getting replay file: ", err, verbose)
				return
			}

			if long {
				fmt.Println(replayFilename, " : ", replay.Description)
			} else {
				fmt.Println(replayFilename)
			}
		}
	}
}
