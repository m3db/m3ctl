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

package stub

import (
	"fmt"
	"time"

	"github.com/m3db/m3ctl/r2"
	"github.com/m3db/m3metrics/policy"
	"github.com/m3db/m3x/instrument"
)

type mappingRuleHistories map[string][]*r2.MappingRule
type rollupRuleHistories map[string][]*r2.RollupRule

type stubData struct {
	Namespaces        []string
	ErrorNamespace    string
	ConflictNamespace string
	RuleSets          map[string]r2.RuleSet
	MappingHistory    map[string]mappingRuleHistories
	RollupHistory     map[string]rollupRuleHistories
}

func makePolicy(s string) policy.Policy {
	p, _ := policy.ParsePolicy(s)
	return p
}

var (
	cutoverTimestamp = time.Now().UnixNano()
	dummyData        = stubData{
		ErrorNamespace:    "errNs",
		ConflictNamespace: "conflictNs",
		Namespaces:        []string{"ns1", "ns2"},
		RuleSets: map[string]r2.RuleSet{
			"ns1": r2.RuleSet{
				Namespace:    "ns1",
				Version:      1,
				CutoverNanos: cutoverTimestamp,
				MappingRules: []r2.MappingRule{
					r2.MappingRule{
						ID:           "mr_id1",
						Name:         "mr1",
						Tombstoned:   false,
						CutoverNanos: cutoverTimestamp,
						RawFilters: map[string]string{
							"tag1": "val1",
							"tag2": "val2",
						},
						Policies: []policy.Policy{
							makePolicy("1m:10d"),
							makePolicy("10m:30d"),
						},
					},
					r2.MappingRule{
						ID:           "mr_id2",
						Name:         "mr2",
						Tombstoned:   false,
						CutoverNanos: cutoverTimestamp,
						RawFilters: map[string]string{
							"tag2": "val2",
						},
						Policies: []policy.Policy{
							makePolicy("1m:10d"),
						},
					},
				},
				RollupRules: []r2.RollupRule{
					r2.RollupRule{
						ID:           "rr_id1",
						Name:         "rr1",
						Tombstoned:   false,
						CutoverNanos: cutoverTimestamp,
						RawFilters: map[string]string{
							"tag1": "val1",
							"tag2": "val2",
						},
						Targets: []r2.RollupTarget{
							r2.RollupTarget{
								Name: "testTarget",
								Tags: []string{"tag1", "tag2"},
								Policies: []policy.Policy{
									makePolicy("1m:10d|Min"),
								},
							},
						},
					},
					r2.RollupRule{
						ID:           "rr_id2",
						Name:         "rr2",
						Tombstoned:   false,
						CutoverNanos: cutoverTimestamp,
						RawFilters: map[string]string{
							"tag1": "val1",
						},
						Targets: []r2.RollupTarget{
							r2.RollupTarget{
								Name: "testTarget",
								Tags: []string{"tag1", "tag2"},
								Policies: []policy.Policy{
									makePolicy("1m:30d|Min"),
								},
							},
						},
					},
				},
			},
			"ns2": r2.RuleSet{
				Namespace:    "ns2",
				Version:      1,
				CutoverNanos: cutoverTimestamp,
				MappingRules: []r2.MappingRule{},
				RollupRules: []r2.RollupRule{
					r2.RollupRule{
						ID:           "rr_id3",
						Name:         "rr1",
						Tombstoned:   false,
						CutoverNanos: cutoverTimestamp,
						RawFilters: map[string]string{
							"tag1": "val1",
							"tag2": "val2",
						},
						Targets: []r2.RollupTarget{
							r2.RollupTarget{
								Name: "testTarget",
								Tags: []string{"tag1", "tag2"},
								Policies: []policy.Policy{
									makePolicy("1m:10d|Min,Max"),
								},
							},
							r2.RollupTarget{
								Name: "testTarget",
								Tags: []string{"tag1", "tag2"},
								Policies: []policy.Policy{
									makePolicy("1m:10d|P999"),
								},
							},
						},
					},
				},
			},
		},
		MappingHistory: map[string]mappingRuleHistories{
			"ns1": nil,
			"ns2": nil,
		},
		RollupHistory: map[string]rollupRuleHistories{
			"ns1": rollupRuleHistories{
				"rr_id3": []*r2.RollupRule{
					&r2.RollupRule{
						ID:           "rr_id3",
						Name:         "rr1",
						Tombstoned:   false,
						CutoverNanos: cutoverTimestamp,
						RawFilters: map[string]string{
							"tag1": "val1",
							"tag2": "val2",
						},
						Targets: []r2.RollupTarget{
							r2.RollupTarget{
								Name: "testTarget",
								Tags: []string{"tag1", "tag2"},
								Policies: []policy.Policy{
									makePolicy("1m:10d|Min,Max"),
								},
							},
							r2.RollupTarget{
								Name: "testTarget",
								Tags: []string{"tag1", "tag2"},
								Policies: []policy.Policy{
									makePolicy("1m:10d|P999"),
								},
							},
						},
					},
					&r2.RollupRule{
						ID:           "rr_id3",
						Name:         "rr1",
						Tombstoned:   false,
						CutoverNanos: cutoverTimestamp,
						RawFilters: map[string]string{
							"tag1": "val1",
						},
						Targets: []r2.RollupTarget{
							r2.RollupTarget{
								Name: "testTarget",
								Tags: []string{"tag1", "tag2"},
								Policies: []policy.Policy{
									makePolicy("1m:10d|Min,Max"),
								},
							},
						},
					},
				},
			},
			"ns2": nil,
		},
	}
)

