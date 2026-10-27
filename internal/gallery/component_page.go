package gallery

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/color"
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// ComponentParam is one usable parameter of a component's options.
type ComponentParam struct {
	Name, Type, Default, Description string
}

// ComponentExample is one named demo under a page's Examples section, such as
// "Variants" or "With Icons". Its title and description render outside its
// own bordered panel. Examples share the page's State and StateOf buckets, so
// a follow-up split must give each demo its own keys if their state must not
// interact. The page keeps one Code snippet today; per-example code is not
// modeled yet.
type ComponentExample struct {
	Title, Description string
	Demo               func(c *ui.Context, s *State)
}

// ComponentPage is one component's documentation page, shaped like a Coss
// docs page: live demos, a Go usage snippet, its usable parameters and the
// theme tokens it consumes. Examples split demos into named frames.
type ComponentPage struct {
	ID, Title, Group, Description string
	Order                         int
	Examples                      []ComponentExample
	Code                          string
	// Import is the package the usage bar names. Empty derives the import
	// block from the qualifiers in Code. A full "import ..." statement is
	// used as it is.
	Import string
	Params []ComponentParam
	Tokens []string
}

// codeImportQualifier matches a package qualifier in a usage snippet, so the
// bar can name the packages the snippet actually uses.
var codeImportQualifier = regexp.MustCompile(`\b(mygoui|form|list|navigation|overlay|date|color|data|theme|lucide)\.`)

// qualifierImports maps a snippet qualifier to its module path.
var qualifierImports = map[string]string{
	"mygoui":     mygouiImportPath,
	"form":       mygouiImportPath + "/form",
	"list":       mygouiImportPath + "/list",
	"navigation": mygouiImportPath + "/navigation",
	"overlay":    mygouiImportPath + "/overlay",
	"date":       mygouiImportPath + "/date",
	"color":      mygouiImportPath + "/color",
	"data":       mygouiImportPath + "/data",
	"theme":      mygouiImportPath + "/theme",
	"lucide":     "github.com/Cyenoch/mygoicons/lucide",
}

// importStatement is the import block shown in the Usage bar and copied by
// the header action. An explicit Import wins; otherwise the block names the
// module root and every package whose qualifier appears in Code. One path
// stays on one line, several become a parenthesized block.
func (p ComponentPage) importStatement() string {
	path := strings.TrimSpace(p.Import)
	if path != "" {
		if strings.HasPrefix(path, "import ") {
			return path
		}
		return `import ` + strconv.Quote(path)
	}
	seen := map[string]bool{mygouiImportPath: true}
	paths := []string{mygouiImportPath}
	for _, match := range codeImportQualifier.FindAllStringSubmatch(p.Code, -1) {
		importPath, ok := qualifierImports[match[1]]
		if !ok || seen[importPath] {
			continue
		}
		seen[importPath] = true
		paths = append(paths, importPath)
	}
	sort.Strings(paths)
	if len(paths) == 1 {
		return `import ` + strconv.Quote(paths[0])
	}
	var b strings.Builder
	b.WriteString("import (\n")
	for _, importPath := range paths {
		b.WriteString("\t")
		b.WriteString(strconv.Quote(importPath))
		b.WriteString("\n")
	}
	b.WriteString(")")
	return b.String()
}

var (
	componentMu   sync.RWMutex
	componentByID = map[string]ComponentPage{}
)

// componentPage returns the component documentation behind a page ID.
func componentPage(id string) (ComponentPage, bool) {
	componentMu.RLock()
	defer componentMu.RUnlock()
	p, ok := componentByID[id]
	return p, ok
}

// mygouiImportPath is the module root the usage bar always names.
const mygouiImportPath = "github.com/Cyenoch/mygoui"

// RegisterComponent registers p as a gallery page. The page ID is the
// component ID, so the sidebar, search and tests address it directly.
func RegisterComponent(p ComponentPage) {
	p.Tokens = cossTokens(p)
	componentMu.Lock()
	componentByID[p.ID] = p
	componentMu.Unlock()
	RegisterPage(Page{
		ID:          p.ID,
		Title:       p.Title,
		Description: p.Description,
		Group:       p.Group,
		Order:       p.Order,
		Build: func(c *ui.Context, s *State) ui.Element {
			return buildComponentPage(c, s, p)
		},
	})
}

// buildComponentPage shows one workbench view while keeping the page's
// examples, source and reference data together.
func buildComponentPage(c *ui.Context, s *State, p ComponentPage) ui.Element {
	chrome := docsChromeOf(s, p.ID)
	return ui.Column(c).FillWidth().MinWidth(0).Gap(theme.Scale(c, theme.Space16)).
		Children(func() {
			switch s.Panel {
			case 1:
				ui.Row(c).FillWidth().Wrap().Gap(theme.Space8).AlignItems(ui.Center).Children(func() {
					docsCopyPill(c, &chrome.MarkdownCopiedAt, "Copy Markdown", lucide.FileText(),
						func() string { return markdownDoc(p) })
					docsCopyPill(c, &chrome.ImportCopiedAt, "Copy import", lucide.Code(),
						func() string { return p.importStatement() })
					if url := sourceURL(p.ID); url != "" {
						services := c.Services()
						mygoui.Link(c, mygoui.LinkOptions{
							Label: "Coss source", Icon: lucide.ExternalLink(),
							OnClick: func() { services.OpenURL(url) },
						})
					}
				})
				usageBar(c, chrome, p)
				panel(c, func() { exampleCode(c, s, p) })
			case 2:
				apiTable(c, s, p.Params)
			case 3:
				tokenCards(c, s, p.Tokens)
			default:
				examplePanel(c, s, p)
			}
		})
}

