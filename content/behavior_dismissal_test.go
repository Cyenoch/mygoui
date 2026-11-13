package content

import (
	"fmt"
	"testing"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// Escape in a viewer's empty search field is left to the application, so a
// dialog around the viewer still closes.
func TestViewerSearchEscapeOnEmptyQueryReachesAncestor(t *testing.T) {
	viewers := map[string]func(c *ui.Context){
		"log": func(c *ui.Context) {
			LogViewer(c, LogViewerOptions{
				Line: func(i int) string { return fmt.Sprintf("entry %d", i) }, Count: 20,
				Height: 200, Toolbar: true,
			})
		},
		"json": func(c *ui.Context) {
			_, _ = JSONViewer(c, JSONViewerOptions{Source: jsonSample, Toolbar: true})
		},
	}
	for name, viewer := range viewers {
		t.Run(name, func(t *testing.T) {
			escapes := 0
			tt := ui.NewTester(func(c *ui.Context) {
				theme.Apply(c)
				ui.Column(c).Fill().Padding(16).Children(func() { viewer(c) }).
					OnShortcut(0, ui.KeyEscape, func() { escapes++ })
			}, 480, 480)
			if err := tt.Click("Search"); err != nil {
				t.Fatal(err)
			}
			tt.Key(0, ui.KeyEscape)
			if escapes != 1 {
				t.Fatalf("Escape in an empty search reached the ancestor %d times, want 1", escapes)
			}
		})
	}
}
