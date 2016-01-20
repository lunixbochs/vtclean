package main

import (
	"flag"
	"github.com/lunixbochs/vtclean"
	"io"
	"os"
)

func main() {
	color := flag.Bool("color", false, "enable color")
	flag.Parse()

	stdout := vtclean.NewWriter(os.Stdout, *color)
	io.Copy(stdout, os.Stdin)
}
