package mujica

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func callStackSample() []DebugThread {
	return []DebugThread{
		{ID: "t1", Name: "main", Frames: []StackFrame{
			{"f1", "court.Seal", "court/seal.go", 14},
			{"f2", "court.Decree", "court/decree.go", 88},
			{"f3", "main.main", "main.go", 9},
		}},
		{ID: "t2", Name: "herald worker", Frames: []StackFrame{{"f4", "herald.Run", "herald/run.go", 31}}},
	}
}

type callProbe struct {
	threads  []DebugThread
	s        CallStackState
	disabled bool
	changed  int
}

func (p *callProbe) view(c *ui.Context) {
	Use(c, Settings{})
	r := CallStack(c, p.threads, &p.s, CallStackOptions{Disabled: p.disabled})
	r.Element.Height(320)
	if r.Changed() {
		p.changed++
	}
}

func TestCallStackChooseFrame(t *testing.T) {
	p := &callProbe{threads: callStackSample(), s: CallStackState{Current: "f1"}}
	tt := ui.NewTester(p.view, 460, 340)
	if p.s.selected != 1 {
		t.Fatalf("current frame row %d, want 1", p.s.selected)
	}
	if err := tt.Click("court.Decree"); err != nil {
		t.Fatal(err)
	}
	if p.s.Current != "f2" || p.changed != 1 {
		t.Fatalf("current %q changed %d", p.s.Current, p.changed)
	}
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyDown)
	if p.s.Current != "f3" || p.s.selected != 4 {
		t.Fatalf("keys: current %q row %d", p.s.Current, p.s.selected)
	}
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyDown)
	if p.s.Current != "f4" {
		t.Fatalf("keys did not pass the thread header: %q", p.s.Current)
	}
}

func TestCallStackFoldThread(t *testing.T) {
	p := &callProbe{threads: callStackSample(), s: CallStackState{Current: "f4"}}
	tt := ui.NewTester(p.view, 460, 340)
	if err := tt.Click("Fold main"); err != nil {
		t.Fatal(err)
	}
	if tt.HasText("court.Seal") || !tt.HasText("herald.Run") || p.s.Current != "f4" {
		t.Fatalf("fold: %q current %q", tt.Texts(), p.s.Current)
	}
	if !tt.HasText("depth 3") {
		t.Fatal("thread header lacks its count")
	}
}

func TestCallStackCallerSetsCurrent(t *testing.T) {
	p := &callProbe{threads: callStackSample()}
	tt := ui.NewTester(p.view, 460, 340)
	p.s.Current = "f3"
	tt.Frame()
	if p.s.selected != 3 || p.changed != 0 {
		t.Fatalf("row %d changed %d", p.s.selected, p.changed)
	}
}

func TestCallStackDisabledAndEmpty(t *testing.T) {
	p := &callProbe{threads: callStackSample(), s: CallStackState{Current: "f1"}, disabled: true}
	tt := ui.NewTester(p.view, 460, 340)
	_ = tt.Click("court.Decree")
	if p.s.Current != "f1" {
		t.Fatal("disabled stack changed frame")
	}
	tt = ui.NewTester((&callProbe{}).view, 460, 340)
	if !tt.HasText("Not paused") {
		t.Fatalf("texts %q", tt.Texts())
	}
}

func TestCallStackSnapshot(t *testing.T) {
	snapshot(t, "call_stack", 420, 230, func(c *ui.Context) {
		p := ui.Local(c.Root(), "snap", func() *callProbe { return &callProbe{threads: callStackSample(), s: CallStackState{Current: "f2"}} })
		ui.Column(c).Padding(12).Children(func() {
			CallStack(c, (*p).threads, &(*p).s, CallStackOptions{}).Element.Height(206)
		})
	})
}
