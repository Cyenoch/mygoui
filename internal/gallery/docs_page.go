package gallery

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// docsEntry is one heading of a component page, as the page index lists it.
type docsEntry struct {
	ID, Title string
	Nested    bool
}

// apiSection is the reference of one type: its Go doc and one row per field.
// A section without a name holds the rows that belong to no listed type.
type apiSection struct {
	Name, Doc string
	Rows      []ComponentParam
}

// docsPageState is one component page's retained reading state: copy
// feedback, the derived reference and where its sections were last laid out.
type docsPageState struct {
	MarkdownCopiedAt time.Time
	InstallCopiedAt  time.Time
	ImportCopiedAt   time.Time
	Index            ui.ScrollState

	derived bool
	api     []apiSection
	entries []docsEntry
	imports string
	// tops is each section's top edge in the scrolled content, in DIPs. It
	// comes from the previous frame's layout: origin is where that frame put
	// the content's own top edge in the window, and offset how far this
	// build finds the content scrolled.
	tops           map[string]float32
	origin, offset float32
}

// exampleViewState is the Preview or Code choice of one example.
type exampleViewState struct{ Tab int }

func docsPageOf(s *State, p ComponentPage) *docsPageState {
	doc := StateOf[docsPageState](s, "docs:"+p.ID)
	if !doc.derived {
		doc.derived = true
		doc.api = p.apiSections()
		doc.entries = p.entries(doc.api)
		doc.imports = p.importStatement()
		doc.tops = map[string]float32{}
	}
	return doc
}

// codeTypeName matches a qualified exported name in a usage snippet.
var codeTypeName = regexp.MustCompile(`\b[a-z]+\.[A-Z]\w*`)

// paramNames splits a combined row name such as "Min, Max" into field names.
func paramNames(name string) []string {
	return strings.FieldsFunc(name, func(r rune) bool { return r == ',' || r == '/' || r == ' ' })
}

// apiSections derives the page's reference from the generated type table.
// It lists the struct the page describes most fields of among those its
// snippet names, then the named structs of the same family (MenuOptions,
// MenuItemOptions). A Params row adds its default and description to the
// first listed field it names; rows naming no field lead as their own table.
func (p ComponentPage) apiSections() []apiSection {
	described := map[string]ComponentParam{}
	for _, param := range p.Params {
		for _, name := range paramNames(param.Name) {
			described[name] = param
		}
	}
	var named []string
	primary, best := "", -1
	seen := map[string]bool{}
	for _, name := range codeTypeName.FindAllString(p.Code, -1) {
		typ, ok := apiTypes[name]
		if !ok || seen[name] || len(typ.Fields) == 0 {
			continue
		}
		seen[name] = true
		named = append(named, name)
		score := 0
		for _, field := range typ.Fields {
			if _, ok := described[field.Name]; ok {
				score++
			}
		}
		if score > best {
			primary, best = name, score
		}
	}
	var sections []apiSection
	used := map[string]bool{}
	if primary != "" {
		family := strings.TrimSuffix(primary, "Options")
		names := []string{primary}
		for _, name := range named {
			if name != primary && !strings.HasSuffix(family, ".") && strings.HasPrefix(name, family) {
				names = append(names, name)
			}
		}
		taken := map[string]bool{}
		for _, name := range names {
			typ := apiTypes[name]
			section := apiSection{Name: name, Doc: typ.Doc}
			for _, field := range typ.Fields {
				row := ComponentParam{Name: field.Name, Type: field.Type, Description: field.Doc}
				if param, ok := described[field.Name]; ok && !taken[field.Name] {
					taken[field.Name], used[param.Name] = true, true
					row.Description = param.Description
					if param.Name == field.Name {
						row.Default = param.Default
					}
				}
				section.Rows = append(section.Rows, row)
			}
			sections = append(sections, section)
		}
	}
	var rest []ComponentParam
	for _, param := range p.Params {
		if !used[param.Name] {
			rest = append(rest, param)
		}
	}
	if len(rest) > 0 {
		sections = append([]apiSection{{Rows: rest}}, sections...)
	}
	return sections
}

