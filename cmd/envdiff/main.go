package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/user/envdiff/internal/differ"
	"github.com/user/envdiff/internal/parser"
	"github.com/user/envdiff/internal/reporter"
)

func main() {
	var (
		format  = flag.String("format", "text", "Output format: text or json")
		leftLabel  = flag.String("left-label", "", "Label for the left file (defaults to filename)")
		rightLabel = flag.String("right-label", "", "Label for the right file (defaults to filename)")
	)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: envdiff [flags] <left.env> <right.env>\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		flag.Usage()
		os.Exit(1)
	}

	leftPath, rightPath := args[0], args[1]

	left, err := parser.ParseFile(leftPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading %s: %v\n", leftPath, err)
		os.Exit(1)
	}

	right, err := parser.ParseFile(rightPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading %s: %v\n", rightPath, err)
		os.Exit(1)
	}

	if *leftLabel == "" {
		*leftLabel = leftPath
	}
	if *rightLabel == "" {
		*rightLabel = rightPath
	}

	diffs := differ.Compare(left, right)
	rpt := reporter.NewReport(*leftLabel, *rightLabel, diffs)

	switch *format {
	case "json":
		if err := reporter.WriteJSON(os.Stdout, rpt); err != nil {
			fmt.Fprintf(os.Stderr, "error writing JSON: %v\n", err)
			os.Exit(1)
		}
	case "text":
		if err := reporter.WriteText(os.Stdout, rpt); err != nil {
			fmt.Fprintf(os.Stderr, "error writing text: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "unknown format %q — use \"text\" or \"json\"\n", *format)
		os.Exit(1)
	}

	if len(diffs) > 0 {
		os.Exit(2)
	}
}