// Operator contains the data necessary to implement stubbed out implementations for various r2 operations.
type store struct {
	data  *stubData
	iOpts instrument.Options
}

// NewStore creates a new stub
func NewStore(iOpts instrument.Options) r2.Store {
	return &store{data: &dummyData, iOpts: iOpts}
}

// FetchNamespaces ...
func (s *store) FetchNamespaces() ([]string, int, error) {
	return s.data.Namespaces, 1, nil
}

// CreateNamespace ...
func (s *store) CreateNamespace(ns r2.NamespaceData) (string, error) {
	switch namespaceID := ns.ID; namespaceID {
	case s.data.ErrorNamespace:
		return "", r2.NewInternalError(fmt.Sprintf("could not create namespace: %s", namespaceID))
	case s.data.ConflictNamespace:
		return "", r2.NewVersionError(fmt.Sprintf("namespaces version mismatch"))
	default:
		for _, n := range s.data.Namespaces {
			if namespaceID == n {
				return "", r2.NewConflictError(fmt.Sprintf("namespace %s already exists", namespaceID))
			}
		}
		return namespaceID, nil
	}
}

// DeleteNamespace ...
func (s *store) DeleteNamespace(ns r2.NamespaceData) error {
	switch namespaceID := ns.ID; namespaceID {
	case s.data.ErrorNamespace:
		return r2.NewInternalError("could not delete namespace")
	case s.data.ConflictNamespace:
		return r2.NewVersionError("namespace version mismatch")
	default:
		for _, n := range s.data.Namespaces {
			if namespaceID == n {
				return nil
			}
		}
		return r2.NewConflictError(fmt.Sprintf("namespace %s doesn't exist", namespaceID))
	}
}

// FetchRuleSet ...
func (s *store) FetchRuleSet(namespaceID string) (*r2.RuleSet, error) {
	switch namespaceID {
	case s.data.ErrorNamespace:
		return nil, r2.NewInternalError(fmt.Sprintf("could not fetch namespace: %s", namespaceID))
	default:
		for _, n := range s.data.Namespaces {
			if namespaceID == n {
				rs := s.data.RuleSets[namespaceID]
				return &rs, nil
			}
		}
		return nil, r2.NewNotFoundError(fmt.Sprintf("namespace %s doesn't exist", namespaceID))
	}
}

// FetchMappingRule ...
func (s *store) FetchMappingRule(namespaceID string, mappingRuleID string) (*r2.MappingRule, error) {
	switch namespaceID {
	case s.data.ErrorNamespace:
		return nil, r2.NewInternalError(fmt.Sprintf("could not fetch mappingRule: %s in namespace: %s", namespaceID, mappingRuleID))
	default:
		rs, exists := s.data.RuleSets[namespaceID]
		if !exists {
			return nil, r2.NewNotFoundError(fmt.Sprintf("namespace %s doesn't exist", namespaceID))
		}

		for _, m := range rs.MappingRules {
			if mappingRuleID == m.ID {
				return &m, nil
			}
		}
		return nil, r2.NewNotFoundError(fmt.Sprintf("mappingRule: %s doesn't exist in Namespace: %s", mappingRuleID, namespaceID))
	}
}

