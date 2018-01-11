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
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/m3db/m3ctl/auth"
	mservice "github.com/m3db/m3ctl/service"
	"github.com/m3db/m3metrics/policy"
	"github.com/m3db/m3metrics/rules"
	"github.com/m3db/m3x/clock"
	"github.com/m3db/m3x/instrument"

	"github.com/gorilla/mux"
	"github.com/pborman/uuid"
	"gopkg.in/go-playground/validator.v9"
)

const (
	namespacePath     = "/namespaces"
	mappingRulePrefix = "mapping-rules"
	rollupRulePrefix  = "rollup-rules"

	namespaceIDVar = "namespaceID"
	ruleIDVar      = "ruleID"

	nanosPerMilli = int64(time.Millisecond / time.Nanosecond)
)

var (
	namespacePrefix     = fmt.Sprintf("%s/{%s}", namespacePath, namespaceIDVar)
	validateRuleSetPath = fmt.Sprintf("%s/{%s}/ruleset/validate", namespacePath, namespaceIDVar)

	mappingRuleRoot        = fmt.Sprintf("%s/%s", namespacePrefix, mappingRulePrefix)
	mappingRuleWithIDPath  = fmt.Sprintf("%s/{%s}", mappingRuleRoot, ruleIDVar)
	mappingRuleHistoryPath = fmt.Sprintf("%s/history", mappingRuleWithIDPath)

	rollupRuleRoot        = fmt.Sprintf("%s/%s", namespacePrefix, rollupRulePrefix)
	rollupRuleWithIDPath  = fmt.Sprintf("%s/{%s}", rollupRuleRoot, ruleIDVar)
	rollupRuleHistoryPath = fmt.Sprintf("%s/history", rollupRuleWithIDPath)

	errNilRequest = errors.New("Nil request")
)

type route struct {
	path   string
	method string
}

var authorizationRegistry = map[route]auth.AuthorizationType{
	// This validation route should only require read access.
	{path: validateRuleSetPath, method: http.MethodPost}: auth.AuthorizationTypeReadOnly,
}

func sendResponse(w http.ResponseWriter, data []byte, status int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err := w.Write(data)
	return err
}

// TODO(dgromov): Make this return a list of validation errors
func parseRequest(s interface{}, body io.ReadCloser) error {
	if err := json.NewDecoder(body).Decode(s); err != nil {
		return NewBadInputError(fmt.Sprintf("Malformed Json: %s", err.Error()))
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
		return NewBadInputError(fmt.Sprintf("Required: [%v]", strings.Join(required, ", ")))
	}
	return nil
}

func writeAPIResponse(w http.ResponseWriter, code int, msg string) error {
	j, err := json.Marshal(apiResponse{Code: code, Message: msg})
	if err != nil {
		return err
	}
	return sendResponse(w, j, code)
}

type r2HandlerFunc func(http.ResponseWriter, *http.Request) error

type r2Handler struct {
	iOpts instrument.Options
	auth  auth.HTTPAuthService
}

func (h r2Handler) wrap(authType auth.AuthorizationType, fn r2HandlerFunc) http.Handler {
	f := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			h.handleError(w, err)
		}
	})
	return h.auth.NewAuthHandler(authType, f, writeAPIResponse)
}

func (h r2Handler) handleError(w http.ResponseWriter, opError error) {
	h.iOpts.Logger().Errorf(opError.Error())

	var err error
	switch opError.(type) {
	case conflictError:
		err = writeAPIResponse(w, http.StatusConflict, opError.Error())
	case badInputError:
		err = writeAPIResponse(w, http.StatusBadRequest, opError.Error())
	case versionError:
		err = writeAPIResponse(w, http.StatusConflict, opError.Error())
	case notFoundError:
		err = writeAPIResponse(w, http.StatusNotFound, opError.Error())
	case authError:
		err = writeAPIResponse(w, http.StatusUnauthorized, opError.Error())
	default:
		err = writeAPIResponse(w, http.StatusInternalServerError, opError.Error())
	}

	// Getting here means that the error handling failed. Trying to convey what was supposed to happen.
	if err != nil {
		msg := fmt.Sprintf("Could not generate error response for: %s", opError.Error())
		h.iOpts.Logger().Errorf(msg)
		http.Error(w, msg, http.StatusInternalServerError)
	}
}

