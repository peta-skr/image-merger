package main

import (
	"flag"
	"fmt"
	"image-meger/collector"
	"log"
	"os"
	"strings"
)

type multiFlag []string

func (m *multiFlag) String() string { return strings.Join(*m, ",") }
func (m *multiFlag) Set(v string) error {
	*m = append(*m, v)
	return nil
}

func main() {
	var srcs multiFlag
	var dst string
	var dryRun bool

	flag.Var(&srcs, "src", "source directory (can be specified multiple times)")
	flag.StringVar(&dst, "dst", "", "destination directory")
	flag.BoolVar(&dryRun, "dry-run", false, "print plan only")
	flag.Parse()

	if len(srcs) == 0 || dst == "" {
		fmt.Println("Usage: img-collector --src <dir> [--src <dir> ...] --dst <dir> [--dry-run]")
		os.Exit(1)
	}

	opt := collector.Options{
		Sources: srcs,
		Dest:    dst,
		DryRun:  dryRun,
	}

	res, err := collector.Run(opt)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Done. copied=%d skipped=%d failed=%d\n", res.Copied, res.Skipped, res.Failed)
}
