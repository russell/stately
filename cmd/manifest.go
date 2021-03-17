/*
Copyright © 2021 Russell Sim <russell.sim@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/russell/stately/pkg/stately/actions"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// manifestCmd represents the manifest command
var manifestCmd = &cobra.Command{
	Use:   "manifest",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger, _ := zap.NewDevelopment()
		sugar := logger.Sugar()
		stateFile, _ := cmd.Flags().GetString("state-file")
		outputDir, _ := cmd.Flags().GetString("output-dir")
		input, _ := cmd.Flags().GetString("input")
		options := actions.ManifestOptions{
			StateFile:       stateFile,
			InputFile:       input,
			OutputDirectory: outputDir,
			Logger:          sugar,
		}
		err := actions.Manifest(&options)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(manifestCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// manifestCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// manifestCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	manifestCmd.Flags().StringP("state-file", "s", ".stately-files.yaml", "The state file to use")
	manifestCmd.Flags().StringP("output-dir", "o", "", "The location to copy to")
	manifestCmd.Flags().StringP("input", "", "", "The input file or - for stdin")
	manifestCmd.MarkFlagRequired("output-dir")
}