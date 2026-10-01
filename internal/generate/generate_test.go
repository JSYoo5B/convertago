package generate

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func samplePackage(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	root, err := filepath.Abs("../..")
	if err != nil {
		t.Fatal(err)
	}
	module := "module example.com/sample\n\ngo 1.23.0\n\nrequire " + modulePath + " v0.0.0\nreplace " + modulePath + " => " + filepath.ToSlash(root) + "\n"
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(module), 0644); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"input.go", "parity_test.go"} {
		data, err := os.ReadFile(filepath.Join("testdata/sample", name))
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, name), data, 0644); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

func TestGeneratedAccessorsCompileAndMatchReflection(t *testing.T) {
	dir := samplePackage(t)
	if err := Run(dir, []string{"Notice", "VerboseNotice", "LayoutNotice"}, "zz_convertago.gen.go"); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "zz_convertago.gen.go")
	first, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(first), "value.Name") || strings.Contains(string(first), `"reflect"`) || strings.Contains(string(first), "reflect.Value") {
		t.Fatal("known fields must use direct access")
	}
	if err := Run(dir, []string{"Notice", "VerboseNotice", "LayoutNotice"}, "zz_convertago.gen.go"); err != nil {
		t.Fatal(err)
	}
	second, err := os.ReadFile(path)
	if err != nil || string(first) != string(second) {
		t.Fatalf("generation must be deterministic: %v", err)
	}
	command := exec.Command("go", "generate", ".")
	command.Dir = dir
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("go generate directive: %v\n%s", err, output)
	}
	command = exec.Command("go", "test", "-mod=readonly", "-race", "-count=1", ".")
	command.Dir = dir
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("generated package: %v\n%s", err, output)
	}
}

func TestRegenerationAfterSourceChange(t *testing.T) {
	dir := samplePackage(t)
	if err := Run(dir, []string{"Notice"}, "zz_convertago.gen.go"); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "input.go")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	data = []byte(strings.ReplaceAll(string(data), "Name", "Person"))
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatal(err)
	}
	if err := Run(dir, []string{"Notice"}, "zz_convertago.gen.go"); err != nil {
		t.Fatalf("stale generated access must not block regeneration: %v", err)
	}
	generated, err := os.ReadFile(filepath.Join(dir, "zz_convertago.gen.go"))
	if err != nil || !strings.Contains(string(generated), "value.Person") || strings.Contains(string(generated), "value.Name") {
		t.Fatalf("regenerated field access was not updated: %v", err)
	}
}

func TestInvalidTagsAreCaughtBeforeWriting(t *testing.T) {
	dir := samplePackage(t)
	data, err := os.ReadFile(filepath.Join(dir, "input.go"))
	if err != nil {
		t.Fatal(err)
	}
	data = []byte(strings.Replace(string(data), `slack:"header"`, `slack:"haeder;optional"`, 1))
	if err := os.WriteFile(filepath.Join(dir, "input.go"), data, 0644); err != nil {
		t.Fatal(err)
	}
	err = Run(dir, []string{"Notice"}, "zz_convertago.gen.go")
	if err == nil || !strings.Contains(err.Error(), "slack $.Title") || !strings.Contains(err.Error(), "haeder") {
		t.Fatalf("error = %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "zz_convertago.gen.go")); !os.IsNotExist(err) {
		t.Fatal("failed generation must not write a partial file")
	}
}

func TestGeneratorPreservesHandwrittenFiles(t *testing.T) {
	dir := samplePackage(t)
	if err := Run(dir, []string{"Notice"}, "input.go"); err == nil || !strings.Contains(err.Error(), "non-generated") {
		t.Fatalf("error = %v", err)
	}
}

func TestGeneratorRejectsUnsupportedNativeNesting(t *testing.T) {
	dir := samplePackage(t)
	path := filepath.Join(dir, "input.go")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	data = []byte(strings.Replace(string(data), "type GoogleColumn struct {\n\tText string `googlechat:\"textParagraph\"`", "type GoogleColumn struct {\n\tText string `googlechat:\"divider\"`", 1))
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatal(err)
	}
	err = Run(dir, []string{"LayoutNotice"}, "zz_convertago.gen.go")
	if err == nil || !strings.Contains(err.Error(), "column does not accept divider") || !strings.Contains(err.Error(), "$.Cards[].Card.Sections[].Columns.Columns[].Text") {
		t.Fatalf("error=%v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "zz_convertago.gen.go")); !os.IsNotExist(err) {
		t.Fatal("invalid nesting wrote output")
	}
}
