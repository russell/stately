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
package models

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

type FormatType string

const (
	Yaml FormatType = "yaml"
	Json            = "json"
	Raw             = "raw"
)

type InstallerType string

const (
	Symlink InstallerType = "symlink"
	Write                 = "write"
	None                  = "none"
)

type ManifestFileHeader struct {
	NoHeader    bool   `json:"-"`
	Prefix      string `json:"prefix"`
	LinesPrefix string `json:"linePrefix"`
	Suffix      string `json:"suffix"`
}

type ManifestFile struct {
	Path         string             `json:"-"`
	Install      InstallerType      `json:"install"`
	HeaderLines  []string           `json:"headerLines"`
	HeaderFormat ManifestFileHeader `json:"headerFormat"`
	Format       FormatType         `json:"format"`
	Contents     string             `json:"contents"`
	Executable   bool               `json:"executable"`
}

func (f *ManifestFile) ManifestFile(destination string) (err error) {
	dest := filepath.Join(destination, f.Path)
	if err := os.MkdirAll(filepath.Dir(dest), os.ModePerm); err != nil {
		return err
	}

	if err := ioutil.WriteFile(dest, []byte(f.Contents), 0644); err != nil {
		return err
	}

	return nil
}
