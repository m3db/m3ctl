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

const (
	emptyErrorMsg = ""
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
	return s.updateHelper.NewUpdateMetadata(time.Now().UnixNano(), "testUser")
}

func (s *store) FetchNamespaces() ([]string, int, error) {
	nss, err := s.ruleStore.ReadNamespaces()
	if err != nil {
		return nil, 0, nil
	}

	// Todo(dgromov): Use NamespaceView here.
	namespaces := nss.Namespaces()
	var namespaceIds []string
	for _, n := range namespaces {
		if !n.Tombstoned() {
			namespaceIds = append(namespaceIds, string(n.Name()))
		}
	}
	return namespaceIds, nss.Version(), nil
}

func (s *store) CreateNamespace(namespaceID string) (string, error) {
	nss, err := s.ruleStore.ReadNamespaces()
	if err != nil {
		return namespaceID, err
	}

	revived, err := nss.AddNamespace(namespaceID)
	if err != nil {
		return "", r2.NewConflictError(err.Error())
	}

	rs := rules.NewEmptyRuleSet(namespaceID, s.updateMeta())

	if revived {
		rawRs, err := s.ruleStore.ReadRuleSet(namespaceID)
		if err != nil {
			return "", r2.NewInternalError(err.Error())
		}
		rs = rawRs.ToMutableRuleSet().Clone()
		if err = rs.Revive(s.updateMeta()); err != nil {
			return "", r2.NewInternalError(err.Error())
		}
	}

	if err = s.ruleStore.WriteAll(nss, rs); err != nil {
		return "", r2.NewInternalError(err.Error())
	}

	return namespaceID, nil
}

func (s *store) DeleteNamespace(namespaceID string) error {
	nss, err := s.ruleStore.ReadNamespaces()
	if err != nil {
		return s.handleUpstreamError(err)
	}
	rs, err := s.ruleStore.ReadRuleSet(namespaceID)

	if err != nil {
		return s.handleUpstreamError(err)
	}

	err = nss.DeleteNamespace(namespaceID, rs.Version())
	if err != nil {
		return s.handleUpstreamError(err)
	}

	mutable := rs.ToMutableRuleSet().Clone()
	err = mutable.Delete(s.updateMeta())
	if err != nil {
		return s.handleUpstreamError(err)
	}

	if err = s.ruleStore.WriteAll(nss, mutable); err != nil {
		return s.handleUpstreamError(err)
	}

	return nil
}

func (s *store) FetchRuleSet(namespaceID string) (*r2.CurrentRuleSet, error) {
	rs, err := s.ruleStore.ReadRuleSet(namespaceID)
	if err != nil {
		return nil, s.handleUpstreamError(err)
	}
	return r2.NewCurrentRuleSetView(rs)
}

func (s *store) FetchMappingRule(namespaceID string, ruleID string) (*rules.MappingRuleView, error) {
	crs, err := s.FetchRuleSet(namespaceID)
	if err != nil {
		return nil, err
	}

	mr, exists := crs.MappingRules[ruleID]
	if !exists {
		return nil, s.mappingRuleNotFoundError(namespaceID, ruleID)
	}
	return mr, nil
}

func (s *store) CreateMappingRule(namespaceID string, mrv *rules.MappingRuleView) (*rules.MappingRuleView, error) {
	rs, err := s.ruleStore.ReadRuleSet(namespaceID)
	if err != nil {
		return nil, s.handleUpstreamError(err)
	}

	mutable := rs.ToMutableRuleSet().Clone()
	newID, err := mutable.AddMappingRule(*mrv, s.updateMeta())
	if err != nil {
		return nil, s.handleUpstreamError(err)
	}

	err = s.ruleStore.WriteRuleSet(mutable)
	if err != nil {
		return nil, s.handleUpstreamError(err)
	}

	crs, err := s.FetchRuleSet(namespaceID)
	if err != nil {
		return nil, s.handleUpstreamError(err)
	}

	mr, exists := crs.MappingRules[newID]
	if !exists {
		return nil, s.mappingRuleNotFoundError(namespaceID, newID)
	}

	return mr, nil
}

func (s *store) UpdateMappingRule(namespaceID, ruleID string, mrv *rules.MappingRuleView) (*rules.MappingRuleView, error) {
	rs, err := s.ruleStore.ReadRuleSet(namespaceID)
	if err != nil {
		return nil, s.handleUpstreamError(err)
	}

	mutable := rs.ToMutableRuleSet().Clone()
	err = mutable.UpdateMappingRule(*mrv, s.updateMeta())
	if err != nil {
		return nil, s.handleUpstreamError(err)
	}

	err = s.ruleStore.WriteRuleSet(mutable)
	if err != nil {
		return nil, s.handleUpstreamError(err)
	}

	crs, err := s.FetchRuleSet(namespaceID)
	if err != nil {
		return nil, s.handleUpstreamError(err)
	}

	mr, exists := crs.MappingRules[ruleID]
	if !exists {
		return nil, s.mappingRuleNotFoundError(namespaceID, ruleID)
	}

	return mr, nil
}

