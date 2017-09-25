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

func writeAPIResponse(w http.ResponseWriter, code int, msg string) error {
	j, err := json.Marshal(apiResponse{Code: code, Message: msg})
	if err != nil {
		return err
	}
	if err := newPayload(j, code).write(w); err != nil {
		return err
	}
	return nil
}

type r2Handler func(http.ResponseWriter, *http.Request) error

func (fn r2Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := fn(w, r); err != nil {
		handleError(w, err)
	}
}

func handleError(w http.ResponseWriter, opError error) {
	var err error
	switch opError.(type) {
	case ConflictError:
		err = writeAPIResponse(w, http.StatusConflict, opError.Error())
	case BadInputError:
		err = writeAPIResponse(w, http.StatusBadRequest, opError.Error())
	case VersionError:
		err = writeAPIResponse(w, http.StatusConflict, opError.Error())
	case NotFoundError:
		err = writeAPIResponse(w, http.StatusNotFound, opError.Error())
	default:
		err = writeAPIResponse(w, http.StatusInternalServerError, opError.Error())
	}

	// Getting here means that the error handling failed. Trying to convey what was supposed to happen.
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("Could not generate error response for: %s", opError.Error()),
			http.StatusInternalServerError,
		)
	}

}

func (s *service) write(w http.ResponseWriter, data interface{}) error {
	if j, err := json.Marshal(data); err == nil {
		return newPayload(j, http.StatusOK).write(w)
	}
	return writeAPIResponse(w, http.StatusInternalServerError, "Could not create response object")
}

func (s *service) fetchNamespaces(w http.ResponseWriter, _ *http.Request) error {
	names, v, err := s.store.FetchNamespaces()
	if err != nil {
		return writeAPIResponse(w, http.StatusInternalServerError, "Could not read namespaces")
	}

	nss := make([]namespaceJSON, len(names))
	for i, n := range names {
		nss[i] = namespaceJSON{ID: n}
	}
	return s.write(w, namespacesJSON{Version: v, Namespaces: nss})
}

func (s *service) fetchNamespace(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)

	rs, err := s.store.FetchRuleSet(vars["namespaceID"])
	if err != nil {
		return err
	}
	return s.write(w, rs.ruleSetJSON())
}

func (s *service) createNamespace(w http.ResponseWriter, r *http.Request) error {
	var n namespaceJSON
	if errs := parseRequest(&n, r.Body); len(errs) > 0 {
		return writeAPIResponse(w, http.StatusBadRequest, errs)
	}

	id, err := s.store.CreateNamespace(n.ID)
	if err != nil {
		return err
	}
	return writeAPIResponse(w, http.StatusCreated, fmt.Sprintf("Created namespace %s", id))
}

func (s *service) deleteNamespace(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)

	namespaceID := vars["namespaceID"]
	if err := s.store.DeleteNamespace(namespaceID); err != nil {
		return err
	}
	return writeAPIResponse(w, http.StatusOK, fmt.Sprintf("Deleted namespace %s", namespaceID))
}

func (s *service) fetchMappingRule(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)

	mr, err := s.store.FetchMappingRule(vars["namespaceID"], vars["ruleID"])
	if err != nil {
		return err
	}
	return s.write(w, newMappingRuleJSON(mr))
}

func (s *service) createMappingRule(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)

	var mrj mappingRuleJSON
	if errs := parseRequest(&mrj, r.Body); len(errs) > 0 {
		return writeAPIResponse(w, http.StatusBadRequest, errs)
	}

	mr, err := s.store.CreateMappingRule(
		vars["namespaceID"],
		mrj.mappingRuleView(),
	)

	if err != nil {
		return err
	}
	return s.write(w, newMappingRuleJSON(mr))
}

func (s *service) updateMappingRule(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	var err error

	var mrj mappingRuleJSON
	if errs := parseRequest(&mrj, r.Body); len(errs) > 0 {
		return writeAPIResponse(w, http.StatusBadRequest, errs)
	}
	mr, err := s.store.UpdateMappingRule(
		vars["namespaceID"],
		vars["ruleID"],
		mrj.mappingRuleView(),
	)

	if err != nil {
		return err
	}
	return s.write(w, newMappingRuleJSON(mr))
}

func (s *service) deleteMappingRule(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)

	namespaceID := vars["namespaceID"]
	mappingRuleID := vars["ruleID"]

	if err := s.store.DeleteMappingRule(namespaceID, mappingRuleID); err != nil {
		return err
	}

	return writeAPIResponse(w, http.StatusOK,
		fmt.Sprintf("Deleted mapping rule: %s in namespace %s", mappingRuleID, namespaceID))
}

func (s *service) fetchMappingRuleHistory(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)

	hist, err := s.store.FetchMappingRuleHistory(vars["namespaceID"], vars["ruleID"])
	if err != nil {
		return err
	}

	views := make([]mappingRuleJSON, len(hist))
	for i, mappingRule := range hist {
		x := mappingRule
		views[i] = newMappingRuleJSON(x)
	}
	return s.write(w, views)
}

