package mygoui_test

import (
	"math"
	"testing"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/overlay"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func TestIconButtonCentersGlyphWithoutStretching(t *testing.T) {
	for _, size := range []struct {
		name  string
		value mygoui.ButtonSize
	}{
		{name: "extra small", value: mygoui.ButtonIconXS}, {name: "small", value: mygoui.ButtonIconSM},
		{name: "default", value: mygoui.ButtonIcon}, {name: "large", value: mygoui.ButtonIconLG},
	} {
		for _, scale := range []struct {
			name  string
			value float32
		}{{name: "normal", value: 1}, {name: "large text", value: 2}} {
			t.Run(size.name+"/"+scale.name, func(t *testing.T) {
				tt := ui.NewTester(func(c *ui.Context) {
					theme.Apply(c)
					ui.Column(c).Fill().Padding(30).Children(func() {
						mygoui.Button(c, mygoui.ButtonOptions{Label: "Close sample", Icon: lucide.X(),
							Size: size.value, Variant: mygoui.ButtonSecondary, OnClick: func() {}})
					})
				}, 360, 180)
				tt.SetPreferences(ui.Preferences{TextScale: scale.value, ReduceMotion: true})
				tt.Frame()
				r, ok := tt.Find("Close sample")
				if !ok {
					t.Fatal("missing icon action")
				}
				if math.Abs(float64(r.W-r.H)) > 0.5 {
					t.Fatalf("parent stretched icon target: %+v", r)
				}
				assertIconCentered(t, tt, r)
			})
		}
	}
}

func TestModalCloseUsesCenteredIconRatherThanFontBaseline(t *testing.T) {
	open := true
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		overlay.Dialog(c, &open, overlay.DialogOptions{Title: "Edit profile"}, func() { ui.Text(c, "Readable body") })
	}, 600, 400)
	tt.SetPreferences(ui.Preferences{TextScale: 2, ReduceMotion: true})
	tt.Frame()
	r, ok := tt.Find("Close")
	if !ok {
		t.Fatal("missing modal close")
	}
	assertIconCentered(t, tt, r)
	if err := tt.Click("Close"); err != nil {
		t.Fatal(err)
	}
	if open {
		t.Fatal("centered close did not dismiss")
	}
}

func assertIconCentered(t *testing.T, tt *ui.Tester, r ui.Rect) {
	t.Helper()
	left, right, top, bottom := int(r.X+r.W), int(r.X), int(r.Y+r.H), int(r.Y)
	for y := int(r.Y + 3); y < int(r.Y+r.H-3); y++ {
		for x := int(r.X + 3); x < int(r.X+r.W-3); x++ {
			p := tt.Image().RGBAAt(x, y)
			if p.R > 220 && p.G > 220 && p.B > 220 {
				left, right, top, bottom = min(left, x), max(right, x), min(top, y), max(bottom, y)
			}
		}
	}
	if left > right {
		t.Fatal("icon glyph was not painted")
	}
	x, y := float64(left+right+1)/2, float64(top+bottom+1)/2
	if math.Abs(x-float64(r.X+r.W/2)) > 1 || math.Abs(y-float64(r.Y+r.H/2)) > 1 {
		t.Fatalf("glyph center=(%.1f,%.1f), target center=(%.1f,%.1f), glyph box=(%d,%d)-(%d,%d)", x, y, r.X+r.W/2, r.Y+r.H/2, left, top, right, bottom)
	}
}
