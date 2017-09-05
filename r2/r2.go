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

	"github.com/m3db/m3ctl/services/r2ctl/server"
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

// Controller sets up the handlers for rules.
type r2Router struct {
	iOpts      instrument.Options
	handlers   Handlers
	rootPrefix string
}

// NewRouter creates a new r2 router using some given handlers.
func NewRouter(iOpts instrument.Options, rootPrefix string, service Service) server.Router {
	newHandlers := &handlers{
		iOpts:   iOpts,
		service: service,
	}

	return r2Router{iOpts: iOpts, handlers: newHandlers, rootPrefix: rootPrefix}
}

// RegisterHandlers registers rule handlers
func (r r2Router) RegisterHandlers(router *mux.Router) {
	log := r.iOpts.Logger()
	// Get all namespaces
	router.HandleFunc(r.fullPath(namespacePrefix), r.handlers.ReadNamespaces).Methods(http.MethodGet)

	fullNamespacePath := r.fullPath(namespacePath)
	// Ruleset actions
	router.HandleFunc(fullNamespacePath, r.handlers.ReadNamespace).Methods(http.MethodGet)
	router.HandleFunc(fullNamespacePath, r.handlers.CreateNamespace).Methods(http.MethodPost)
	router.HandleFunc(fullNamespacePath, r.handlers.DeleteNamespace).Methods(http.MethodDelete)

	// Mapping Rule actions
	router.HandleFunc(r.fullPath(mappingRuleRoot), r.handlers.CreateMappingRule).Methods(http.MethodPost)

	fullMappingRuleWithIDPath := r.fullPath(mappingRuleWithIDPath)
	router.HandleFunc(fullMappingRuleWithIDPath, r.handlers.ReadMappingRule).Methods(http.MethodGet)
	router.HandleFunc(fullMappingRuleWithIDPath, r.handlers.UpdateMappingRule).Methods(http.MethodPut, http.MethodPatch)
	router.HandleFunc(fullMappingRuleWithIDPath, r.handlers.DeleteMappingRule).Methods(http.MethodDelete)

	// Mapping Rule history
	router.HandleFunc(r.fullPath(mappingRuleHistoryPath), r.handlers.ReadMappingRuleHistory).Methods(http.MethodGet)

	// Rollup Rule actions
	router.HandleFunc(r.fullPath(rollupRuleRoot), r.handlers.CreateRollupRule).Methods(http.MethodPost)

	fullRollupRuleWithIDPath := r.fullPath(rollupRuleWithIDPath)
	router.HandleFunc(fullRollupRuleWithIDPath, r.handlers.ReadRollupRule).Methods(http.MethodGet)
	router.HandleFunc(fullRollupRuleWithIDPath, r.handlers.UpdateRollupRule).Methods(http.MethodPut, http.MethodPatch)
	router.HandleFunc(fullRollupRuleWithIDPath, r.handlers.DeleteRollupRule).Methods(http.MethodDelete)

	// Rollup Rule history
	router.HandleFunc(r.fullPath(rollupRuleHistoryPath), r.handlers.ReadRollupRuleHistory).Methods(http.MethodGet)

	log.Infof("Registered rules endpoints")
}

func (r r2Router) fullPath(path string) string {
	return fmt.Sprintf("/%s/%s", r.rootPrefix, path)
}
