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
	"bufio"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
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
	Path          string                        `json:"-"`
	Install       InstallerType                 `json:"install"`
	HeaderLines   []string                      `json:"headerLines"`
	HeaderFormat  ManifestFileHeader            `json:"headerFormat"`
	Format        FormatType                    `json:"format"`
	ContentString string                        `json:"-"`
	ContentObj    map[interface{}]interface{}   `json:"-"`
	ContentArray  []map[interface{}]interface{} `json:"-"`
	Executable    bool                          `json:"executable"`
}

func (f *ManifestFile) ManifestFile(destination string) (loc string, err error) {
	if f.Install == None {
		return "", nil
	}

	dest := filepath.Join(destination, f.Path)
	if err := os.MkdirAll(filepath.Dir(dest), os.ModePerm); err != nil {
		return "", err
	}

	if f.Format == Yaml {
		f.WriteYaml(dest)
	} else if f.Format == Json {
		f.WriteJson(dest)
	} else {
		f.WriteRaw(dest)
	}

	return dest, nil
}

func (f *ManifestFile) HasHeader() bool {
	return !f.HeaderFormat.NoHeader
}

func (f *ManifestFile) Header() (header string) {
	if f.HeaderFormat.NoHeader {
		return ""
	}

	var prefix string
	if f.HeaderFormat.LinesPrefix != "" {
		prefix = f.HeaderFormat.LinesPrefix
	} else if f.Format == Yaml {
		prefix = "# "
	} else {
		prefix = ""
	}

	var b strings.Builder
	b.WriteString(f.HeaderFormat.Prefix + "\n")
	for _, s := range f.HeaderLines {
		b.WriteString(prefix + s + "\n")
	}
	b.WriteString(f.HeaderFormat.Suffix + "\n")
	return b.String()
}

func (f *ManifestFile) WriteYaml(destination string) error {
	file, err := os.OpenFile(destination, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := bufio.NewWriter(file)

	enc := yaml.NewEncoder(writer)
	writer.WriteString(f.Header())
	for _, e := range f.ContentArray {
		enc.Encode(&e)
	}
	writer.Flush()
	return nil
}

func (f *ManifestFile) WriteJson(destination string) error {
	data, err := json.MarshalIndent(f.ContentObj, "", "\t")
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(destination, data, 0644); err != nil {
		return err
	}
	return nil
}

func (f *ManifestFile) WriteRaw(destination string) error {
	file, err := os.OpenFile(destination, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	scanner := bufio.NewScanner(strings.NewReader(f.ContentString))
	scanner.Scan()

	// Write the header after and shebang
	if strings.HasPrefix(scanner.Text(), "#!") {
		writer.WriteString(scanner.Text() + "\n")
		writer.WriteString(f.Header())
	} else {
		writer.WriteString(f.Header())
		writer.WriteString(scanner.Text() + "\n")
	}

	for scanner.Scan() {
		writer.WriteString(scanner.Text() + "\n")
	}
	writer.Flush()
	return nil
}
