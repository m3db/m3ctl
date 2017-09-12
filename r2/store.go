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
// THE SOFTWARE

package r2

import "github.com/m3db/m3metrics/policy"

// MappingRule represents a configuration of a mapping rule at a single point in time.
type MappingRule struct {
	ID           string
	Name         string
	Tombstoned   bool
	CutoverNanos int64
	RawFilters   map[string]string
	Policies     []policy.Policy
}

// RollupTarget represents a single rollup target but with types for display purposes.
type RollupTarget struct {
	Name     string          `json:"name"`
	Tags     []string        `json:"tags"`
	Policies []policy.Policy `json:"policies"`
}

// RollupRule represents a configuration of a rollup rule at a single point in time.
type RollupRule struct {
	ID           string
	Name         string
	Tombstoned   bool
	CutoverNanos int64
	RawFilters   map[string]string
	Targets      []RollupTarget
}

// RuleSet represents a latest state slice of a m3metrics ruleset.
type RuleSet struct {
	UUID         string
	Namespace    string
	Version      int
	CutoverNanos int64
	MappingRules []MappingRule
	RollupRules  []RollupRule
}

// Store is a construct that can perform operations against a backing rule store.
type Store interface {
	// Fetch namespaces and their version
	FetchNamespaces() ([]string, int, error)

	CreateNamespace(ns NamespaceData) (string, error)

	DeleteNamespace(ns NamespaceData) error

	FetchRuleSet(namespaceID string) (*RuleSet, error)

	FetchMappingRule(namespaceID, mappingRuleID string) (*MappingRule, error)

	CreateMappingRule(namespaceID string, mrd MappingRuleData) (*MappingRule, error)

	UpdateMappingRule(namespaceID, mappingRuleID string, mrd MappingRuleData) (*MappingRule, error)

	DeleteMappingRule(namespaceID, mappingRuleID string, ruleSetVersion int) error

	FetchMappingRuleHistory(namespaceID, mappingRuleID string) ([]*MappingRule, error)

	FetchRollupRule(namespaceID, rollupRuleID string) (*RollupRule, error)

	CreateRollupRule(namespaceID string, rrd RollupRuleData) (*RollupRule, error)

	UpdateRollupRule(namespaceID, rollupRuleID string, rrd RollupRuleData) (*RollupRule, error)

	DeleteRollupRule(namespaceID, rollupRuleID string, ruleSetVersion int) error

	FetchRollupRuleHistory(namespaceID, rollupRuleID string) ([]*RollupRule, error)
}
