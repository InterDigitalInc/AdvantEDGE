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

// replayRmCmd represents the replay rm command
var replayRmCmd = &cobra.Command{
	Use:   "rm [filename]",
	Short: "Deletes one/all replay files",
	Long:  "Deletes one/all replay files.",
	Args:  cobra.RangeArgs(0, 1),
	Example: `  # Deletes all the replay files stored
  meepctl replay rm -a
  # Deletes one replay file
  meepctl replay rm <filename>`,
	Run: func(cmd *cobra.Command, args []string) {

		v, _ := cmd.Flags().GetBool("verbose")
		a, _ := cmd.Flags().GetBool("all")
		if v {
			fmt.Println("Replay rm called")
			fmt.Println("[flag] verbose:", v)
			fmt.Println("[flag] all:", a)
		}
		filename := ""
		if len(args) == 1 {
			filename = args[0]
		} else { //no args
			if !a {
				_ = cmd.Help()
				return
			}
		}
		replayDelete(cmd, filename)
	},
}

func init() {
	setSandboxFlag(replayRmCmd)
	replayRmCmd.Flags().BoolP("all", "a", false, "Removes all replay files")
	replayCmd.AddCommand(replayRmCmd)
}

func replayDelete(cobraCmd *cobra.Command, filename string) {

	verbose, _ := cobraCmd.Flags().GetBool("verbose")
	all, _ := cobraCmd.Flags().GetBool("all")

	client, err := createClient(getBasePath(cobraCmd))
	if err != nil {
		printError("Error creating client: ", err, verbose)
		return
	}

	if filename == "" && all {
		_, err = client.EventReplayApi.DeleteReplayFileList(context.TODO())

	} else {
		_, err = client.EventReplayApi.DeleteReplayFile(context.TODO(), filename)
		if all {
			if verbose {
				fmt.Println("[flag] used ignored")
			}
		}
	}
	if err != nil {
		printError("Error: ", err, verbose)
	} else {
		if verbose {
			fmt.Println("Command successful")
		}
	}
}
