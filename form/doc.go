// Package form composes caller-owned Go values with maintained MyGo inputs.
//
// Appearance and interaction sources are the official Base component pages,
// https://ui.shadcn.com/r/styles/base-rhea/{input,textarea,field,label,
// input-group,input-otp,checkbox,radio-group,switch,slider,native-select}.json,
// and shadcn's canonical registry/styles/style-rhea.css. These specify visual
// and observable behavior; no React providers, slots or editing code is ported.
//
// Field delegates label activation and AX naming/help/error association to
// MyGo Field. Input and Textarea retain its editor, undo, selection, composition
// and secure masking. OTP adds ASCII committed-text filtering and visual guides
// around one editor, preserving native caret and selection. Choice controls use
// MyGo bases; RadioGroup keeps native roving focus and selection. NativeSelect
// opens an operating-system menu; groups appear as native submenus.
//
// RangeSlider preserves MyGo's horizontal two-thumb control, native painting,
// ticks and independent focus rings. Its toolkit currently offers no vertical
// range orientation; single Slider supports both orientations. MyGo offers
// native read-only AX state for text only. Other read-only controls preserve
// their caller's value with a temporary binding and describe that policy in AX.
// All geometry and colors use the shared theme; live values never read a static
// palette. Native visual and operating-system gesture verification is separate
// from source implementation.
package form
