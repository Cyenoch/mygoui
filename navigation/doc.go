// Package navigation implements HeroUI v3's navigation surfaces natively:
// accordions and disclosures, tabs, breadcrumbs, pagination, toolbars and
// drop-down menus.
//
// Every constructor takes a *ui.Context and returns a build-scoped
// ui.Element. Expansion, selection and page state is caller-owned through
// pointers, or kept internally when the option is nil so a caller can opt out
// of state plumbing. Coordinated collections take typed slices:
// AccordionOptions.Items, DisclosureGroupOptions.Items,
// TabsOptions.Items and MenuOptions.Sections.
//
// Platform limits: a breadcrumb overflow menu uses the native system menu, so
// its placement and look are owned by the platform. Drop-down and submenu
// panels are drawn in the window (ui.PopoverBase); the toolkit owns edge
// flipping and does not virtualize long menus.
package navigation
