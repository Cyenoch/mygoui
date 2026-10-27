package form

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestComposerSubmitsCallerDraftAndKeepsNewlines(t *testing.T) {
	var draft string
	var sent []string
	tt := formTextTester(t, func(c *ui.Context) {
		Composer(c, &draft, ComposerOptions{
			Label: "Message", SubmitLabel: "Send",
			OnSubmit: func(text string) { sent = append(sent, text); draft = "" },
		})
	})
	formTextClick(t, tt, "Message")
	tt.Type("hello")
	tt.Key(ui.Shift, ui.KeyEnter)
	tt.Type("world")
	tt.Key(0, ui.KeyEnter)
	if len(sent) != 1 || sent[0] != "hello\nworld" || draft != "" {
		t.Fatalf("submitted %q, draft %q", sent, draft)
	}
	tt.Type("again")
	if draft != "again" {
		t.Fatalf("typing after submit = %q", draft)
	}
	if err := tt.Click("Send"); err != nil {
		t.Fatal(err)
	}
	if len(sent) != 2 || sent[1] != "again" {
		t.Fatalf("button submitted %q", sent)
	}
}

func TestComposerDoesNotSubmitDisabledOrEmptyDraft(t *testing.T) {
	var draft string
	var calls int
	tt := formTextTester(t, func(c *ui.Context) {
		Composer(c, &draft, ComposerOptions{
			Label: "Message", OnSubmit: func(string) { calls++ },
		})
	})
	formTextClick(t, tt, "Message")
	tt.Key(0, ui.KeyEnter)
	if calls != 0 || draft != "" {
		t.Fatalf("empty submit: calls=%d draft=%q", calls, draft)
	}
}

func TestComposerReadOnlyPreservesDraft(t *testing.T) {
	draft := "existing"
	var changes, submissions int
	tt := formTextTester(t, func(c *ui.Context) {
		Composer(c, &draft, ComposerOptions{
			Label: "Message", ReadOnly: true,
			OnChange: func() { changes++ }, OnSubmit: func(string) { submissions++ },
		})
	})
	formTextClick(t, tt, "Message")
	tt.Type("new")
	tt.Key(0, ui.KeyEnter)
	if draft != "existing" || changes != 0 || submissions != 0 {
		t.Fatalf("read-only: draft %q changes %d submissions %d", draft, changes, submissions)
	}
}

func TestComposerDisabledDoesNotAcceptInput(t *testing.T) {
	draft := "existing"
	var changes, submissions int
	tt := formTextTester(t, func(c *ui.Context) {
		Composer(c, &draft, ComposerOptions{
			Label: "Message", Disabled: true,
			OnChange: func() { changes++ }, OnSubmit: func(string) { submissions++ },
		})
	})
	formTextClick(t, tt, "Message")
	tt.Type("new")
	tt.Key(0, ui.KeyEnter)
	if draft != "existing" || changes != 0 || submissions != 0 {
		t.Fatalf("disabled: draft %q changes %d submissions %d", draft, changes, submissions)
	}
}

func TestComposerSubmitButtonPreservesEditorFocus(t *testing.T) {
	var draft string
	var sent []string
	tt := formTextTester(t, func(c *ui.Context) {
		Composer(c, &draft, ComposerOptions{
			Label: "Message", SubmitLabel: "Send",
			OnSubmit: func(text string) { sent = append(sent, text); draft = "" },
		})
	})
	formTextClick(t, tt, "Message")
	tt.Type("hello")
	formTextClick(t, tt, "Send")
	tt.Type("again")
	if len(sent) != 1 || sent[0] != "hello" || draft != "again" {
		t.Fatalf("sent %q, typing after submit = %q", sent, draft)
	}
}
