package conversation

import (
	"fmt"
	"math"
	"testing"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func TestMessageMeasuredFooterPreservesAvatarAndBubbleAlignment(t *testing.T) {
	for _, width := range []int{390, 800} {
		for _, scale := range []float32{1, 2} {
			for _, direction := range []ui.Direction{ui.LeftToRight, ui.RightToLeft} {
				for _, align := range []Align{Start, End} {
					name := fmt.Sprintf("width=%d/scale=%g/direction=%d/align=%d", width, scale, direction, align)
					t.Run(name, func(t *testing.T) {
						footer := "Delivered"
						footerActions := 0
						tt := ui.NewTester(func(c *ui.Context) {
							theme.Apply(c)
							ui.Column(c).FillWidth().Padding(16).Children(func() {
								Message(c, MessageOptions{Align: align, Label: "Message", Header: "Oliver", Footer: "Replaced footer",
									Avatar: func() { ui.Box(c).Size(32, 32).Label("Oliver avatar") },
									FooterContent: func() {
										ui.Text(c, footer).Label("Delivery details")
										ui.ButtonBase(c).Label("Delivery action").OnClick(func() { footerActions++ }).Children(func() {
											ui.Text(c, "Details")
										})
									},
								}, func() {
									MessageGroup(c, func() {
										Bubble(c, BubbleOptions{Align: align, Variant: BubbleSecondary, Label: "First bubble"}, func() {
											ui.Text(c, "This longer first bubble wraps when the conversation is narrow.")
										})
										Bubble(c, BubbleOptions{Align: align, Variant: BubbleSecondary, Label: "Last bubble"}, func() {
											ui.Text(c, "A short reply")
										})
									})
								})
							})
						}, width, 1200)
						tt.SetPreferences(ui.Preferences{TextScale: scale, ReduceMotion: true})
						tt.SetDirection(direction)
						first := messageRect(t, tt, "First bubble")
						last := messageRect(t, tt, "Last bubble")
						avatar := messageRect(t, tt, "Oliver avatar")
						message := messageRect(t, tt, "Message")
						oldFooter := messageRect(t, tt, "Delivery details")
						assertMessageNear(t, avatar.Y+avatar.H, last.Y+last.H, "avatar bottom and final bubble")
						if first.Y+first.H > last.Y || first.W <= last.W {
							t.Fatalf("consecutive bubbles must stack and retain variable width: first=%+v last=%+v", first, last)
						}
						for _, r := range []ui.Rect{first, last, avatar, oldFooter} {
							if r.X < 15 || r.X+r.W > float32(width)-15 {
								t.Fatalf("message escaped narrow viewport: %+v at width %d", r, width)
							}
						}
						atRight := (align == End) != (direction == ui.RightToLeft)
						if atRight {
							if avatar.X < last.X+last.W {
								t.Fatalf("trailing avatar crossed content: avatar=%+v bubble=%+v", avatar, last)
							}
							assertMessageNear(t, first.X+first.W, last.X+last.W, "trailing grouped bubble edge")
						} else {
							if avatar.X+avatar.W > last.X {
								t.Fatalf("leading avatar crossed content: avatar=%+v bubble=%+v", avatar, last)
							}
							assertMessageNear(t, first.X, last.X, "leading grouped bubble edge")
						}
						footer = "Delivered to Oliver\nRead on another device\nAdditional delivery information remains readable when it wraps across several lines."
						tt.Frame()
						if got := messageRect(t, tt, "Oliver avatar"); got != avatar {
							t.Fatalf("footer height moved avatar: before=%+v after=%+v", avatar, got)
						}
						if got := messageRect(t, tt, "Last bubble"); got != last {
							t.Fatalf("footer height moved bubble: before=%+v after=%+v", last, got)
						}
						newFooter := messageRect(t, tt, "Delivery details")
						if newFooter.H <= oldFooter.H || messageRect(t, tt, "Message").H <= message.H {
							t.Fatal("multiline footer did not grow the measured message")
						}
						if newFooter.Y < last.Y+last.H || tt.HasText("Replaced footer") {
							t.Fatal("rich footer overlapped content or failed to replace string metadata")
						}
						if err := tt.Click("Delivery action"); err != nil {
							t.Fatal(err)
						}
						if footerActions != 1 {
							t.Fatal("rich footer action did not retain native input")
						}
					})
				}
			}
		}
	}
}

func TestMessageOptionalPartsAndRichHeader(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Column(c).FillWidth().Padding(16).Gap(24).Children(func() {
			Message(c, MessageOptions{Header: "Replaced header", HeaderContent: func() {
				ui.Text(c, "Oliver · Today").Label("Rich header")
			}}, func() {
				Bubble(c, BubbleOptions{Variant: BubbleSecondary, Label: "Without avatar"}, func() { ui.Text(c, "Hello") })
			})
			Message(c, MessageOptions{Align: End}, func() {
				Bubble(c, BubbleOptions{Align: End, Label: "Without metadata"}, func() { ui.Text(c, "Hello again") })
			})
		})
	}, 390, 400)
	tt.SetPreferences(ui.Preferences{TextScale: 2, ReduceMotion: true})
	header, bubble := messageRect(t, tt, "Rich header"), messageRect(t, tt, "Without avatar")
	if header.Y+header.H > bubble.Y || tt.HasText("Replaced header") {
		t.Fatal("rich header overlapped bubble or failed to replace metadata")
	}
	assertMessageNear(t, bubble.X, 16, "avatar-free leading content")
	end := messageRect(t, tt, "Without metadata")
	assertMessageNear(t, end.X+end.W, 374, "avatar-free trailing content")
}

