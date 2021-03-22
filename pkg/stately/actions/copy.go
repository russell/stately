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
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofrs/flock"
	"github.com/russell/stately/pkg/stately/config"
	"go.uber.org/zap"
)

type CopyOptions struct {
	SourcePaths     []string
	StateFile       string
	StripPrefix     string
	TargetName      string
	FollowSymlinks  bool
	OutputDirectory string
	FileMode        os.FileMode
	Logger          *zap.SugaredLogger
}

func Copy(o *CopyOptions) error {
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

	cb := func(s string, d string) {
		dest, _ := filepath.Abs(d)
		rel, _ := filepath.Rel(stateFileDir, dest)
		newFiles = append(newFiles, config.StateFile{Path: rel})
	}

	for _, s := range o.SourcePaths {
		var dest string
		dest = filepath.Join(o.OutputDirectory, strings.TrimPrefix(s, o.StripPrefix))
		// dircopy, won't call the skip method on single file copies.
		if err := o.Copy(s, dest, cb); err != nil {
			return err
		}
	}

	newState := config.NewStateConfig()
	newState.Targets[o.TargetName] = config.StateTarget{Files: newFiles}
	newState.WriteToFile(o.StateFile)
	config.Cleanup(stateFile, o.TargetName, currentState, newState, o.Logger)
	return nil
}

func (o *CopyOptions) Copy(src string, dest string, cb func(string, string)) (err error) {
	stat, err := os.Lstat(src)
	if err != nil {
		return fmt.Errorf("ERROR: File doesn't exist: %s", src)
	} else if stat.IsDir() {
		return o.CopyDirectory(src, dest, cb)
	} else if stat.Mode()&os.ModeSymlink != 0 {
		return o.CopySymlink(src, dest, cb)
	} else if stat.Mode()&os.ModeNamedPipe != 0 {
		return fmt.Errorf("ERROR: NamedPipes aren't supported: %s", src)
	} else {
		cb(src, dest)
		return o.CopyFile(src, dest)
	}
}

func (o *CopyOptions) CopyDirectory(src string, dest string, cb func(string, string)) (err error) {
	files, err := ioutil.ReadDir(src)

	for _, f := range files {
		if err := o.Copy(filepath.Join(src, f.Name()), filepath.Join(dest, f.Name()), cb); err != nil {
			return err
		}
	}
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

	o.Logger.Debugf("Copying: %s -> %s", src, dest)

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

	if o.FileMode.Perm() != os.FileMode(0).Perm() {
		if err := os.Chmod(destination.Name(), o.FileMode|0200); err != nil {
			return err
		}
	}

	return nil
}

func (o *CopyOptions) CopySymlink(src string, dest string, cb func(string, string)) (err error) {
	if o.FollowSymlinks == true {
		target, err := os.Readlink(src)
		if err != nil {
			return err
		}
		// If the end of the symlink isn't a
		stat, err := os.Lstat(target)
		if err != nil {
			return fmt.Errorf("ERROR: File doesn't exist: %s", src)
		} else if stat.IsDir() {
			return o.CopyDirectory(src, dest, cb)
		} else if stat.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("ERROR: Multiple symlinks aren't supported: %s", src)
		} else if stat.Mode()&os.ModeNamedPipe != 0 {
			return fmt.Errorf("ERROR: NamedPipes aren't supported: %s", src)
		} else {
			cb(src, dest)
			return o.CopyFile(target, dest)
		}
	}
	return fmt.Errorf("ERROR: Symlinks aren't supported: %s", src)

}

func (o *CopyOptions) LockFile(fileLock *flock.Flock) (bool, error) {
	locked, err := fileLock.TryLock()
	if err != nil {
		o.Logger.Debugf("Trying to acquire lock of %s", o.StateFile+".lock")
		time.Sleep(100 * time.Millisecond)
		return o.LockFile(fileLock)
	}
	return locked, nil
}
