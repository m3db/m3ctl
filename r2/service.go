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
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/m3db/m3ctl/services/r2ctl/server"
	"github.com/m3db/m3x/instrument"

	"github.com/gorilla/mux"
)

const (
	namespacePrefix   = "namespace/{namespaceID}"
	mappingRulePrefix = "mappings"
	rollupRulePrefix  = "rollups"
	namespacePath     = "namespaces"
)

var (
	mappingRuleRoot        = fmt.Sprintf("%s/%s", namespacePrefix, mappingRulePrefix)
	mappingRuleWithIDPath  = fmt.Sprintf("%s/{ruleID}", mappingRuleRoot)
	mappingRuleHistoryPath = fmt.Sprintf("%s/history", mappingRuleWithIDPath)

	rollupRuleRoot        = fmt.Sprintf("%s/%s", namespacePrefix, rollupRulePrefix)
	rollupRuleWithIDPath  = fmt.Sprintf("%s/{ruleID}", rollupRuleRoot)
	rollupRuleHistoryPath = fmt.Sprintf("%s/history", rollupRuleWithIDPath)
)

// service handles all of the endpoints for r2.
type service struct {
	iOpts      instrument.Options
	rootPrefix string
	store      Store
}

// NewService creates a new r2 service using a given store.
func NewService(iOpts instrument.Options, rootPrefix string, store Store) server.Service {
	return &service{iOpts: iOpts, store: store, rootPrefix: rootPrefix}
}

type payload struct {
	data   []byte
	status int
}

func newPayload(d []byte, status int) payload {
	return payload{data: d, status: status}
}

func (p payload) write(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(p.status)
	w.Write(p.data)
}

func writeAPIResponse(w http.ResponseWriter, code int, msg string) {
	j, _ := json.Marshal(apiResponse{Code: code, Message: msg})
	newPayload(j, code).write(w)
}

func handleError(w http.ResponseWriter, opError error) {
	switch opError.(type) {
	case ConflictError:
		writeAPIResponse(w, http.StatusConflict, opError.Error())
	case BadInputError:
		writeAPIResponse(w, http.StatusBadRequest, opError.Error())
	case VersionError:
		writeAPIResponse(w, http.StatusConflict, opError.Error())
	case NotFoundError:
		writeAPIResponse(w, http.StatusNotFound, opError.Error())
	default:
		writeAPIResponse(w, http.StatusInternalServerError, opError.Error())
	}
}

func (s *service) readNamespaces(w http.ResponseWriter, _ *http.Request) {
	nss, v, err := s.store.FetchNamespaces()
	if err != nil {
		writeAPIResponse(w, http.StatusInternalServerError, "Could not read namespaces")
		return
	}

	j, err := json.Marshal(namespaces{
		Version:    v,
		Namespaces: nss,
	})

	if err != nil {
		writeAPIResponse(w, http.StatusInternalServerError, "Could not create response object")
		return
	}

	newPayload(j, http.StatusOK).write(w)
}

func (s *service) readNamespace(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespaceID, exists := vars["namespaceID"]
	if !exists {
		writeAPIResponse(w, http.StatusBadRequest, "Could not parse namespaceId")
		return
	}

	rs, err := s.store.FetchRuleSet(namespaceID)
	if err != nil {
		handleError(w, err)
		return
	}

	j, err := json.Marshal(newRuleSetView(*rs))
	if err != nil {
		writeAPIResponse(w, http.StatusInternalServerError, "Could not create response object")
		return
	}

	newPayload(j, http.StatusOK).write(w)
}

func (s *service) createNamespace(w http.ResponseWriter, r *http.Request) {
	var cnb createNamespaceData
	if err := json.NewDecoder(r.Body).Decode(&cnb); err != nil {
		writeAPIResponse(w, http.StatusBadRequest, "Unable to parse new namespace data.")
		return
	}

	nsd := NewNamespaceData(cnb.ID, cnb.NamespacesVersion)
	newName, err := s.store.CreateNamespace(nsd)
	if err == nil {
		writeAPIResponse(w, http.StatusCreated, fmt.Sprintf("Created namespace %s", newName))
	} else {
		handleError(w, err)
	}
}

func (s *service) deleteNamespace(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespaceID, exists := vars["namespaceID"]
	if !exists {
		writeAPIResponse(w, http.StatusBadRequest, "Could not parse namespaceId")
		return
	}

	var cdn deleteNamespaceData
	if err := json.NewDecoder(r.Body).Decode(&cdn); err != nil {
		writeAPIResponse(w, http.StatusBadRequest, "Unable to parse update data.")
		return
	}
	nsd := NewNamespaceData(namespaceID, cdn.NamespacesVersion)
	err := s.store.DeleteNamespace(nsd)
	if err == nil {
		writeAPIResponse(w, http.StatusOK, fmt.Sprintf("Deleted namespace %s", namespaceID))
	} else {
		handleError(w, err)
	}
}

