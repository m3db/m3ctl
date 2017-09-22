// Copyright (c) 2017 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package r2

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/m3db/m3metrics/policy"
	"github.com/m3db/m3metrics/rules"
)

const (
	tagFilterListSeparator = " "
	tagFilterSeparator     = ":"
)

// CurrentRuleSet represents the current state of a rule set. A rule with all the latest snapshots
// of its ruleset.
type CurrentRuleSet struct {
	Namespace    string
	Version      int
	CutoverNanos int64
	MappingRules []*rules.MappingRuleView
	RollupRules  []*rules.RollupRuleView
}

// nolint
func newCurrentRuleSetView(r rules.RuleSet) (*CurrentRuleSet, error) {
	mrs, err := r.MappingRules()
	if err != nil {
		return nil, err
	}

	rrs, err := r.RollupRules()
	if err != nil {
		return nil, err
	}

	currentMappingRules := newMappingRules(mrs)
	currentRollupRules := newRollupRules(rrs)

	return &CurrentRuleSet{
		Namespace:    string(r.Namespace()),
		Version:      r.Version(),
		CutoverNanos: r.CutoverNanos(),
		MappingRules: currentMappingRules,
		RollupRules:  currentRollupRules,
	}, nil
}

type ruleSetJSON struct {
	Namespace     string            `json:"id"`
	Version       int               `json:"version"`
	CutoverMillis int64             `json:"cutoverMillis"`
	MappingRules  []mappingRuleJSON `json:"mappingRules"`
	RollupRules   []rollupRuleJSON  `json:"rollupRules"`
}

func (c CurrentRuleSet) ruleSetJSON() ruleSetJSON {
	mrJSON := make([]mappingRuleJSON, len(c.MappingRules))
	for i, m := range c.MappingRules {
		mrJSON[i] = newMappingRuleJSON(m)
	}
	rrJSON := make([]rollupRuleJSON, len(c.RollupRules))
	for i, r := range c.RollupRules {
		rrJSON[i] = newRollupRuleJSON(r)
	}
	return ruleSetJSON{
		Namespace:     c.Namespace,
		Version:       c.Version,
		CutoverMillis: nanoToMilli(c.CutoverNanos),
		MappingRules:  mrJSON,
		RollupRules:   rrJSON,
	}
}

// nolint
func newMappingRules(mrs rules.MappingRules) []*rules.MappingRuleView {
	result := make([]*rules.MappingRuleView, len(mrs))
	idx := 0
	for _, m := range mrs {
		if len(m) == 0 {
			continue
		}
		result[idx] = m[len(m)-1]
		idx++
	}
	return result
}

// nolint
func newRollupRules(rrs rules.RollupRules) []*rules.RollupRuleView {
	result := make([]*rules.RollupRuleView, len(rrs))
	idx := 0
	for _, r := range rrs {
		if len(r) == 0 {
			continue
		}
		result[idx] = r[len(r)-1]
		idx++
	}
	return result
}

// API response
type apiResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type namespaceJSON struct {
	ID string `json:"id" validate:"required"`
}

// Namespace Level
type namespacesJSON struct {
	Version    int             `json:"version"`
	Namespaces []namespaceJSON `json:"namespaces"`
}

type filters map[string]string

func (f filters) MarshalJSON() ([]byte, error) {
	tmp := make([]string, 0, len(f))
	for k, v := range f {
		tmp = append(tmp, fmt.Sprintf("%s:%s", k, v))
	}
	marshalled := strconv.Quote(strings.Join(tmp, " "))
	return []byte(marshalled), nil
}

func (f *filters) UnmarshalJSON(data []byte) error {
	str := string(data)
	unquoted, err := strconv.Unquote(str)
	if err != nil {
		return err
	}

	*f = make(map[string]string)
	tagPairs := strings.Split(unquoted, tagFilterListSeparator)
	for _, p := range tagPairs {
		parts := strings.Split(p, tagFilterSeparator)
		if len(parts) != 2 {
			return NewBadInputError(fmt.Sprintf("invalid tag filter pair %s", p))
		}
		(*f)[parts[0]] = parts[1]
	}
	return nil
}

type mappingRuleJSON struct {
	ID             string          `json:"id,omitempty"`
	Name           string          `json:"name" validate:"required"`
	CutoverMillies int64           `json:"cutoverMillies,omitempty"`
	Filters        filters         `json:"filter" validate:"required"`
	Policies       []policy.Policy `json:"policies" validate:"required"`
}