func docsCopyPill(c *ui.Context, at *time.Time, label string, icon *ui.SVG, text func() string) {
	shown, glyph := label, icon
	if !at.IsZero() {
		elapsed := c.Now().Sub(*at)
		if elapsed < codeCopyFeedback {
			shown, glyph = "Copied", lucide.Check()
			c.After(codeCopyFeedback - elapsed)
		} else {
			*at = time.Time{}
		}
	}
	services := c.Services()
	mygoui.Button(c, mygoui.ButtonOptions{
		Label: shown, Icon: glyph, Variant: mygoui.ButtonSecondary,
		Size: theme.ButtonSmall,
		OnClick: func() {
			services.WriteClipboard(text())
			*at = c.Now()
			c.Invalidate()
		},
	})
}

// usageBar is the import block under Usage: one rounded bar, mono text and an
// icon-only copy control.
func usageBar(c *ui.Context, chrome *docsChromeState, p ComponentPage) ui.Element {
	t := theme.Colors(c)
	ts := theme.Code()
	ts.Size = 13
	e := ui.Row(c).FillWidth().MinWidth(0).MinHeight(theme.Scale(c, 44)).
		AlignItems(ui.Center).Gap(theme.Scale(c, theme.Space8)).
		Padding(0, theme.Scale(c, 10), 0, theme.Scale(c, theme.Space16)).
		Radius(theme.Scale(c, theme.Space12)).Clip().Background(t.SurfaceSecondary).
		TextColor(t.SurfaceSecondaryForeground).BorderWidth(theme.Hairline).
		BorderColor(t.Border).Label("Usage import")
	e.Children(func() {
		ui.Row(c).Grow(1).Basis(0).MinWidth(0).AlignItems(ui.Center).Children(func() {
			style.Text(ui.Text(c, p.importStatement()), c, ts).
				TextColor(t.Foreground).Selectable().NoWrap()
		})
		copyIconButton(c, &chrome.UsageCopiedAt, "Copy import line",
			func() string { return p.importStatement() })
	})
	return e
}

// demoPanel is the bordered, clipped surface around one demo body. Its title
// and description are rendered outside it by exampleFrame.
func demoPanel(c *ui.Context, label string, body func()) ui.Element {
	t := theme.Colors(c)
	panel := ui.Column(c).FillWidth().MinWidth(0).Clip().
		Radius(theme.RadiusScale().XL).Background(t.Card).
		TextColor(t.SurfaceForeground).BorderWidth(theme.Hairline).
		BorderColor(t.Border).Label(label)
	panel.Children(func() {
		ui.Column(c).FillWidth().MinWidth(0).Padding(theme.Scale(c, theme.Space24)).
			AlignItems(ui.Center).Children(body)
	})
	return panel
}

// exampleFrame keeps each live specimen in its own keyed surface.
func exampleFrame(c *ui.Context, s *State, example ComponentExample) ui.Element {
	return ui.Column(c.Key(example.Title)).FillWidth().MinWidth(0).Gap(theme.Scale(c, theme.Space8)).
		Children(func() {
			mygoui.Typography(c, example.Title, mygoui.TypographyOptions{Type: theme.TypeH5})
			demoPanel(c, "Example panel: "+example.Title, func() {
				if example.Demo != nil {
					example.Demo(c, s)
				}
			})
		})
}

// examplePanel renders all live examples, including off-screen scroll content.
func examplePanel(c *ui.Context, s *State, p ComponentPage) ui.Element {
	return ui.Column(c).FillWidth().MinWidth(0).Gap(theme.Scale(c, theme.Space24)).Children(func() {
		for _, example := range p.Examples {
			exampleFrame(c, s, example)
		}
	})
}

// panel is a Coss documentation surface: 16 DIP padding, 12 DIP radius,
// a hairline border and the page surface.
func panel(c *ui.Context, body func()) ui.Element {
	t := theme.Colors(c)
	return ui.Column(c).FillWidth().MinWidth(0).Padding(theme.Scale(c, theme.Space16)).
		Radius(theme.Scale(c, theme.Space12)).Clip().Background(t.Surface).
		TextColor(t.SurfaceForeground).BorderWidth(theme.Hairline).
		BorderColor(t.Border).Children(body)
}

// markdownDoc renders a component page as a Markdown document for the
// "Copy Markdown" action.
func markdownDoc(p ComponentPage) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# %s\n\n%s\n", p.Title, p.Description)
	if url := sourceURL(p.ID); url != "" {
		fmt.Fprintf(&b, "\n[Coss source](%s)\n", url)
	} else {
		b.WriteString("\nNative addition.\n")
	}
	if code := strings.TrimSpace(p.Code); code != "" {
		fmt.Fprintf(&b, "\n## Usage\n\n```go\n%s\n```\n", code)
	}
	if len(p.Params) > 0 {
		b.WriteString("\n## API Reference\n\n| Name | Type | Default | Description |\n")
		b.WriteString("| --- | --- | --- | --- |\n")
		for _, param := range p.Params {
			fmt.Fprintf(&b, "| %s | %s | %s | %s |\n",
				mdCell(param.Name), mdCell(param.Type), mdCell(param.Default), mdCell(param.Description))
		}
	}
	if len(p.Tokens) > 0 {
		b.WriteString("\n## Theme tokens\n\n")
		for _, token := range p.Tokens {
			role := tokenRoles[token]
			if role == "" {
				role = "Theme token."
			}
			fmt.Fprintf(&b, "- `%s` — %s\n", token, role)
		}
	}
	return b.String()
}

