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

	"github.com/m3db/m3x/instrument"

	"github.com/m3db/m3ctl/handler"
)

const (
	getAllURL          = "/rules/namespaces"
	createOrDeleteURL  = "/rules/namespace"
	modifyNamespaceURL = "/rules/namespace/"
)

// Controller sets up the handlers for rules
type controller struct {
	instruments instrument.Options
}

// NewController creates a new rules controller
func NewController(instrument instrument.Options) handler.Controller {
	return &controller{instruments: instrument}
}

func (c *controller) handleNotSupported(w http.ResponseWriter, r *http.Request, supported []string) {
	err := handler.MethodNotSupportedResponse(w, r, []string{http.MethodGet})
	if err != nil {
		c.instruments.Logger().Errorf("%v", err)
	}
}

// RegisterHandlers registers rule handlers
func (c *controller) RegisterHandlers(mux *http.ServeMux) error {
	log := c.instruments.Logger()
	mux.HandleFunc(getAllURL, c.handleGetAll)
	mux.HandleFunc(createOrDeleteURL, c.handleCreateDelete)
	mux.HandleFunc(modifyNamespaceURL, c.handleModify)
	log.Infof("Registered rules endpoints")
	return nil
}

func (c *controller) handleGetAll(w http.ResponseWriter, r *http.Request) {
	handler.InitJSONResponse(w)
	switch r.Method {
	case http.MethodGet:
		getNamespaces(w, r)
	default:
		c.handleNotSupported(w, r, []string{http.MethodGet})
	}
}

func (c *controller) handleCreateDelete(w http.ResponseWriter, r *http.Request) {
	handler.InitJSONResponse(w)
	switch r.Method {
	case http.MethodPost:
		createNamespace(w, r)
	case http.MethodDelete:
		deleteNamespace(w, r)
	default:
		c.handleNotSupported(w, r, []string{http.MethodPost, http.MethodDelete})
	}
}

func (c *controller) handleModify(w http.ResponseWriter, r *http.Request) {
	handler.InitJSONResponse(w)
	switch r.Method {
	case http.MethodGet:
		getRulesInNamespace(w, r)
	case http.MethodPut:
		setRulesInNamespace(w, r)
	default:
		c.handleNotSupported(w, r, []string{http.MethodGet, http.MethodPut})
	}
}
