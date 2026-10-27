// Package mygoui provides native MyGo components following coss ui.
//
// Apply theme.Apply at the start of a view. Callers own model values and widget
// state per window; constructors consume a build-scoped *ui.Context and return
// build-scoped ui.Element values. Retain ui.Handle for persistent focus and
// geometry, never elements or custom parts from an earlier build.
//
// Composition uses typed options, content callbacks and keyed collections.
// Related controls live in form, navigation, overlay, date, data and other
// family packages. A card composes ordinary native children:
//
//	mygoui.Card(c, mygoui.CardOptions{}, func() {
//		mygoui.CardHeader(c, func() {
//			mygoui.CardTitle(c, "Project")
//		})
//		mygoui.CardContent(c, func() {
//			mygoui.Paragraph(c, "Ready to deploy", false)
//		})
//	})
//
// Colors are read live through theme.Colors. Host preferences configure
// semantic roles, radii and body/heading/mono families; applications register
// their own fonts explicitly. Importing this library does not initialize an
// application or start I/O or background jobs.
//
// ScrollArea wraps native scrolling. ScrollFade provides subtree alpha masking,
// shadow and opaque-surface edge treatments through native paint and effects.
// docs/components.md maps the coss catalog and native limits, docs/theme.md
// describes customization, and docs/references.md records source provenance.
package mygoui