// mdCell keeps a Markdown table cell on one row.
func mdCell(text string) string {
	text = strings.ReplaceAll(text, "|", "\\|")
	return strings.ReplaceAll(text, "\n", " ")
}

// exampleCode is the code area of the Examples panel: a compact header with
// an icon-only copy control, a line-number gutter and highlighted lines that
// scroll horizontally. Long snippets can fold to six lines under a fade.
func exampleCode(c *ui.Context, s *State, p ComponentPage) {
	code := strings.TrimRight(p.Code, "\n")
	if strings.TrimSpace(code) == "" {
		return
	}
	t := theme.Colors(c)
	lines := strings.Split(code, "\n")
	long := len(lines) > codeCollapseLines
	st := StateOf[codeScrollState](s, "code:"+p.ID)
	visible := lines
	collapsed := long && st.Collapsed
	if collapsed {
		visible = lines[:codeCollapseLines]
	}
	ui.Column(c).FillWidth().MinWidth(0).Label("Go code").Children(func() {
		ui.Row(c).FillWidth().MinWidth(0).MinHeight(theme.Scale(c, 40)).
			AlignItems(ui.Center).Justify(ui.SpaceBetween).Gap(theme.Scale(c, theme.Space8)).
			Padding(0, theme.Scale(c, 10), 0, theme.Scale(c, theme.Space16)).
			BorderWidth(0, 0, theme.Hairline, 0).BorderColor(t.Border).
			Label("Go code header").Children(func() {
			mygoui.Typography(c, "Go", mygoui.TypographyOptions{
				Type: theme.TypeBodyXS, Weight: theme.WeightMedium,
			})
			copyIconButton(c, &st.CopiedAt, "Copy code", func() string { return p.Code })
		})
		body := ui.Column(c).FillWidth().MinWidth(0).Clip()
		body.Children(func() {
			bottom := theme.Scale(c, theme.Space16)
			if collapsed {
				bottom = theme.Scale(c, theme.Space32)
			}
			ui.ScrollHorizontal(c).TrackScroll(&st.Scroll).MinWidth(0).Children(func() {
				ui.Column(c).Shrink(0).
					Padding(theme.Scale(c, theme.Space16), theme.Scale(c, theme.Space16),
						bottom, theme.Scale(c, theme.Space16)).
					Children(func() {
						for i, line := range visible {
							codeRow(c, i+1, line)
						}
					})
			})
			if collapsed {
				fade := ui.Row(c).Absolute().Left(0).Right(0).Bottom(0).
					Height(theme.Scale(c, 56))
				fade.Draw(func(painter *ui.Painter, r ui.Rect) {
					clear := t.Surface
					clear.A = 0
					painter.FillGradient(r, ui.LinearGradient{
						Angle: 180,
						Stops: []ui.GradientStop{{Color: clear, At: 0}, {Color: t.Surface, At: 1}},
					}, 0)
				})
				ui.Row(c).Absolute().Left(0).Right(0).Bottom(theme.Scale(c, theme.Space6)).
					Justify(ui.Center).Children(func() {
					mygoui.Button(c, mygoui.ButtonOptions{
						Label: "Expand code", Icon: lucide.ChevronDown(),
						Variant: mygoui.ButtonSecondary, Size: theme.ButtonSmall,
						OnClick: func() {
							st.Collapsed = false
							c.Invalidate()
						},
					})
				})
			} else if long {
				ui.Row(c).FillWidth().MinWidth(0).Justify(ui.Center).
					Padding(theme.Space4, theme.Scale(c, theme.Space16),
						theme.Scale(c, theme.Space12), theme.Scale(c, theme.Space16)).
					Children(func() {
						mygoui.Button(c, mygoui.ButtonOptions{
							Label: "Collapse code", Icon: lucide.ChevronUp(),
							Variant: mygoui.ButtonGhost, Size: theme.ButtonSmall,
							OnClick: func() {
								st.Collapsed = true
								c.Invalidate()
							},
						})
					})
			}
		})
	})
}

// copyIconButton is an icon-only clipboard control with a tooltip. The
// accessible label names what it copies, so tests and assistive technology
// can tell several copy buttons apart.
func copyIconButton(c *ui.Context, at *time.Time, label string, text func() string) {
	shown, icon, tip := label, lucide.Copy(), "Copy"
	if !at.IsZero() {
		elapsed := c.Now().Sub(*at)
		if elapsed < codeCopyFeedback {
			shown, icon, tip = "Copied", lucide.Check(), "Copied"
			c.After(codeCopyFeedback - elapsed)
		} else {
			*at = time.Time{}
		}
	}
	mygoui.Button(c, mygoui.ButtonOptions{
		Label: shown, Icon: icon, IconOnly: true, Variant: mygoui.ButtonGhost,
		Size: theme.ButtonSmall, Tooltip: tip,
		OnClick: func() {
			c.WriteClipboard(text())
			*at = c.Now()
			c.Invalidate()
		},
	}).AlignSelf(ui.Center)
}

