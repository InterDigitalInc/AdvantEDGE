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
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ghodss/yaml"
	"github.com/roymx/viper"
	"github.com/spf13/cobra"

	ce "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-ctrl-engine-client"
)

// replayCmd represents the replay command
var replayCmd = &cobra.Command{
	Use:   "replay <action> [filename]",
	Short: "Use and manage auto-replay feature",
	Long: `AdvantEDGE supports creation and usage of auto-replay files.

Replay files are 'script-like' files that are maintained in AdvantEDGE document store.
Replay file contain sequence of events that can be automatically replayed following the specific time sequence of events.

Multiple actions can be performed on replay files.

Valid actions:
  * ls: gets a list of replay files name
  * cat: gets a replay file content
  * rm:  deletes one/all replay files
  * import: copies local yaml file to the replay store
  * export: copies a replay store file content into a local yaml file
  * generate: auto-creates a new replay file from scenario events
  * start: executes auto-replay file
  * stop: stops execution of an auto-replay`,

	Example: `  # Gets all the replay files name stored
  meepctl replay ls
  # Gets one replay file content
  meepctl replay cat <filename>
  # Deletes all replay files
  meepctl replay rm
  # Deletes one replay file
  meepctl replay rm <filename>
  # Creates a replay file based on YAML file
  meepctl replay import <yaml-file-name.yaml> <replay-file-name>
  # Creates a YAML file based on a replay file
  meepctl replay export <replay-file-name> <yaml-file-name.yaml>
  # Creates a replay file using latest events of a scenario
  meepctl replay generate <filename> <scenarioname>
  # Starts auto-replay
  meepctl replay start <filename>
  # Plays a replay file in an infinite loop
  meepctl replay start <filename> -l
  # Stops a replay file
  meepctl replay stop <filename>`,

	Args:      cobra.RangeArgs(1, 3),
	ValidArgs: []string{"ls", "cat", "rm", "cp", "export", "generate", "start", "stop"},
	Run: func(cmd *cobra.Command, args []string) {

		action := args[0]
		filename := ""
		arg2 := "" //may differ based on action
		nbArgs := len(args)
		if nbArgs > 1 {
			filename = args[1]
		}
		if nbArgs > 2 {
			arg2 = args[2]
		}

		v, _ := cmd.Flags().GetBool("verbose")
		l, _ := cmd.Flags().GetBool("loop")
		if v {
			fmt.Println("Replay called")
			fmt.Println("[arg]  action:", action)
			fmt.Println("[flag] verbose:", v)
			fmt.Println("[flag] loop:", l)
		}

		switch action {
		case "ls":
			replayGetAll(cmd)
		case "cat":
			if nbArgs != 2 {
				fmt.Println("This command requires 1 argument")
				return
			}
			replayGet(cmd, filename)
		case "rm":
			if nbArgs != 2 {
				fmt.Println("This command requires 1 argument")
				return
			}
			replayDelete(cmd, filename)
		case "import":
			if nbArgs != 3 {
				fmt.Println("This command requires 2 argument")
				return
			}
			replayAdd(cmd, filename, arg2)
		case "export":
			if nbArgs != 3 {
				fmt.Println("This command requires 2 argument")
				return
			}
			replayExport(cmd, filename, arg2)
		case "generate":
			if nbArgs != 3 {
				fmt.Println("This command requires 2 argument")
				return
			}
			replayAddFromScenario(cmd, filename, arg2)
		case "start":
			if nbArgs != 2 {
				fmt.Println("This command requires 1 argument")
				return
			}
			replayPlay(cmd, filename)
		case "stop":
			if nbArgs != 2 {
				fmt.Println("This command requires 1 argument")
				return
			}
			replayStop(cmd, filename)
		default:
			fmt.Println("Action ", action, ", not supported")
		}
	},
}

func init() {
	rootCmd.AddCommand(replayCmd)

	replayCmd.Flags().BoolP("loop", "l", false, "Enables replay files to loop indefinitely")
}

func createClient(path string) (*ce.APIClient, error) {
	// Create & store client for App REST API
	ceClientCfg := ce.NewConfiguration()
	ceClientCfg.BasePath = path
	ceClient := ce.NewAPIClient(ceClientCfg)
	if ceClient == nil {
		err := errors.New("Failed to create ctrl-engine REST API client")
		return nil, err
	}
	return ceClient, nil
}

func replayGet(cobraCmd *cobra.Command, filename string) {
	verbose, _ := cobraCmd.Flags().GetBool("verbose")

	client, err := createClient(getBasePath())
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
		fmt.Println(string(json))
	}
}

func replayGetAll(cobraCmd *cobra.Command) {
	verbose, _ := cobraCmd.Flags().GetBool("verbose")

	client, err := createClient(getBasePath())
	if err != nil {
		printError("Error creating client: ", err, verbose)
		return
	}

	replayFileNameList, _, err := client.EventReplayApi.GetReplayFileList(context.TODO())
	if err != nil {
		printError("Error: ", err, verbose)
	} else {
		json, err := json.Marshal(replayFileNameList)
		if err != nil {
			printError("Error: ", err, verbose)
		}
		fmt.Println(string(json))
	}
}

func replayDelete(cobraCmd *cobra.Command, filename string) {

	verbose, _ := cobraCmd.Flags().GetBool("verbose")

	client, err := createClient(getBasePath())
	if err != nil {
		printError("Error creating client: ", err, verbose)
		return
	}

	if filename == "" {
		_, err = client.EventReplayApi.DeleteReplayFileList(context.TODO())

	} else {
		_, err = client.EventReplayApi.DeleteReplayFile(context.TODO(), filename)
	}
	if err != nil {
		printError("Error: ", err, verbose)
	} else {
		if verbose {
			fmt.Println("Command successful")
		}
	}
}

func replayAdd(cobraCmd *cobra.Command, yamlFilename string, replayFilename string) {
	verbose, _ := cobraCmd.Flags().GetBool("verbose")
	client, err := createClient(getBasePath())
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

	var replay ce.Replay
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

func printError(errorString string, err error, verbose bool) {
	if verbose {
		fmt.Println(errorString, err)
	} else {
		fmt.Println("Command failed, use -v for more details")
	}
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

func replayAddFromScenario(cobraCmd *cobra.Command, filename string, scenarioName string) {
	verbose, _ := cobraCmd.Flags().GetBool("verbose")
	client, err := createClient(getBasePath())
	if err != nil {
		printError("Error creating client: ", err, verbose)
		return
	}

	var scenarioInfo ce.ScenarioInfo
	scenarioInfo.Name = scenarioName
	_, err = client.EventReplayApi.CreateReplayFileFromScenarioExec(context.TODO(), filename, scenarioInfo)

	if err != nil {
		printError("Error: ", err, verbose)
	} else {
		if verbose {
			fmt.Println("Command successful")
		}
	}
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

func replayStop(cobraCmd *cobra.Command, filename string) {
	verbose, _ := cobraCmd.Flags().GetBool("verbose")

	client, err := createClient(getBasePath())
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

func getBasePath() string {

	host := viper.GetString("node.ip")
	reqString := "http://" + host + ":30000/v1"
	return reqString
}
