package form

import (
	"slices"
	"testing"

	"github.com/Cyenoch/mygoui/locale"
	"github.com/egoist/mygo/ui"
)

func tagsTester(t *testing.T, tags *[]string, o TagsInputOptions) *ui.Tester {
	t.Helper()
	o.Label = "Tags"
	return formTextTester(t, func(c *ui.Context) { TagsInput(c, tags, o) })
}

func TestTagsInputCommitsOnEnterAndComma(t *testing.T) {
	var tags []string
	changes := 0
	tester := tagsTester(t, &tags, TagsInputOptions{OnChange: func() { changes++ }})
	formTextClick(t, tester, "Tags")
	tester.Type("go")
	tester.Key(0, ui.KeyEnter)
	tester.Type(" ui ,nat")
	if want := []string{"go", "ui"}; !slices.Equal(tags, want) || changes != 2 {
		t.Fatalf("tags = %q after %d changes, want %q after 2", tags, changes, want)
	}
	tester.Type("ive")
	tester.Key(0, ui.KeyEnter)
	if want := []string{"go", "ui", "native"}; !slices.Equal(tags, want) {
		t.Fatalf("tags = %q, want %q: the text after a comma stays in the editor", tags, want)
	}
	if _, ok := tester.AccessByRole(ui.RoleButton, locale.M().Remove+" native"); !ok {
		t.Fatal("a token has no named remove button")
	}
}

func TestTagsInputBackspaceRemovesLastOnlyWhenEmpty(t *testing.T) {
	tags := []string{"a", "b"}
	tester := tagsTester(t, &tags, TagsInputOptions{})
	formTextClick(t, tester, "Tags")
	tester.Type("x")
	tester.Key(0, ui.KeyBackspace)
	if len(tags) != 2 {
		t.Fatalf("Backspace over typed text removed a tag: %q", tags)
	}
	tester.Key(0, ui.KeyBackspace)
	if want := []string{"a"}; !slices.Equal(tags, want) {
		t.Fatalf("tags = %q, want %q", tags, want)
	}
}

func TestTagsInputRemoveButtonRemovesItsToken(t *testing.T) {
	tags := []string{"a", "b", "c"}
	changes := 0
	tester := tagsTester(t, &tags, TagsInputOptions{OnChange: func() { changes++ }})
	formTextClick(t, tester, locale.M().Remove+" b")
	if want := []string{"a", "c"}; !slices.Equal(tags, want) || changes != 1 {
		t.Fatalf("tags = %q after %d changes, want %q after 1", tags, changes, want)
	}
}

func TestTagsInputKeepsAndFlagsARefusedDraft(t *testing.T) {
	tags := []string{"a"}
	tester := tagsTester(t, &tags, TagsInputOptions{})
	formTextClick(t, tester, "Tags")
	tester.Type("a")
	tester.Key(0, ui.KeyEnter)
	if info, _ := tester.AccessByRole(ui.RoleTextField, "Tags"); info.Value != "a" || len(tags) != 1 {
		t.Fatalf("a duplicate left draft %q and tags %q, want the draft kept", info.Value, tags)
	}
	tester.Type("b")
	tester.Key(0, ui.KeyEnter)
	if want := []string{"a", "ab"}; !slices.Equal(tags, want) {
		t.Fatalf("tags = %q, want %q after the draft changed", tags, want)
	}
	// A comma refuses a duplicate the same way.
	tester.Type("a,")
	if info, _ := tester.AccessByRole(ui.RoleTextField, "Tags"); info.Value != "a" || len(tags) != 2 {
		t.Fatalf("a comma after a duplicate left draft %q and tags %q", info.Value, tags)
	}
}

func TestTagsInputStopsAtMax(t *testing.T) {
	tags := []string{"a", "b"}
	tester := tagsTester(t, &tags, TagsInputOptions{Max: 2})
	formTextClick(t, tester, "Tags")
	tester.Type("c")
	tester.Key(0, ui.KeyEnter)
	if len(tags) != 2 {
		t.Fatalf("tags = %q, want no tag past Max", tags)
	}
	tester.Key(0, ui.KeyBackspace)
	if want := []string{"a"}; !slices.Equal(tags, want) {
		t.Fatalf("Backspace at Max left %q, want %q", tags, want)
	}
	tester.Type("c")
	tester.Key(0, ui.KeyEnter)
	if want := []string{"a", "c"}; !slices.Equal(tags, want) {
		t.Fatalf("tags = %q, want %q once there is room", tags, want)
	}
}

