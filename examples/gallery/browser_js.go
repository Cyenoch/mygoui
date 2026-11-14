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

// touchScroll turns a finger's drag on the canvas into wheel events along
// the drag's axis. MyGo's surface reads a touch as a mouse press and sets
// touch-action: none, so without it nothing scrolls on a phone. The press
// the drag began is canceled; the scroll keeps the place the finger went
// down, so one container scrolls throughout, and glides on after a flick.
//
// A drag stays MyGo's when the finger pressed a control that follows it:
// one under a drag cursor (a divider, a grip, a colour slider), one that
// says so through window.mygouiDragHeld (internal/touchhint: a slider, a
// carousel), or any drag begun after holding still.
const touchScroll = `(() => {
  const slop = 8, hold = 300, settle = 325;
  const drags = /grab|move|resize|crosshair/;
  let id = null, target, sx, sy, lx, ly, t0, lt, vx = 0, vy = 0, frames = 0;
  let pan = false, wide = false, own = false, raf = 0;
  const wheel = (dx, dy) => target.dispatchEvent(new WheelEvent('wheel',
    {clientX: sx, clientY: sy, deltaX: dx, deltaY: dy, bubbles: true, cancelable: true}));
  const count = () => { frames++; if (id !== null && !pan) requestAnimationFrame(count); };
  addEventListener('pointerdown', e => {
    if (e.pointerType !== 'touch' || !e.isPrimary || e.target.tagName !== 'CANVAS') return;
    cancelAnimationFrame(raf);
    id = e.pointerId; target = e.target; pan = false; frames = 0;
    sx = lx = e.clientX; sy = ly = e.clientY; t0 = e.timeStamp;
    window.mygouiDragHeld = false;
    requestAnimationFrame(count);
  }, true);
  addEventListener('pointermove', e => {
    if (e.pointerId !== id) return;
    if (!pan) {
      const dx = e.clientX - sx, dy = e.clientY - sy;
      // Two frames let the content say what the finger pressed.
      if (Math.hypot(dx, dy) < slop || frames < 2) return;
      if (e.timeStamp - t0 > hold || window.mygouiDragHeld || drags.test(target.style.cursor)) { id = null; return; }
      pan = true; wide = Math.abs(dx) > Math.abs(dy); lt = e.timeStamp;
      own = true;
      target.dispatchEvent(new PointerEvent('pointercancel', {pointerId: id, pointerType: 'touch',
        isPrimary: true, clientX: e.clientX, clientY: e.clientY, bubbles: true}));
      own = false;
    }
    e.stopImmediatePropagation();
    const dx = wide ? lx - e.clientX : 0, dy = wide ? 0 : ly - e.clientY;
    const dt = Math.max(1, e.timeStamp - lt);
    vx = dx / dt; vy = dy / dt; lt = e.timeStamp; lx = e.clientX; ly = e.clientY;
    wheel(dx, dy);
  }, true);
  const end = e => {
    if (own || e.pointerId !== id) return;
    id = null;
    if (!pan) return;
    e.stopImmediatePropagation();
    if (e.type !== 'pointerup' || e.timeStamp - lt > 80) return;
    let last = performance.now();
    const glide = now => {
      const dt = now - last, k = Math.exp(-dt / settle); last = now;
      vx *= k; vy *= k;
      if (Math.hypot(vx, vy) < 0.02) return;
      wheel(vx * dt, vy * dt);
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
