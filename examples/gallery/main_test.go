package main

import (
	"runtime"
	"testing"

	"github.com/Cyenoch/mygoui/internal/gallery"
	"github.com/egoist/mygo"
)

func TestWindowOptionsKeepNativeTitleBar(t *testing.T) {
	opts := windowOptions(gallery.New())
	if opts.TitleBarStyle != mygo.TitleBarHiddenInset || opts.TitleBarHeight != gallery.TitleBarHeight ||
		opts.Frameless || opts.HideWindowButtons || opts.DisableResize || opts.DisableMove || opts.DisableFullScreen {
		t.Fatalf("native title bar is not configured: %+v", opts)
	}
	if runtime.GOOS == "darwin" {
		want := mygo.Point{X: 16, Y: 19}
		if opts.TrafficLightPosition == nil || *opts.TrafficLightPosition != want {
			t.Fatalf("macOS traffic-light position: %v, want %v", opts.TrafficLightPosition, want)
		}
	} else if opts.TrafficLightPosition != nil {
		t.Fatalf("unexpected macOS traffic-light position on %s", runtime.GOOS)
	}
}
