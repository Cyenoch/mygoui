package form

import (
	"path/filepath"
	"strings"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// UploadStatus is the caller's state of one listed file.
type UploadStatus uint8

const (
	UploadReady UploadStatus = iota
	UploadUploading
	UploadDone
	UploadFailed
)

// UploadFile is one row of the caller-owned file list. Detail is the muted
// line under the name, such as a size or an error; Progress runs from 0 to 1
// while uploading.
type UploadFile struct {
	Name     string
	Detail   string
	Status   UploadStatus
	Progress float64
}

// UploadOptions configures a file drop zone. The library opens no dialog and
// reads no file: OnFiles receives the dropped paths, and OnChoose together
// with ChooseLabel adds a button for the application's own picker (see
// PickFiles).
//
// Accept lists file extensions, with or without the dot; empty accepts
// every file. Without Multiple only the first accepted file is reported.
// Files renders below the zone, each with a remove button when OnRemove is
// set.
type UploadOptions struct {
	Title       string
	Description string
	ChooseLabel string
	Accept      []string
	Multiple    bool
	Disabled    bool
	Files       []UploadFile
	OnFiles     func(paths []string)
	OnChoose    func()
	OnRemove    func(index int)
}

// Upload is a dashed zone that takes files dropped from other applications
// and highlights while they hover over it.
func Upload(c *ui.Context, o UploadOptions) ui.Element {
	root := ui.Column(c).FillWidth().MinWidth(0).Gap(theme.Space12)
	root.Children(func() {
		uploadZone(c, o)
		if len(o.Files) == 0 {
			return
		}
		ui.Column(c).FillWidth().MinWidth(0).Gap(theme.Space8).Children(func() {
			for i, f := range o.Files {
				uploadFileRow(c.Key(i), f, o.Disabled, func() { o.OnRemove(i) }, o.OnRemove != nil)
			}
		})
	})
	return root
}

func uploadZone(c *ui.Context, o UploadOptions) {
	t := theme.Colors(c)
	zone := ui.Column(c).FillWidth().MinWidth(0).AlignItems(ui.Center).Justify(ui.Center).
		Gap(theme.Space12).Padding(theme.Space24).Radius(theme.RadiusScale().XL).
		BorderStyle(ui.BorderDashed).Role(ui.RoleGroup).Label(o.Title)
	zone = style.Disabled(zone, o.Disabled)
	// A disabled zone is no drop target, so MyGo hands it nothing.
	if zone.FileDragOver() {
		zone.Border(theme.Hairline, t.Ring).Background(t.Accent)
	} else {
		zone.Border(theme.Hairline, t.Input)
	}
	if paths := uploadAccepted(zone.DroppedFiles(), o.Accept, o.Multiple); len(paths) > 0 && o.OnFiles != nil {
		o.OnFiles(paths)
	}
	zone.Children(func() {
		icon := ui.Box(c).Size(theme.Scale(c, 36), theme.Scale(c, 36)).Center().Shrink(0).
			Radius(theme.RadiusScale().MD).Border(theme.Hairline, t.Border).Background(t.Card)
		style.Shadow(icon, c, theme.SurfaceShadows(c.Theme().Dark)).Children(func() {
			ui.Icon(c, lucide.Upload()).Size(theme.Scale(c, 18), theme.Scale(c, 18)).TextColor(t.Foreground)
		})
		if o.Title != "" || o.Description != "" {
			ui.Column(c).MinWidth(0).AlignItems(ui.Center).Gap(theme.Space2).Children(func() {
				if o.Title != "" {
					style.Text(ui.Text(c, o.Title), c, theme.SM()).FontWeight(theme.WeightMedium).
						TextColor(t.Foreground).TextAlign(ui.Center)
				}
				if o.Description != "" {
					style.Text(ui.Text(c, o.Description), c, theme.XS()).
						TextColor(t.MutedForeground).TextAlign(ui.Center)
				}
			})
		}
		if o.OnChoose != nil && o.ChooseLabel != "" {
			mygoui.Button(c, mygoui.ButtonOptions{
				Label: o.ChooseLabel, Variant: mygoui.ButtonOutline, Size: theme.ButtonSmall,
				OnClick: o.OnChoose,
			}).AlignSelf(ui.Center)
		}
	})
}

// uploadAccepted keeps the paths whose extension is accepted, and only the
// first of them unless multiple.
func uploadAccepted(paths, accept []string, multiple bool) []string {
	var out []string
	for _, path := range paths {
		if !uploadAccepts(path, accept) {
			continue
		}
		out = append(out, path)
		if !multiple {
			break
		}
	}
	return out
}

func uploadAccepts(path string, accept []string) bool {
	if len(accept) == 0 {
		return true
	}
	ext := strings.TrimPrefix(filepath.Ext(path), ".")
	for _, a := range accept {
		if strings.EqualFold(ext, strings.TrimPrefix(a, ".")) {
			return true
		}
	}
	return false
}

func uploadFileRow(c *ui.Context, f UploadFile, disabled bool, remove func(), removable bool) {
	t := theme.Colors(c)
	iconSize := theme.Scale(c, theme.IconMD)
	row := ui.Row(c).FillWidth().MinWidth(0).AlignItems(ui.Center).Gap(theme.Space12).
		Padding(theme.Space8, theme.Space12).Radius(theme.RadiusScale().LG).
		Border(theme.Hairline, t.Border)
	row.Children(func() {
		ui.Icon(c, lucide.File()).Size(iconSize, iconSize).Shrink(0).TextColor(t.MutedForeground)
		ui.Column(c).Grow(1).MinWidth(0).Gap(theme.Space2).Children(func() {
			style.Text(ui.Text(c, f.Name), c, theme.SM()).FontWeight(theme.WeightMedium).
				TextColor(t.Foreground).SingleLine().MinWidth(0)
			if f.Detail != "" {
				ink := t.MutedForeground
				if f.Status == UploadFailed {
					ink = t.DestructiveForeground
				}
				style.Text(ui.Text(c, f.Detail), c, theme.XS()).TextColor(ink).SingleLine().MinWidth(0)
			}
			if f.Status == UploadUploading {
				mygoui.ProgressBar(c, mygoui.ProgressOptions{
					Value: f.Progress, Max: 1, Size: mygoui.ProgressSmall,
				})
			}
		})
		switch f.Status {
		case UploadDone:
			ui.Icon(c, mygoui.StatusIcon(mygoui.ColorSuccess)).Size(iconSize, iconSize).Shrink(0).
				TextColor(t.SuccessForeground)
		case UploadFailed:
			ui.Icon(c, mygoui.StatusIcon(mygoui.ColorDanger)).Size(iconSize, iconSize).Shrink(0).
				TextColor(t.DestructiveForeground)
		}
		if removable {
			mygoui.Button(c, mygoui.ButtonOptions{
				Label: locale.M().Remove + " " + f.Name, Icon: lucide.X(), IconOnly: true,
				Variant: mygoui.ButtonGhost, Size: theme.ButtonExtraSmall,
				Disabled: disabled, OnClick: remove,
			})
		}
	})
}
