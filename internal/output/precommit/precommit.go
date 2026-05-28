// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package precommit implements output to .pre-commit-config.yaml.
package precommit

import (
	"fmt"
	"io"

	"go.yaml.in/yaml/v4"

	"github.com/siderolabs/kres/internal/output"
)

const filename = ".pre-commit-config.yaml"

// Output implements .pre-commit-config.yaml generation.
type Output struct {
	output.FileAdapter

	defaultInstallHookTypes []string
	repos                   []*Repo

	enabled bool
}

// config is the YAML shape of the top-level .pre-commit-config.yaml document.
// Field order here drives key order in the generated file.
type config struct {
	DefaultInstallHookTypes []string `yaml:"default_install_hook_types,omitempty"`
	Repos                   []*Repo  `yaml:"repos"`
}

// NewOutput creates new .pre-commit-config.yaml output.
func NewOutput() *Output {
	output := &Output{}

	output.FileWriter = output

	return output
}

// Compile implements [output.TypedWriter] interface.
func (o *Output) Compile(compiler Compiler) error {
	return compiler.CompilePreCommit(o)
}

// Enable should be called to enable config generation.
func (o *Output) Enable() {
	o.enabled = true
}

// SetDefaultInstallHookTypes sets the git hook types that `pre-commit install`
// wires up by default (e.g. "pre-commit", "commit-msg"). Without this, users
// have to pass --hook-type for any hook stage beyond the default "pre-commit".
func (o *Output) SetDefaultInstallHookTypes(types []string) {
	o.defaultInstallHookTypes = types
}

// Repo appends a new repo entry to the config and returns it for further configuration.
func (o *Output) Repo(url string) *Repo {
	r := &Repo{URL: url}

	o.repos = append(o.repos, r)

	return r
}

// Filenames implements output.FileWriter interface.
func (o *Output) Filenames() []string {
	if !o.enabled {
		return nil
	}

	return []string{filename}
}

// GenerateFile implements output.FileWriter interface.
func (o *Output) GenerateFile(_ string, w io.Writer) error {
	return o.config(w)
}

func (o *Output) config(w io.Writer) error {
	if _, err := w.Write([]byte(output.Preamble("# "))); err != nil {
		return err
	}

	encoder := yaml.NewEncoder(w)
	encoder.SetIndent(2)

	if err := encoder.Encode(config{
		DefaultInstallHookTypes: o.defaultInstallHookTypes,
		Repos:                   o.repos,
	}); err != nil {
		return fmt.Errorf("failed to encode pre-commit config: %w", err)
	}

	return nil
}

// Compiler is implemented by project blocks which support .pre-commit-config.yaml generation.
type Compiler interface {
	CompilePreCommit(*Output) error
}
