package gallery

import (
	"cmp"
	"fmt"
	"slices"
	"strings"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/data"
	"github.com/Cyenoch/mygoui/navigation"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// tableDeployment is one in-memory fixture row.
type tableDeployment struct {
	ID          string
	Service     string
	Environment string
	Status      string
	Requests    int
}

// tableState owns the fixtures, the table's caller-owned state and the demo
// switches shown above the table.
type tableState struct {
	ready           bool
	Rows            []tableDeployment
	Table           data.TableState
	SortTable       data.TableState
	EmptyTable      data.TableState
	LockedTable     data.TableState
	SecondaryTable  data.TableState
	VisibilityTable data.TableState
	PageTable       data.TableState
	Page            int
	ShowEnvironment bool
	ShowRequests    bool
	Note            string
}

const tableLockedRow = 6

func (st *tableState) ensure() {
	if st.ready {
		return
	}
	st.ready = true
	st.Rows = []tableDeployment{
		{ID: "gateway", Service: "Gateway", Environment: "prod", Status: "Healthy", Requests: 9210},
		{ID: "billing", Service: "Billing", Environment: "prod", Status: "Degraded", Requests: 4312},
		{ID: "ledger", Service: "Ledger", Environment: "prod", Status: "Healthy", Requests: 1180},
		{ID: "search", Service: "Search", Environment: "staging", Status: "Healthy", Requests: 7740},
		{ID: "identity", Service: "Identity", Environment: "prod", Status: "Paused", Requests: 640},
		{ID: "notifications", Service: "Notifications", Environment: "staging", Status: "Healthy", Requests: 2305},
		{ID: "media", Service: "Media", Environment: "dev", Status: "Failed", Requests: 312},
		{ID: "analytics", Service: "Analytics", Environment: "prod", Status: "Healthy", Requests: 5521},
		{ID: "webhook", Service: "Webhook", Environment: "staging", Status: "Degraded", Requests: 1907},
		{ID: "scheduler", Service: "Scheduler", Environment: "prod", Status: "Healthy", Requests: 860},
		{ID: "reports", Service: "Reports", Environment: "dev", Status: "Paused", Requests: 95},
		{ID: "session", Service: "Session", Environment: "prod", Status: "Healthy", Requests: 12100},
		{ID: "audit", Service: "Audit", Environment: "staging", Status: "Healthy", Requests: 440},
		{ID: "exports", Service: "Exports", Environment: "dev", Status: "Healthy", Requests: 158},
		{ID: "payments", Service: "Payments", Environment: "prod", Status: "Degraded", Requests: 3320},
		{ID: "ingest", Service: "Ingest", Environment: "prod", Status: "Healthy", Requests: 6680},
		{ID: "recommendations", Service: "Recommendations", Environment: "staging", Status: "Healthy", Requests: 2290},
		{ID: "search-index", Service: "Search index", Environment: "dev", Status: "Failed", Requests: 720},
	}
	st.Table.Selection = data.TableSelection{
		Mode: data.TableSelectionMultiple,
		Row:  -1,
	}
	st.Table.Selection.Rows.Add("gateway")
	st.Table.Selection.Rows.Add("search")
	st.Table.Revision = 1
	st.SortTable.Sort = data.TableSort{Column: "service"}
	st.SortTable.Selection = data.TableSelection{Mode: data.TableSelectionNone, Row: -1}
	st.EmptyTable.Selection = data.TableSelection{Mode: data.TableSelectionNone, Row: -1}
	st.LockedTable.Selection = data.TableSelection{Mode: data.TableSelectionNone, Row: -1}
	st.SecondaryTable.Selection = data.TableSelection{Mode: data.TableSelectionNone, Row: -1}
	st.VisibilityTable.Selection = data.TableSelection{Mode: data.TableSelectionNone, Row: -1}
	st.PageTable.Selection = data.TableSelection{Mode: data.TableSelectionNone, Row: -1}
	st.Page = 1
	st.ShowEnvironment = true
	st.ShowRequests = true
	st.Note = "No row action yet"
}

func (st *tableState) columns(c *ui.Context, narrow bool) []data.TableColumnOptions {
	return tableColumns(c, narrow, st.Rows)
}

// tableColumns declares the four demo columns plus the controls the examples
// add. Widths leave room for HeroUI's header type and a soft chip at every
// text scale; the sum of MinWidth is the table's minimum width, so a narrower
// panel scrolls the table sideways instead of squeezing the headers. Below the
// gallery's compact breakpoint the optional columns are dropped instead.
func tableColumns(c *ui.Context, narrow bool, rows []tableDeployment) []data.TableColumnOptions {
	columns := []data.TableColumnOptions{
		{
			ID: "service", Header: "Service", Sortable: true, Resizable: true,
			MinWidth: theme.Scale(c, 140),
			Cell: func(c *ui.Context, cell data.TableCell) ui.Element {
				return ui.Text(c, rows[cell.Row.Index].Service)
			},
		},
		{
			ID: "status", Header: "Status", Sortable: true, Resizable: true,
			Width: theme.Scale(c, 104), MinWidth: theme.Scale(c, 104),
			Cell: func(c *ui.Context, cell data.TableCell) ui.Element {
				status := rows[cell.Row.Index].Status
				return mygoui.Chip(c, mygoui.ChipOptions{
					Label: status, Color: tableStatusColor(status), Variant: mygoui.ChipSoft,
					Size: mygoui.ChipSmall,
				})
			},
		},
	}
	if !narrow {
		columns = append(columns,
			data.TableColumnOptions{
				ID: "environment", Header: "Environment", Sortable: true, Resizable: true, Hideable: true,
				MinWidth: theme.Scale(c, 110),
				Cell: func(c *ui.Context, cell data.TableCell) ui.Element {
					return ui.Text(c, rows[cell.Row.Index].Environment)
				},
			},
			data.TableColumnOptions{
				ID: "requests", Header: "Requests", Sortable: true, Resizable: true, Hideable: true, Align: ui.End,
				MinWidth: theme.Scale(c, 88),
				Cell: func(c *ui.Context, cell data.TableCell) ui.Element {
					return ui.Text(c, fmt.Sprint(rows[cell.Row.Index].Requests))
				},
			},
		)
	}
	return columns
}

// tableExampleOptions builds the shared fixture table for one example.
func (st *tableState) tableExampleOptions(c *ui.Context, state *data.TableState, count int) data.TableOptions {
	sort := state.Sort
	width, _ := c.Size()
	narrow := width < theme.Scale(c, 980)
	return data.TableOptions{
		Label:   "Deployments",
		Variant: data.TablePrimary,
		Columns: st.columns(c, narrow),
		Rows:    count,
		State:   state,
		RowKey:  func(row int) any { return st.Rows[row].ID },
		Compare: func(a, b data.TableRow) int {
			left, right := st.Rows[a.Index], st.Rows[b.Index]
			switch sort.Column {
			case "service":
				return strings.Compare(left.Service, right.Service)
			case "environment":
				return strings.Compare(left.Environment, right.Environment)
			case "status":
				return strings.Compare(left.Status, right.Status)
			case "requests":
				return cmp.Compare(left.Requests, right.Requests)
			}
			return 0
		},
	}
}

// tableExampleFrame sizes one table and its caption inside an example panel.
func tableExampleFrame(c *ui.Context, options data.TableOptions, caption func() string) {
	tableExampleFrameHeight(c, options, 320, caption)
}

// tableExampleFrameHeight sizes a table of a given height above its caption.
func tableExampleFrameHeight(c *ui.Context, options data.TableOptions, height float32, caption func() string) {
	ui.Column(c).FillWidth().MinWidth(0).Gap(theme.Space12).Children(func() {
		ui.Column(c).Height(theme.Scale(c, height)).MinHeight(0).Children(func() {
			data.Table(c, options).Grow(1)
		})
		mygoui.Typography(c, caption(), mygoui.TypographyOptions{
			Type: theme.TypeBodyXS, Muted: true,
		})
	})
}

// tablePageSize mirrors HeroUI's four-row demo pages.
const tablePageSize = 4

// tableDeploymentLess orders deployments for one sort column.
func tableDeploymentLess(sort data.TableSort) func(a, b tableDeployment) int {
	return func(a, b tableDeployment) int {
		var order int
		switch sort.Column {
		case "service":
			order = strings.Compare(a.Service, b.Service)
		case "environment":
			order = strings.Compare(a.Environment, b.Environment)
		case "status":
			order = strings.Compare(a.Status, b.Status)
		case "requests":
			order = cmp.Compare(a.Requests, b.Requests)
		}
		if sort.Descending {
			return -order
		}
		return order
	}
}

// tablePaginatedExample mirrors HeroUI's paginated table: four columns, four
// rows a page, and the footer summary with previous, page and next controls.
func tablePaginatedExample(c *ui.Context, st *tableState) {
	all := st.Rows[:8]
	sorted := append([]tableDeployment(nil), all...)
	if st.PageTable.Sort.Column != "" {
		slices.SortStableFunc(sorted, tableDeploymentLess(st.PageTable.Sort))
	}
	pages := (len(sorted) + tablePageSize - 1) / tablePageSize
	st.Page = min(max(st.Page, 1), pages)
	start := (st.Page - 1) * tablePageSize
	end := min(start+tablePageSize, len(sorted))
	pageRows := sorted[start:end]
	width, _ := c.Size()
	narrow := width < theme.Scale(c, 980)
	options := data.TableOptions{
		Label:   "Deployments",
		Variant: data.TablePrimary,
		Columns: tableColumns(c, narrow, pageRows),
		Rows:    len(pageRows),
		State:   &st.PageTable,
		RowKey:  func(row int) any { return pageRows[row].ID },
		OnSort:  func(data.TableSort) { st.Page = 1 },
		Footer: func(c *ui.Context) ui.Element {
			return navigation.Pagination(c, navigation.PaginationOptions{
				Page:          &st.Page,
				Pages:         pages,
				Size:          navigation.PaginationSmall,
				Label:         "Deployments pages",
				PreviousLabel: "Prev",
				Summary: func() {
					ui.Text(c, fmt.Sprintf("%d to %d of %d results", start+1, end, len(sorted)))
				},
			})
		},
	}
	tableExampleFrameHeight(c, options, 320, func() string {
		return fmt.Sprintf("Page %d of %d; %d deployments in four-row pages.", st.Page, pages, len(sorted))
	})
}

// tableSelectionExample chooses rows and runs a row action.
func tableSelectionExample(c *ui.Context, st *tableState) {
	options := st.tableExampleOptions(c, &st.Table, len(st.Rows))
	options.OnRowActivate = func(row data.TableRow) {
		st.Note = "Opened " + st.Rows[row.Index].Service
	}
	options.RowActions = func(c *ui.Context, row data.TableRow) ui.Element {
		name := st.Rows[row.Index].Service
		return mygoui.Button(c, mygoui.ButtonOptions{
			Label: "Restart " + name, Icon: lucide.Play(), IconOnly: true,
			Variant: mygoui.ButtonGhost, Size: theme.ButtonSmall,
			OnClick: func() { st.Note = "Restarted " + name },
		})
	}
	options.ActionsHeader = "Run"
	options.ActionsWidth = theme.Scale(c, 64)
	tableExampleFrame(c, options, func() string {
		return "Selection: " + st.selectionSummary()
	})
}

// tableSortingExample starts sorted by service.
func tableSortingExample(c *ui.Context, st *tableState) {
	options := st.tableExampleOptions(c, &st.SortTable, len(st.Rows))
	tableExampleFrame(c, options, func() string {
		if st.SortTable.Sort.Column == "" {
			return "Unsorted; click a header to sort"
		}
		return "Sorted by " + st.SortTable.Sort.Column
	})
}

// tableEmptyExample shows the empty state.
func tableEmptyExample(c *ui.Context, st *tableState) {
	options := st.tableExampleOptions(c, &st.EmptyTable, 0)
	options.Empty = func(c *ui.Context) ui.Element {
		return mygoui.EmptyState(c, mygoui.EmptyStateOptions{
			Title:       "No deployments",
			Description: "Deployments you create appear here.",
			Icon:        lucide.Box(),
		})
	}
	tableExampleFrame(c, options, func() string {
		return "The empty builder replaces the row area."
	})
}

// tableDisabledExample disables one row.
func tableDisabledExample(c *ui.Context, st *tableState) {
	options := st.tableExampleOptions(c, &st.LockedTable, len(st.Rows))
	options.RowDisabled = func(row int) bool { return row == tableLockedRow }
	tableExampleFrame(c, options, func() string {
		return fmt.Sprintf("Row %d is disabled and cannot be selected.", tableLockedRow+1)
	})
}

// tableSecondaryExample shows the flush variant, whose header is a standalone
// card over a table with no frame.
func tableSecondaryExample(c *ui.Context, st *tableState) {
	options := st.tableExampleOptions(c, &st.SecondaryTable, len(st.Rows))
	options.Variant = data.TableSecondary
	tableExampleFrame(c, options, func() string {
		return "Secondary: no card, a standalone header, row separators."
	})
}

// tableVisibilityExample hides columns through the caller-owned state; the
// hidden columns keep their widths and sort state while they are out of view.
func tableVisibilityExample(c *ui.Context, st *tableState) {
	if st.ShowEnvironment {
		st.VisibilityTable.ShowColumn("environment")
	} else {
		st.VisibilityTable.HideColumn("environment")
	}
	if st.ShowRequests {
		st.VisibilityTable.ShowColumn("requests")
	} else {
		st.VisibilityTable.HideColumn("requests")
	}
	options := st.tableExampleOptions(c, &st.VisibilityTable, len(st.Rows))
	// This example keeps every column so the toggles always have something
	// to hide; a narrow window scrolls the table sideways instead.
	options.Columns = st.columns(c, false)
	ui.Column(c).FillWidth().MinWidth(0).Gap(theme.Space12).Children(func() {
		ui.Row(c).Gap(theme.Space8).Children(func() {
			mygoui.ToggleButton(c, mygoui.ToggleButtonOptions{
				Label: "Environment", Size: theme.ButtonSmall, Selected: &st.ShowEnvironment,
			})
			mygoui.ToggleButton(c, mygoui.ToggleButtonOptions{
				Label: "Requests", Size: theme.ButtonSmall, Selected: &st.ShowRequests,
			})
		})
		tableExampleFrame(c, options, func() string {
			var hidden []string
			if !st.ShowEnvironment {
				hidden = append(hidden, "Environment")
			}
			if !st.ShowRequests {
				hidden = append(hidden, "Requests")
			}
			if len(hidden) == 0 {
				return "All columns shown; hide one with the toggles above."
			}
			return "Hidden: " + strings.Join(hidden, ", ")
		})
	})
}

func (st *tableState) selectionSummary() string {
	selection := st.Table.Selection
	switch selection.Mode {
	case data.TableSelectionNone:
		return "disabled"
	case data.TableSelectionSingle:
		if selection.Row < 0 {
			return "no row"
		}
		return fmt.Sprintf("row %d", selection.Row+1)
	}
	var names []string
	for _, row := range st.Rows {
		if selection.HasKey(row.ID) {
			names = append(names, row.Service)
		}
	}
	if len(names) == 0 {
		return "no rows"
	}
	if len(names) > 3 {
		return fmt.Sprintf("%d rows", len(names))
	}
	return strings.Join(names, ", ")
}

func tableStatusColor(status string) mygoui.StatusColor {
	switch status {
	case "Healthy":
		return mygoui.ColorSuccess
	case "Degraded":
		return mygoui.ColorWarning
	case "Failed":
		return mygoui.ColorDanger
	}
	return mygoui.ColorDefault
}
