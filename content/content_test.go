package content

import (
	"bytes"
	"errors"
	"fmt"
	"image"
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
	if err := unknown.Click("Copy snippet"); err != nil {
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

// TestDiffViewerCountsFollowTheTexts pins the diff's input contract: the
// header counts the changed lines of the caller's texts, identical texts
// show an empty state, and the split layout adds the old side.
func TestDiffViewerCountsFollowTheTexts(t *testing.T) {
	o := DiffViewerOptions{
		Old: "alpha\nbeta\ngamma\n", New: "alpha\nBETA\ngamma\ndelta\n", Title: "notes.txt",
	}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Column(c).Fill().Padding(16).Children(func() { DiffViewer(c, o) })
	}, 640, 360)
	if !tt.HasText("notes.txt") || !tt.HasText("+2") || !tt.HasText("−1") {
		t.Fatalf("header = %q", tt.Texts())
	}
	if _, ok := tt.Access("Old text"); ok {
		t.Fatal("a unified diff exposes an old side")
	}
	o.New, o.Layout = o.Old, DiffSplit
	tt.Frame()
	if !tt.HasText("No changes") || tt.HasText("+0") {
		t.Fatalf("identical texts show %q", tt.Texts())
	}
	o.New = "alpha\nBETA\n"
	tt.Frame()
	for _, side := range []string{"Old text", "New text"} {
		if _, ok := tt.Access(side); !ok {
			t.Fatalf("split diff has no %q", side)
		}
	}
}

// rowShows reports whether a text of the frame is exactly s, as a key of the
// tree is: a closed object's preview also holds its keys.
func rowShows(tt *ui.Tester, s string) bool {
	for _, text := range tt.Texts() {
		if text == s {
			return true
		}
	}
	return false
}

const jsonSample = `{"zeta": 1, "alpha": {"name": "Ada \"L\"", "tags": ["x", "y"]}, "ok": true}`

func jsonTester(o *JSONViewerOptions, err *error) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Column(c).Fill().Padding(16).Children(func() { _, *err = JSONViewer(c, *o) })
	}, 480, 480)
}

// TestJSONViewerKeepsOrderAndExpansion pins the tree: members keep the
// document's order, nested values stay hidden until opened, and nodes open
// by path survive a new document.
func TestJSONViewerKeepsOrderAndExpansion(t *testing.T) {
	var st JSONViewerState
	var err error
	o := JSONViewerOptions{Source: jsonSample, State: &st, ExpandDepth: -1}
	tt := jsonTester(&o, &err)
	if err != nil {
		t.Fatal(err)
	}
	zeta, zok := tt.Find("zeta")
	alpha, aok := tt.Find("alpha")
	if !zok || !aok || zeta.Y >= alpha.Y {
		t.Fatalf("members out of document order: zeta=%+v alpha=%+v in %q", zeta, alpha, tt.Texts())
	}
	if rowShows(tt, "name") {
		t.Fatal("a nested member shows before its object is opened")
	}
	st.ExpandAll()
	tt.Frame()
	if !rowShows(tt, "name") || !tt.HasText(`"y"`) {
		t.Fatalf("expanded tree = %q", tt.Texts())
	}
	o.Source = `{"alpha": {"name": "Grace", "tags": ["z"]}, "extra": {"deep": 1}}`
	tt.Frame()
	if err != nil {
		t.Fatal(err)
	}
	if !tt.HasText(`"Grace"`) || !tt.HasText(`"z"`) {
		t.Fatalf("open nodes did not survive the new document: %q", tt.Texts())
	}
	if rowShows(tt, "deep") {
		t.Fatal("a node new to the document opened by itself")
	}
	st.CollapseAll()
	tt.Frame()
	if rowShows(tt, "name") || !tt.HasText("alpha") {
		t.Fatalf("collapsed tree = %q", tt.Texts())
	}
	for _, source := range []string{`{}`, `{"alpha": 1}`} {
		o.Source = source
		tt.Frame()
	}
	if !tt.HasText("alpha") {
		t.Fatalf("the root closed after an empty document: %q", tt.Texts())
	}
}

