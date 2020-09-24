package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("list or count subcommand is required")
		os.Exit(1)
	}

	cf := flag.String("config", "/opt/scanner/config.json", "path for the configuration file")

	flag.Parse()

	if *cf == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	scan(*cf)
}