func defaultAuthorizationTypeForHTTPMethod(method string) auth.AuthorizationType {
	switch method {
	case http.MethodGet:
		return auth.AuthorizationTypeReadOnly
	case http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch:
		return auth.AuthorizationTypeReadWrite
	default:
		return auth.AuthorizationTypeUnknown
	}
}

func registerRoute(router *mux.Router, path, method string, h r2Handler, hf r2HandlerFunc) error {
	authType, exists := authorizationRegistry[route{path: path, method: method}]
	if !exists {
		var err error
		if authType = defaultAuthorizationTypeForHTTPMethod(method); authType == auth.AuthorizationTypeUnknown {
			return err
		}
	}
	fn := h.wrap(authType, hf)
	router.Handle(path, fn).Methods(method)
	return nil
}

// service handles all of the endpoints for r2.
type service struct {
	rootPrefix  string
	store       Store
	authService auth.HTTPAuthService
	iOpts       instrument.Options
	nowFn       clock.NowFn
	metrics     serviceMetrics
}

// NewService creates a new r2 service using a given store.
func NewService(
	rootPrefix string,
	authService auth.HTTPAuthService,
	store Store,
	iOpts instrument.Options,
	clockOpts clock.Options,
) mservice.Service {
	return &service{
		rootPrefix:  rootPrefix,
		store:       store,
		authService: authService,
		iOpts:       iOpts,
		nowFn:       clockOpts.NowFn(),
		metrics:     newServiceMetrics(iOpts.MetricsScope()),
	}
}

func (s *service) URLPrefix() string { return s.rootPrefix }

func (s *service) RegisterHandlers(router *mux.Router) error {
	log := s.iOpts.Logger()
	h := r2Handler{s.iOpts, s.authService}

	// Namespaces actions
	if err := registerRoute(router, namespacePath, http.MethodGet, h, s.fetchNamespaces); err != nil {
		return err
	}
	if err := registerRoute(router, namespacePath, http.MethodPost, h, s.createNamespace); err != nil {
		return err
	}

	// Ruleset actions
	if err := registerRoute(router, namespacePrefix, http.MethodGet, h, s.fetchNamespace); err != nil {
		return err
	}
	if err := registerRoute(router, namespacePrefix, http.MethodDelete, h, s.deleteNamespace); err != nil {
		return err
	}
	if err := registerRoute(router, validateRuleSetPath, http.MethodPost, h, s.validateNamespace); err != nil {
		return err
	}

	// Mapping Rule actions
	if err := registerRoute(router, mappingRuleRoot, http.MethodPost, h, s.createMappingRule); err != nil {
		return err
	}

	if err := registerRoute(router, mappingRuleWithIDPath, http.MethodGet, h, s.fetchMappingRule); err != nil {
		return err
	}
	if err := registerRoute(router, mappingRuleWithIDPath, http.MethodPut, h, s.updateMappingRule); err != nil {
		return err
	}
	if err := registerRoute(router, mappingRuleWithIDPath, http.MethodDelete, h, s.deleteMappingRule); err != nil {
		return err
	}

	// Mapping Rule history
	if err := registerRoute(router, mappingRuleHistoryPath, http.MethodGet, h, s.fetchMappingRuleHistory); err != nil {
		return err
	}

	// Rollup Rule actions
	if err := registerRoute(router, rollupRuleRoot, http.MethodPost, h, s.createRollupRule); err != nil {
		return err
	}

	if err := registerRoute(router, rollupRuleWithIDPath, http.MethodGet, h, s.fetchRollupRule); err != nil {
		return err
	}
	if err := registerRoute(router, rollupRuleWithIDPath, http.MethodPut, h, s.updateRollupRule); err != nil {
		return err
	}
	if err := registerRoute(router, rollupRuleWithIDPath, http.MethodDelete, h, s.deleteRollupRule); err != nil {
		return err
	}

	if err := registerRoute(router, rollupRuleHistoryPath, http.MethodGet, h, s.fetchRollupRuleHistory); err != nil {
		return err
	}

	log.Infof("Registered rules endpoints")
	return nil
}

