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

	"github.com/russell/stately/pkg/stately/models"
	"go.uber.org/zap"
)

type ManifestOptions struct {
	StateFile       string
	InputFile       string
	OutputDirectory string
	Logger          *zap.SugaredLogger
}

func Manifest(o *ManifestOptions) error {
	manifests, err  := models.NewManifestContainerFromStdin()
	if err != nil {
		return err
	}

	for _, file := range manifests.Files {
		if err := file.ManifestFile(o.OutputDirectory); err != nil {
			return err
		}
	}

	return nil
}
