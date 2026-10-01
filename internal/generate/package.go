package generate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

type packageInfo struct {
	Dir        string
	Name       string
	ImportPath string
	GoFiles    []string
	CgoFiles   []string
	Export     string
}

func loadPackage(dir, output string) (*types.Package, error) {
	command := exec.Command("go", "list", "-mod=readonly", "-e", "-export", "-deps", "-json", ".")
	command.Dir = dir
	var stderr bytes.Buffer
	command.Stderr = &stderr
	data, err := command.Output()
	if err != nil {
		return nil, fmt.Errorf("go list: %w: %s", err, stderr.String())
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	exports := make(map[string]string)
	var current packageInfo
	for {
		var info packageInfo
		if err := decoder.Decode(&info); err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		exports[info.ImportPath] = info.Export
		current = info
	}
	if len(current.CgoFiles) != 0 {
		return nil, fmt.Errorf("convertago: generation in a cgo package is not supported")
	}
	fset := token.NewFileSet()
	var files []*ast.File
	for _, name := range current.GoFiles {
		if name == output {
			continue
		}
		file, err := parser.ParseFile(fset, filepath.Join(current.Dir, name), nil, parser.ParseComments)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}
	lookup := func(path string) (io.ReadCloser, error) {
		file := exports[path]
		if file == "" {
			return nil, fmt.Errorf("no compiled export data for %s", path)
		}
		return os.Open(file)
	}
	config := types.Config{Importer: importer.ForCompiler(fset, "gc", lookup), Sizes: types.SizesFor("gc", runtime.GOARCH)}
	return config.Check(current.ImportPath, fset, files, nil)
}
