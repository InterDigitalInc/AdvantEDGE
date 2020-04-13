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
	"os"

	"github.com/InterDigitalInc/AdvantEDGE/go-apps/meepctl/utils"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/roymx/viper"
	"github.com/spf13/cobra"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "meepctl",
	Short: "meepctl - CLI application to control the AdvantEDGE platform",
	Long: `CLI application to control the AdvantEDGE platform
Find more information [here](https://github.com/InterDigitalInc/AdvantEDGE/blob/master/docs/meepctl/meepctl.md)`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	// Make sure Config is initialized
	if !utils.CfgInitialized {
		fmt.Println("Failed to initialize configuration")
		os.Exit(1)
	}

	// Run command
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
