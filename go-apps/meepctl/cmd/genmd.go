// Copyright © 2019 InterDigital, Inc
// This file is part of {{ .appName }}.

package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

// genmdCmd represents the genmd command
var genmdCmd = &cobra.Command{
	Use:   "genmd",
	Short: "Generate markdown files for meepctl",
	Long:  `Generate markdown files for meepctl`,
	Run: func(cmd *cobra.Command, args []string) {
		err := doc.GenMarkdownTree(rootCmd, "/tmp")
		if err != nil {
			log.Fatal(err)
			return
		}
		fmt.Println("Markdown files can be found in /tmp folder")
	},
}

func init() {
	rootCmd.AddCommand(genmdCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// genmdCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// genmdCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
