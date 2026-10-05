package mujica

import (
	"image/color"
	"math"
	"testing"

	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
	"mujicaui/theme"
)

var sparkBarValues = []float64{3, -2, 4, 1, -3}

func TestSparkBarWidthAndColors(t *testing.T) {
	values := sparkBarValues
	h := &chartHarness{build: func(c *ui.Context) ChartResult {
		return SparkBar(c, values, SparkBarOptions{Frame: ChartFrame{Label: "Alms"}, BarWidth: 10})
	}}
	tt := ui.NewTester(h.view, 300, 200)
	r, _ := tt.Find("Alms")
	if r.W != 5*11+6 || chartReadout(tt) != "-3" {
		t.Fatalf("box %v readout %q", r, chartReadout(tt))
	}
	k := theme.Light()
	at := func(i int, below bool) color.RGBA {
		x := r.X + 3 + (r.W-6)*float32(i)/4
		y := r.Y + r.H*0.3
		if below {
			y = r.Y + r.H*0.65
		}
		return color.RGBAModel.Convert(tt.Image().At(int(x), int(y))).(color.RGBA)
	}
	if at(4, true) != toRGBA(k.Danger) || at(0, false) != toRGBA(k.TextMuted) {
		t.Fatalf("last fall %v, first rise %v", at(4, true), at(0, false))
	}
	tabTo(t, tt, "Alms")
	tt.Key(0, ui.KeyHome)
	if chartReadout(tt) != "3" || at(0, false) != toRGBA(internal.Palette(k)[0]) {
		t.Fatalf("read bar not highlighted: %q", chartReadout(tt))
	}
	defer func() {
		if recover() == nil {
			t.Fatal("no panic for a negative width")
		}
	}()
	ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		SparkBar(c, values, SparkBarOptions{Frame: ChartFrame{Label: "x"}, BarWidth: -1})
	}, 100, 100)
}

func TestSparkBarSnapshot(t *testing.T) {
	snapshot(t, "spark_bar", 360, 80, func(c *ui.Context) {
		ui.Row(c).Padding(16).Gap(24).Children(func() {
			SparkBar(c, sparkBarValues, SparkBarOptions{Frame: ChartFrame{Label: "Alms"}})
			SparkBar(c, []float64{1, 3, 2, 5, 4, 6, 8}, SparkBarOptions{Frame: ChartFrame{Label: "Bells"}, BarWidth: 8})
		})
	})
}

func TestSparkBarShrinkingAndDisabled(t *testing.T) {
	values := []float64{1, 2, 3, 4}
	disabled := false
	var r ChartResult
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		ui.Column(c).Disabled(disabled).Children(func() {
			r = SparkBar(c, values, SparkBarOptions{Frame: ChartFrame{Label: "Bars"}})
		})
	}, 320, 100)
	tabTo(t, tt, "Bars")
	tt.Key(0, ui.KeyEnd)
	values = []float64{7}
	tt.Frame()
	if r.Point != 0 || !tt.HasText("7") {
		t.Fatalf("shortened data: point %d, text %q", r.Point, tt.Texts())
	}
	disabled = true
	tt.Frame()
	tt.Key(0, ui.KeyLeft)
	if r.Point != -1 {
		t.Fatalf("disabled point %d", r.Point)
	}
}

func TestSparkBarLargeData(t *testing.T) {
	values := make([]float64, 100_000)
	for i := range values {
		values[i] = float64(i%17 - 8)
	}
	var built int
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		res := SparkBar(c, values, SparkBarOptions{Frame: ChartFrame{Label: "Large bars"}})
		memo := ui.Local(res.Element, "mujica.sparkline.data", func() sparkMemo { return sparkMemo{} })
		built = len(memo.indices)
	}, 320, 100)
	tt.Frame()
	r, _ := tt.Find("Large bars")
	if built == 0 || built > 4*int(r.W)+4 {
		t.Fatalf("built %d bars for width %v", built, r.W)
	}
	if r.W > 320 {
		t.Fatalf("unbounded plot width %v", r.W)
	}
	tabTo(t, tt, "Large bars")
	tt.Key(0, ui.KeyEnd)
	if !tt.HasText("-3") {
		t.Fatalf("last value %q", tt.Texts())
	}
}

func TestSparkBarFiniteWidthAndStates(t *testing.T) {
	charts2MustPanic(t, func(c *ui.Context) {
		SparkBar(c, []float64{1}, SparkBarOptions{Frame: ChartFrame{Label: "Bars"}, BarWidth: float32(math.Inf(1))})
	})
	status := DataReady
	retries := 0
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		r := SparkBar(c, []float64{2, -1, 3}, SparkBarOptions{Frame: ChartFrame{Label: "Bars", Status: status, Height: 40}})
		if r.Retried() {
			retries++
		}
	}, 320, 100)
	ready, _ := tt.Find("Bars")
	status = DataFailed
	tt.Frame()
	if err := tt.Click("Retry"); err != nil || retries != 1 {
		t.Fatalf("retry %v count %d", err, retries)
	}
	status = DataLoading
	tt.Frame()
	if ready.W != 28 || ready.H != 40 {
		t.Fatalf("ready plot %v", ready)
	}
	status = DataReady
	tt.Frame()
	tt.Move(319, 99)
	tabTo(t, tt, "Bars")
	tt.Key(0, ui.KeyEnd)
	if !tt.HasText("3") {
		t.Fatalf("restored data %q", tt.Texts())
	}
}

func TestSparkBarMinimumHitWidth(t *testing.T) {
	for _, values := range [][]float64{nil, {1}} {
		var result ChartResult
		tt := ui.NewTester(func(c *ui.Context) {
			Use(c, Settings{})
			result = SparkBar(c, values, SparkBarOptions{Frame: ChartFrame{Label: "Bars"}})
		}, 200, 100)
		tt.Frame()
		if result.Element.Bounds().W < 28 {
			t.Fatalf("hit width %v", result.Element.Bounds().W)
		}
	}
}
