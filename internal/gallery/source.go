package gallery

import _ "embed"

//go:embed specimen.go
var specimenSource string

//go:embed navigation.go
var navigationSource string

//go:embed data_layout.go
var dataLayoutSource string

//go:embed datecharts.go
var dateChartSource string

//go:embed message.go
var messageSource string

var messagePreviewSource = specimenSource + "\n\n// File: message.go — compiled conversation specimen\n" + messageSource

var navigationPreviewSource = specimenSource + "\n\n// File: navigation.go — compiled sample implementation\n" + navigationSource
var dataLayoutPreviewSource = specimenSource + "\n\n// File: data_layout.go — compiled sample implementation\n" + dataLayoutSource
var dateChartPreviewSource = specimenSource + "\n\n// File: datecharts.go — compiled sample state and examples\n" + dateChartSource

// sourceFor shows the actual compiled default specimens and their fixture state.
// Embedding the implementation avoids a second, hand-maintained set of examples.
func sourceFor(id ComponentID) string {
	switch id {
	case "message":
		return messagePreviewSource
	case "accordion", "breadcrumb", "collapsible", "combobox", "command", "navigation-menu", "pagination", "select", "sidebar", "tabs":
		return navigationPreviewSource
	case "alert", "progress", "direction", "resizable", "scroll-area", "carousel", "table", "data-table":
		return dataLayoutPreviewSource
	case "calendar", "date-picker", "chart":
		return dateChartPreviewSource
	}
	return specimenSource
}
