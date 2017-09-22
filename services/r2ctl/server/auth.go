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

package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/context"
)

const (
	// UserIDField is the name of the variable that simple auth will put in the mux context
	UserIDField = "userID"
)

// HTTPAuthService defines how to handle requests for various http authentication methods.
type HTTPAuthService interface {
	// User gets some identifying information about the user from the http request.
	User(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}

// SimpleAuth is a basic implementation
type SimpleAuth struct {
	userNameHeader string
}

// NewSimpleAuth...
func NewSimpleAuth(name string) SimpleAuth {
	return SimpleAuth{userNameHeader: name}
}

// User looks for a header defining a user name. If it can't find it, spits back not authorized.
func (a SimpleAuth) User(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	userID := r.Header.Get(a.userNameHeader)

	if userID == "" {
		http.Error(rw, fmt.Sprintf("No %s Provided", a.userNameHeader), 401)
	} else {
		// http.Error(rw, fmt.Sprintf("%s Provided", userID), 401)
		context.Set(r, "userID", userID)
		next(rw, r)
	}
}
