package fonts_test

import (
	"fmt"
	"math"
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/Cyenoch/mygoui/fonts"
	"github.com/egoist/mygo/ui"
)

func TestRegister(t *testing.T) {
	if runtime.GOOS == "js" {
		if err := fonts.Register(); err == nil || !strings.Contains(err.Error(), "browser font registration") {
			t.Fatalf("browser registration error = %v", err)
		}
		return
	}

	const sample = "MWa"
	faces := []struct {
		name, style, reference string
		italic                 bool
		data                   []byte
	}{
		{name: "CalSansVF.ttf", style: "Regular", reference: "Cal Sans Test Regular"},
		{name: "CalSansVF-Italic.ttf", style: "Italic", reference: "Cal Sans Test Italic", italic: true},
	}
	for i := range faces {
		face := &faces[i]
		data, err := os.ReadFile(face.name)
		if err != nil {
			t.Fatalf("read %s: %v", face.name, err)
		}
		face.data = data
	}

	const callers = 8
	results := make(chan error, callers)
	for range callers {
		go func() { results <- fonts.Register() }()
	}
	for range callers {
		if err := <-results; err != nil {
			t.Fatalf("register Cal Sans: %v", err)
		}
	}

	found := false
	for _, family := range ui.FontFamilies() {
		if family.Name == fonts.Family && family.Registered {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("%q is not listed as a registered font family", fonts.Family)
	}

	weights := []int{400, 700}
	shaped := make([][][]ui.Glyph, len(faces))
	for i, face := range faces {
		shaped[i] = make([][]ui.Glyph, len(weights))
		for j, weight := range weights {
			name := fmt.Sprintf("%s/%d", face.style, weight)
			if !t.Run(name, func(t *testing.T) {
				shaped[i][j] = ui.Shape(sample, ui.Font{
					Family: fonts.Family, Size: 64, Weight: weight, Italic: face.italic,
				})
				if len(shaped[i][j]) == 0 {
					t.Fatal("no shaped glyphs")
				}
				tester := ui.NewTester(func(c *ui.Context) {
					ui.Column(c).Fill().Background(ui.RGB(255, 255, 255)).Children(func() {
						text := ui.Text(c, sample).Font(fonts.Family).FontSize(64).
							FontWeight(weight).TextColor(ui.RGB(0, 0, 0))
						if face.italic {
							text.Italic()
						}
					})
				}, 256, 128)
				if _, ok := tester.Find(sample); !ok {
					t.Fatal("text is missing from the rendered view")
				}
				img := tester.Image()
				for y := img.Rect.Min.Y; y < img.Rect.Max.Y; y++ {
					for x := img.Rect.Min.X; x < img.Rect.Max.X; x++ {
						pixel := img.RGBAAt(x, y)
						if pixel.R < 255 || pixel.G < 255 || pixel.B < 255 {
							return
						}
					}
				}
				t.Fatal("text drew no visible ink")
			}) {
				return
			}
		}
	}

	// Compare with each original face registered under a separate family.
	// Capture the package's shapes first: registration can also expose a
	// font under its original family name.
	for i, face := range faces {
		if err := ui.RegisterFont(face.data, face.reference); err != nil {
			t.Fatalf("register reference %s: %v", face.name, err)
		}
		for j, weight := range weights {
			reference := ui.Shape(sample, ui.Font{
				Family: face.reference, Size: 64, Weight: weight, Italic: face.italic,
			})
			got := shaped[i][j]
			if len(got) != len(reference) {
				t.Fatalf("%s/%d produced %d glyphs, reference produced %d", face.style, weight, len(got), len(reference))
			}
			for k := range got {
				if got[k].Cluster != reference[k].Cluster || got[k].Runes != reference[k].Runes ||
					math.Abs(float64(got[k].Advance-reference[k].Advance)) > 0.01 {
					t.Fatalf("%s/%d glyph %d = %+v, reference %+v", face.style, weight, k, got[k], reference[k])
				}
			}
		}
		if math.Abs(float64(shaped[i][0][0].Advance-shaped[i][1][0].Advance)) < 0.1 {
			t.Errorf("%s text did not change advance between weights 400 and 700", face.style)
		}
	}
}