// codeRow is one source line: a right-aligned number in a fixed gutter, then
// the highlighted line.
func codeRow(c *ui.Context, number int, line string) {
	t := theme.Colors(c)
	ts := theme.Code()
	ts.Size = 13
	label := strconv.Itoa(number)
	ui.Row(c).AlignItems(ui.Start).Gap(theme.Space12).Shrink(0).Children(func() {
		ui.Column(c).Width(theme.Scale(c, 24)).Shrink(0).AlignItems(ui.End).Children(func() {
			gutter := style.Text(ui.Text(c, label), c, ts).TextColor(t.Muted)
			gutter.Label("line " + label).TextAlign(ui.End).SingleLine()
		})
		if line == "" {
			line = " "
		}
		ui.RichText(c, goSpans(c, line)...).
			Font(theme.Mono).FontSize(theme.Scale(c, 13)).
			FixedLineHeight(theme.Scale(c, 20)).Selectable().
			TextColor(t.Foreground)
	})
}

// goKeywords are the Go words the gallery snippets paint in the accent color.
var goKeywords = map[string]bool{
	"package": true, "import": true, "func": true, "type": true, "struct": true,
	"var": true, "const": true, "return": true, "if": true, "else": true,
	"for": true, "range": true, "nil": true, "true": true, "false": true,
}

// goSpans splits one line into Highlight spans: comments muted italic,
// strings success, numbers warning, keywords accent and everything else the
// foreground color.
func goSpans(c *ui.Context, line string) []ui.Span {
	t := theme.Colors(c)
	// A source line usually splits into a handful of tokens; sizing the slice
	// up front keeps a build from growing it three or four times per line.
	spans := make([]ui.Span, 0, 8)
	add := func(text string, color ui.Color, italic bool) {
		if text == "" {
			return
		}
		spans = append(spans, ui.Span{Text: text, Color: color, Italic: italic})
	}
	identStart := func(b byte) bool {
		return b == '_' || b >= 'a' && b <= 'z' || b >= 'A' && b <= 'Z'
	}
	identPart := func(b byte) bool {
		return identStart(b) || b >= '0' && b <= '9'
	}
	for i := 0; i < len(line); {
		if strings.HasPrefix(line[i:], "//") {
			add(line[i:], t.Muted, true)
			break
		}
		ch := line[i]
		switch {
		case ch == '"' || ch == '`':
			quote := ch
			j := i + 1
			for j < len(line) {
				if quote == '"' && line[j] == '\\' {
					j += 2
					continue
				}
				if line[j] == quote {
					j++
					break
				}
				j++
			}
			if j > len(line) {
				j = len(line)
			}
			add(line[i:j], t.Success, false)
			i = j
		case ch >= '0' && ch <= '9':
			j := i + 1
			for j < len(line) && (line[j] >= '0' && line[j] <= '9' || line[j] == '.') {
				j++
			}
			add(line[i:j], t.Warning, false)
			i = j
		case identStart(ch):
			j := i + 1
			for j < len(line) && identPart(line[j]) {
				j++
			}
			word := line[i:j]
			if goKeywords[word] {
				add(word, t.Primary, false)
			} else {
				add(word, t.Foreground, false)
			}
			i = j
		default:
			j := i + 1
			for j < len(line) && !identStart(line[j]) && (line[j] < '0' || line[j] > '9') &&
				line[j] != '"' && line[j] != '`' && !strings.HasPrefix(line[j:], "//") {
				j++
			}
			add(line[i:j], t.Foreground, false)
			i = j
		}
	}
	return spans
}

// codeScrollState is the example code area: its horizontal scroll so long
// lines do not truncate, its expanded/collapsed choice and the time its copy
// control last reported success, so feedback survives frames.
type codeScrollState struct {
	Scroll    ui.ScrollState
	CopiedAt  time.Time
	Collapsed bool
}

// codeCopyFeedback is how long the copy button shows "Copied".
const codeCopyFeedback = 1500 * time.Millisecond

// codeCollapseLines is where a long snippet folds: six source lines stay
// visible until the reader expands it.
const codeCollapseLines = 6

// apiTable lists a component's parameters. Above 760 DIP of page width it
// keeps the Name/Type/Default/Description table; below it stacks one block
// per parameter so long values stay readable.
func apiTable(c *ui.Context, s *State, params []ComponentParam) {
	if len(params) == 0 {
		return
	}
	width := pageWidth(c, s)
	if width < theme.Scale(c, 600) {
		apiStacked(c, params)
		return
	}
	// The panel's padding is the only thing between the page column and the
	// columns; measuring it once keeps every row from re-measuring the page.
	available := width - 2*theme.Scale(c, theme.Space16)
	panel(c, func() {
		ui.Column(c).FillWidth().MinWidth(0).Children(func() {
			apiRow(c, ComponentParam{Name: "Name", Type: "Type", Default: "Default", Description: "Description"}, true, available)
			for _, param := range params {
				mygoui.Separator(c, mygoui.SeparatorOptions{Color: mygoui.SeparatorSecondary})
				apiRow(c, param, false, available)
			}
		})
	})
}