func (s *store) DeleteMappingRule(namespaceID, mappingRuleID string) error {
	rs, err := s.ruleStore.ReadRuleSet(namespaceID)
	if err != nil {
		return s.handleUpstreamError(err)
	}

	mutable := rs.ToMutableRuleSet().Clone()
	err = mutable.DeleteMappingRule(mappingRuleID, s.updateMeta())
	if err != nil {
		return s.handleUpstreamError(err)
	}

	err = s.ruleStore.WriteRuleSet(mutable)
	if err != nil {
		return s.handleUpstreamError(err)
	}

	return nil
}

func (s *store) FetchMappingRuleHistory(namespaceID, ruleID string) ([]*rules.MappingRuleView, error) {
	rs, err := s.ruleStore.ReadRuleSet(namespaceID)
	if err != nil {
		return nil, s.handleUpstreamError(err)
	}

	mrs, err := rs.MappingRules()
	if err != nil {
		return nil, s.handleUpstreamError(err)
	}

	hist, exists := mrs[ruleID]
	if !exists {
		return nil, s.mappingRuleNotFoundError(namespaceID, ruleID)
	}

	return hist, nil
}

func (s *store) FetchRollupRule(namespaceID string, ruleID string) (*rules.RollupRuleView, error) {
	crs, err := s.FetchRuleSet(namespaceID)
	if err != nil {
		return nil, s.handleUpstreamError(err)
	}

	rr, exists := crs.RollupRules[ruleID]
	if !exists {
		return nil, s.rollupRuleNotFoundError(namespaceID, ruleID)
	}
	return rr, nil
}

func (s *store) CreateRollupRule(namespaceID string, rrv *rules.RollupRuleView) (*rules.RollupRuleView, error) {
	rs, err := s.ruleStore.ReadRuleSet(namespaceID)
	if err != nil {
		return nil, s.handleUpstreamError(err)
	}

	mutable := rs.ToMutableRuleSet().Clone()
	newID, err := mutable.AddRollupRule(*rrv, s.updateMeta())
	if err != nil {
		return nil, s.handleUpstreamError(err)
	}

	err = s.ruleStore.WriteRuleSet(mutable)
	if err != nil {
		return nil, s.handleUpstreamError(err)
	}

	crs, err := s.FetchRuleSet(namespaceID)
	if err != nil {
		return nil, s.handleUpstreamError(err)
	}

	rr, exists := crs.RollupRules[newID]
	if !exists {
		return nil, s.rollupRuleNotFoundError(namespaceID, newID)
	}

	return rr, nil
}

func (s *store) UpdateRollupRule(namespaceID, rollupRuleID string, rrv *rules.RollupRuleView) (*rules.RollupRuleView, error) {
	rs, err := s.ruleStore.ReadRuleSet(namespaceID)
	if err != nil {
		return nil, s.handleUpstreamError(err)
	}

	mutable := rs.ToMutableRuleSet().Clone()
	err = mutable.UpdateRollupRule(*rrv, s.updateMeta())
	if err != nil {
		return nil, s.handleUpstreamError(err)
	}

	err = s.ruleStore.WriteRuleSet(mutable)
	if err != nil {
		return nil, s.handleUpstreamError(err)
	}

	crs, err := s.FetchRuleSet(namespaceID)
	if err != nil {
		return nil, s.handleUpstreamError(err)
	}

	mr, exists := crs.RollupRules[rollupRuleID]
	if !exists {
		return nil, s.rollupRuleNotFoundError(namespaceID, rollupRuleID)
	}

	return mr, nil
}

func (s *store) DeleteRollupRule(namespaceID, rollupRuleID string) error {
	rs, err := s.ruleStore.ReadRuleSet(namespaceID)
	if err != nil {
		return s.handleUpstreamError(err)
	}

	mutable := rs.ToMutableRuleSet().Clone()
	err = mutable.DeleteRollupRule(rollupRuleID, s.updateMeta())
	if err != nil {
		return s.handleUpstreamError(err)
	}

	err = s.ruleStore.WriteRuleSet(mutable)
	if err != nil {
		return s.handleUpstreamError(err)
	}

	return nil
}

func (s *store) FetchRollupRuleHistory(namespaceID, rollupRuleID string) ([]*rules.RollupRuleView, error) {
	rs, err := s.ruleStore.ReadRuleSet(namespaceID)
	if err != nil {
		return nil, s.handleUpstreamError(err)
	}

	rrs, err := rs.RollupRules()
	if err != nil {
		return nil, s.handleUpstreamError(err)
	}

	hist, exists := rrs[rollupRuleID]
	if !exists {
		return nil, s.rollupRuleNotFoundError(namespaceID, rollupRuleID)
	}
	return hist, nil
}

func (s *store) rollupRuleNotFoundError(namespaceID, ruleID string) error {
	return r2.NewNotFoundError(
		fmt.Sprintf("rollup rule: %s doesn't exist in Namespace: %s",
			ruleID,
			namespaceID,
		),
	)
}

func (s *store) mappingRuleNotFoundError(namespaceID, ruleID string) error {
	return r2.NewNotFoundError(
		fmt.Sprintf("mapping rule: %s doesn't exist in Namespace: %s",
			ruleID,
			namespaceID,
		),
	)
}

func (s *store) handleUpstreamError(err error) error {
	// Todo(dgromov): handle validation error
	switch err.(type) {
	case rules.RuleConflictError:
		return r2.NewConflictError(err.Error())
	default:
		return r2.NewInternalError(err.Error())
	}
}
