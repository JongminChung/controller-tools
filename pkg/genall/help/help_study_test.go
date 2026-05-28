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

package help

import (
	"testing"

	"sigs.k8s.io/controller-tools/pkg/markers"
)

func TestForArgumentStudy(t *testing.T) {
	arg := markers.Argument{
		Type: markers.SliceType,
		ItemType: &markers.Argument{
			Type: markers.StringType,
		},
		Optional: true,
	}

	docArg := ForArgument(arg)
	if got, want := docArg.Type, "slice"; got != want {
		t.Fatalf("ForArgument().Type = %q, want %q", got, want)
	}
	if got, want := docArg.TypeString(), "[]string"; got != want {
		t.Fatalf("ForArgument().TypeString() = %q, want %q", got, want)
	}
	if !docArg.Optional {
		t.Fatal("ForArgument().Optional = false, want true")
	}
}

func TestForDefinitionAndByCategoryStudy(t *testing.T) {
	type markerArgs struct {
		Name string
		Age  int
	}

	def, err := markers.MakeDefinition("demo:marker", markers.DescribesType, markerArgs{})
	if err != nil {
		t.Fatalf("MakeDefinition() error = %v", err)
	}

	help := &markers.DefinitionHelp{
		Category:     "demo",
		DetailedHelp: markers.DetailedHelp{Summary: "study marker"},
		FieldHelp: map[string]markers.DetailedHelp{
			"name": {Summary: "human name"},
			"age":  {Summary: "human age"},
		},
	}

	doc := ForDefinition(def, help)
	if got, want := len(doc.Fields), 2; got != want {
		t.Fatalf("len(ForDefinition().Fields) = %d, want %d", got, want)
	}
	if got, want := doc.Fields[0].Name, "age"; got != want {
		t.Fatalf("first field = %q, want %q", got, want)
	}

	reg := &markers.Registry{}
	if err := reg.Register(def); err != nil {
		t.Fatalf("Register() error = %v", err)
	}
	reg.AddHelp(def, help)

	categories := ByCategory(reg, SortByCategory)
	if got, want := len(categories), 1; got != want {
		t.Fatalf("len(ByCategory()) = %d, want %d", got, want)
	}
	if got, want := categories[0].Category, "demo"; got != want {
		t.Fatalf("category = %q, want %q", got, want)
	}
	if got, want := categories[0].Markers[0].Name, "demo:marker"; got != want {
		t.Fatalf("marker name = %q, want %q", got, want)
	}
}