// apiStacked is the narrow-width API: the name beside its wrapping code pills,
// then the description below, so long values stay readable without columns.
func apiStacked(c *ui.Context, params []ComponentParam) {
	panel(c, func() {
		ui.Column(c).FillWidth().MinWidth(0).Children(func() {
			for i, param := range params {
				if i > 0 {
					mygoui.Separator(c, mygoui.SeparatorOptions{Color: mygoui.SeparatorSecondary})
				}
				ui.Column(c).FillWidth().MinWidth(0).Gap(theme.Space6).
					Padding(theme.Space12, 0, theme.Space12, 0).Children(func() {
					ui.Row(c).FillWidth().MinWidth(0).Wrap().AlignItems(ui.Start).
						Gap(theme.Space6).Children(func() {
						apiParamName(c, param.Name)
						if param.Type != "" {
							codePill(c, param.Type, theme.Colors(c).Foreground)
						}
						if param.Default != "" {
							codePill(c, param.Default, defaultInk(c, param.Default))
						}
					})
					apiDescription(c, param.Description)
				})
			}
		})
	})
}

// apiParamName paints a parameter's name as compact mono medium ink.
func apiParamName(c *ui.Context, name string) {
	ts := theme.Code()
	ts.Weight = theme.WeightMedium
	style.Text(ui.Text(c, name), c, ts).TextColor(theme.Colors(c).Foreground)
}

// apiDescription paints a parameter's supporting copy.
func apiDescription(c *ui.Context, description string) {
	mygoui.Typography(c, description, mygoui.TypographyOptions{
		Type: theme.TypeBodySM, Muted: true,
	})
}

func apiRow(c *ui.Context, p ComponentParam, header bool, available float32) {
	// Wide columns: Name 20%, Type 30%, Default 20%, Description the rest.
	gaps := float32(3 * theme.Space12)
	nameW := available * 0.20
	typeW := available * 0.30
	defW := available * 0.20
	ui.Row(c).FillWidth().MinWidth(0).AlignItems(ui.Start).Gap(theme.Space12).
		Padding(theme.Space12, 0, theme.Space12, 0).Children(func() {
		ui.Column(c).Width(nameW).Shrink(0).MinWidth(0).Children(func() {
			if header {
				apiHeaderText(c, p.Name)
				return
			}
			apiParamName(c, p.Name)
		})
		ui.Column(c).Width(typeW).Shrink(0).MinWidth(0).Children(func() {
			if header {
				apiHeaderText(c, p.Type)
				return
			}
			codePill(c, p.Type, theme.Colors(c).Foreground)
		})
		ui.Column(c).Width(defW).Shrink(0).MinWidth(0).Children(func() {
			if header {
				apiHeaderText(c, p.Default)
				return
			}
			codePill(c, p.Default, defaultInk(c, p.Default))
		})
		ui.Column(c).Width(max(0, available-nameW-typeW-defW-gaps)).Shrink(0).MinWidth(0).
			Children(func() {
				if header {
					apiHeaderText(c, p.Description)
					return
				}
				apiDescription(c, p.Description)
			})
	})
}

func apiHeaderText(c *ui.Context, text string) {
	mygoui.Typography(c, text, mygoui.TypographyOptions{
		Type: theme.TypeBodyXS, Weight: theme.WeightMedium, Muted: true,
	})
}

// codePill is an inline-code chip: 12/16 mono on the Default background. It
// hugs its value, wraps onto more lines rather than ellipsizing when the
// column is narrow, stays selectable and carries the full value in a tooltip.
func codePill(c *ui.Context, text string, ink ui.Color) ui.Element {
	if text == "" {
		return ui.Row(c)
	}
	t := theme.Colors(c)
	ts := theme.XS()
	ts.Mono = true
	e := ui.Row(c).MinWidth(0).MaxWidthPercent(100).AlignSelf(ui.Start).
		AlignItems(ui.Center).Clip().
		Padding(theme.Space2, theme.Space6, theme.Space2, theme.Space6).
		Radius(theme.Scale(c, 6)).Background(t.Default)
	e.Children(func() {
		style.Text(ui.Text(c, text), c, ts).MaxWidthPercent(100).
			Selectable().TextColor(ink)
	})
	e.Tooltip(text)
	return e
}

// defaultInk paints literal defaults in the warning color and everything
// else in the foreground.
func defaultInk(c *ui.Context, value string) ui.Color {
	t := theme.Colors(c)
	switch value {
	case "true", "false", "nil", "":
		return t.Warning
	}
	if value[0] >= '0' && value[0] <= '9' || value[0] == '"' {
		return t.Warning
	}
	return t.Foreground
}

