package feedback

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
)

func TestStatusIndicatorWords(t *testing.T) {
	want := map[Status]string{
		StatusSuccess: "Healthy", StatusWarning: "Degraded", StatusError: "Failing",
		StatusOffline: "Offline", StatusBusy: "Working",
	}
	for s, word := range want {
		tt := ui.NewTester(func(c *ui.Context) {
			core.Use(c, core.Settings{})
			StatusIndicator(c, s, StatusIndicatorOptions{})
		}, 300, 100)
		if !tt.HasText(word) {
			t.Errorf("%v: texts %q", s, tt.Texts())
		}
	}
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{Locale: core.ZhCN})
		StatusIndicator(c, StatusOffline, StatusIndicatorOptions{})
		StatusIndicator(c, StatusSuccess, StatusIndicatorOptions{Label: "Database online"})
	}, 300, 100)
	if !tt.HasText("离线") || !tt.HasText("Database online") || tt.HasText("正常") {
		t.Fatalf("texts %q", tt.Texts())
	}
}

func TestStatusIndicatorShapesDiffer(t *testing.T) {
	// The same tone would still differ by shape: render each in one color.
	var imgs [][]byte
	for _, s := range []Status{StatusSuccess, StatusWarning, StatusError, StatusOffline} {
		tt := ui.NewTester(func(c *ui.Context) {
			core.Use(c, core.Settings{})
			StatusIndicator(c, s, StatusIndicatorOptions{Label: "Same"})
		}, 120, 40)
		imgs = append(imgs, encode(t, tt.Image()))
	}
	for i := range imgs {
		for j := i + 1; j < len(imgs); j++ {
			if string(imgs[i]) == string(imgs[j]) {
				t.Errorf("statuses %d and %d draw alike", i, j)
			}
		}
	}
}

func TestStatusIndicatorPanicsOnUnknown(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic")
		}
	}()
	ui.NewTester(func(c *ui.Context) { core.Use(c, core.Settings{}); StatusIndicator(c, 42, StatusIndicatorOptions{}) }, 300, 100)
}

func TestStatusIndicatorSnapshot(t *testing.T) {
	snapshot(t, "status_indicator", 560, 80, func(c *ui.Context) {
		ui.Row(c).Padding(24).Gap(20).Children(func() {
			for _, s := range []Status{StatusSuccess, StatusWarning, StatusError, StatusOffline, StatusBusy} {
				StatusIndicator(c, s, StatusIndicatorOptions{})
			}
		})
	})
}