func TestTagsInputCommitsDraftWhenFocusLeaves(t *testing.T) {
	for _, keep := range []bool{false, true} {
		var tags []string
		tester := formTextTester(t, func(c *ui.Context) {
			Input(c, new(string), InputOptions{Label: "Other"})
			TagsInput(c, &tags, TagsInputOptions{Label: "Tags", KeepDraftOnBlur: keep})
		})
		formTextClick(t, tester, "Tags")
		tester.Type("draft")
		formTextClick(t, tester, "Other")
		tester.Frame()
		if got := len(tags) == 1; got == keep {
			t.Fatalf("KeepDraftOnBlur=%v left tags %q", keep, tags)
		}
	}
}

func TestTagsInputPasteSplitsAtCommasTabsAndLines(t *testing.T) {
	var tags []string
	tester := tagsTester(t, &tags, TagsInputOptions{})
	formTextClick(t, tester, "Tags")
	tester.SetClipboard("alpha\nbeta, gamma\tdelta")
	tester.Key(ui.Cmd, ui.KeyV)
	tester.Frame()
	if want := []string{"alpha", "beta", "gamma", "delta"}; !slices.Equal(tags, want) {
		t.Fatalf("tags = %q, want %q", tags, want)
	}
	tester.SetClipboard("plain")
	tester.Key(ui.Cmd, ui.KeyV)
	if len(tags) != 4 {
		t.Fatalf("a paste without separators committed a tag: %q", tags)
	}
}

func TestTagsInputKeyboardWalksTheTokens(t *testing.T) {
	tags := []string{"a", "b", "c"}
	tester := tagsTester(t, &tags, TagsInputOptions{})
	formTextClick(t, tester, "Tags")
	tester.Key(0, ui.KeyLeft)
	tester.Key(0, ui.KeyLeft)
	tester.Key(0, ui.KeyBackspace)
	if want := []string{"a", "c"}; !slices.Equal(tags, want) {
		t.Fatalf("Backspace on the middle token left %q, want %q", tags, want)
	}
	tester.Key(0, ui.KeyRight)
	tester.Key(0, ui.KeyRight)
	tester.Type("z")
	tester.Key(0, ui.KeyEnter)
	if want := []string{"a", "c", "z"}; !slices.Equal(tags, want) {
		t.Fatalf("typing after walking out gave %q, want %q: the focus did not return to the editor", tags, want)
	}
	tester.Key(0, ui.KeyLeft)
	tester.Key(0, ui.KeyDelete)
	tester.Key(0, ui.KeyRight)
	tester.Type("y")
	tester.Key(0, ui.KeyEnter)
	if want := []string{"a", "c", "y"}; !slices.Equal(tags, want) {
		t.Fatalf("Delete on the last token then typing gave %q, want %q", tags, want)
	}
}

func TestTagsInputDeleteWalksDuplicateTokens(t *testing.T) {
	tags := []string{"a", "a", "b"}
	tester := tagsTester(t, &tags, TagsInputOptions{AllowDuplicates: true})
	formTextClick(t, tester, "Tags")
	tester.Key(0, ui.KeyLeft)
	tester.Key(0, ui.KeyLeft)
	tester.Key(0, ui.KeyLeft)
	tester.Key(0, ui.KeyDelete)
	if want := []string{"a", "b"}; !slices.Equal(tags, want) {
		t.Fatalf("Delete on the first duplicate left %q, want %q", tags, want)
	}
}

func TestTagsInputRefusedPiecesStayInTheEditor(t *testing.T) {
	tags := []string{"dup"}
	tester := tagsTester(t, &tags, TagsInputOptions{})
	formTextClick(t, tester, "Tags")
	tester.SetClipboard("dup, b")
	tester.Key(ui.Cmd, ui.KeyV)
	tester.Frame()
	if info, _ := tester.AccessByRole(ui.RoleTextField, "Tags"); !slices.Equal(tags, []string{"dup", "b"}) || info.Value != "dup" {
		t.Fatalf("paste left tags %q and draft %q, want the refused dup kept in the editor", tags, info.Value)
	}
}
