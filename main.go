package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	htmlCmd := flag.NewFlagSet("html", flag.ExitOnError)
	htmlTemplate := htmlCmd.String("template", "", "template")

	formatCmd := flag.NewFlagSet("format", flag.ExitOnError)

	if len(os.Args) < 2 {
		fmt.Println("expected 'html' or 'format' subcommands")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "html":
		htmlCmd.Parse(os.Args[2:])
		renderAsHTML(*htmlTemplate)
	case "format":
		formatCmd.Parse(os.Args[2:])
		files := formatCmd.Args()
		format(files)
	default:
		fmt.Println("expected 'html' or 'format' subcommands")
		os.Exit(1)
	}
}