func TestBubbleReactionOverlayPreservesFootprintAndNativeActions(t *testing.T) {
	for _, scale := range []float32{1, 2} {
		for _, direction := range []ui.Direction{ui.LeftToRight, ui.RightToLeft} {
			for _, above := range []bool{false, true} {
				t.Run(fmt.Sprintf("scale=%g/direction=%d/above=%t", scale, direction, above), func(t *testing.T) {
					selected, disabled, visible, prepend := false, false, true, false
					count, calls := 1, 0
					tt := ui.NewTester(func(c *ui.Context) {
						theme.Apply(c)
						ui.Column(c).FillWidth().Padding(40).Children(func() {
							reactions := []Reaction{}
							if visible {
								if prepend {
									reactions = append(reactions, Reaction{Key: "eyes", Label: "Eyes", AccessibilityLabel: "Watching"})
								}
								reactions = append(reactions, Reaction{Key: "thumbs", Label: "Like", AccessibilityLabel: "Thumbs up",
									Count: count, Selected: selected, Disabled: disabled, OnClick: func() { selected = !selected; calls++ }})
							}
							Bubble(c, BubbleOptions{Variant: BubbleSecondary, Label: "Reaction bubble", Reactions: reactions,
								ReactionsAbove: above, ReserveReactionSpace: true}, func() {
								ui.Text(c, "A reply with a reaction badge.").Label("Reply text")
							})
							ui.Text(c, "Next message").Label("Next message")
						})
					}, 390, 600)
					tt.SetPreferences(ui.Preferences{TextScale: scale, ReduceMotion: true})
					tt.SetDirection(direction)
					bubble, next := messageRect(t, tt, "Reaction bubble"), messageRect(t, tt, "Next message")
					text, badge := messageRect(t, tt, "Reply text"), messageRect(t, tt, "Thumbs up, 1")
					if above {
						if badge.Y+badge.H > text.Y {
							t.Fatalf("top reaction covers text: badge=%+v text=%+v", badge, text)
						}
					} else if badge.Y < text.Y+text.H || badge.Y+badge.H > next.Y {
						t.Fatalf("bottom reaction covers text or following message: badge=%+v text=%+v next=%+v", badge, text, next)
					}
					if err := tt.Click("Thumbs up, 1"); err != nil {
						t.Fatal(err)
					}
					if calls != 1 || !selected || !tt.Focused("Thumbs up, 1") {
						t.Fatal("reaction did not activate and focus the named native toggle")
					}
					count, prepend = 999, true
					tt.Frame()
					if !tt.Focused("Thumbs up, 999") {
						t.Fatal("stable reaction identity lost focus when another badge was prepended")
					}
					tt.Key(0, ui.KeySpace)
					if calls != 2 || selected {
						t.Fatal("focused reaction did not toggle through native keyboard input")
					}
					if messageRect(t, tt, "Reaction bubble") != bubble || messageRect(t, tt, "Next message") != next {
						t.Fatal("reaction count or selected state changed the bubble footprint")
					}
					disabled = true
					tt.Frame()
					if err := tt.Click("Thumbs up, 999"); err != nil {
						t.Fatal(err)
					}
					if calls != 2 {
						t.Fatal("disabled reaction invoked callback")
					}
					visible = false
					tt.Frame()
					if messageRect(t, tt, "Reaction bubble") != bubble || messageRect(t, tt, "Next message") != next {
						t.Fatal("removing reactions discarded explicitly reserved attachment footprint")
					}
				})
			}
		}
	}
}

func messageRect(t *testing.T, tt *ui.Tester, label string) ui.Rect {
	t.Helper()
	r, ok := tt.Find(label)
	if !ok || r.W <= 0 || r.H <= 0 {
		t.Fatalf("missing visible %q: %+v", label, r)
	}
	return r
}

func assertMessageNear(t *testing.T, got, want float32, label string) {
	t.Helper()
	if math.Abs(float64(got-want)) > 1 {
		t.Fatalf("%s: got %g want %g", label, got, want)
	}
}
