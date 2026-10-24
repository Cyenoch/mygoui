// Package mygoui is a native HeroUI v3 component library for MyGo.
//
// Constructors receive a build-scoped *ui.Context and return a build-scoped
// ui.Element. Callers own model state; the toolkit owns input, focus, editing,
// selection, animation and accessibility. Components read the live HeroUI
// palette from package theme and never start background work or I/O.
//
// Composition is Go-native: roots take typed options and a body callback for
// parts. A card, for example, composes CardHeader, CardContent and CardFooter
// inside one body:
//
//	mygoui.Card(c, mygoui.CardOptions{}, func() {
//		mygoui.CardHeader(c, func() {
//			mygoui.CardTitle(c, "Project")
//			mygoui.CardDescription(c, "Deployed 3 minutes ago")
//		})
//		mygoui.CardContent(c, func() { /* ... */ })
//	})
//
// Reference implementations are pinned under references/; docs/references.md
// records how to review upstream changes.
package mygoui
