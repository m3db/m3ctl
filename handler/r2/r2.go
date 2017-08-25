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
	r2v1Prefix      = "/r2/v1"
	namespacePrefix = "namespace/{namespaceId}"
)

var (
	namespacesURL          = fmt.Sprintf("%s/namespaces", r2v1Prefix)
	namespaceURL           = fmt.Sprintf("%s/%s", r2v1Prefix, namespacePrefix)
	mappingRuleWithIDURL   = "/r2/v1/namespace/{namespaceId}/mapping/{ruleId}"
	mappingRuleHistoryURL  = "/r2/v1/namespace/{namespaceId}/mapping/{ruleId}/history"
	mappingRuleWithNameURL = "/r2/v1/namespace/{namespaceId}/mapping/{ruleName}"
	rollupRuleWithIDURL    = "/r2/v1/namespace/{namespaceId}/rollup/{ruleId}"
	rollupRuleHistoryURL   = "/r2/v1/namespace/{namespaceId}/rollup/{ruleId}/history"
	rollupRuleWithNameURL  = "/r2/v1/namespace/{namespaceId}/rollup/{ruleName}"
)

// Controller sets up the handlers for rules
type r2Controller struct {
	iOpts instrument.Options
}

// type r2api interface {
// 	readNamespaces(http.ResponseWriter, *http.Request)
// 	readNamespace(http.ResponseWriter, *http.Request)
// 	createNamespace(http.ResponseWriter, *http.Request)
// 	deleteNamespace(http.ResponseWriter, *http.Request)

// 	readMappingRule(http.ResponseWriter, *http.Request)
// 	createMappingRule(http.ResponseWriter, *http.Request)
// 	updateMappingRule(http.ResponseWriter, *http.Request)
// 	deleteMappingRule(http.ResponseWriter, *http.Request)
// 	readMappingRuleHistory(http.ResponseWriter, *http.Request)

// 	readRollupRule(http.ResponseWriter, *http.Request)
// 	createRollupRule(http.ResponseWriter, *http.Request)
// 	updateRollupRule(http.ResponseWriter, *http.Request)
// 	deleteRollupRule(http.ResponseWriter, *http.Request)
// 	readRollupRuleHistory(http.ResponseWriter, *http.Request)
// }

// NewController creates a new rules controller
func NewController(iOpts instrument.Options) handler.Controller {
	return &r2Controller{iOpts: iOpts}
}

func (c r2Controller) readNamespaces(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, "Placeholder for readNamespaces")
}

func (c r2Controller) readNamespace(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, "Placeholder for readNamespace")
}

func (c r2Controller) createNamespace(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, "Placeholder for createNamespace")
}

func (c r2Controller) deleteNamespace(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, "Placeholder for deleteNamespace")
}

func (c r2Controller) readMappingRule(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, "Placeholder for readMappingRule")
}

func (c r2Controller) createMappingRule(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, "Placeholder for createMappingRule")
}

func (c r2Controller) updateMappingRule(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, "Placeholder for updateMappingRule")
}

func (c r2Controller) deleteMappingRule(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, "Placeholder for deleteMappingRule")
}

func (c r2Controller) readMappingRuleHistory(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, "Placeholder for readMappingRuleHistory")
}

func (c r2Controller) readRollupRule(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, "Placeholder for readRollupRule")
}

func (c r2Controller) createRollupRule(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, "Placeholder for createRollupRule")
}

func (c r2Controller) updateRollupRule(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, "Placeholder for updateRollupRule")
}

func (c r2Controller) deleteRollupRule(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, "Placeholder for deleteRollupRule")
}

func (c r2Controller) readRollupRuleHistory(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, "Placeholder for readRollupRuleHistory")
}

// RegisterHandlers registers rule handlers
func (c r2Controller) RegisterHandlers(router *mux.Router) {
	log := c.iOpts.Logger()
	// Get all namespaces
	router.HandleFunc(namespacesURL, c.readNamespaces).Methods(http.MethodGet)

	// Get current state of a ruleset
	router.HandleFunc(namespaceURL, c.readNamespace).Methods(http.MethodGet)

	// Create/Delete namespace
	router.HandleFunc(namespaceURL, c.createNamespace).Methods(http.MethodPost)
	router.HandleFunc(namespaceURL, c.deleteNamespace).Methods(http.MethodDelete)

	// Mapping Rule actions
	router.HandleFunc(mappingRuleWithNameURL, c.createMappingRule).Methods(http.MethodPost)
	router.HandleFunc(mappingRuleWithIDURL, c.readMappingRule).Methods(http.MethodGet)
	router.HandleFunc(mappingRuleWithIDURL, c.updateMappingRule).Methods(http.MethodPut, http.MethodPatch)
	router.HandleFunc(mappingRuleWithIDURL, c.deleteMappingRule).Methods(http.MethodDelete)

	// Mapping Rule history
	router.HandleFunc(mappingRuleHistoryURL, c.readMappingRuleHistory).Methods(http.MethodGet)

	// Rollup Rule actions
	router.HandleFunc(rollupRuleWithNameURL, c.createRollupRule).Methods(http.MethodPost)
	router.HandleFunc(rollupRuleWithIDURL, c.readRollupRule).Methods(http.MethodGet)
	router.HandleFunc(rollupRuleWithIDURL, c.updateRollupRule).Methods(http.MethodPut, http.MethodPatch)
	router.HandleFunc(rollupRuleWithIDURL, c.deleteRollupRule).Methods(http.MethodDelete)

	// Rollup Rule history
	router.HandleFunc(rollupRuleHistoryURL, c.readRollupRuleHistory).Methods(http.MethodGet)

	log.Infof("Registered rules endpoints")
}
