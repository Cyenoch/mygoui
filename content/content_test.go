package content

import (
	"bytes"
	"errors"
	"image/color"
	"strings"
	"testing"
	"time"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func TestCopyButtonCopiesAndRevertsFeedback(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Column(c).Fill().Padding(24).Children(func() {
			CopyButton(c, CopyButtonOptions{TextFunc: func() string { return "go get example.com/x" }})
		})
	}, 240, 120)
	if err := tt.Click("Copy"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if got := tt.Clipboard(); got != "go get example.com/x" {
		t.Fatalf("clipboard = %q", got)
	}
	if _, ok := tt.Find("Copied"); !ok {
		t.Fatal("no copied state after a successful copy")
	}
	if got := tt.Announcements(); len(got) != 1 || got[0] != "Copied" {
		t.Fatalf("announcements = %q", got)
	}
	tt.Advance(2 * time.Second)
	if _, ok := tt.Find("Copied"); ok {
		t.Fatal("copied state did not revert")
	}
	if _, ok := tt.Find("Copy"); !ok {
		t.Fatal("copy label did not return")
	}
}

func TestCopyButtonReportsClipboardFailure(t *testing.T) {
	refused := errors.New("clipboard refused")
	var got error
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Column(c).Fill().Padding(24).Children(func() {
			CopyButton(c, CopyButtonOptions{Text: "secret", OnError: func(err error) { got = err }})
		})
	}, 240, 120)
	tt.FailClipboard(refused)
	if err := tt.Click("Copy"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if !errors.Is(got, refused) {
		t.Fatalf("OnError got %v", got)
	}
	if _, ok := tt.Find("Copied"); ok {
		t.Fatal("copied state shown although nothing was copied")
	}
}

func TestSnippetCopiesTextWithoutPrefix(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Column(c).Fill().Padding(24).Children(func() {
			Snippet(c, SnippetOptions{Prefix: "$", Text: "npm install"})
		})
	}, 320, 120)
	if !tt.HasText("$") || !tt.HasText("npm install") {
		t.Fatalf("snippet texts = %q", tt.Texts())
	}
	if err := tt.Click("Copy"); err != nil {
		t.Fatal(err)
	}
	if got := tt.Clipboard(); got != "npm install" {
		t.Fatalf("clipboard = %q", got)
	}
}

func TestHighlighterUnknownLanguageRendersPlainAndCopies(t *testing.T) {
	const code = "func main() { return \"hi\" } // 42"
	render := func(language string) *ui.Tester {
		return ui.NewTester(func(c *ui.Context) {
			theme.Apply(c)
			ui.Column(c).Fill().Padding(24).Children(func() {
				Highlighter(c, HighlighterOptions{Code: "\n" + code + "\n", Language: language, Title: "snippet"})
			})
		}, 480, 200)
	}
	plain, unknown, known := render(""), render("basic"), render("go")
	if !bytes.Equal(unknown.Image().Pix, plain.Image().Pix) {
		t.Fatal("an unknown language is not rendered as plain code")
	}
	if bytes.Equal(known.Image().Pix, plain.Image().Pix) {
		t.Fatal("a known language is not coloured")
	}
	if err := unknown.Click("Copy"); err != nil {
		t.Fatal(err)
	}
	if got := unknown.Clipboard(); got != code {
		t.Fatalf("clipboard = %q", got)
	}
}

