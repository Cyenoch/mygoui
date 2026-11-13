package form

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

// Escape in an empty search field is left to the surrounding window, such as
// a dialog that closes on it.
func TestSearchFieldEscapeOnEmptyValueReachesAncestor(t *testing.T) {
	value := ""
	escapes := 0
	tt := formTextTester(t, func(c *ui.Context) {
		ui.Column(c).Children(func() {
			SearchField(c, &value, SearchFieldOptions{Label: "Find"})
		}).OnShortcut(0, ui.KeyEscape, func() { escapes++ })
	})
	formTextClick(t, tt, "Find")
	tt.Key(0, ui.KeyEscape)
	if escapes != 1 {
		t.Fatalf("Escape in an empty search reached the ancestor %d times, want 1", escapes)
	}
}

// Escape discards an unfinished tag and keeps the committed ones.
func TestTagsInputEscapeClearsDraft(t *testing.T) {
	tags := []string{"go"}
	tt := formTextTester(t, func(c *ui.Context) {
		TagsInput(c, &tags, TagsInputOptions{Label: "Tags"})
	})
	formTextClick(t, tt, "Tags")
	tt.Type("dra")
	tt.Key(0, ui.KeyEscape)
	info, _ := tt.AccessByRole(ui.RoleTextField, "Tags")
	if info.Value != "" || len(tags) != 1 {
		t.Fatalf("after Escape the draft is %q and tags %v", info.Value, tags)
	}
}

// Escape leaves editing with the focus on the text, so Enter edits again.
func TestEditableTextEscapeReturnsFocusToText(t *testing.T) {
	name, changes := "Draft", 0
	tt := editableTester(&name, &changes)
	formTextClick(t, tt, "Name")
	tt.Type("x")
	tt.Key(0, ui.KeyEscape)
	tt.Frame()
	if !tt.Focused("Name") {
		t.Fatal("the focus did not return to the text after Escape")
	}
	tt.Key(0, ui.KeyEnter)
	tt.Frame()
	if info, ok := tt.AccessByRole(ui.RoleTextField, "Name"); !ok || !info.Focused {
		t.Fatal("Enter on the text did not start editing again")
	}
}
