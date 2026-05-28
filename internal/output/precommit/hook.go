// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package precommit

// Hook represents a single entry under repos[*].hooks in .pre-commit-config.yaml.
type Hook struct { //nolint:govet
	ID                     string   `yaml:"id"`
	Alias                  string   `yaml:"alias,omitempty"`
	Name                   string   `yaml:"name,omitempty"`
	Args                   []string `yaml:"args,omitempty"`
	Files                  string   `yaml:"files,omitempty"`
	Exclude                string   `yaml:"exclude,omitempty"`
	Stages                 []string `yaml:"stages,omitempty"`
	AdditionalDependencies []string `yaml:"additional_dependencies,omitempty"`
	AlwaysRun              bool     `yaml:"always_run,omitempty"`

	// PassFilenames is a pointer because pre-commit defaults it to true, while
	// the Go zero value of bool is false. nil means "leave at pre-commit's
	// default", *false suppresses filename arguments (typical for `make`-style
	// entries).
	PassFilenames *bool `yaml:"pass_filenames,omitempty"`

	// For local hooks.
	Language string `yaml:"language,omitempty"`
	Entry    string `yaml:"entry,omitempty"`
}

// WithAlias sets an alternate hook id (used to disambiguate multiple runs of the same hook).
func (h *Hook) WithAlias(alias string) *Hook {
	h.Alias = alias

	return h
}

// WithName overrides the hook's display name.
func (h *Hook) WithName(name string) *Hook {
	h.Name = name

	return h
}

// WithArgs sets additional CLI arguments passed to the hook.
func (h *Hook) WithArgs(args ...string) *Hook {
	h.Args = args

	return h
}

// WithFiles sets a regex matching files the hook should run against.
func (h *Hook) WithFiles(files string) *Hook {
	h.Files = files

	return h
}

// WithExclude sets a regex matching files the hook should skip.
func (h *Hook) WithExclude(exclude string) *Hook {
	h.Exclude = exclude

	return h
}

// WithStages restricts the hook to the given git stages (e.g. "pre-commit", "pre-push").
func (h *Hook) WithStages(stages ...string) *Hook {
	h.Stages = stages

	return h
}

// WithAdditionalDependencies declares packages installed alongside the hook.
func (h *Hook) WithAdditionalDependencies(deps ...string) *Hook {
	h.AdditionalDependencies = deps

	return h
}

// WithAlwaysRun forces the hook to run even when no matching files have changed.
func (h *Hook) WithAlwaysRun() *Hook {
	h.AlwaysRun = true

	return h
}

// WithPassFilenames overrides whether pre-commit appends the matched filenames
// to the hook's command line. Set to false for entries (e.g. `make` targets)
// that don't accept positional file arguments.
func (h *Hook) WithPassFilenames(pass bool) *Hook {
	h.PassFilenames = &pass

	return h
}

// WithLanguage sets the language used to run a local hook (e.g. "system", "python").
func (h *Hook) WithLanguage(language string) *Hook {
	h.Language = language

	return h
}

// WithEntry sets the command/entry point invoked by a local hook.
func (h *Hook) WithEntry(entry string) *Hook {
	h.Entry = entry

	return h
}
