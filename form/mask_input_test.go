package form

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

// maskTester builds a focused masked input over value.
func maskTester(t *testing.T, mask string, value *string, o MaskInputOptions) *ui.Tester {
	t.Helper()
	o.Mask, o.Label = mask, "Masked"
	tester := formTextTester(t, func(c *ui.Context) { MaskInput(c, value, o) })
	formTextClick(t, tester, "Masked")
	return tester
}

// maskShown is the text the input shows.
func maskShown(t *testing.T, tester *ui.Tester) string {
	t.Helper()
	info, ok := tester.AccessByRole(ui.RoleTextField, "Masked")
	if !ok {
		t.Fatal("the masked input is not rendered")
	}
	return info.Value
}

func maskKeys(tester *ui.Tester, key ui.Key, n int) {
	for range n {
		tester.Key(0, key)
	}
}

func TestMaskInputWritesLiteralsAndStoresRawValue(t *testing.T) {
	var value string
	changes := 0
	tester := maskTester(t, "99/99/9999", &value, MaskInputOptions{OnChange: func() { changes++ }})
	for _, ch := range "12x251999777" {
		tester.Type(string(ch))
	}
	if value != "12251999" || maskShown(t, tester) != "12/25/1999" {
		t.Fatalf("value %q shown %q, want 12251999 shown 12/25/1999", value, maskShown(t, tester))
	}
	if changes != 8 {
		t.Fatalf("%d changes for eight accepted digits", changes)
	}
}

func TestMaskInputClassesAndPaste(t *testing.T) {
	var value string
	tester := maskTester(t, "AA-9999", &value, MaskInputOptions{})
	tester.Type("1ab-12-345")
	if value != "ab1234" || maskShown(t, tester) != "ab-1234" {
		t.Fatalf("value %q shown %q", value, maskShown(t, tester))
	}
}

func TestMaskInputInsertsInTheMiddleAndKeepsCaret(t *testing.T) {
	value := "1234"
	tester := maskTester(t, "99/99/9999", &value, MaskInputOptions{})
	tester.Key(0, ui.KeyHome)
	tester.Key(0, ui.KeyRight)
	tester.Type("9")
	if value != "19234" || maskShown(t, tester) != "19/23/4" {
		t.Fatalf("value %q shown %q, want 19234 shown 19/23/4", value, maskShown(t, tester))
	}
	// The caret followed the typed digit across the literal.
	tester.Type("8")
	if value != "198234" || maskShown(t, tester) != "19/82/34" {
		t.Fatalf("value %q shown %q, want 198234 shown 19/82/34", value, maskShown(t, tester))
	}
}

func TestMaskInputDeletesAcrossLiterals(t *testing.T) {
	value := "123456"
	tester := maskTester(t, "99/99/9999", &value, MaskInputOptions{})
	// 12/34/56: put the caret after the first slash.
	tester.Key(0, ui.KeyHome)
	maskKeys(tester, ui.KeyRight, 3)
	tester.Key(0, ui.KeyBackspace)
	if value != "13456" || maskShown(t, tester) != "13/45/6" {
		t.Fatalf("Backspace after a literal gave %q shown %q, want 13456 shown 13/45/6", value, maskShown(t, tester))
	}
	// The caret now sits after the 1; Delete removes the 3 and the next
	// Delete, before the slash, removes the digit beyond it.
	tester.Key(0, ui.KeyDelete)
	if value != "1456" || maskShown(t, tester) != "14/56" {
		t.Fatalf("Delete gave %q shown %q, want 1456 shown 14/56", value, maskShown(t, tester))
	}
	tester.Key(0, ui.KeyRight)
	tester.Key(0, ui.KeyDelete)
	if value != "146" || maskShown(t, tester) != "14/6" {
		t.Fatalf("Delete before a literal gave %q shown %q, want 146 shown 14/6", value, maskShown(t, tester))
	}
}

func TestMaskInputTrailingLiteralClosesAndReopens(t *testing.T) {
	var value string
	tester := maskTester(t, "(999)", &value, MaskInputOptions{})
	tester.Type("123")
	if maskShown(t, tester) != "(123)" {
		t.Fatalf("a full mask shows %q, want (123)", maskShown(t, tester))
	}
	tester.Key(0, ui.KeyBackspace)
	if value != "12" || maskShown(t, tester) != "(12" {
		t.Fatalf("Backspace at the end gave %q shown %q, want 12 shown (12", value, maskShown(t, tester))
	}
}

func TestMaskInputFollowsCallerValue(t *testing.T) {
	value := "0101"
	tester := maskTester(t, "99/99/9999", &value, MaskInputOptions{})
	value = "31122024"
	tester.Frame()
	if maskShown(t, tester) != "31/12/2024" {
		t.Fatalf("shown %q after the caller set the value", maskShown(t, tester))
	}
	if got := MaskFormat("99/99/9999", "3112"); got != "31/12" {
		t.Fatalf("MaskFormat = %q, want 31/12", got)
	}
}

