package theme

import (
	"math"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestHostPreferencesComposeWithSystem(t *testing.T) {
	previous := preferenceSnapshot()
	t.Cleanup(func() { SetPreferences(previous) })
	for _, tc := range []struct {
		name       string
		app        Preferences
		os         ui.Preferences
		wantScale  float32
		wantMotion bool
	}{
		{name: "defaults", os: ui.Preferences{TextScale: 1}, wantScale: 1},
		{name: "numeric size", app: Preferences{FontSize: 18}, os: ui.Preferences{TextScale: 1}, wantScale: 1.125},
		{name: "combined scale", app: Preferences{FontSize: 20}, os: ui.Preferences{TextScale: 1.2}, wantScale: 1.5},
		{name: "panels disabled preserve progress", os: ui.Preferences{TextScale: 1}, wantScale: 1},
		{name: "OS motion cannot be disabled", app: Preferences{PanelAnimationDuration: 400 * time.Millisecond}, os: ui.Preferences{TextScale: 1, ReduceMotion: true}, wantScale: 1, wantMotion: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			SetPreferences(tc.app)
			var scale, font float32
			var motion bool
			var os ui.Preferences
			tt := ui.NewTester(func(c *ui.Context) {
				scale, motion, os = TextScale(c), ReduceMotion(c), c.Preferences()
				font = Apply(c).FontSize
				ui.Text(c, "Readable")
			}, 300, 100)
			tt.SetPreferences(tc.os)
			tt.Frame()
			if math.Abs(float64(scale-tc.wantScale)) > 0.0001 || motion != tc.wantMotion {
				t.Fatalf("effective preferences = (%g, %v), want (%g, %v)", scale, motion, tc.wantScale, tc.wantMotion)
			}
			if math.Abs(float64(font-16*tc.wantScale)) > 0.0001 {
				t.Errorf("root font = %g, want %g without compounding", font, 16*tc.wantScale)
			}
			if os != tc.os {
				t.Errorf("read-only OS preferences changed: %+v, want %+v", os, tc.os)
			}
		})
	}
}

func TestInterfaceFamilyAndPanelDuration(t *testing.T) {
	previous := preferenceSnapshot()
	t.Cleanup(func() { SetPreferences(previous) })
	SetPreferences(Preferences{FontFamily: "monospace", FontSize: 16, PanelAnimationDuration: 200 * time.Millisecond})
	var family string
	var root, explicit float32
	var duration time.Duration
	tt := ui.NewTester(func(c *ui.Context) {
		applied := Apply(c)
		family, root = applied.Font, applied.FontSize
		explicit, duration = Scale(c, 14), PanelAnimationDuration(c)
		ui.Text(c, "Readable").FontSize(explicit)
	}, 400, 160)
	if family != "monospace, "+Sans || root != 16 || explicit != 14 || duration != 200*time.Millisecond {
		t.Fatalf("source defaults/family = (%q, %g, %g, %s)", family, root, explicit, duration)
	}
	tt.SetPreferences(ui.Preferences{TextScale: 1.5, ReduceMotion: true})
	if explicit != 21 || root != 24 || duration != 0 {
		t.Fatalf("OS accessibility = (explicit %g, root %g, panel %s), want (21, 24, 0)", explicit, root, duration)
	}
	SetPreferences(Preferences{FontSize: 20})
	tt.SetPreferences(ui.Preferences{TextScale: 1})
	if explicit != 17.5 || root != 20 || duration != 0 || family != Sans {
		t.Fatalf("numeric edit/reset = (%q, %g, %g, %s)", family, root, explicit, duration)
	}
}

func TestNumericInterfaceFontChangesGlyphLayout(t *testing.T) {
	previous := preferenceSnapshot()
	t.Cleanup(func() { SetPreferences(previous) })
	SetPreferences(Preferences{FontSize: 16})
	const text = "Readable interface text"
	tt := ui.NewTester(func(c *ui.Context) {
		Apply(c)
		ui.Text(c, text).Font(FontFamily(c)).FontSize(Scale(c, 14)).
			FixedLineHeight(Scale(c, 20)).NoWrap().AlignSelf(ui.Start)
	}, 500, 160)
	before, ok := tt.Find(text)
	if !ok {
		t.Fatal("interface glyphs missing")
	}
	for _, size := range []float32{18, 20} {
		SetPreferences(Preferences{FontSize: size})
		tt.Frame()
		after, ok := tt.Find(text)
		if !ok || after.H <= before.H || after.W <= before.W {
			t.Fatalf("size %g layout = %+v, baseline %+v", size, after, before)
		}
	}
	SetPreferences(Preferences{FontSize: 16, FontFamily: "monospace"})
	tt.Frame()
	mono, ok := tt.Find(text)
	if !ok || mono.W == before.W {
		t.Fatalf("chosen family did not change glyph advances: system %+v, monospace %+v", before, mono)
	}
}
