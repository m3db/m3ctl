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

package rules

import (
	"net/http"

	"github.com/m3db/m3ctl/handler"
	"github.com/m3db/m3x/instrument"
)

var (
	getAllURL          = "/rules/namespaces"
	createOrDeleteURL  = "/rules/namespace"
	modifyNamespaceURL = "/rules/namespace/"
)

// Controller sets up the handlers for rules
type Controller struct{}

// RegisterHandlers registers rule handlers
func (c Controller) RegisterHandlers(mux *http.ServeMux, inst instrument.Options) error {
	log := inst.Logger()
	mux.HandleFunc(getAllURL, handleGetAll)
	mux.HandleFunc(createOrDeleteURL, handleCreateDelete)
	mux.HandleFunc(modifyNamespaceURL, handleModify)
	log.Infof("Registered rules endpoints")
	return nil
}

func handleGetAll(w http.ResponseWriter, r *http.Request) {
	handler.InitJSONResponse(w)
	switch r.Method {
	case http.MethodGet:
		getNamespaces(w, r)
	default:
		handler.MethodNotSupportedResponse(w, r, []string{http.MethodGet})
	}
}

func handleCreateDelete(w http.ResponseWriter, r *http.Request) {
	handler.InitJSONResponse(w)
	switch r.Method {
	case http.MethodPost:
		createNamespace(w, r)
	case http.MethodDelete:
		deleteNamespace(w, r)
	default:
		handler.MethodNotSupportedResponse(w, r, []string{http.MethodPost, http.MethodDelete})
	}
}

func handleModify(w http.ResponseWriter, r *http.Request) {
	handler.InitJSONResponse(w)
	switch r.Method {
	case http.MethodGet:
		getRulesInNamespace(w, r)
	case http.MethodPut:
		setRulesInNamespace(w, r)
	default:
		handler.MethodNotSupportedResponse(w, r, []string{http.MethodGet, http.MethodPut})
	}
}