func (s *service) Close() { s.store.Close() }

type routeFunc func(s *service, r *http.Request) (data interface{}, err error)

func (s *service) handleRoute(rf routeFunc, r *http.Request, namespace string) (interface{}, error) {
	if r == nil {
		return nil, errNilRequest
	}
	start := s.nowFn()
	data, err := rf(s, r)
	s.metrics.recordMetric(r.RequestURI, r.Method, namespace, time.Since(start), err)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *service) sendResponse(w http.ResponseWriter, statusCode int, data interface{}) error {
	if j, err := json.Marshal(data); err == nil {
		return sendResponse(w, j, statusCode)
	}
	return writeAPIResponse(w, http.StatusInternalServerError, "could not create response object")
}

func (s *service) fetchNamespaces(w http.ResponseWriter, r *http.Request) error {
	data, err := s.handleRoute(fetchNamespaces, r, "")
	if err != nil {
		return err
	}
	return s.sendResponse(w, http.StatusOK, data)
}

func (s *service) fetchNamespace(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	data, err := s.handleRoute(fetchNamespace, r, vars[namespaceIDVar])
	if err != nil {
		return err
	}
	return s.sendResponse(w, http.StatusOK, data)
}

func (s *service) createNamespace(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	data, err := s.handleRoute(createNamespace, r, vars[namespaceIDVar])
	if err != nil {
		return err
	}
	return s.sendResponse(w, http.StatusCreated, data)
}

func (s *service) validateNamespace(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	data, err := s.handleRoute(validateRuleSet, r, vars[namespaceIDVar])
	if err != nil {
		return err
	}
	return writeAPIResponse(w, http.StatusOK, data.(string))
}

func (s *service) deleteNamespace(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	data, err := s.handleRoute(deleteNamespace, r, vars[namespaceIDVar])
	if err != nil {
		return err
	}
	return writeAPIResponse(w, http.StatusOK, data.(string))
}

func (s *service) fetchMappingRule(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	data, err := s.handleRoute(fetchMappingRule, r, vars[namespaceIDVar])
	if err != nil {
		return err
	}
	return s.sendResponse(w, http.StatusOK, data)
}

func (s *service) createMappingRule(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	data, err := s.handleRoute(createMappingRule, r, vars[namespaceIDVar])
	if err != nil {
		return err
	}
	return s.sendResponse(w, http.StatusCreated, data)
}

func (s *service) updateMappingRule(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	data, err := s.handleRoute(updateMappingRule, r, vars[namespaceIDVar])
	if err != nil {
		return err
	}
	return s.sendResponse(w, http.StatusOK, data)
}

func (s *service) deleteMappingRule(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	data, err := s.handleRoute(deleteMappingRule, r, vars[namespaceIDVar])
	if err != nil {
		return err
	}
	return writeAPIResponse(w, http.StatusOK, data.(string))
}

func (s *service) fetchMappingRuleHistory(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	data, err := s.handleRoute(fetchMappingRuleHistory, r, vars[namespaceIDVar])
	if err != nil {
		return err
	}
	return s.sendResponse(w, http.StatusOK, data)
}

func (s *service) fetchRollupRule(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	data, err := s.handleRoute(fetchRollupRule, r, vars[namespaceIDVar])
	if err != nil {
		return err
	}
	return s.sendResponse(w, http.StatusOK, data)
}

func (s *service) createRollupRule(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	data, err := s.handleRoute(createRollupRule, r, vars[namespaceIDVar])
	if err != nil {
		return err
	}
	return s.sendResponse(w, http.StatusCreated, data)
}

func (s *service) updateRollupRule(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	data, err := s.handleRoute(updateRollupRule, r, vars[namespaceIDVar])
	if err != nil {
		return err
	}
	return s.sendResponse(w, http.StatusOK, data)
}

func (s *service) deleteRollupRule(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	data, err := s.handleRoute(deleteRollupRule, r, vars[namespaceIDVar])
	if err != nil {
		return err
	}
	return writeAPIResponse(w, http.StatusOK, data.(string))
}

func (s *service) fetchRollupRuleHistory(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	data, err := s.handleRoute(fetchRollupRuleHistory, r, vars[namespaceIDVar])
	if err != nil {
		return err
	}
	return s.sendResponse(w, http.StatusOK, data)
}

