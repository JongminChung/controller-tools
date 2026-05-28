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

package loader

import (
	"errors"
	"go/ast"
	"testing"
)

func TestNonVendorPathStudy(t *testing.T) {
	got := NonVendorPath("example.com/project/vendor/sigs.k8s.io/controller-tools/pkg/loader")
	if want := "sigs.k8s.io/controller-tools/pkg/loader"; got != want {
		t.Fatalf("NonVendorPath() = %q, want %q", got, want)
	}
}

func TestParseAstTagStudy(t *testing.T) {
	tag := &ast.BasicLit{Value: "`json:\"name,omitempty\" yaml:\"name\"`"}
	parsed := ParseAstTag(tag)

	if got, want := parsed.Get("json"), "name,omitempty"; got != want {
		t.Fatalf("json tag = %q, want %q", got, want)
	}
	if got, want := parsed.Get("yaml"), "name"; got != want {
		t.Fatalf("yaml tag = %q, want %q", got, want)
	}
}

func TestMaybeErrListStudy(t *testing.T) {
	if err := MaybeErrList(nil); err != nil {
		t.Fatalf("MaybeErrList(nil) = %v, want nil", err)
	}

	err := MaybeErrList([]error{errors.New("first"), errors.New("second")})
	var list ErrList
	if !errors.As(err, &list) {
		t.Fatalf("MaybeErrList() = %T, want ErrList", err)
	}
	if got, want := len(list), 2; got != want {
		t.Fatalf("len(ErrList) = %d, want %d", got, want)
	}
}
