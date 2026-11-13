package gallery

import (
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
	"strings"

	"github.com/Cyenoch/mygoui/form"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// absorbedFormsState is the caller-owned state behind the absorbed inputs.
type absorbedFormsState struct {
	ready bool

	title, emptyTitle string
	titleCommits      int

	hotkey, clearable form.Hotkey

	rating, halfRating float64

	layout  string
	layouts []string

	files     []form.UploadFile
	pickError string
	statuses  []form.UploadFile
}

func absorbedFormsStateOf(s *State) *absorbedFormsState {
	st := StateOf[absorbedFormsState](s, "absorbed-forms")
	if !st.ready {
		st.ready = true
		st.title = "Quarterly report"
		st.hotkey = form.Hotkey{Mods: ui.Ctrl, Key: ui.KeyK}
		st.clearable = form.Hotkey{Mods: ui.Ctrl | ui.Shift, Key: ui.KeyP}
		st.rating, st.halfRating = 3, 2.5
		st.layout = "Dawn"
		st.layouts = []string{"grid"}
		st.statuses = []form.UploadFile{
			{Name: "brief.pdf", Detail: "1.2 MB", Status: form.UploadDone},
			{Name: "mockups.zip", Detail: "48 MB", Status: form.UploadUploading, Progress: 0.6},
			{Name: "notes.txt", Detail: "Upload failed", Status: form.UploadFailed},
		}
	}
	return st
}

func absorbedFormsDemo(body func(c *ui.Context, s *State, st *absorbedFormsState)) func(*ui.Context, *State) {
	return func(c *ui.Context, s *State) { body(c, s, absorbedFormsStateOf(s)) }
}

func init() {
	RegisterComponent(ComponentPage{
		ID: "editabletext", Title: "EditableText", Group: "Forms", Order: 230,
		Description: "Text that turns into an input of the same size when activated.",
		Examples: []ComponentExample{
			{Title: "Inline edit", Description: "Enter or blur commits; Escape drops the draft.",
				Demo: absorbedFormsDemo(editableTextBasicBody)},
			{Title: "Placeholder", Demo: absorbedFormsDemo(func(c *ui.Context, _ *State, st *absorbedFormsState) {
				form.EditableText(c, &st.emptyTitle, form.EditableTextOptions{
					Placeholder: "Untitled", AriaLabel: "Document title",
				})
			})},
			{Title: "Disabled", Demo: absorbedFormsDemo(func(c *ui.Context, _ *State, st *absorbedFormsState) {
				form.EditableText(c, &st.title, form.EditableTextOptions{AriaLabel: "Locked title", Disabled: true})
			})},
		},
		Code: `form.EditableText(c, &name, form.EditableTextOptions{
    Placeholder: "Untitled",
    AriaLabel:   "Project name",
    OnChange:    save,
})`,
		Params: []ComponentParam{
			{"Placeholder", "string", "\"\"", "Shown while the value is empty."},
			{"AriaLabel", "string", "\"\"", "Accessible name of the text and its input."},
			{"Disabled", "bool", "false", "Dims the text and stops editing."},
			{"FullWidth", "bool", "false", "Fills the parent instead of the default input width."},
			{"OnChange", "func()", "nil", "Runs after a commit that changed the value."},
		},
		Tokens: []string{"Foreground", "FieldPlaceholder", "Accent", "Input", "Ring"},
	})

	RegisterComponent(ComponentPage{
		ID: "hotkeyinput", Title: "HotkeyInput", Group: "Forms", Order: 240,
		Description: "Records a keyboard shortcut and shows it as key caps.",
		Examples: []ComponentExample{
			{Title: "Recorder", Description: "Click, then press a key with a modifier.",
				Demo: absorbedFormsDemo(func(c *ui.Context, _ *State, st *absorbedFormsState) {
					form.HotkeyInput(c, &st.hotkey, form.HotkeyInputOptions{
						Placeholder: "Press shortcut", AriaLabel: "Search shortcut",
					})
				})},
			{Title: "Clearable", Demo: absorbedFormsDemo(func(c *ui.Context, _ *State, st *absorbedFormsState) {
				form.HotkeyInput(c, &st.clearable, form.HotkeyInputOptions{
					Placeholder: "Press shortcut", AriaLabel: "Palette shortcut", Clearable: true,
				})
			})},
			{Title: "Disabled", Demo: absorbedFormsDemo(func(c *ui.Context, _ *State, st *absorbedFormsState) {
				form.HotkeyInput(c, &st.hotkey, form.HotkeyInputOptions{AriaLabel: "Locked shortcut", Disabled: true})
			})},
		},
		Code: `form.HotkeyInput(c, &shortcut, form.HotkeyInputOptions{
    Placeholder: "Press shortcut",
    AriaLabel:   "Search shortcut",
    Clearable:   true,
    OnChange:    save,
})`,
		Params: []ComponentParam{
			{"Placeholder", "string", "\"\"", "Shown while no shortcut is recorded."},
			{"AriaLabel", "string", "\"\"", "Accessible name of the recorder."},
			{"Clearable", "bool", "false", "Adds a clear button; Backspace and Delete clear too."},
			{"Disabled", "bool", "false", "Dims the control and stops recording."},
			{"FullWidth", "bool", "false", "Fills the parent instead of the default input width."},
			{"OnChange", "func()", "nil", "Runs after the shortcut is recorded or cleared."},
		},
		Tokens: []string{"Input", "Ring", "FieldPlaceholder", "Accent", "Foreground"},
	})

	RegisterComponent(ComponentPage{
		ID: "rate", Title: "Rate", Group: "Forms", Order: 250,
		Description: "A star rating set by pointer or arrow keys.",
		Examples: []ComponentExample{
			{Title: "Stars", Demo: absorbedFormsDemo(func(c *ui.Context, _ *State, st *absorbedFormsState) {
				rateBody(c, &st.rating, form.RateOptions{AriaLabel: "Rating"})
			})},
			{Title: "Half steps", Description: "A click on the current rating clears it.",
				Demo: absorbedFormsDemo(func(c *ui.Context, _ *State, st *absorbedFormsState) {
					rateBody(c, &st.halfRating, form.RateOptions{
						AriaLabel: "Half-step rating", AllowHalf: true, AllowClear: true,
					})
				})},
			{Title: "Read only", Demo: absorbedFormsDemo(func(c *ui.Context, _ *State, st *absorbedFormsState) {
				average := 4.5
				form.Rate(c, &average, form.RateOptions{AriaLabel: "Average rating", ReadOnly: true})
			})},
		},
		Code: `form.Rate(c, &rating, form.RateOptions{
    AriaLabel:  "Rating",
    AllowHalf:  true,
    AllowClear: true,
    OnChange:   save,
})`,
		Params: []ComponentParam{
			{"AriaLabel", "string", "\"\"", "Accessible name of the rating."},
			{"Count", "int", "5", "Number of stars."},
			{"AllowHalf", "bool", "false", "Steps by half a star."},
			{"AllowClear", "bool", "false", "A click on the current rating clears it."},
			{"ReadOnly", "bool", "false", "Shows the rating without taking input."},
			{"Disabled", "bool", "false", "Dims the stars and stops input."},
			{"OnChange", "func()", "nil", "Runs after the rating changed."},
		},
		Tokens: []string{"Warning", "Input", "Ring"},
	})

	RegisterComponent(ComponentPage{
		ID: "imageselect", Title: "ImageSelect", Group: "Forms", Order: 260,
		Description: "Picture tiles chosen like radio buttons or check boxes.",
		Examples: []ComponentExample{
			{Title: "Single", Description: "Arrows move the choice.",
				Demo: absorbedFormsDemo(imageSelectSingleBody)},
			{Title: "Multiple", Description: "Each tile toggles; Content draws the picture.",
				Demo: absorbedFormsDemo(imageSelectMultipleBody)},
		},
		Code: `form.ImageSelect(c, form.ImageSelectOptions{
    AriaLabel: "Layout",
    Items: []form.ImageSelectItem{
        {Value: "grid", Label: "Grid", Image: grid},
        {Value: "rows", Label: "Rows", Image: rows},
    },
    Value:    &layout,
    OnChange: save,
})`,
		Params: []ComponentParam{
			{"AriaLabel", "string", "\"\"", "Accessible name of the group."},
			{"Items", "[]form.ImageSelectItem", "nil", "Tiles: Value, Label, Image or Content, Disabled."},
			{"Value", "*string", "nil", "Single choice, like a radio group."},
			{"Values", "*[]string", "nil", "Multiple choice, like check boxes; wins over Value."},
			{"Width", "float32", "144", "Picture width in DIPs."},
			{"Height", "float32", "90", "Picture height in DIPs."},
			{"Disabled", "bool", "false", "Dims the tiles and stops input."},
			{"OnChange", "func()", "nil", "Runs after the choice changed."},
		},
		Tokens: []string{"Primary", "PrimaryForeground", "Border", "Ring", "MutedBackground", "MutedForeground", "Foreground"},
	})

	RegisterComponent(ComponentPage{
		ID: "upload", Title: "Upload", Group: "Forms", Order: 270,
		Description: "A drop zone over a caller-owned file list.",
		Examples: []ComponentExample{
			{Title: "Drop zone", Demo: absorbedFormsDemo(uploadZoneBody)},
			{Title: "File states", Description: "Uploading, done and failed rows.",
				Demo: absorbedFormsDemo(func(c *ui.Context, _ *State, st *absorbedFormsState) {
					formsFieldStack(c, func() {
						form.Upload(c, form.UploadOptions{
							Title: "Drop files here", Multiple: true, Files: st.statuses,
							OnFiles:  func(paths []string) { st.statuses = append(st.statuses, uploadRows(paths)...) },
							OnRemove: func(i int) { st.statuses = slices.Delete(st.statuses, i, i+1) },
						})
					})
				})},
		},
		Code: `accept := []string{"png", "jpg"}
form.Upload(c, form.UploadOptions{
    Title:       "Drop images here",
    Description: "PNG or JPG",
    ChooseLabel: "Browse",
    Accept:      accept,
    Multiple:    true,
    Files:       files,
    OnFiles:     addFiles,
    OnRemove:    removeFile,
    // PickFiles blocks: run it on a worker and publish through the window.
    OnChoose: func() {
        go func() {
            paths, err := form.PickFiles(form.PickFilesOptions{
                Parent: window, Accept: accept, Multiple: true,
            })
            window.Update(func() { pickErr = err; addFiles(paths) })
        }()
    },
})`,
		Params: []ComponentParam{
			{"Title", "string", "\"\"", "Zone heading and accessible name."},
			{"Description", "string", "\"\"", "Muted line under the title."},
			{"ChooseLabel", "string", "\"\"", "Label of the picker button, shown with OnChoose."},
			{"Accept", "[]string", "nil", "Accepted extensions, with or without the dot; empty accepts all."},
			{"Multiple", "bool", "false", "Reports every accepted file instead of the first."},
			{"Disabled", "bool", "false", "Dims the zone and takes no drops."},
			{"Files", "[]form.UploadFile", "nil", "Rows under the zone: Name, Detail, Status, Progress."},
			{"OnFiles", "func(paths []string)", "nil", "Receives the dropped paths."},
			{"OnChoose", "func()", "nil", "Opens the application's picker; see form.PickFiles."},
			{"OnRemove", "func(index int)", "nil", "Adds a remove button to each row."},
		},
		Tokens: []string{"Input", "Ring", "Accent", "Border", "Card", "Foreground", "MutedForeground", "SuccessForeground", "DestructiveForeground"},
	})
}

func editableTextBasicBody(c *ui.Context, _ *State, st *absorbedFormsState) {
	ui.Column(c).Gap(theme.Space8).AlignItems(ui.Center).Children(func() {
		form.EditableText(c, &st.title, form.EditableTextOptions{
			AriaLabel: "Report title", OnChange: func() { st.titleCommits++ },
		})
		formsMutedText(c, strconv.Itoa(st.titleCommits)+" commits")
	})
}

func rateBody(c *ui.Context, value *float64, o form.RateOptions) {
	ui.Row(c).Gap(theme.Space12).AlignItems(ui.Center).Children(func() {
		form.Rate(c, value, o)
		formsMutedText(c, strconv.FormatFloat(*value, 'f', -1, 64))
	})
}

func imageSelectSingleBody(c *ui.Context, _ *State, st *absorbedFormsState) {
	pictures := demoPictures()
	items := make([]form.ImageSelectItem, len(pictures))
	for i, p := range pictures {
		items[i] = form.ImageSelectItem{Value: p.Label, Label: p.Label, Image: p.Source}
	}
	items[len(items)-1].Disabled = true
	form.ImageSelect(c, form.ImageSelectOptions{AriaLabel: "Wallpaper", Items: items, Value: &st.layout}).
		Justify(ui.Center)
}

func imageSelectMultipleBody(c *ui.Context, _ *State, st *absorbedFormsState) {
	t := theme.Colors(c)
	// layoutTile sketches a layout as blocks of cols by rows.
	layoutTile := func(cols, rows int) func() {
		return func() {
			ui.Column(c).Fill().Padding(theme.Space12).Gap(theme.Space6).Children(func() {
				for range rows {
					ui.Row(c).FillWidth().Grow(1).Gap(theme.Space6).Children(func() {
						for range cols {
							ui.Box(c).Grow(1).Basis(0).AlignSelf(ui.Stretch).
								Radius(theme.RadiusScale().SM).Background(t.Border)
						}
					})
				}
			})
		}
	}
	form.ImageSelect(c, form.ImageSelectOptions{
		AriaLabel: "Views",
		Items: []form.ImageSelectItem{
			{Value: "grid", Label: "Grid", Content: layoutTile(3, 2)},
			{Value: "rows", Label: "Rows", Content: layoutTile(1, 3)},
			{Value: "columns", Label: "Columns", Content: layoutTile(3, 1)},
		},
		Values: &st.layouts, Width: 120, Height: 76,
	}).Justify(ui.Center)
}

func uploadZoneBody(c *ui.Context, s *State, st *absorbedFormsState) {
	accept := []string{"png", "jpg", "pdf"}
	add := func(paths []string) { st.files = append(st.files, uploadRows(paths)...) }
	o := form.UploadOptions{
		Title: "Drop files here", Description: "PNG, JPG or PDF",
		Accept: accept, Multiple: true, Files: st.files,
		OnFiles:  add,
		OnRemove: func(i int) { st.files = slices.Delete(st.files, i, i+1) },
	}
	// Browsers have no native dialog, and a view built without a window has
	// nothing to publish the result through.
	if window := s.Window; window != nil && runtime.GOOS != "js" {
		o.ChooseLabel = "Browse"
		o.OnChoose = func() {
			go func() {
				paths, err := form.PickFiles(form.PickFilesOptions{Parent: window, Accept: accept, Multiple: true})
				window.Update(func() {
					st.pickError = ""
					if err != nil {
						st.pickError = err.Error()
					}
					add(paths)
				})
			}()
		}
	}
	formsFieldStack(c, func() {
		form.Upload(c, o)
		form.FieldError(c, st.pickError)
	})
}

// uploadRows lists paths by file name, with the extension as the detail.
func uploadRows(paths []string) []form.UploadFile {
	rows := make([]form.UploadFile, len(paths))
	for i, path := range paths {
		rows[i] = form.UploadFile{
			Name:   filepath.Base(path),
			Detail: strings.ToUpper(strings.TrimPrefix(filepath.Ext(path), ".")),
			Status: form.UploadDone,
		}
	}
	return rows
}
