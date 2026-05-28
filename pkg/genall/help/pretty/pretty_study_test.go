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

package pretty

import (
	"bytes"
	"testing"

	"github.com/fatih/color"
	"sigs.k8s.io/controller-tools/pkg/genall/help"
)

func renderStudySpan(t *testing.T, span Span) string {
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

func TestFieldSyntaxHelpStudy(t *testing.T) {
	field := help.FieldHelp{
		Name: "paths",
		Argument: help.Argument{
			Type:     "slice",
			ItemType: &help.Argument{Type: "string"},
			Optional: true,
		},
	}

	if got, want := renderStudySpan(t, FieldSyntaxHelp(field)), "[paths=<[]string>]"; got != want {
		t.Fatalf("FieldSyntaxHelp() = %q, want %q", got, want)
	}
}

func TestMarkerSyntaxHelpStudy(t *testing.T) {
	marker := help.MarkerDoc{
		Name: "kubebuilder:printcolumn",
		Fields: []help.FieldHelp{
			{
				Name: "name",
				Argument: help.Argument{
					Type: "string",
				},
			},
			{
				Name: "type",
				Argument: help.Argument{
					Type: "string",
				},
			},
			{
				Name: "priority",
				Argument: help.Argument{
					Type:     "int",
					Optional: true,
				},
			},
		},
	}

	got := renderStudySpan(t, MarkerSyntaxHelp(marker))
	want := "+kubebuilder:printcolumn:name=<string>,type=<string>[,priority=<int>]"
	if got != want {
		t.Fatalf("MarkerSyntaxHelp() = %q, want %q", got, want)
	}
}

func TestIndentedStudy(t *testing.T) {
	if got, want := renderStudySpan(t, Indented(1, Text("line1\nline2"))), "\tline1\n\tline2"; got != want {
		t.Fatalf("Indented() = %q, want %q", got, want)
	}
}

func TestMarkersSummaryStudy(t *testing.T) {
	markers := []help.MarkerDoc{
		{
			Name:   "paths",
			Target: "package",
			DetailedHelp: help.DetailedHelp{
				Summary: "selects packages to load",
			},
			Fields: []help.FieldHelp{
				{
					Argument: help.Argument{
						Type:     "slice",
						ItemType: &help.Argument{Type: "string"},
					},
				},
			},
		},
		{
			Name:   "output:crd:artifacts",
			Target: "package",
			DetailedHelp: help.DetailedHelp{
				Summary: "writes CRD output",
			},
			Fields: []help.FieldHelp{
				{
					Name: "config",
					Argument: help.Argument{
						Type: "string",
					},
				},
			},
		},
	}

	got := renderStudySpan(t, MarkersSummary("Generator Options", markers))
	want := "\nGenerator Options\n\n+paths=<[]string>                      package  selects packages to load  \n+output:crd:artifacts:config=<string>  package  writes CRD output         \n\n"
	if got != want {
		t.Fatalf("MarkersSummary() = %q, want %q", got, want)
	}
}

func TestMarkersDetailsStudy(t *testing.T) {
	marker := help.MarkerDoc{
		Name:   "kubebuilder:printcolumn",
		Target: "type",
		DetailedHelp: help.DetailedHelp{
			Summary: "adds a printer column",
			Details: "The column is shown by kubectl get output.",
		},
		Fields: []help.FieldHelp{
			{
				Name: "name",
				Argument: help.Argument{
					Type: "string",
				},
				DetailedHelp: help.DetailedHelp{
					Summary: "column header",
					Details: "This is the title users see in tables.",
				},
			},
			{
				Name: "priority",
				Argument: help.Argument{
					Type:     "int",
					Optional: true,
				},
				DetailedHelp: help.DetailedHelp{
					Summary: "visibility priority",
				},
			},
		},
	}

	t.Run("compact", func(t *testing.T) {
		got := renderStudySpan(t, MarkersDetails(false, "CRD Markers", []help.MarkerDoc{marker}))
		want := "\nCRD Markers\n\n\n+kubebuilder:printcolumn type\n\tadds a printer column\n\tname=<string>     column header        \n\t[priority=<int>]  visibility priority  \n"
		if got != want {
			t.Fatalf("MarkersDetails(false) = %q, want %q", got, want)
		}
	})

	t.Run("full detail", func(t *testing.T) {
		got := renderStudySpan(t, MarkersDetails(true, "CRD Markers", []help.MarkerDoc{marker}))
		want := "\nCRD Markers\n\n\n+kubebuilder:printcolumn type\n\tadds a printer column\n\tThe column is shown by kubectl get output.\n\n\tname=<string>\n\t\tcolumn header\n\t\tThis is the title users see in tables.\n\n\t[priority=<int>]\n\t\tvisibility priority\n"
		if got != want {
			t.Fatalf("MarkersDetails(true) = %q, want %q", got, want)
		}
	})
}

func TestTableCalculatorStudy(t *testing.T) {
	calc := &TableCalculator{Padding: 2, MaxWidth: 6}
	calc.AddRowSizes(3, 10)
	calc.AddRowSizes(4, 1)

	widths := calc.ColumnWidths()
	if got, want := len(widths), 2; got != want {
		t.Fatalf("len(ColumnWidths()) = %d, want %d", got, want)
	}
	if got, want := widths[0], 6; got != want {
		t.Fatalf("first width = %d, want %d", got, want)
	}
	if got, want := widths[1], 6; got != want {
		t.Fatalf("second width = %d, want %d", got, want)
	}
}