func (s *service) readMappingRule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespaceID, exists := vars["namespaceID"]
	if !exists {
		writeAPIResponse(w, http.StatusBadRequest, "Could not parse namespaceID")
		return
	}

	mappingRuleID, exists := vars["ruleID"]
	if !exists {
		writeAPIResponse(w, http.StatusBadRequest, "Could not parse mappingRuleID")
		return
	}

	mappingRule, err := s.store.FetchMappingRule(namespaceID, mappingRuleID)
	if err != nil {
		handleError(w, err)
		return
	}
	view := newMappingRuleView(*mappingRule)

	j, err := json.Marshal(view)
	if err != nil {
		writeAPIResponse(w, http.StatusInternalServerError, "Could not create response object")
		return
	}

	newPayload(j, http.StatusOK).write(w)
}

func (s *service) createMappingRule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespaceID, exists := vars["namespaceID"]
	if !exists {
		writeAPIResponse(w, http.StatusBadRequest, "Could not parse namespaceID")
		return
	}

	var mrd MappingRuleData
	if err := json.NewDecoder(r.Body).Decode(&mrd); err != nil {
		writeAPIResponse(w, http.StatusBadRequest, "Unable to parse new mapping rule data.")
		return
	}

	mappingRule, err := s.store.CreateMappingRule(namespaceID, mrd)
	if err != nil {
		handleError(w, err)
		return
	}
	s.iOpts.Logger().Infof("%v", mappingRule.RawFilters)
	view := newMappingRuleView(*mappingRule)

	j, err := json.Marshal(view)
	if err != nil {
		writeAPIResponse(w, http.StatusInternalServerError, "Could not create response object")
		return
	}

	newPayload(j, http.StatusOK).write(w)
}

func (s *service) updateMappingRule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespaceID, exists := vars["namespaceID"]
	if !exists {
		writeAPIResponse(w, http.StatusBadRequest, "Could not parse namespaceID")
		return
	}

	mappingRuleID, exists := vars["ruleID"]
	if !exists {
		writeAPIResponse(w, http.StatusBadRequest, "Could not parse mappingRuleID")
		return
	}
	var mrd MappingRuleData
	if err := json.NewDecoder(r.Body).Decode(&mrd); err != nil {
		writeAPIResponse(w, http.StatusBadRequest, "Unable to parse new mapping rule data.")
		return
	}

	err := s.store.UpdateMappingRule(namespaceID, mappingRuleID, mrd)
	if err == nil {
		writeAPIResponse(w, http.StatusOK,
			fmt.Sprintf("Updated mapping rule: %s in namespace %s", mappingRuleID, namespaceID))
	} else {
		handleError(w, err)
	}
}

func (s *service) deleteMappingRule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespaceID, exists := vars["namespaceID"]
	if !exists {
		writeAPIResponse(w, http.StatusBadRequest, "Could not parse namespaceID")
		return
	}

	mappingRuleID, exists := vars["ruleID"]
	if !exists {
		writeAPIResponse(w, http.StatusBadRequest, "Could not parse mappingRuleID")
		return
	}

	err := s.store.DeleteMappingRule(namespaceID, mappingRuleID, 0)
	if err == nil {
		writeAPIResponse(w, http.StatusOK,
			fmt.Sprintf("Deleted mapping rule: %s in namespace %s", mappingRuleID, namespaceID))
	} else {
		handleError(w, err)
	}
}

func (s *service) readMappingRuleHistory(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, "Placeholder for readMappingRuleHistory")
}

func (s *service) readRollupRule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespaceID, exists := vars["namespaceID"]
	if !exists {
		writeAPIResponse(w, http.StatusBadRequest, "Could not parse namespaceID")
		return
	}

	rollupRuleID, exists := vars["ruleID"]
	if !exists {
		writeAPIResponse(w, http.StatusBadRequest, "Could not parse rollupRuleID")
		return
	}

	rollupRule, err := s.store.FetchRollupRule(namespaceID, rollupRuleID)
	if err != nil {
		handleError(w, err)
		return
	}
	view := newRollupRuleView(*rollupRule)

	j, err := json.Marshal(view)
	if err != nil {
		writeAPIResponse(w, http.StatusInternalServerError, "Could not create response object")
		return
	}

	newPayload(j, http.StatusOK).write(w)
}

func (s *service) createRollupRule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespaceID, exists := vars["namespaceID"]
	if !exists {
		writeAPIResponse(w, http.StatusBadRequest, "Could not parse namespaceID")
		return
	}

	var rrd RollupRuleData
	if err := json.NewDecoder(r.Body).Decode(&rrd); err != nil {
		writeAPIResponse(w, http.StatusBadRequest, "Unable to parse new rollup rule data.")
		return
	}

	rollupRule, err := s.store.CreateRollupRule(namespaceID, rrd)

	if err != nil {
		handleError(w, err)
		return
	}

	view := newRollupRuleView(*rollupRule)

	j, err := json.Marshal(view)
	if err != nil {
		writeAPIResponse(w, http.StatusInternalServerError, "Could not create response object")
		return
	}
	newPayload(j, http.StatusOK).write(w)
}

