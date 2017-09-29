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
package auth

import (
	"context"
	"fmt"
	"net/http"
)

// SimpleAuthConfig holds this configuration necessary for a simple auth implementation
type SimpleAuthConfig struct {
	// This is an http header that identifies the user performing the operation
	UserNameHeader string `yaml:"userNameHeader" validate:"nonzero"`
}

// SimpleAuth is an naive authentication method that just checks for a HTTPHeader identifying the user
// performing the action.
type SimpleAuth struct {
	userNameHeader string
}

// NewSimpleAuth defines an authentication method that just checks for the existence of a valid user name header
func NewSimpleAuth(name string) SimpleAuth {
	return SimpleAuth{userNameHeader: name}
}

// WithCheck looks for a header defining a user name. If it finds it, runs the actual http handler passed as a parameter.
// Otherwise, it returns an Unathorized http response.
func (a SimpleAuth) WithCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get(a.userNameHeader)
		if userID == "" {
			http.Error(w, fmt.Sprintf("Must provide header: [%s]", a.userNameHeader), http.StatusUnauthorized)
			return
		}
		ctx := a.SetUser(r.Context(), userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a SimpleAuth) SetUser(parent context.Context, userID string) context.Context {
	return context.WithValue(parent, UserIDField, userID)
}

// GetUser fetches the ID of an api caller from the global context.
func (a SimpleAuth) GetUser(ctx context.Context) (string, error) {
	id := ctx.Context().Value(UserIDField)
	if id == "" {
		return "", fmt.Errorf("Couldn't identify user. %s is not set.", UserIDField)
	}

	s.iOpts.Logger().Infof("%v", id)
}
