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
	"sort"

	dircopy "github.com/otiai10/copy"
	"go.uber.org/zap"

	"github.com/russell/stately/pkg/stately/config"
)

type CopyOptions struct {
	SourcePaths     []string
	StateFile       string
	OutputDirectory string
	Logger *zap.SugaredLogger
}

func Copy(o *CopyOptions) {

	state := config.NewStateConfig()
	state.OutputDirectory = o.OutputDirectory

	var newFiles []config.StateFile

	opt := dircopy.Options{
		Skip: func(src string) (bool, error) {
			newFiles = append(newFiles, config.StateFile{ Path: src })
			return false, nil
		},
		PreserveTimes: true,
	}

	for _, s := range o.SourcePaths {
		err := dircopy.Copy(s, o.OutputDirectory, opt)
		if err != nil {
			o.Logger.Infow(fmt.Sprint(err))
		}
	}

	state.Files = newFiles
	state.WriteToFile(o.StateFile)
}
