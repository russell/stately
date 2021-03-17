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
		file := sfile.(map[string]interface{})

		mFile := ManifestFile{
			Path: path,
		}

		switch executable := file["executable"].(type) {
		case bool:
			mFile.Executable = executable
		default:
			return ManifestContainer{}, fmt.Errorf("Missing required field 'executable' for file %s", path)
		}

		switch install := file["install"].(type) {
		case string:
			mFile.Install = InstallerType(strings.ToLower(install))
		default:
			return ManifestContainer{}, fmt.Errorf("Missing required field 'install' for file %s", path)
		}

		switch format := file["format"].(type) {
		case string:
			switch FormatType(strings.ToLower(format)) {
			case Yaml, Json, Raw:
				mFile.Format = FormatType(strings.ToLower(format))
			default:
				return ManifestContainer{}, fmt.Errorf("Invalid 'format' '%s' for file %s", format, path)
			}
		default:
			return ManifestContainer{}, fmt.Errorf("Missing required field 'format' for file %s", path)
		}

		switch contents := file["contents"].(type) {
		case string:
			mFile.Contents = contents
		default:
			return ManifestContainer{}, fmt.Errorf("Missing required field 'contents' for file %s", path)
		}

		switch headerLines := files["headerLines"].(type) {
		case []string:
			mFile.HeaderLines = headerLines
		}

		switch headerFormat := files["headerFormat"].(type) {
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

		manifests.Files = append(manifests.Files, mFile)
	}
	return manifests, nil
}
