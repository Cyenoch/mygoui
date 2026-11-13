package effect

import (
	"fmt"
	"image/png"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// TestEffectSpecimens writes each effect, on its shader path (interpreted
// on the CPU) and in its fallback, in both appearances, to
// MYGOUI_EFFECT_SPECIMEN_DIR for inspection. It asserts nothing and shows
// no GPU presentation.
func TestEffectSpecimens(t *testing.T) {
	dir := os.Getenv("MYGOUI_EFFECT_SPECIMEN_DIR")
	if dir == "" {
		t.Skip("MYGOUI_EFFECT_SPECIMEN_DIR is not set")
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	save := func(name string, tt *ui.Tester) {
		t.Helper()
		f, err := os.Create(filepath.Join(dir, name+".png"))
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		if err := png.Encode(f, tt.Image()); err != nil {
			t.Fatal(err)
		}
	}
	const scale = 2
	for _, dark := range []bool{false, true} {
		mode := map[bool]string{false: "light", true: "dark"}[dark]
		page := func(w, h int, body func(c *ui.Context)) *ui.Tester {
			tt := ui.NewTester(func(c *ui.Context) {
				theme.ApplyMode(c, dark)
				ui.Column(c).Fill().Padding(16).Background(theme.Colors(c).Background).Children(func() { body(c) })
			}, w, h)
			tt.SetDark(dark)
			tt.SetScale(scale)
			return tt
		}

		for _, v := range []struct {
			name   string
			o      MaskOptions
			forced bool
		}{
			{"shader", MaskOptions{Edges: EdgeBottom, Size: 48, Radius: 10}, true},
			{"fallback-shadow", MaskOptions{Edges: EdgeBottom, Size: 48, Radius: 10}, false},
			{"fallback-surface", MaskOptions{Edges: EdgeBottom, Size: 48, Radius: 10, Surface: theme.ColorsForMode(dark).Background}, false},
			{"shader-sides", MaskOptions{Edges: EdgeLeft | EdgeRight, Radius: 10}, true},
		} {
			tt := page(320, 220, func(c *ui.Context) {
				tk := theme.Colors(c)
				from, to := tk.InfoSoft.Over(tk.Background), tk.SuccessSoft.Over(tk.Background)
				ui.Column(c).FillWidth().Radius(10).Draw(func(p *ui.Painter, r ui.Rect) {
					p.FillGradient(r, ui.LinearGradient{From: from, To: to, Angle: 90}, 10)
				}).Children(func() {
					list := ui.Column(c).FillWidth().Height(188).Padding(12).Gap(8).Radius(10)
					list.Children(func() {
						for i := range 6 {
							ui.Row(c.Key(i)).FillWidth().Padding(8, 12).Radius(8).
								Background(tk.Surface).BorderWidth(1).BorderColor(tk.Border).Children(func() {
								ui.Text(c, fmt.Sprintf("Release 2.%d", 6-i)).TextColor(tk.Foreground)
							})
						}
					})
					Mask(c, list, v.o)
					ui.Local(list, maskKey{}, newMask).forced = v.forced
				})
			})
			tt.Frame()
			save("mask-"+v.name+"-"+mode, tt)
		}

		placeholders := func(c *ui.Context, o ShimmerOptions, forced bool) {
			tk := theme.Colors(c)
			card := ui.Row(c).FillWidth().Padding(12).Gap(12).Radius(10).AlignItems(ui.Start).
				Background(tk.Surface).BorderWidth(1).BorderColor(tk.Border)
			card.Children(func() {
				ui.Box(c).Size(40, 40).Radius(20).Background(tk.MutedBackground)
				ui.Column(c).Grow(1).Gap(8).Children(func() {
					for _, w := range []float32{180, 240, 120} {
						ui.Box(c.Key(w)).Size(w, 12).Radius(6).Background(tk.MutedBackground)
					}
				})
			})
			Shimmer(c, card, o)
			ui.Local(card, shimmerKey{}, newShimmer).forced = forced
		}
		period := theme.Default().Effects.ShimmerPeriod
		for _, v := range []struct {
			name   string
			o      ShimmerOptions
			forced bool
			at     time.Duration
		}{
			{"sweep-early", ShimmerOptions{}, true, period * 4 / 10},
			{"sweep-mid", ShimmerOptions{}, true, period / 2},
			{"pulse-mid", ShimmerOptions{Pulse: true}, true, period / 2},
			{"rest", ShimmerOptions{}, false, period / 2},
		} {
			tt := page(360, 110, func(c *ui.Context) { placeholders(c, v.o, v.forced) })
			into := time.Duration(tt.Now().UnixNano() % int64(period))
			tt.Advance((v.at - into + period) % period)
			tt.Frame()
			save("shimmer-"+v.name+"-"+mode, tt)
		}

		duration := theme.Default().Motion.Reveal
		for _, v := range []struct {
			name   string
			o      RevealOptions
			forced bool
			at     time.Duration
		}{
			{"down-start", RevealOptions{}, true, 0},
			{"down-quarter", RevealOptions{}, true, duration / 4},
			{"down-half", RevealOptions{}, true, duration / 2},
			{"down-done", RevealOptions{}, true, duration},
			{"forward-quarter", RevealOptions{Direction: RevealForward}, true, duration / 4},
			{"fallback-quarter", RevealOptions{}, false, duration / 4},
		} {
			tt := page(320, 170, func(c *ui.Context) {
				tk := theme.Colors(c)
				card := ui.Column(c).FillWidth().Padding(16).Gap(8).Radius(10).
					Background(tk.Surface).BorderWidth(1).BorderColor(tk.Border)
				card.Children(func() {
					ui.Text(c, "Summary").TextColor(tk.Foreground)
					for i, line := range []string{
						"Build finished in 42 seconds.", "128 tests passed, none skipped.",
						"Artifacts were uploaded to the", "release bucket and signed.",
					} {
						ui.Text(c.Key(i), line).TextColor(tk.MutedForeground)
					}
				})
				Reveal(c, card, v.o)
				ui.Local(card, revealKey{}, newReveal).forced = v.forced
			})
			tt.Frame()
			if v.at > 0 {
				tt.Advance(v.at)
			}
			save("reveal-"+v.name+"-"+mode, tt)
		}
	}
}
