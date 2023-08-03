package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	content, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
		os.Exit(1)
	}
	fmt.Fprint(os.Stdout, os.ExpandEnv(string(content)))
}
