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

type multiBackendR2Api struct {
	iOpts instrument.Options
	// backends -- Backend List
}

func newMultiBackendR2Api(iOpts instrument.Options) multiBackendR2Api {
	return multiBackendR2Api{iOpts: iOpts}
}

func (c multiBackendR2Api) readNamespaces(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, "Placeholder for readNamespaces")
}

func (c multiBackendR2Api) readNamespace(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, "Placeholder for readNamespace")
}

func (c multiBackendR2Api) createNamespace(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, "Placeholder for createNamespace")
}

func (c multiBackendR2Api) deleteNamespace(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, "Placeholder for deleteNamespace")
}

func (c multiBackendR2Api) readMappingRule(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, "Placeholder for readMappingRule")
}

func (c multiBackendR2Api) createMappingRule(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, "Placeholder for createMappingRule")
}

func (c multiBackendR2Api) updateMappingRule(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, "Placeholder for updateMappingRule")
}

func (c multiBackendR2Api) deleteMappingRule(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, "Placeholder for deleteMappingRule")
}

func (c multiBackendR2Api) readMappingRuleHistory(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, "Placeholder for readMappingRuleHistory")
}

func (c multiBackendR2Api) readRollupRule(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, "Placeholder for readRollupRule")
}

func (c multiBackendR2Api) createRollupRule(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, "Placeholder for createRollupRule")
}

func (c multiBackendR2Api) updateRollupRule(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, "Placeholder for updateRollupRule")
}

func (c multiBackendR2Api) deleteRollupRule(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, "Placeholder for deleteRollupRule")
}

func (c multiBackendR2Api) readRollupRuleHistory(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, "Placeholder for readRollupRuleHistory")
}
