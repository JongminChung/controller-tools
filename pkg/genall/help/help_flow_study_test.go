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

package help_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/fatih/color"
	genhelp "sigs.k8s.io/controller-tools/pkg/genall/help"
	prettyhelp "sigs.k8s.io/controller-tools/pkg/genall/help/pretty"
	"sigs.k8s.io/controller-tools/pkg/markers"
)

type printColumnStudyArgs struct {
	Name     string
	Priority int `marker:",optional"`
	Type     string
}

func renderHelpStudySpan(t *testing.T, span prettyhelp.Span) string {
	t.Helper()

	previous := color.NoColor
	color.NoColor = true
	defer func() {
		color.NoColor = previous
	}()

	var out bytes.Buffer
	if err := span.WriteTo(&out); err != nil {
		t.Fatalf("WriteTo() error = %v", err)
	}
	return out.String()
}

func newPrintColumnStudyInput(t *testing.T) (*markers.Definition, *markers.DefinitionHelp) {
	t.Helper()

	def, err := markers.MakeDefinition("kubebuilder:printcolumn", markers.DescribesType, printColumnStudyArgs{})
	if err != nil {
		t.Fatalf("MakeDefinition() error = %v", err)
	}

	defHelp := &markers.DefinitionHelp{
		Category: "kubebuilder",
		DetailedHelp: markers.DetailedHelp{
			Summary: "adds a printer column",
			Details: "Shows extra columns in kubectl get output.",
		},
		FieldHelp: map[string]markers.DetailedHelp{
			"Name": {
				Summary: "column header",
				Details: "This is the title shown above the rendered values.",
			},
			"Priority": {
				Summary: "visibility priority",
			},
			"Type": {
				Summary: "column data type",
			},
		},
	}

	return def, defHelp
}

func TestStudyInputToDocShape(t *testing.T) {
	def, defHelp := newPrintColumnStudyInput(t)

	doc := genhelp.ForDefinition(def, defHelp)

	if got, want := doc.Name, "kubebuilder:printcolumn"; got != want {
		t.Fatalf("doc.Name = %q, want %q", got, want)
	}
	if got, want := doc.Category, "kubebuilder"; got != want {
		t.Fatalf("doc.Category = %q, want %q", got, want)
	}
	if got, want := doc.Target, "type"; got != want {
		t.Fatalf("doc.Target = %q, want %q", got, want)
	}
	if got, want := len(doc.Fields), 3; got != want {
		t.Fatalf("len(doc.Fields) = %d, want %d", got, want)
	}
	if got, want := doc.Fields[0].Name, "name"; got != want {
		t.Fatalf("doc.Fields[0].Name = %q, want %q", got, want)
	}
	if got, want := doc.Fields[0].Summary, "column header"; got != want {
		t.Fatalf("doc.Fields[0].Summary = %q, want %q", got, want)
	}
	if got, want := doc.Fields[1].Name, "priority"; got != want {
		t.Fatalf("doc.Fields[1].Name = %q, want %q", got, want)
	}
	if !doc.Fields[1].Optional {
		t.Fatal("doc.Fields[1].Optional = false, want true")
	}
	if got, want := doc.Fields[2].Name, "type"; got != want {
		t.Fatalf("doc.Fields[2].Name = %q, want %q", got, want)
	}
}

func TestStudyInputToPrettySyntax(t *testing.T) {
	def, defHelp := newPrintColumnStudyInput(t)

	doc := genhelp.ForDefinition(def, defHelp)
	got := renderHelpStudySpan(t, prettyhelp.MarkerSyntaxHelp(doc))
	want := "+kubebuilder:printcolumn:name=<string>[,priority=<int>],type=<string>"
	if got != want {
		t.Fatalf("MarkerSyntaxHelp() = %q, want %q", got, want)
	}
}

func TestStudyRegistryToSummaryOutput(t *testing.T) {
	def, defHelp := newPrintColumnStudyInput(t)

	reg := &markers.Registry{}
	if err := reg.Register(def); err != nil {
		t.Fatalf("Register() error = %v", err)
	}
	reg.AddHelp(def, defHelp)

	categories := genhelp.ByCategory(reg, genhelp.SortByCategory)
	if got, want := len(categories), 1; got != want {
		t.Fatalf("len(categories) = %d, want %d", got, want)
	}

	summary := renderHelpStudySpan(t, prettyhelp.MarkersSummary(categories[0].Category, categories[0].Markers))
	checks := []string{
		"kubebuilder",
		"+kubebuilder:printcolumn:name=<string>[,priority=<int>],type=<string>",
		"type",
		"adds a printer column",
	}

	for _, check := range checks {
		if !strings.Contains(summary, check) {
			t.Fatalf("summary output = %q, missing %q", summary, check)
		}
	}
}

func TestStudyRegistryToDetailedOutput(t *testing.T) {
	def, defHelp := newPrintColumnStudyInput(t)

	reg := &markers.Registry{}
	if err := reg.Register(def); err != nil {
		t.Fatalf("Register() error = %v", err)
	}
	reg.AddHelp(def, defHelp)

	categories := genhelp.ByCategory(reg, genhelp.SortByCategory)
	details := renderHelpStudySpan(t, prettyhelp.MarkersDetails(true, categories[0].Category, categories[0].Markers))
	checks := []string{
		"+kubebuilder:printcolumn type",
		"\tadds a printer column",
		"\tShows extra columns in kubectl get output.",
		"\tname=<string>",
		"\t\tcolumn header",
		"\t[priority=<int>]",
		"\t\tvisibility priority",
		"\ttype=<string>",
		"\t\tcolumn data type",
	}

	for _, check := range checks {
		if !strings.Contains(details, check) {
			t.Fatalf("details output = %q, missing %q", details, check)
		}
	}
}
