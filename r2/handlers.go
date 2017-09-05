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

	"github.com/m3db/m3x/instrument"
)

// API holds all the state for a multi backend api  instance.
type handlers struct {
	iOpts   instrument.Options
	service Service
}

// Handlers exposes the http functions available to any implementation of an r2api.
type Handlers interface {
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

func (h *handlers) ReadNamespaces(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, "Placeholder for readNamespace")
}

func (h *handlers) ReadNamespace(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, "Placeholder for readNamespace")
}

func (h *handlers) CreateNamespace(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, "Placeholder for createNamespace")
}

func (h *handlers) DeleteNamespace(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, "Placeholder for deleteNamespace")
}

func (h *handlers) ReadMappingRule(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, "Placeholder for readMappingRule")
}

func (h *handlers) CreateMappingRule(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, "Placeholder for createMappingRule")
}

func (h *handlers) UpdateMappingRule(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, "Placeholder for updateMappingRule")
}

func (h *handlers) DeleteMappingRule(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, "Placeholder for deleteMappingRule")
}

func (h *handlers) ReadMappingRuleHistory(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, "Placeholder for readMappingRuleHistory")
}

func (h *handlers) ReadRollupRule(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, "Placeholder for readRollupRule")
}

func (h *handlers) CreateRollupRule(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, "Placeholder for createRollupRule")
}

func (h *handlers) UpdateRollupRule(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, "Placeholder for updateRollupRule")
}

func (h *handlers) DeleteRollupRule(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, "Placeholder for deleteRollupRule")
}

func (h *handlers) ReadRollupRuleHistory(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, "Placeholder for readRollupRuleHistory")
}
