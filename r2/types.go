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
	"strings"

	"github.com/m3db/m3metrics/policy"
)

// API response
type apiResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Namespace Level
type namespaces struct {
	Version    int      `json:"version"`
	Namespaces []string `json:"namespaces"`
}

type createNamespaceBody struct {
	ID                string `json:"id"`
	NamespacesVersion int    `json:"namespacesVersion"`
}

type deleteNamespaceBody struct {
	NamespacesVersion int `json:"namespacesVersion"`
}

// RuleSet View
type filter map[string]string

func (f filter) string() string {
	tmp := make([]string, 0, len(f))
	for k, v := range f {
		tmp = append(tmp, fmt.Sprintf("%s:%s", k, v))
	}
	return strings.Join(tmp, " ")
}

type ruleSetView struct {
	Namespace     string            `json:"namespace"`
	Version       int               `json:"version"`
	CutoverMillis int64             `json:"cutoverMillis"`
	MappingRules  []mappingRuleView `json:"mappingRules"`
	RollupRules   []rollupRuleView  `json:"rollupRules"`
}

func newMappingRuleViews(mrs []MappingRule) []mappingRuleView {
	result := make([]mappingRuleView, len(mrs))
	for i, m := range mrs {
		result[i] = newMappingRuleView(m)
	}
	return result
}

func newRollupRuleViews(rrs []RollupRule) []rollupRuleView {
	result := make([]rollupRuleView, len(rrs))
	for i, r := range rrs {
		result[i] = newRollupRuleView(r)
	}
	return result
}

func newRuleSetView(rs RuleSet) ruleSetView {
	mrs := newMappingRuleViews(rs.MappingRules)
	rrs := newRollupRuleViews(rs.RollupRules)

	return ruleSetView{
		Namespace:     rs.Namespace,
		Version:       rs.Version,
		CutoverMillis: rs.CutoverNanos / 1000,
		MappingRules:  mrs,
		RollupRules:   rrs,
	}
}

// Mapping Rule
type mappingRuleView struct {
	mappingRuleData

	ID             string `json:"id"`
	CutoverMillies int64  `json:"cutoverMillies"`
}

func newMappingRuleView(mr MappingRule) mappingRuleView {
	mrd := mappingRuleData{
		Name:     mr.Name,
		Filter:   filter(mr.RawFilters).string(),
		Policies: mr.Policies,
	}

	return mappingRuleView{
		mappingRuleData: mrd,
		ID:              mr.ID,
		CutoverMillies:  mr.CutoverNanos / 1000,
	}
}

type mappingRuleData struct {
	Name     string          `json:"name"`
	Filter   string          `json:"filter"`
	Policies []policy.Policy `json:"policies"`
}

type mappingRuleUpdate struct { // nolint
	mappingRuleData

	ID             string `json:"id"`
	RuleSetVersion int    `json:"rulesetVersion"`
}

type mappingRuleHistory struct { // nolint
	States []mappingRuleView `json:"states"`
}

// Rollup Rule
type rollupRuleView struct {
	rollupRuleData

	ID             string `json:"id"`
	CutoverMillies int64  `json:"cutoverMillies"`
}

func newRollupRuleView(rr RollupRule) rollupRuleView {
	targets := make([]rollupTargetView, len(rr.Targets))
	for i, t := range rr.Targets {
		targets[i] = rollupTargetView(t)
	}

	rrd := rollupRuleData{
		Name:    rr.Name,
		Filter:  filter(rr.RawFilters).string(),
		Targets: targets,
	}

	return rollupRuleView{
		rollupRuleData: rrd,
		ID:             rr.ID,
		CutoverMillies: rr.CutoverNanos / 1000,
	}
}

type rollupRuleData struct {
	Name    string             `json:"name"`
	Filter  string             `json:"filter"`
	Targets []rollupTargetView `json:"targets"`
}

type rollupTargetView struct {
	Name     string          `json:"name"`
	Tags     []string        `json:"tags"`
	Policies []policy.Policy `json:"policies"`
}

type rollupRuleUpdate struct { // nolint
	rollupRuleData

	ID             string `json:"id"`
	RuleSetVersion int    `json:"rulesetVersion"`
}

type rollupRuleHistory struct { // nolint
	States []rollupRuleView `json:"states"`
}

// Errors

// ConflictError represents either a version mismatch writing data or a data conflict issue.
type ConflictError struct{ msg string }

// NewConflictError creates a new Conflict Error
func NewConflictError(msg string) ConflictError { return ConflictError{msg: msg} }

func (c ConflictError) Error() string { return c.msg }

// InternalError represents an unexpected server error.
type InternalError struct{ msg string }

// NewInternalError creates a new Internal Error
func NewInternalError(msg string) InternalError { return InternalError{msg: msg} }

func (c InternalError) Error() string { return c.msg }

// VersionError represents a mismatch in the Namespaces or Ruleset version specified in the request
// and the latest one.
type VersionError struct{ msg string }

// NewVersionError creates a new Version Error
func NewVersionError(msg string) VersionError { return VersionError{msg: msg} }

func (c VersionError) Error() string { return c.msg }
