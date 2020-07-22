// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package header

import (
	"bufio"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-safeweb/safehttp"
	"github.com/google/safehtml"
)

type dispatcher struct{}

func (dispatcher) Write(rw http.ResponseWriter, resp safehttp.Response) error {
	switch x := resp.(type) {
	case safehtml.HTML:
		_, err := rw.Write([]byte(x.String()))
		return err
	default:
		panic("not a safe response type")
	}
}

func (dispatcher) ExecuteTemplate(rw http.ResponseWriter, t safehttp.Template, data interface{}) error {
	return nil
}

type responseWriter struct {
	header http.Header
	writer io.Writer
	status int
}

func newResponseWriter(w io.Writer) *responseWriter {
	return &responseWriter{header: http.Header{}, writer: w, status: http.StatusOK}
}

func (r *responseWriter) Header() http.Header {
	return r.header
}

func (r *responseWriter) WriteHeader(statusCode int) {
	r.status = statusCode
}

func (r *responseWriter) Write(data []byte) (int, error) {
	return r.writer.Write(data)
}

func TestAccessIncomingHeaders(t *testing.T) {
	m := safehttp.NewMachinery(func(rw safehttp.ResponseWriter, r *safehttp.IncomingRequest) safehttp.Result {
		if got, want := r.Header.Get("A"), "B"; got != want {
			t.Errorf(`r.Header.Get("A") got: %v want: %v`, got, want)
		}
		return rw.Write(safehtml.HTMLEscaped("hello"))
	}, &dispatcher{})

	request := "GET / HTTP/1.1\r\n" +
		"Host: localhost:8080\r\n" +
		"A: B\r\n\r\n"
	req, err := http.ReadRequest(bufio.NewReader(strings.NewReader(request)))
	if err != nil {
		t.Fatalf("http.ReadRequest() got err: %v", err)
	}

	b := &strings.Builder{}
	rw := newResponseWriter(b)

	m.HandleRequest(rw, req)
}

func TestChangingResponseHeaders(t *testing.T) {
	m := safehttp.NewMachinery(func(rw safehttp.ResponseWriter, _ *safehttp.IncomingRequest) safehttp.Result {
		rw.Header().Set("pIZZA", "Pasta")
		return rw.Write(safehtml.HTMLEscaped("hello"))
	}, &dispatcher{})

	req := httptest.NewRequest("GET", "/", nil)

	b := &strings.Builder{}
	rw := newResponseWriter(b)

	m.HandleRequest(rw, req)

	want := []string{"Pasta"}
	if diff := cmp.Diff(want, rw.Header()["Pizza"]); diff != "" {
		t.Errorf(`resp.Header["Pizza"] mismatch (-want +got):\n%s`, diff)
	}
}