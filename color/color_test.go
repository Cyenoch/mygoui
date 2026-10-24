package color

import (
	"math"
	"testing"

	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func near(a, b, tolerance float64) bool { return math.Abs(a-b) <= tolerance }

func TestParseAndFormatHex(t *testing.T) {
	tests := []struct {
		text string
		want ui.Color
		ok   bool
	}{
		{"#ff0000", ui.RGBA(255, 0, 0, 1), true},
		{"ff0000", ui.RGBA(255, 0, 0, 1), true},
		{"#0f0", ui.RGBA(0, 255, 0, 1), true},
		{"#0f08", ui.Color{R: 0, G: 255, B: 0, A: 0x88}, true},
		{"#33669980", ui.Color{R: 0x33, G: 0x66, B: 0x99, A: 0x80}, true},
		{"  #AbCdEf  ", ui.RGBA(0xab, 0xcd, 0xef, 1), true},
		{"#ff000", ui.Color{}, false},
		{"#gg0000", ui.Color{}, false},
		{"", ui.Color{}, false},
		{"#ff0000ff00", ui.Color{}, false},
	}
	for _, tt := range tests {
		got, err := ParseHex(tt.text)
		if tt.ok != (err == nil) {
			t.Fatalf("ParseHex(%q) error = %v, want ok %v", tt.text, err, tt.ok)
		}
		if tt.ok && got != tt.want {
			t.Fatalf("ParseHex(%q) = %+v, want %+v", tt.text, got, tt.want)
		}
	}
	if got := FormatHex(ui.Color{R: 0x33, G: 0x66, B: 0x99, A: 0x80}); got != "#33669980" {
		t.Fatalf("FormatHex translucent = %q", got)
	}
	if got := FormatHex(ui.RGB(1, 2, 3)); got != "#010203" {
		t.Fatalf("FormatHex opaque = %q", got)
	}
}

func TestColorStateKeepsHueThroughBlack(t *testing.T) {
	s := NewColorState(ui.Hex("#ff0000"))
	s.SetBrightness(0)
	if s.Hue() != 0 || s.Saturation() != 1 {
		t.Fatalf("black lost its hue: h=%v s=%v", s.Hue(), s.Saturation())
	}
	s.SetBrightness(1)
	if s.Value != ui.Hex("#ff0000") {
		t.Fatalf("restored %v, want red", s.Value)
	}
}

func TestColorAreaDragChangesColor(t *testing.T) {
	s := NewColorState(ui.Hex("#ff0000"))
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ColorArea(c, ColorAreaOptions{Label: "Fill color", State: s})
	}, 320, 320)
	r, ok := tt.Find("Fill color")
	if !ok {
		t.Fatal("area is not labeled")
	}
	if r.W > 224.5 || r.W < 100 || !near(float64(r.W), float64(r.H), 0.5) {
		t.Fatalf("area bounds %v x %v, want a square no wider than 224", r.W, r.H)
	}
	tt.Press(r.X+r.W*0.9, r.Y+r.H*0.1)
	tt.Move(r.X+r.W*0.25, r.Y+r.H*0.75)
	tt.Release(r.X+r.W*0.25, r.Y+r.H*0.75)
	if s.Value == ui.Hex("#ff0000") {
		t.Fatalf("drag left the color red")
	}
	if !near(s.Saturation(), 0.25, 0.03) || !near(s.Brightness(), 0.25, 0.03) {
		t.Fatalf("saturation/brightness = %v/%v, want about .25/.25",
			s.Saturation(), s.Brightness())
	}
}

func TestColorSliderChannelUpdates(t *testing.T) {
	s := NewColorState(ui.Hex("#ff0000"))
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ColorSlider(c, ColorSliderOptions{
			Label: "Hue", State: s, Channel: ChannelHue, ShowOutput: true,
		})
	}, 320, 100)
	// The label row is 20 DIPs and the track begins 4 DIPs below it.
	tt.ClickAt(160, 34)
	if !near(s.Hue(), 190, 15) {
		t.Fatalf("hue = %v, want the middle of the track", s.Hue())
	}
	if s.Saturation() != 1 || s.Brightness() != 1 {
		t.Fatalf("hue slider changed saturation/brightness: %v/%v",
			s.Saturation(), s.Brightness())
	}
}

