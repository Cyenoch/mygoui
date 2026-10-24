package gallery

import (
	"fmt"
	"reflect"
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

// ComponentPage is one component's documentation page, shaped like a HeroUI
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

// mygouiImport is the module root's single-line import.
const mygouiImport = `import "` + mygouiImportPath + `"`

// RegisterComponent registers p as a gallery page. The page ID is the
// component ID, so the sidebar, search and tests address it directly.
func RegisterComponent(p ComponentPage) {
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

// buildComponentPage renders the docs sections every component page shares.
// The section order matches the HeroUI docs: usage, examples, the API table
// and the theme tokens the component consumes.
func buildComponentPage(c *ui.Context, s *State, p ComponentPage) ui.Element {
	chrome := docsChromeOf(s, p.ID)
	return ui.Column(c).FillWidth().MinWidth(0).Gap(theme.Scale(c, theme.Space32)).
		Children(func() {
			docsSection(c, chrome, "usage", "Usage", func() {
				usageBar(c, chrome, p)
			})
			docsSection(c, chrome, "examples", "Examples", func() {
				examplePanel(c, s, p)
			})
			docsSection(c, chrome, "api", "API Reference", func() {
				ui.Column(c).Gap(theme.Scale(c, theme.Space12)).Children(func() {
					mygoui.Typography(c, p.Title+" options", mygoui.TypographyOptions{
						Type: theme.TypeBodySM, Muted: true,
					})
					apiTable(c, p.Params)
				})
			})
			docsSection(c, chrome, "tokens", "Theme tokens", func() {
				tokenCards(c, p.Tokens)
			})
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
		Radius(theme.Scale(c, 20)).Background(t.Surface).
		TextColor(t.SurfaceForeground).BorderWidth(theme.Hairline).
		BorderColor(t.Border).Label(label)
	panel.Children(func() {
		ui.Column(c).FillWidth().MinWidth(0).Padding(theme.Scale(c, theme.Space40)).
			AlignItems(ui.Center).Children(body)
	})
	return panel
}

// exampleFrame renders one named demo: its title and description outside,
// then its own panel containing only the demo.
func exampleFrame(c *ui.Context, s *State, label, title, description string,
	demo func(c *ui.Context, s *State)) (ui.Element, ui.Element) {
	var panel ui.Element
	frame := ui.Column(c).FillWidth().MinWidth(0).Gap(theme.Scale(c, theme.Space12)).
		Children(func() {
			ui.Column(c).Gap(theme.Scale(c, theme.Space4)).Children(func() {
				mygoui.Typography(c, title, mygoui.TypographyOptions{Type: theme.TypeH5})
				if description != "" {
					mygoui.Typography(c, description, mygoui.TypographyOptions{
						Type: theme.TypeBodySM, Muted: true,
					})
				}
			})
			panel = demoPanel(c, "Example panel: "+label, func() {
				if demo != nil {
					demo(c, s)
				}
			})
		})
	return frame, panel
}

// examplePanel renders a page's examples and its shared code. Each
// ComponentExample gets its own panel. A single frame carries the code area
// under a divider; two or more frames share one code panel under a small
// "Code" label. The full snippet renders once.
func examplePanel(c *ui.Context, s *State, p ComponentPage) ui.Element {
	t := theme.Colors(c)
	var lastPanel ui.Element
	outer := ui.Column(c).FillWidth().MinWidth(0).Gap(theme.Scale(c, theme.Space40))
	outer.Children(func() {
		for _, example := range p.Examples {
			_, panel := exampleFrame(c, s, example.Title, example.Title,
				example.Description, example.Demo)
			lastPanel = panel
		}

		if strings.TrimSpace(p.Code) == "" {
			return
		}
		if len(p.Examples) <= 1 && lastPanel.Valid() {
			lastPanel.Children(func() {
				mygoui.Separator(c, mygoui.SeparatorOptions{})
				exampleCode(c, s, p)
			})
			return
		}
		ui.Column(c).FillWidth().MinWidth(0).Gap(theme.Scale(c, theme.Space12)).
			Children(func() {
				mygoui.Typography(c, "Code", mygoui.TypographyOptions{Type: theme.TypeH5})
				panel := ui.Column(c).FillWidth().MinWidth(0).Clip().
					Radius(theme.Scale(c, 20)).Background(t.Surface).
					TextColor(t.SurfaceForeground).BorderWidth(theme.Hairline).
					BorderColor(t.Border).Label("Example code panel")
				panel.Children(func() { exampleCode(c, s, p) })
			})
	})
	return outer
}

// panel is a HeroUI documentation surface: 16 DIP padding, 12 DIP radius,
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
// scroll horizontally. Long snippets fold to six lines under a fade.
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
	collapsed := long && !st.Expanded
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
							st.Expanded = true
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
								st.Expanded = false
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
	ui.Row(c).AlignItems(ui.Start).Gap(theme.Space12).Shrink(0).Children(func() {
		ui.Column(c).Width(theme.Scale(c, 24)).Shrink(0).AlignItems(ui.End).Children(func() {
			gutter := style.Text(ui.Text(c, strconv.Itoa(number)), c, ts).TextColor(t.Muted)
			gutter.Label("line " + strconv.Itoa(number)).TextAlign(ui.End).SingleLine()
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
	var spans []ui.Span
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
				add(word, t.Accent, false)
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
	Scroll   ui.ScrollState
	CopiedAt time.Time
	Expanded bool
}

// codeCopyFeedback is how long the copy button shows "Copied".
const codeCopyFeedback = 1500 * time.Millisecond

// codeCollapseLines is where a long snippet folds: six source lines stay
// visible until the reader expands it.
const codeCollapseLines = 6

// apiTable lists a component's parameters. Above 760 DIP of page width it
// keeps the Name/Type/Default/Description table; below it stacks one block
// per parameter so long values stay readable.
func apiTable(c *ui.Context, params []ComponentParam) {
	if len(params) == 0 {
		return
	}
	if pageWidth(c) < theme.Scale(c, 600) {
		apiStacked(c, params)
		return
	}
	panel(c, func() {
		ui.Column(c).FillWidth().MinWidth(0).Children(func() {
			apiRow(c, ComponentParam{Name: "Name", Type: "Type", Default: "Default", Description: "Description"}, true)
			for _, param := range params {
				mygoui.Separator(c, mygoui.SeparatorOptions{Color: mygoui.SeparatorSecondary})
				apiRow(c, param, false)
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
							codePill(c, param.Type, theme.Colors(c).Accent)
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

func apiRow(c *ui.Context, p ComponentParam, header bool) {
	// Wide columns: Name 20%, Type 30%, Default 20%, Description the rest.
	available := pageWidth(c) - 2*theme.Scale(c, theme.Space16)
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
			codePill(c, p.Type, theme.Colors(c).Accent)
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
	"Accent":                     "Primary action and focus emphasis.",
	"AccentForeground":           "Text and icons on an Accent fill.",
	"AccentHover":                "Accent fill under the pointer.",
	"AccentSoft":                 "Selected or softly highlighted surface.",
	"AccentSoftForeground":       "Text on an AccentSoft surface.",
	"Backdrop":                   "Dims the page behind a modal.",
	"Background":                 "Window background behind every surface.",
	"Border":                     "Outlines that separate surfaces and controls.",
	"Danger":                     "Destructive actions and invalid state.",
	"Default":                    "Neutral control surface, chips and keycaps.",
	"DefaultForeground":          "Text and icons on the neutral Default surface.",
	"DefaultHover":               "Neutral surface under the pointer.",
	"DefaultSoft":                "Quiet neutral fill for secondary controls.",
	"DefaultSoftForeground":      "Text on a DefaultSoft surface.",
	"DefaultSoftHover":           "DefaultSoft fill under the pointer.",
	"FieldBackground":            "Input and text-field fill.",
	"FieldBorder":                "Input outline at rest.",
	"FieldBorderHover":           "Input outline under the pointer.",
	"FieldBorderFocus":           "Input outline while focused.",
	"FieldFocus":                 "Input outline and caret while focused.",
	"FieldForeground":            "Text typed into a field.",
	"FieldHover":                 "Input fill under the pointer.",
	"FieldPlaceholder":           "Placeholder text inside an empty field.",
	"Focus":                      "Focus ring around the keyboard target.",
	"Foreground":                 "Primary body and heading text.",
	"Link":                       "Inline link text.",
	"Muted":                      "Secondary and supporting text.",
	"Overlay":                    "Popover, menu and tooltip surface.",
	"OverlayForeground":          "Text on the Overlay surface.",
	"Scrollbar":                  "Scrollbar thumb.",
	"Segment":                    "Selected segment in a control.",
	"SegmentForeground":          "Text on a Segment.",
	"Separator":                  "Divider rule between sections.",
	"Success":                    "Positive status and valid state.",
	"SuccessForeground":          "Text on a Success fill.",
	"SuccessSoft":                "Positive status fill at low emphasis.",
	"SuccessSoftForeground":      "Text on a SuccessSoft surface.",
	"Surface":                    "Base card and panel fill.",
	"SurfaceForeground":          "Text on a Surface.",
	"SurfaceSecondary":           "Secondary panel, section and input surface.",
	"SurfaceSecondaryForeground": "Text on a SurfaceSecondary surface.",
	"SurfaceShadows":             "Elevation shadow layers for cards and overlays.",
	"SurfaceTertiary":            "Tertiary, quietest surface.",
	"Warning":                    "Caution status and attention state.",
	"WarningForeground":          "Text on a Warning fill.",
	"WarningSoft":                "Caution status fill at low emphasis.",
	"WarningSoftForeground":      "Text on a WarningSoft surface.",
}

var paletteColorType = reflect.TypeOf(ui.Color{})

// paletteColor resolves a token name to the matching theme.Palette field, so
// the gallery shows the colors the library actually ships. It reports false
// for tokens that are not palette colors, such as shadows.
func paletteColor(p theme.Palette, token string) (ui.Color, bool) {
	field := reflect.ValueOf(p).FieldByName(token)
	if !field.IsValid() || field.Type() != paletteColorType {
		return ui.Color{}, false
	}
	return field.Interface().(ui.Color), true
}

// tokenCards shows a component's theme tokens as cards: the live swatch, its
// name and hex, the light and dark variants and what the token is for.
func tokenCards(c *ui.Context, tokens []string) {
	if len(tokens) == 0 {
		return
	}
	width := pageWidth(c)
	two := width >= theme.Scale(c, 600)
	gap := float32(theme.Space12)
	cardWidth := width
	if two {
		cardWidth = (width - gap) / 2
	}
	ui.Row(c).FillWidth().MinWidth(0).Wrap().Gap(gap).Children(func() {
		for _, token := range tokens {
			tokenCard(c, token, cardWidth, two)
		}
	})
}

// tokenCard is one token: a current-mode swatch with its mono name and hex,
// L/D variants and the token's role. Unknown tokens say so instead of
// inventing a color.
func tokenCard(c *ui.Context, token string, width float32, grid bool) {
	t := theme.Colors(c)
	light, lightOK := paletteColor(theme.ColorsForMode(false), token)
	dark, darkOK := paletteColor(theme.ColorsForMode(true), token)
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
				tokenMode(c, "L", light)
				tokenMode(c, "D", dark)
			})
		}
		mygoui.Typography(c, role, mygoui.TypographyOptions{
			Type: theme.TypeBodyXS, Muted: true,
		})
	})
}

// tokenMode is one light or dark swatch with its letter and hex value.
func tokenMode(c *ui.Context, label string, value ui.Color) {
	t := theme.Colors(c)
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
