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
	"net/http"
	"testing"

	"github.com/m3db/m3ctl/auth"

	"github.com/stretchr/testify/require"
)

func TestDefaultAuthorizationTypeForHTTPMethodGet(t *testing.T) {
	require.EqualValues(
		t,
		auth.AuthorizationTypeReadOnly,
		defaultAuthorizationTypeForHTTPMethod(http.MethodGet),
	)
}
func TestDefaultAuthorizationTypeForHTTPMethodPost(t *testing.T) {
	require.EqualValues(
		t,
		auth.AuthorizationTypeReadWrite,
		defaultAuthorizationTypeForHTTPMethod(http.MethodPost),
	)
}

func TestDefaultAuthorizationTypeForHTTPMethodPut(t *testing.T) {
	require.EqualValues(
		t,
		auth.AuthorizationTypeReadWrite,
		defaultAuthorizationTypeForHTTPMethod(http.MethodPut),
	)
}

func TestDefaultAuthorizationTypeForHTTPMethodPatch(t *testing.T) {
	require.EqualValues(
		t,
		auth.AuthorizationTypeReadWrite,
		defaultAuthorizationTypeForHTTPMethod(http.MethodPatch),
	)
}

func TestDefaultAuthorizationTypeForHTTPMethodDelete(t *testing.T) {
	require.EqualValues(
		t,
		auth.AuthorizationTypeReadWrite,
		defaultAuthorizationTypeForHTTPMethod(http.MethodDelete),
	)
}

func TestDefaultAuthorizationTypeForHTTPMethodUnrecognizedMethod(t *testing.T) {
	require.EqualValues(
		t,
		auth.AuthorizationTypeUnknown,
		defaultAuthorizationTypeForHTTPMethod(http.MethodOptions),
	)
}
