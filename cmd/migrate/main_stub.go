//go:build !sqlite

package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Fprintln(os.Stderr, "migrate: this binary was built without SQLite support; rebuild with -tags sqlite")
	os.Exit(1)
}
