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
	Short: "Convert a JSON object into some files",
	Long: `Takes a JSON object format and converts this
into files on disk.
`,
	Run: func(cmd *cobra.Command, args []string) {
		logger, _ := zap.NewDevelopment()
		sugar := logger.Sugar()
		stateFile, _ := cmd.Flags().GetString("state-file")
		name, _ := cmd.Flags().GetString("name")
		outputDir, _ := cmd.Flags().GetString("output-dir")
		input, _ := cmd.Flags().GetString("input")
		options := actions.ManifestOptions{
			StateFile:       stateFile,
			InputFile:       input,
			TargetName:      name,
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

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	manifestCmd.Flags().StringP("state-file", "s", ".stately-files.yaml", "The state file to use")
	manifestCmd.Flags().StringP("output-dir", "o", cwd, "The location to copy to")
	manifestCmd.Flags().StringP("name", "n", "default", "The name of the file set to track")
	manifestCmd.Flags().StringP("input", "", "", "The input file or - for stdin")
}
