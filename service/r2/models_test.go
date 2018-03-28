package r2

import (
	"testing"

	"github.com/m3db/m3metrics/policy"
	"github.com/m3db/m3metrics/rules/models"

	"github.com/stretchr/testify/require"
)

func TestNewMappingRuleHistoryJSON(t *testing.T) {
	id := "id"
	hist := []*models.MappingRuleView{
		testMappingRuleView(id, "name1"),
		testMappingRuleView(id, "name2"),
	}
	expected := mappingRuleHistoryJSON{
		MappingRules: []models.MappingRule{
			{
				ID:                  id,
				Name:                "name1",
				Filter:              "filter",
				Policies:            []policy.Policy{},
				CutoverMillis:       0,
				LastUpdatedBy:       "",
				LastUpdatedAtMillis: 0,
			},
			{
				ID:                  id,
				Name:                "name2",
				Filter:              "filter",
				Policies:            []policy.Policy{},
				CutoverMillis:       0,
				LastUpdatedBy:       "",
				LastUpdatedAtMillis: 0,
			},
		},
	}
	require.EqualValues(t, expected, newMappingRuleHistoryJSON(hist))
}

func TestNewRollupRuleHistoryJSON(t *testing.T) {
	id := "id"
	targets := []models.RollupTargetView{
		*testRollupTargetView("target1"),
		*testRollupTargetView("target2"),
	}
	hist := []*models.RollupRuleView{
		testRollupRuleView(id, "name1", targets),
		testRollupRuleView(id, "name2", targets),
	}
	expected := rollupRuleHistoryJSON{
		RollupRules: []models.RollupRule{
			{
				ID:     id,
				Name:   "name1",
				Filter: "filter",
				Targets: []models.RollupTarget{
					{
						Name:     "target1",
						Tags:     []string{"tag"},
						Policies: []policy.Policy{},
					},
					{
						Name:     "target2",
						Tags:     []string{"tag"},
						Policies: []policy.Policy{},
					},
				},
				CutoverMillis:       0,
				LastUpdatedBy:       "",
				LastUpdatedAtMillis: 0,
			},
			{
				ID:     id,
				Name:   "name2",
				Filter: "filter",
				Targets: []models.RollupTarget{
					{
						Name:     "target1",
						Tags:     []string{"tag"},
						Policies: []policy.Policy{},
					},
					{
						Name:     "target2",
						Tags:     []string{"tag"},
						Policies: []policy.Policy{},
					},
				},
				CutoverMillis:       0,
				LastUpdatedBy:       "",
				LastUpdatedAtMillis: 0,
			},
		},
	}
	require.EqualValues(t, expected, newRollupRuleHistoryJSON(hist))
}

// Tests Setup
func testNamespaceView(name string) *models.NamespaceView {
	return &models.NamespaceView{
		Name:              name,
		ForRuleSetVersion: 1,
	}
}

func testNamespacesView(namespaceNames ...string) *models.NamespacesView {
	namespaces := make([]*models.NamespaceView, len(namespaceNames))
	for i, name := range namespaceNames {
		namespaces[i] = testNamespaceView(name)
	}
	return &models.NamespacesView{
		Version:    1,
		Namespaces: namespaces,
	}
}

func testMappingRuleView(id, name string) *models.MappingRuleView {
	return &models.MappingRuleView{
		ID:       id,
		Name:     name,
		Filter:   "filter",
		Policies: []policy.Policy{},
	}
}

func testRollupTargetView(name string) *models.RollupTargetView {
	return &models.RollupTargetView{
		Name:     name,
		Tags:     []string{"tag"},
		Policies: []policy.Policy{},
	}
}

func testRollupRuleView(id, name string, targets []models.RollupTargetView) *models.RollupRuleView {
	return &models.RollupRuleView{
		ID:      id,
		Name:    name,
		Filter:  "filter",
		Targets: targets,
	}
}
