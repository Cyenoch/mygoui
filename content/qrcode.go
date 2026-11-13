package content

import (
	"fmt"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
	"rsc.io/qr"
)

// QRLevel is how much of a QR code may be damaged and still read.
type QRLevel uint8

const (
	// QRMedium recovers from about 15% damage.
	QRMedium QRLevel = iota
	// QRLow recovers from about 7%, for the smallest code.
	QRLow
	// QRQuartile recovers from about 25%.
	QRQuartile
	// QRHigh recovers from about 30%.
	QRHigh
)

// QRCodeOptions configures a QR code. Size is the side of the square in
// DIPs, quiet zone included (default 160). Foreground and Background
// default to the light palette's foreground and background in both
// appearances, since many scanners cannot read light modules on dark. Icon builds
// content over the middle of the code and raises Level to QRHigh so the
// covered modules stay recoverable. Label is the accessible name and
// defaults to Value.
type QRCodeOptions struct {
	Value      string
	Size       float32
	Level      QRLevel
	Foreground ui.Color
	Background ui.Color
	Icon       func()
	Label      string
}

const (
	qrDefaultSize = 160
	// qrQuietModules is the standard margin around the symbol.
	qrQuietModules = 4
	// qrIconRatio is the icon's side relative to the code's; QRHigh
	// tolerates the modules it covers.
	qrIconRatio = 0.22
)

// qrState caches the encoded matrix by value and level, and the path of
// its modules by the box and scale it was last painted in.
type qrState struct {
	encoded bool
	value   string
	level   QRLevel
	err     error
	side    int
	modules []bool

	path  *ui.Path
	rect  ui.Rect
	scale float32
}

type qrStateKey struct{}

// QRCode draws Value as a QR code with crisp modules. When Value cannot be
// encoded, as when it exceeds what a QR code holds at the level, it returns
// the error and an empty element instead of a partial code.
func QRCode(c *ui.Context, o QRCodeOptions) (ui.Element, error) {
	root := ui.Box(c).Shrink(0)
	level := o.Level
	if o.Icon != nil {
		level = QRHigh
	}
	st := ui.Local(root, qrStateKey{}, func() qrState { return qrState{} })
	if !st.encoded || st.value != o.Value || st.level != level {
		*st = encodeQR(o.Value, level)
	}
	if st.err != nil {
		return root, st.err
	}
	t := theme.ColorsForMode(false)
	size, fg, bg := o.Size, o.Foreground, o.Background
	if size <= 0 {
		size = qrDefaultSize
	}
	if fg.A == 0 {
		fg = t.Foreground
	}
	if bg.A == 0 {
		bg = t.Background
	}
	label := o.Label
	if label == "" {
		label = o.Value
	}
	root.Size(size, size).Radius(theme.RadiusScale().LG).Clip().Background(bg).
		Role(ui.RoleImage).Label(label)
	root.Draw(func(p *ui.Painter, r ui.Rect) {
		if st.path == nil || st.rect != r || st.scale != p.Scale() {
			st.path, st.rect, st.scale = qrPath(st.modules, st.side, r, p.Scale()), r, p.Scale()
		}
		p.FillPath(st.path, fg)
	})
	if o.Icon != nil {
		side := size * qrIconRatio
		root.Children(func() {
			ui.Box(c).Absolute().Left((size-side)/2).Top((size-side)/2).Size(side, side).
				Radius(theme.RadiusScale().MD).Clip().Background(bg).Center().Children(o.Icon)
		})
	}
	return root, nil
}

func encodeQR(value string, level QRLevel) qrState {
	st := qrState{encoded: true, value: value, level: level}
	code, err := qr.Encode(value, [...]qr.Level{QRMedium: qr.M, QRLow: qr.L, QRQuartile: qr.Q, QRHigh: qr.H}[level])
	if err != nil {
		st.err = fmt.Errorf("content: encode %d-byte QR value: %w", len(value), err)
		return st
	}
	st.side = code.Size
	st.modules = make([]bool, code.Size*code.Size)
	for y := range code.Size {
		for x := range code.Size {
			st.modules[y*code.Size+x] = code.Black(x, y)
		}
	}
	return st
}

// qrPath outlines each horizontal run of dark modules inside r, leaving the
// quiet zone clear. Module edges snap to device pixels so neighbouring
// modules neither gap nor blur.
func qrPath(modules []bool, side int, r ui.Rect, scale float32) *ui.Path {
	unit := r.W / float32(side+2*qrQuietModules)
	snap := func(origin float32, k int) float32 {
		v := origin + float32(k+qrQuietModules)*unit
		return float32(int(v*scale+0.5)) / scale
	}
	path := new(ui.Path)
	for y := range side {
		top, bottom := snap(r.Y, y), snap(r.Y, y+1)
		for x := 0; x < side; {
			if !modules[y*side+x] {
				x++
				continue
			}
			end := x
			for end < side && modules[y*side+end] {
				end++
			}
			left, right := snap(r.X, x), snap(r.X, end)
			path.MoveTo(left, top).LineTo(right, top).LineTo(right, bottom).LineTo(left, bottom).Close()
			x = end
		}
	}
	return path
}
