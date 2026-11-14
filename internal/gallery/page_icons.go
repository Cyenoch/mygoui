package gallery

import (
	"fmt"
	"strings"

	"github.com/Cyenoch/mygoicons/catalog"
	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/form"
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/overlay"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

const (
	iconCellMinWidth = 96
	iconCellHeight   = 76
	iconListHeight   = 380
)

// iconsState is the icon browser: the search the reader typed and the
// matches derived from it. It lists Lucide, the collection the library
// draws with. Twemoji would add 16 MB to the browser gallery and its own
// attribution to every distributed copy.
type iconsState struct {
	set   catalog.Collection
	query string

	// filtered names what results was derived from.
	filtered string
	results  []catalog.Icon
	width    float32
	list     ui.ListState
	// built counts the rows of the last pass, for the virtualization test.
	built int
}

func iconsStateOf(s *State) *iconsState {
	st := StateOf[iconsState](s, "icons")
	if st.set.Prefix == "" {
		st.set = catalog.Lucide()
	}
	return st
}

// filter derives the drawings matching the search, once per query. An alias matches by its
// own name and ID and yields the drawing it repeats, once.
func (st *iconsState) filter() {
	set := st.set
	if st.results != nil && st.query == st.filtered {
		return
	}
	st.filtered = st.query
	if st.results == nil {
		st.results = make([]catalog.Icon, 0, len(set.Icons))
	}
	st.results = st.results[:0]
	query := strings.ToLower(strings.TrimSpace(st.query))
	byName := make(map[string]catalog.Icon, len(set.Icons))
	for _, icon := range set.Icons {
		byName[icon.Name] = icon
	}
	shown := make(map[string]bool, len(set.Icons))
	for _, icon := range set.Icons {
		if !strings.Contains(strings.ToLower(icon.Name), query) && !strings.Contains(icon.ID, query) {
			continue
		}
		if icon.AliasOf != "" {
			icon = byName[icon.AliasOf]
		}
		if !shown[icon.Name] {
			shown[icon.Name] = true
			st.results = append(st.results, icon)
		}
	}
	st.list.ScrollTo(0, ui.Start)
}

// iconBrowser searches a collection and copies an icon's Go expression.
func iconBrowser(c *ui.Context, s *State, st *iconsState) {
	st.filter()
	set := st.set
	t := theme.Colors(c)
	ui.Column(c).FillWidth().MinWidth(0).Gap(theme.Space12).Children(func() {
		form.SearchField(c, &st.query, form.SearchFieldOptions{
			AriaLabel: "Search icons", Placeholder: "Go name or Iconify ID", FullWidth: true,
		})
		mygoui.Typography(c, fmt.Sprintf("%d icons · %s %s · %s", len(st.results), set.Title, set.Version, set.License),
			mygoui.TypographyOptions{Type: theme.TypeBodyXS, Muted: true})
		if len(st.results) == 0 {
			mygoui.EmptyState(c, mygoui.EmptyStateOptions{
				Icon: lucide.SearchX(), Title: "No icons found",
			})
			return
		}
		gap := theme.Scale(c, theme.Space4)
		columns := max(1, int((st.width+gap)/(theme.Scale(c, iconCellMinWidth)+gap)))
		rows := (len(st.results) + columns - 1) / columns
		st.built = 0
		st.list.Key = func(row int) any { return st.results[row*columns].Name }
		ui.List(c, &st.list, rows, func(row int) {
			st.built++
			ui.Row(c).FillWidth().Height(theme.Scale(c, iconCellHeight)).Gap(gap).Children(func() {
				for column := range columns {
					index := row*columns + column
					if index >= len(st.results) {
						ui.Box(c).Grow(1).Basis(0)
						continue
					}
					iconCell(c, s, set, st.results[index])
				}
			})
		}).FillWidth().Height(theme.Scale(c, iconListHeight)).Label("Icons").
			Border(theme.Hairline, t.Border).Radius(theme.RadiusScale().LG).
			DrawOver(func(_ *ui.Painter, r ui.Rect) {
				if r.W != st.width {
					st.width = r.W
					c.Invalidate()
				}
			})
	})
}

// iconCell is one drawing with its Go name. Activating it copies the
// expression that returns it.
func iconCell(c *ui.Context, s *State, set catalog.Collection, icon catalog.Icon) {
	t := theme.Colors(c)
	radius := theme.RadiusScale().MD
	expression := set.Prefix + "." + icon.Name + "()"
	cell := ui.ButtonBase(c).Grow(1).Basis(0).MinWidth(0).FillHeight().Radius(radius).
		Padding(theme.Scale(c, theme.Space4)).Gap(theme.Scale(c, theme.Space8)).
		Justify(ui.Center).AlignItems(ui.Center).
		Label(icon.Name).Description(set.Prefix + ":" + icon.ID).TextColor(t.Foreground)
	cell = style.Background(cell, radius, ui.Transparent, t.Accent, t.Accent)
	cell.Children(func() {
		ui.Column(c).FillWidth().MinWidth(0).Gap(theme.Scale(c, theme.Space8)).AlignItems(ui.Center).Children(func() {
			ui.Icon(c, icon.SVG()).Size(theme.Scale(c, 20), theme.Scale(c, 20))
			style.Text(ui.Text(c, icon.Name), c, theme.XS()).TextColor(t.Muted).SingleLine().MaxWidth(theme.Scale(c, iconCellMinWidth-8))
		})
	})
	cell.OnClick(func() {
		c.WriteClipboard(expression)
		overlay.Toast(c, &overlaysStateOf(s).toasts, overlay.ToastOptions{
			ID: "icon-copied", Title: "Copied " + expression,
		})
	})
}

// iconSizesExample shows that an icon takes its size from the element and
// its color from the text around it.
func iconSizesExample(c *ui.Context) {
	t := theme.Colors(c)
	ui.Row(c).Wrap().Gap(theme.Space24).AlignItems(ui.Center).Children(func() {
		for _, size := range []float32{14, 16, 20, 24, 32} {
			ui.Icon(c, lucide.Sparkles()).Size(theme.Scale(c, size), theme.Scale(c, size)).
				Label(fmt.Sprintf("Sparkles at %.0f", size))
		}
		for _, ink := range []ui.Color{t.Muted, t.Primary, t.Destructive} {
			ui.Row(c).TextColor(ink).Children(func() {
				ui.Icon(c, lucide.Heart()).Size(theme.Scale(c, 24), theme.Scale(c, 24)).Label("Heart")
			})
		}
	})
}

// iconComponentsExample passes icons to the components that take one.
func iconComponentsExample(c *ui.Context) {
	ui.Row(c).Wrap().Gap(theme.Space12).AlignItems(ui.Center).Children(func() {
		mygoui.Button(c, mygoui.ButtonOptions{Label: "Download", Icon: lucide.Download()})
		mygoui.Button(c, mygoui.ButtonOptions{
			Label: "Settings", Icon: lucide.Settings(), IconOnly: true, Variant: mygoui.ButtonOutline,
		})
		mygoui.Button(c, mygoui.ButtonOptions{
			Label: "Delete", Icon: lucide.Trash2(), Variant: mygoui.ButtonDanger,
		})
	})
}

func init() {
	RegisterComponent(ComponentPage{
		ID: "icons", Title: "Icons", Group: "Foundations", Order: 15,
		Description: "Lucide icons from mygoicons, as typed Go functions.",
		Examples: []ComponentExample{
			{
				Title:       "Browse",
				Description: "Search by Go name or Iconify ID; activating an icon copies its expression.",
				Demo:        func(c *ui.Context, s *State) { iconBrowser(c, s, iconsStateOf(s)) },
			},
			{
				Title:       "Size and color",
				Description: "The element sets the size; the text color sets the stroke.",
				Demo:        func(c *ui.Context, s *State) { iconSizesExample(c) },
			},
			{
				Title:       "In components",
				Description: "Components with an Icon option take the same value.",
				Demo:        func(c *ui.Context, s *State) { iconComponentsExample(c) },
			},
		},
		Import: `import (
    "github.com/Cyenoch/mygoicons/lucide"
    "github.com/Cyenoch/mygoui"
    "github.com/egoist/mygo/ui"
)`,
		Code: `ui.Icon(c, lucide.Search()).Size(16, 16).Label("Search")

mygoui.Button(c, mygoui.ButtonOptions{
    Label: "Download", Icon: lucide.Download(),
})`,
		Params: []ComponentParam{
			{"Icon", "*ui.SVG", "nil", "Any mygoicons accessor, such as lucide.Download()."},
			{"IconOnly", "bool", "false", "Shows the icon alone; Label stays the accessible name."},
		},
		Tokens: []string{"Foreground", "Muted", "Accent"},
	})
}
