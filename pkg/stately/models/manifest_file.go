package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
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
