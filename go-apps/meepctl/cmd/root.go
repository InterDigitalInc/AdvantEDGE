// Copyright © 2019 InterDigital, Inc
// This file is part of {{ .appName }}.

package cmd

import (
	"fmt"
	"github.com/InterDigitalInc/AdvantEDGE/go-apps/meepctl/utils"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/roymx/viper"
	"github.com/spf13/cobra"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "meepctl",
	Short: "meepctl - CLI application to control the AdvantEDGE platform",
	Long: `
meepctl - CLI application to control the AdvantEDGE platform

  Find more information [here](https://kopsvas19p.interdigital.com/wbu-tep/AdvantEDGE/blob/develop/docs/meepctl/meepctl.md)
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	_, err := utils.ConfigCreateIfNotExist()
	if err != nil {
		fmt.Println("Error accessing config file at $(HOME)/.meepctl.yaml")
		fmt.Println(err)
		os.Exit(1)
	}
	utils.ConfigValidate("")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	//rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.meepctl.yaml)")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Display debug information")
	rootCmd.PersistentFlags().BoolP("time", "t", false, "Display timing information")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("verbose", "v", false, "Display debug information")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".meepctl" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".meepctl")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using meepctl config file:", viper.ConfigFileUsed())
	}
}
