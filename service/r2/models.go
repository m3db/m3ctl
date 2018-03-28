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
	"github.com/m3db/m3metrics/rules/models"
)

type mappingRuleHistoryJSON struct {
	MappingRules []models.MappingRule `json:"mappingRules"`
}

func newMappingRuleHistoryJSON(hist []*models.MappingRuleView) mappingRuleHistoryJSON {
	mappingRules := make([]models.MappingRule, len(hist))
	for i, view := range hist {
		mappingRules[i] = models.NewMappingRule(view)
	}
	return mappingRuleHistoryJSON{MappingRules: mappingRules}
}

type rollupRuleHistoryJSON struct {
	RollupRules []models.RollupRule `json:"rollupRules"`
}

func newRollupRuleHistoryJSON(hist []*models.RollupRuleView) rollupRuleHistoryJSON {
	rollupRules := make([]models.RollupRule, len(hist))
	for i, view := range hist {
		rollupRules[i] = models.NewRollupRule(view)
	}
	return rollupRuleHistoryJSON{RollupRules: rollupRules}
}
