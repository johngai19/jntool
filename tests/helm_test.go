package helm_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/johngai19/jntool/jntool/helm"
	"gopkg.in/yaml.v3"
)

func writeTempYAML(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "values.yaml")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp YAML: %v", err)
	}
	return path
}

func TestValuesCommandJSON(t *testing.T) {
	yamlContent := `
foo: bar
nested:
  key: value
`
	file := writeTempYAML(t, yamlContent)
	var buf bytes.Buffer

	cmd := helm.HelmCmd
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	cmd.SetArgs([]string{"values", file, "-o", "json"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() failed: %v, output: %s", err, buf.String())
	}

	var data map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &data); err != nil {
		t.Fatalf("invalid JSON output: %v\n%s", err, buf.String())
	}

	want := map[string]interface{}{
		"foo": "bar",
		"nested": map[string]interface{}{
			"key": "value",
		},
	}
	if !reflect.DeepEqual(data, want) {
		t.Errorf("JSON output = %#v; want %#v", data, want)
	}
}

func TestValuesCommandYAML(t *testing.T) {
	yamlContent := `
foo: bar
nested:
  list:
    - a
    - b
`
	file := writeTempYAML(t, yamlContent)
	var buf bytes.Buffer

	cmd := helm.HelmCmd
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	cmd.SetArgs([]string{"values", file, "-o", "yaml"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() failed: %v, output: %s", err, buf.String())
	}

	var data map[string]interface{}
	if err := yaml.Unmarshal(buf.Bytes(), &data); err != nil {
		t.Fatalf("invalid YAML output: %v\n%s", err, buf.String())
	}

	want := map[string]interface{}{
		"foo": "bar",
		"nested": map[string]interface{}{
			"list": []interface{}{"a", "b"},
		},
	}
	if !reflect.DeepEqual(data, want) {
		t.Errorf("YAML output = %#v; want %#v", data, want)
	}
}

func TestValuesCommandUnsupportedFormat(t *testing.T) {
	yamlContent := `foo: bar`
	file := writeTempYAML(t, yamlContent)
	var buf bytes.Buffer

	cmd := helm.HelmCmd
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	cmd.SetArgs([]string{"values", file, "-o", "xml"})

	err := cmd.Execute()
	if err == nil {
		t.Fatalf("expected error for unsupported format, got nil; output: %s", buf.String())
	}
}