func (p ComponentPage) entries(api []apiSection) []docsEntry {
	entries := []docsEntry{{ID: "installation", Title: "Installation"}}
	if strings.TrimSpace(p.Code) != "" {
		entries = append(entries, docsEntry{ID: "usage", Title: "Usage"})
	}
	if len(api) > 0 {
		entries = append(entries, docsEntry{ID: "api", Title: "API Reference"})
		for _, section := range api {
			if section.Name != "" {
				entries = append(entries, docsEntry{ID: "api:" + section.Name, Title: section.Name, Nested: true})
			}
		}
	}
	if len(p.Examples) > 1 {
		entries = append(entries, docsEntry{ID: "examples", Title: "Examples"})
		for _, example := range p.Examples[1:] {
			entries = append(entries, docsEntry{ID: "example:" + example.Title, Title: example.Title, Nested: true})
		}
	}
	if len(p.Tokens) > 0 {
		entries = append(entries, docsEntry{ID: "tokens", Title: "Theme tokens"})
	}
	return entries
}

// buildComponentPage is one component's document: its lead example, how to
// install and use it, its reference, the remaining examples and its tokens.
func buildComponentPage(c *ui.Context, s *State, p ComponentPage) ui.Element {
	doc := docsPageOf(s, p)
	// Bounds are the previous frame's, laid out with the previous offset.
	doc.origin = s.Content.Bounds(c).Y - doc.offset
	doc.offset = s.Content.Y
	return ui.Column(c).FillWidth().MinWidth(0).Gap(theme.Scale(c, theme.Space40)).
		Children(func() {
			ui.Column(c).FillWidth().MinWidth(0).Gap(theme.Scale(c, theme.Space24)).Children(func() {
				docsActions(c, s, doc, p)
				if len(p.Examples) > 0 {
					exampleFrame(c, s, doc, p, p.Examples[0], true)
				}
			})
			for _, entry := range doc.entries {
				if entry.Nested {
					continue
				}
				docsSection(c, doc, entry, func() {
					switch entry.ID {
					case "installation":
						codeBlock(c, s, codeBlockOptions{
							Key: p.ID + ":install", Code: installCommand,
							Label: "Install command", CopyLabel: "Copy install command", Plain: true,
						})
					case "usage":
						codeBlock(c, s, codeBlockOptions{
							Key: p.ID + ":import", Code: doc.imports,
							Label: "Usage import", CopyLabel: "Copy import line",
						})
						codeBlock(c, s, codeBlockOptions{Key: p.ID, Code: p.Code})
					case "api":
						for _, section := range doc.api {
							apiReference(c, s, doc, section)
						}
					case "examples":
						for _, example := range p.Examples[1:] {
							exampleFrame(c, s, doc, p, example, false)
						}
					case "tokens":
						tokenCards(c, s, p.Tokens)
					}
				})
			}
		})
}

// docsActions are the page-level actions under the title.
func docsActions(c *ui.Context, s *State, doc *docsPageState, p ComponentPage) {
	ui.Row(c).FillWidth().Wrap().Gap(theme.Space8).AlignItems(ui.Center).Children(func() {
		if len(doc.api) > 0 {
			mygoui.Button(c, mygoui.ButtonOptions{
				Label: "API Reference", Icon: lucide.BookOpen(),
				Variant: mygoui.ButtonSecondary, Size: theme.ButtonSmall,
				OnClick: func() { scrollToSection(c, s, doc, "api") },
			})
		}
		docsCopyPill(c, &doc.MarkdownCopiedAt, "Copy Markdown", lucide.FileText(),
			func() string { return markdownDoc(p) })
		if url := sourceURL(p.ID); url != "" {
			services := c.Services()
			mygoui.Link(c, mygoui.LinkOptions{
				Label: referenceLabel(p.ID), Icon: lucide.ExternalLink(),
				OnClick: func() { services.OpenURL(url) },
			}).AlignSelf(ui.Center)
		}
	})
}

// docsSectionInset is the space kept above a section the reader jumps to.
const docsSectionInset = 16

// scrollToSection brings a section's heading to the top of the content.
func scrollToSection(c *ui.Context, s *State, doc *docsPageState, id string) {
	top, ok := doc.tops[id]
	if !ok {
		return
	}
	s.Content.Y = max(0, top-theme.Scale(c, docsSectionInset))
	c.Invalidate()
}

// docsSection is a headed part of the page. It records where the content
// holds it, which the page index reads to follow and move the reader.
func docsSection(c *ui.Context, doc *docsPageState, entry docsEntry, body func()) {
	section := ui.Column(c.Key("section:" + entry.ID)).FillWidth().MinWidth(0).
		Gap(theme.Scale(c, theme.Space12))
	if r := section.Bounds(); r.W > 0 {
		doc.tops[entry.ID] = r.Y - doc.origin
	} else {
		delete(doc.tops, entry.ID)
	}
	heading := theme.TypeH4
	if entry.Nested {
		heading = theme.TypeH5
	}
	section.Children(func() {
		mygoui.Typography(c, entry.Title, mygoui.TypographyOptions{Type: heading})
		body()
	})
}

