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

import "net/http"
import "fmt"

var (
	globalLevel    = "/rules/namespaces"
	namespaceLevel = "/rules/namespace/"
)

// Controller sets up the handlers for rules
type Controller struct{}

func handleGlobals(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintf(w, "%v: In global handler. ", r.Method)
	switch r.Method {
	case http.MethodGet:
		getNamespaces(w, r)
	case http.MethodPut:
		createNamespace(w, r)
	case http.MethodDelete:
		deleteNamespace(w, r)
	default:
		fmt.Fprintf(w, "Not defined")
	}
}

func handleNamespace(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintf(w, "%v: In namespace handler", r.Method)
	switch r.Method {
	case http.MethodGet:
		getRulesInNamespace(w, r)
	case http.MethodPut:
		setRulesInNamespace(w, r)
	default:
		fmt.Fprintf(w, "Not defined")
	}
}

// RegisterHandlers registers namespace level handlers
func (r Controller) RegisterHandlers(mux *http.ServeMux) error {
	mux.HandleFunc(globalLevel, handleGlobals)
	mux.HandleFunc(namespaceLevel, handleNamespace)
	return nil
}