// tokenRoles explains what each theme token is for, so a token card is
// self-explanatory rather than a bare name.
var tokenRoles = map[string]string{
	"Primary":                    "Primary action fill; neutral by default.",
	"PrimaryForeground":          "Text and icons on a Primary fill.",
	"Secondary":                  "Low-emphasis control fill.",
	"SecondaryForeground":        "Text and icons on a Secondary fill.",
	"Card":                       "Card and panel background.",
	"CardForeground":             "Text and icons on a Card.",
	"Popover":                    "Floating menu and popover background.",
	"PopoverForeground":          "Text and icons inside a Popover.",
	"MutedBackground":            "Quiet background; CSS --muted.",
	"MutedForeground":            "Secondary text; CSS --muted-foreground.",
	"Destructive":                "Destructive action fill.",
	"DestructiveForeground":      "Destructive status text.",
	"Input":                      "Input outline at rest.",
	"Ring":                       "Keyboard focus ring.",
	"Chart1":                     "First chart series.",
	"Chart2":                     "Second chart series.",
	"Chart3":                     "Third chart series.",
	"Chart4":                     "Fourth chart series.",
	"Chart5":                     "Fifth chart series.",
	"Sidebar":                    "Sidebar background.",
	"SidebarForeground":          "Sidebar supporting text.",
	"SidebarPrimary":             "Sidebar primary action fill.",
	"SidebarPrimaryForeground":   "Text on a SidebarPrimary fill.",
	"SidebarAccent":              "Selected or hovered navigation row.",
	"SidebarAccentForeground":    "Text on a SidebarAccent fill.",
	"SidebarBorder":              "Sidebar boundary and separators.",
	"SidebarRing":                "Sidebar keyboard focus ring.",
	"Code":                       "Code block background.",
	"CodeForeground":             "Code block text.",
	"CodeHighlight":              "Highlighted source line.",
	"Accent":                     "Subtle hover and interaction fill.",
	"AccentForeground":           "Text and icons on an Accent fill.",
	"AccentHover":                "Accent fill under the pointer.",
	"AccentSoft":                 "Selected or softly highlighted surface.",
	"AccentSoftForeground":       "Text on an AccentSoft surface.",
	"AccentSoftHover":            "AccentSoft fill under the pointer.",
	"Backdrop":                   "Dims the page behind a modal.",
	"Background":                 "Window background behind every surface.",
	"BackgroundInverse":          "Canvas turned over for contrast, as tooltips use.",
	"BackgroundSecondary":        "Quiet alternative canvas, one step from Background.",
	"BackgroundTertiary":         "Third canvas level for grouped background regions.",
	"Border":                     "Outlines that separate surfaces and controls.",
	"BorderSecondary":            "Stronger surface outline for emphasis.",
	"BorderTertiary":             "Strongest surface outline.",
	"Danger":                     "Destructive actions and invalid state.",
	"DangerForeground":           "Text and icons on a Danger fill.",
	"DangerHover":                "Danger fill under the pointer.",
	"DangerSoft":                 "Destructive action at low emphasis.",
	"DangerSoftForeground":       "Text on a DangerSoft surface.",
	"DangerSoftHover":            "DangerSoft fill under the pointer.",
	"Default":                    "Neutral control surface, chips and keycaps.",
	"DefaultForeground":          "Text and icons on the neutral Default surface.",
	"DefaultHover":               "Neutral surface under the pointer.",
	"DefaultSoft":                "Quiet neutral fill for secondary controls.",
	"DefaultSoftForeground":      "Text on a DefaultSoft surface.",
	"DefaultSoftHover":           "DefaultSoft fill under the pointer.",
	"FieldBackground":            "Input and text-field fill.",
	"FieldBorder":                "Native derived input outline.",
	"FieldBorderHover":           "Input outline under the pointer.",
	"FieldBorderFocus":           "Input outline while focused.",
	"FieldFocus":                 "Input fill while focused.",
	"FieldForeground":            "Text typed into a field.",
	"FieldHover":                 "Input fill under the pointer.",
	"FieldPlaceholder":           "Placeholder text inside an empty field.",
	"Focus":                      "Ring around the focused control.",
	"Foreground":                 "Primary body and heading text.",
	"Info":                       "Informational status color.",
	"InfoForeground":             "Informational status text.",
	"InfoHover":                  "Info fill under the pointer.",
	"InfoSoft":                   "Informational status at low emphasis.",
	"InfoSoftForeground":         "Text on an InfoSoft surface.",
	"InfoSoftHover":              "InfoSoft fill under the pointer.",
	"Link":                       "Inline link text.",
	"Muted":                      "Secondary and supporting text.",
	"Overlay":                    "Popover, menu and tooltip surface.",
	"OverlayForeground":          "Text on the Overlay surface.",
	"Scrollbar":                  "Scrollbar thumb.",
	"Segment":                    "Selected segment in a control.",
	"SegmentForeground":          "Text on a Segment.",
	"Selection":                  "Selected row and selected-text fill.",
	"Separator":                  "Divider rule between sections.",
	"SeparatorSecondary":         "Stronger divider rule.",
	"SeparatorTertiary":          "Strongest divider rule.",
	"Success":                    "Positive status and valid state.",
	"SuccessForeground":          "Positive status text.",
	"SuccessHover":               "Success fill under the pointer.",
	"SuccessSoft":                "Positive status fill at low emphasis.",
	"SuccessSoftForeground":      "Text on a SuccessSoft surface.",
	"SuccessSoftHover":           "SuccessSoft fill under the pointer.",
	"Surface":                    "Base card and panel fill.",
	"SurfaceForeground":          "Text on a Surface.",
	"SurfaceHover":               "Surface fill under the pointer.",
	"SurfaceSecondary":           "Secondary panel, section and input surface.",
	"SurfaceSecondaryForeground": "Text on a SurfaceSecondary surface.",
	"SurfaceShadows":             "Elevation shadow layers for cards and overlays.",
	"SurfaceTertiary":            "Tertiary, quietest surface.",
	"SurfaceTertiaryForeground":  "Text on a SurfaceTertiary surface.",
	"Warning":                    "Caution status and attention state.",
	"WarningForeground":          "Caution status text.",
	"WarningHover":               "Warning fill under the pointer.",
	"WarningSoft":                "Caution status fill at low emphasis.",
	"WarningSoftForeground":      "Text on a WarningSoft surface.",
	"WarningSoftHover":           "WarningSoft fill under the pointer.",
}

