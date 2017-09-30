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
	"fmt"
	"time"

	"github.com/m3db/m3ctl/r2"
	"github.com/m3db/m3metrics/rules"
	"github.com/m3db/m3x/instrument"
)

type store struct {
	iOpts        instrument.Options
	updateHelper rules.RuleSetUpdateHelper
	ruleStore    rules.Store
}

// NewStore returns a new service that knows how to talk to a kv backed r2 store.
func NewStore(iOpts instrument.Options, uh rules.RuleSetUpdateHelper, rs rules.Store) r2.Store {
	return &store{iOpts: iOpts, updateHelper: uh, ruleStore: rs}
}

func (s *store) updateMeta() rules.UpdateMetadata {
	return s.updateHelper.NewUpdateMetadata(time.Now().UnixNano())
}

func (s *store) FetchNamespaces() ([]string, int, error) {
	nss, err := s.ruleStore.ReadNamespaces()
	if err != nil {
		return nil, 0, nil
	}
	namespaces := nss.Namespaces()
	namespaceIds := make([]string, len(namespaces))
	for i, n := range nss.Namespaces() {
		namespaceIds[i] = string(n.Name())
	}
	return namespaceIds, nss.Version(), nil
}

func (s *store) CreateNamespace(namespaceID string) (string, error) {
	nss, err := s.ruleStore.ReadNamespaces()
	if err != nil {
		return namespaceID, err
	}

	err = nss.AddNamespace(namespaceID)
	if err != nil {
		return "", r2.NewConflictError(err.Error())
	}

	rs := rules.NewEmptyRuleSet(namespaceID, s.updateMeta())
	if err = s.ruleStore.WriteAll(nss, rs); err != nil {
		return "", r2.NewInternalError(err.Error())
	}

	return namespaceID, nil
}

func (s *store) DeleteNamespace(namespaceID string) error {
	nss, err := s.ruleStore.ReadNamespaces()
	if err != nil {
		return r2.NewInternalError(err.Error())
	}
	rs, err := s.ruleStore.ReadRuleSet(namespaceID)

	if err != nil {
		return r2.NewNotFoundError(err.Error())
	}

	err = nss.DeleteNamespace(namespaceID, rs.Version())
	if err != nil {
		return r2.NewInternalError(err.Error())
	}

	mutable := rs.ToMutableRuleSet().Clone()
	err = mutable.Delete(s.updateMeta())
	if err != nil {
		return r2.NewInternalError(err.Error())
	}

	if err = s.ruleStore.WriteAll(nss, mutable); err != nil {
		return r2.NewInternalError(err.Error())
	}

	return nil
}

func (s *store) FetchRuleSet(namespaceID string) (*r2.CurrentRuleSet, error) {
	rs, err := s.ruleStore.ReadRuleSet(namespaceID)
	if err != nil {
		return nil, r2.NewNotFoundError(err.Error())
	}
	return r2.NewCurrentRuleSetView(rs)
}

func (s *store) FetchMappingRule(namespaceID string, mappingRuleID string) (*rules.MappingRuleView, error) {
	crs, err := s.FetchRuleSet(namespaceID)
	if err != nil {
		return nil, err
	}

	mr, exists := crs.MappingRules[mappingRuleID]
	if !exists {
		return nil, r2.NewNotFoundError(fmt.Sprintf("mapping rule: %s doesn't exist in namespace: %s", mappingRuleID, namespaceID))
	}
	return mr, nil
}

func (s *store) CreateMappingRule(namespaceID string, mrv *rules.MappingRuleView) (*rules.MappingRuleView, error) {
	rs, err := s.ruleStore.ReadRuleSet(namespaceID)
	if err != nil {
		return nil, r2.NewNotFoundError(err.Error())
	}

	mutable := rs.ToMutableRuleSet().Clone()
	newID, err := mutable.AddMappingRule(*mrv, s.updateMeta())
	if err != nil {
		return nil, r2.NewInternalError(err.Error())
	}

	crs, err := s.FetchRuleSet(namespaceID)
	if err != nil {
		return nil, err
	}

	mr, exists := crs.MappingRules[newID]
	if !exists {
		return nil, r2.NewNotFoundError(fmt.Sprintf("mapping rule: %s doesn't exist in namespace: %s", newID, namespaceID))
	}

	return mr, nil
}

func (s *store) UpdateMappingRule(namespaceID, mappingRuleID string, mrv *rules.MappingRuleView) (*rules.MappingRuleView, error) {
	rs, err := s.ruleStore.ReadRuleSet(namespaceID)
	if err != nil {
		return nil, r2.NewNotFoundError(err.Error())
	}

	mutable := rs.ToMutableRuleSet().Clone()
	err = mutable.UpdateMappingRule(*mrv, s.updateMeta())
	if err != nil {
		return nil, r2.NewInternalError(err.Error())
	}

	crs, err := s.FetchRuleSet(namespaceID)
	if err != nil {
		return nil, err
	}

	mr, exists := crs.MappingRules[mappingRuleID]
	if !exists {
		return nil, r2.NewNotFoundError(fmt.Sprintf("mapping rule: %s doesn't exist in namespace: %s", mappingRuleID, namespaceID))
	}

	return mr, nil
}

