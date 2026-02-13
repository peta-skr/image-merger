package main

import (
	"fmt"
	"image-meger/collector"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: image-collector <src> <dst>")
		return
	}

	src := os.Args[1]
	dst := os.Args[2]

	err := collector.CollectImages(src, dst)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Done!")
}
