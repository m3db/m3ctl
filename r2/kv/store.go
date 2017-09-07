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

package kv

import (
	"github.com/m3db/m3ctl/r2"
	"github.com/m3db/m3metrics/rules"
)

type store struct {
	cfg StoreConfig
}

// StoreConfig is all the configuration state needed by a kv service.
type StoreConfig struct {
	UpdateHelper rules.RuleSetUpdateHelper
}

// NewStore returns a new service that knows how to talk to a kv backed r2 store.
func NewStore(cfg StoreConfig) r2.Store {
	return &store{cfg: cfg}
}

func (s *store) FetchNamespaces() ([]string, int, error) {
	return nil, 0, nil
}

func (s *store) CreateNamespace(ns r2.NamespaceData) (string, error) {
	return "", nil
}

func (s *store) DeleteNamespace(ns r2.NamespaceData) error {
	return nil
}

func (s *store) FetchRuleSet(namespaceID string) (*r2.RuleSet, error) {
	return nil, nil
}

func (s *store) FetchMappingRule(namespaceID string, ruleSetID string) (*r2.MappingRule, error) {
	return nil, nil
}

func (s *store) CreateMappingRule(namespaceID string, mrd r2.MappingRuleData) (*r2.MappingRule, error) {
	return nil, nil
}

func (s *store) UpdateMappingRule(namespaceID, mappingRuleID string, mrd r2.MappingRuleData) error {
	return nil
}

func (s *store) DeleteMappingRule(namespaceID, mappingRuleID string, ruleSetVersion int) error {
	return nil
}

func (s *store) FetchRollupRule(namespaceID string, rollupRuleID string) (*r2.RollupRule, error) {
	return nil, nil
}

func (s *store) CreateRollupRule(namespaceID string, rrd r2.RollupRuleData) (*r2.RollupRule, error) {
	return nil, nil
}

func (s *store) UpdateRollupRule(namespaceID, rollupRuleID string, rrd r2.RollupRuleData) error {
	return nil
}

func (s *store) DeleteRollupRule(namespaceID, rollupRuleID string, ruleSetVersion int) error {
	return nil
}
func (s *store) FetchRollupRuleHistory(namespaceID, rollupRuleID string) ([]*RollupRule, error) {
	return nil, nil
}
