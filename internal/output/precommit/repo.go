// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package precommit

// Repo represents a single entry under repos: in .pre-commit-config.yaml.
type Repo struct {
	URL   string  `yaml:"repo"`
	Rev   string  `yaml:"rev,omitempty"`
	Hooks []*Hook `yaml:"hooks"`
}

// Revision sets the rev field of the repo.
func (r *Repo) Revision(rev string) *Repo {
	r.Rev = rev

	return r
}

// Hook appends a new hook to the repo and returns it for further configuration.
func (r *Repo) Hook(id string) *Hook {
	h := &Hook{ID: id}

	r.Hooks = append(r.Hooks, h)

	return h
}