func (s *service) newUpdateOptions(r *http.Request) (UpdateOptions, error) {
	uOpts := NewUpdateOptions()
	author, err := s.authService.GetUser(r.Context())
	if err != nil {
		return uOpts, nil
	}
	return uOpts.SetAuthor(author), nil
}

func newRuleSetJSON(latest *rules.RuleSetSnapshot) ruleSetJSON {
	var mrJSON []mappingRuleJSON
	for _, m := range latest.MappingRules {
		mrJSON = append(mrJSON, newMappingRuleJSON(m))
	}
	var rrJSON []rollupRuleJSON
	for _, r := range latest.RollupRules {
		rrJSON = append(rrJSON, newRollupRuleJSON(r))
	}
	return ruleSetJSON{
		Namespace:     latest.Namespace,
		Version:       latest.Version,
		CutoverMillis: latest.CutoverNanos / nanosPerMilli,
		MappingRules:  mrJSON,
		RollupRules:   rrJSON,
	}
}

type apiResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type namespaceJSON struct {
	ID                string `json:"id" validate:"required"`
	ForRuleSetVersion int    `json:"forRuleSetVersion"`
}

func newNamespaceJSON(nv *rules.NamespaceView) namespaceJSON {
	return namespaceJSON{
		ID:                nv.Name,
		ForRuleSetVersion: nv.ForRuleSetVersion,
	}
}

type namespacesJSON struct {
	Version    int             `json:"version"`
	Namespaces []namespaceJSON `json:"namespaces"`
}

func newNamespacesJSON(nss *rules.NamespacesView) namespacesJSON {
	views := make([]namespaceJSON, len(nss.Namespaces))
	for i, namespace := range nss.Namespaces {
		views[i] = newNamespaceJSON(namespace)
	}
	return namespacesJSON{
		Version:    nss.Version,
		Namespaces: views,
	}
}

type mappingRuleJSON struct {
	ID                  string          `json:"id,omitempty"`
	Name                string          `json:"name" validate:"required"`
	CutoverMillis       int64           `json:"cutoverMillis,omitempty"`
	Filter              string          `json:"filter" validate:"required"`
	Policies            []policy.Policy `json:"policies" validate:"required"`
	LastUpdatedBy       string          `json:"lastUpdatedBy"`
	LastUpdatedAtMillis int64           `json:"lastUpdatedAtMillis"`
}

func (m mappingRuleJSON) mappingRuleView() *rules.MappingRuleView {
	return &rules.MappingRuleView{
		ID:       m.ID,
		Name:     m.Name,
		Filter:   m.Filter,
		Policies: m.Policies,
	}
}

func newMappingRuleJSON(mrv *rules.MappingRuleView) mappingRuleJSON {
	return mappingRuleJSON{
		ID:                  mrv.ID,
		Name:                mrv.Name,
		Filter:              mrv.Filter,
		Policies:            mrv.Policies,
		CutoverMillis:       mrv.CutoverNanos / nanosPerMilli,
		LastUpdatedBy:       mrv.LastUpdatedBy,
		LastUpdatedAtMillis: mrv.LastUpdatedAtNanos / nanosPerMilli,
	}
}

type mappingRuleHistoryJSON struct {
	MappingRules []mappingRuleJSON `json:"mappingRules"`
}

func newMappingRuleHistoryJSON(hist []*rules.MappingRuleView) mappingRuleHistoryJSON {
	views := make([]mappingRuleJSON, len(hist))
	for i, mappingRule := range hist {
		views[i] = newMappingRuleJSON(mappingRule)
	}
	return mappingRuleHistoryJSON{MappingRules: views}
}

type rollupTargetJSON struct {
	Name     string          `json:"name" validate:"required"`
	Tags     []string        `json:"tags" validate:"required"`
	Policies []policy.Policy `json:"policies" validate:"required"`
}

func (t rollupTargetJSON) rollupTargetView() rules.RollupTargetView {
	return rules.RollupTargetView{
		Name:     t.Name,
		Tags:     t.Tags,
		Policies: t.Policies,
	}
}

