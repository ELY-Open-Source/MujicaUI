package feedback

import (
	"bytes"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
)

// motionFrames builds view and returns its tester and a counter of built frames.
func motionFrames(view func(c *ui.Context), w, h int) (*ui.Tester, *int) {
	n := new(int)
	tt := ui.NewTester(func(c *ui.Context) {
		*n++
		core.Use(c, core.Settings{})
		view(c)
	}, w, h)
	return tt, n
}

// motionBusy reports whether one more frame keeps asking for frames.
func motionBusy(tt *ui.Tester, n *int) bool {
	before := *n
	tt.Frame()
	return *n-before > 2
}

// motionCount counts texts and names equal to s.
func motionCount(tt *ui.Tester, s string) int {
	n := 0
	for _, x := range tt.Texts() {
		if x == s {
			n++
		}
	}
	return n
}

var loaderCases = map[string]func(c *ui.Context, o LoaderOptions) ui.Element{
	"dots": DotsLoader, "bars": BarsLoader, "pulse": PulseLoader, "orbit": OrbitLoader, "wave": WaveLoader,
}

func testLoaderBehavior(t *testing.T, loader func(*ui.Context, LoaderOptions) ui.Element) {
	show := true
	tt, n := motionFrames(func(c *ui.Context) {
		if show {
			loader(c, LoaderOptions{Size: 32, Label: "Thinking"})
		}
	}, 200, 80)
	if _, ok := tt.Find("Thinking"); !ok {
		t.Fatal("loader not named")
	}
	if motionCount(tt, "Thinking") != 1 {
		t.Fatal("label shown without ShowLabel")
	}
	a := encode(t, tt.Image())
	tt.Advance(120 * time.Millisecond)
	if bytes.Equal(a, encode(t, tt.Image())) || !motionBusy(tt, n) {
		t.Fatal("loader does not move")
	}
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	if motionBusy(tt, n) || motionCount(tt, "Thinking") != 2 {
		t.Fatalf("reduced motion: still repainting or no static label: %q", tt.Texts())
	}
	a = encode(t, tt.Image())
	tt.Advance(120 * time.Millisecond)
	if !bytes.Equal(a, encode(t, tt.Image())) {
		t.Fatal("loader moved under reduced motion")
	}
	tt.SetPreferences(ui.Preferences{TextScale: 1})
	show = false
	tt.Frame()
	if motionBusy(tt, n) {
		t.Fatal("hidden loader keeps repainting")
	}
}

func TestLoadersInvalidSize(t *testing.T) {
	for name, loader := range loaderCases {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("%s: negative size did not panic", name)
				}
			}()
			ui.Render(func(c *ui.Context) {
				core.Use(c, core.Settings{})
				loader(c, LoaderOptions{Size: -1})
			}, 100, 50, 1)
		}()
	}
}

func loaderSnapshot(t *testing.T, name string, loader func(c *ui.Context, o LoaderOptions) ui.Element) {
	snapshot(t, name, 480, 90, func(c *ui.Context) {
		ui.Row(c).Padding(20).Gap(24).AlignItems(ui.Center).Children(func() {
			loader(c, LoaderOptions{})
			loader(c, LoaderOptions{Size: 28, Label: "Thinking", ShowLabel: true})
			loader(c, LoaderOptions{Size: 20, Color: core.Tokens(c).Info})
		})
	})
}