// paletteColor resolves a token name to the matching theme.Palette field, so
// the gallery shows the colors the library actually ships. It reports false
// for tokens that are not palette colors, such as shadows. The name switch is
// one hash lookup per card build; reflection here allocated and scanned every
// field for each of tens of cards on every frame.
func paletteColor(p *theme.Palette, token string) (ui.Color, bool) {
	switch token {
	case "Primary":
		return p.Primary, true
	case "PrimaryForeground":
		return p.PrimaryForeground, true
	case "Secondary":
		return p.Secondary, true
	case "SecondaryForeground":
		return p.SecondaryForeground, true
	case "Card":
		return p.Card, true
	case "CardForeground":
		return p.CardForeground, true
	case "Popover":
		return p.Popover, true
	case "PopoverForeground":
		return p.PopoverForeground, true
	case "MutedBackground":
		return p.MutedBackground, true
	case "MutedForeground":
		return p.MutedForeground, true
	case "Destructive":
		return p.Destructive, true
	case "DestructiveForeground":
		return p.DestructiveForeground, true
	case "Input":
		return p.Input, true
	case "Ring":
		return p.Ring, true
	case "Chart1":
		return p.Chart1, true
	case "Chart2":
		return p.Chart2, true
	case "Chart3":
		return p.Chart3, true
	case "Chart4":
		return p.Chart4, true
	case "Chart5":
		return p.Chart5, true
	case "Sidebar":
		return p.Sidebar, true
	case "SidebarForeground":
		return p.SidebarForeground, true
	case "SidebarPrimary":
		return p.SidebarPrimary, true
	case "SidebarPrimaryForeground":
		return p.SidebarPrimaryForeground, true
	case "SidebarAccent":
		return p.SidebarAccent, true
	case "SidebarAccentForeground":
		return p.SidebarAccentForeground, true
	case "SidebarBorder":
		return p.SidebarBorder, true
	case "SidebarRing":
		return p.SidebarRing, true
	case "Code":
		return p.Code, true
	case "CodeForeground":
		return p.CodeForeground, true
	case "CodeHighlight":
		return p.CodeHighlight, true
	case "Accent":
		return p.Accent, true
	case "AccentForeground":
		return p.AccentForeground, true
	case "AccentHover":
		return p.AccentHover, true
	case "AccentSoft":
		return p.AccentSoft, true
	case "AccentSoftForeground":
		return p.AccentSoftForeground, true
	case "AccentSoftHover":
		return p.AccentSoftHover, true
	case "Backdrop":
		return p.Backdrop, true
	case "Background":
		return p.Background, true
	case "BackgroundInverse":
		return p.BackgroundInverse, true
	case "BackgroundSecondary":
		return p.BackgroundSecondary, true
	case "BackgroundTertiary":
		return p.BackgroundTertiary, true
	case "Border":
		return p.Border, true
	case "BorderSecondary":
		return p.BorderSecondary, true
	case "BorderTertiary":
		return p.BorderTertiary, true
	case "Danger":
		return p.Danger, true
	case "DangerForeground":
		return p.DangerForeground, true
	case "DangerHover":
		return p.DangerHover, true
	case "DangerSoft":
		return p.DangerSoft, true
	case "DangerSoftForeground":
		return p.DangerSoftForeground, true
	case "DangerSoftHover":
		return p.DangerSoftHover, true
	case "Default":
		return p.Default, true
	case "DefaultForeground":
		return p.DefaultForeground, true
	case "DefaultHover":
		return p.DefaultHover, true
	case "DefaultSoft":
		return p.DefaultSoft, true
	case "DefaultSoftForeground":
		return p.DefaultSoftForeground, true
	case "DefaultSoftHover":
		return p.DefaultSoftHover, true
	case "FieldBackground":
		return p.FieldBackground, true
	case "FieldBorder":
		return p.FieldBorder, true
	case "FieldBorderFocus":
		return p.FieldBorderFocus, true
	case "FieldBorderHover":
		return p.FieldBorderHover, true
	case "FieldFocus":
		return p.FieldFocus, true
	case "FieldForeground":
		return p.FieldForeground, true
	case "FieldHover":
		return p.FieldHover, true
	case "FieldPlaceholder":
		return p.FieldPlaceholder, true
	case "Focus":
		return p.Focus, true
	case "Foreground":
		return p.Foreground, true
	case "Info":
		return p.Info, true
	case "InfoForeground":
		return p.InfoForeground, true
	case "InfoHover":
		return p.InfoHover, true
	case "InfoSoft":
		return p.InfoSoft, true
	case "InfoSoftForeground":
		return p.InfoSoftForeground, true
	case "InfoSoftHover":
		return p.InfoSoftHover, true
	case "Link":
		return p.Link, true
	case "Muted":
		return p.Muted, true
	case "Overlay":
		return p.Overlay, true
	case "OverlayForeground":
		return p.OverlayForeground, true
	case "Scrollbar":
		return p.Scrollbar, true
	case "Segment":
		return p.Segment, true
	case "SegmentForeground":
		return p.SegmentForeground, true
	case "Selection":
		return p.Selection, true
	case "Separator":
		return p.Separator, true
	case "SeparatorSecondary":
		return p.SeparatorSecondary, true
	case "SeparatorTertiary":
		return p.SeparatorTertiary, true
	case "Success":
		return p.Success, true
	case "SuccessForeground":
		return p.SuccessForeground, true
	case "SuccessHover":
		return p.SuccessHover, true
	case "SuccessSoft":
		return p.SuccessSoft, true
	case "SuccessSoftForeground":
		return p.SuccessSoftForeground, true
	case "SuccessSoftHover":
		return p.SuccessSoftHover, true
	case "Surface":
		return p.Surface, true
	case "SurfaceForeground":
		return p.SurfaceForeground, true
	case "SurfaceHover":
		return p.SurfaceHover, true
	case "SurfaceSecondary":
		return p.SurfaceSecondary, true
	case "SurfaceSecondaryForeground":
		return p.SurfaceSecondaryForeground, true
	case "SurfaceTertiary":
		return p.SurfaceTertiary, true
	case "SurfaceTertiaryForeground":
		return p.SurfaceTertiaryForeground, true
	case "Warning":
		return p.Warning, true
	case "WarningForeground":
		return p.WarningForeground, true
	case "WarningHover":
		return p.WarningHover, true
	case "WarningSoft":
		return p.WarningSoft, true
	case "WarningSoftForeground":
		return p.WarningSoftForeground, true
	case "WarningSoftHover":
		return p.WarningSoftHover, true
	}
	return ui.Color{}, false
}

