package mujica

import (
	"bytes"
	"slices"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestSpinnerLabelAndSize(t *testing.T) {
	var mark ui.Rect
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{Locale: ZhCN})
		Spinner(c, SpinnerOptions{})
		Spinner(c, SpinnerOptions{Size: 32, Label: "Fetching mail", ShowLabel: true})
	}, 400, 300)
	r, ok := tt.Find("加载中")
	if !ok || r.W != 16 || r.H != 16 {
		t.Fatalf("default mark %v %v", r, ok)
	}
	if n := len(slices.DeleteFunc(tt.Texts(), func(s string) bool { return s != "加载中" })); n != 1 {
		t.Fatalf("hidden label painted: %d", n)
	}
	mark, ok = tt.Find("Fetching mail")
	if !ok || mark.W != 32 {
		t.Fatalf("sized mark %v %v", mark, ok)
	}
	n := 0
	for _, s := range tt.Texts() {
		if s == "Fetching mail" {
			n++
		}
	}
	if n != 2 {
		t.Fatalf("label as name and text: %d of %q", n, tt.Texts())
	}
}

func TestSpinnerMotion(t *testing.T) {
	show := true
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		if show {
			Spinner(c, SpinnerOptions{Size: 40})
		}
	}, 100, 100)
	a := encode(t, tt.Image())
	time.Sleep(150 * time.Millisecond)
	tt.Frame()
	if bytes.Equal(a, encode(t, tt.Image())) {
		t.Fatal("spinner did not turn")
	}
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	a = encode(t, tt.Image())
	time.Sleep(150 * time.Millisecond)
	tt.Frame()
	if !bytes.Equal(a, encode(t, tt.Image())) {
		t.Fatal("spinner turned under reduced motion")
	}
	show = false
	tt.Frame()
	if _, ok := tt.Find("Loading"); ok {
		t.Fatal("hidden spinner still built")
	}
}

func TestSpinnerSnapshot(t *testing.T) {
	snapshot(t, "spinner", 360, 100, func(c *ui.Context) {
		ui.Row(c).Padding(24).Gap(24).AlignItems(ui.Center).Children(func() {
			Spinner(c, SpinnerOptions{})
			Spinner(c, SpinnerOptions{Size: 24})
			Spinner(c, SpinnerOptions{Size: 20, Label: "Syncing", ShowLabel: true})
			Button(c, "Saving", ButtonOptions{Loading: true})
		})
	})
}