func TestHighlightColoursTokensAndKeepsSource(t *testing.T) {
	colors := syntaxColors{
		comment: ui.RGB(1, 0, 0), str: ui.RGB(2, 0, 0),
		number: ui.RGB(3, 0, 0), keyword: ui.RGB(4, 0, 0),
	}
	tests := []struct {
		language, src string
		want          map[string]ui.Color
	}{
		{"go", "func main() { // 日本語 \"x\"\n\tx := `a\nb` + \"q\\\"\" + 'c' /* k */ + 0x1F\n\tx.type = nil\n}",
			map[string]ui.Color{
				"func": colors.keyword, "// 日本語 \"x\"": colors.comment, "`a\nb`": colors.str,
				`"q\""`: colors.str, "'c'": colors.str, "/* k */": colors.comment, "0x1F": colors.number, "nil": colors.keyword,
			}},
		{"ts", "const a = 'it' // c\nlet b = `t\n${a}`", map[string]ui.Color{
			"const": colors.keyword, "let": colors.keyword, "'it'": colors.str, "// c": colors.comment, "`t\n${a}`": colors.str}},
		{"json", `{"key": "value", "n": 12.5, "ok": true}`, map[string]ui.Color{
			`"value"`: colors.str, "12.5": colors.number, "true": colors.keyword}},
		{"bash", "if [ $# -gt 0 ]; then echo \"hi\" # note\nfi", map[string]ui.Color{
			"if": colors.keyword, "0": colors.number, "then": colors.keyword, `"hi"`: colors.str, "# note": colors.comment, "fi": colors.keyword}},
		{"python", "def f():\n    '''doc\n    string'''\n    return don't\nx = 'y'", map[string]ui.Color{
			"def": colors.keyword, "'''doc\n    string'''": colors.str, "return": colors.keyword, "'y'": colors.str}},
		{"yaml", "no-cache: yes # c\nname: \"x\"", map[string]ui.Color{
			"yes": colors.keyword, "# c": colors.comment, `"x"`: colors.str}},
		{"html", "<!-- c --><a href=\"u\">it's 5</a>", map[string]ui.Color{
			"<!-- c -->": colors.comment, "a": colors.keyword, `"u"`: colors.str}},
		{"css", "a { width: 12px; /* c */ }", map[string]ui.Color{
			"12px": colors.number, "/* c */": colors.comment}},
	}
	for _, tc := range tests {
		t.Run(tc.language, func(t *testing.T) {
			s, ok := syntaxFor(tc.language)
			if !ok {
				t.Fatal("language not known")
			}
			spans := highlight(tc.src, s, colors)
			var joined strings.Builder
			coloured := map[string]ui.Color{}
			for _, span := range spans {
				joined.WriteString(span.Text)
				if span.Color.A != 0 {
					coloured[span.Text] = span.Color
				}
			}
			if joined.String() != tc.src {
				t.Fatalf("spans lost source text:\n got %q\nwant %q", joined.String(), tc.src)
			}
			if len(coloured) != len(tc.want) {
				t.Errorf("coloured %d tokens, want %d: %v", len(coloured), len(tc.want), coloured)
			}
			for text, want := range tc.want {
				if coloured[text] != want {
					t.Errorf("token %q colour = %v, want %v", text, coloured[text], want)
				}
			}
		})
	}
}

func TestMarkdownReportsLinksAndHighlightsCodeBlocks(t *testing.T) {
	var opened string
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Column(c).Fill().Padding(24).Children(func() {
			Markdown(c, MarkdownOptions{
				Source: "# Title\n\nRead the [docs](https://example.com/docs).\n\n```go\nfunc main() {}\n```\n",
				OnLink: func(url string) { opened = url },
			})
		})
	}, 520, 360)
	if err := tt.Click("docs"); err != nil {
		t.Fatal(err)
	}
	if opened != "https://example.com/docs" {
		t.Fatalf("OnLink got %q", opened)
	}
	if len(tt.OpenedURLs()) != 0 {
		t.Fatalf("link opened natively despite OnLink: %v", tt.OpenedURLs())
	}
	if err := tt.Click("Copy"); err != nil {
		t.Fatalf("code block has no Highlighter copy action: %v", err)
	}
	if got := tt.Clipboard(); got != "func main() {}" {
		t.Fatalf("clipboard = %q", got)
	}
}

func TestQRCodeReportsUnencodableValue(t *testing.T) {
	var err error
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Column(c).Fill().Children(func() {
			_, err = QRCode(c, QRCodeOptions{Value: strings.Repeat("x", 4000), Label: "Login code"})
		})
	}, 240, 240)
	if err == nil {
		t.Fatal("no error for a value beyond QR capacity")
	}
	if _, ok := tt.Find("Login code"); ok {
		t.Fatal("a code was drawn for an unencodable value")
	}
}

func TestQRCodePaintsModulesInsideQuietZone(t *testing.T) {
	fg, bg := ui.RGB(0, 0, 0), ui.RGB(255, 255, 255)
	var err error
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Column(c).Fill().Children(func() {
			_, err = QRCode(c, QRCodeOptions{
				Value: "https://example.com", Size: 200, Foreground: fg, Background: bg, Label: "Site",
			})
		})
	}, 240, 240)
	if err != nil {
		t.Fatal(err)
	}
	box, ok := tt.Find("Site")
	if !ok {
		t.Fatal("QR code has no accessible label")
	}
	code, _ := encodeQRForTest(t, "https://example.com")
	unit := box.W / float32(code+2*qrQuietModules)
	img := tt.Image()
	scale := float32(img.Bounds().Dx()) / 240
	at := func(mx, my float32) color.RGBA {
		return img.RGBAAt(int((box.X+mx*unit)*scale), int((box.Y+my*unit)*scale))
	}
	// The quiet zone stays clear; the finder pattern's corner module is
	// dark, its inner ring light and its centre dark.
	if got := at(2, 2); got.R != 255 {
		t.Errorf("quiet zone pixel = %v", got)
	}
	if got := at(4.5, 4.5); got.R != 0 {
		t.Errorf("finder corner pixel = %v", got)
	}
	if got := at(5.5, 5.5); got.R != 255 {
		t.Errorf("finder ring pixel = %v", got)
	}
	if got := at(7.5, 7.5); got.R != 0 {
		t.Errorf("finder centre pixel = %v", got)
	}
}

func encodeQRForTest(t *testing.T, value string) (int, []bool) {
	t.Helper()
	st := encodeQR(value, QRMedium)
	if st.err != nil {
		t.Fatal(st.err)
	}
	return st.side, st.modules
}
