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

package kv

import (
	"fmt"
	"net/http"

	"github.com/m3db/m3ctl/handler"
	"github.com/m3db/m3ctl/handler/r2"
	"github.com/m3db/m3x/instrument"
)

// API holds all the state for a multi backend api  instance.
type API struct {
	iOpts  instrument.Options
	config r2.APIConfig
}

// NewController creates a new rules controller using multiple backends
func NewController(iOpts instrument.Options, rootPrefix string, cfg r2.APIConfig) handler.Controller {
	newAPI := &API{
		iOpts:  iOpts,
		config: cfg,
	}
	return r2.NewController(iOpts, newAPI, rootPrefix)
}

// ReadNamespaces fetches all the namespaces.
func (a *API) ReadNamespaces(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, "Placeholder for readNamespaces")
}

// ReadNamespace fetches the current state of a namespace.
func (a *API) ReadNamespace(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, "Placeholder for readNamespace")
}

// CreateNamespace creates a namespace and new ruleset and persists it.
func (a *API) CreateNamespace(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, "Placeholder for createNamespace")
}

// DeleteNamespace deletes a namespace and its underlying ruleset.
func (a *API) DeleteNamespace(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, "Placeholder for deleteNamespace")
}

// ReadMappingRule fetches the current state of a mapping rule.
func (a *API) ReadMappingRule(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, "Placeholder for readMappingRule")
}

// CreateMappingRule creates or revives a mapping rule.
func (a *API) CreateMappingRule(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, "Placeholder for createMappingRule")
}

// UpdateMappingRule updates a mapping rule to a new state.
func (a *API) UpdateMappingRule(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, "Placeholder for updateMappingRule")
}

// DeleteMappingRule deletes a mapping rule.
func (a *API) DeleteMappingRule(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, "Placeholder for deleteMappingRule")
}

// ReadMappingRuleHistory fetches all of the snapshots for a given mapping rule.
func (a *API) ReadMappingRuleHistory(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, "Placeholder for readMappingRuleHistory")
}

// ReadRollupRule fetches the current state of a rollup rule.
func (a *API) ReadRollupRule(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, "Placeholder for readRollupRule")
}

// CreateRollupRule creates or revives a rollup rule.
func (a *API) CreateRollupRule(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, "Placeholder for createRollupRule")
}

// UpdateRollupRule updates a rollup rule to a new state.
func (a *API) UpdateRollupRule(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, "Placeholder for updateRollupRule")
}

// DeleteRollupRule deletes a rollup rule.
func (a *API) DeleteRollupRule(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, "Placeholder for deleteRollupRule")
}

// ReadRollupRuleHistory fetches all of the snapshots for a given rollup rule.
func (a *API) ReadRollupRuleHistory(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, "Placeholder for readRollupRuleHistory")
}
