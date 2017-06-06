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
// THE SOFTWARE

package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMethodNotSupportedResponse(t *testing.T) {
	rr := httptest.NewRecorder()
	goodMethods := []string{http.MethodDelete, http.MethodConnect}

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req := http.Request{Method: "GET", RequestURI: "/badMethod"}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		MethodNotSupportedResponse(w, r, goodMethods)
	})

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, &req)

	rawResult := make([]byte, rr.Body.Len())
	_, err := rr.Body.Read(rawResult)
	require.NoError(t, err, "Encountered error parsing reponse")

	var actualResult NotSupportedResponse
	json.Unmarshal(rawResult, &actualResult)

	expectedResult := NotSupportedResponse{
		Error:            "Method not supported",
		URI:              "/badMethod",
		MethodName:       "GET",
		SupportedMethods: goodMethods,
	}

	require.Equal(t, rr.Code, http.StatusMethodNotAllowed)
	require.Equal(t, expectedResult, actualResult)
}
