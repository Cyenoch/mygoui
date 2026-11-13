// Package locale owns library-level accessible labels. Applications set their
// language once; the library never registers a language or reads the host
// locale behind the caller's back.
package locale

import "sync/atomic"

// Messages holds user-visible library text that is not supplied by the caller.
type Messages struct {
	Loading       string
	Close         string
	Cancel        string
	Confirm       string
	Done          string
	Clear         string
	Copy          string
	Copied        string
	Search        string
	ShowPassword  string
	HidePassword  string
	Increase      string
	Decrease      string
	Previous      string
	Next          string
	MorePages     string
	More          string
	Remove        string
	Add           string
	Browse        string
	Key           string
	Value         string
	DuplicateKey  string
	Notifications string
	Required      string
	Invalid       string
	NoResults     string
	Today         string
	Selected      string
	Empty         string
	Resize        string
	Reorder       string
	Back          string
	Modified      string
	Page          string
	ZoomIn        string
	ZoomOut       string
	ResetZoom     string
	ShowMore      string
	ShowLess      string
	ExpandAll     string
	CollapseAll   string
	CopyValue     string
	CopyPath      string
	JumpToLatest  string
	CopyLine      string
	CopyAll       string
	WrapLines     string
	LevelDebug    string
	LevelInfo     string
	LevelWarning  string
	LevelError    string
	NoChanges     string
	SideBySide    string
	CloseOthers   string
	CloseToRight  string
	AllTabs       string
	// TitleBar names a window's title bar; ToggleSidebar and
	// ToggleInspector name the buttons showing its side panes.
	TitleBar        string
	ToggleSidebar   string
	ToggleInspector string
	// ChooseYear names a calendar's year grid, OpenCalendar the button
	// opening a date picker's calendar, and StartDate and EndDate the
	// editors of a date range.
	ChooseYear   string
	OpenCalendar string
	StartDate    string
	EndDate      string
	// NoData is what an empty chart shows, and Range names a chart's
	// range slider.
	NoData string
	Range  string
	// Log names a log viewer, and CopyNamed is a fmt format naming what a
	// copy button copies.
	Log       string
	CopyNamed string
	// StepOf is a fmt format of a step's position and the step count.
	StepOf string
	// CloseNamed and RemoveNamed are fmt formats naming what a close or
	// remove button acts on.
	CloseNamed  string
	RemoveNamed string
	// VerificationCode names a one-time code input and Submit a composer's
	// send button; Uploaded and UploadFailed name a file's upload status.
	VerificationCode string
	Submit           string
	Uploaded         string
	UploadFailed     string
	// Sidebar names a sidebar's navigation, Success a status that went
	// well, and MoreCount is a fmt format of how many items a group shows
	// no room for. MovedTo is a fmt format of an item's label, its new
	// place and the item count, read out when it is reordered.
	Sidebar   string
	Success   string
	MoreCount string
	MovedTo   string
	// Carousel and Slide are what assistive technology calls a carousel
	// and its slides; Breadcrumbs names a trail of links without Label.
	Carousel    string
	Slide       string
	Breadcrumbs string
}

var current atomic.Pointer[Messages]

func init() { Set(Messages{}) }

