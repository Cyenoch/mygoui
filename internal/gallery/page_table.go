package gallery

import (
	"cmp"
	"fmt"
	"strings"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/data"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// tableDeployment is one in-memory fixture row.
type tableDeployment struct {
	ID          string
	Service     string
	Environment string
	Status      string
	Owner       string
	Requests    int
}

// tableState owns the fixtures, the table's caller-owned state and the demo
// switches shown above the table.
type tableState struct {
	ready       bool
	Rows        []tableDeployment
	Table       data.TableState
	SortTable   data.TableState
	EmptyTable  data.TableState
	LockedTable data.TableState
	Note        string
}

const tableLockedRow = 6

func (st *tableState) ensure() {
	if st.ready {
		return
	}
	st.ready = true
	st.Rows = []tableDeployment{
		{ID: "gateway", Service: "Gateway", Environment: "prod", Status: "Healthy", Owner: "Ana", Requests: 9210},
		{ID: "billing", Service: "Billing", Environment: "prod", Status: "Degraded", Owner: "Rui", Requests: 4312},
		{ID: "ledger", Service: "Ledger", Environment: "prod", Status: "Healthy", Owner: "Mei", Requests: 1180},
		{ID: "search", Service: "Search", Environment: "staging", Status: "Healthy", Owner: "Ilya", Requests: 7740},
		{ID: "identity", Service: "Identity", Environment: "prod", Status: "Paused", Owner: "Nina", Requests: 640},
		{ID: "notifications", Service: "Notifications", Environment: "staging", Status: "Healthy", Owner: "Owen", Requests: 2305},
		{ID: "media", Service: "Media", Environment: "dev", Status: "Failed", Owner: "Trina", Requests: 312},
		{ID: "analytics", Service: "Analytics", Environment: "prod", Status: "Healthy", Owner: "Sam", Requests: 5521},
		{ID: "webhook", Service: "Webhook", Environment: "staging", Status: "Degraded", Owner: "Ada", Requests: 1907},
		{ID: "scheduler", Service: "Scheduler", Environment: "prod", Status: "Healthy", Owner: "Jun", Requests: 860},
		{ID: "reports", Service: "Reports", Environment: "dev", Status: "Paused", Owner: "Lea", Requests: 95},
		{ID: "session", Service: "Session", Environment: "prod", Status: "Healthy", Owner: "Kai", Requests: 12100},
		{ID: "audit", Service: "Audit", Environment: "staging", Status: "Healthy", Owner: "Ravi", Requests: 440},
		{ID: "exports", Service: "Exports", Environment: "dev", Status: "Healthy", Owner: "Zoe", Requests: 158},
		{ID: "payments", Service: "Payments", Environment: "prod", Status: "Degraded", Owner: "Mina", Requests: 3320},
		{ID: "ingest", Service: "Ingest", Environment: "prod", Status: "Healthy", Owner: "Theo", Requests: 6680},
		{ID: "recommendations", Service: "Recommendations", Environment: "staging", Status: "Healthy", Owner: "Iris", Requests: 2290},
		{ID: "search-index", Service: "Search index", Environment: "dev", Status: "Failed", Owner: "Luis", Requests: 720},
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
	st.Note = "No row action yet"
}

func (st *tableState) columns(c *ui.Context, narrow bool) []data.TableColumnOptions {
	rows := st.Rows
	columns := []data.TableColumnOptions{
		{
			ID: "service", Header: "Service", Sortable: true, Resizable: true, MinWidth: theme.Scale(c, 120),
			Cell: func(c *ui.Context, cell data.TableCell) ui.Element {
				return ui.Text(c, rows[cell.Row.Index].Service)
			},
		},
		{
			ID: "status", Header: "Status", Sortable: true, Width: theme.Scale(c, 88),
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
				ID: "environment", Header: "Environment", Sortable: true, Width: theme.Scale(c, 100),
				Cell: func(c *ui.Context, cell data.TableCell) ui.Element {
					return ui.Text(c, rows[cell.Row.Index].Environment)
				},
			},
			data.TableColumnOptions{
				ID: "requests", Header: "Requests", Sortable: true, Align: ui.End, Width: theme.Scale(c, 84),
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
	ui.Column(c).FillWidth().MinWidth(0).Gap(theme.Space12).Children(func() {
		ui.Column(c).Height(theme.Scale(c, 320)).MinHeight(0).Children(func() {
			data.Table(c, options).Grow(1)
		})
		mygoui.Typography(c, caption(), mygoui.TypographyOptions{
			Type: theme.TypeBodyXS, Muted: true,
		})
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
			Description: "This table is configured with zero rows.",
			Icon:        lucide.Play(),
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
