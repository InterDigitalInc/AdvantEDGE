/*
 * Copyright (c) 2019  InterDigital Communications, Inc
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
	"fmt"
	"time"

	"github.com/InterDigitalInc/AdvantEDGE/go-apps/meepctl/utils"
	"github.com/spf13/cobra"
)

// configKibana represents the configKibana command
var configKibana = &cobra.Command{
	Use:   "kibana",
	Short: "Configures Kibana (index pattern, saved objects such as dashboards, visualisations, etc.)",
	Long: `Configures Kibana (index pattern, saved objects such as dashboards, visualisations, etc.)
Any Kibana saved object will be overwritten in the process if the object Id are the same 
meepctl config kibana.`,
	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if !utils.ConfigValidate("") {
			fmt.Println("Fix configuration issues")
			return
		}

		t, _ := cmd.Flags().GetBool("time")
		start := time.Now()
		utils.DeployKibanaDashboards(cmd)

		elapsed := time.Since(start)
		if t {
			fmt.Println("Took ", elapsed.Round(time.Millisecond).String())
		}
	},
}

func init() {
	configCmd.AddCommand(configKibana)
}