func (s *service) updateRollupRule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespaceID, exists := vars["namespaceID"]
	if !exists {
		writeAPIResponse(w, http.StatusBadRequest, "Could not parse namespaceID")
		return
	}

	rollupRuleID, exists := vars["ruleID"]
	if !exists {
		writeAPIResponse(w, http.StatusBadRequest, "Could not parse rollupRuleID")
		return
	}

	var rrd RollupRuleData
	if err := json.NewDecoder(r.Body).Decode(&rrd); err != nil {
		writeAPIResponse(w, http.StatusBadRequest, "Unable to parse new rollup rule data.")
		return
	}
	err := s.store.UpdateRollupRule(namespaceID, rollupRuleID, rrd)

	if err == nil {
		writeAPIResponse(w, http.StatusOK, fmt.Sprintf("Updated rollup rule: %s in namespace %s", rollupRuleID, namespaceID))
	} else {
		handleError(w, err)
	}
}

func (s *service) deleteRollupRule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespaceID, exists := vars["namespaceID"]
	if !exists {
		writeAPIResponse(w, http.StatusBadRequest, "Could not parse namespaceID")
		return
	}

	rollupRuleID, exists := vars["ruleID"]
	if !exists {
		writeAPIResponse(w, http.StatusBadRequest, "Could not parse rollupRuleID")
		return
	}

	err := s.store.DeleteRollupRule(namespaceID, rollupRuleID, 0)

	if err == nil {
		writeAPIResponse(w, http.StatusOK, fmt.Sprintf("Deleted rollup rule: %s in namespace %s", rollupRuleID, namespaceID))
	} else {
		handleError(w, err)
	}
}

func (s *service) readRollupRuleHistory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespaceID, exists := vars["namespaceID"]
	if !exists {
		writeAPIResponse(w, http.StatusBadRequest, "Could not parse namespaceID")
		return
	}

	rollupRuleID, exists := vars["ruleID"]
	if !exists {
		writeAPIResponse(w, http.StatusBadRequest, "Could not parse rollupRuleID")
		return
	}

	hist, err := s.store.FetchRollupRuleHistory(namespaceID, rollupRuleID)
	if err != nil {
		handleError(w, err)
		return
	}
	views := make([]rollupRuleView, len(hist))
	for i, rollupRule := range hist {
		x := rollupRule
		views[i] = newRollupRuleView(*x)
	}

	j, err := json.Marshal(views)
	if err != nil {
		writeAPIResponse(w, http.StatusInternalServerError, "Could not create response object")
		return
	}

	newPayload(j, http.StatusOK).write(w)
}

// RegisterHandlers registers rule handler.
func (s *service) RegisterHandlers(router *mux.Router) {
	log := s.iOpts.Logger()
	// Namespaces action
	fullNamespacesPrefix := s.fullPath("namespaces")
	router.HandleFunc(fullNamespacesPrefix, s.readNamespaces).Methods(http.MethodGet)
	router.HandleFunc(fullNamespacesPrefix, s.createNamespace).Methods(http.MethodPost)

	// Ruleset actions
	fullNamespacePath := s.fullPath(namespacePrefix)
	router.HandleFunc(fullNamespacePath, s.readNamespace).Methods(http.MethodGet)
	router.HandleFunc(fullNamespacePath, s.deleteNamespace).Methods(http.MethodDelete)

	// Mapping Rule actions
	router.HandleFunc(s.fullPath(mappingRuleRoot), s.createMappingRule).Methods(http.MethodPost)

	fullMappingRuleWithIDPath := s.fullPath(mappingRuleWithIDPath)
	router.HandleFunc(fullMappingRuleWithIDPath, s.readMappingRule).Methods(http.MethodGet)
	router.HandleFunc(fullMappingRuleWithIDPath, s.updateMappingRule).Methods(http.MethodPut, http.MethodPatch)
	router.HandleFunc(fullMappingRuleWithIDPath, s.deleteMappingRule).Methods(http.MethodDelete)

	// Mapping Rule history
	router.HandleFunc(s.fullPath(mappingRuleHistoryPath), s.readMappingRuleHistory).Methods(http.MethodGet)

	// Rollup Rule actions
	router.HandleFunc(s.fullPath(rollupRuleRoot), s.createRollupRule).Methods(http.MethodPost)

	fullRollupRuleWithIDPath := s.fullPath(rollupRuleWithIDPath)
	router.HandleFunc(fullRollupRuleWithIDPath, s.readRollupRule).Methods(http.MethodGet)
	router.HandleFunc(fullRollupRuleWithIDPath, s.updateRollupRule).Methods(http.MethodPut, http.MethodPatch)
	router.HandleFunc(fullRollupRuleWithIDPath, s.deleteRollupRule).Methods(http.MethodDelete)

	// Rollup Rule history
	router.HandleFunc(s.fullPath(rollupRuleHistoryPath), s.readRollupRuleHistory).Methods(http.MethodGet)

	log.Infof("Registered rules endpoints")
}

func (s *service) fullPath(path string) string {
	return fmt.Sprintf("/%s/%s", s.rootPrefix, path)
}
