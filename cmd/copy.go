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

// copyCmd represents the copy command
var copyCmd = &cobra.Command{
	Use:   "copy",
	Short: "Copy files into an output directory",
	Long: `Copy files from one directory to another`,
	Run: func(cmd *cobra.Command, args []string) {
		logger, _ := zap.NewDevelopment()
		sugar := logger.Sugar()
		stateFile, _ := cmd.Flags().GetString("state-file")
		outputDir, _ := cmd.Flags().GetString("output-dir")
		name, _ := cmd.Flags().GetString("name")
		followSymlinks, _ := cmd.Flags().GetBool("follow-symlinks")
		stripPrefix, _ := cmd.Flags().GetString("strip-prefix")
		options := actions.CopyOptions{
			SourcePaths:     args,
			StateFile:       stateFile,
			StripPrefix:     stripPrefix,
			FollowSymlinks:  followSymlinks,
			TargetName:      name,
			OutputDirectory: outputDir,
			Logger:          sugar,
		}
		err := actions.Copy(&options)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(copyCmd)

	copyCmd.Flags().StringP("state-file", "s", ".stately-files.yaml", "The state file to use")
	copyCmd.Flags().StringP("output-dir", "o", "", "The location to copy to")
	copyCmd.Flags().StringP("name", "n", "default", "The name of the file set to track")
	copyCmd.Flags().StringP("strip-prefix", "", "", "Remove the prefix from output paths")
	copyCmd.Flags().BoolP("follow-symlinks", "L", false, "Copy the files instead of their symlinks")
	copyCmd.MarkFlagRequired("output-dir")
}
