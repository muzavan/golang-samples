// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/cloudrunci"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestDiagramService(t *testing.T) {
	tc := testutil.EndToEndTest(t)
	service := cloudrunci.NewService("diagram", tc.ProjectID)
	if err := service.Deploy(); err != nil {
		t.Fatalf("could not deploy %s: %v", service.Name, err)
	}
	defer service.Clean()

	requestPath := "/diagram.png?" + url.QueryEscape("digraph G { A -> {B, C, D} -> {F} }")
	req, err := service.NewRequest("GET", requestPath)
	if err != nil {
		t.Errorf("service.NewRequest: %q", err)
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Errorf("Get: %v", err)
		return
	}
	defer resp.Body.Close()

	if got := resp.StatusCode; got != 200 {
		t.Errorf("response status: got %d, want %d", got, 200)
	}

	want := "image/png"
	if got := resp.Header.Get("Content-Type"); got != want {
		t.Errorf("response Content-Type: got %q, want %s", got, want)
	}

	want = "public, max-age=86400"
	if got := resp.Header.Get("Cache-Control"); got != want {
		t.Errorf("response Cache-Control: got %q, want %q", got, want)
	}
}