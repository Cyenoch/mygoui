package form

import (
	"fmt"
	"testing"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func TestInputGroupTypesAndRunsSuffixAction(t *testing.T) {
	var value string
	cleared := 0
	tester := formTextTester(t, func(c *ui.Context) {
		InputGroup(c, &value, InputGroupOptions{
			Label:       "Amount",
			Placeholder: "0.00",
			Prefix: func() {
				ui.Text(c, "$")
			},
			Suffix: func() {
				ui.ButtonBase(c).Label("Clear").Padding(8).OnClick(func() {
					value = ""
					cleared++
				})
			},
		})
	})
	formTextClick(t, tester, "Amount")
	tester.Type("42")
	if value != "42" {
		t.Fatalf("input group committed %q", value)
	}
	formTextClick(t, tester, "Clear")
	if value != "" || cleared != 1 {
		t.Fatalf("suffix action: value %q, cleared %d", value, cleared)
	}
}

func TestInputGroupBlockPrefixKeepsEditorAsFieldControl(t *testing.T) {
	var value string
	activated := 0
	tester := formTextTester(t, func(c *ui.Context) {
		InputGroup(c, &value, InputGroupOptions{
			Label: "Amount",
			BlockPrefix: func() {
				ui.Text(c, "Block above")
				ui.Button(c, "Add").OnClick(func() { activated++ })
			},
		})
	})
	block, ok := tester.Find("Block above")
	if !ok {
		t.Fatal("input group shows no block prefix")
	}
	// Click the label text at the field's top edge: the label names and
	// focuses the editor, not the Add button that sits above it.
	tester.ClickAt(24, 24)
	caret, focused := tester.TextCaret()
	tester.Type("42")
	if value != "42" || activated != 0 {
		t.Fatalf("label click did not focus the editor: value %q, activated %d", value, activated)
	}
	if !focused || block.Y >= caret.Y {
		t.Fatalf("block prefix is not above the editor: block %v, caret %v", block, caret)
	}
}

func TestInputGroupAddonFocusesEditor(t *testing.T) {
	for _, addon := range []string{"prefix", "suffix", "block prefix", "block suffix"} {
		t.Run(addon, func(t *testing.T) {
			var value string
			tester := formTextTester(t, func(c *ui.Context) {
				o := InputGroupOptions{AriaLabel: "Amount", FullWidth: true}
				content := func() { ui.Text(c, "Addon") }
				switch addon {
				case "prefix":
					o.Prefix = content
				case "suffix":
					o.Suffix = content
				case "block prefix":
					o.BlockPrefix = content
				case "block suffix":
					o.BlockSuffix = content
				}
				InputGroup(c, &value, o)
			})
			formTextClick(t, tester, "Addon")
			tester.Type("42")
			if value != "42" {
				t.Fatalf("typing after addon click = %q", value)
			}
		})
	}
}

func TestInputGroupBlockAddonsAreRows(t *testing.T) {
	var value string
	tester := formTextTester(t, func(c *ui.Context) {
		InputGroup(c, &value, InputGroupOptions{
			AriaLabel: "Editor",
			FullWidth: true,
			BlockPrefix: func() {
				ui.Text(c, "Above one")
				ui.Text(c, "Above two")
			},
			BlockSuffix: func() {
				ui.Text(c, "Below one")
				ui.Text(c, "Below two")
			},
		})
	})
	for _, pair := range [][2]string{{"Above one", "Above two"}, {"Below one", "Below two"}} {
		first, firstOK := tester.Find(pair[0])
		second, secondOK := tester.Find(pair[1])
		if !firstOK || !secondOK {
			t.Fatalf("missing block addon content: %q, %q", pair[0], pair[1])
		}
		if first.X+first.W > second.X+0.5 || first.Y != second.Y {
			t.Errorf("block addon content is not in a row: %q %+v, %q %+v", pair[0], first, pair[1], second)
		}
	}
}

func TestInputGroupInlineAddonContentHasSpace(t *testing.T) {
	var value string
	tester := formTextTester(t, func(c *ui.Context) {
		InputGroup(c, &value, InputGroupOptions{
			AriaLabel: "Editor", FullWidth: true,
			Prefix: func() {
				ui.Text(c, "Prefix one")
				ui.Text(c, "Prefix two")
			},
			Suffix: func() {
				ui.Text(c, "Suffix one")
				ui.Text(c, "Suffix two")
			},
		})
	})
	for _, pair := range [][2]string{{"Prefix one", "Prefix two"}, {"Suffix one", "Suffix two"}} {
		first, firstOK := tester.Find(pair[0])
		second, secondOK := tester.Find(pair[1])
		if !firstOK || !secondOK || second.X-first.X-first.W < theme.Space8-0.5 {
			t.Errorf("addon content has no spacing: %q %+v, %q %+v", pair[0], first, pair[1], second)
		}
	}
}

func TestInputGroupBlockAddonDoesNotForceFullWidth(t *testing.T) {
	var value string
	var normal, block ui.Handle
	var normalBounds, blockBounds ui.Rect
	frames := 0
	tester := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Fill().Padding(16).AlignItems(ui.Start).Children(func() {
			InputGroup(c, &value, InputGroupOptions{AriaLabel: "Plain"}).Bind(&normal)
			InputGroup(c, &value, InputGroupOptions{
				AriaLabel:   "With block",
				BlockPrefix: func() { ui.Text(c, "Preview") },
			}).Bind(&block)
		})
		if frames > 0 {
			normalBounds, blockBounds = normal.Bounds(c), block.Bounds(c)
		}
		frames++
	}, 520, 240)
	tester.Frame()
	if blockBounds.W > normalBounds.W+1 {
		t.Fatalf("block addon forces full width: plain %+v, block %+v", normalBounds, blockBounds)
	}
}

