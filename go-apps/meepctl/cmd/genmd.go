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
	"fmt"
	"log"
	"os"

	"github.com/roymx/viper"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

// genmdCmd represents the genmd command
var genmdCmd = &cobra.Command{
	Use:   "genmd",
	Short: "Generate markdown files for meepctl",
	Long:  `Generate markdown files for meepctl`,
	Run: func(cmd *cobra.Command, args []string) {
		outDir := viper.GetString("meep.gitdir") + "/docs/meepctl"
		if _, err := os.Stat(outDir); os.IsNotExist(err) {
			// default outdir not found ... use /tmp
			outDir = "/tmp"
		}
		err := doc.GenMarkdownTree(rootCmd, outDir)
		if err != nil {
			log.Fatal(err)
			return
		}
		fmt.Println("Markdown files can be found in ", outDir, " folder")
	},
}

func init() {
	rootCmd.AddCommand(genmdCmd)
}
