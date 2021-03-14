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
package actions

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	dircopy "github.com/otiai10/copy"
	"go.uber.org/zap"

	"github.com/russell/stately/pkg/stately/config"
)

type CopyOptions struct {
	SourcePaths     []string
	StateFile       string
	OutputDirectory string
	Logger          *zap.SugaredLogger
}

func Copy(o *CopyOptions) (error) {
	var currentState config.StateConfig

	newState := config.NewStateConfig()
	newState.OutputDirectory = o.OutputDirectory

	if _, err := os.Stat(o.StateFile); err == nil {
		currentState, _ = config.NewStateConfigFromFile(o.StateFile)

		// Check that output directory matches
		if currentState.OutputDirectory != newState.OutputDirectory{
			return fmt.Errorf("Output directory in state file doesn't match argument '%s' != %s",
				currentState.OutputDirectory,
				newState.OutputDirectory)
		}
	}

	var newFiles []config.StateFile

	var src os.FileInfo
	var srcPath string
	opt := dircopy.Options{
		Skip: func(s string) (bool, error) {
			// Skip recording directories
			if info, _ := os.Lstat(s); info.IsDir() {
				return false, nil
			}

			o.Logger.Debugf("Copying: %s", s)

			if src.IsDir() {
				s = strings.TrimPrefix(s, srcPath)
			}
			if s[0] == '/' {
				s = s[1:]
			}
			newFiles = append(newFiles, config.StateFile{Path: s})
			return false, nil
		},
		PreserveTimes: true,
	}

	for _, s := range o.SourcePaths {
		var dest string
		src, _ = os.Lstat(s)
		srcPath = s
		if src.IsDir() {
			dest = o.OutputDirectory
		} else {
			dest = filepath.Join(o.OutputDirectory, path.Base(s))
			// dircopy, won't call the skip method on single file copies.
			newFiles = append(newFiles, config.StateFile{Path: path.Base(s)})
			o.Logger.Debugf("Copying: %s", path.Base(s))
		}
		err := dircopy.Copy(s, dest, opt)
		if err != nil {
			o.Logger.Errorw(fmt.Sprint(err))
		}
	}

	newState.Files = newFiles

	// Calculate what should be deleted
	toDelete := make(map[string]bool)
	for _, s := range currentState.Files {
		toDelete[s.Path] = true
	}

	for _, s := range newState.Files {
		toDelete[s.Path] = false
	}

	for file, delete := range toDelete {
		if delete == false {
			continue
		}
		filePath := filepath.Join(o.OutputDirectory, file)
		o.Logger.Debugf("Deleting: %s", filePath)
		if err := os.Remove(filePath); err != nil {
			o.Logger.Infof("Couldn't delete: %s", filePath)
		}
	}

	newState.WriteToFile(o.StateFile)
	return nil
}