// Equal neighbours give a text comparison nothing to place the edit by.
func TestMaskInputPlacesEditsAmongEqualCharacters(t *testing.T) {
	value := "12"
	tester := maskTester(t, "9999", &value, MaskInputOptions{})
	tester.Key(0, ui.KeyHome)
	tester.Key(0, ui.KeyRight)
	tester.Type("2")
	tester.Type("5")
	if value != "1252" {
		t.Fatalf("typing 2 then 5 after the 1 of 12 gave %q, want 1252", value)
	}
	value = "1234"
	tester = maskTester(t, "99--99", &value, MaskInputOptions{})
	tester.Key(0, ui.KeyHome)
	maskKeys(tester, ui.KeyRight, 3)
	tester.Key(0, ui.KeyBackspace)
	if value != "134" {
		t.Fatalf("Backspace between two equal literals gave %q, want 134", value)
	}
}

func TestMaskInputRefusesRemovalThatWouldDropLaterCharacters(t *testing.T) {
	value := "ab12"
	changes := 0
	tester := maskTester(t, "AA-99", &value, MaskInputOptions{OnChange: func() { changes++ }})
	tester.Key(0, ui.KeyHome)
	tester.Key(0, ui.KeyRight)
	tester.Key(0, ui.KeyBackspace)
	if value != "ab12" || maskShown(t, tester) != "ab-12" || changes != 0 {
		t.Fatalf("value %q shown %q after %d changes, want ab12 kept", value, maskShown(t, tester), changes)
	}
	// The caret stayed after the a: removal from the end still works.
	tester.Key(0, ui.KeyEnd)
	tester.Key(0, ui.KeyBackspace)
	if value != "ab1" {
		t.Fatalf("Backspace at the end gave %q, want ab1", value)
	}
}

func TestMaskInputUndoAndRedoStepThroughRawValues(t *testing.T) {
	var value string
	changes := 0
	tester := maskTester(t, "99/99", &value, MaskInputOptions{OnChange: func() { changes++ }})
	tester.Type("1")
	tester.Type("2")
	tester.Type("3")
	tester.Key(0, ui.KeyBackspace)
	if value != "12" {
		t.Fatalf("value %q after typing 123 and Backspace", value)
	}
	tester.Key(ui.Cmd, ui.KeyZ)
	if value != "123" || maskShown(t, tester) != "12/3" {
		t.Fatalf("the first undo gave %q shown %q, want the deleted digit back", value, maskShown(t, tester))
	}
	tester.Key(ui.Cmd, ui.KeyZ)
	if value != "" || maskShown(t, tester) != "" {
		t.Fatalf("the second undo gave %q shown %q, want the whole typed run undone", value, maskShown(t, tester))
	}
	tester.Key(ui.Cmd|ui.Shift, ui.KeyZ)
	if value != "123" || maskShown(t, tester) != "12/3" {
		t.Fatalf("redo gave %q shown %q", value, maskShown(t, tester))
	}
	if changes != 7 {
		t.Fatalf("%d changes for four edits, two undos and a redo, want 7", changes)
	}
	tester.Type("4")
	tester.Key(ui.Cmd|ui.Shift, ui.KeyZ)
	if value != "1234" {
		t.Fatalf("redo after a new edit changed the value to %q", value)
	}
}

func TestMaskInputPasteKeepsTheLastDigitsOfALongNumber(t *testing.T) {
	var value string
	tester := maskTester(t, "(999) 999-9999", &value, MaskInputOptions{})
	tester.SetClipboard("+1 (415) 555-0100")
	tester.Key(ui.Cmd, ui.KeyV)
	tester.Frame()
	if value != "4155550100" || maskShown(t, tester) != "(415) 555-0100" {
		t.Fatalf("pasted number gave %q shown %q, want 4155550100", value, maskShown(t, tester))
	}
	tester.SetClipboard("(212) 555-0199")
	tester.Key(ui.Cmd, ui.KeyA)
	tester.Key(ui.Cmd, ui.KeyV)
	tester.Frame()
	if value != "2125550199" {
		t.Fatalf("a paste over a selection gave %q", value)
	}
}

func TestMaskInputTypingReplacesTheSelection(t *testing.T) {
	value := "123456"
	tester := maskTester(t, "99/99/9999", &value, MaskInputOptions{})
	tester.Key(ui.Cmd, ui.KeyA)
	tester.Type("7")
	if value != "7" || maskShown(t, tester) != "7" {
		t.Fatalf("typing over a full selection gave %q shown %q", value, maskShown(t, tester))
	}
}

func TestMaskInputKeepsCallerValueOnRejectedKey(t *testing.T) {
	value := "12ab"
	changes := 0
	tester := maskTester(t, "9999", &value, MaskInputOptions{OnChange: func() { changes++ }})
	tester.Key(0, ui.KeyEnd)
	tester.Type("x")
	if value != "12ab" || changes != 0 {
		t.Fatalf("a rejected key left %q after %d changes", value, changes)
	}
}

func TestMaskInputLongPasteReplacesAPartlyFilledValue(t *testing.T) {
	var value string
	tester := maskTester(t, "(999) 999-9999", &value, MaskInputOptions{})
	tester.Type("12")
	tester.SetClipboard("+1 (555) 123-4567")
	tester.Key(ui.Cmd, ui.KeyV)
	tester.Frame()
	if value != "5551234567" {
		t.Fatalf("a full number pasted after 12 gave %q, want 5551234567", value)
	}
}
