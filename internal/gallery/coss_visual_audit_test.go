package gallery

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"slices"
	"testing"
	"time"

	"github.com/Cyenoch/mygoui/fonts"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

const auditWidth, auditHeight = 720, 720

// auditSpec names actual gallery controls; no component is reconstructed for the audit.
type auditSpec struct {
	example, target        string
	disabled, disabledText string
	open, openText         string
	motionExample          string
	openByHover            bool
	requiresOpen           bool
	selected               string
	motion                 bool
}

type auditCapture struct {
	State        string `json:"state"`
	Example      string `json:"example"`
	ExamplePage  string `json:"example_page_id,omitempty"`
	Target       string `json:"target,omitempty"`
	Path         string `json:"path"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	CanvasX      int    `json:"canvas_x"`
	CanvasY      int    `json:"canvas_y"`
	ReduceMotion bool   `json:"reduce_motion"`
}

type auditCoverage struct {
	Status string `json:"status"`
	Reason string `json:"reason,omitempty"`
	Path   string `json:"path,omitempty"`
}

type auditComponent struct {
	Slug       string                   `json:"coss_slug"`
	Page       string                   `json:"native_page_id"`
	Example    string                   `json:"default_example"`
	Mode       string                   `json:"mode"`
	Source     string                   `json:"source_url"`
	SourceDemo string                   `json:"source_demo_url,omitempty"`
	Captures   []auditCapture           `json:"captures"`
	States     map[string]auditCoverage `json:"states"`
	NativeMenu []string                 `json:"native_menu_items,omitempty"`
}

type auditManifest struct {
	CanvasWidth   int                 `json:"canvas_width"`
	CanvasHeight  int                 `json:"canvas_height"`
	Note          string              `json:"note"`
	Components    []auditComponent    `json:"components"`
	ContactSheets []auditContactSheet `json:"contact_sheets"`
}

type auditContactBounds struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

type auditContactCell struct {
	Path        string             `json:"capture_path"`
	Example     string             `json:"example"`
	ExamplePage string             `json:"example_page_id,omitempty"`
	Target      string             `json:"target,omitempty"`
	Bounds      auditContactBounds `json:"bounds"`
}

type auditContactRow struct {
	Page  string           `json:"native_page_id"`
	Slug  string           `json:"coss_slug"`
	State string           `json:"state"`
	Light auditContactCell `json:"light"`
	Dark  auditContactCell `json:"dark"`
}

type auditContactSheet struct {
	Kind   string            `json:"kind"`
	Page   int               `json:"page"`
	Path   string            `json:"path"`
	Width  int               `json:"width"`
	Height int               `json:"height"`
	Rows   []auditContactRow `json:"rows"`
}

type auditContactPair struct {
	page, slug  string
	light, dark auditCapture
}

type auditRunner struct {
	tester         *ui.Tester
	state          *State
	example        ComponentExample
	initial        *image.RGBA
	selectedTarget string
}

// Every entry selects a named control within the actual first example unless
// a more representative example is specified. Empty targets are static demos.
var auditSpecs = map[string]auditSpec{
	"accordion":         {target: "General", open: "General", openText: "The project name", disabled: "Disabled section", disabledText: "Audit log", motion: true},
	"alert":             {},
	"alertdialog":       {target: "Delete project", open: "Delete project", openText: "Delete permanently", requiresOpen: true, motion: true},
	"autocomplete":      {target: "Destination", open: "Destination", openText: "Tokyo", selected: "Tokyo"},
	"avatar":            {},
	"badge":             {},
	"breadcrumbs":       {target: "Workspace"},
	"button":            {target: "Default", disabled: "Disabled", disabledText: "Unavailable"},
	"calendar":          {target: "Next", selected: "calendar day"},
	"card":              {},
	"checkbox":          {target: "Release notes", selected: "Release notes", disabled: "Disabled", disabledText: "Locked by policy"},
	"checkboxgroup":     {target: "SMS", selected: "SMS"},
	"combobox":          {target: "Language", open: "Language", openText: "Standard library first.", selected: "Rust"},
	"command":           {target: "Search actions", disabledText: "Archive project"},
	"composer":          {target: "Message draft"},
	"contextmenu":       {target: "Document actions"},
	"datepicker":        {target: "Open deploy calendar", open: "Open deploy calendar", openText: "Previous", selected: "calendar day", disabled: "Disabled", disabledText: "Archived deploy", motion: true},
	"dialog":            {target: "Open dialog", open: "Open dialog", openText: "Rename workspace", requiresOpen: true, motion: true},
	"disclosure":        {target: "Release notes", open: "Release notes", openText: "Release notes are open.", disabled: "Disabled trigger", disabledText: "Archived projects", motion: true},
	"drawer":            {target: "Open Right drawer", open: "Open Right drawer", openText: "Right drawer", requiresOpen: true, motion: true},
	"emptystate":        {example: "Default"},
	"field":             {target: "Repository"},
	"fieldset":          {target: "Security alerts", disabled: "Disabled", disabledText: "Weekly digest"},
	"form":              {target: "Email"},
	"frame":             {target: "Save project"},
	"group":             {target: "Files"},
	"input":             {target: "Project name", disabled: "Disabled", disabledText: "Account ID"},
	"inputgroup":        {target: "Price"},
	"inputotp":          {target: "Verification code", disabled: "Disabled", disabledText: "Verification code"},
	"kbd":               {},
	"label":             {disabled: "Disabled", disabledText: "Managed by your organization"},
	"menu":              {target: "Reports", selected: "Reports", disabledText: "Archived items", motion: true},
	"meter":             {},
	"numberfield":       {target: "Increase", disabled: "Disabled", disabledText: "Seats"},
	"pagination":        {target: "Next", selected: "Next"},
	"popover":           {target: "Show keyboard shortcuts", open: "Show keyboard shortcuts", openText: "Got it", requiresOpen: true, motion: true},
	"previewcard":       {target: "Ada Lovelace", open: "Ada Lovelace", openText: "Follow updates", openByHover: true, requiresOpen: true, motion: true},
	"progressbar":       {target: "Advance", motionExample: "Indeterminate", motion: true},
	"radiogroup":        {target: "Free", selected: "Free", disabled: "Disabled", disabledText: "Accessible"},
	"scrollarea":        {},
	"segmentedcontrol":  {target: "Yearly", selected: "Yearly", disabled: "Content view", disabledText: "Map", motion: true},
	"select":            {target: "Fruit", open: "Apple", openText: "Banana", selected: "Banana", motion: true},
	"separator":         {},
	"sheet":             {target: "Open right sheet", open: "Open right sheet", openText: "Workspace settings", requiresOpen: true, motion: true},
	"sidebar":           {target: "Projects", selected: "Projects", disabledText: "Archived"},
	"skeleton":          {motionExample: "Animations", motion: true},
	"slider":            {target: "Volume", selected: "Volume", disabled: "Disabled", disabledText: "Locked level"},
	"spinner":           {motion: true},
	"switch":            {target: "Wi-Fi", selected: "Wi-Fi", disabled: "Disabled", disabledText: "Locked by policy", motion: true},
	"table":             {target: "Next", disabled: "Disabled row"},
	"tabs":              {target: "Password", selected: "Password", disabledText: "Danger", motion: true},
	"textarea":          {target: "Bio"},
	"toast":             {target: "Default toast", open: "Default toast", openText: "Draft saved", requiresOpen: true},
	"togglebutton":      {target: "Default off", selected: "Default off", disabled: "Disabled", disabledText: "Unavailable"},
	"togglebuttongroup": {target: "Left", selected: "Left"},
	"toolbar":           {target: "Bold"},
	"tooltip":           {example: "Arrow and rich content", target: "Copy share link", open: "Copy share link", openText: "Expires in 7 days", openByHover: true, requiresOpen: true, motion: true},
}

func TestCossVisualAudit(t *testing.T) {
	dir := os.Getenv("MYGOUI_COSS_VISUAL_AUDIT_DIR")
	if dir == "" {
		t.Skip("set MYGOUI_COSS_VISUAL_AUDIT_DIR to collect the native visual audit")
	}
	if err := fonts.Register(); err != nil {
		t.Fatal(err)
	}
	var err error
	dir, err = filepath.Abs(dir)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	before := theme.PreferencesSnapshot()
	t.Cleanup(func() { theme.Configure(before) })

	ids := make([]string, 0, len(cossSources))
	for id := range cossSources {
		ids = append(ids, id)
	}
	slices.Sort(ids)
	missing, extra := auditCatalogMismatch(cossSources, auditSpecs)
	if len(missing) != 0 || len(extra) != 0 {
		t.Fatalf("audit specs do not match catalog: missing IDs %v; extra IDs %v", missing, extra)
	}
	manifest := auditManifest{CanvasWidth: auditWidth, CanvasHeight: auditHeight,
		Note:          "Headless native gallery examples; screenshots are evidence for review, not upstream comparisons or GPU/accessibility verification.",
		ContactSheets: []auditContactSheet{}}
	var failures []string
	for _, id := range ids {
		spec, ok := auditSpecs[id]
		page, found := componentPage(id)
		if !ok || !found || len(page.Examples) == 0 {
			failures = append(failures, fmt.Sprintf("%s: missing audit spec or gallery page", id))
			continue
		}
		for _, mode := range []string{"light", "dark"} {
			entry := collectAudit(dir, page, spec, mode, &failures)
			manifest.Components = append(manifest.Components, entry)
		}
	}
	defaults, states, err := auditContactPairs(ids, manifest.Components)
	if err != nil {
		failures = append(failures, err.Error())
	} else if len(failures) == 0 {
		manifest.ContactSheets, err = writeAuditContactSheets(dir, defaults, states)
		if err != nil {
			failures = append(failures, fmt.Sprintf("contact sheets: %v", err))
		}
	}
	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	manifestPath := filepath.Join(dir, "manifest.json")
	if err := os.WriteFile(manifestPath, append(data, '\n'), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Logf("manifest: %s; components: %d; mode entries: %d; contact sheets: %d",
		manifestPath, len(ids), len(manifest.Components), len(manifest.ContactSheets))
	for _, failure := range failures {
		t.Error(failure)
	}
}

func auditCatalogMismatch(catalog map[string]string, specs map[string]auditSpec) (missing, extra []string) {
	for id := range catalog {
		if _, ok := specs[id]; !ok {
			missing = append(missing, id)
		}
	}
	for id := range specs {
		if _, ok := catalog[id]; !ok {
			extra = append(extra, id)
		}
	}
	slices.Sort(missing)
	slices.Sort(extra)
	return missing, extra
}

func TestAuditCatalogMismatch(t *testing.T) {
	catalog := map[string]string{"existing": "existing", "missing": "missing"}
	specs := map[string]auditSpec{"existing": {}, "extra": {}}
	missing, extra := auditCatalogMismatch(catalog, specs)
	if !slices.Equal(missing, []string{"missing"}) || !slices.Equal(extra, []string{"extra"}) {
		t.Fatalf("missing %v, extra %v; want [missing] and [extra]", missing, extra)
	}
}

func collectAudit(dir string, page ComponentPage, spec auditSpec, mode string, failures *[]string) auditComponent {
	example := page.Examples[0].Title
	if spec.example != "" {
		example = spec.example
	}
	e := auditComponent{Slug: cossSources[page.ID], Page: page.ID, Example: example, Mode: mode,
		Source: sourceURL(page.ID),
		States: map[string]auditCoverage{}}
	if page.ID == "sidebar" {
		e.SourceDemo = cossSourceRoot + "components/docs-sidebar.tsx"
	}
	for _, state := range []string{"default", "hover", "keyboard_focus", "selected", "disabled", "open", "motion"} {
		e.States[state] = auditCoverage{Status: "not_exercised", Reason: "Not exercised by this collector for this example."}
	}
	if page.ID == "autocomplete" || page.ID == "combobox" {
		e.States["motion"] = auditCoverage{Status: "not_exercised",
			Reason: "The first-open headless frames at 50 and 120 ms were identical; popup motion is not established."}
		e.States["disabled"] = auditCoverage{Status: "not_exercised",
			Reason: "The gallery has no disabled example for this control."}
	}
	if page.ID == "textarea" {
		e.States["disabled"] = auditCoverage{Status: "not_exercised",
			Reason: "The gallery has a read-only example, but no disabled TextArea example."}
	}
	fail := func(state string, err error) {
		message := fmt.Sprintf("%s/%s/%s/%s: %v", page.ID, mode, example, state, err)
		*failures = append(*failures, message)
		e.States[state] = auditCoverage{Status: "missing", Reason: err.Error()}
	}
	newRunner := func(title string, reduced bool) (*auditRunner, error) {
		return newAuditRunner(page, title, mode == "dark", reduced)
	}
	capture := func(r *auditRunner, state, target string, full, reduced bool) (string, error) {
		path := filepath.Join(dir, mode, e.Slug, state+".png")
		bounds, err := writeAuditPNG(path, r, full)
		if err != nil {
			return "", err
		}
		item := auditCapture{
			State: state, Example: r.example.Title, Target: target, Path: path,
			Width: bounds.Dx(), Height: bounds.Dy(), CanvasX: bounds.Min.X, CanvasY: bounds.Min.Y,
			ReduceMotion: reduced,
		}
		if r.state.Selected != page.ID {
			item.ExamplePage = r.state.Selected
		}
		e.Captures = append(e.Captures, item)
		return path, nil
	}

	r, err := newRunner(example, true)
	if err != nil {
		fail("default", err)
		return e
	}
	if page.ID == "contextmenu" {
		path, err := capture(r, "default", spec.target, false, true)
		if err != nil {
			fail("default", err)
		} else {
			e.States["default"] = auditCoverage{Status: "captured", Path: path,
				Reason: "Closed trigger only; the OS-owned menu itself is not in this screenshot."}
		}
		menuErr := r.tester.RightClick(spec.target)
		if menuErr == nil {
			e.NativeMenu = r.tester.Menu()
			if !slices.Contains(e.NativeMenu, "Rename") || !slices.Contains(e.NativeMenu, "Duplicate") {
				menuErr = fmt.Errorf("native menu items missing: %v", e.NativeMenu)
			}
		}
		if menuErr != nil {
			fail("open", menuErr)
		} else {
			reason := "The OS draws this native context menu; ui.NewTester exposes items but does not rasterize the menu."
			e.States["open"] = auditCoverage{Status: "unsupported", Reason: reason}
		}
	} else {
		if spec.requiresOpen {
			if _, err := capture(r, "trigger", spec.open, false, true); err != nil {
				fail("trigger", err)
			}
			if err := auditOpen(r, spec); err != nil {
				fail("default", err)
			} else if path, err := capture(r, "default", spec.open, true, true); err != nil {
				fail("default", err)
			} else {
				e.States["default"] = auditCoverage{Status: "captured", Path: path}
				e.States["open"] = auditCoverage{Status: "captured", Path: path, Reason: "The default component is visible when opened."}
			}
		} else if path, err := capture(r, "default", "", false, true); err != nil {
			fail("default", err)
		} else {
			e.States["default"] = auditCoverage{Status: "captured", Path: path}
		}
	}

	if spec.target != "" {
		for _, state := range []string{"hover", "keyboard_focus"} {
			target := spec.target
			if state == "keyboard_focus" {
				switch page.ID {
				case "menu":
					target = "Overview"
				case "radiogroup":
					target = "Pro"
				case "segmentedcontrol":
					target = "Monthly"
				case "tabs":
					target = "Overview"
				}
			}
			r, err := newRunner(example, true)
			if err == nil {
				if state == "hover" {
					err = auditHover(r, target)
				} else {
					err = auditFocus(r, target)
				}
			}
			if err != nil {
				fail(state, err)
				continue
			}
			path, err := capture(r, state, target, false, true)
			if err != nil {
				fail(state, err)
			} else {
				entry := auditCoverage{Status: "captured", Path: path}
				if page.ID == "contextmenu" {
					entry.Reason = "This is the trigger; the OS-owned menu is not rasterized."
				} else if state == "hover" && bytes.Equal(r.initial.Pix, r.tester.Image().Pix) {
					entry.Reason = "The pointer is over the target; this example has no distinct hover pixels."
				}
				e.States[state] = entry
			}
		}
	}
	if spec.open != "" && !spec.requiresOpen {
		r, err := newRunner(example, true)
		if err == nil {
			err = auditOpen(r, spec)
		}
		if err != nil {
			fail("open", err)
		} else if path, err := capture(r, "open", spec.open, true, true); err != nil {
			fail("open", err)
		} else {
			e.States["open"] = auditCoverage{Status: "captured", Path: path}
		}
	}
	if page.ID == "menu" {
		dropdown, found := componentPage("dropdown")
		var r *auditRunner
		var err error
		if !found {
			err = fmt.Errorf("gallery dropdown page is missing")
		} else {
			r, err = newAuditRunner(dropdown, "Radio menu", mode == "dark", true)
		}
		if err == nil {
			err = auditOpen(r, auditSpec{open: "Sort by", openText: "Last modified"})
		}
		if err == nil && !navigationStateOf(r.state).sortOpen {
			err = fmt.Errorf("dropdown's caller-owned menu state did not open")
		}
		if err != nil {
			fail("open", fmt.Errorf("dropdown/Radio menu: %w", err))
		} else if path, err := capture(r, "open", "Sort by", true, true); err != nil {
			fail("open", fmt.Errorf("dropdown/Radio menu: %w", err))
		} else {
			e.States["open"] = auditCoverage{Status: "captured", Path: path,
				Reason: "The dropdown page's Radio menu example opens a real Menu from its Sort by trigger."}
		}
	}
	if spec.disabled != "" || spec.disabledText != "" {
		title := example
		if spec.disabled != "" {
			title = spec.disabled
		}
		r, err := newRunner(title, true)
		if err == nil && spec.disabledText != "" && !r.tester.HasText(spec.disabledText) {
			err = fmt.Errorf("disabled example %q lacks %q", title, spec.disabledText)
		}
		if err == nil && page.ID == "table" {
			box, ok := r.tester.Find("Deployments")
			if !ok {
				err = fmt.Errorf("disabled table viewport missing")
			} else {
				r.tester.Scroll(box.X+box.W/2, box.Y+box.H/2, 0, 160)
				r.tester.Frame()
				row, found := r.tester.Find("Media")
				if !found || row.Y < box.Y || row.Y+row.H > box.Y+box.H {
					err = fmt.Errorf("disabled Media row is not visible after scrolling: %+v", row)
				}
			}
		}
		if err != nil {
			fail("disabled", err)
		} else if path, err := capture(r, "disabled", spec.disabledText, false, true); err != nil {
			fail("disabled", err)
		} else {
			e.States["disabled"] = auditCoverage{Status: "captured", Path: path, Reason: "Disabled control appears in this named example."}
		}
	}
	if page.ID == "calendar" || page.ID == "select" {
		title := example
		if page.ID == "calendar" {
			title = "Bounded dates"
		}
		r, err := newRunner(title, true)
		target := ""
		if err == nil {
			if page.ID == "calendar" {
				target, err = auditDisabledCalendar(r)
			} else {
				target, err = auditDisabledSelect(r, spec)
			}
		}
		if err != nil {
			fail("disabled", fmt.Errorf("example %q: %w", title, err))
		} else if path, err := capture(r, "disabled", target, page.ID == "select", true); err != nil {
			fail("disabled", fmt.Errorf("example %q: %w", title, err))
		} else {
			e.States["disabled"] = auditCoverage{Status: "captured", Path: path}
		}
	}
	if spec.selected != "" {
		state := "selected"
		if page.ID == "slider" {
			state = "value_changed"
			e.States["selected"] = auditCoverage{Status: "not_applicable", Reason: "The slider carries a continuous value."}
		}
		r, err := newRunner(example, true)
		if err == nil {
			err = auditSelect(r, page.ID, spec)
		}
		if err != nil {
			fail(state, err)
		} else if path, err := capture(r, state, r.selectedTarget, page.ID == "datepicker", true); err != nil {
			fail(state, err)
		} else {
			e.States[state] = auditCoverage{Status: "captured", Path: path}
		}
	}
	if page.ID == "pagination" {
		r, err := newRunner(example, true)
		if err == nil {
			err = r.tester.Click("1")
		}
		if err == nil && navigationStateOf(r.state).pageDefault != 1 {
			err = fmt.Errorf("page 1 did not become selected")
		}
		if err == nil {
			err = r.tester.Click("Previous")
		}
		if err == nil && navigationStateOf(r.state).pageDefault != 1 {
			err = fmt.Errorf("disabled Previous changed the page")
		}
		if err != nil {
			fail("disabled", err)
		} else if path, err := capture(r, "disabled", "Previous", false, true); err != nil {
			fail("disabled", err)
		} else {
			e.States["disabled"] = auditCoverage{Status: "captured", Path: path}
		}
	}
	if page.ID == "textarea" {
		r, err := newRunner("Locked", true)
		if err == nil {
			err = auditFocus(r, "Meeting notes")
		}
		if err == nil {
			before := formsPageState(r.state).bio
			r.tester.Type(" changed")
			if formsPageState(r.state).bio != before {
				err = fmt.Errorf("read-only Meeting notes accepted typed input")
			}
		}
		if err != nil {
			fail("read_only", fmt.Errorf("example %q: %w", "Locked", err))
		} else if path, err := capture(r, "read_only", "Meeting notes", false, true); err != nil {
			fail("read_only", fmt.Errorf("example %q: %w", "Locked", err))
		} else {
			e.States["read_only"] = auditCoverage{Status: "captured", Path: path,
				Reason: "The gallery's Locked example is read-only, not disabled."}
		}
	}
	if page.ID == "checkbox" || page.ID == "switch" {
		state := "unchecked"
		if page.ID == "switch" {
			state = "off"
		}
		r, err := newRunner(example, true)
		if err == nil {
			err = r.tester.Click(spec.selected)
		}
		if err == nil && ((page.ID == "checkbox" && formsPageState(r.state).terms) ||
			(page.ID == "switch" && formsPageState(r.state).wifi)) {
			err = fmt.Errorf("one click did not change caller-owned state to %s", state)
		}
		if err != nil {
			fail(state, err)
		} else if path, err := capture(r, state, spec.selected, false, true); err != nil {
			fail(state, err)
		} else {
			e.States[state] = auditCoverage{Status: "captured", Path: path}
		}
	}
	if page.ID == "combobox" {
		r, err := newRunner(example, true)
		if err == nil {
			err = auditOpen(r, spec)
		}
		if err == nil {
			r.tester.Key(ui.Cmd, ui.KeyA)
			r.tester.Type("rust")
			if !r.tester.HasText("Rust") || r.tester.HasText("TypeScript") {
				err = fmt.Errorf("typing rust did not filter options: %q", r.tester.Texts())
			}
		}
		if err != nil {
			fail("filtered", err)
		} else if path, err := capture(r, "filtered", spec.target, true, true); err != nil {
			fail("filtered", err)
		} else {
			e.States["filtered"] = auditCoverage{Status: "captured", Path: path}
		}
	}
	if page.ID == "toolbar" {
		r, err := newRunner(example, true)
		if err == nil {
			var box ui.Rect
			var ok bool
			box, ok = r.tester.Find("Bold")
			if !ok {
				err = fmt.Errorf("toolbar button Bold is missing")
			} else {
				r.tester.Press(box.X+box.W/2, box.Y+box.H/2)
				r.tester.Frame()
			}
		}
		if err != nil {
			fail("pressed", err)
		} else if path, err := capture(r, "pressed", "Bold", false, true); err != nil {
			fail("pressed", err)
		} else {
			e.States["pressed"] = auditCoverage{Status: "captured", Path: path}
			box, _ := r.tester.Find("Bold")
			r.tester.Release(box.X+box.W/2, box.Y+box.H/2)
			if navigationStateOf(r.state).toolbarAction != "Bold" {
				fail("activated", fmt.Errorf("toolbar action Bold did not run"))
			} else if path, err := capture(r, "activated", "Bold", false, true); err != nil {
				fail("activated", err)
			} else {
				e.States["activated"] = auditCoverage{Status: "captured", Path: path}
			}
		}
		e.States["selected"] = auditCoverage{Status: "not_applicable", Reason: "This toolbar example uses action buttons, not persistent toggles."}
	}
	var alternateTitle, alternateTarget string
	switch page.ID {
	case "card":
		alternateTitle, alternateTarget = "Interactive", "Open project details"
	case "alert":
		alternateTitle, alternateTarget = "With action", "Retry deploy"
	case "emptystate":
		alternateTitle, alternateTarget = "With action", "Create search"
	}
	if alternateTitle != "" {
		for _, state := range []string{"hover", "keyboard_focus"} {
			r, err := newRunner(alternateTitle, true)
			if err == nil {
				if state == "hover" {
					err = auditHover(r, alternateTarget)
				} else {
					err = auditFocus(r, alternateTarget)
				}
			}
			if err != nil {
				fail(state, fmt.Errorf("example %q: %w", alternateTitle, err))
			} else if path, err := capture(r, state, alternateTarget, false, true); err != nil {
				fail(state, fmt.Errorf("example %q: %w", alternateTitle, err))
			} else {
				e.States[state] = auditCoverage{Status: "captured", Path: path}
			}
		}
	}
	if page.ID == "scrollarea" || page.ID == "progressbar" {
		r, err := newRunner(example, true)
		target := "Activity"
		state := "scrolled"
		if err == nil && page.ID == "scrollarea" {
			box, ok := r.tester.Find(target)
			if !ok {
				err = fmt.Errorf("scroll viewport %q is missing", target)
			} else {
				r.tester.Scroll(box.X+box.W/2, box.Y+box.H/2, 0, 120)
				r.tester.Frame()
				if cossStateOf(r.state, "scrollarea").Scroll[0].Y <= 0 {
					err = fmt.Errorf("scroll input did not change the caller-owned offset")
				}
			}
		}
		if page.ID == "progressbar" {
			target, state = "Advance", "value_changed"
			if err == nil {
				before := displayStateOf(r.state).Progress
				err = r.tester.Click(target)
				if err == nil && displayStateOf(r.state).Progress == before {
					err = fmt.Errorf("Advance did not change progress")
				}
			}
		}
		if err != nil {
			fail(state, err)
		} else if path, err := capture(r, state, target, false, true); err != nil {
			fail(state, err)
		} else {
			e.States[state] = auditCoverage{Status: "captured", Path: path}
		}
	}
	if page.ID == "table" {
		for _, variant := range []struct{ title, state, target string }{
			{"Selection", "selected", "Billing"}, {"Sorting", "sorted", "Status"},
		} {
			r, err := newRunner(variant.title, true)
			if err == nil {
				err = r.tester.Click(variant.target)
			}
			if err == nil {
				st := StateOf[tableState](r.state, "table")
				if variant.state == "selected" && !st.Table.Selection.HasKey("billing") {
					err = fmt.Errorf("Billing row was not selected")
				} else if variant.state == "sorted" && st.SortTable.Sort.Column != "status" {
					err = fmt.Errorf("Status column did not sort the table")
				}
			}
			if err != nil {
				fail(variant.state, fmt.Errorf("example %q: %w", variant.title, err))
			} else if path, err := capture(r, variant.state, variant.target, false, true); err != nil {
				fail(variant.state, fmt.Errorf("example %q: %w", variant.title, err))
			} else {
				e.States[variant.state] = auditCoverage{Status: "captured", Path: path}
			}
		}
	}
	if spec.motion {
		title := example
		if spec.motionExample != "" {
			title = spec.motionExample
		}
		if page.ID == "skeleton" || page.ID == "spinner" || page.ID == "progressbar" {
			quiet, quietErr := newRunner(title, true)
			if quietErr == nil {
				first := sha256.Sum256(quiet.tester.Image().Pix)
				quiet.tester.Advance(600 * time.Millisecond)
				if sha256.Sum256(quiet.tester.Image().Pix) != first {
					quietErr = fmt.Errorf("reduced-motion example %q changed during 600 ms", title)
				}
			}
			if quietErr != nil {
				fail("reduced_motion", quietErr)
			} else if path, quietErr := capture(quiet, "reduced_motion", "", false, true); quietErr != nil {
				fail("reduced_motion", quietErr)
			} else {
				e.States["reduced_motion"] = auditCoverage{Status: "captured", Path: path}
			}
		}
		r, err := newRunner(title, false)
		if err == nil {
			if spec.openByHover {
				err = auditHover(r, spec.open)
				if err == nil && page.ID == "tooltip" {
					r.tester.Advance(600 * time.Millisecond)
				}
				if err == nil && !r.tester.HasText(spec.openText) {
					err = fmt.Errorf("hover did not start the popup containing %q", spec.openText)
				}
			} else if spec.open != "" {
				err = auditOpen(r, spec)
			} else if page.ID == "switch" {
				err = r.tester.Click(spec.selected)
				if err == nil && formsPageState(r.state).wifi {
					err = fmt.Errorf("switch click did not start the off transition")
				}
			} else if spec.selected != "" {
				err = auditSelect(r, page.ID, spec)
			}
		}
		if err != nil {
			fail("motion", fmt.Errorf("example %q: %w", title, err))
		} else {
			var hashes [][32]byte
			var paths []string
			steps := []struct {
				name  string
				delay time.Duration
			}{{"motion-050ms", 50 * time.Millisecond}, {"motion-120ms", 70 * time.Millisecond}, {"motion-600ms", 480 * time.Millisecond}}
			if page.ID == "menu" {
				steps[0].name, steps[0].delay = "motion-025ms", 25*time.Millisecond
				steps[1].name, steps[1].delay = "motion-060ms", 35*time.Millisecond
				steps[2].name, steps[2].delay = "motion-300ms", 240*time.Millisecond
			}
			if page.ID == "skeleton" {
				const period = 2 * time.Second
				phase := time.Duration(r.tester.Now().UnixNano() % int64(period))
				r.tester.Advance(period - phase)
				steps[0].name, steps[0].delay = "motion-250ms", 250*time.Millisecond
				steps[1].name, steps[1].delay = "motion-700ms", 450*time.Millisecond
				steps[2].name, steps[2].delay = "motion-1200ms", 500*time.Millisecond
			}
			for _, step := range steps {
				r.tester.Advance(step.delay)
				frame := r.tester.Image()
				hashes = append(hashes, sha256.Sum256(frame.Pix))
				path, captureErr := capture(r, step.name, spec.target, spec.open != "", false)
				if captureErr != nil {
					err = captureErr
					break
				}
				paths = append(paths, path)
			}
			if err != nil {
				fail("motion", err)
			} else if hashes[0] == hashes[1] || hashes[1] == hashes[2] || hashes[0] == hashes[2] {
				fail("motion", fmt.Errorf("two or more motion frames have identical pixels"))
			} else {
				reason := fmt.Sprintf("See %s, %s and %s; reduced motion is off.", steps[0].name, steps[1].name, steps[2].name)
				if page.ID == "select" {
					reason += " The trigger chevron moves; this does not establish a popup fade."
				}
				if page.ID == "skeleton" {
					reason += " Times are measured from the start of the 2-second animation cycle."
				}
				e.States["motion"] = auditCoverage{Status: "captured", Path: paths[0], Reason: reason}
			}
		}
	}
	return e
}

func newAuditRunner(page ComponentPage, title string, dark, reduced bool) (*auditRunner, error) {
	var example ComponentExample
	for _, candidate := range page.Examples {
		if candidate.Title == title {
			example = candidate
			break
		}
	}
	if example.Demo == nil {
		return nil, fmt.Errorf("example %q has no demo", title)
	}
	s := New()
	s.Selected, s.Dark = page.ID, dark
	r := &auditRunner{state: s, example: example}
	r.tester = ui.NewTester(func(c *ui.Context) {
		applyAppearance(c, s)
		overlaysStateOf(s).toasts.View(c)
		palette := theme.Colors(c)
		ui.Column(c).Fill().Padding(100, 0, 0, 0).
			Background(palette.Background).TextColor(palette.Foreground).
			Children(func() { demoPanel(c, "audit panel", func() { example.Demo(c, s) }) })
	}, auditWidth, auditHeight)
	r.tester.SetPreferences(ui.Preferences{TextScale: 1, ReduceMotion: reduced})
	initial := r.tester.Image()
	r.initial = image.NewRGBA(initial.Bounds())
	copy(r.initial.Pix, initial.Pix)
	return r, nil
}

func auditHover(r *auditRunner, target string) error {
	box, ok := r.tester.Find(target)
	if !ok || box.W <= 0 || box.H <= 0 || box.X < 0 || box.Y < 0 || box.X+box.W > auditWidth || box.Y+box.H > auditHeight {
		return fmt.Errorf("hover target %q is missing or outside canvas: %+v", target, box)
	}
	x, y := box.X+box.W/2, box.Y+box.H/2
	switch r.state.Selected {
	case "autocomplete", "combobox", "select":
		// These examples put a separate label with the same name above the editor.
		y = box.Y + box.H + 16
	case "slider":
		x = box.X + (auditWidth-2*box.X)*float32(formsPageState(r.state).volume/100)
		y = box.Y + box.H + 14
	}
	r.tester.Move(x, y)
	r.tester.Frame()
	return nil
}

func auditFocus(r *auditRunner, target string) error {
	for range 32 {
		if r.tester.Focused(target) {
			return nil
		}
		r.tester.Key(0, ui.KeyTab)
	}
	return fmt.Errorf("keyboard Tab did not focus %q", target)
}

func auditOpen(r *auditRunner, spec auditSpec) error {
	if spec.openByHover {
		if err := auditHover(r, spec.open); err != nil {
			return err
		}
		r.tester.Advance(700 * time.Millisecond)
	} else if err := r.tester.Click(spec.open); err != nil {
		return err
	}
	if !r.tester.HasText(spec.openText) {
		return fmt.Errorf("opening %q did not show %q; texts: %q", spec.open, spec.openText, r.tester.Texts())
	}
	return nil
}

func auditDisabledCalendar(r *auditRunner) (string, error) {
	st := StateOf[datesState](r.state, "dates")
	before := st.Cal.Value
	for day := st.Min; !day.After(st.Max); day = day.AddDate(0, 0, 1) {
		if !datesWeekend(day) || day.Month() != before.Month() {
			continue
		}
		label := day.Format("Monday, January 2, 2006")
		box, ok := r.tester.Find(label)
		if !ok || box.Y < 0 || box.Y+box.H > auditHeight {
			continue
		}
		if err := r.tester.Click(label); err != nil {
			return "", err
		}
		if !st.Cal.Value.Equal(before) {
			return "", fmt.Errorf("disabled weekend %q changed the selected date", label)
		}
		return label, nil
	}
	return "", fmt.Errorf("calendar has no visible disabled weekend within its bounds")
}

func auditDisabledSelect(r *auditRunner, spec auditSpec) (string, error) {
	if err := auditOpen(r, spec); err != nil {
		return "", err
	}
	box, ok := r.tester.Find("Banana")
	if !ok {
		return "", fmt.Errorf("select popup has no scrollable Banana row")
	}
	r.tester.Scroll(box.X+box.W/2, box.Y+box.H/2, 0, 100)
	r.tester.Frame()
	if err := r.tester.Click("Dragon fruit"); err != nil {
		return "", err
	}
	st := selectionPageState(r.state)
	if st.fruitSelect.List.SelectedKey != "apple" || !st.fruitSelect.Open {
		return "", fmt.Errorf("disabled Dragon fruit changed the selection or closed the popup")
	}
	return "Dragon fruit", nil
}

func auditPickerSelect(r *auditRunner, spec auditSpec) error {
	st := StateOf[datesState](r.state, "dates")
	previous := st.Picker
	if err := auditOpen(r, spec); err != nil {
		return err
	}
	for day := st.Min; !day.After(st.Max); day = day.AddDate(0, 0, 1) {
		if datesWeekend(day) || day.Equal(previous) || day.Month() != previous.Month() {
			continue
		}
		label := "Deploy " + day.Format("Monday, January 2, 2006")
		box, ok := r.tester.Find(label)
		if !ok || box.Y < 0 || box.Y+box.H > auditHeight {
			continue
		}
		if err := r.tester.Click(label); err != nil {
			return err
		}
		if !st.Picker.Equal(day) || st.PickerOpen {
			return fmt.Errorf("date picker click %q did not commit the date and close the popup", label)
		}
		if err := auditOpen(r, spec); err != nil {
			return err
		}
		r.selectedTarget = label
		return nil
	}
	return fmt.Errorf("date picker has no visible alternative selectable day")
}

func auditSelect(r *auditRunner, id string, spec auditSpec) error {
	r.selectedTarget = spec.selected
	if id == "calendar" {
		st := StateOf[datesState](r.state, "dates")
		previous := st.CalDefault.Value
		first := time.Date(previous.Year(), previous.Month(), 1, 0, 0, 0, 0, previous.Location())
		for day := first; day.Month() == previous.Month(); day = day.AddDate(0, 0, 1) {
			if day.Equal(previous) {
				continue
			}
			label := day.Format("Monday, January 2, 2006")
			box, ok := r.tester.Find(label)
			if !ok || box.Y < 0 || box.Y+box.H > auditHeight {
				continue
			}
			if err := r.tester.Click(label); err != nil {
				return err
			}
			if !st.CalDefault.Value.Equal(day) {
				return fmt.Errorf("calendar click %q did not change the caller-owned date", label)
			}
			r.selectedTarget = label
			return nil
		}
		return fmt.Errorf("calendar has no visible alternative selectable day")
	}
	if id == "datepicker" {
		return auditPickerSelect(r, spec)
	}
	if id == "slider" {
		before := formsPageState(r.state).volume
		if err := auditFocus(r, spec.selected); err != nil {
			return err
		}
		r.tester.Key(0, ui.KeyRight)
		if formsPageState(r.state).volume == before {
			return fmt.Errorf("slider Right key did not change the caller-owned value")
		}
		return nil
	}
	if id == "select" || id == "autocomplete" || id == "combobox" {
		if err := auditOpen(r, spec); err != nil {
			return err
		}
	}
	if err := r.tester.Click(spec.selected); err != nil {
		return err
	}
	if id == "checkbox" || id == "switch" {
		if (id == "checkbox" && formsPageState(r.state).terms) ||
			(id == "switch" && formsPageState(r.state).wifi) {
			return fmt.Errorf("first click did not change %q to off", spec.selected)
		}
		if err := r.tester.Click(spec.selected); err != nil {
			return err
		}
	}
	valid := false
	switch id {
	case "autocomplete":
		valid = selectionPageState(r.state).auto.List.SelectedKey == "tokyo"
	case "checkbox":
		valid = formsPageState(r.state).terms
	case "checkboxgroup":
		valid = formsPageState(r.state).channels[1]
	case "combobox":
		valid = selectionPageState(r.state).combo.List.SelectedKey == "rust"
	case "menu":
		valid = slices.Contains(navigationStateOf(r.state).menuSelected, "Reports")
	case "pagination":
		valid = navigationStateOf(r.state).pageDefault == 3
	case "radiogroup":
		valid = formsPageState(r.state).plan == "free"
	case "segmentedcontrol":
		valid = cossStateOf(r.state, "segment-billing").Choice == 1
	case "select":
		valid = selectionPageState(r.state).fruitSelect.List.SelectedKey == "banana"
	case "sidebar":
		valid = cossStateOf(r.state, id).Name == "projects"
	case "switch":
		valid = formsPageState(r.state).wifi
	case "tabs":
		valid = navigationStateOf(r.state).tabsPrimary == 1
	case "togglebutton":
		valid = buttonsStateOf(r.state).DefaultOff
	case "togglebuttongroup":
		valid = buttonsStateOf(r.state).AlignLeft
	}
	if !valid {
		return fmt.Errorf("%q did not produce the expected caller-owned selection", spec.selected)
	}
	return nil
}

func writeAuditPNG(path string, r *auditRunner, full bool) (image.Rectangle, error) {
	img := r.tester.Image()
	panel, ok := r.tester.Find("audit panel")
	if !ok || panel.W != auditWidth || panel.H <= 0 {
		return image.Rectangle{}, fmt.Errorf("preview panel has invalid bounds: %+v", panel)
	}
	bounds := image.Rect(0, int(math.Floor(float64(panel.Y))), auditWidth, int(math.Ceil(float64(panel.Y+panel.H))))
	if !bounds.In(img.Bounds()) {
		return image.Rectangle{}, fmt.Errorf("preview panel is clipped by the canvas: %+v", panel)
	}
	if full {
		// Centered modals and toast cards can be cropped to their labeled
		// component. Edge panels need the entire window height.
		label := map[string]string{
			"alertdialog": "Delete project?", "dialog": "Rename workspace", "toast": "Draft saved",
		}[r.state.Selected]
		if box, found := r.tester.Find(label); label != "" && found && box.H > 0 {
			bounds.Min.Y = max(0, int(math.Floor(float64(box.Y)))-20)
			bounds.Max.Y = min(auditHeight, int(math.Ceil(float64(box.Y+box.H)))+20)
		} else {
			for y := 0; y < auditHeight; y++ {
				start, end := y*img.Stride, y*img.Stride+auditWidth*4
				if !bytes.Equal(img.Pix[start:end], r.initial.Pix[start:end]) {
					bounds.Min.Y = min(bounds.Min.Y, max(0, y-16))
					bounds.Max.Y = max(bounds.Max.Y, min(auditHeight, y+17))
				}
			}
		}
	}
	crop := img.SubImage(bounds)
	var buf bytes.Buffer
	if err := png.Encode(&buf, crop); err != nil {
		return image.Rectangle{}, err
	}
	decoded, err := png.Decode(bytes.NewReader(buf.Bytes()))
	if err != nil || decoded.Bounds().Dx() != auditWidth || decoded.Bounds().Dy() == 0 {
		return image.Rectangle{}, fmt.Errorf("PNG is unreadable or has invalid dimensions: %v", err)
	}
	first := decoded.At(0, 0)
	varied := false
	for y := 0; y < decoded.Bounds().Dy() && !varied; y++ {
		for x := 0; x < decoded.Bounds().Dx(); x++ {
			if decoded.At(x, y) != first {
				varied = true
				break
			}
		}
	}
	if !varied {
		return image.Rectangle{}, fmt.Errorf("preview contains only one color")
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return image.Rectangle{}, err
	}
	if err := os.WriteFile(path, buf.Bytes(), 0o644); err != nil {
		return image.Rectangle{}, err
	}
	saved, err := os.ReadFile(path)
	if err != nil {
		return image.Rectangle{}, fmt.Errorf("read written PNG: %w", err)
	}
	if len(saved) == 0 || !bytes.Equal(saved, buf.Bytes()) {
		return image.Rectangle{}, fmt.Errorf("written PNG is empty or differs from the encoded image")
	}
	if _, err := png.DecodeConfig(bytes.NewReader(saved)); err != nil {
		return image.Rectangle{}, fmt.Errorf("written PNG is unreadable: %w", err)
	}
	return bounds, nil
}

func auditContactPairs(ids []string, components []auditComponent) ([]auditContactPair, [][]auditContactPair, error) {
	byPage := make(map[string]map[string]auditComponent, len(ids))
	for _, component := range components {
		if _, ok := cossSources[component.Page]; !ok {
			return nil, nil, fmt.Errorf("unexpected audit component %q", component.Page)
		}
		if component.Mode != "light" && component.Mode != "dark" {
			return nil, nil, fmt.Errorf("%s: unexpected mode %q", component.Page, component.Mode)
		}
		if byPage[component.Page] == nil {
			byPage[component.Page] = map[string]auditComponent{}
		}
		if _, exists := byPage[component.Page][component.Mode]; exists {
			return nil, nil, fmt.Errorf("%s/%s: duplicate mode entry", component.Page, component.Mode)
		}
		byPage[component.Page][component.Mode] = component
	}
	defaults := make([]auditContactPair, 0, len(ids))
	var states [][]auditContactPair
	for _, id := range ids {
		light, hasLight := byPage[id]["light"]
		dark, hasDark := byPage[id]["dark"]
		if !hasLight || !hasDark {
			return nil, nil, fmt.Errorf("%s: missing mode entry (light=%t, dark=%t)", id, hasLight, hasDark)
		}
		captures := make(map[string]map[string]auditCapture, 2)
		for _, component := range []auditComponent{light, dark} {
			captures[component.Mode] = make(map[string]auditCapture, len(component.Captures))
			for _, capture := range component.Captures {
				if _, exists := captures[component.Mode][capture.State]; exists {
					return nil, nil, fmt.Errorf("%s/%s: duplicate %q capture", id, component.Mode, capture.State)
				}
				if capture.Path == "" || capture.Width != auditWidth || capture.Height <= 0 {
					return nil, nil, fmt.Errorf("%s/%s/%s: empty or invalid capture", id, component.Mode, capture.State)
				}
				captures[component.Mode][capture.State] = capture
			}
			capture, exists := captures[component.Mode]["default"]
			coverage := component.States["default"]
			if !exists || coverage.Status != "captured" || coverage.Path != capture.Path {
				return nil, nil, fmt.Errorf("%s/%s: missing nonempty captured default", id, component.Mode)
			}
		}
		if light.Slug != dark.Slug || light.Slug != cossSources[id] {
			return nil, nil, fmt.Errorf("%s: catalog slug %q does not match light %q and dark %q", id, cossSources[id], light.Slug, dark.Slug)
		}
		defaultPair := auditContactPair{page: id, slug: light.Slug,
			light: captures["light"]["default"], dark: captures["dark"]["default"]}
		defaults = append(defaults, defaultPair)
		componentStates := []auditContactPair{defaultPair}
		for _, capture := range light.Captures {
			other, ok := captures["dark"][capture.State]
			if !ok {
				return nil, nil, fmt.Errorf("%s: %s capture missing in dark mode", id, capture.State)
			}
			if capture.State != "default" {
				componentStates = append(componentStates, auditContactPair{
					page: id, slug: light.Slug, light: capture, dark: other,
				})
			}
		}
		for _, capture := range dark.Captures {
			if _, ok := captures["light"][capture.State]; !ok {
				return nil, nil, fmt.Errorf("%s: %s capture missing in light mode", id, capture.State)
			}
		}
		if len(componentStates) > 1 {
			states = append(states, componentStates)
		}
	}
	return defaults, states, nil
}

const (
	auditSheetMargin    = 24
	auditSheetWidth     = 2*auditWidth + 3*auditSheetMargin
	auditSheetHeader    = 76
	auditSheetRowHeader = 50
	auditSheetRowGap    = 20
	auditSheetMaxHeight = 2400
	auditSheetMaxRows   = 4
)

func writeAuditContactSheets(dir string, defaults []auditContactPair, states [][]auditContactPair) ([]auditContactSheet, error) {
	sheets := []auditContactSheet{}
	sheets, err := appendAuditContactSheets(dir, "defaults", defaults, sheets)
	if err != nil {
		return sheets, err
	}
	for _, group := range states {
		sheets, err = appendAuditContactSheets(dir, "states", group, sheets)
		if err != nil {
			return sheets, err
		}
	}
	return sheets, nil
}

func appendAuditContactSheets(
	dir, kind string,
	pairs []auditContactPair,
	sheets []auditContactSheet,
) ([]auditContactSheet, error) {
	for page := 1; len(pairs) > 0; page++ {
		height, count := auditSheetHeader, 0
		for count < len(pairs) && count < auditSheetMaxRows {
			pair := pairs[count]
			rowHeight := auditSheetRowHeader + max(pair.light.Height, pair.dark.Height) + auditSheetRowGap
			if count > 0 && height+rowHeight > auditSheetMaxHeight {
				break
			}
			height += rowHeight
			count++
		}
		name := fmt.Sprintf("defaults-%02d.png", page)
		if kind == "states" {
			name = filepath.Join("states", fmt.Sprintf("%s-%02d.png", pairs[0].slug, page))
		}
		path := filepath.Join(dir, "contactsheets", name)
		sheet, err := writeAuditContactSheet(path, kind, page, height, pairs[:count])
		if err != nil {
			return sheets, err
		}
		sheets = append(sheets, sheet)
		pairs = pairs[count:]
	}
	return sheets, nil
}

func writeAuditContactSheet(
	path, kind string,
	page, height int,
	pairs []auditContactPair,
) (auditContactSheet, error) {
	title := fmt.Sprintf("COSS native - defaults - page %d", page)
	if kind == "states" {
		title = fmt.Sprintf("COSS native - states: %s / %s - page %d", pairs[0].page, pairs[0].slug, page)
	}
	img := ui.Render(func(c *ui.Context) {
		ui.Column(c).Fill().Background(ui.RGB(246, 247, 249)).Children(func() {
			auditSheetText(c, auditSheetMargin, 10, title)
			auditSheetText(c, auditSheetMargin, 35, "Actual gallery captures at native scale; light and dark are shown side by side.")
			y := auditSheetHeader
			for _, pair := range pairs {
				label := fmt.Sprintf("native: %s    coss: %s    state: %s", pair.page, pair.slug, pair.light.State)
				auditSheetText(c, auditSheetMargin, y, label)
				auditSheetText(c, auditSheetMargin, y+23, "LIGHT  |  example: "+auditExampleLabel(pair.light))
				auditSheetText(c, 2*auditSheetMargin+auditWidth, y+23, "DARK  |  example: "+auditExampleLabel(pair.dark))
				y += auditSheetRowHeader + max(pair.light.Height, pair.dark.Height) + auditSheetRowGap
			}
		})
	}, auditSheetWidth, height, 1)

	sheet := auditContactSheet{Kind: kind, Page: page, Path: path, Width: auditSheetWidth, Height: height}
	y := auditSheetHeader
	for _, pair := range pairs {
		row := auditContactRow{Page: pair.page, Slug: pair.slug, State: pair.light.State}
		var err error
		row.Light, err = drawAuditContactCell(img, pair.light, auditSheetMargin, y+auditSheetRowHeader)
		if err != nil {
			return sheet, fmt.Errorf("%s/light/%s: %w", pair.page, pair.light.State, err)
		}
		row.Dark, err = drawAuditContactCell(img, pair.dark, 2*auditSheetMargin+auditWidth, y+auditSheetRowHeader)
		if err != nil {
			return sheet, fmt.Errorf("%s/dark/%s: %w", pair.page, pair.dark.State, err)
		}
		sheet.Rows = append(sheet.Rows, row)
		y += auditSheetRowHeader + max(pair.light.Height, pair.dark.Height) + auditSheetRowGap
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return sheet, fmt.Errorf("encode %s: %w", path, err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return sheet, err
	}
	if err := os.WriteFile(path, buf.Bytes(), 0o644); err != nil {
		return sheet, err
	}
	saved, err := os.ReadFile(path)
	if err != nil {
		return sheet, fmt.Errorf("read contact sheet %s: %w", path, err)
	}
	config, err := png.DecodeConfig(bytes.NewReader(saved))
	if err != nil || len(saved) == 0 || config.Width != sheet.Width || config.Height != sheet.Height {
		return sheet, fmt.Errorf("contact sheet %s is empty, unreadable or has invalid dimensions: %v", path, err)
	}
	return sheet, nil
}

func auditSheetText(c *ui.Context, x, y int, text string) {
	ui.Text(c, text).Absolute().Left(float32(x)).Top(float32(y)).
		Font(fonts.Family).FontSize(16).TextColor(ui.RGB(34, 42, 53)).SingleLine()
}

func auditExampleLabel(capture auditCapture) string {
	if capture.ExamplePage == "" {
		return capture.Example
	}
	return capture.Example + " (" + capture.ExamplePage + ")"
}

func drawAuditContactCell(img *image.RGBA, capture auditCapture, x, y int) (auditContactCell, error) {
	data, err := os.ReadFile(capture.Path)
	if err != nil {
		return auditContactCell{}, err
	}
	src, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		return auditContactCell{}, fmt.Errorf("decode %s: %w", capture.Path, err)
	}
	if src.Bounds().Dx() != capture.Width || src.Bounds().Dy() != capture.Height {
		return auditContactCell{}, fmt.Errorf("capture %s dimensions do not match the manifest", capture.Path)
	}
	bounds := image.Rect(x, y, x+capture.Width, y+capture.Height)
	draw.Draw(img, bounds.Inset(-1), image.NewUniform(color.RGBA{R: 198, G: 203, B: 211, A: 255}), image.Point{}, draw.Src)
	draw.Draw(img, bounds, src, src.Bounds().Min, draw.Src)
	return auditContactCell{
		Path: capture.Path, Example: capture.Example, ExamplePage: capture.ExamplePage, Target: capture.Target,
		Bounds: auditContactBounds{X: x, Y: y, Width: capture.Width, Height: capture.Height},
	}, nil
}
