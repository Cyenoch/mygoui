// Command apidoc writes the gallery's generated API table. Run it from the
// module root: go run ./internal/apidoc/cmd/apidoc.
package main

import (
	"fmt"
	"os"

	"github.com/Cyenoch/mygoui/internal/apidoc"
)

func main() {
	source, err := apidoc.Generate(".")
	if err == nil {
		err = os.WriteFile("internal/gallery/api_gen.go", source, 0o644)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "apidoc:", err)
		os.Exit(1)
	}
}