func (s *store) DeleteMappingRule(namespaceID, mappingRuleID string) error {
	rs, err := s.ruleStore.ReadRuleSet(namespaceID)
	if err != nil {
		return r2.NewNotFoundError(err.Error())
	}

	mutable := rs.ToMutableRuleSet().Clone()
	err = mutable.DeleteMappingRule(mappingRuleID, s.updateMeta())
	if err != nil {
		return r2.NewInternalError(err.Error())
	}

	return nil
}

func (s *store) FetchMappingRuleHistory(namespaceID, mappingRuleID string) ([]*rules.MappingRuleView, error) {
	rs, err := s.ruleStore.ReadRuleSet(namespaceID)
	if err != nil {
		return nil, r2.NewNotFoundError(err.Error())
	}

	mrs, err := rs.MappingRules()
	if err != nil {
		return nil, r2.NewInternalError(err.Error())
	}

	hist, exists := mrs[mappingRuleID]
	if !exists {
		return nil, r2.NewNotFoundError(fmt.Sprintf("mappingRule: %s doesn't exist in Namespace: %s", mappingRuleID, namespaceID))
	}

	return hist, nil
}

func (s *store) FetchRollupRule(namespaceID string, rollupRuleID string) (*rules.RollupRuleView, error) {
	crs, err := s.FetchRuleSet(namespaceID)
	if err != nil {
		return nil, err
	}

	rr, exists := crs.RollupRules[rollupRuleID]
	if !exists {
		return nil, r2.NewNotFoundError(fmt.Sprintf("rollup rule: %s doesn't exist in namespace: %s", rollupRuleID, namespaceID))
	}
	return rr, nil
}

func (s *store) CreateRollupRule(namespaceID string, rrv *rules.RollupRuleView) (*rules.RollupRuleView, error) {
	rs, err := s.ruleStore.ReadRuleSet(namespaceID)
	if err != nil {
		return nil, r2.NewNotFoundError(err.Error())
	}

	mutable := rs.ToMutableRuleSet().Clone()
	newID, err := mutable.AddRollupRule(*rrv, s.updateMeta())
	if err != nil {
		return nil, r2.NewInternalError(err.Error())
	}

	crs, err := s.FetchRuleSet(namespaceID)
	if err != nil {
		return nil, err
	}

	rr, exists := crs.RollupRules[newID]
	if !exists {
		return nil, r2.NewNotFoundError(fmt.Sprintf("rollup rule: %s doesn't exist in namespace: %s", newID, namespaceID))
	}

	return rr, nil
}

func (s *store) UpdateRollupRule(namespaceID, rollupRuleID string, rrv *rules.RollupRuleView) (*rules.RollupRuleView, error) {
	rs, err := s.ruleStore.ReadRuleSet(namespaceID)
	if err != nil {
		return nil, r2.NewNotFoundError(err.Error())
	}

	mutable := rs.ToMutableRuleSet().Clone()
	err = mutable.UpdateRollupRule(*rrv, s.updateMeta())
	if err != nil {
		return nil, r2.NewInternalError(err.Error())
	}

	crs, err := s.FetchRuleSet(namespaceID)
	if err != nil {
		return nil, err
	}

	mr, exists := crs.RollupRules[rollupRuleID]
	if !exists {
		return nil, r2.NewNotFoundError(fmt.Sprintf("rollup rule: %s doesn't exist in namespace: %s", rollupRuleID, namespaceID))
	}

	return mr, nil
}

func (s *store) DeleteRollupRule(namespaceID, rollupRuleID string) error {
	rs, err := s.ruleStore.ReadRuleSet(namespaceID)
	if err != nil {
		return r2.NewNotFoundError(err.Error())
	}

	mutable := rs.ToMutableRuleSet().Clone()
	err = mutable.DeleteRollupRule(rollupRuleID, s.updateMeta())
	if err != nil {
		return r2.NewInternalError(err.Error())
	}

	return nil
}

func (s *store) FetchRollupRuleHistory(namespaceID, rollupRuleID string) ([]*rules.RollupRuleView, error) {
	rs, err := s.ruleStore.ReadRuleSet(namespaceID)
	if err != nil {
		return nil, r2.NewNotFoundError(err.Error())
	}

	rrs, err := rs.RollupRules()
	if err != nil {
		return nil, r2.NewInternalError(err.Error())
	}

	hist, exists := rrs[rollupRuleID]
	if !exists {
		return nil, r2.NewNotFoundError(fmt.Sprintf("rollup rule: %s doesn't exist in Namespace: %s", rollupRuleID, namespaceID))
	}

	return hist, nil
}
