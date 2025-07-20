package main

import (
	"fmt"
	"time"

	"github.com/ShivangSrivastava/fynd/cmd"
)

func main() {
	start := time.Now()
	cmd.Execute()
	fmt.Printf("\n📊 Runtime: %v\n", time.Since(start))
}
