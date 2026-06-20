package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/fran-codigo/tui-images/internal/tui"
)

var version = "dev"

func main() {
	quality := flag.Int("q", 75, "Compression quality (1-100)")
	showVersion := flag.Bool("v", false, "Show version")
	flag.Parse()

	if *showVersion {
		fmt.Printf("tui-images %s\n", version)
		os.Exit(0)
	}

	if *quality < 1 || *quality > 100 {
		fmt.Fprintf(os.Stderr, "Error: quality must be between 1 and 100\n")
		os.Exit(1)
	}

	if err := tui.Run(*quality); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
