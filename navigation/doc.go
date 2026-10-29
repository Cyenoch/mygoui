// Package navigation provides Coss-inspired native navigation controls:
// accordions, disclosures, tabs, segmented controls, breadcrumbs,
// pagination, toolbars, menus, context menus, and searchable commands.
//
// Elements and contexts are build-scoped. Keep caller-owned expansion,
// selection, query, and focus state per window; some controls also support
// component-local state when their option is nil.
//
// ContextMenu uses MyGo's native system menu; PopoverMenu draws the context
// menu inside the window in the current native theme. Both keep toolkit-owned
// input, selection semantics and submenu navigation. Breadcrumb overflow uses
// the system menu. Drop-down menus and submenus are drawn in the window and
// use MyGo's anchored popover placement.
package navigation
