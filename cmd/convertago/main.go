// Command convertago generates direct field accessors for tagged structs.
//
// Run it from the source package, usually through go:generate:
//
//	convertago -type Notice[,Other...] [-output zz_convertago.gen.go]
//
// The -type flag lists the root struct types. The -output flag names the
// generated file in the source package and defaults to zz_convertago.gen.go.
// See docs/generation.md for the supported types and the generated contract.
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/JSYoo5B/convertago/internal/generate"
)

func main() {
	names := flag.String("type", "", "comma-separated struct type names")
	output := flag.String("output", "zz_convertago.gen.go", "generated Go filename")
	flag.Parse()
	if *names == "" || flag.NArg() != 0 {
		flag.Usage()
		os.Exit(2)
	}
	if err := generate.Run(".", strings.Split(*names, ","), *output); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