// Set publishes library messages. Empty fields keep the English default.
func Set(m Messages) {
	d := Messages{
		Loading: "Loading", Close: "Close", Clear: "Clear", Search: "Search",
		Copy: "Copy", Copied: "Copied",
		Cancel: "Cancel", Confirm: "Confirm", Done: "Done",
		ShowPassword: "Show characters", HidePassword: "Conceal characters",
		Increase: "Increase", Decrease: "Decrease", Previous: "Previous", Next: "Next",
		MorePages: "More pages", More: "More", Remove: "Remove", Notifications: "Notifications",
		Add: "Add", Browse: "Browse", Key: "Key", Value: "Value", DuplicateKey: "Duplicate key",
		Required: "Required", Invalid: "Invalid value", NoResults: "No results",
		Today: "Today", Selected: "Selected", Empty: "Nothing here",
		Resize: "Resize", Reorder: "Reorder",
		Back: "Back", Modified: "Modified", Page: "Page",
		ZoomIn: "Zoom in", ZoomOut: "Zoom out", ResetZoom: "Reset zoom",
		ShowMore: "Show more", ShowLess: "Show less",
		ExpandAll: "Expand all", CollapseAll: "Collapse all",
		CopyValue: "Copy value", CopyPath: "Copy path", JumpToLatest: "Jump to latest",
		CopyLine: "Copy line", CopyAll: "Copy all", WrapLines: "Wrap lines",
		LevelDebug: "Debug", LevelInfo: "Info", LevelWarning: "Warning", LevelError: "Error",
		CloseOthers: "Close others", CloseToRight: "Close to the right", AllTabs: "All tabs",
		NoChanges: "No changes", SideBySide: "Side by side", StepOf: "Step %d of %d",
		CloseNamed: "Close %s", RemoveNamed: "Remove %s",
		TitleBar: "Title bar", ToggleSidebar: "Toggle sidebar", ToggleInspector: "Toggle inspector",
		ChooseYear: "Choose year", OpenCalendar: "Open calendar", StartDate: "Start date", EndDate: "End date",
		NoData: "No data", Range: "Range", Log: "Log", CopyNamed: "Copy %s",
		VerificationCode: "Verification code", Submit: "Submit",
		Uploaded: "Uploaded", UploadFailed: "Upload failed",
		Sidebar: "Sidebar", Success: "Success", MoreCount: "%d more",
		MovedTo:  "%s, %d of %d",
		Carousel: "carousel", Slide: "slide", Breadcrumbs: "Breadcrumbs",
	}
	override := func(dst *string, src string) {
		if src != "" {
			*dst = src
		}
	}
	override(&d.Loading, m.Loading)
	override(&d.Close, m.Close)
	override(&d.Cancel, m.Cancel)
	override(&d.Confirm, m.Confirm)
	override(&d.Done, m.Done)
	override(&d.Clear, m.Clear)
	override(&d.Copy, m.Copy)
	override(&d.Copied, m.Copied)
	override(&d.Search, m.Search)
	override(&d.ShowPassword, m.ShowPassword)
	override(&d.HidePassword, m.HidePassword)
	override(&d.Increase, m.Increase)
	override(&d.Decrease, m.Decrease)
	override(&d.Previous, m.Previous)
	override(&d.Next, m.Next)
	override(&d.MorePages, m.MorePages)
	override(&d.More, m.More)
	override(&d.Remove, m.Remove)
	override(&d.Add, m.Add)
	override(&d.Browse, m.Browse)
	override(&d.Key, m.Key)
	override(&d.Value, m.Value)
	override(&d.DuplicateKey, m.DuplicateKey)
	override(&d.Notifications, m.Notifications)
	override(&d.Required, m.Required)
	override(&d.Invalid, m.Invalid)
	override(&d.NoResults, m.NoResults)
	override(&d.Today, m.Today)
	override(&d.Selected, m.Selected)
	override(&d.Empty, m.Empty)
	override(&d.Resize, m.Resize)
	override(&d.Reorder, m.Reorder)
	override(&d.Back, m.Back)
	override(&d.Modified, m.Modified)
	override(&d.Page, m.Page)
	override(&d.ZoomIn, m.ZoomIn)
	override(&d.ZoomOut, m.ZoomOut)
	override(&d.ResetZoom, m.ResetZoom)
	override(&d.ShowMore, m.ShowMore)
	override(&d.ShowLess, m.ShowLess)
	override(&d.ExpandAll, m.ExpandAll)
	override(&d.CollapseAll, m.CollapseAll)
	override(&d.CopyValue, m.CopyValue)
	override(&d.CopyPath, m.CopyPath)
	override(&d.JumpToLatest, m.JumpToLatest)
	override(&d.CopyLine, m.CopyLine)
	override(&d.CopyAll, m.CopyAll)
	override(&d.WrapLines, m.WrapLines)
	override(&d.LevelDebug, m.LevelDebug)
	override(&d.LevelInfo, m.LevelInfo)
	override(&d.LevelWarning, m.LevelWarning)
	override(&d.LevelError, m.LevelError)
	override(&d.CloseOthers, m.CloseOthers)
	override(&d.CloseToRight, m.CloseToRight)
	override(&d.AllTabs, m.AllTabs)
	override(&d.StepOf, m.StepOf)
	override(&d.CloseNamed, m.CloseNamed)
	override(&d.RemoveNamed, m.RemoveNamed)
	override(&d.NoChanges, m.NoChanges)
	override(&d.SideBySide, m.SideBySide)
	override(&d.TitleBar, m.TitleBar)
	override(&d.ToggleSidebar, m.ToggleSidebar)
	override(&d.ToggleInspector, m.ToggleInspector)
	override(&d.ChooseYear, m.ChooseYear)
	override(&d.OpenCalendar, m.OpenCalendar)
	override(&d.StartDate, m.StartDate)
	override(&d.EndDate, m.EndDate)
	override(&d.NoData, m.NoData)
	override(&d.Range, m.Range)
	override(&d.Log, m.Log)
	override(&d.CopyNamed, m.CopyNamed)
	override(&d.VerificationCode, m.VerificationCode)
	override(&d.Submit, m.Submit)
	override(&d.Uploaded, m.Uploaded)
	override(&d.UploadFailed, m.UploadFailed)
	override(&d.Sidebar, m.Sidebar)
	override(&d.Success, m.Success)
	override(&d.MoreCount, m.MoreCount)
	override(&d.MovedTo, m.MovedTo)
	override(&d.Carousel, m.Carousel)
	override(&d.Slide, m.Slide)
	override(&d.Breadcrumbs, m.Breadcrumbs)
	current.Store(&d)
}

// M returns the published messages.
func M() Messages {
	m := current.Load()
	if m == nil {
		Set(Messages{})
		m = current.Load()
	}
	return *m
}