func (m mappingRuleJSON) mappingRuleView() rules.MappingRuleView {
	return rules.MappingRuleView{
		ID:       m.ID,
		Name:     m.Name,
		Filters:  m.Filters,
		Policies: m.Policies,
	}
}

func newMappingRuleJSON(mrv *rules.MappingRuleView) mappingRuleJSON {
	return mappingRuleJSON{
		ID:             mrv.ID,
		Name:           mrv.Name,
		Filters:        mrv.Filters,
		Policies:       mrv.Policies,
		CutoverMillies: nanoToMilli(mrv.CutoverNanos),
	}
}

type mappingRuleHistory struct { // nolint
	States []mappingRuleJSON `json:"states"`
}

// RollupTargetData ...
type rollupTargetJSON struct {
	Name     string          `json:"name" validate:"required"`
	Tags     []string        `json:"tags" validate:"required"`
	Policies []policy.Policy `json:"policies" validate:"required"`
}

func (t rollupTargetJSON) rollupTargetView() rules.RollupTargetView {
	return rules.RollupTargetView{
		Name:     t.Name,
		Tags:     t.Tags,
		Policies: t.Policies,
	}
}

func newRollupTargetJSON(t rules.RollupTargetView) rollupTargetJSON {
	return rollupTargetJSON{
		Name:     t.Name,
		Tags:     t.Tags,
		Policies: t.Policies,
	}
}

type rollupRuleJSON struct {
	ID             string             `json:"id,omitempty"`
	Name           string             `json:"name" validate:"required"`
	Filters        filters            `json:"filter" validate:"required"`
	Targets        []rollupTargetJSON `json:"targets" validate:"required,dive,required"`
	CutoverMillies int64              `json:"cutoverMillies,omitempty"`
}

func newRollupRuleJSON(rrv *rules.RollupRuleView) rollupRuleJSON {
	targets := make([]rollupTargetJSON, len(rrv.Targets))
	for i, t := range rrv.Targets {
		targets[i] = newRollupTargetJSON(t)
	}
	return rollupRuleJSON{
		ID:             rrv.ID,
		Name:           rrv.Name,
		Filters:        rrv.Filters,
		Targets:        targets,
		CutoverMillies: nanoToMilli(rrv.CutoverNanos),
	}
}

func (r rollupRuleJSON) rollupRuleView() rules.RollupRuleView {
	targets := make([]rules.RollupTargetView, len(r.Targets))
	for i, t := range r.Targets {
		targets[i] = t.rollupTargetView()
	}

	return rules.RollupRuleView{
		ID:      r.ID,
		Name:    r.Name,
		Filters: r.Filters,
		Targets: targets,
	}
}

type rollupRuleHistory struct { // nolint
	States []rollupRuleJSON `json:"states"`
}

func nanoToMilli(nanos int64) int64 {
	return nanos / 1000000
}

// ConflictError represents either a version mismatch writing data or a data conflict issue.
type ConflictError struct{ msg string }

// NewConflictError creates a new Conflict Error
func NewConflictError(msg string) ConflictError { return ConflictError{msg: msg} }

func (e ConflictError) Error() string { return e.msg }

// InternalError represents an unexpected server error.
type InternalError struct{ msg string }

// NewInternalError creates a new Internal Error
func NewInternalError(msg string) InternalError { return InternalError{msg: msg} }

func (e InternalError) Error() string { return e.msg }

// VersionError represents a mismatch in the Namespaces or Ruleset version specified in the request
// and the latest one.
type VersionError struct{ msg string }

// NewVersionError creates a new Version Error
func NewVersionError(msg string) VersionError { return VersionError{msg: msg} }

func (e VersionError) Error() string { return e.msg }

// BadInputError represents an error due to malformed or invalid metrics.
type BadInputError struct{ msg string }

// NewBadInputError creates a new Bad Input Error.
func NewBadInputError(msg string) BadInputError { return BadInputError{msg: msg} }

func (e BadInputError) Error() string { return e.msg }

// NotFoundError represents an error due to malformed or invalid metrics.
type NotFoundError struct{ msg string }

// NewNotFoundError creates a new not found Error.
func NewNotFoundError(msg string) NotFoundError { return NotFoundError{msg: msg} }

func (e NotFoundError) Error() string { return e.msg }
