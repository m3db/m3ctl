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
	"fmt"
	"net/http"

	"github.com/m3db/m3ctl/handler"
	"github.com/m3db/m3metrics/rules"
	"github.com/m3db/m3x/instrument"

	"github.com/gorilla/mux"
)

const (
	namespacePrefix   = "namespace/{namespaceId}"
	mappingRulePrefix = "mapping"
	rollupRulePrefix  = "rollup"
	namespacePath     = "namespaces"
)

var (
	mappingRuleRoot        = fmt.Sprintf("%s/%s", namespacePrefix, mappingRulePrefix)
	mappingRuleWithIDPath  = fmt.Sprintf("%s/{ruleId}", mappingRuleRoot)
	mappingRuleHistoryPath = fmt.Sprintf("%s/history", mappingRuleWithIDPath)

	rollupRuleRoot        = fmt.Sprintf("%s/%s", namespacePrefix, rollupRulePrefix)
	rollupRuleWithIDPath  = fmt.Sprintf("%s/{ruleId}", rollupRuleRoot)
	rollupRuleHistoryPath = fmt.Sprintf("%s/history", rollupRuleWithIDPath)
)

// Controller sets up the handlers for rules
type r2Controller struct {
	iOpts      instrument.Options
	api        API
	rootPrefix string
}

// APIConfig is all the configuration state needed by every api
type APIConfig struct {
	UpdateHelper rules.RuleSetUpdateHelper
}

// NewController creates a new r2Controller instance.
func NewController(iOpts instrument.Options, api API, rootPrefix string) handler.Controller {
	return r2Controller{iOpts: iOpts, api: api, rootPrefix: rootPrefix}
}

// API exposes the http functions available to any implementation of an r2api.
type API interface {
	// ReadNamespaces fetches all the namespaces.
	ReadNamespaces(http.ResponseWriter, *http.Request)

	// ReadNamespace fetches the current state of a namespace.
	ReadNamespace(http.ResponseWriter, *http.Request)

	// CreateNamespace creates a namespace and new ruleset and persists it.
	CreateNamespace(http.ResponseWriter, *http.Request)

	// DeleteNamespace deletes a namespace and its underlying ruleset.
	DeleteNamespace(http.ResponseWriter, *http.Request)

	// ReadMappingRule fetches the current state of a mapping rule.
	ReadMappingRule(http.ResponseWriter, *http.Request)

	// CreateMappingRule creates or revives a mapping rule.
	CreateMappingRule(http.ResponseWriter, *http.Request)

	// UpdateMappingRule updates a mapping rule to a new state.
	UpdateMappingRule(http.ResponseWriter, *http.Request)

	// DeleteMappingRule deletes a mapping rule.
	DeleteMappingRule(http.ResponseWriter, *http.Request)

	// ReadMappingRuleHistory fetches all of the snapshots for a given mapping rule.
	ReadMappingRuleHistory(http.ResponseWriter, *http.Request)

	// ReadRollupRule fetches the current state of a rollup rule.
	ReadRollupRule(http.ResponseWriter, *http.Request)

	// CreateRollupRule creates or revives a rollup rule.
	CreateRollupRule(http.ResponseWriter, *http.Request)

	// UpdateRollupRule updates a rollup rule to a new state.
	UpdateRollupRule(http.ResponseWriter, *http.Request)

	// DeleteRollupRule deletes a rollup rule.
	DeleteRollupRule(http.ResponseWriter, *http.Request)

	// ReadRollupRuleHistory fetches all of the snapshots for a given rollup rule.
	ReadRollupRuleHistory(http.ResponseWriter, *http.Request)
}

// RegisterHandlers registers rule handlers
func (c r2Controller) RegisterHandlers(router *mux.Router) {
	log := c.iOpts.Logger()
	// Get all namespaces
	router.HandleFunc(c.fullPath(namespacePrefix), c.api.ReadNamespaces).Methods(http.MethodGet)

	fullNamespacePath := c.fullPath(namespacePath)
	// Ruleset actions
	router.HandleFunc(fullNamespacePath, c.api.ReadNamespace).Methods(http.MethodGet)
	router.HandleFunc(fullNamespacePath, c.api.CreateNamespace).Methods(http.MethodPost)
	router.HandleFunc(fullNamespacePath, c.api.DeleteNamespace).Methods(http.MethodDelete)

	// Mapping Rule actions
	router.HandleFunc(c.fullPath(mappingRuleRoot), c.api.CreateMappingRule).Methods(http.MethodPost)

	fullMappingRuleWithIDPath := c.fullPath(mappingRuleWithIDPath)
	router.HandleFunc(fullMappingRuleWithIDPath, c.api.ReadMappingRule).Methods(http.MethodGet)
	router.HandleFunc(fullMappingRuleWithIDPath, c.api.UpdateMappingRule).Methods(http.MethodPut, http.MethodPatch)
	router.HandleFunc(fullMappingRuleWithIDPath, c.api.DeleteMappingRule).Methods(http.MethodDelete)

	// Mapping Rule history
	router.HandleFunc(c.fullPath(mappingRuleHistoryPath), c.api.ReadMappingRuleHistory).Methods(http.MethodGet)

	// Rollup Rule actions
	router.HandleFunc(c.fullPath(rollupRuleRoot), c.api.CreateRollupRule).Methods(http.MethodPost)

	fullRollupRuleWithIDPath := c.fullPath(rollupRuleWithIDPath)
	router.HandleFunc(fullRollupRuleWithIDPath, c.api.ReadRollupRule).Methods(http.MethodGet)
	router.HandleFunc(fullRollupRuleWithIDPath, c.api.UpdateRollupRule).Methods(http.MethodPut, http.MethodPatch)
	router.HandleFunc(fullRollupRuleWithIDPath, c.api.DeleteRollupRule).Methods(http.MethodDelete)

	// Rollup Rule history
	router.HandleFunc(c.fullPath(rollupRuleHistoryPath), c.api.ReadRollupRuleHistory).Methods(http.MethodGet)

	log.Infof("Registered rules endpoints")
}

func (c r2Controller) fullPath(path string) string {
	return fmt.Sprintf("/%s/%s", c.rootPrefix, path)
}
