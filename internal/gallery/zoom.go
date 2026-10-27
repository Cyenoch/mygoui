package gallery

import "github.com/egoist/mygo/ui"

var zoomLevels = [...]float32{0.5, 0.67, 0.75, 0.8, 0.9, 1, 1.1, 1.25, 1.5, 1.75, 2, 2.5, 3}

const zoomDefault = 5

func zoom(c *ui.Context, s *State) {
	switch {
	case c.Shortcut(ui.Cmd, ui.KeyEqual), c.Shortcut(ui.Cmd|ui.Shift, ui.KeyEqual):
		s.zoomStep = min(s.zoomStep+1, len(zoomLevels)-1-zoomDefault)
	case c.Shortcut(ui.Cmd, ui.KeyMinus):
		s.zoomStep = max(s.zoomStep-1, -zoomDefault)
	case c.Shortcut(ui.Cmd, ui.Key0):
		s.zoomStep = 0
	}
	c.SetZoom(zoomLevels[zoomDefault+s.zoomStep])
}
