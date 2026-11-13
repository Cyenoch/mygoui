package overlay

import "github.com/Cyenoch/mygoui/theme"

// The motion tests time against the built-in theme's overlay durations.
var (
	popupMotionDuration  = theme.Default().Motion.Popup
	dialogMotionDuration = theme.Default().Motion.Dialog
	drawerMotionDuration = theme.Default().Motion.Drawer
)
