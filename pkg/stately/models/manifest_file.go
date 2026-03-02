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
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/google/go-jsonnet/formatter"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

type FormatType string

const (
	Yaml    FormatType = "yaml"
	Json    FormatType = "json"
	Jsonnet FormatType = "jsonnet"
	Raw     FormatType = "raw"
)

type InstallerType string

const (
	Symlink      InstallerType = "symlink"
	Write        InstallerType = "write"
	MergeSection InstallerType = "mergesection"
	None         InstallerType = "none"
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
	Content      interface{}        `json:"-"`
	Executable   bool               `json:"executable"`
	SectionStart string             `json:"sectionStart,omitempty"`
	SectionEnd   string             `json:"sectionEnd,omitempty"`
}

func (f *ManifestFile) ManifestFile(destination string, Logger *zap.SugaredLogger) (loc string, err error) {
	if f.Install == None {
		return "", nil
	}

	dest := filepath.Join(destination, f.Path)
	if err := os.MkdirAll(filepath.Dir(dest), os.ModePerm); err != nil {
		return "", err
	}

	if f.Install == Symlink {
		Logger.Infof("Installing as Symlink is not supported %s", f.Path)
	}

	if f.Install == MergeSection {
		if err := f.MergeSectionRaw(dest); err != nil {
			return "", err
		}
		return dest, nil
	}

	if f.Install == Write || f.Install == Symlink {
		switch f.Format {
		case Yaml:
			f.WriteYaml(dest)
		case Json:
			f.WriteJson(dest)
		case Jsonnet:
			f.WriteJsonnet(dest)
		default:
			f.WriteRaw(dest)
		}
		return dest, nil
	}

	return "", fmt.Errorf("Unknown install type '%s' for file '%s'", f.Install, f.Path)
}

func (f *ManifestFile) HasHeader() bool {
	return !f.HeaderFormat.NoHeader
}

func (f *ManifestFile) Permissions() os.FileMode {
	if f.Executable {
		return 0755
	}
	return 0644
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
	file, err := os.OpenFile(destination, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, f.Permissions())
	if err != nil {
		return err
	}
	defer file.Close()
	writer := bufio.NewWriter(file)

	enc := yaml.NewEncoder(writer)
	writer.WriteString(f.Header())
	contents := reflect.ValueOf(f.Content)
	switch contents.Kind() {
	case reflect.Slice:
		for _, e := range f.Content.([]interface{}) {
			enc.Encode(&e)
		}
	case reflect.String:
		writer.WriteString(f.Content.(string))
	default:
		enc.Encode(f.Content)
	}
	writer.Flush()
	return nil
}

func (f *ManifestFile) WriteJson(destination string) error {
	file, err := os.OpenFile(destination, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, f.Permissions())
	if err != nil {
		return err
	}
	defer file.Close()
	var data []byte
	data, err = json.MarshalIndent(f.Content, "", "  ")
	if err != nil {
		return err
	}
	file.Write(data)
	return nil
}

func (f *ManifestFile) WriteJsonnet(destination string) error {
	file, err := os.OpenFile(destination, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, f.Permissions())
	if err != nil {
		return err
	}
	defer file.Close()

	content := fmt.Sprintf("%s", f.Content)

	// leave an empty file if the len is 0
	if len(content) == 0 {
		return nil
	}

	formatterOptions := formatter.DefaultOptions()
	formattedJsonnet, err := formatter.Format(f.Path, content, formatterOptions)
	if err != nil {
		return err
	}

	writer := bufio.NewWriter(file)
	writer.WriteString(f.Header())
	writer.WriteString(formattedJsonnet)
	writer.Flush()
	return nil
}

func (f *ManifestFile) WriteRaw(destination string) error {
	file, err := os.OpenFile(destination, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, f.Permissions())
	if err != nil {
		return err
	}
	defer file.Close()

	var content string

	c := reflect.ValueOf(f.Content)
	switch c.Kind() {
	case reflect.String:
		content = f.Content.(string)
	default:
		content = string(fmt.Sprintf("%s", f.Content))
	}

	// leave an empty file if the len is 0
	if len(content) == 0 {
		return nil
	}

	writer := bufio.NewWriter(file)
	scanner := bufio.NewScanner(strings.NewReader(content))
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

// MergeSectionRaw manages a delimited section within an existing file.
// Content between SectionStart and SectionEnd markers is replaced, while
// content outside the markers is preserved. If the file does not exist,
// it is created with only the managed section. If the file exists but
// contains no markers, the managed section is appended.
func (f *ManifestFile) MergeSectionRaw(destination string) error {
	// Validate section markers.
	startTrimmed := strings.TrimSpace(f.SectionStart)
	endTrimmed := strings.TrimSpace(f.SectionEnd)
	if startTrimmed == "" || endTrimmed == "" {
		return fmt.Errorf("invalid section markers: SectionStart and SectionEnd must be non-empty")
	}
	if startTrimmed == endTrimmed {
		return fmt.Errorf("invalid section markers: SectionStart and SectionEnd must be different")
	}

	// Build the managed content string.
	var managedContent string
	c := reflect.ValueOf(f.Content)
	switch c.Kind() {
	case reflect.String:
		managedContent = f.Content.(string)
	default:
		managedContent = fmt.Sprintf("%s", f.Content)
	}

	// Ensure managed content ends with a newline so the end marker
	// always appears on its own line.
	if len(managedContent) > 0 && !strings.HasSuffix(managedContent, "\n") {
		managedContent += "\n"
	}

	managedBlock := f.SectionStart + "\n" + managedContent + f.SectionEnd + "\n"

	existing, err := os.ReadFile(destination)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("reading %s: %w", destination, err)
		}
		// File doesn't exist — write managed block only.
		return os.WriteFile(destination, []byte(managedBlock), f.Permissions())
	}

	content := string(existing)
	lines := strings.Split(content, "\n")

	// Find the first start marker, then the first end marker after it.
	startIdx := -1
	endIdx := -1
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if startIdx == -1 && trimmed == startTrimmed {
			startIdx = i
			continue
		}
		if startIdx != -1 && trimmed == endTrimmed {
			endIdx = i
			break
		}
	}

	var result string
	if startIdx == -1 || endIdx == -1 {
		// No valid managed section found — append.
		if len(content) > 0 && !strings.HasSuffix(content, "\n") {
			content += "\n"
		}
		if len(content) > 0 && !strings.HasSuffix(content, "\n\n") {
			content += "\n"
		}
		result = content + managedBlock
	} else {
		// Replace existing managed section.
		before := strings.Join(lines[:startIdx], "\n")
		after := ""
		if endIdx+1 < len(lines) {
			after = strings.Join(lines[endIdx+1:], "\n")
		}
		result = before
		if len(before) > 0 && !strings.HasSuffix(before, "\n") {
			result += "\n"
		}
		result += managedBlock
		if len(after) > 0 {
			result += after
		}
	}

	return os.WriteFile(destination, []byte(result), f.Permissions())
}