func TestColorFieldParsesAndRejects(t *testing.T) {
	s := NewColorState(ui.Hex("#000000"))
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ColorField(c, ColorFieldOptions{Label: "Hex color", State: s})
	}, 320, 120)
	if err := tt.Click("Hex color"); err != nil {
		t.Fatal(err)
	}
	tt.Command("selectAll")
	tt.Type("#00ff00")
	if s.Value != ui.Hex("#00ff00") {
		t.Fatalf("field parsed %v, want green", s.Value)
	}
	tt.Command("selectAll")
	tt.Type("#zzzzzz")
	if s.Value != ui.Hex("#00ff00") {
		t.Fatalf("invalid text changed the color to %v", s.Value)
	}
	if !tt.HasText(locale.M().Invalid) {
		t.Fatalf("invalid text showed no error; screen: %v", tt.Texts())
	}
}

func TestColorSwatchPickerSelects(t *testing.T) {
	red, green, blue := ui.Hex("#ef4444"), ui.Hex("#22c55e"), ui.Hex("#3b82f6")
	value := red
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ColorSwatchPicker(c, ColorSwatchPickerOptions{
			Label:  "Palette",
			Colors: []ui.Color{red, green, blue},
			Names:  []string{"Red", "Green", "Blue"},
			Value:  &value,
		})
	}, 320, 120)
	if err := tt.Click("Blue"); err != nil {
		t.Fatal(err)
	}
	if value != blue {
		t.Fatalf("picker chose %v, want blue", value)
	}
}

func TestColorPickerSyncsValue(t *testing.T) {
	s := NewColorState(ui.Hex("#ff0000"))
	open := true
	blue := ui.Hex("#3b82f6")
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ColorPicker(c, ColorPickerOptions{
			Label: "Fill", State: s, Open: &open,
			Swatches:    []ui.Color{blue},
			SwatchNames: []string{"Blue"},
			HideAlpha:   true,
		})
	}, 520, 760)
	r, ok := tt.Find("Color area")
	if !ok {
		t.Fatalf("picker popover is not showing; screen: %v", tt.Texts())
	}
	tt.ClickAt(r.X+r.W*0.5, r.Y+r.H*0.5)
	if !near(s.Saturation(), 0.5, 0.05) || !near(s.Brightness(), 0.5, 0.05) {
		t.Fatalf("area set s/b = %v/%v, want about .5/.5",
			s.Saturation(), s.Brightness())
	}
	if err := tt.Click("Blue"); err != nil {
		t.Fatal(err)
	}
	if s.Value != blue {
		t.Fatalf("swatch set %v, want blue", s.Value)
	}
	if got := s.Hex(); got != "#3b82f6" {
		t.Fatalf("hex after swatch = %q", got)
	}
}

func TestColorSwatchPickerArrowRoving(t *testing.T) {
	red, green := ui.Hex("#ef4444"), ui.Hex("#22c55e")
	value := red
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ColorSwatchPicker(c, ColorSwatchPickerOptions{
			Label:  "Palette",
			Colors: []ui.Color{red, green},
			Names:  []string{"Red", "Green"},
			Value:  &value,
		})
	}, 320, 120)
	if err := tt.Click("Red"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyRight)
	if value != green {
		t.Fatalf("arrow chose %v, want green", value)
	}
}

func TestColorPickerTriggerOpens(t *testing.T) {
	s := NewColorState(ui.Hex("#ff0000"))
	open := false
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ColorPicker(c, ColorPickerOptions{Label: "Fill", State: s, Open: &open})
	}, 520, 760)
	if err := tt.Click("Fill"); err != nil {
		t.Fatal(err)
	}
	if !open {
		t.Fatal("trigger did not open the picker")
	}
	if _, ok := tt.Find("Color area"); !ok {
		t.Fatalf("open picker shows no area; screen: %v", tt.Texts())
	}
}
