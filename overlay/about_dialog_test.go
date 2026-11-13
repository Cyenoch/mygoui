package overlay

import (
	"strings"
	"testing"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func TestAboutDialogKeepsIdentityInPlaceWhileSectionsScroll(t *testing.T) {
	open, visited := false, 0
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Button(c, "About").OnClick(func() { open = true })
		AboutDialog(c, &open, AboutDialogOptions{
			Icon: func() { ui.Box(c).Size(48, 48).Label("App icon") },
			Name: "Notes", Version: "Version 2.4.1", Description: "Plain text, everywhere.",
			Links: []AboutLink{{Label: "Website", OnClick: func() { visited++ }}},
			Sections: []AboutSection{
				{Title: "Credits", Body: "Ada Lovelace"},
				{Title: "Licence", Body: strings.Repeat("Permission is hereby granted.\n", 80)},
			},
		})
	}, 800, 480)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	if err := tt.Click("About"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if _, ok := tt.AccessByRole(ui.RoleDialog, "Notes"); !ok {
		t.Fatal("the dialog is not named by the application")
	}
	if err := tt.Click("Website"); err != nil {
		t.Fatal(err)
	}
	if visited != 1 || !open {
		t.Fatalf("link: visited=%d open=%v", visited, open)
	}
	name, _ := tt.Find("Notes")
	credits, _ := tt.Find("Credits")
	tt.Wheel(credits.X+10, credits.Y+10, 0, 600)
	tt.Frame()
	after, _ := tt.Find("Notes")
	scrolled, _ := tt.Find("Credits")
	if after != name || name.Y < 0 {
		t.Fatalf("name moved from %+v to %+v", name, after)
	}
	if scrolled.Y >= credits.Y {
		t.Fatalf("sections did not scroll: %+v then %+v", credits, scrolled)
	}
	tt.Key(0, ui.KeyEscape)
	tt.Frame()
	if open || !tt.Focused("About") {
		t.Fatalf("Escape: open=%v opener focused=%v", open, tt.Focused("About"))
	}
}

// A bug report needs the version, and a licence is long enough to push
// everything else out of sight.
func TestAboutDialogCopiesInfoAndFoldsLongSections(t *testing.T) {
	open := true
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		AboutDialog(c, &open, AboutDialogOptions{
			Name: "Notes", Version: "Version 2.4.1", CopyInfo: "Notes 2.4.1 (318)",
			Sections: []AboutSection{{Title: "Licence", Body: "Permission is hereby granted.", Collapsible: true}},
		})
	}, 800, 480)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	tt.Frame()
	if tt.HasText("Permission is hereby granted.") {
		t.Fatal("a collapsible section shows its body before it is opened")
	}
	if err := tt.Click("Copy"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if got := tt.Clipboard(); got != "Notes 2.4.1 (318)" {
		t.Fatalf("copied %q", got)
	}
	if err := tt.Click("Licence"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if !tt.HasText("Permission is hereby granted.") {
		t.Fatal("the section did not open")
	}
}