// apiReference is one type of the reference: its name, Go doc and fields.
func apiReference(c *ui.Context, s *State, doc *docsPageState, section apiSection) {
	if section.Name == "" {
		apiTable(c, s, section.Rows)
		return
	}
	docsSection(c, doc, docsEntry{ID: "api:" + section.Name, Title: section.Name, Nested: true}, func() {
		if section.Doc != "" {
			mygoui.Typography(c, section.Doc, mygoui.TypographyOptions{Type: theme.TypeBodySM, Muted: true})
		}
		apiTable(c, s, section.Rows)
	})
}

// exampleFrame is one live example. One whose source the page files spell
// out shows it behind a Code tab. The lead example opens the page untitled;
// the others are indexed sections.
func exampleFrame(c *ui.Context, s *State, doc *docsPageState, p ComponentPage, example ComponentExample, lead bool) {
	view := StateOf[exampleViewState](s, "example:"+p.ID+":"+example.Title)
	code := exampleSource(p.ID, example.Title)
	body := func() {
		if example.Description != "" && !lead {
			mygoui.Typography(c, example.Description, mygoui.TypographyOptions{Type: theme.TypeBodySM, Muted: true})
		}
		if code != "" {
			viewTabs(c, &view.Tab, "Example view: "+example.Title)
			if view.Tab == 1 {
				codeBlock(c, s, codeBlockOptions{Key: p.ID + ":" + example.Title, Code: code, Numbered: true})
				return
			}
		}
		demoPanel(c, "Example panel: "+example.Title, func() {
			if example.Demo != nil {
				example.Demo(c, s)
			}
		})
	}
	if lead {
		ui.Column(c.Key(example.Title)).FillWidth().MinWidth(0).
			Gap(theme.Scale(c, theme.Space12)).Children(body)
		return
	}
	docsSection(c, doc, docsEntry{ID: "example:" + example.Title, Title: example.Title, Nested: true}, body)
}

var viewTabLabels = [...]string{"Preview", "Code"}

// viewTabs switches one example between its live preview and its source.
func viewTabs(c *ui.Context, selected *int, label string) {
	t := theme.Colors(c)
	parts := ui.SegmentedBase(c.Key("view"), selected, len(viewTabLabels))
	parts.Track.AlignSelf(ui.Start).MinWidth(0).MaxWidthPercent(100).Wrap().
		Padding(theme.Scale(c, 2)).Gap(theme.Scale(c, 2)).
		Radius(theme.RadiusScale().LG).Background(t.MutedBackground).Label(label).
		Children(func() {
			for i, text := range viewTabLabels {
				segment := parts.Segment(i).Shrink(0).Height(theme.Scale(c, 30)).
					PaddingX(theme.Scale(c, 10)).Radius(theme.RadiusScale().MD).
					FocusRing(true).Label(text)
				ink := t.MutedForeground.Alpha(0.72)
				if segment.Hovered() {
					ink = t.MutedForeground
				}
				if i == *selected {
					ink = t.Foreground
					face := t.Background
					if c.Theme().Dark {
						face = t.Input
					}
					segment = style.Shadow(segment.Background(face), c, theme.SurfaceShadows(c.Theme().Dark))
				}
				segment.Children(func() {
					style.Text(ui.Text(c, text), c, theme.SM()).
						FontWeight(theme.WeightMedium).TextColor(ink).SingleLine()
				})
			}
		})
}

// docsIndexReach is how far below the content's top edge a heading still
// counts as the section being read: just past where a jump leaves it.
const docsIndexReach = docsSectionInset + 8

// activeSection is the last section whose heading has passed the top of the
// content, or the last section once the content is scrolled to its end.
func activeSection(c *ui.Context, s *State, doc *docsPageState) string {
	if len(doc.entries) == 0 {
		return ""
	}
	if s.Content.MaxY > 0 && s.Content.Y >= s.Content.MaxY-1 {
		return doc.entries[len(doc.entries)-1].ID
	}
	line := s.Content.Y + theme.Scale(c, docsIndexReach)
	active := ""
	for _, entry := range doc.entries {
		if top, ok := doc.tops[entry.ID]; ok && top <= line {
			active = entry.ID
		}
	}
	return active
}

