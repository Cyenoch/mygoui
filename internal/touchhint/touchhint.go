// Package touchhint tells a page which presses belong to a control.
//
// MyGo's web surface reads a finger as a mouse press and scrolls nothing
// from it, so a page turns a finger's drag into scrolling itself. It cannot
// see what the finger pressed: a control that follows a held pointer and
// shows no drag cursor (a slider, a carousel) says so here, and the page
// leaves that drag to it. Outside a browser Held does nothing.
package touchhint
