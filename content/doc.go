// Package content provides native components that copy, colour and render
// text: a copy button, a one-line snippet, a highlighted code block, themed
// Markdown and a QR code.
//
// Elements and contexts are build-scoped. Source text, code and QR values
// are caller-owned; components keep only derived data (highlighted spans,
// the encoded QR matrix) and transient copied feedback in element state.
// Nothing here starts goroutines or performs I/O beyond writing the
// clipboard on activation.
package content
