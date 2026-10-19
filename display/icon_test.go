package display

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
)

func TestIconSizeAndLabel(t *testing.T) {
	var def, big ui.Rect
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Row(c).AlignItems(ui.Start).Children(func() {
			def = Icon(c, "star", IconOptions{}).Bounds()
			big = Icon(c, "bell", IconOptions{Size: 24, Tone: IconAccent, Label: "Alerts"}).Bounds()
		})
	}, 200, 100)
	tt.Frame()
	if def.W != 16 || def.H != 16 || big.W != 24 || big.H != 24 {
		t.Fatalf("sizes %v %v", def, big)
	}
	if _, ok := tt.Find("Alerts"); !ok {
		t.Fatal("label not exposed")
	}
}

func TestIconMissingFallbackIsVisibleAndLogged(t *testing.T) {
	var buf bytes.Buffer
	old := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&buf, nil)))
	defer slog.SetDefault(old)
	var r ui.Rect
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Row(c).AlignItems(ui.Start).Children(func() {
			r = Icon(c, "no-such-icon", IconOptions{Size: 20}).Bounds()
		})
	}, 200, 100)
	tt.Frame()
	tt.Frame()
	if r.W != 20 || r.H != 20 {
		t.Fatalf("fallback size %v", r)
	}
	if _, ok := tt.Find("Missing icon: no-such-icon"); !ok {
		t.Fatal("fallback not labeled")
	}
	if n := strings.Count(buf.String(), "mujica: icon missing"); n != 1 {
		t.Fatalf("logged %d times: %s", n, buf.String())
	}
	if !strings.Contains(buf.String(), "name=no-such-icon") {
		t.Fatalf("log lacks name: %s", buf.String())
	}
}

func TestIconBadOptionsPanic(t *testing.T) {
	for name, opts := range map[string]IconOptions{"": {}, "star": {Size: -1}, "x": {Tone: IconTone(99)}} {
		func() {
			defer func() {
				if recover() == nil {
					t.Fatalf("%q %+v: no panic", name, opts)
				}
			}()
			tt := ui.NewTester(func(c *ui.Context) {
				core.Use(c, core.Settings{})
				Icon(c, name, opts)
			}, 100, 100)
			tt.Frame()
		}()
	}
}

func TestIconSnapshot(t *testing.T) {
	snapshot(t, "icon", 360, 80, func(c *ui.Context) {
		ui.Row(c).Padding(24).Gap(14).AlignItems(ui.Center).Children(func() {
			Icon(c, "star", IconOptions{})
			Icon(c, "star", IconOptions{Size: 20, Tone: IconMuted})
			Icon(c, "star", IconOptions{Size: 24, Tone: IconAccent})
			Icon(c, "circle-check", IconOptions{Size: 20, Tone: IconSuccess})
			Icon(c, "triangle-alert", IconOptions{Size: 20, Tone: IconWarning})
			Icon(c, "circle-x", IconOptions{Size: 20, Tone: IconDanger})
			Icon(c, "info", IconOptions{Size: 20, Tone: IconInfo})
			Icon(c, "no-such-icon", IconOptions{Size: 20})
		})
	})
}
