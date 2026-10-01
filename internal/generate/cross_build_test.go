package generate

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestGeneratedAccessorsCrossCompile(t *testing.T) {
	dir := samplePackage(t)
	input := "package sample\n\nimport \"unsafe\"\n\ntype CrossNotice struct {\n\tWord [unsafe.Sizeof(int(0))]byte `kakaowork:\"text\" slack:\"rich_text\" googlechat:\"textParagraph\"`\n\tPointer *[unsafe.Sizeof(uintptr(0))]byte `kakaowork:\"text\" slack:\"rich_text\" googlechat:\"textParagraph\"`\n\tOS [platformBytes]byte `kakaowork:\"text\" slack:\"rich_text\" googlechat:\"textParagraph\"`\n}\n"
	if err := os.WriteFile(filepath.Join(dir, "cross.go"), []byte(input), 0644); err != nil {
		t.Fatal(err)
	}
	for name, data := range map[string]string{
		"platform_linux.go":   "package sample\nconst platformBytes = 1\n",
		"platform_windows.go": "package sample\nconst platformBytes = 2\n",
	} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(data), 0644); err != nil {
			t.Fatal(err)
		}
	}
	for _, target := range []struct{ os, arch string }{
		{"linux", "386"},
		{"windows", "amd64"},
	} {
		t.Run(target.os+"/"+target.arch, func(t *testing.T) {
			t.Setenv("GOOS", target.os)
			t.Setenv("GOARCH", target.arch)
			t.Setenv("CGO_ENABLED", "0")
			output := "zz_cross_" + target.os + "_" + target.arch + ".go"
			if err := Run(dir, []string{"CrossNotice"}, output); err != nil {
				t.Fatal(err)
			}
			command := exec.Command("go", "test", "-mod=readonly", "-c", "-o", filepath.Join(t.TempDir(), "sample.test"), ".")
			command.Dir = dir
			if output, err := command.CombinedOutput(); err != nil {
				t.Fatalf("target compilation: %v\n%s", err, output)
			}
		})
	}
}

func TestTargetArchitectureFromGoEnvironmentFile(t *testing.T) {
	dir := samplePackage(t)
	environment := filepath.Join(t.TempDir(), "goenv")
	if err := os.WriteFile(environment, []byte("GOARCH=386\nGOOS=linux\nCGO_ENABLED=0\n"), 0644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("GOENV", environment)
	t.Setenv("GOARCH", "")
	t.Setenv("GOOS", "")
	t.Setenv("CGO_ENABLED", "")
	if err := os.WriteFile(filepath.Join(dir, "overflow.go"), []byte("package sample\nvar targetInteger int = 1 << 40\n"), 0644); err != nil {
		t.Fatal(err)
	}
	err := Run(dir, []string{"Notice"}, "zz_convertago.gen.go")
	if err == nil || !strings.Contains(err.Error(), "overflows") {
		t.Fatalf("32-bit integer overflow must fail generation: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "zz_convertago.gen.go")); !os.IsNotExist(err) {
		t.Fatal("failed target type checking must not write output")
	}
}

func TestHostGeneratedAccessorsCrossCompile(t *testing.T) {
	dir := samplePackage(t)
	if err := Run(dir, []string{"Notice", "VerboseNotice", "LayoutNotice"}, "zz_convertago.gen.go"); err != nil {
		t.Fatal(err)
	}
	for _, target := range []struct{ os, arch string }{
		{"linux", "386"},
		{"linux", "arm64"},
		{"windows", "amd64"},
	} {
		t.Run(target.os+"/"+target.arch, func(t *testing.T) {
			t.Setenv("GOOS", target.os)
			t.Setenv("GOARCH", target.arch)
			t.Setenv("CGO_ENABLED", "0")
			command := exec.Command("go", "test", "-mod=readonly", "-c", "-o", filepath.Join(t.TempDir(), "sample.test"), ".")
			command.Dir = dir
			if output, err := command.CombinedOutput(); err != nil {
				t.Fatalf("host-generated target compilation: %v\n%s", err, output)
			}
		})
	}
}