// TestJSONViewerCopiesValueAndPath pins the row menu: it copies a string's
// characters, not its JSON text, and the node's path.
func TestJSONViewerCopiesValueAndPath(t *testing.T) {
	var st JSONViewerState
	var err error
	o := JSONViewerOptions{Source: jsonSample, State: &st, ExpandDepth: -1}
	tt := jsonTester(&o, &err)
	st.ExpandAll()
	tt.Frame()
	for item, want := range map[string]string{"Copy value": `Ada "L"`, "Copy path": "$.alpha.name"} {
		if err := tt.RightClick("name"); err != nil {
			t.Fatal(err)
		}
		if err := tt.ChooseMenuItem(item); err != nil {
			t.Fatalf("%s: %v in %q", item, err, tt.Menu())
		}
		tt.Frame()
		if got := tt.Clipboard(); got != want {
			t.Fatalf("%s wrote %q, want %q", item, got, want)
		}
	}
}

// TestJSONViewerReportsInvalidDocument pins the failure contract: invalid
// JSON, and data after the document, return an error naming where, and
// the next valid document clears it.
func TestJSONViewerReportsInvalidDocument(t *testing.T) {
	var err error
	o := JSONViewerOptions{Source: `{"a": [1, 2,]}`}
	tt := jsonTester(&o, &err)
	if err == nil || !strings.Contains(err.Error(), "byte") {
		t.Fatalf("trailing comma: err = %v", err)
	}
	if !tt.HasText(err.Error()) {
		t.Fatalf("the viewer does not show its error: %q", tt.Texts())
	}
	o.Source = `{"a": 1} 2`
	tt.Frame()
	if err == nil {
		t.Fatal("data after the document parsed")
	}
	o.Source = `{"a": 1}`
	tt.Frame()
	if err != nil || !tt.HasText("a") {
		t.Fatalf("valid document: err = %v, texts %q", err, tt.Texts())
	}
}

// TestJSONViewerBuildsOnlyRowsInView pins virtualization: a fully open
// document of ten thousand values builds a screenful of rows.
func TestJSONViewerBuildsOnlyRowsInView(t *testing.T) {
	var b strings.Builder
	b.WriteByte('[')
	for i := range 10000 {
		if i > 0 {
			b.WriteByte(',')
		}
		fmt.Fprintf(&b, "%d", i)
	}
	b.WriteByte(']')
	var st JSONViewerState
	var err error
	o := JSONViewerOptions{Source: b.String(), State: &st}
	tt := jsonTester(&o, &err)
	st.ExpandAll()
	tt.Frame()
	if err != nil {
		t.Fatal(err)
	}
	if n := len(tt.Texts()); n < 10 || n > 200 {
		t.Fatalf("an open 10k array builds %d texts", n)
	}
}

// TestLogViewerFollowsUntilScrolledAway pins the tail: new lines stay in
// view while the log is at its end, scrolling up holds the place and
// offers a way back, and taking it follows again.
func TestLogViewerFollowsUntilScrolledAway(t *testing.T) {
	var st LogViewerState
	n := 200
	line := func(i int) string { return fmt.Sprintf("entry %d", i) }
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Column(c).Fill().Padding(16).Children(func() {
			LogViewer(c, LogViewerOptions{
				Line: line, Count: n, Follow: true, LineNumbers: true, Height: 200, State: &st,
			})
		})
	}, 480, 320)
	tt.Frame()
	if !tt.HasText("entry 199") {
		t.Fatalf("a following log does not start at its end: %q", tt.Texts())
	}
	if _, ok := tt.Find("Jump to latest"); ok {
		t.Fatal("the way back shows while the log is at its end")
	}
	n = 210
	tt.Frame()
	tt.Frame()
	if !tt.HasText("entry 209") {
		t.Fatalf("the log did not follow new lines: %q", tt.Texts())
	}
	tt.Scroll(200, 120, 0, -600)
	tt.Frame()
	if tt.HasText("entry 209") {
		t.Fatal("the wheel did not scroll the log")
	}
	n = 220
	tt.Frame()
	tt.Frame()
	if tt.HasText("entry 219") {
		t.Fatal("the log followed new lines after the user scrolled away")
	}
	if err := tt.Click("Jump to latest"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	tt.Frame()
	if !tt.HasText("entry 219") {
		t.Fatalf("the way back did not reach the newest line: %q", tt.Texts())
	}
	if _, ok := tt.Find("Jump to latest"); ok {
		t.Fatal("the way back stays after reaching the end")
	}
}

