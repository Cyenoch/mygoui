package main

import (
	"strings"
	"syscall/js"
)

// appKeys are the letters the gallery takes with the primary modifier: text
// editing and its search.
const appKeys = "acxvzyk"

// keepBrowserShortcuts leaves the browser its own shortcuts, as reload and
// find. MyGo's surface cancels every key pressed with the primary modifier,
// so the others are kept from it.
func keepBrowserShortcuts() {
	window := js.Global()
	mac := strings.Contains(window.Get("navigator").Get("userAgent").String(), "Mac")
	window.Call("addEventListener", "keydown", js.FuncOf(func(_ js.Value, args []js.Value) any {
		e := args[0]
		primary := e.Get("ctrlKey").Bool()
		if mac {
			primary = e.Get("metaKey").Bool()
		}
		key := strings.ToLower(e.Get("key").String())
		if primary && len(key) == 1 && !strings.Contains(appKeys, key) {
			e.Call("stopPropagation")
		}
		return nil
	}), true)
}

// touchScroll turns a finger's vertical drag on the canvas into wheel
// events. MyGo's surface reads a touch as a mouse press and sets
// touch-action: none, so without it nothing scrolls on a phone. The press
// the drag began is canceled; the scroll keeps the place the finger went
// down, so one container scrolls throughout, and glides on after a flick.
// A sideways drag, or one begun after holding still, stays MyGo's: sliders,
// dividers and reordering keep working.
const touchScroll = `(() => {
  const slop = 8, hold = 300, settle = 325;
  let id = null, target, sx, sy, ly, t0, lt, vy = 0, pan = false, own = false, raf = 0;
  const wheel = dy => target.dispatchEvent(new WheelEvent('wheel',
    {clientX: sx, clientY: sy, deltaY: dy, bubbles: true, cancelable: true}));
  addEventListener('pointerdown', e => {
    if (e.pointerType !== 'touch' || !e.isPrimary || e.target.tagName !== 'CANVAS') return;
    cancelAnimationFrame(raf);
    id = e.pointerId; target = e.target; pan = false;
    sx = e.clientX; sy = ly = e.clientY; t0 = e.timeStamp;
  }, true);
  addEventListener('pointermove', e => {
    if (e.pointerId !== id) return;
    if (!pan) {
      const dx = e.clientX - sx, dy = e.clientY - sy;
      if (Math.hypot(dx, dy) < slop) return;
      if (e.timeStamp - t0 > hold || Math.abs(dx) > Math.abs(dy)) { id = null; return; }
      pan = true; lt = e.timeStamp;
      own = true;
      target.dispatchEvent(new PointerEvent('pointercancel', {pointerId: id, pointerType: 'touch',
        isPrimary: true, clientX: e.clientX, clientY: e.clientY, bubbles: true}));
      own = false;
    }
    e.stopImmediatePropagation();
    const dy = ly - e.clientY, dt = Math.max(1, e.timeStamp - lt);
    vy = dy / dt; lt = e.timeStamp; ly = e.clientY;
    wheel(dy);
  }, true);
  const end = e => {
    if (own || e.pointerId !== id) return;
    id = null;
    if (!pan) return;
    e.stopImmediatePropagation();
    if (e.type !== 'pointerup' || e.timeStamp - lt > 80) return;
    let last = performance.now();
    const glide = now => {
      const dt = now - last; last = now;
      vy *= Math.exp(-dt / settle);
      if (Math.abs(vy) < 0.02) return;
      wheel(vy * dt);
      raf = requestAnimationFrame(glide);
    };
    raf = requestAnimationFrame(glide);
  };
  addEventListener('pointerup', end, true);
  addEventListener('pointercancel', end, true);
})()`

// scrollWithTouch lets a finger scroll the gallery in a phone's browser.
func scrollWithTouch() {
	js.Global().Call("eval", touchScroll)
}
