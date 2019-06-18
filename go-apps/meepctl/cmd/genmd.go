/*
 * Copyright (c) 2019
 * InterDigital Communications, Inc.
 * All rights reserved.
 *
 * The information provided herein is the proprietary and confidential
 * information of InterDigital Communications, Inc.
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