func logFilterSample(i int) string {
	return fmt.Sprintf("entry %d %s", i, [...]string{"alpha", "BETA", "gamma"}[i%3])
}

// TestLogViewerQueryAndLevelsListMatchingLines pins the filter: only the
// lines holding the query, whatever its case, are listed, under their own
// numbers, and a level set lists the lines of those levels.
func TestLogViewerQueryAndLevelsListMatchingLines(t *testing.T) {
	var st LogViewerState
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Column(c).Fill().Padding(16).Children(func() {
			LogViewer(c, LogViewerOptions{
				Line: logFilterSample, Count: 9, LineNumbers: true, State: &st,
				Level: func(i int) LogLevel { return LogLevel(1 + i%2) },
			})
		})
	}, 480, 320)
	st.Query = "beta"
	tt.Frame()
	if !tt.HasText("entry 1 BETA") || !tt.HasText("entry 7 BETA") || tt.HasText("entry 0 alpha") {
		t.Fatalf("the query did not narrow the lines: %q", tt.Texts())
	}
	if !tt.HasText("2") || !tt.HasText("8") {
		t.Fatalf("the lines lost their own numbers: %q", tt.Texts())
	}
	st.Query, st.Levels = "", LevelSet(LogInfo)
	tt.Frame()
	if !tt.HasText("entry 1 BETA") || tt.HasText("entry 0 alpha") {
		t.Fatalf("the level set did not narrow the lines: %q", tt.Texts())
	}
	st.Query = "nothing like this"
	tt.Frame()
	if !tt.HasText("No results") {
		t.Fatalf("an empty listing says nothing: %q", tt.Texts())
	}
}

// TestLogViewerFilterFollowsAppendedLines pins a filtered tail: a new line
// that matches shows at once, one that does not never does, and the log is
// not scanned again while nothing changes.
func TestLogViewerFilterFollowsAppendedLines(t *testing.T) {
	var st LogViewerState
	lines := make([]string, 100)
	var reads [200]int
	for i := range lines {
		lines[i] = fmt.Sprintf("entry %d", i)
	}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Column(c).Fill().Padding(16).Children(func() {
			LogViewer(c, LogViewerOptions{
				Line:  func(i int) string { reads[i]++; return lines[i] },
				Count: len(lines), Follow: true, Height: 120, State: &st,
			})
		})
	}, 480, 320)
	st.Query = "entry 9"
	tt.Frame()
	if !tt.HasText("entry 99") {
		t.Fatalf("a filtered log does not start at its end: %q", tt.Texts())
	}
	reads = [200]int{}
	tt.Frame()
	tt.Frame()
	if reads[50] != 0 {
		t.Fatal("an unchanged log was scanned again")
	}
	lines = append(lines, "entry 100", "entry 9 again")
	tt.Frame()
	tt.Frame()
	if !tt.HasText("entry 9 again") || tt.HasText("entry 100") {
		t.Fatalf("the appended lines were not filtered: %q", tt.Texts())
	}
	if reads[50] != 0 {
		t.Fatal("an appended line scanned the whole log again")
	}
}

// TestLogViewerToolbarSearches pins the toolbar: the shortcut reaches the
// search field, typing narrows the log, Escape clears it, and the toolbar
// at rest asks for no frame.
func TestLogViewerToolbarSearches(t *testing.T) {
	var st LogViewerState
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Column(c).Fill().Padding(16).Children(func() {
			LogViewer(c, LogViewerOptions{
				Line: logFilterSample, Count: 30, Toolbar: true, Height: 200, State: &st,
			})
		})
	}, 480, 320)
	tt.Frame()
	if tt.HasText("/ 30") {
		t.Fatalf("the count shows without a filter: %q", tt.Texts())
	}
	if line, ok := tt.Find("entry 0 alpha"); ok {
		tt.ClickAt(line.X+line.W/2, line.Y+line.H/2)
	} else {
		t.Fatalf("no first line in %q", tt.Texts())
	}
	tt.Key(ui.Cmd, ui.KeyF)
	tt.Frame()
	tt.Type("gamma")
	tt.Frame()
	if st.Query != "gamma" || tt.HasText("entry 0 alpha") || !tt.HasText("10 / 30") {
		t.Fatalf("typing did not narrow the log: query %q, %q", st.Query, tt.Texts())
	}
	tt.Key(0, ui.KeyEscape)
	tt.Frame()
	if st.Query != "" || !tt.HasText("entry 0 alpha") {
		t.Fatalf("Escape did not clear the query: %q", st.Query)
	}
	tt.Advance(time.Second)
	if tt.FramePending() {
		t.Fatal("the toolbar keeps asking for frames at rest")
	}
}

