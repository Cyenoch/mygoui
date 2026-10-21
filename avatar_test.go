package mygoui_test

import (
	"image"
	"testing"

	"github.com/Cyenoch/mygoui"
	myui "github.com/egoist/mygo/ui"
)

func TestAvatarTransitionsBetweenLoadedImageAndFallback(t *testing.T) {
	var photo *myui.Bitmap
	tt := myui.NewTester(func(c *myui.Context) {
		mygoui.Avatar(c, mygoui.AvatarOptions{Name: "Ada Lovelace", Image: photo})
	}, 160, 100)
	if !tt.HasText("AL") {
		t.Fatal("typed-nil image did not display initials")
	}
	photo = myui.NewBitmap(image.NewRGBA(image.Rect(0, 0, 8, 8)))
	tt.Frame()
	if tt.HasText("AL") {
		t.Fatal("loaded image left fallback visible")
	}
	photo = nil
	tt.Frame()
	if !tt.HasText("AL") {
		t.Fatal("removed image did not restore fallback")
	}
}

func TestAvatarInitialsRespectUnicodeAndExplicitFallback(t *testing.T) {
	cases := []struct{ name, fallback, want string }{
		{name: "Élodie Durand", want: "ÉD"},
		{name: "山田 太郎", want: "山太"},
		{name: "  Ada   Lovelace  ", want: "AL"},
		{name: "Ada Adams", want: "AA"},
		{name: "", want: "?"},
		{name: "Ada Lovelace", fallback: "AD", want: "AD"},
	}
	for _, tc := range cases {
		t.Run(tc.name+tc.fallback, func(t *testing.T) {
			tt := myui.NewTester(func(c *myui.Context) {
				mygoui.Avatar(c, mygoui.AvatarOptions{Name: tc.name, Fallback: tc.fallback})
			}, 160, 100)
			if !tt.HasText(tc.want) {
				t.Fatalf("missing initials %q in %v", tc.want, tt.Texts())
			}
		})
	}
}
