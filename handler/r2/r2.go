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
	mappingRuleRoot         = fmt.Sprintf("%s/%s", namespacePrefix, mappingRulePrefix)
	mappingRuleWithIDPath   = fmt.Sprintf("%s/{ruleId}", mappingRuleRoot)
	mappingRuleWithNamePath = fmt.Sprintf("%s/{ruleName}", mappingRuleRoot)
	mappingRuleHistoryPath  = fmt.Sprintf("%s/history", mappingRuleWithIDPath)

	rollupRuleRoot         = fmt.Sprintf("%s/%s", namespacePrefix, rollupRulePrefix)
	rollupRuleWithIDPath   = fmt.Sprintf("%s/{ruleId}", rollupRuleRoot)
	rollupRuleWithNamePath = fmt.Sprintf("%s/{ruleName}", rollupRuleRoot)
	rollupRuleHistoryPath  = fmt.Sprintf("%s/history", rollupRuleWithIDPath)
)

// Controller sets up the handlers for rules
type r2Controller struct {
	iOpts      instrument.Options
	api        r2Api
	rootPrefix string
}

// NewMultiBackendController creates a new rules controller
func NewMultiBackendController(iOpts instrument.Options, rootPrefix string) handler.Controller {
	return r2Controller{iOpts: iOpts, api: newMultiBackendR2Api(iOpts), rootPrefix: rootPrefix}
}

type r2Api interface {
	readNamespaces(http.ResponseWriter, *http.Request)
	readNamespace(http.ResponseWriter, *http.Request)
	createNamespace(http.ResponseWriter, *http.Request)
	deleteNamespace(http.ResponseWriter, *http.Request)

	readMappingRule(http.ResponseWriter, *http.Request)
	createMappingRule(http.ResponseWriter, *http.Request)
	updateMappingRule(http.ResponseWriter, *http.Request)
	deleteMappingRule(http.ResponseWriter, *http.Request)
	readMappingRuleHistory(http.ResponseWriter, *http.Request)

	readRollupRule(http.ResponseWriter, *http.Request)
	createRollupRule(http.ResponseWriter, *http.Request)
	updateRollupRule(http.ResponseWriter, *http.Request)
	deleteRollupRule(http.ResponseWriter, *http.Request)
	readRollupRuleHistory(http.ResponseWriter, *http.Request)
}

// RegisterHandlers registers rule handlers
func (c r2Controller) RegisterHandlers(router *mux.Router) {
	log := c.iOpts.Logger()
	// Get all namespaces
	router.HandleFunc(c.fullPath(namespacePrefix), c.api.readNamespaces).Methods(http.MethodGet)

	fullNamespacePath := c.fullPath(namespacePath)
	// Ruleset actions
	router.HandleFunc(fullNamespacePath, c.api.readNamespace).Methods(http.MethodGet)
	router.HandleFunc(fullNamespacePath, c.api.createNamespace).Methods(http.MethodPost)
	router.HandleFunc(fullNamespacePath, c.api.deleteNamespace).Methods(http.MethodDelete)

	// Mapping Rule actions
	router.HandleFunc(c.fullPath(mappingRuleWithNamePath), c.api.createMappingRule).Methods(http.MethodPost)

	fullMappingRuleWithIDPath := c.fullPath(mappingRuleWithIDPath)
	router.HandleFunc(fullMappingRuleWithIDPath, c.api.readMappingRule).Methods(http.MethodGet)
	router.HandleFunc(fullMappingRuleWithIDPath, c.api.updateMappingRule).Methods(http.MethodPut, http.MethodPatch)
	router.HandleFunc(fullMappingRuleWithIDPath, c.api.deleteMappingRule).Methods(http.MethodDelete)

	// Mapping Rule history
	router.HandleFunc(c.fullPath(mappingRuleHistoryPath), c.api.readMappingRuleHistory).Methods(http.MethodGet)

	// Rollup Rule actions
	router.HandleFunc(c.fullPath(rollupRuleWithNamePath), c.api.createRollupRule).Methods(http.MethodPost)

	fullRollupRuleWithIDPath := c.fullPath(rollupRuleWithIDPath)
	router.HandleFunc(fullRollupRuleWithIDPath, c.api.readRollupRule).Methods(http.MethodGet)
	router.HandleFunc(fullRollupRuleWithIDPath, c.api.updateRollupRule).Methods(http.MethodPut, http.MethodPatch)
	router.HandleFunc(fullRollupRuleWithIDPath, c.api.deleteRollupRule).Methods(http.MethodDelete)

	// Rollup Rule history
	router.HandleFunc(c.fullPath(rollupRuleHistoryPath), c.api.readRollupRuleHistory).Methods(http.MethodGet)

	log.Infof("Registered rules endpoints")
}

func (c r2Controller) fullPath(path string) string {
	return fmt.Sprintf("/%s/%s", c.rootPrefix, path)
}