// TestLogViewerJumpButtonFadesIn pins the way back arriving: after the user
// scrolls away it passes through a half-shown frame before it settles.
func TestLogViewerJumpButtonFadesIn(t *testing.T) {
	var st LogViewerState
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Column(c).Fill().Padding(16).Children(func() {
			LogViewer(c, LogViewerOptions{
				Line: func(i int) string { return fmt.Sprintf("entry %d", i) }, Count: 200,
				Follow: true, Height: 200, State: &st,
			})
		})
	}, 480, 320)
	tt.Advance(time.Second)
	if _, ok := tt.Find("Jump to latest"); ok {
		t.Fatal("the way back shows while the log is at its end")
	}
	tt.Scroll(200, 120, 0, -600)
	tt.Advance(50 * time.Millisecond)
	mid := *tt.Image()
	mid.Pix = append([]uint8(nil), mid.Pix...)
	tt.Advance(time.Second)
	rect, ok := tt.Find("Jump to latest")
	if !ok {
		t.Fatal("the way back never showed")
	}
	// How far the button's box is from the dark log behind it.
	strength := func(img *image.RGBA) (sum int) {
		for y := int(rect.Y); y < int(rect.Y+rect.H); y++ {
			for x := int(rect.X); x < int(rect.X+rect.W); x++ {
				p := img.RGBAAt(x, y)
				sum += int(p.R) + int(p.G) + int(p.B)
			}
		}
		return sum
	}
	hidden := strength(&mid)
	shown := strength(tt.Image())
	if hidden == shown {
		t.Fatal("the way back did not change between the half and the settled frame")
	}
	if tt.FramePending() {
		t.Fatal("the settled button keeps asking for frames")
	}
}

// A double click on a container's row opens and closes it, as its arrow
// does.
func TestJSONViewerDoubleClickTogglesAContainer(t *testing.T) {
	var st JSONViewerState
	var err error
	o := JSONViewerOptions{Source: jsonSample, State: &st, ExpandDepth: -1}
	tt := jsonTester(&o, &err)
	alpha, ok := tt.Find("alpha")
	if !ok {
		t.Fatalf("no alpha row in %q", tt.Texts())
	}
	x, y := alpha.X+alpha.W/2, alpha.Y+alpha.H/2
	tt.ClickAt(x, y)
	tt.ClickAt(x, y)
	if !rowShows(tt, "name") {
		t.Fatalf("a double click did not open the object: %q", tt.Texts())
	}
	tt.Advance(time.Second)
	tt.ClickAt(x, y)
	tt.ClickAt(x, y)
	if rowShows(tt, "name") {
		t.Fatal("a second double click did not close the object")
	}
}

