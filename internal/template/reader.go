/*
 *
 * Copyright 2023 @ Linying Assad <linying@apache.org>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 * /
 */

package template

import (
	"bytes"
	"embed"
	"html/template"
	"k8s.io/klog/v2"
)

//go:embed *
var content embed.FS

// Read reads the content of a file from the go embed fs
func Read(path string) (string, error) {
	data, err := content.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// ReadOrPanic reads the content of a file from the embed fs and panics if it fails
func ReadOrPanic(path string) string {
	data, err := content.ReadFile(path)
	if err != nil {
		klog.Error("Fail to read embed file: " + path)
		panic(err)
	}
	return string(data)
}

// NewTemplateOrPanic creates a new template from the embed fs and panics if it fails
func NewTemplateOrPanic(name string, path string) *template.Template {
	fileContent := ReadOrPanic(path)
	tmpl, err := template.New(name).Parse(fileContent)
	if err != nil {
		klog.Error("Fail to parse template: " + path)
		panic(err)
	}
	return tmpl
}

// ExecTemplate executes a template with the given data
func ExecTemplate(tmpl *template.Template, data any) (string, error) {
	var buf bytes.Buffer
	err := tmpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