// CreateMappingRule ...
func (s *store) CreateMappingRule(namespaceID string, mrd r2.MappingRuleData) (*r2.MappingRule, error) {
	switch namespaceID {
	case s.data.ErrorNamespace:
		return nil, r2.NewInternalError("could not create mapping rule")
	case s.data.ConflictNamespace:
		return nil, r2.NewVersionError("namespaces version mismatch")
	default:
		rs, exists := s.data.RuleSets[namespaceID]
		if !exists {
			return nil, r2.NewNotFoundError(fmt.Sprintf("namespace %s doesn't exist", namespaceID))
		}

		for _, m := range rs.MappingRules {
			if mrd.Name == m.Name {
				return nil, r2.NewConflictError(fmt.Sprintf("mapping rule: %s already exists in namespace: %s", mrd.Name, namespaceID))
			}
		}

		return &r2.MappingRule{
			ID:           "new",
			Name:         mrd.Name,
			Tombstoned:   false,
			CutoverNanos: time.Now().UnixNano(),
			RawFilters:   mrd.Filter,
			Policies:     mrd.Policies,
		}, nil
	}
}

// UpdateMappingRule ...
func (s *store) UpdateMappingRule(namespaceID, mappingRuleID string, mrd r2.MappingRuleData) (*r2.MappingRule, error) {
	switch namespaceID {
	case s.data.ErrorNamespace:
		return nil, r2.NewInternalError("could not update mapping rule.")
	case s.data.ConflictNamespace:
		return nil, r2.NewVersionError("namespaces version mismatch")
	default:
		rs, exists := s.data.RuleSets[namespaceID]
		if !exists {
			return nil, r2.NewNotFoundError(fmt.Sprintf("namespace %s doesn't exist", namespaceID))
		}

		for _, m := range rs.MappingRules {
			if mappingRuleID == m.ID {
				return &r2.MappingRule{
					ID:           "new",
					Name:         mrd.Name,
					Tombstoned:   false,
					CutoverNanos: time.Now().UnixNano(),
					RawFilters:   mrd.Filter,
					Policies:     mrd.Policies,
				}, nil
			}
		}
		return nil, r2.NewNotFoundError(fmt.Sprintf("mapping rule: %s doesn't exist in namespace: %s", mappingRuleID, namespaceID))
	}
}

// DeleteMappingRule ...
func (s *store) DeleteMappingRule(namespaceID, mappingRuleID string, ruleSetVersion int) error {
	switch namespaceID {
	case s.data.ErrorNamespace:
		return r2.NewInternalError("could not delete mapping rule.")
	case s.data.ConflictNamespace:
		return r2.NewVersionError("namespaces version mismatch")
	default:
		rs, exists := s.data.RuleSets[namespaceID]
		if !exists {
			return r2.NewNotFoundError(fmt.Sprintf("namespace %s doesn't exist", namespaceID))
		}
		for _, m := range rs.MappingRules {
			if mappingRuleID == m.ID {
				return nil
			}
		}
		return r2.NewNotFoundError(fmt.Sprintf("mapping rule: %s doesn't exist in namespace: %s", mappingRuleID, namespaceID))
	}
}

// FetchMappingRuleHistory ...
func (s *store) FetchMappingRuleHistory(namespaceID, mappingRuleID string) ([]*r2.MappingRule, error) {
	switch namespaceID {
	case s.data.ErrorNamespace:
		return nil, r2.NewInternalError(fmt.Sprintf("Could not fetch mappingRuleID: %s in namespace: %s", namespaceID, mappingRuleID))
	default:
		ns, exists := s.data.MappingHistory[namespaceID]
		if !exists {
			return nil, r2.NewNotFoundError(fmt.Sprintf("namespace %s doesn't exist", namespaceID))
		}
		hist, exists := ns[mappingRuleID]
		if !exists {
			return nil, r2.NewNotFoundError(fmt.Sprintf("mappingRule: %s doesn't exist in Namespace: %s", mappingRuleID, namespaceID))
		}
		return hist, nil
	}
}

// FetchRollupRule ...
func (s *store) FetchRollupRule(namespaceID, rollupRuleID string) (*r2.RollupRule, error) {
	switch namespaceID {
	case s.data.ErrorNamespace:
		return nil, r2.NewInternalError(fmt.Sprintf("Could not fetch rollupRule: %s in namespace: %s", namespaceID, rollupRuleID))
	default:
		rs, exists := s.data.RuleSets[namespaceID]
		if !exists {
			return nil, r2.NewNotFoundError(fmt.Sprintf("namespace %s doesn't exist", namespaceID))
		}
		for _, r := range rs.RollupRules {
			if rollupRuleID == r.ID {
				return &r, nil
			}
		}
		return nil, r2.NewNotFoundError(fmt.Sprintf("rollupRule: %s doesn't exist in Namespace: %s", rollupRuleID, namespaceID))
	}
}

