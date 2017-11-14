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

// SimpleAuthConfig holds this configuration necessary for a simple auth implementation.
type SimpleAuthConfig struct {
	// AuthN configuration
	SimpleAuthNConfig SimpleAuthNConfig `yaml:"authN"`
	// AuthZ configuration
	SimpleAuthZConfig SimpleAuthZConfig `yaml:"authZ"`
}

// SimpleAuthNConfig holds this configuration necessary for a simple authentication implementation.
type SimpleAuthNConfig struct {
	// This is an http header that identifies the user performing the operation
	UserIDHeader string `yaml:"userIDHeader" validate:"nonzero"`
}

// SimpleAuthZConfig holds this configuration necessary for a simple authorization implementation.
type SimpleAuthZConfig struct {
	// This indicates whether reads should use a read whitelist
	WhitelistReads bool `yaml:"whitelistReads,omitempty"`
	// This indicates whether writes should use a read whitelist
	WhitelistWrites bool `yaml:"whitelistWrites,omitempty"`
	// This is a list of users that are allowed to perform read operations
	ReadWhitelistedUserIDs []string `yaml:"readWhitelistedUserIDs,omitempty"`
	// This is a list of users that are allowed to perform write operations
	WriteWhitelistedUserIDs []string `yaml:"readWhitelistedUserIDs,omitempty"`
}

// NewSimpleAuth creates a new simple auth instance given using the provided config.
func (ac SimpleAuthConfig) NewSimpleAuth() HTTPAuthService {
	return simpleAuth{
		authN: simpleAuthN{
			userIDHeader: ac.SimpleAuthNConfig.UserIDHeader,
		},
		authZ: simpleAuthZ{
			whitelistReads:          ac.SimpleAuthZConfig.WhitelistReads,
			whitelistWrites:         ac.SimpleAuthZConfig.WhitelistWrites,
			readWhitelistedUserIDs:  ac.SimpleAuthZConfig.ReadWhitelistedUserIDs,
			writeWhitelistedUserIDs: ac.SimpleAuthZConfig.WriteWhitelistedUserIDs,
		},
	}
}

type simpleAuth struct {
	authN simpleAuthN
	authZ simpleAuthZ
}

type simpleAuthN struct {
	userIDHeader string
}

type simpleAuthZ struct {
	whitelistReads          bool
	whitelistWrites         bool
	readWhitelistedUserIDs  []string
	writeWhitelistedUserIDs []string
}

// Authenticate looks for a header defining a user name. If it finds it, runs the actual http handler passed as a parameter.
// Otherwise, it returns an Unauthorized http response.
func (a simpleAuth) NewAuthHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get(a.authN.userIDHeader)
		err := a.authN.authenticate(userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		err = a.authZ.authorize(r.Method, userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
		}

		ctx := a.SetUser(r.Context(), userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// SetUser sets the user making the changes to the api.
func (a simpleAuth) SetUser(parent context.Context, userID string) context.Context {
	return context.WithValue(parent, UserIDField, userID)
}

// GetUser fetches the ID of an api caller from the global context.
func (a simpleAuth) GetUser(ctx context.Context) (string, error) {
	id := ctx.Value(UserIDField)
	if id == nil {
		return "", fmt.Errorf("couldn't identify user")
	}
	return id.(string), nil
}

func (a simpleAuthN) authenticate(userID string) error {
	if userID == "" {
		return fmt.Errorf("Must provide header: [%s]", a.userIDHeader)
	}
	return nil
}

func (a simpleAuthZ) authorize(httpMethod, userID string) error {
	switch httpMethod {
	case "GET", "":
		return a.authorizeUser(a.whitelistReads, a.readWhitelistedUserIDs, userID)
	case "POST", "PUT", "PATCH", "DELETE":
		return a.authorizeUser(a.whitelistWrites, a.writeWhitelistedUserIDs, userID)
	default:
		return nil
	}
}

func (a simpleAuthZ) authorizeUser(useWhitelist bool, whitelistedUsers []string, userID string) error {
	if !useWhitelist {
		return nil
	}

	isAuthorized := false
	for _, u := range whitelistedUsers {
		if u == userID {
			isAuthorized = true
			break
		}
	}
	if !isAuthorized {
		return fmt.Errorf("Supplied userID: [%s] is not authorized", userID)
	}

	return nil
}
