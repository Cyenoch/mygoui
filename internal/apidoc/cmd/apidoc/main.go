// Command apidoc writes the gallery's generated API and example tables. Run
// it from the module root: go run ./internal/apidoc/cmd/apidoc.
package main

import (
	"fmt"
	"os"

	"github.com/Cyenoch/mygoui/internal/apidoc"
)

func main() {
	for file, generate := range map[string]func(string) ([]byte, error){
		"internal/gallery/api_gen.go":            apidoc.Generate,
		"internal/gallery/example_source_gen.go": apidoc.GenerateExamples,
	} {
		source, err := generate(".")
		if err == nil {
			err = os.WriteFile(file, source, 0o644)
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, "apidoc:", file+":", err)
			os.Exit(1)
		}
	}
}
