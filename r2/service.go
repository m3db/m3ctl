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
	"io"
	"net/http"
	"reflect"
	"strings"

	"github.com/m3db/m3ctl/services/r2ctl/server"
	"github.com/m3db/m3x/instrument"

	"github.com/gorilla/mux"
	"gopkg.in/go-playground/validator.v9"
)

const (
	namespacePrefix   = "namespaces/{namespaceID}"
	mappingRulePrefix = "mapping-rules"
	rollupRulePrefix  = "rollup-rules"
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

func (p payload) write(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(p.status)
	if _, err := w.Write(p.data); err != nil {
		return fmt.Errorf("cannot write http payload. %s", err.Error())
	}
	return nil
}

func writeAPIResponse(w http.ResponseWriter, code int, msg string) error {
	j, err := json.Marshal(apiResponse{Code: code, Message: msg})
	if err != nil {
		panic(fmt.Sprintf("cannot create response for [%d]: %s", code, msg))
	}
	return newPayload(j, code).write(w)
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

// TODO(dgromov): Make this return a list of validation errors
func parseRequest(s interface{}, body io.ReadCloser) string {
	if err := json.NewDecoder(body).Decode(s); err != nil {
		return fmt.Sprintf("Malformed Json: %s", err.Error())
	}

	// Invoking the validation explictely to have control over the format of the error output.
	validate := validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		parts := strings.SplitN(fld.Tag.Get("json"), ",", 2)
		if len(parts) > 0 {
			return parts[0]
		}
		return fld.Name
	})

	var required []string
	if err := validate.Struct(s); err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			if e.ActualTag() == "required" {
				required = append(required, e.Namespace())
			}
		}
	}

	if len(required) > 0 {
		return fmt.Sprintf("Required: [%v]", strings.Join(required, ", "))
	}
	return ""
}

func (s *service) write(w http.ResponseWriter, data interface{}) {
	if j, err := json.Marshal(data); err != nil {
		writeAPIResponse(w, http.StatusInternalServerError, "Could not create response object")
	} else {
		newPayload(j, http.StatusOK).write(w)
	}
}

func (s *service) fetchNamespaces(w http.ResponseWriter, _ *http.Request) {
	names, v, err := s.store.FetchNamespaces()
	if err != nil {
		writeAPIResponse(w, http.StatusInternalServerError, "Could not read namespaces")
		return
	}
	nss := make([]namespaceJSON, len(names))
	for i, n := range names {
		nss[i] = namespaceJSON{ID: n}
	}
	s.write(w, namespacesJSON{Version: v, Namespaces: nss})
}

func (s *service) fetchNamespace(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if rs, err := s.store.FetchRuleSet(vars["namespaceID"]); err == nil {
		s.write(w, rs.ruleSetJSON())
	} else {
		handleError(w, err)
	}
}

func (s *service) createNamespace(w http.ResponseWriter, r *http.Request) {
	var n namespaceJSON
	errs := parseRequest(&n, r.Body)
	if len(errs) > 0 {
		writeAPIResponse(w, http.StatusBadRequest, errs)
		return
	}

	id, err := s.store.CreateNamespace(n.ID)
	if err == nil {
		writeAPIResponse(w, http.StatusCreated, fmt.Sprintf("Created namespace %s", id))
	} else {
		handleError(w, err)
	}
}

func (s *service) deleteNamespace(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	namespaceID := vars["namespaceID"]
	if err := s.store.DeleteNamespace(namespaceID); err == nil {
		writeAPIResponse(w, http.StatusOK, fmt.Sprintf("Deleted namespace %s", namespaceID))
	} else {
		handleError(w, err)
	}
}

func (s *service) fetchMappingRule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if mr, err := s.store.FetchMappingRule(vars["namespaceID"], vars["ruleID"]); err == nil {
		s.write(w, newMappingRuleJSON(mr))
	} else {
		handleError(w, err)
	}
}

func (s *service) createMappingRule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var mrj mappingRuleJSON
	if errs := parseRequest(&mrj, r.Body); len(errs) > 0 {
		writeAPIResponse(w, http.StatusBadRequest, errs)
	} else {
		if mr, err := s.store.CreateMappingRule(vars["namespaceID"], mrj.mappingRuleView()); err == nil {
			s.write(w, newMappingRuleJSON(mr))
		} else {
			handleError(w, err)
		}
	}
}

func (s *service) updateMappingRule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var mrj mappingRuleJSON
	if errs := parseRequest(&mrj, r.Body); len(errs) > 0 {
		writeAPIResponse(w, http.StatusBadRequest, errs)
	} else {
		if mr, err := s.store.UpdateMappingRule(vars["namespaceID"], vars["ruleID"], mrj.mappingRuleView()); err == nil {
			s.write(w, newMappingRuleJSON(mr))
		} else {
			handleError(w, err)
		}
	}
}

func (s *service) deleteMappingRule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	namespaceID := vars["namespaceID"]
	mappingRuleID := vars["ruleID"]
	if err := s.store.DeleteMappingRule(namespaceID, mappingRuleID); err == nil {
		writeAPIResponse(w, http.StatusOK,
			fmt.Sprintf("Deleted mapping rule: %s in namespace %s", mappingRuleID, namespaceID))
	} else {
		handleError(w, err)
	}
}