// A query opens the objects above its match and chooses it; Cmd+C writes
// the chosen row's value, and Next steps to the following match.
func TestJSONViewerSearchRevealsMatchesAndCopies(t *testing.T) {
	var st JSONViewerState
	var err error
	o := JSONViewerOptions{Source: jsonSample, State: &st, ExpandDepth: -1, Toolbar: true}
	tt := jsonTester(&o, &err)
	if rowShows(tt, "name") {
		t.Fatal("nested member shows before the search")
	}
	st.Query = "name"
	tt.Frame()
	tt.Frame()
	if !rowShows(tt, "name") || !tt.HasText("1 / 1") {
		t.Fatalf("the match is not revealed: %q", tt.Texts())
	}
	if err := tt.Click("name"); err != nil {
		t.Fatal(err)
	}
	tt.Key(ui.Cmd, ui.KeyC)
	if got := tt.Clipboard(); got != `Ada "L"` {
		t.Fatalf("Cmd+C wrote %q", got)
	}
	st.Query = "a"
	tt.Frame()
	if !tt.HasText("1 / 4") {
		t.Fatalf("matches of a: %q", tt.Texts())
	}
	if err := tt.Click("Next"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if !tt.HasText("2 / 4") {
		t.Fatalf("Next did not step: %q", tt.Texts())
	}
}

// Next and Previous scroll to the following change and say which it is.
func TestDiffViewerStepsThroughChanges(t *testing.T) {
	var old, next strings.Builder
	for i := range 60 {
		fmt.Fprintf(&old, "line %d\n", i)
		if i == 5 || i == 50 {
			fmt.Fprintf(&next, "changed %d\n", i)
			continue
		}
		fmt.Fprintf(&next, "line %d\n", i)
	}
	o := DiffViewerOptions{Old: old.String(), New: next.String(), Title: "f.txt", Height: 200}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Column(c).Fill().Padding(16).Children(func() { DiffViewer(c, o) })
	}, 640, 360)
	for _, want := range []string{"1 / 2", "2 / 2", "1 / 2"} {
		if err := tt.Click("Next"); err != nil {
			t.Fatal(err)
		}
		tt.Frame()
		if !tt.HasText(want) {
			t.Fatalf("want %q after Next: %q", want, tt.Texts())
		}
	}
	tt.Key(ui.Shift, ui.KeyF7)
	tt.Frame()
	if !tt.HasText("2 / 2") {
		t.Fatalf("Shift+F7 did not step back: %q", tt.Texts())
	}
}

// A new document under the same query keeps the match the user stepped to.
func TestJSONViewerSearchKeepsItsMatchWhenTheSourceChanges(t *testing.T) {
	var st JSONViewerState
	var err error
	o := JSONViewerOptions{Source: jsonSample, State: &st, ExpandDepth: -1, Toolbar: true}
	tt := jsonTester(&o, &err)
	st.Query = "a"
	tt.Frame()
	_ = tt.Click("Next")
	tt.Frame()
	o.Source = strings.Replace(jsonSample, `"zeta": 1`, `"zeta": 2`, 1)
	tt.Frame()
	if !tt.HasText("2 / 4") {
		t.Fatalf("a new source moved the match: %q", tt.Texts())
	}
}

// Raising Revision under a filter lists again without throwing the reader
// back to the top, and a plain line follows the Info toggle.
func TestLogViewerFilterKeepsScrollOnRevisionAndListsPlainAsInfo(t *testing.T) {
	var st LogViewerState
	o := LogViewerOptions{Line: logFilterSample, Count: 90, Height: 120, State: &st}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Column(c).Fill().Padding(16).Children(func() { LogViewer(c, o) })
	}, 480, 320)
	st.Query = "entry"
	tt.Frame()
	st.List.ScrollTo(60, ui.Start)
	tt.Frame()
	st.Revision++
	tt.Frame()
	if tt.HasText("entry 0 alpha") || !tt.HasText("entry 6") {
		t.Fatalf("a new revision scrolled to the top: %q", tt.Texts())
	}
	st.Query, st.Levels = "", LevelSet(LogInfo)
	o.Level = func(i int) LogLevel {
		if i%2 == 0 {
			return LogPlain
		}
		return LogError
	}
	if !st.Levels.Shows(LogPlain) || st.Levels.Shows(LogError) {
		t.Fatal("a plain line does not follow Info")
	}
}

// Cmd+F still reaches the search field when a search finds nothing.
func TestLogViewerSearchShortcutWorksWithoutResults(t *testing.T) {
	var st LogViewerState
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Column(c).Fill().Padding(16).Children(func() {
			LogViewer(c, LogViewerOptions{Line: logFilterSample, Count: 9, Toolbar: true, State: &st})
		})
	}, 480, 320)
	st.Query = "nothing like this"
	tt.Frame()
	if err := tt.Click("Wrap lines"); err != nil {
		t.Fatal(err)
	}
	tt.Key(ui.Cmd, ui.KeyF)
	tt.Type("x")
	if !strings.Contains(st.Query, "x") {
		t.Fatalf("Cmd+F did not focus the search field: %q", st.Query)
	}
}