func (s *service) fetchRollupRule(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)

	rr, err := s.store.FetchRollupRule(vars["namespaceID"], vars["ruleID"])
	if err != nil {
		return err
	}
	return s.write(w, newRollupRuleJSON(rr))
}

func (s *service) createRollupRule(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	namespaceID := vars["namespaceID"]

	var rrj rollupRuleJSON
	if errs := parseRequest(&rrj, r.Body); len(errs) > 0 {
		return writeAPIResponse(w, http.StatusBadRequest, errs)
	}

	rr, err := s.store.CreateRollupRule(namespaceID, rrj.rollupRuleView())
	if err != nil {
		return err
	}

	return s.write(w, newRollupRuleJSON(rr))
}

func (s *service) updateRollupRule(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)

	var rrj rollupRuleJSON
	if errs := parseRequest(&rrj, r.Body); len(errs) > 0 {
		return writeAPIResponse(w, http.StatusBadRequest, errs)
	}

	rr, err := s.store.UpdateRollupRule(
		vars["namespaceID"],
		vars["ruleID"],
		rrj.rollupRuleView(),
	)

	if err != nil {
		return err
	}
	return s.write(w, newRollupRuleJSON(rr))
}

func (s *service) deleteRollupRule(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)

	namespaceID := vars["namespaceID"]
	rollupRuleID := vars["ruleID"]
	if err := s.store.DeleteRollupRule(namespaceID, rollupRuleID); err != nil {
		return err
	}
	return writeAPIResponse(w, http.StatusOK,
		fmt.Sprintf("Deleted rollup rule: %s in namespace %s", rollupRuleID, namespaceID))
}

func (s *service) fetchRollupRuleHistory(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)

	hist, err := s.store.FetchRollupRuleHistory(
		vars["namespaceID"],
		vars["ruleID"],
	)

	if err != nil {
		return err
	}

	views := make([]rollupRuleJSON, len(hist))
	for i, rollupRule := range hist {
		x := rollupRule
		views[i] = newRollupRuleJSON(x)
	}
	return s.write(w, views)
}

// RegisterHandlers registers rule handler.
func (s *service) RegisterHandlers(router *mux.Router) {
	log := s.iOpts.Logger()
	// Namespaces action
	fullNamespacesPrefix := s.fullPath(namespacePath)
	router.Handle(fullNamespacesPrefix, r2Handler(s.fetchNamespaces)).Methods(http.MethodGet)
	router.Handle(fullNamespacesPrefix, r2Handler(s.createNamespace)).Methods(http.MethodPost)

	// Ruleset actions
	fullNamespacePath := s.fullPath(namespacePrefix)
	router.Handle(fullNamespacePath, r2Handler(s.fetchNamespace)).Methods(http.MethodGet)
	router.Handle(fullNamespacePath, r2Handler(s.deleteNamespace)).Methods(http.MethodDelete)

	// Mapping Rule actions
	router.Handle(s.fullPath(mappingRuleRoot), r2Handler(s.createMappingRule)).Methods(http.MethodPost)

	fullMappingRuleWithIDPath := s.fullPath(mappingRuleWithIDPath)
	router.Handle(fullMappingRuleWithIDPath, r2Handler(s.fetchMappingRule)).Methods(http.MethodGet)
	router.Handle(fullMappingRuleWithIDPath, r2Handler(s.updateMappingRule)).Methods(http.MethodPut, http.MethodPatch)
	router.Handle(fullMappingRuleWithIDPath, r2Handler(s.deleteMappingRule)).Methods(http.MethodDelete)

	// Mapping Rule history
	router.Handle(s.fullPath(mappingRuleHistoryPath), r2Handler(s.fetchMappingRuleHistory)).Methods(http.MethodGet)

	// Rollup Rule actions
	router.Handle(s.fullPath(rollupRuleRoot), r2Handler(s.createRollupRule)).Methods(http.MethodPost)

	fullRollupRuleWithIDPath := s.fullPath(rollupRuleWithIDPath)
	router.Handle(fullRollupRuleWithIDPath, r2Handler(s.fetchRollupRule)).Methods(http.MethodGet)
	router.Handle(fullRollupRuleWithIDPath, r2Handler(s.updateRollupRule)).Methods(http.MethodPut, http.MethodPatch)
	router.Handle(fullRollupRuleWithIDPath, r2Handler(s.deleteRollupRule)).Methods(http.MethodDelete)

	// Rollup Rule history
	router.Handle(s.fullPath(rollupRuleHistoryPath), r2Handler(s.fetchRollupRuleHistory)).Methods(http.MethodGet)

	log.Infof("Registered rules endpoints")
}

func (s *service) fullPath(path string) string {
	return fmt.Sprintf("/%s/%s", s.rootPrefix, path)
}
