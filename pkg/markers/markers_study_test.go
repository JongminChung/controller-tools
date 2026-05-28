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

package markers

import "testing"

func TestRegistryLookupStudy(t *testing.T) {
	reg := &Registry{}

	if err := reg.Define("demo:flag", DescribesPackage, struct{}{}); err != nil {
		t.Fatalf("Define() error = %v", err)
	}
	if err := reg.Define("demo:flag:item", DescribesPackage, ""); err != nil {
		t.Fatalf("Define() error = %v", err)
	}

	short := reg.Lookup("+demo:flag", DescribesPackage)
	if short == nil || short.Name != "demo:flag" {
		t.Fatalf("Lookup(+demo:flag) = %#v, want demo:flag", short)
	}

	anon := reg.Lookup("+demo:flag:item=value", DescribesPackage)
	if anon == nil || anon.Name != "demo:flag:item" {
		t.Fatalf("Lookup(+demo:flag:item=value) = %#v, want demo:flag:item", anon)
	}
}

func TestArgumentTypeStringStudy(t *testing.T) {
	arg := Argument{
		Type: SliceType,
		ItemType: &Argument{
			Type:     MapType,
			ItemType: &Argument{Type: StringType},
		},
	}

	if got, want := arg.TypeString(), "[]map[string]string"; got != want {
		t.Fatalf("TypeString() = %q, want %q", got, want)
	}
}