func newRollupTargetJSON(t rules.RollupTargetView) rollupTargetJSON {
	return rollupTargetJSON{
		Name:     t.Name,
		Tags:     t.Tags,
		Policies: t.Policies,
	}
}

type rollupRuleJSON struct {
	ID                  string             `json:"id,omitempty"`
	Name                string             `json:"name" validate:"required"`
	Filter              string             `json:"filter" validate:"required"`
	Targets             []rollupTargetJSON `json:"targets" validate:"required,dive,required"`
	CutoverMillis       int64              `json:"cutoverMillis,omitempty"`
	LastUpdatedBy       string             `json:"lastUpdatedBy"`
	LastUpdatedAtMillis int64              `json:"lastUpdatedAtMillis"`
}

func newRollupRuleJSON(rrv *rules.RollupRuleView) rollupRuleJSON {
	targets := make([]rollupTargetJSON, len(rrv.Targets))
	for i, t := range rrv.Targets {
		targets[i] = newRollupTargetJSON(t)
	}
	return rollupRuleJSON{
		ID:                  rrv.ID,
		Name:                rrv.Name,
		Filter:              rrv.Filter,
		Targets:             targets,
		CutoverMillis:       rrv.CutoverNanos / nanosPerMilli,
		LastUpdatedBy:       rrv.LastUpdatedBy,
		LastUpdatedAtMillis: rrv.LastUpdatedAtNanos / nanosPerMilli,
	}
}

func (r rollupRuleJSON) rollupRuleView() *rules.RollupRuleView {
	targets := make([]rules.RollupTargetView, len(r.Targets))
	for i, t := range r.Targets {
		targets[i] = t.rollupTargetView()
	}

	return &rules.RollupRuleView{
		ID:      r.ID,
		Name:    r.Name,
		Filter:  r.Filter,
		Targets: targets,
	}
}

type rollupRuleHistoryJSON struct {
	RollupRules []rollupRuleJSON `json:"rollupRules"`
}

func newRollupRuleHistoryJSON(hist []*rules.RollupRuleView) rollupRuleHistoryJSON {
	views := make([]rollupRuleJSON, len(hist))
	for i, rollupRule := range hist {
		views[i] = newRollupRuleJSON(rollupRule)
	}
	return rollupRuleHistoryJSON{RollupRules: views}
}

type ruleSetJSON struct {
	Namespace     string            `json:"id"`
	Version       int               `json:"version"`
	CutoverMillis int64             `json:"cutoverMillis"`
	MappingRules  []mappingRuleJSON `json:"mappingRules"`
	RollupRules   []rollupRuleJSON  `json:"rollupRules"`
}

type idGenType int

const (
	generateID idGenType = iota
	dontGenerateID
)

// ruleSetSnapshot create a RuleSetSnapshot from a rulesetJSON. If the ruleSetJSON has no IDs
// for any of its mapping rules or rollup rules, it generates missing IDs and sets as a string UUID
// string so they can be stored in a mapping (id -> rule).
func (r ruleSetJSON) ruleSetSnapshot(idGenType idGenType) (*rules.RuleSetSnapshot, error) {
	mappingRules := make(map[string]*rules.MappingRuleView, len(r.MappingRules))
	for _, mr := range r.MappingRules {
		id := mr.ID
		if id == "" {
			if idGenType == dontGenerateID {
				return nil, fmt.Errorf("can't convert rulesetJSON to ruleSetSnapshot, no mapping rule id for %v", mr)
			}
			id = uuid.New()
			mr.ID = id
		}
		mappingRules[id] = mr.mappingRuleView()
	}

	rollupRules := make(map[string]*rules.RollupRuleView, len(r.RollupRules))
	for _, rr := range r.RollupRules {
		id := rr.ID
		if id == "" {
			if idGenType == dontGenerateID {
				return nil, fmt.Errorf("can't convert rulesetJSON to ruleSetSnapshot, no rollup rule id for %v", rr)
			}
			id = uuid.New()
			rr.ID = id
		}
		rollupRules[id] = rr.rollupRuleView()
	}

	return &rules.RuleSetSnapshot{
		Namespace:    r.Namespace,
		Version:      r.Version,
		MappingRules: mappingRules,
		RollupRules:  rollupRules,
	}, nil
}
