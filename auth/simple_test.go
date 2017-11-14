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
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	testConfig = SimpleAuthConfig{
		SimpleAuthNConfig: SimpleAuthNConfig{
			UserIDHeader: "testHeader",
		},
		SimpleAuthZConfig: SimpleAuthZConfig{
			WhitelistReads:          true,
			WhitelistWrites:         false,
			ReadWhitelistedUserIDs:  []string{"testUser"},
			WriteWhitelistedUserIDs: []string{},
		},
	}
)

func TestNewSimpleAuth(t *testing.T) {
	an := testConfig.NewSimpleAuth().(simpleAuth).authN
	az := testConfig.NewSimpleAuth().(simpleAuth).authZ
	require.Equal(t, an.userIDHeader, "testHeader")
	require.Equal(t, az.whitelistReads, true)
	require.Equal(t, az.whitelistWrites, false)
	require.Equal(t, az.readWhitelistedUserIDs, []string{"testUser"})
	require.Equal(t, az.writeWhitelistedUserIDs, []string{})
}

func TestSetUser(t *testing.T) {
	a := testConfig.NewSimpleAuth()
	ctx := context.Background()
	require.Nil(t, ctx.Value(UserIDField))
	ctx = a.SetUser(ctx, "foo")
	require.Equal(t, "foo", ctx.Value(UserIDField).(string))
}

func TestGetUser(t *testing.T) {
	a := testConfig.NewSimpleAuth()
	ctx := context.Background()

	id, err := a.GetUser(ctx)
	require.Empty(t, id)
	require.Error(t, err)

	ctx = a.SetUser(ctx, "foo")
	id, err = a.GetUser(ctx)
	require.Equal(t, "foo", id)
	require.NoError(t, err)
}

func TestSimpleAuthNAuthenticate(t *testing.T) {
	authN := simpleAuthN{
		userIDHeader: "foo",
	}

	require.Nil(t, authN.authenticate("bar"))
	require.EqualError(t, authN.authenticate(""), "Must provide header: [foo]")
}

func TestSimpleAuthZAuthorize(t *testing.T) {
	authZ := simpleAuthZ{
		whitelistReads:          true,
		whitelistWrites:         false,
		readWhitelistedUserIDs:  []string{"foo", "bar"},
		writeWhitelistedUserIDs: []string{"foo", "bar"},
	}

	require.Nil(t, authZ.authorize("GET", "foo"))
	require.Nil(t, authZ.authorize("POST", "foo"))
	require.Nil(t, authZ.authorize("OPTIONS", "foo"))
	require.EqualError(t, authZ.authorize("", "baz"), "Supplied userID: [baz] is not authorized")
}

func TestHealthCheck(t *testing.T) {
	a := testConfig.NewSimpleAuth()
	f := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v, err := a.GetUser(r.Context())
		require.NoError(t, err)
		require.Equal(t, "testHeader", v)
	})

	wrappedCall := a.NewAuthHandler(f)
	wrappedCall.ServeHTTP(httptest.NewRecorder(), &http.Request{})
}