func TestInputGroupEditorFillsSpaceBetweenAddons(t *testing.T) {
	for _, width := range []int{240, 520} {
		t.Run(fmt.Sprintf("width-%d", width), func(t *testing.T) {
			var value string
			tester := ui.NewTester(func(c *ui.Context) {
				ui.Column(c).Fill().Padding(16).Children(func() {
					InputGroup(c, &value, InputGroupOptions{
						AriaLabel: "Editor", FullWidth: true,
						Prefix: func() { ui.Text(c, "Prefix") },
						Suffix: func() { ui.Text(c, "Suffix") },
					})
				})
			}, width, 120)
			prefix, prefixOK := tester.Find("Prefix")
			editor, ok := tester.Find("Editor")
			suffix, suffixOK := tester.Find("Suffix")
			if !prefixOK || !ok || !suffixOK || editor.W <= 0 ||
				prefix.X+prefix.W > editor.X+0.5 || editor.X+editor.W > suffix.X+0.5 {
				t.Fatalf("addons overlap editor: prefix %+v, editor %+v, suffix %+v", prefix, editor, suffix)
			}
			if editor.X+editor.W < suffix.X-0.5 || suffix.X+suffix.W < float32(width-32) {
				t.Fatalf("editor and suffix do not fill the shell: editor %+v, suffix %+v", editor, suffix)
			}
			tester.ClickAt(editor.X+editor.W-2, editor.Y+editor.H/2)
			tester.Type("x")
			if value != "x" {
				t.Fatalf("editor's trailing space does not accept input: %q", value)
			}
		})
	}
}

func TestInputGroupDisabledStopsEditorAndAddonActions(t *testing.T) {
	value := "fixed"
	activated := 0
	tester := formTextTester(t, func(c *ui.Context) {
		InputGroup(c, &value, InputGroupOptions{
			AriaLabel: "Editor", Disabled: true, FullWidth: true,
			Prefix: func() { ui.Text(c, "Prefix") },
			Suffix: func() { ui.Button(c, "Action").OnClick(func() { activated++ }) },
		})
	})
	formTextClick(t, tester, "Prefix")
	tester.Type("x")
	formTextClick(t, tester, "Action")
	tester.Key(0, ui.KeySpace)
	if value != "fixed" || activated != 0 || tester.Focused("Editor") || tester.Focused("Action") {
		t.Fatalf("disabled group accepted input: value %q, actions %d, editor focused %v, action focused %v",
			value, activated, tester.Focused("Editor"), tester.Focused("Action"))
	}
}

func TestInputGroupAddonActionKeepsItsFocus(t *testing.T) {
	var value string
	activated := 0
	var shell ui.Rect
	tester := formTextTester(t, func(c *ui.Context) {
		e := InputGroup(c, &value, InputGroupOptions{
			AriaLabel: "Editor",
			FullWidth: true,
			Suffix: func() {
				ui.Button(c, "Action").OnClick(func() { activated++ })
			},
		})
		shell = e.LayoutBounds()
	})
	tester.Frame()
	img := tester.Image()
	x, y := int(shell.X)-2, int(shell.Y+shell.H/2)
	face := img.RGBAAt(x, y)
	formTextClick(t, tester, "Action")
	if activated != 1 || !tester.Focused("Action") {
		t.Fatalf("addon action lost activation or focus: activated %d, action focused %v, editor focused %v",
			activated, tester.Focused("Action"), tester.Focused("Editor"))
	}
	if got := tester.Image().RGBAAt(x, y); got != face {
		t.Fatalf("addon button focus painted the editor shell ring: %v -> %v", face, got)
	}
	formTextClick(t, tester, "Editor")
	if got := tester.Image().RGBAAt(x, y); got == face {
		t.Fatalf("focused editor is missing its shared shell ring: shell=%+v sample=(%d,%d) focused=%v pixel=%v", shell, x, y, tester.Focused("Editor"), got)
	}
}

func TestInputGroupErrorPaintsInvalidShell(t *testing.T) {
	var errorValue, invalidValue string
	var errorEditor, invalidEditor ui.Rect
	var colors theme.Palette
	tester := formTextTester(t, func(c *ui.Context) {
		colors = theme.Colors(c)
		InputGroup(c, &errorValue, InputGroupOptions{AriaLabel: "Error editor", Error: "Required", FullWidth: true})
		InputGroup(c, &invalidValue, InputGroupOptions{AriaLabel: "Invalid editor", Invalid: true, FullWidth: true})
	})
	var ok bool
	if errorEditor, ok = tester.Find("Error editor"); !ok {
		t.Fatal("error editor missing")
	}
	if invalidEditor, ok = tester.Find("Invalid editor"); !ok {
		t.Fatal("invalid editor missing")
	}
	img := tester.Image()
	wantFormHoverPixel(t, img, int(invalidEditor.X), int(invalidEditor.Y+invalidEditor.H/2),
		colors.Destructive.Alpha(0.36).Over(colors.Background), "invalid border")
	errorBorder := img.RGBAAt(int(errorEditor.X), int(errorEditor.Y+errorEditor.H/2))
	invalidBorder := img.RGBAAt(int(invalidEditor.X), int(invalidEditor.Y+invalidEditor.H/2))
	if errorBorder != invalidBorder {
		t.Fatalf("error without Invalid shows a different border: error %v, invalid %v", errorBorder, invalidBorder)
	}
}
