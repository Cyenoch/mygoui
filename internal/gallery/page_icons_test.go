package gallery

import (
	"fmt"
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestIconBrowserFindsAndCopiesIcons(t *testing.T) {
	state := New()
	state.Selected = "icons"
	tt := ui.NewTester(func(c *ui.Context) { View(c, state) }, 1180, 900)
	tt.SetPreferences(ui.Preferences{TextScale: 1, ReduceMotion: true})
	tt.Frame()
	tt.Frame()
	st := iconsStateOf(state)
	all := len(st.results)
	if all < 1500 || st.built == 0 || st.built > 12 {
		t.Fatalf("built %d rows for %d icons; only rows near the viewport may be built", st.built, all)
	}

	if err := tt.Click("Search icons"); err != nil {
		t.Fatal(err)
	}
	tt.Type("arrow-right")
	tt.Frame()
	if len(st.results) == 0 || len(st.results) >= all {
		t.Fatalf("searching an Iconify ID left %d of %d icons", len(st.results), all)
	}
	if !tt.HasText("ArrowRight") {
		t.Fatal("the search did not show ArrowRight")
	}
	if count := fmt.Sprintf("%d icons", len(st.results)); !tt.HasText(count) {
		t.Fatalf("the list does not report %q", count)
	}
	if err := tt.Click("ArrowRight"); err != nil {
		t.Fatal(err)
	}
	if got := tt.Clipboard(); got != "lucide.ArrowRight()" {
		t.Fatalf("activating an icon copied %q", got)
	}
	tt.Frame()
	if !tt.HasText("Copied lucide.ArrowRight()") {
		t.Fatal("the copy was not reported")
	}

	st.query = "no-such-icon-anywhere"
	tt.Frame()
	if !tt.HasText("No icons found") {
		t.Fatal("a search without matches showed no empty state")
	}
}

// An alias is found by its own name but shown as the drawing it repeats.
func TestIconBrowserShowsEachDrawingOnce(t *testing.T) {
	st := iconsStateOf(New())
	st.query = "alert-circle"
	st.filter()
	seen := map[string]bool{}
	for _, icon := range st.results {
		if icon.AliasOf != "" || seen[icon.Name] {
			t.Fatalf("%s is listed as an alias or twice", icon.Name)
		}
		seen[icon.Name] = true
	}
	if len(st.results) == 0 {
		t.Fatal("an alias's ID found nothing")
	}
}
