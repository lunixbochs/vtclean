package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/lunixbochs/vtclean"
	"os"
)

func main() {
	color := flag.Bool("color", false, "enable color")
	flag.Parse()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		fmt.Println(vtclean.Clean(scanner.Text(), *color))
	}
}
