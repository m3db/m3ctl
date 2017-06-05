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
// THE SOFTWARE

package health

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/m3db/m3ctl/handler"
	"github.com/m3db/m3x/instrument"
)

const (
	ok   = "OK"
	fail = "FAIL"
)

var (
	healthURL = "/health"
)

type health struct {
	Host         string        `json:"host"`
	Timestamp    time.Time     `json:"timestamp"`
	ResponseTime time.Duration `json:"response_time"`
	Status       string        `json:"status"`
}

// Controller sets up the handlers for health endpoint
type Controller struct{}

// Returns health check status msg plus additional message
func healthCheck() string {
	return ok
}

func gatherHostInfo() string {
	host, err := os.Hostname()
	if err != nil {
		host = "unknown"
	}
	return host
}

// M3CtlHealthCheck is a health check for the m3ctl server
func m3ctlHealthCheck(w http.ResponseWriter, r *http.Request) {
	handler.InitJSONResponse(w)
	start := time.Now()
	host := gatherHostInfo()
	h := health{Host: host, Timestamp: start}
	status := healthCheck()
	h.ResponseTime = time.Since(start)

	h.Status = status
	if status != ok {
		w.WriteHeader(http.StatusInternalServerError)
	}
	body, _ := json.Marshal(h)

	w.Write(body)
}

// RegisterHandlers registers health handlers
func (c Controller) RegisterHandlers(mux *http.ServeMux, inst instrument.Options) error {
	log := inst.Logger()
	mux.HandleFunc(healthURL, m3ctlHealthCheck)
	log.Infof("Registered health endpoints")
	return nil
}