func (s *service) fetchMappingRuleHistory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if hist, err := s.store.FetchMappingRuleHistory(vars["namespaceID"], vars["ruleID"]); err == nil {
		views := make([]mappingRuleJSON, len(hist))
		for i, mappingRule := range hist {
			x := mappingRule
			views[i] = newMappingRuleJSON(x)
		}
		s.write(w, views)
	} else {
		handleError(w, err)
	}
}

func (s *service) fetchRollupRule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if rr, err := s.store.FetchRollupRule(vars["namespaceID"], vars["ruleID"]); err == nil {
		s.write(w, newRollupRuleJSON(rr))
	} else {
		handleError(w, err)
	}
}

func (s *service) createRollupRule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespaceID := vars["namespaceID"]

	var rrj rollupRuleJSON
	if errs := parseRequest(&rrj, r.Body); len(errs) > 0 {
		writeAPIResponse(w, http.StatusBadRequest, errs)
	} else {
		if rr, err := s.store.CreateRollupRule(namespaceID, rrj.rollupRuleView()); err == nil {
			s.write(w, newRollupRuleJSON(rr))
		} else {
			handleError(w, err)
		}
	}
}

func (s *service) updateRollupRule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var rrj rollupRuleJSON
	if errs := parseRequest(&rrj, r.Body); len(errs) > 0 {
		writeAPIResponse(w, http.StatusBadRequest, errs)
	} else {
		if rr, err := s.store.UpdateRollupRule(vars["namespaceID"], vars["ruleID"], rrj.rollupRuleView()); err == nil {
			s.write(w, newRollupRuleJSON(rr))
		} else {
			handleError(w, err)
		}
	}
}

func (s *service) deleteRollupRule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	namespaceID := vars["namespaceID"]
	rollupRuleID := vars["ruleID"]
	if err := s.store.DeleteRollupRule(namespaceID, rollupRuleID); err == nil {
		writeAPIResponse(w, http.StatusOK,
			fmt.Sprintf("Deleted rollup rule: %s in namespace %s", rollupRuleID, namespaceID))
	} else {
		handleError(w, err)
	}
}

func (s *service) fetchRollupRuleHistory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if hist, err := s.store.FetchRollupRuleHistory(vars["namespaceID"], vars["ruleID"]); err == nil {
		views := make([]rollupRuleJSON, len(hist))
		for i, rollupRule := range hist {
			x := rollupRule
			views[i] = newRollupRuleJSON(x)
		}
		s.write(w, views)
	} else {
		handleError(w, err)
	}
}

// RegisterHandlers registers rule handler.
func (s *service) RegisterHandlers(router *mux.Router) {
	log := s.iOpts.Logger()
	// Namespaces action
	fullNamespacesPrefix := s.fullPath(namespacePath)
	router.HandleFunc(fullNamespacesPrefix, s.fetchNamespaces).Methods(http.MethodGet)
	router.HandleFunc(fullNamespacesPrefix, s.createNamespace).Methods(http.MethodPost)

	// Ruleset actions
	fullNamespacePath := s.fullPath(namespacePrefix)
	router.HandleFunc(fullNamespacePath, s.fetchNamespace).Methods(http.MethodGet)
	router.HandleFunc(fullNamespacePath, s.deleteNamespace).Methods(http.MethodDelete)

	// Mapping Rule actions
	router.HandleFunc(s.fullPath(mappingRuleRoot), s.createMappingRule).Methods(http.MethodPost)

	fullMappingRuleWithIDPath := s.fullPath(mappingRuleWithIDPath)
	router.HandleFunc(fullMappingRuleWithIDPath, s.fetchMappingRule).Methods(http.MethodGet)
	router.HandleFunc(fullMappingRuleWithIDPath, s.updateMappingRule).Methods(http.MethodPut, http.MethodPatch)
	router.HandleFunc(fullMappingRuleWithIDPath, s.deleteMappingRule).Methods(http.MethodDelete)

	// Mapping Rule history
	router.HandleFunc(s.fullPath(mappingRuleHistoryPath), s.fetchMappingRuleHistory).Methods(http.MethodGet)

	// Rollup Rule actions
	router.HandleFunc(s.fullPath(rollupRuleRoot), s.createRollupRule).Methods(http.MethodPost)

	fullRollupRuleWithIDPath := s.fullPath(rollupRuleWithIDPath)
	router.HandleFunc(fullRollupRuleWithIDPath, s.fetchRollupRule).Methods(http.MethodGet)
	router.HandleFunc(fullRollupRuleWithIDPath, s.updateRollupRule).Methods(http.MethodPut, http.MethodPatch)
	router.HandleFunc(fullRollupRuleWithIDPath, s.deleteRollupRule).Methods(http.MethodDelete)

	// Rollup Rule history
	router.HandleFunc(s.fullPath(rollupRuleHistoryPath), s.fetchRollupRuleHistory).Methods(http.MethodGet)

	log.Infof("Registered rules endpoints")
}

func (s *service) fullPath(path string) string {
	return fmt.Sprintf("/%s/%s", s.rootPrefix, path)
}
