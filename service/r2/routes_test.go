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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/m3db/m3metrics/rules"
	"github.com/m3db/m3metrics/rules/models/changes"

	"github.com/m3db/m3ctl/auth"
	"github.com/m3db/m3ctl/service/r2/mocks"
	"github.com/m3db/m3ctl/service/r2/store"
	"github.com/m3db/m3metrics/rules/models"
	"github.com/m3db/m3x/clock"
	"github.com/m3db/m3x/instrument"
	"github.com/uber-go/tally"

	"github.com/stretchr/testify/require"
)

func TestHandleRoute(t *testing.T) {
	s := newTestService(nil)
	r := newTestGetRequest()
	expected := models.NewNamespaces(&models.NamespacesView{})
	actual, err := s.handleRoute(fetchNamespaces, r, newTestInstrumentMethodMetrics())
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestHandleRouteNilRequest(t *testing.T) {
	s := newTestService(nil)
	_, err := s.handleRoute(fetchNamespaces, nil, newTestInstrumentMethodMetrics())
	require.EqualError(t, err, errNilRequest.Error())
}
func TestFetchNamespacesSuccess(t *testing.T) {
	expected := models.NewNamespaces(&models.NamespacesView{})
	actual, err := fetchNamespaces(newTestService(nil), newTestGetRequest())
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestFetchNamespaceSuccess(t *testing.T) {
	expected := models.NewRuleSet(&models.RuleSetSnapshotView{})
	actual, err := fetchNamespace(newTestService(nil), newTestGetRequest())
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestValidateRuleSetSuccess(t *testing.T) {
	expected := "Ruleset is valid"
	actual, err := validateRuleSet(newTestService(nil), newTestPostRequest([]byte(`{}`)))
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestCreateNamespaceSuccess(t *testing.T) {
	expected := models.NewNamespace(&models.NamespaceView{})
	actual, err := createNamespace(newTestService(nil), newTestPostRequest([]byte(`{"id": "id"}`)))
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestDeleteNamespaceSuccess(t *testing.T) {
	expected := fmt.Sprintf("Deleted namespace %s", "")
	actual, err := deleteNamespace(newTestService(nil), newTestDeleteRequest())
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestFetchMappingRuleSuccess(t *testing.T) {
	expected := models.NewMappingRule(&models.MappingRuleView{})
	actual, err := fetchMappingRule(newTestService(nil), newTestGetRequest())
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestCreateMappingRuleSuccess(t *testing.T) {
	expected := models.NewMappingRule(&models.MappingRuleView{})
	actual, err := createMappingRule(newTestService(nil), newTestPostRequest(
		[]byte(`{"filter": "key:val", "name": "name", "policies": []}`),
	))
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestUpdateMappingRuleSuccess(t *testing.T) {
	expected := models.NewMappingRule(&models.MappingRuleView{})
	actual, err := updateMappingRule(newTestService(nil), newTestPutRequest(
		[]byte(`{"filter": "key:val", "name": "name", "policies": []}`),
	))
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestDeleteMappingRuleSuccess(t *testing.T) {
	expected := fmt.Sprintf("Deleted mapping rule: %s in namespace %s", "", "")
	actual, err := deleteMappingRule(newTestService(nil), newTestDeleteRequest())
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestFetchMappingRuleHistorySuccess(t *testing.T) {
	expected := models.NewMappingRuleSnapshots(make([]*models.MappingRuleView, 0))
	actual, err := fetchMappingRuleHistory(newTestService(nil), newTestGetRequest())
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestFetchRollupRuleSuccess(t *testing.T) {
	expected := models.NewRollupRule(&models.RollupRuleView{})
	actual, err := fetchRollupRule(newTestService(nil), newTestGetRequest())
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestCreateRollupRuleSuccess(t *testing.T) {
	expected := models.NewRollupRule(&models.RollupRuleView{})
	actual, err := createRollupRule(newTestService(nil), newTestPostRequest(
		[]byte(`{"filter": "key:val", "name": "name", "targets": []}`),
	))
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestUpdateRollupRuleSuccess(t *testing.T) {
	expected := models.NewRollupRule(&models.RollupRuleView{})
	actual, err := updateRollupRule(newTestService(nil), newTestPutRequest(
		[]byte(`{"filter": "key:val", "name": "name", "targets": []}`),
	))
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestDeleteRollupRuleSuccess(t *testing.T) {
	expected := fmt.Sprintf("Deleted rollup rule: %s in namespace %s", "", "")
	actual, err := deleteRollupRule(newTestService(nil), newTestDeleteRequest())
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestFetchRollupRuleHistorySuccess(t *testing.T) {
	expected := models.NewRollupRuleSnapshots([]*models.RollupRuleView{})
	actual, err := fetchRollupRuleHistory(newTestService(nil), newTestGetRequest())
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestBulkUpdateRuleSet(t *testing.T) {
	namespaceID := "testNamespace"
	originalRuleSet := newRuleSet(namespaceID, 1)
	rrvs, _ := originalRuleSet.RollupRules()
	var rrIDs []string
	for id := range rrvs {
		rrIDs = append(rrIDs, id)
	}
	mrvs, _ := originalRuleSet.MappingRules()
	var mrIDs []string
	for id := range mrvs {
		mrIDs = append(mrIDs, id)
	}

	bulkReqBody := newBulkReqBody()
	bulkRequestWithUpdates(&bulkReqBody, rrIDs[0], mrIDs[0])
	bulkRequestWithDeletes(&bulkReqBody, rrIDs[1], mrIDs[1])
	bodyBytes, err := json.Marshal(bulkReqBody)
	require.NoError(t, err)
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("/namespaces/%s/ruleset/bulk", namespaceID),
		bytes.NewBuffer(bodyBytes),
	)
	require.NoError(t, err)
	req = mux.SetURLVars(
		req,
		map[string]string{
			"namespaceID": namespaceID,
		},
	)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	storeMock := mocks.NewMockStore(ctrl)
	storeMock.EXPECT().FetchRuleSet(namespaceID).Return(originalRuleSet, nil)
	storeMock.EXPECT().UpdateRuleSet(gomock.Any()).Do(func(mrs rules.MutableRuleSet) {
		latest, err := mrs.Latest()
		if err != nil {
			t.Fail()
		}
		rs := models.NewRuleSet(latest)
		rs.Sort()

		actualRRNames := make([]string, 0, len(rs.RollupRules))
		for _, rr := range rs.RollupRules {
			actualRRNames = append(actualRRNames, rr.Name)
		}
		actualMRNames := make([]string, 0, len(rs.MappingRules))
		for _, rr := range rs.MappingRules {
			actualMRNames = append(actualMRNames, rr.Name)
		}

		expectedRRNames := []string{
			"rollupRule3",
			"updatedRollupRule",
		}
		expectedMRNames := []string{
			"mappingRule3",
			"updatedMappingRule",
		}
		require.Equal(t, expectedMRNames, actualMRNames)
		require.Equal(t, expectedRRNames, actualRRNames)
	}).Return(newRuleSet(namespaceID, 2), nil)

	service := newTestService(storeMock)
	resp, err := bulkUpdateRuleSet(service, req)
	require.NoError(t, err)
	typedResp := resp.(models.RuleSet)
	require.Equal(t, typedResp.Version, 2)
}

func TestBulkUpdateRuleSetVersionMismatch(t *testing.T) {
	namespaceID := "testNamespace"
	bulkReqBody := newBulkReqBody()
	bodyBytes, err := json.Marshal(bulkReqBody)
	require.NoError(t, err)
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("/namespaces/%s/ruleset/bulk", namespaceID),
		bytes.NewBuffer(bodyBytes),
	)
	require.NoError(t, err)
	req = mux.SetURLVars(
		req,
		map[string]string{
			"namespaceID": namespaceID,
		},
	)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	storeMock := mocks.NewMockStore(ctrl)
	storeMock.EXPECT().FetchRuleSet(namespaceID).Return(newRuleSet(namespaceID, 4), nil)

	service := newTestService(storeMock)
	_, err = bulkUpdateRuleSet(service, req)
	require.Error(t, err)
	require.IsType(t, NewConflictError(""), err)
}

func TestBulkUpdateRuleSetKVFetchFailure(t *testing.T) {
	namespaceID := "testNamespace"
	bulkReqBody := newBulkReqBody()
	bodyBytes, err := json.Marshal(bulkReqBody)
	require.NoError(t, err)
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("/namespaces/%s/ruleset/bulk", namespaceID),
		bytes.NewBuffer(bodyBytes),
	)
	require.NoError(t, err)
	req = mux.SetURLVars(
		req,
		map[string]string{
			"namespaceID": namespaceID,
		},
	)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	storeMock := mocks.NewMockStore(ctrl)
	storeMock.EXPECT().FetchRuleSet(namespaceID).Return(nil, NewInternalError("KV error"))
	service := newTestService(storeMock)
	_, err = bulkUpdateRuleSet(service, req)
	require.Error(t, err)
	require.IsType(t, NewInternalError(""), err)
}

func TestBulkUpdateRuleSetKVUpdateFailure(t *testing.T) {
	namespaceID := "testNamespace"
	bulkReqBody := newBulkReqBody()
	bodyBytes, err := json.Marshal(bulkReqBody)
	require.NoError(t, err)
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("/namespaces/%s/ruleset/bulk", namespaceID),
		bytes.NewBuffer(bodyBytes),
	)
	require.NoError(t, err)
	req = mux.SetURLVars(
		req,
		map[string]string{
			"namespaceID": namespaceID,
		},
	)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	storeMock := mocks.NewMockStore(ctrl)
	storeMock.EXPECT().FetchRuleSet(namespaceID).Return(newRuleSet(namespaceID, 1), nil).Times(1)
	storeMock.EXPECT().UpdateRuleSet(gomock.Any()).Return(nil, NewInternalError("KV error"))

	service := newTestService(storeMock)
	_, err = bulkUpdateRuleSet(service, req)
	require.Error(t, err)
	require.IsType(t, NewInternalError(""), err)
}

func TestBulkUpdateRuleSetBadInput(t *testing.T) {
	namespaceID := "testNamespace"
	originalRuleSet := newRuleSet(namespaceID, 1)
	rrvs, _ := originalRuleSet.RollupRules()
	var rrIDs []string
	for id := range rrvs {
		rrIDs = append(rrIDs, id)
	}
	mrvs, _ := originalRuleSet.MappingRules()
	var mrIDs []string
	for id := range mrvs {
		mrIDs = append(mrIDs, id)
	}

	bulkReqBody := newBulkReqBody()
	bulkRequestWithUpdates(&bulkReqBody, rrIDs[0], mrIDs[0])
	bodyBytes, err := json.Marshal(bulkReqBody)
	require.NoError(t, err)
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("/namespaces/%s/ruleset/bulk", namespaceID),
		bytes.NewBuffer(bodyBytes),
	)
	require.NoError(t, err)
	req = mux.SetURLVars(
		req,
		map[string]string{
			"namespaceID": namespaceID,
		},
	)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	storeMock := mocks.NewMockStore(ctrl)
	storeMock.EXPECT().FetchRuleSet(namespaceID).Return(newRuleSet(namespaceID, 1), nil).Times(1)

	service := newTestService(storeMock)
	_, err = bulkUpdateRuleSet(service, req)
	require.Error(t, err)
	require.IsType(t, NewBadInputError(""), err)
}

func TestApplyChangesToRuleSetUpdateRuleFailure(t *testing.T) {
	originalRuleSet := newRuleSet("validNamepspace", 1)
	bulkReqBody := newBulkReqBody()
	bulkRequestWithUpdates(&bulkReqBody, "invalidRRID", "")
	changes := bulkReqBody.RuleSetChanges

	_, err := applyChangesToRuleSet(
		changes,
		originalRuleSet.ToMutableRuleSet(),
		store.NewUpdateOptions(),
		rules.NewRuleSetUpdateHelper(time.Minute),
	)
	require.Error(t, err)
	require.Equal(t, "cannot update rule invalidRRID: rule not found", err.Error())

	originalRuleSet = newRuleSet("validNamepspace", 1)
	bulkReqBody = newBulkReqBody()
	bulkRequestWithUpdates(&bulkReqBody, "", "invalidMRID")
	changes = bulkReqBody.RuleSetChanges
	_, err = applyChangesToRuleSet(
		changes,
		originalRuleSet.ToMutableRuleSet(),
		store.NewUpdateOptions(),
		rules.NewRuleSetUpdateHelper(time.Minute),
	)
	require.Error(t, err)
	require.Equal(t, "cannot update rule invalidMRID: rule not found", err.Error())
}

func TestApplyChangesToRuleSetDeleteRuleFailure(t *testing.T) {
	originalRuleSet := newRuleSet("validNamepspace", 1)
	bulkReqBody := newBulkReqBody()
	bulkRequestWithDeletes(&bulkReqBody, "invalidRRID", "")
	changes := bulkReqBody.RuleSetChanges

	_, err := applyChangesToRuleSet(
		changes,
		originalRuleSet.ToMutableRuleSet(),
		store.NewUpdateOptions(),
		rules.NewRuleSetUpdateHelper(time.Minute),
	)
	require.Error(t, err)
	require.Equal(t, "cannot delete rule invalidRRID: rule not found", err.Error())

	originalRuleSet = newRuleSet("validNamepspace", 1)
	bulkReqBody = newBulkReqBody()
	bulkRequestWithDeletes(&bulkReqBody, "", "invalidMRID")
	changes = bulkReqBody.RuleSetChanges
	_, err = applyChangesToRuleSet(
		changes,
		originalRuleSet.ToMutableRuleSet(),
		store.NewUpdateOptions(),
		rules.NewRuleSetUpdateHelper(time.Minute),
	)
	require.Error(t, err)
	require.Equal(t, "cannot delete rule invalidMRID: rule not found", err.Error())
}

func TestApplyChangesToRuleSetNoChanges(t *testing.T) {
	namespaceID := "testNamespace"

	bulkReqBody := bulkRequest{
		RuleSetVersion: 1,
	}
	bodyBytes, err := json.Marshal(bulkReqBody)
	require.NoError(t, err)
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("/namespaces/%s/ruleset/bulk", namespaceID),
		bytes.NewBuffer(bodyBytes),
	)
	require.NoError(t, err)
	req = mux.SetURLVars(
		req,
		map[string]string{
			"namespaceID": namespaceID,
		},
	)

	service := newTestService(nil)
	_, err = bulkUpdateRuleSet(service, req)
	require.Error(t, err)
	require.IsType(t, NewBadInputError(""), err)
}

func TestBulkUpdateRuleSetMappingRuleMissingOp(t *testing.T) {
	namespaceID := "testNamespace"
	bulkReqBody := bulkRequest{
		RuleSetVersion: 1,
		RuleSetChanges: changes.RuleSetChanges{
			MappingRuleChanges: []changes.MappingRuleChange{
				changes.MappingRuleChange{},
			},
		},
	}
	bodyBytes, err := json.Marshal(bulkReqBody)
	require.NoError(t, err)
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("/namespaces/%s/ruleset/bulk", namespaceID),
		bytes.NewBuffer(bodyBytes),
	)
	require.NoError(t, err)
	req = mux.SetURLVars(
		req,
		map[string]string{
			"namespaceID": namespaceID,
		},
	)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	storeMock := mocks.NewMockStore(ctrl)
	storeMock.EXPECT().FetchRuleSet(namespaceID).Return(newRuleSet(namespaceID, 1), nil).Times(1)

	service := newTestService(storeMock)
	_, err = bulkUpdateRuleSet(service, req)
	require.Error(t, err)
	require.IsType(t, NewBadInputError(""), err)
}

func TestBulkUpdateRuleSetRollupRuleMissingOp(t *testing.T) {
	namespaceID := "testNamespace"
	bulkReqBody := bulkRequest{
		RuleSetVersion: 1,
		RuleSetChanges: changes.RuleSetChanges{
			RollupRuleChanges: []changes.RollupRuleChange{
				changes.RollupRuleChange{},
			},
		},
	}
	bodyBytes, err := json.Marshal(bulkReqBody)
	require.NoError(t, err)
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("/namespaces/%s/ruleset/bulk", namespaceID),
		bytes.NewBuffer(bodyBytes),
	)
	require.NoError(t, err)
	req = mux.SetURLVars(
		req,
		map[string]string{
			"namespaceID": namespaceID,
		},
	)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	storeMock := mocks.NewMockStore(ctrl)
	storeMock.EXPECT().FetchRuleSet(namespaceID).Return(newRuleSet(namespaceID, 1), nil).Times(1)

	service := newTestService(storeMock)
	_, err = bulkUpdateRuleSet(service, req)
	require.Error(t, err)
	require.IsType(t, NewBadInputError(""), err)
}

func newTestService(store store.Store) *service {
	if store == nil {
		store = newMockStore()
	}
	iOpts := instrument.NewOptions()
	return &service{
		metrics:     newServiceMetrics(iOpts.MetricsScope(), iOpts.MetricsSamplingRate()),
		nowFn:       clock.NewOptions().NowFn(),
		store:       store,
		authService: auth.NewNoopAuth(),
		logger:      iOpts.Logger(),
	}
}

func newTestGetRequest() *http.Request {
	req, _ := http.NewRequest("GET", "/route", nil)
	return req.WithContext(context.Background())
}

func newTestPostRequest(bodyBuff []byte) *http.Request {
	req, _ := http.NewRequest("POST", "/route", bytes.NewReader(bodyBuff))
	return req.WithContext(context.Background())
}

func newTestPutRequest(bodyBuff []byte) *http.Request {
	req, _ := http.NewRequest("PUT", "/route", bytes.NewReader(bodyBuff))
	return req.WithContext(context.Background())
}

func newTestDeleteRequest() *http.Request {
	req, _ := http.NewRequest("DELETE", "/route", nil)
	return req.WithContext(context.Background())
}

func newTestInstrumentMethodMetrics() instrument.MethodMetrics {
	return instrument.NewMethodMetrics(tally.NoopScope, "testRoute", 1.0)
}

func ptr(s string) *string {
	return &s
}

func newBulkReqBody() bulkRequest {
	return bulkRequest{
		RuleSetVersion: 1,
		RuleSetChanges: changes.RuleSetChanges{
			Namespace: "testNamespace",
			RollupRuleChanges: []changes.RollupRuleChange{
				changes.RollupRuleChange{
					Op: changes.AddOp,
					RuleData: &models.RollupRule{
						Name: "rollupRule3",
					},
				},
			},
			MappingRuleChanges: []changes.MappingRuleChange{
				changes.MappingRuleChange{
					Op: changes.AddOp,
					RuleData: &models.MappingRule{
						Name: "mappingRule3",
					},
				},
			},
		},
	}
}

func bulkRequestWithUpdates(req *bulkRequest, rrIDToUpdate, mrIDToUpdate string) {
	if rrIDToUpdate != "" {
		req.RuleSetChanges.RollupRuleChanges = append(
			req.RuleSetChanges.RollupRuleChanges,
			changes.RollupRuleChange{
				Op:     changes.ChangeOp,
				RuleID: ptr(rrIDToUpdate),
				RuleData: &models.RollupRule{
					ID:   rrIDToUpdate,
					Name: "updatedRollupRule",
				},
			},
		)
	}

	if mrIDToUpdate != "" {
		req.RuleSetChanges.MappingRuleChanges = append(
			req.RuleSetChanges.MappingRuleChanges,
			changes.MappingRuleChange{
				Op:     changes.ChangeOp,
				RuleID: ptr(mrIDToUpdate),
				RuleData: &models.MappingRule{
					ID:   mrIDToUpdate,
					Name: "updatedMappingRule",
				},
			},
		)
	}
}

func bulkRequestWithDeletes(req *bulkRequest, rrIDToDelete, mrIDToDelete string) {
	if rrIDToDelete != "" {
		req.RuleSetChanges.RollupRuleChanges = append(
			req.RuleSetChanges.RollupRuleChanges,
			changes.RollupRuleChange{
				Op:     changes.DeleteOp,
				RuleID: ptr(rrIDToDelete),
			},
		)
	}

	if mrIDToDelete != "" {
		req.RuleSetChanges.MappingRuleChanges = append(
			req.RuleSetChanges.MappingRuleChanges,
			changes.MappingRuleChange{
				Op:     changes.DeleteOp,
				RuleID: ptr(mrIDToDelete),
			},
		)
	}
}

func newRuleSet(nsID string, version int) rules.RuleSet {
	helper := rules.NewRuleSetUpdateHelper(time.Minute)
	// For testing all updates happen at the 0 epoch
	meta := helper.NewUpdateMetadata(0, "originalAuthor")

	mrs := rules.NewEmptyRuleSet(nsID, rules.UpdateMetadata{})
	mrs.AddRollupRule(
		models.RollupRuleView{
			Name: "rollupRule1",
		},
		meta,
	)
	mrs.AddRollupRule(
		models.RollupRuleView{
			Name: "rollupRule2",
		},
		meta,
	)
	mrs.AddMappingRule(
		models.MappingRuleView{
			Name: "mappingRule1",
		},
		meta,
	)
	mrs.AddMappingRule(
		models.MappingRuleView{
			Name: "mappingRule2",
		},
		meta,
	)
	schema, _ := mrs.Schema()
	rs, _ := rules.NewRuleSetFromSchema(version, schema, rules.NewOptions())
	return rs
}

type mockStore struct{}

func newMockStore() store.Store {
	return mockStore{}
}

func (s mockStore) FetchNamespaces() (*models.NamespacesView, error) {
	return &models.NamespacesView{}, nil
}

func (s mockStore) ValidateRuleSet(rs *models.RuleSetSnapshotView) error {
	return nil
}

func (s mockStore) CreateNamespace(namespaceID string, uOpts store.UpdateOptions) (*models.NamespaceView, error) {
	return &models.NamespaceView{}, nil
}

func (s mockStore) DeleteNamespace(namespaceID string, uOpts store.UpdateOptions) error {
	return nil
}

func (s mockStore) FetchRuleSet(namespaceID string) (rules.RuleSet, error) {
	return newRuleSet("testNamespace", 1), nil
}

func (s mockStore) FetchRuleSetSnapshot(namespaceID string) (*models.RuleSetSnapshotView, error) {
	return &models.RuleSetSnapshotView{}, nil
}

func (s mockStore) FetchMappingRule(namespaceID, mappingRuleID string) (*models.MappingRuleView, error) {
	return &models.MappingRuleView{}, nil
}

func (s mockStore) CreateMappingRule(namespaceID string, mrv *models.MappingRuleView, uOpts store.UpdateOptions) (*models.MappingRuleView, error) {
	return &models.MappingRuleView{}, nil
}

func (s mockStore) UpdateMappingRule(namespaceID, mappingRuleID string, mrv *models.MappingRuleView, uOpts store.UpdateOptions) (*models.MappingRuleView, error) {
	return &models.MappingRuleView{}, nil
}

func (s mockStore) DeleteMappingRule(namespaceID, mappingRuleID string, uOpts store.UpdateOptions) error {
	return nil
}

func (s mockStore) FetchMappingRuleHistory(namespaceID, mappingRuleID string) ([]*models.MappingRuleView, error) {
	return make([]*models.MappingRuleView, 0), nil
}

func (s mockStore) FetchRollupRule(namespaceID, rollupRuleID string) (*models.RollupRuleView, error) {
	return &models.RollupRuleView{}, nil
}

func (s mockStore) CreateRollupRule(namespaceID string, rrv *models.RollupRuleView, uOpts store.UpdateOptions) (*models.RollupRuleView, error) {
	return &models.RollupRuleView{}, nil
}

func (s mockStore) UpdateRollupRule(namespaceID, rollupRuleID string, rrv *models.RollupRuleView, uOpts store.UpdateOptions) (*models.RollupRuleView, error) {
	return &models.RollupRuleView{}, nil
}

func (s mockStore) DeleteRollupRule(namespaceID, rollupRuleID string, uOpts store.UpdateOptions) error {
	return nil
}

func (s mockStore) FetchRollupRuleHistory(namespaceID, rollupRuleID string) ([]*models.RollupRuleView, error) {
	return make([]*models.RollupRuleView, 0), nil
}

func (s mockStore) UpdateRuleSet(rs rules.MutableRuleSet) (rules.RuleSet, error) {
	return s.FetchRuleSet(string(rs.Namespace()))
}

func (s mockStore) Close() {}
