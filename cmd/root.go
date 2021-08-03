/*
Copyright © 2021 Manuel Kieweg <mail@manuelkieweg.de>
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice,
   this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice,
   this list of conditions and the following disclaimer in the documentation
   and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its contributors
   may be used to endorse or promote products derived from this software
   without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
POSSIBILITY OF SUCH DAMAGE.
*/
package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/mkieweg/syncdeps"
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string
var baseline string
var target string
var stdout bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "syncdeps",
	Short: "A brief description of your application",
	RunE: func(cmd *cobra.Command, args []string) error {
		baselineFile, err := os.Open(baseline)
		if err != nil {
			return err
		}
		targetFile, err := os.Open(target)
		if err != nil {
			return err
		}
		baselineDeps, err := syncdeps.ScanFile(baselineFile)
		if err != nil {
			return err
		}
		targetDeps, err := syncdeps.ScanFile(targetFile)
		if err != nil {
			return err
		}
		deps := syncdeps.Compare(baselineDeps, targetDeps)
		var w io.Writer
		if stdout {
			w = os.Stdout
		} else {
			w, err = os.OpenFile("go.mod", os.O_APPEND, os.ModeAppend)
			if err != nil {
				return err
			}
		}

		return syncdeps.WriteGoMod(w, deps)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.syncdeps.yaml)")
	rootCmd.Flags().StringVar(&baseline, "baseline", "", "path to the baseline go.sum")
	rootCmd.Flags().StringVar(&target, "target", "", "path to the target go.sum")
	rootCmd.Flags().BoolVar(&stdout, "stdout", true, "print output to stdout")
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

		// Search config in home directory with name ".syncdeps" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".syncdeps")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
