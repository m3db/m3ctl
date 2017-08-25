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

	"github.com/m3db/m3ctl/r2"
)

type stubRuleSet struct {
}

type stubData struct {
	Namespaces        []string
	RuleSets          map[string]r2.RuleSet
	ErrorNamespace    string
	ConflictNamespace string
}

var (
	errNamespace = "errNs"
	conflicatNs  = "conflictNs"
	dummyData    = stubData{
		Namespaces: []string{"ns1", "ns2"},
		RuleSets: map[string]r2.RuleSet{
			"ns1": r2.RuleSet{},
			"ns2": r2.RuleSet{},
		},
	}
)

// Operator contains the data necessary to implement stubbed out implementations for various r2 operations.
type store struct {
	data *stubData
}

func NewStore() r2.Operator {
	return &store{data: &dummyData}
}

// FetchNamespaces ...
func (s *store) FetchNamespaces() ([]string, int, error) {
	return o.data.Namespaces, 1, nil
}

// CreateNamespace ...
func (s *store) CreateNamespace(ns r2.NamespaceData) (string, error) {
	switch namespaceID := ns.ID; namespaceID {
	case errNamespace:
		return "", r2.NewInternalError(fmt.Sprintf("Could not create namespace: %s", namespaceID))
	case conflicatNs:
		return "", r2.NewVersionError(fmt.Sprintf("Namespaces version mismatch"))
	default:
		for _, n := range o.data.Namespaces {
			if namespaceID == n {
				return "", r2.NewConflictError(fmt.Sprintf("Namespace %s already exists", namespaceID))
			}
		}
		return namespaceID, nil
	}
}

// DeleteNamespace ...
func (s *store) DeleteNamespace(ns r2.NamespaceData) error {
	switch namespaceID := ns.ID; namespaceID {
	case errNamespace:
		return r2.NewInternalError("Could not delete namespace")
	case conflicatNs:
		return r2.NewVersionError("Namespace version mismatch")
	default:
		for _, n := range o.data.Namespaces {
			if namespaceID == n {
				return nil
			}
		}
		return r2.NewConflictError(fmt.Sprintf("Namespace %s doesn't exists", namespaceID))
	}
}

// FetchRuleSet ...
func (s *store) FetchRuleSet(namespaceID string) (*r2.RuleSet, error) {
	switch namespaceID {
	case errNamespace:
		return nil, r2.NewInternalError(fmt.Sprintf("Could not fetch namespace: %s", namespaceID))
	case conflicatNs:
		return nil, r2.NewVersionError("Namespace version mismatch")
	default:
		for _, n := range o.data.Namespaces {
			if namespaceID == n {
				rs, _ := o.data.RuleSets[namespaceID]
				return &rs, nil
			}
		}
		return nil, r2.NewConflictError(fmt.Sprintf("Namespace %s doesn't exists", namespaceID))
	}
}
