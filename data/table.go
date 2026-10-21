// Package data composes maintained native virtual tables with typed row and
// column identities. It does not implement a second layout or selection engine.
package data

import (
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

type RowID string
type ColumnID string

type Row struct {
	ID    RowID
	Cells []string
}

type Column struct {
	ID                        ColumnID
	Title                     string
	Width, MinWidth, MaxWidth float32
	Align                     ui.Align
	Fixed                     bool
	// RowHeader selects whole rows in native cell-selection mode.
	RowHeader bool
	Compare   func(a, b Row) int
}

type TableState struct {
	List        ui.ListState
	Selection   ui.Selection[RowID]
	cursor      int
	initialized bool
	columns     []ui.TableColumn
}

type TableOptions struct {
	Label, Caption, Footer string
	Rows                   []Row
	Columns                []Column
	Disabled               bool
	Cell                   func(*ui.Context, Row, Column)
	// RowLabel supplies the visible row name for type-ahead and accessibility.
	// Custom cells may use it when their text differs from Row.Cells.
	RowLabel          func(Row) string
	OnActivate        func(RowID)
	OnSelectionChange func()
	// OnReorder receives stable IDs and the ID to insert before; empty is end.
	// The caller updates Rows and, for DataTable, its Revision.
	OnReorder func(rows []RowID, before RowID)
}

// Table is a bounded native virtual table with header, caption and footer.
// Set Height or Grow on the returned container. Every Row.ID and Column.ID must
// be nonempty and unique; selection survives row sorting, filtering and paging.
// Native ListState.Columns retains user column order and widths. Enable cell
// selection or drag selection with state.List.CellSelect or DragSelect. These
// native modes cannot be combined with OnReorder. Custom cells can inspect
// state.List.CellSelected to preserve selection contrast.
func Table(c *ui.Context, state *TableState, o TableOptions) ui.Element {
	if !state.initialized {
		state.cursor, state.initialized = -1, true
	}
	if cap(state.columns) < len(o.Columns) {
		state.columns = make([]ui.TableColumn, len(o.Columns))
	} else {
		state.columns = state.columns[:len(o.Columns)]
	}
	for i, col := range o.Columns {
		state.columns[i] = ui.TableColumn{ID: string(col.ID), Title: col.Title, Width: col.Width,
			MinWidth: col.MinWidth, MaxWidth: col.MaxWidth, Align: col.Align,
			Fixed: col.Fixed, RowHeader: col.RowHeader,
			Sortable: !col.RowHeader && (state.List.Sort != nil || state.List.SortBy != nil)}
	}
	state.List.Selected, state.List.Selection = &state.cursor, &state.Selection
	state.List.Key = func(i int) any { return o.Rows[i].ID }
	state.List.Label = func(i int) string {
		if o.RowLabel != nil {
			return o.RowLabel(o.Rows[i])
		}
		if len(o.Rows[i].Cells) != 0 {
			return o.Rows[i].Cells[0]
		}
		return string(o.Rows[i].ID)
	}
	state.List.Reorder = nil
	if o.OnReorder != nil && !o.Disabled {
		state.List.Reorder = func(indices []int, before int) {
			ids := make([]RowID, len(indices))
			for i, index := range indices {
				ids[i] = o.Rows[index].ID
			}
			var target RowID
			if before < len(o.Rows) {
				target = o.Rows[before].ID
			}
			o.OnReorder(ids, target)
		}
	}
	t := theme.Colors(c)
	return ui.Column(c).MinWidth(0).MinHeight(0).Gap(theme.Space8).Disabled(o.Disabled).
		TextColor(t.Foreground).Children(func() {
		if o.Caption != "" {
			ui.Text(c, o.Caption).FontSize(theme.Scale(c, theme.Small().Size)).TextColor(t.MutedForeground)
		}
		ui.Table(c.Key("table"), &state.List, state.columns, len(o.Rows), func(row, col int) {
			if o.Cell != nil {
				o.Cell(c, o.Rows[row], o.Columns[col])
			} else if col < len(o.Rows[row].Cells) {
				ui.Text(c, o.Rows[row].Cells[col]).FontSize(theme.Scale(c, theme.Body().Size))
			}
		}).Grow(1).Basis(0).MinHeight(0).MinWidth(0).Label(o.Label).
			Border(theme.HairlineWidth, t.Border).Radius(theme.Radii().XXL).
			OnChange(func() {
				if o.OnSelectionChange != nil {
					o.OnSelectionChange()
				}
			}).OnSubmit(func() {
			if o.OnActivate != nil && state.cursor >= 0 && state.cursor < len(o.Rows) {
				o.OnActivate(o.Rows[state.cursor].ID)
			}
		})
		if o.Footer != "" {
			ui.Text(c, o.Footer).FontSize(theme.Scale(c, theme.Small().Size)).TextColor(t.MutedForeground)
		}
	})
}
