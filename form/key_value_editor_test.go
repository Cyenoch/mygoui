package form

import (
	"slices"
	"testing"

	"github.com/Cyenoch/mygoui/locale"
	"github.com/egoist/mygo/ui"
)

// keyValueFixture applies the editor's reports the way a caller does.
type keyValueFixture struct {
	pairs []KeyValue
	o     KeyValueEditorOptions
}

func (f *keyValueFixture) tester(t *testing.T) *ui.Tester {
	t.Helper()
	return formTextTester(t, func(c *ui.Context) {
		o := f.o
		o.Pairs = f.pairs
		o.OnChange = func(i int, p KeyValue) { f.pairs[i] = p }
		o.OnAdd = func() { f.pairs = append(f.pairs, KeyValue{}) }
		o.OnRemove = func(i int) { f.pairs = slices.Delete(f.pairs, i, i+1) }
		KeyValueEditor(c, o)
	})
}

func TestKeyValueEditorReportsEditsByRow(t *testing.T) {
	f := &keyValueFixture{pairs: []KeyValue{{Key: "HOME", Value: "/root"}, {Key: "SHELL", Value: ""}}}
	tester := f.tester(t)
	formTextClick(t, tester, "Value 2")
	tester.Type("/bin/sh")
	want := []KeyValue{{Key: "HOME", Value: "/root"}, {Key: "SHELL", Value: "/bin/sh"}}
	if !slices.Equal(f.pairs, want) {
		t.Fatalf("pairs = %+v, want %+v", f.pairs, want)
	}
}

func TestKeyValueEditorNeverWritesCallerPairs(t *testing.T) {
	pairs := []KeyValue{{Key: "HOME", Value: "/root"}}
	var reported []KeyValue
	tester := formTextTester(t, func(c *ui.Context) {
		KeyValueEditor(c, KeyValueEditorOptions{
			Pairs:    pairs,
			OnChange: func(_ int, p KeyValue) { reported = append(reported, p) },
		})
	})
	formTextClick(t, tester, "Key 1")
	tester.Key(0, ui.KeyEnd)
	tester.Type("S")
	if pairs[0] != (KeyValue{Key: "HOME", Value: "/root"}) {
		t.Fatalf("the editor wrote the caller's pair: %+v", pairs[0])
	}
	if len(reported) == 0 || reported[0] != (KeyValue{Key: "HOMES", Value: "/root"}) {
		t.Fatalf("reports = %+v", reported)
	}
	// An edit the caller declines is dropped, not kept in the editor.
	tester.Frame()
	if info, _ := tester.AccessByRole(ui.RoleTextField, "Key 1"); info.Value != "HOME" {
		t.Fatalf("the editor shows %q, want the caller's HOME", info.Value)
	}
}

func TestKeyValueEditorAddFocusesNewRowAndRemoveDropsItsRow(t *testing.T) {
	f := &keyValueFixture{pairs: []KeyValue{{Key: "A", Value: "1"}, {Key: "B", Value: "2"}}}
	tester := f.tester(t)
	formTextClick(t, tester, locale.M().Add)
	tester.Frame()
	tester.Type("C")
	if len(f.pairs) != 3 || f.pairs[2].Key != "C" {
		t.Fatalf("typing after Add gave %+v, want a third row keyed C", f.pairs)
	}
	formTextClick(t, tester, locale.M().Remove+" A")
	want := []KeyValue{{Key: "B", Value: "2"}, {Key: "C", Value: ""}}
	if !slices.Equal(f.pairs, want) {
		t.Fatalf("pairs = %+v, want %+v", f.pairs, want)
	}
	if info, _ := tester.AccessByRole(ui.RoleTextField, "Key 1"); info.Value != "B" {
		t.Fatalf("the first row shows %q after the removal, want B", info.Value)
	}
}

func TestKeyValueEditorMarksDuplicateKeys(t *testing.T) {
	f := &keyValueFixture{
		pairs: []KeyValue{{Key: "PATH", Value: "/bin"}, {Key: "HOME", Value: "/root"}, {Key: "PATH", Value: "/usr/bin"}, {Key: "", Value: ""}, {Key: "", Value: ""}},
		o:     KeyValueEditorOptions{MarkDuplicateKeys: true},
	}
	tester := f.tester(t)
	for i, want := range []bool{true, false, true, false, false} {
		name := locale.M().Key + " " + string(rune('1'+i))
		info, ok := tester.AccessByRole(ui.RoleTextField, name)
		if !ok || info.Invalid != want {
			t.Fatalf("%s: found %v, invalid %v, want %v", name, ok, info.Invalid, want)
		}
	}
}