// docsIndex lists the page's sections beside the content, marks the one
// being read and scrolls to the one the reader picks.
func docsIndex(c *ui.Context, s *State, p ComponentPage) {
	doc := docsPageOf(s, p)
	active := activeSection(c, s, doc)
	pad := docsPadding(c)
	ui.Column(c.Key("index:"+p.ID)).Width(theme.Scale(c, docsIndexWidth)).Shrink(0).FillHeight().
		Padding(pad, pad, 0, 0).Gap(theme.Scale(c, theme.Space8)).Children(func() {
		mygoui.Typography(c, "On this page", mygoui.TypographyOptions{
			Type: theme.TypeBodyXS, Weight: theme.WeightMedium,
		})
		mygoui.ScrollArea(c, &doc.Index, mygoui.ScrollAreaOptions{
			Axis: mygoui.ScrollVertical, Label: "Page index",
		}, func() {
			ui.Column(c).FillWidth().MinWidth(0).Children(func() {
				for _, entry := range doc.entries {
					indexRow(c, entry, entry.ID == active, func() {
						scrollToSection(c, s, doc, entry.ID)
					})
				}
			})
		}).Grow(1).Basis(0).MinHeight(0)
	})
}

// indexRow is one index link; the section being read takes the foreground.
func indexRow(c *ui.Context, entry docsEntry, active bool, onClick func()) {
	t := theme.Colors(c)
	ts := theme.SM()
	start := float32(0)
	if entry.Nested {
		ts = theme.XS()
		start = theme.Scale(c, theme.Space12)
	}
	row := ui.ButtonBase(c.Key(entry.ID)).FillWidth().MinWidth(0).
		MinHeight(theme.Scale(c, 28)).Padding(0, 0, 0, start).
		Radius(theme.RadiusScale().SM).Justify(ui.Start).AlignItems(ui.Center).
		FocusRing(true).Shrink(0).Role(ui.RoleButton).Checked(active).
		Label("Section: " + entry.Title)
	ink := t.MutedForeground
	if active {
		ink = t.Foreground
		ts.Weight = theme.WeightMedium
		row.Description(locale.M().Selected)
	} else if row.Hovered() {
		ink = t.Foreground
	}
	row.Children(func() {
		style.Text(ui.Text(c, entry.Title), c, ts).TextColor(ink).SingleLine()
	})
	row.OnClick(onClick)
}

// markdownDoc renders a component page as a Markdown document for the
// "Copy Markdown" action.
func markdownDoc(p ComponentPage) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# %s\n\n%s\n", p.Title, p.Description)
	if url := sourceURL(p.ID); url != "" {
		fmt.Fprintf(&b, "\n[%s](%s)\n", referenceLabel(p.ID), url)
	} else {
		b.WriteString("\nNative addition.\n")
	}
	fmt.Fprintf(&b, "\n## Installation\n\n```sh\n%s\n```\n", installCommand)
	if code := strings.TrimSpace(p.Code); code != "" {
		fmt.Fprintf(&b, "\n## Usage\n\n```go\n%s\n\n%s\n```\n", p.importStatement(), code)
	}
	if api := p.apiSections(); len(api) > 0 {
		b.WriteString("\n## API Reference\n")
		for _, section := range api {
			if section.Name != "" {
				fmt.Fprintf(&b, "\n### %s\n", section.Name)
				if section.Doc != "" {
					fmt.Fprintf(&b, "\n%s\n", section.Doc)
				}
			}
			b.WriteString("\n| Prop | Type | Default | Description |\n| --- | --- | --- | --- |\n")
			for _, row := range section.Rows {
				fmt.Fprintf(&b, "| %s | %s | %s | %s |\n",
					mdCell(row.Name), mdCell(row.Type), mdCell(row.Default), mdCell(row.Description))
			}
		}
	}
	if len(p.Examples) > 1 {
		b.WriteString("\n## Examples\n")
		for _, example := range p.Examples[1:] {
			fmt.Fprintf(&b, "\n### %s\n", example.Title)
			if example.Description != "" {
				fmt.Fprintf(&b, "\n%s\n", example.Description)
			}
			if code := exampleSource(p.ID, example.Title); code != "" {
				fmt.Fprintf(&b, "\n```go\n%s\n```\n", code)
			}
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
