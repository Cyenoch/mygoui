package mygoui_test

import (
	"image"
	"testing"
	"time"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// colorSpan is how far two pixels are apart, summed over the channels.
func colorSpan(a, b image.RGBA, x, y int) int {
	pa, pb := a.RGBAAt(x, y), b.RGBAAt(x, y)
	return absDiff(pa.R, pb.R) + absDiff(pa.G, pb.G) + absDiff(pa.B, pb.B)
}

// TestStateChangesFadeThroughAnIntermediateFrame pins "state changes move":
// when the caller flips the state, the pixel that changes most passes
// through a colour between the two before it settles, motion preferences
// or not, and a control built already in the new state shows it at once and
// asks for no frames.
func TestStateChangesFadeThroughAnIntermediateFrame(t *testing.T) {
	source := redImage(t)
	cases := []struct {
		name  string
		dark  bool
		build func(c *ui.Context, on bool)
	}{
		{"ToggleButton selected", false, func(c *ui.Context, on bool) {
			mygoui.ToggleButton(c, mygoui.ToggleButtonOptions{Label: "Bold", Selected: &on})
		}},
		{"Tag selected", false, func(c *ui.Context, on bool) {
			state := &mygoui.TagGroupState{Selected: map[string]bool{"a": on}}
			mygoui.TagGroup(c, mygoui.TagGroupOptions{Selection: mygoui.TagSelectMultiple, Variant: mygoui.TagSurface, State: state},
				[]mygoui.TagOptions{{Key: "a", Label: "Go"}})
		}},
		{"Chip status", false, func(c *ui.Context, on bool) {
			status := mygoui.ColorDefault
			if on {
				status = mygoui.ColorDanger
			}
			mygoui.Chip(c, mygoui.ChipOptions{Label: "Build", Color: status, Variant: mygoui.ChipPrimary})
		}},
		{"Badge variant", false, func(c *ui.Context, on bool) {
			variant := mygoui.BadgeSecondary
			if on {
				variant = mygoui.BadgeDestructive
			}
			mygoui.Badge(c, mygoui.BadgeOptions{Label: "3", Variant: variant})
		}},
		{"StatusDot status", false, func(c *ui.Context, on bool) {
			status := mygoui.ColorDefault
			if on {
				status = mygoui.ColorDanger
			}
			mygoui.StatusDot(c, mygoui.StatusDotOptions{Status: status})
		}},
		{"Sidebar item selected", false, func(c *ui.Context, on bool) {
			selected, collapsed := "", false
			if on {
				selected = "home"
			}
			mygoui.Sidebar(c, mygoui.SidebarOptions{
				Selected: &selected, Collapsed: &collapsed, Width: 160,
				Groups: []mygoui.SidebarGroupOptions{{Items: []mygoui.SidebarItemOptions{
					{ID: "home", Label: "Home", Icon: lucide.House()},
				}}},
			})
		}},
		{"ProgressBar colour", false, func(c *ui.Context, on bool) {
			color := mygoui.IndicatorSuccess
			if on {
				color = mygoui.IndicatorDanger
			}
			mygoui.ProgressBar(c, mygoui.ProgressOptions{Value: 60, Color: color})
		}},
		{"Statistic delta ink", false, func(c *ui.Context, on bool) {
			status := mygoui.ColorDefault
			if on {
				status = mygoui.ColorDanger
			}
			mygoui.Statistic(c, mygoui.StatisticOptions{Value: "9", Delta: "1", DeltaStatus: status})
		}},
		{"Timeline marker status", false, func(c *ui.Context, on bool) {
			status := mygoui.ColorDefault
			if on {
				status = mygoui.ColorDanger
			}
			mygoui.Timeline(c, mygoui.TimelineOptions{Items: []mygoui.TimelineItem{{Title: "Deploy", Status: status}}})
		}},
		{"Button icon swap", false, func(c *ui.Context, on bool) {
			icon := lucide.Copy()
			if on {
				icon = lucide.Check()
			}
			mygoui.Button(c, mygoui.ButtonOptions{Label: "Copy", Icon: icon, IconOnly: true, Variant: mygoui.ButtonOutline})
		}},
		{"Banner appears", false, func(c *ui.Context, on bool) {
			if on {
				mygoui.Banner(c, mygoui.BannerOptions{Text: "Heads up", Status: mygoui.ColorDanger})
			}
		}},
		{"Alert appears", false, func(c *ui.Context, on bool) {
			if on {
				mygoui.Alert(c, mygoui.AlertOptions{Title: "Heads up", Status: mygoui.ColorDanger})
			}
		}},
		{"Image arrives", false, func(c *ui.Context, on bool) {
			o := mygoui.ImageOptions{Label: "Photo", Width: 64, Height: 64, Status: mygoui.ImageLoading}
			if on {
				o.Source, o.Status = source, mygoui.ImageReady
			}
			mygoui.Image(c, o)
		}},
		{"Avatar image arrives", false, func(c *ui.Context, on bool) {
			o := mygoui.AvatarOptions{Name: "Ada Lovelace", Size: mygoui.AvatarLarge}
			if on {
				o.Image = source
			}
			mygoui.Avatar(c, o)
		}},
	}
	for _, tc := range cases {
		for _, reduced := range []bool{false, true} {
			name := tc.name
			if reduced {
				name += " reduced motion"
			}
			t.Run(name, func(t *testing.T) {
				on := false
				newTester := func(start bool) *ui.Tester {
					on = start
					tt := ui.NewTester(func(c *ui.Context) {
						theme.ApplyMode(c, tc.dark)
						ui.Column(c).Fill().Padding(24).Background(theme.Colors(c).Background).Children(func() {
							tc.build(c, on)
						})
					}, 240, 160)
					tt.SetPreferences(ui.Preferences{ReduceMotion: reduced, TextScale: 1})
					return tt
				}
				snapshot := func(tt *ui.Tester) image.RGBA {
					img := *tt.Image()
					img.Pix = append([]uint8(nil), img.Pix...)
					return img
				}
				tt := newTester(false)
				tt.Advance(time.Second)
				off := snapshot(tt)
				on = true
				tt.Frame()
				tt.Advance(75 * time.Millisecond)
				mid := snapshot(tt)
				tt.Advance(time.Second)
				settled := snapshot(tt)

				bx, by, best := 0, 0, 0
				for y := range settled.Rect.Dy() {
					for x := range settled.Rect.Dx() {
						if d := colorSpan(off, settled, x, y); d > best {
							bx, by, best = x, y, d
						}
					}
				}
				if best < 30 {
					t.Fatalf("the state change paints nothing: largest pixel change %d", best)
				}
				if done := float64(colorSpan(off, mid, bx, by)) / float64(best); done < 0.05 || done > 0.95 {
					t.Fatalf("at 75 ms the pixel (%d, %d) is %.0f%% of the way, want a frame in between", bx, by, done*100)
				}
				if tt.FramePending() {
					t.Fatal("the settled control keeps asking for frames")
				}

				built := newTester(true)
				if d := colorSpan(settled, snapshot(built), bx, by); d > 6 {
					t.Fatalf("built in the new state, the pixel is off by %d on the first frame", d)
				}
				if built.FramePending() {
					t.Fatal("built in the new state, the control animates in")
				}
			})
		}
	}
}

// TestBlocksRevealWithHeightThenSettle pins content opening and closing: the
// parts below a spoiler or an expander row pass through an in
// between position and then settle, and reduced motion settles at once.
func TestBlocksRevealWithHeightThenSettle(t *testing.T) {
	cases := []struct {
		name  string
		build func(c *ui.Context, open *bool)
		// below is the text that sits under the block.
		below string
		// toggle flips the block.
		toggle func(tt *ui.Tester, open *bool)
	}{
		{"Spoiler", func(c *ui.Context, open *bool) {
			mygoui.Spoiler(c, mygoui.SpoilerOptions{Expanded: open, MaxHeight: 40}, func() {
				for i := range 8 {
					ui.Text(c, "line").Key(i)
				}
			})
			ui.Text(c, "Below")
		}, "Below", func(tt *ui.Tester, open *bool) { *open = true }},
		{"Settings expander", func(c *ui.Context, open *bool) {
			mygoui.SettingsGroup(c, mygoui.SettingsGroupOptions{Rows: []mygoui.SettingsRow{{
				Title: "Advanced", Expanded: open,
				Rows: []mygoui.SettingsRow{{Title: "Proxy"}, {Title: "DNS"}},
			}}})
			ui.Text(c, "Below")
		}, "Below", func(tt *ui.Tester, open *bool) { *open = true }},
	}
	for _, tc := range cases {
		for _, reduced := range []bool{false, true} {
			name := tc.name
			if reduced {
				name += " reduced motion"
			}
			t.Run(name, func(t *testing.T) {
				open := false
				tt := ui.NewTester(func(c *ui.Context) {
					theme.ApplyMode(c, false)
					ui.Column(c).Fill().Padding(16).Gap(8).Children(func() { tc.build(c, &open) })
				}, 400, 400)
				tt.SetPreferences(ui.Preferences{ReduceMotion: reduced, TextScale: 1})
				tt.Advance(time.Second)
				from := findRect(t, tt, tc.below).Y
				tc.toggle(tt, &open)
				tt.Frame()
				tt.Advance(100 * time.Millisecond)
				mid := findRect(t, tt, tc.below).Y
				tt.Advance(time.Second)
				to := findRect(t, tt, tc.below).Y
				if to-from < 20 {
					t.Fatalf("opening moved what is below by %v DIPs", to-from)
				}
				if reduced {
					if mid != to {
						t.Fatalf("reduced motion: after 100 ms the block is at %v, want %v at once", mid, to)
					}
				} else if mid <= from+1 || mid >= to-1 {
					t.Fatalf("after 100 ms what is below is at %v, want between %v and %v", mid, from, to)
				}
				if tt.FramePending() {
					t.Fatal("the settled block keeps asking for frames")
				}
			})
		}
	}
}
