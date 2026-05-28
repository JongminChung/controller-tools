/*
Copyright 2026 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package genall

import "testing"

func TestSplitOutputRuleOptionStudy(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantRule string
		wantGen  string
	}{
		{name: "default rule", input: "output:dir", wantRule: "dir"},
		{name: "generator specific rule", input: "output:crd:artifacts", wantRule: "artifacts", wantGen: "crd"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRule, gotGen := splitOutputRuleOption(tt.input)
			if gotRule != tt.wantRule || gotGen != tt.wantGen {
				t.Fatalf("splitOutputRuleOption(%q) = (%q, %q), want (%q, %q)", tt.input, gotRule, gotGen, tt.wantRule, tt.wantGen)
			}
		})
	}
}

func TestTransformRemoveCreationTimestampStudy(t *testing.T) {
	obj := map[string]any{
		"metadata": map[any]any{
			"creationTimestamp": "2026-05-29T00:00:00Z",
			"name":              "demo",
		},
	}

	if err := TransformRemoveCreationTimestamp(obj); err != nil {
		t.Fatalf("TransformRemoveCreationTimestamp() error = %v", err)
	}

	metadata := obj["metadata"].(map[any]any)
	if _, exists := metadata["creationTimestamp"]; exists {
		t.Fatal("creationTimestamp was not removed")
	}
	if got, want := metadata["name"], any("demo"); got != want {
		t.Fatalf("metadata[name] = %#v, want %#v", got, want)
	}
}
