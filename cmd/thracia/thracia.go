package main

import (
	"fmt"
	"os"

	"github.com/delphinus/go-thracia"
)

func main() {
	if err := thracia.New().Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}
	fmt.Println("done successfully")
}