// CreateRollupRule ...
func (s *store) CreateRollupRule(namespaceID string, rrd r2.RollupRuleData) (*r2.RollupRule, error) {
	switch namespaceID {
	case s.data.ErrorNamespace:
		return nil, r2.NewInternalError("could not create rollup rule")
	case s.data.ConflictNamespace:
		return nil, r2.NewVersionError("namespaces version mismatch")
	default:
		rs, exists := s.data.RuleSets[namespaceID]
		if !exists {
			return nil, r2.NewNotFoundError(fmt.Sprintf("namespace %s doesn't exist", namespaceID))
		}
		for _, r := range rs.RollupRules {
			if rrd.Name == r.Name {
				return nil, r2.NewConflictError(fmt.Sprintf("rollup rule: %s already exists in namespace: %s", rrd.Name, namespaceID))
			}
		}

		return &r2.RollupRule{
			ID:           "new",
			Name:         rrd.Name,
			Tombstoned:   false,
			CutoverNanos: time.Now().UnixNano(),
			RawFilters:   rrd.Filter,
			Targets:      rrd.Targets,
		}, nil
	}
}

// UpdateRollupRule ...
func (s *store) UpdateRollupRule(namespaceID, rollupRuleID string, rrd r2.RollupRuleData) (*r2.RollupRule, error) {
	switch namespaceID {
	case s.data.ErrorNamespace:
		return nil, r2.NewInternalError("could not update rollup rule.")
	case s.data.ConflictNamespace:
		return nil, r2.NewVersionError("namespaces version mismatch")
	default:
		rs, exists := s.data.RuleSets[namespaceID]
		if !exists {
			return nil, r2.NewNotFoundError(fmt.Sprintf("namespace %s doesn't exist", namespaceID))
		}
		for _, r := range rs.RollupRules {
			if rollupRuleID == r.ID {
				return &r2.RollupRule{
					ID:           "new",
					Name:         rrd.Name,
					Tombstoned:   false,
					CutoverNanos: time.Now().UnixNano(),
					RawFilters:   rrd.Filter,
					Targets:      rrd.Targets,
				}, nil
			}
		}
		return nil, r2.NewNotFoundError(fmt.Sprintf("rollup rule: %s doesn't exist in namespace: %s", rollupRuleID, namespaceID))
	}
}

// DeleteRollupRule ...
func (s *store) DeleteRollupRule(namespaceID, rollupRuleID string, ruleSetVersion int) error {
	switch namespaceID {
	case s.data.ErrorNamespace:
		return r2.NewInternalError("could not delete rollup rule.")
	case s.data.ConflictNamespace:
		return r2.NewVersionError("namespaces version mismatch")
	default:
		rs, exists := s.data.RuleSets[namespaceID]
		if !exists {
			return r2.NewNotFoundError(fmt.Sprintf("namespace %s doesn't exist", namespaceID))
		}
		for _, r := range rs.RollupRules {
			if rollupRuleID == r.ID {
				return nil
			}
		}
		return r2.NewNotFoundError(fmt.Sprintf("rollup rule: %s doesn't exist in namespace: %s", rollupRuleID, namespaceID))
	}
}

// FetchRollupRuleHistory ...
func (s *store) FetchRollupRuleHistory(namespaceID, rollupRuleID string) ([]*r2.RollupRule, error) {
	switch namespaceID {
	case s.data.ErrorNamespace:
		return nil, r2.NewInternalError(fmt.Sprintf("Could not fetch rollupRule: %s in namespace: %s", namespaceID, rollupRuleID))
	default:
		ns, exists := s.data.RollupHistory[namespaceID]
		if !exists {
			return nil, r2.NewNotFoundError(fmt.Sprintf("namespace %s doesn't exist", namespaceID))
		}
		hist, exists := ns[rollupRuleID]
		if !exists {
			return nil, r2.NewNotFoundError(fmt.Sprintf("rollupRule: %s doesn't exist in Namespace: %s", rollupRuleID, namespaceID))
		}
		return hist, nil
	}
}
