package main

import (
	"context"
	"fmt"
	"os"

	"github.com/delphinus/go-thracia"
)

func main() {
	ctx := context.Background()
	if err := thracia.Run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}
	fmt.Println("done successfully")
}
