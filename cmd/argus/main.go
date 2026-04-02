package main

import (
	"fmt"
	"os"

	"github.com/gentleman-programming/argus/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
