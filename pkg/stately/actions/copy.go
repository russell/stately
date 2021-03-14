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
	"io"
	"os"
	"path/filepath"

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

	if _, err := os.Stat(o.StateFile); err == nil {
		currentState, _ = config.NewStateConfigFromFile(o.StateFile)
	}

	var newFiles []config.StateFile

	var src os.FileInfo

	for _, s := range o.SourcePaths {
		var dest string
		src, _ = os.Lstat(s)
		if src.IsDir() {
			return fmt.Errorf("ERROR: Only files are supported: %s", s)
		} else {
			dest = filepath.Join(o.OutputDirectory, s)
			// dircopy, won't call the skip method on single file copies.
			newFiles = append(newFiles, config.StateFile{Path: s})
			o.CopyFile(s, dest)
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

func (o *CopyOptions) CopyFile(src string, dest string) (err error) {
	var destination *os.File
	var source *os.File
	var info os.FileInfo

	if info, err = os.Lstat(src); err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(dest), os.ModePerm); err != nil {
		return err
	}

	o.Logger.Debugf("Copying: %s", src)

	destination, err = os.Create(dest)
	if err != nil {
		return err
	}

	if err := os.Chmod(destination.Name(), info.Mode()|0200); err != nil {
		return err
	}

	source, err = os.Open(src)
	if err != nil {
		return err
	}

	if _, err := io.Copy(destination, source); err != nil {
		return err
	}

	return nil
}
