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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type ManifestContainer struct {
	Files []ManifestFile `json:"directories"`
}

func NewManifestContainerFromFile(path string) (ManifestContainer, error) {
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		return ManifestContainer{}, fmt.Errorf("Reading state config '%s': %s", path, err)
	}

	return NewManifestContainerFromBytes(bs)
}

func NewManifestContainerFromStdin() (ManifestContainer, error) {
	bs, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return ManifestContainer{}, fmt.Errorf("Reading state config stdin: %s", err)
	}

	return NewManifestContainerFromBytes(bs)
}

func NewManifestContainerFromBytes(bs []byte) (ManifestContainer, error) {
	var manifests ManifestContainer

	var container map[string]interface{}
	json.Unmarshal(bs, &container)

	files := container["files"].(map[string]interface{})

	for path, sfile := range files {
		switch sfileOrList := sfile.(type) {
		case map[string]interface{}:
			err := manifests.AddFile(path, sfile.(map[string]interface{}))
			if err != nil {
				return ManifestContainer{}, err
			}
		case []interface{}:
			for _, ssfile := range sfile.([]interface{}) {
				switch ssfileOrList := ssfile.(type) {
				case map[string]interface{}:
					ssfilei := ssfile.(map[string]interface{})
					switch ssfilepath := ssfilei["path"].(type) {
					case string:
						subpath := filepath.Join(path, ssfilepath)
						err := manifests.AddFile(subpath, ssfile.(map[string]interface{}))
						if err != nil {
							return ManifestContainer{}, err
						}
					default:
						return ManifestContainer{}, fmt.Errorf("Missing required field 'path' for file in list %s", path)
					}

				default:
					return ManifestContainer{}, fmt.Errorf("Invalid file for %s is unsupported type %s", path, ssfileOrList)
				}
			}
		default:
			return ManifestContainer{}, fmt.Errorf("Invalid file for %s is unsupported type %s", path, sfileOrList)
		}

	}
	return manifests, nil
}

func (m *ManifestContainer) AddFile(path string, file map[string]interface{}) error {
	mFile := ManifestFile{
		Path: path,
	}

	switch executable := file["executable"].(type) {
	case bool:
		mFile.Executable = executable
	default:
		return fmt.Errorf("Missing required field 'executable' for file %s", path)
	}

	switch install := file["install"].(type) {
	case string:
		mFile.Install = InstallerType(strings.ToLower(install))
	default:
		return fmt.Errorf("Missing required field 'install' for file %s", path)
	}

	switch format := file["format"].(type) {
	case string:
		switch FormatType(strings.ToLower(format)) {
		case Yaml, Json, Raw:
			mFile.Format = FormatType(strings.ToLower(format))
		default:
			return fmt.Errorf("Invalid 'format' '%s' for file %s", format, path)
		}
	default:
		return fmt.Errorf("Missing required field 'format' for file %s", path)
	}

	mFile.Content = file["contents"]

	switch headerLines := file["headerLines"].(type) {
	case []interface{}:
		for _, l := range headerLines {
			switch line := l.(type) {
			case string:
				mFile.HeaderLines = append(mFile.HeaderLines, line)
			}
		}
	}

	switch headerFormat := file["headerFormat"].(type) {
	case map[string]interface{}:
		mHeader := ManifestFileHeader{
			Prefix:      headerFormat["prefix"].(string),
			Suffix:      headerFormat["suffix"].(string),
			LinesPrefix: headerFormat["linePrefix"].(string),
		}
		mFile.HeaderFormat = mHeader
	case string:
		mHeader := ManifestFileHeader{
			LinesPrefix: headerFormat,
		}
		mFile.HeaderFormat = mHeader
	case bool:
		if headerFormat == false {
			mHeader := ManifestFileHeader{
				NoHeader: true,
			}
			mFile.HeaderFormat = mHeader
		}
	default:
		mHeader := ManifestFileHeader{
			NoHeader: true,
		}
		mFile.HeaderFormat = mHeader
	}

	m.Files = append(m.Files, mFile)
	return nil
}
