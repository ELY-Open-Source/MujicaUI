package devtools

import (
	"fmt"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
)

var waterfallSample = []WaterfallSpan{
	{ID: "db", Parent: "api", Name: "SELECT guests", Service: "postgres", Start: 40 * time.Millisecond, Duration: 90 * time.Millisecond},
	{ID: "root", Name: "GET /ball", Service: "gateway", Duration: 320 * time.Millisecond},
	{ID: "api", Parent: "root", Name: "ListGuests", Service: "api", Start: 20 * time.Millisecond, Duration: 200 * time.Millisecond},
	{ID: "pay", Parent: "root", Name: "Charge", Service: "payments", Start: 230 * time.Millisecond, Duration: 80 * time.Millisecond, Error: true},
}

func traceTester(spans []WaterfallSpan, changes *int) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if TraceWaterfall(c, spans, TraceWaterfallOptions{}).Changed() {
			*changes++
		}
	}, 700, 460)
}

func TestTraceWaterfallOrderAndDetail(t *testing.T) {
	n := 0
	tt := traceTester(waterfallSample, &n)
	if !tt.HasText("4 spans · 320 ms") {
		t.Fatalf("summary: %q", tt.Texts())
	}
	root, _ := tt.Find("GET /ball")
	api, _ := tt.Find("ListGuests")
	db, _ := tt.Find("SELECT guests")
	pay, _ := tt.Find("Charge")
	if !(root.Y < api.Y && api.Y < db.Y && db.Y < pay.Y && db.X > api.X && api.X > root.X) {
		t.Fatalf("tree order: root %v api %v db %v pay %v", root, api, db, pay)
	}
	_ = tt.Click("Charge")
	if n != 1 || !tt.HasText("payments · starts at 230 ms · takes 80 ms") {
		t.Fatalf("detail: %q", tt.Texts())
	}
	tt.Key(0, ui.KeyUp)
	tt.Key(0, ui.KeyEnter)
	if !tt.HasText("postgres · starts at 40 ms · takes 90 ms") {
		t.Fatalf("keyboard detail: %q", tt.Texts())
	}
}

func TestTraceWaterfallRulesAndMany(t *testing.T) {
	n := 0
	tt := traceTester(nil, &n)
	if !tt.HasText("No spans") {
		t.Fatal("empty")
	}
	many := make([]WaterfallSpan, 100000)
	many[0] = WaterfallSpan{ID: "0", Name: "root", Duration: time.Second}
	for i := 1; i < len(many); i++ {
		many[i] = WaterfallSpan{ID: fmt.Sprint(i), Parent: "0", Name: fmt.Sprint("span ", i), Start: time.Duration(i) * time.Microsecond, Duration: time.Microsecond}
	}
	tt = traceTester(many, &n)
	if got := len(tt.Texts()); got > 400 {
		t.Fatalf("%d texts for 100000 spans", got)
	}
	defer func() {
		if recover() == nil {
			t.Fatal("cycle did not panic")
		}
	}()
	traceTester([]WaterfallSpan{{ID: "a", Parent: "b", Name: "a"}, {ID: "b", Parent: "a", Name: "b"}}, &n)
}

func TestTraceWaterfallSnapshot(t *testing.T) {
	snapshot(t, "trace_waterfall", 680, 280, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() { TraceWaterfall(c, waterfallSample, TraceWaterfallOptions{Height: 180}) })
	})
}

func TestTraceWaterfallAncestorDisabled(t *testing.T) {
	disabled := false
	changes := 0
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Disabled(disabled).Children(func() {
			if TraceWaterfall(c, waterfallSample, TraceWaterfallOptions{}).Changed() {
				changes++
			}
		})
	}, 700, 460)
	_ = tt.Click("Charge")
	disabled = true
	tt.Key(0, ui.KeyUp)
	tt.Key(0, ui.KeyEnter)
	_ = tt.Click("ListGuests")
	if changes != 1 || !tt.HasText("payments · starts at 230 ms · takes 80 ms") {
		t.Fatal("disabled trace selection changed")
	}
}