// tokenCards shows a component's theme tokens as cards: the live swatch, its
// name and hex, the light and dark variants and what the token is for. The
// three palettes are resolved once for every card in the section.
func tokenCards(c *ui.Context, s *State, tokens []string) {
	if len(tokens) == 0 {
		return
	}
	width := pageWidth(c, s)
	two := width >= theme.Scale(c, 600)
	gap := float32(theme.Space12)
	cardWidth := width
	if two {
		cardWidth = (width - gap) / 2
	}
	current, light, dark := theme.Colors(c), theme.ColorsForMode(false), theme.ColorsForMode(true)
	ui.Row(c).FillWidth().MinWidth(0).Wrap().Gap(gap).Children(func() {
		for _, token := range tokens {
			tokenCard(c, &current, &light, &dark, token, cardWidth, two)
		}
	})
}

// tokenCard is one token: a current-mode swatch with its mono name and hex,
// L/D variants and the token's role. Unknown tokens say so instead of
// inventing a color.
func tokenCard(c *ui.Context, t, lightPalette, darkPalette *theme.Palette, token string, width float32, grid bool) {
	light, lightOK := paletteColor(lightPalette, token)
	dark, darkOK := paletteColor(darkPalette, token)
	current, currentOK := paletteColor(t, token)
	known := lightOK && darkOK && currentOK

	role := tokenRoles[token]
	if role == "" {
		if known {
			role = "Theme palette color."
		} else {
			role = "Custom token supplied by the app."
		}
	}
	tip := token + ": " + role
	if known {
		tip += " (L " + color.FormatHex(light) + ", D " + color.FormatHex(dark) + ")"
	}

	card := ui.Column(c).Gap(theme.Space8).MinWidth(0).
		Padding(theme.Scale(c, theme.Space12)).Radius(theme.Scale(c, theme.Space12)).
		Clip().Background(t.SurfaceSecondary).TextColor(t.SurfaceSecondaryForeground).
		BorderWidth(theme.Hairline).BorderColor(t.Border).
		Label("token card " + token).Tooltip(tip)
	if grid {
		card = card.Width(width).Shrink(0)
	} else {
		card = card.FillWidth()
	}
	card.Children(func() {
		ui.Row(c).FillWidth().MinWidth(0).AlignItems(ui.Center).
			Gap(theme.Space8).Children(func() {
			if known {
				color.ColorSwatch(c, color.ColorSwatchOptions{
					Color: current, Size: color.SwatchLG, Shape: color.SwatchSquare,
					Label: token + " current swatch",
				})
			}
			ui.Column(c).Grow(1).Basis(0).MinWidth(0).Gap(theme.Space2).Children(func() {
				ts := theme.XS()
				ts.Mono = true
				ts.Weight = theme.WeightMedium
				style.Text(ui.Text(c, token), c, ts).TextColor(t.Foreground)
				if known {
					ts := theme.XS()
					ts.Mono = true
					style.Text(ui.Text(c, color.FormatHex(current)), c, ts).
						TextColor(t.Muted).Selectable()
				} else {
					mygoui.Typography(c, "custom token", mygoui.TypographyOptions{
						Type: theme.TypeBodyXS, Muted: true,
					})
				}
			})
		})
		if known {
			ui.Row(c).FillWidth().MinWidth(0).Wrap().Gap(theme.Space12).
				AlignItems(ui.Center).Children(func() {
				tokenMode(c, t, "L", light)
				tokenMode(c, t, "D", dark)
			})
		}
		mygoui.Typography(c, role, mygoui.TypographyOptions{
			Type: theme.TypeBodyXS, Muted: true,
		})
	})
}

// tokenMode is one light or dark swatch with its letter and hex value.
func tokenMode(c *ui.Context, t *theme.Palette, label string, value ui.Color) {
	ui.Row(c).Gap(theme.Space4).AlignItems(ui.Center).Shrink(0).Children(func() {
		ts := theme.XS()
		ts.Mono = true
		ts.Weight = theme.WeightMedium
		style.Text(ui.Text(c, label), c, ts).TextColor(t.Muted)
		color.ColorSwatch(c, color.ColorSwatchOptions{
			Color: value, Size: color.SwatchXS, Shape: color.SwatchSquare,
			Label: label + " mode swatch",
		})
		ts = theme.XS()
		ts.Mono = true
		style.Text(ui.Text(c, color.FormatHex(value)), c, ts).
			TextColor(t.Muted).Selectable()
	})
}
