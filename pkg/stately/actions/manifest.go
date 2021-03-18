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
	"github.com/gofrs/flock"
	"github.com/russell/stately/pkg/stately/config"
	"github.com/russell/stately/pkg/stately/models"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"time"
)

type ManifestOptions struct {
	StateFile       string
	InputFile       string
	TargetName      string
	OutputDirectory string
	Logger          *zap.SugaredLogger
}

func Manifest(o *ManifestOptions) error {
	var currentState config.StateConfig

	fileLock := flock.New(o.StateFile)
	o.LockFile(fileLock)
	defer fileLock.Unlock()

	if _, err := os.Stat(o.StateFile); err == nil {
		currentState, _ = config.NewStateConfigFromFile(o.StateFile)
	}

	var newFiles []config.StateFile
	stateFile, _ := filepath.Abs(o.StateFile)
	stateFileDir := filepath.Dir(stateFile)

	// Manifest files
	manifests, err := models.NewManifestContainerFromStdin()
	if err != nil {
		return err
	}

	for _, file := range manifests.Files {
		// Don't install none files
		if file.Install == models.None {
			continue
		}

		dest, err := file.ManifestFile(o.OutputDirectory, o.Logger)
		if err != nil {
			o.Logger.Errorf("%s", err)
		}

		o.Logger.Debugf("Manifesting file: %s", dest)
		dest, _ = filepath.Abs(dest)
		rel, _ := filepath.Rel(stateFileDir, dest)
		newFiles = append(newFiles, config.StateFile{Path: rel})
	}

	newState := config.NewStateConfig()
	newState.Targets = currentState.Targets
	if newState.Targets == nil {
		newState.Targets = make(map[string]config.StateTarget)
	}
	newState.Targets[o.TargetName] = config.StateTarget{Files: newFiles}
	newState.WriteToFile(o.StateFile)
	config.Cleanup(stateFile, o.TargetName, currentState, newState, o.Logger)
	return nil
}

func (o *ManifestOptions) LockFile(fileLock *flock.Flock) (bool, error) {
	locked, err := fileLock.TryLock()
	if err != nil {
		o.Logger.Debugf("Trying to acquire lock of %s", o.StateFile+".lock")
		time.Sleep(100 * time.Millisecond)
		return o.LockFile(fileLock)
	}
	return locked, nil
}
