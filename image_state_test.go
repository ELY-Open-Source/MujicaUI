package mujica

import (
	"errors"
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestImageErrorTextNeedsRoom(t *testing.T) {
	for _, tc := range []struct {
		w    float32
		text bool
	}{{200, true}, {80, false}} {
		ready := false
		tt := ui.NewTester(func(c *ui.Context) {
			Use(c, Settings{})
			if ready {
				Image(c, nil, ImageOptions{Width: tc.w, Height: 60, Err: errors.New("404")})
			}
		}, 300, 200)
		ready = true
		tt.Frame()
		n := 0
		for _, s := range tt.Texts() {
			if s == "Image failed to load" {
				n++
			}
		}
		// One is the box's name; a second is the visible text.
		if got := n == 2; got != tc.text {
			t.Errorf("width %v: error text shown %v, want %v (%v)", tc.w, got, tc.text, tt.Texts())
		}
	}
}
