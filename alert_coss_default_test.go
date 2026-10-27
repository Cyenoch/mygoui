package mygoui_test

import (
	"bytes"
	"testing"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func TestAlertDefaultDoesNotReserveIconSpace(t *testing.T) {
	for _, status := range []mygoui.StatusColor{
		mygoui.ColorDefault, mygoui.ColorAccent, mygoui.ColorInfo, mygoui.ColorSuccess,
		mygoui.ColorWarning, mygoui.ColorDanger,
	} {
		t.Run(alertStatusName(status), func(t *testing.T) {
			tt := ui.NewTester(func(c *ui.Context) {
				theme.ApplyMode(c, false)
				ui.Column(c).Fill().Padding(20).Children(func() {
					mygoui.Alert(c, mygoui.AlertOptions{
						Title: "Notice", Description: "Description", Status: status,
					}).Label("Callout")
				})
			}, 420, 160)
			callout := alertTestRect(t, tt, "Callout")
			title := alertTestRect(t, tt, "Notice")
			description := alertTestRect(t, tt, "Description")
			left := callout.X + theme.Space14 + theme.Hairline
			if title.X != left || description.X != left {
				t.Fatalf("icon-free text starts at %v and %v, want %v", title.X, description.X, left)
			}
		})
	}
}

func TestAlertKeepsActionWorkingWithAndWithoutIcon(t *testing.T) {
	for _, withIcon := range []bool{false, true} {
		name := "without icon"
		if withIcon {
			name = "explicit icon"
		}
		t.Run(name, func(t *testing.T) {
			var calls int
			tt := ui.NewTester(func(c *ui.Context) {
				theme.ApplyMode(c, false)
				ui.Column(c).Fill().Padding(20).Children(func() {
					var icon *ui.SVG
					if withIcon {
						icon = lucide.Rocket()
					}
					mygoui.Alert(c, mygoui.AlertOptions{
						Title: "Notice", Description: "Description", Status: mygoui.ColorDanger, Icon: icon,
						Action: func() {
							ui.ButtonBase(c).Size(80, 32).Label("Retry").
								OnClick(func() { calls++ })
						},
					}).Label("Callout")
				})
			}, 420, 160)
			callout := alertTestRect(t, tt, "Callout")
			title := alertTestRect(t, tt, "Notice")
			description := alertTestRect(t, tt, "Description")
			left := callout.X + theme.Space14 + theme.Hairline
			if withIcon {
				left += theme.IconMD + theme.Space8
			}
			if title.X != left || description.X != left {
				t.Fatalf("text starts at %v and %v, want %v", title.X, description.X, left)
			}
			if err := tt.Click("Retry"); err != nil {
				t.Fatal(err)
			}
			if calls != 1 || !tt.Focused("Retry") {
				t.Fatalf("action calls=%d focused=%v", calls, tt.Focused("Retry"))
			}
			tt.Key(0, ui.KeyEnter)
			if calls != 2 {
				t.Fatalf("keyboard action calls=%d, want 2", calls)
			}
		})
	}
}

func TestAlertStatusIconOptInAndExplicitIconPrecedence(t *testing.T) {
	image := func(showIcon, explicitIcon bool) (*ui.Tester, []byte) {
		tt := ui.NewTester(func(c *ui.Context) {
			theme.ApplyMode(c, false)
			ui.Column(c).Fill().Padding(20).Children(func() {
				var icon *ui.SVG
				if explicitIcon {
					icon = lucide.Rocket()
				}
				mygoui.Alert(c, mygoui.AlertOptions{
					Title: "Notice", Status: mygoui.ColorInfo,
					ShowIcon: showIcon, Icon: icon,
				}).Label("Callout")
			})
		}, 420, 120)
		return tt, tt.Image().Pix
	}
	optedIn, statusIcon := image(true, false)
	callout := alertTestRect(t, optedIn, "Callout")
	title := alertTestRect(t, optedIn, "Notice")
	left := callout.X + theme.Space14 + theme.Hairline + theme.IconMD + theme.Space8
	if title.X != left {
		t.Fatalf("status icon did not reserve its column: title starts at %v, want %v", title.X, left)
	}
	_, customIcon := image(false, true)
	_, customAndOptedIn := image(true, true)
	if !bytes.Equal(customIcon, customAndOptedIn) {
		t.Fatal("ShowIcon replaced the explicit icon")
	}
	if bytes.Equal(statusIcon, customIcon) {
		t.Fatal("status glyph and explicit icon painted identically")
	}
}

func TestAlertKeepsStatusSurfaceColors(t *testing.T) {
	for _, dark := range []bool{false, true} {
		for _, status := range []mygoui.StatusColor{
			mygoui.ColorDefault, mygoui.ColorAccent, mygoui.ColorInfo,
			mygoui.ColorSuccess, mygoui.ColorWarning, mygoui.ColorDanger,
		} {
			name := "light " + alertStatusName(status)
			if dark {
				name = "dark " + alertStatusName(status)
			}
			t.Run(name, func(t *testing.T) {
				var want ui.Color
				tt := ui.NewTester(func(c *ui.Context) {
					theme.ApplyMode(c, dark)
					palette := theme.Colors(c)
					want = palette.Background
					if status != mygoui.ColorDefault {
						want = mygoui.StatusColors(c, status).Solid.Alpha(0.04).Over(want)
					} else if dark {
						want = palette.Input.Alpha(0.32).Over(want)
					}
					ui.Column(c).Fill().Padding(20).Background(palette.Background).Children(func() {
						mygoui.Alert(c, mygoui.AlertOptions{
							Title: "Notice", Description: "Description", Status: status,
						}).Label("Callout")
					})
				}, 420, 160)
				callout := alertTestRect(t, tt, "Callout")
				got := tt.Image().RGBAAt(int(callout.X+callout.W/2), int(callout.Y+callout.H-5))
				diff := func(a, b uint8) int {
					return max(int(a)-int(b), int(b)-int(a))
				}
				if diff(got.R, want.R) > 2 || diff(got.G, want.G) > 2 || diff(got.B, want.B) > 2 {
					t.Fatalf("surface = %v, want %v", got, want)
				}
			})
		}
	}
}

func alertTestRect(t *testing.T, tt *ui.Tester, label string) ui.Rect {
	t.Helper()
	r, ok := tt.Find(label)
	if !ok {
		t.Fatalf("%q is missing: %v", label, tt.Texts())
	}
	return r
}

func alertStatusName(status mygoui.StatusColor) string {
	switch status {
	case mygoui.ColorAccent:
		return "accent"
	case mygoui.ColorInfo:
		return "info"
	case mygoui.ColorSuccess:
		return "success"
	case mygoui.ColorWarning:
		return "warning"
	case mygoui.ColorDanger:
		return "danger"
	default:
		return "default"
	}
}
