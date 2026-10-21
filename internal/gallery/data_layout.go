package gallery

import (
	"strconv"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/data"
	"github.com/Cyenoch/mygoui/layout"
	"github.com/egoist/mygo/ui"
)

type DataLayoutState struct {
	Split                                               layout.ResizableState
	Scroll                                              ui.ScrollState
	Carousel                                            layout.CarouselState
	Table                                               data.TableState
	DataTable                                           data.DataTableState
	Value                                               float64
	Indeterminate, Destructive, Vertical, Disabled, RTL bool
	Revision                                            uint64
	Rows                                                []data.Row
	CellSelection, DragSelection                        bool
}

func NewDataLayoutState() DataLayoutState {
	return DataLayoutState{Value: 40, Rows: append([]data.Row(nil), sampleRows...), DataTable: data.DataTableState{PageSize: 5}}
}

var sampleColumns = []data.Column{{ID: "name", Title: "Name", Width: 180, MinWidth: 120}, {ID: "status", Title: "Status", Width: 120, MinWidth: 90}, {ID: "count", Title: "Files", Width: 80, MinWidth: 60, Align: ui.End, Compare: func(a, b data.Row) int {
	left, _ := strconv.Atoi(a.Cells[2])
	right, _ := strconv.Atoi(b.Cells[2])
	return left - right
}}}
var sampleRows = []data.Row{{ID: "core", Cells: []string{"Core library", "Ready", "24"}}, {ID: "native", Cells: []string{"Native UI", "In progress", "18"}}, {ID: "docs", Cells: []string{"Documentation", "Ready", "12"}}, {ID: "tests", Cells: []string{"Behavior tests", "Ready", "32"}}, {ID: "examples", Cells: []string{"Examples", "In progress", "16"}}, {ID: "assets", Cells: []string{"Assets", "Ready", "8"}}, {ID: "release", Cells: []string{"Release", "Planned", "6"}}}

func dataLayoutSample(c *ui.Context, s *DataLayoutState, id ComponentID, advanced bool) {
	ui.Column(c).FillWidth().MinWidth(0).Gap(12).Children(func() {
		if advanced {
			ui.Row(c).FillWidth().Wrap().Gap(8).Children(func() {
				mygoui.Toggle(c, &s.Disabled, mygoui.ToggleOptions{Label: "Disabled"})
				switch id {
				case "resizable", "carousel":
					mygoui.Toggle(c, &s.Vertical, mygoui.ToggleOptions{Label: "Vertical"})
				case "alert":
					mygoui.Toggle(c, &s.Destructive, mygoui.ToggleOptions{Label: "Destructive"})
				case "progress":
					mygoui.Toggle(c, &s.Indeterminate, mygoui.ToggleOptions{Label: "Indeterminate"})
				case "table", "data-table":
					mygoui.Toggle(c, &s.CellSelection, mygoui.ToggleOptions{Label: "Cell selection"})
					mygoui.Toggle(c, &s.DragSelection, mygoui.ToggleOptions{Label: "Drag selection"})
				}
				mygoui.Button(c, mygoui.ButtonOptions{Label: "Reset data sample", Variant: mygoui.ButtonOutline, OnClick: func() { *s = NewDataLayoutState() }})
				if id == "data-table" || id == "table" {
					mygoui.Button(c, mygoui.ButtonOptions{Label: "Append local row", Variant: mygoui.ButtonOutline, OnClick: func() {
						s.Rows = append(s.Rows, data.Row{ID: data.RowID("sample-" + strconv.Itoa(len(s.Rows))), Cells: []string{"Added sample", "Ready", "1"}})
						s.Revision++
					}})
				}
			})
		}
		switch id {
		case "alert":
			variant := mygoui.AlertDefault
			if s.Destructive {
				variant = mygoui.AlertDestructive
			}
			mygoui.Alert(c, mygoui.AlertOptions{Title: "Keep a local copy", Description: "Your sample changes stay in this window.", Variant: variant, Action: &mygoui.ButtonOptions{Label: "Review", Variant: mygoui.ButtonOutline, Disabled: s.Disabled, OnClick: func() { c.AddToast(ui.Toast{Title: "Local changes reviewed"}) }}})
		case "progress":
			mygoui.Progress(c, mygoui.ProgressOptions{Value: s.Value, Min: 0, Max: 100, Indeterminate: s.Indeterminate, Label: "Sample completion"})
			ui.Row(c).Wrap().Gap(8).Children(func() {
				mygoui.Button(c, mygoui.ButtonOptions{Label: "Decrease", Variant: mygoui.ButtonOutline, Disabled: s.Disabled || s.Value <= 0, OnClick: func() { s.Value = max(0, s.Value-10) }})
				mygoui.Button(c, mygoui.ButtonOptions{Label: "Increase", Variant: mygoui.ButtonOutline, Disabled: s.Disabled || s.Value >= 100, OnClick: func() { s.Value = min(100, s.Value+10) }})
			})
		case "direction":
			mygoui.Toggle(c, &s.RTL, mygoui.ToggleOptions{Label: "Right to left"})
			direction := ui.LeftToRight
			if s.RTL {
				direction = ui.RightToLeft
			}
			layout.Direction(c, direction, func() {
				ui.Row(c).FillWidth().Gap(12).Children(func() { ui.Text(c, "Start"); ui.Text(c, "Content").Grow(1); ui.Text(c, "End") })
			}).FillWidth()
		case "resizable":
			layout.Resizable(c, &s.Split, layout.ResizableOptions{Label: "Sample panels", Vertical: s.Vertical, Disabled: s.Disabled}, func() { ui.Text(c, "First pane").Padding(12) }, func() { ui.Text(c, "Second pane").Padding(12) }).FillWidth().Height(200)
		case "scroll-area":
			layout.ScrollArea(c, &s.Scroll, layout.ScrollAreaOptions{Label: "Sample documents", Orientation: layout.ScrollVertical}, func() {
				for i := range 24 {
					ui.Text(c, "Local document "+strconv.Itoa(i+1)).Padding(8)
				}
			}).FillWidth().Height(180)
		case "carousel":
			layout.Carousel(c, &s.Carousel, layout.CarouselOptions{Label: "Sample pages", Count: 5, Width: 280, Height: 160, Vertical: s.Vertical, Disabled: s.Disabled}, func(index int) {
				ui.Column(c).Fill().Center().Children(func() { ui.Text(c, "Sample "+strconv.Itoa(index+1)).Bold() })
			})
		case "table":
			s.Table.List.CellSelect, s.Table.List.DragSelect = s.CellSelection, s.DragSelection
			data.Table(c, &s.Table, data.TableOptions{Label: "Sample project table", Caption: "Local project files", Columns: sampleColumns, Rows: s.Rows, Disabled: s.Disabled, OnActivate: func(row data.RowID) { c.AddToast(ui.Toast{Title: "Opened " + string(row)}) }}).FillWidth().Height(250)
		case "data-table":
			s.DataTable.Table.List.CellSelect, s.DataTable.Table.List.DragSelect = s.CellSelection, s.DragSelection
			data.DataTable(c, &s.DataTable, data.DataTableOptions{TableOptions: data.TableOptions{Label: "Sample project data", Columns: sampleColumns, Rows: s.Rows, Disabled: s.Disabled, OnActivate: func(row data.RowID) { c.AddToast(ui.Toast{Title: "Opened " + string(row)}) }}, Revision: s.Revision, FilterPlaceholder: "Search projects…", Height: 220})
		}
	})
}