func TestKeyValueEditorDuplicateKeyShowsMessage(t *testing.T) {
	f := &keyValueFixture{
		pairs: []KeyValue{{Key: "A"}, {Key: "A"}, {Key: "B"}},
		o:     KeyValueEditorOptions{MarkDuplicateKeys: true},
	}
	tester := f.tester(t)
	if !tester.HasText(locale.M().DuplicateKey) {
		t.Fatal("a duplicate key shows no message")
	}
	f.pairs[1].Key = "C"
	tester.Frame()
	if tester.HasText(locale.M().DuplicateKey) {
		t.Fatal("the message stays after the duplicate was resolved")
	}
}

func TestKeyValueEditorPasteReportsEnvLines(t *testing.T) {
	var got []KeyValue
	f := &keyValueFixture{pairs: []KeyValue{{}}}
	f.o.OnPaste = func(p []KeyValue) { got = p }
	tester := f.tester(t)
	formTextClick(t, tester, "Key 1")
	tester.SetClipboard("# db\nexport HOST=localhost\n\nPASS=\"s=3\"\n")
	tester.Key(ui.Cmd, ui.KeyV)
	tester.Frame()
	want := []KeyValue{{Key: "HOST", Value: "localhost"}, {Key: "PASS", Value: "s=3"}}
	if !slices.Equal(got, want) || f.pairs[0].Key != "" {
		t.Fatalf("paste reported %+v (row holds %+v), want %+v and nothing typed", got, f.pairs[0], want)
	}
	// A single value with an equals sign is a value, not a row.
	got = nil
	formTextClick(t, tester, "Value 1")
	tester.SetClipboard("abc==")
	tester.Key(ui.Cmd, ui.KeyV)
	tester.Frame()
	if got != nil || f.pairs[0].Value != "abc==" {
		t.Fatalf("a value paste was taken as rows: %+v, value %q", got, f.pairs[0].Value)
	}
}

func TestKeyValueEditorSecretConcealsUntilRevealed(t *testing.T) {
	f := &keyValueFixture{pairs: []KeyValue{{Key: "TOKEN", Value: "hunter2", Secret: true}}}
	tester := f.tester(t)
	if tester.HasText("hunter2") {
		t.Fatal("a secret value is readable")
	}
	formTextClick(t, tester, locale.M().ShowPassword+" TOKEN")
	if info, _ := tester.AccessByRole(ui.RoleTextField, "Value 1"); info.Value != "hunter2" {
		t.Fatalf("revealed value reads %q", info.Value)
	}
}

func TestKeyValueEditorEnterMovesOnAndAddsRow(t *testing.T) {
	f := &keyValueFixture{pairs: []KeyValue{{Key: "A"}}}
	tester := f.tester(t)
	formTextClick(t, tester, "Key 1")
	tester.Key(0, ui.KeyEnter)
	tester.Type("1")
	if f.pairs[0].Value != "1" {
		t.Fatalf("Enter in a key left the value %q, want the typing in the value", f.pairs[0].Value)
	}
	tester.Key(0, ui.KeyEnter)
	tester.Frame()
	tester.Type("B")
	if len(f.pairs) != 2 || f.pairs[1].Key != "B" {
		t.Fatalf("Enter in the last value gave %+v, want a second row keyed B", f.pairs)
	}
}

func TestKeyValueEditorRowIDKeepsEditingStateAcrossRemoval(t *testing.T) {
	f := &keyValueFixture{pairs: []KeyValue{{ID: "a", Key: "A"}, {ID: "b", Key: "B", Value: "x"}}}
	tester := f.tester(t)
	formTextClick(t, tester, "Value 2")
	f.pairs = f.pairs[1:]
	tester.Frame()
	tester.Type("y")
	if f.pairs[0].Value != "xy" {
		t.Fatalf("typing after row A went gave %q, want the focus kept in B's value", f.pairs[0].Value)
	}
}

func TestKeyValueEditorRemovingARowKeepsLaterSecretsConcealed(t *testing.T) {
	f := &keyValueFixture{pairs: []KeyValue{{Key: "A", Value: "one", Secret: true}, {Key: "B", Value: "two", Secret: true}}}
	tester := f.tester(t)
	formTextClick(t, tester, locale.M().ShowPassword+" A")
	formTextClick(t, tester, locale.M().Remove+" A")
	tester.Frame()
	if info, _ := tester.AccessByRole(ui.RoleTextField, "Value 1"); info.Value == "two" {
		t.Fatal("removing a revealed row revealed the secret that moved into its place")
	}
}

func TestKeyValueEditorRemoveMovesFocusToNeighbour(t *testing.T) {
	f := &keyValueFixture{pairs: []KeyValue{{Key: "A"}, {Key: "B"}, {Key: "C"}}}
	tt := f.tester(t)
	remove := func(key string) string { return locale.M().Remove + " " + key }
	for _, step := range []struct{ remove, focus string }{
		{"A", remove("B")}, {"C", remove("B")}, {"B", locale.M().Add},
	} {
		formTextClick(t, tt, remove(step.remove))
		tt.Frame()
		if !tt.Focused(step.focus) {
			t.Fatalf("after removing %s the focus is not on %q", step.remove, step.focus)
		}
	}
}
