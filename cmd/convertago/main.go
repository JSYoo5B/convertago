// Command convertago generates direct field accessors for tagged structs.
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
