// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package precommit_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.yaml.in/yaml/v4"

	"github.com/siderolabs/kres/internal/output/precommit"
)

func TestDisabledOutputProducesNoFiles(t *testing.T) {
	o := precommit.NewOutput()

	assert.Nil(t, o.Filenames())
}

func TestSampleConfigRoundTrip(t *testing.T) {
	o := precommit.NewOutput()
	o.Enable()

	repo := o.Repo("https://github.com/pre-commit/pre-commit-hooks").Revision("v3.2.0")
	repo.Hook("trailing-whitespace")
	repo.Hook("end-of-file-fixer")
	repo.Hook("check-yaml")
	repo.Hook("check-added-large-files")

	require.Equal(t, []string{".pre-commit-config.yaml"}, o.Filenames())

	var buf bytes.Buffer

	require.NoError(t, o.GenerateFile(".pre-commit-config.yaml", &buf))

	var decoded struct {
		Repos []struct {
			Repo  string `yaml:"repo"`
			Rev   string `yaml:"rev"`
			Hooks []struct {
				ID string `yaml:"id"`
			} `yaml:"hooks"`
		} `yaml:"repos"`
	}

	require.NoError(t, yaml.Unmarshal([]byte(stripPreamble(buf.String())), &decoded))

	require.Len(t, decoded.Repos, 1)
	assert.Equal(t, "https://github.com/pre-commit/pre-commit-hooks", decoded.Repos[0].Repo)
	assert.Equal(t, "v3.2.0", decoded.Repos[0].Rev)

	hookIDs := make([]string, len(decoded.Repos[0].Hooks))
	for i, h := range decoded.Repos[0].Hooks {
		hookIDs[i] = h.ID
	}

	assert.Equal(t, []string{
		"trailing-whitespace",
		"end-of-file-fixer",
		"check-yaml",
		"check-added-large-files",
	}, hookIDs)
}

func TestHookBuilders(t *testing.T) {
	o := precommit.NewOutput()
	o.Enable()

	o.Repo("local").
		Hook("go-fmt").
		WithName("go fmt").
		WithLanguage("system").
		WithEntry("gofmt -l -w").
		WithFiles(`\.go$`).
		WithExclude(`^vendor/`).
		WithArgs("-s").
		WithStages("pre-commit", "pre-push").
		WithAdditionalDependencies("golang.org/x/tools/cmd/goimports@latest").
		WithAlwaysRun().
		WithPassFilenames(false).
		WithAlias("go-fmt-strict")

	var buf bytes.Buffer

	require.NoError(t, o.GenerateFile(".pre-commit-config.yaml", &buf))

	var decoded struct {
		Repos []struct {
			Repo  string     `yaml:"repo"`
			Hooks []struct { //nolint:govet
				ID                     string   `yaml:"id"`
				Alias                  string   `yaml:"alias"`
				Name                   string   `yaml:"name"`
				Args                   []string `yaml:"args"`
				Files                  string   `yaml:"files"`
				Exclude                string   `yaml:"exclude"`
				Stages                 []string `yaml:"stages"`
				AdditionalDependencies []string `yaml:"additional_dependencies"`
				AlwaysRun              bool     `yaml:"always_run"`
				PassFilenames          *bool    `yaml:"pass_filenames"`
				Language               string   `yaml:"language"`
				Entry                  string   `yaml:"entry"`
			} `yaml:"hooks"`
		} `yaml:"repos"`
	}

	require.NoError(t, yaml.Unmarshal([]byte(stripPreamble(buf.String())), &decoded))

	require.Len(t, decoded.Repos, 1)
	assert.Equal(t, "local", decoded.Repos[0].Repo)
	require.Len(t, decoded.Repos[0].Hooks, 1)

	h := decoded.Repos[0].Hooks[0]
	assert.Equal(t, "go-fmt", h.ID)
	assert.Equal(t, "go-fmt-strict", h.Alias)
	assert.Equal(t, "go fmt", h.Name)
	assert.Equal(t, []string{"-s"}, h.Args)
	assert.Equal(t, `\.go$`, h.Files)
	assert.Equal(t, `^vendor/`, h.Exclude)
	assert.Equal(t, []string{"pre-commit", "pre-push"}, h.Stages)
	assert.Equal(t, []string{"golang.org/x/tools/cmd/goimports@latest"}, h.AdditionalDependencies)
	assert.True(t, h.AlwaysRun)
	require.NotNil(t, h.PassFilenames)
	assert.False(t, *h.PassFilenames)
	assert.Equal(t, "system", h.Language)
	assert.Equal(t, "gofmt -l -w", h.Entry)
}

// stripPreamble removes the leading kres preamble (# comment lines and the
// blank line that follows it) so the remaining YAML can be unmarshaled.
func stripPreamble(s string) string {
	var (
		out        []string
		inPreamble = true
	)

	for line := range strings.SplitSeq(s, "\n") {
		if inPreamble && (strings.HasPrefix(line, "#") || line == "") {
			continue
		}

		inPreamble = false

		out = append(out, line)
	}

	return strings.Join(out, "\n")
}
