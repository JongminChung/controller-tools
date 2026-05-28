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

package yaml

import (
	"testing"

	yamlv3 "gopkg.in/yaml.v3"
)

func TestValueInMappingStudy(t *testing.T) {
	root, err := ToYAML(map[string]any{
		"spec": map[string]any{
			"replicas": int64(3),
		},
	})
	if err != nil {
		t.Fatalf("ToYAML() error = %v", err)
	}

	mapping, found, err := GetNode(root, "spec")
	if err != nil {
		t.Fatalf("GetNode() error = %v", err)
	}
	if !found {
		t.Fatal("GetNode() did not find spec")
	}

	node, err := ValueInMapping(mapping, "replicas")
	if err != nil {
		t.Fatalf("ValueInMapping() error = %v", err)
	}
	if node == nil {
		t.Fatal("ValueInMapping() returned nil node")
	}
	if node.Value != "3" {
		t.Fatalf("ValueInMapping() value = %q, want %q", node.Value, "3")
	}
}

func TestSetAndDeleteNodeStudy(t *testing.T) {
	root, err := ToYAML(map[string]any{
		"metadata": map[string]any{
			"name": "demo",
		},
	})
	if err != nil {
		t.Fatalf("ToYAML() error = %v", err)
	}

	value, err := ToYAML("default")
	if err != nil {
		t.Fatalf("ToYAML() error = %v", err)
	}
	value, _, err = asCloseAsPossible(value)
	if err != nil {
		t.Fatalf("asCloseAsPossible() error = %v", err)
	}

	if err := SetNode(root, *value, "metadata", "namespace"); err != nil {
		t.Fatalf("SetNode() error = %v", err)
	}

	namespace, found, err := GetNode(root, "metadata", "namespace")
	if err != nil {
		t.Fatalf("GetNode() error = %v", err)
	}
	if !found {
		t.Fatal("GetNode() did not find metadata.namespace")
	}
	if namespace.Value != "default" {
		t.Fatalf("metadata.namespace = %q, want %q", namespace.Value, "default")
	}

	if err := DeleteNode(root, "metadata", "name"); err != nil {
		t.Fatalf("DeleteNode() error = %v", err)
	}

	_, found, err = GetNode(root, "metadata", "name")
	if err != nil {
		t.Fatalf("GetNode() after delete error = %v", err)
	}
	if found {
		t.Fatal("GetNode() still found deleted metadata.name")
	}

	metadata, found, err := GetNode(root, "metadata")
	if err != nil {
		t.Fatalf("GetNode() metadata error = %v", err)
	}
	if !found {
		t.Fatal("GetNode() did not find metadata")
	}
	if metadata.Kind != yamlv3.MappingNode {
		t.Fatalf("metadata.Kind = %v, want %v", metadata.Kind, yamlv3.MappingNode)
	}
}
