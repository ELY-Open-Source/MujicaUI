package layout

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

type splitPaneProbe struct {
	resized, toggled int
}

func splitPaneTester(st *SplitPaneState, opts SplitPaneOptions, p *splitPaneProbe) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := SplitPane(c, st, opts, func() { ui.Text(c, "Files") }, func() { ui.Text(c, "Editor") })
		r.Element.Fill()
		if r.Resized {
			p.resized++
		}
		if r.Toggled {
			p.toggled++
		}
	}, 401, 200)
}

func TestSplitPaneStartsAtHalf(t *testing.T) {
	st := &SplitPaneState{}
	tt := splitPaneTester(st, SplitPaneOptions{}, &splitPaneProbe{})
	if st.Size != 200 {
		t.Fatalf("size %v, want 200", st.Size)
	}
	if r, _ := tt.Find("Editor"); r.X != 201 {
		t.Fatalf("second pane at %v, want 201", r.X)
	}
}

func TestSplitPaneKeyboard(t *testing.T) {
	st := &SplitPaneState{Size: 150}
	p := &splitPaneProbe{}
	tt := splitPaneTester(st, SplitPaneOptions{MinFirst: 100, MinSecond: 120}, p)
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Divider") {
		t.Fatal("Tab did not focus the divider")
	}
	tt.Key(0, ui.KeyRight)
	if st.Size != 166 {
		t.Fatalf("Right moved to %v, want 166", st.Size)
	}
	tt.Key(0, ui.KeyHome)
	if st.Size != 100 {
		t.Fatalf("Home moved to %v, want the minimum 100", st.Size)
	}
	tt.Key(0, ui.KeyLeft)
	if st.Size != 100 || st.Collapsed {
		t.Fatalf("Left past the minimum gave %v collapsed=%v", st.Size, st.Collapsed)
	}
	tt.Key(0, ui.KeyEnd)
	if st.Size != 401-1-120 {
		t.Fatalf("End moved to %v, want 280", st.Size)
	}
	tt.Key(0, ui.KeyEnter)
	if st.Collapsed || p.toggled != 0 {
		t.Fatal("Enter collapsed a pane that is not collapsible")
	}
	if p.resized != 3 {
		t.Fatalf("resized %d times, want 3", p.resized)
	}
}

func TestSplitPaneDragClampsToMinimum(t *testing.T) {
	st := &SplitPaneState{Size: 200}
	tt := splitPaneTester(st, SplitPaneOptions{MinFirst: 80}, &splitPaneProbe{})
	tt.Press(200, 100)
	tt.Move(240, 100)
	if st.Size != 240 {
		t.Fatalf("dragged to %v, want 240", st.Size)
	}
	tt.Move(10, 100)
	tt.Release(10, 100)
	if st.Size != 80 || st.Collapsed {
		t.Fatalf("dragged past the minimum to %v collapsed=%v, want 80", st.Size, st.Collapsed)
	}
	if !tt.Focused("Divider") {
		t.Fatal("drag did not focus the divider")
	}
}

func TestSplitPaneCollapseFirst(t *testing.T) {
	st := &SplitPaneState{Size: 200}
	p := &splitPaneProbe{}
	tt := splitPaneTester(st, SplitPaneOptions{MinFirst: 100, Collapsible: SplitPaneFirst}, p)
	// Dragging under half the minimum collapses; back over it opens again.
	tt.Press(200, 100)
	tt.Move(120, 100)
	if st.Collapsed || st.Size != 120 {
		t.Fatalf("at 120: size %v collapsed=%v", st.Size, st.Collapsed)
	}
	tt.Move(40, 100)
	if !st.Collapsed || tt.HasText("Files") {
		t.Fatalf("under half the minimum: collapsed=%v files shown=%v", st.Collapsed, tt.HasText("Files"))
	}
	if r, _ := tt.Find("Editor"); r.X != 1 {
		t.Fatalf("second pane at %v, want the whole width from 1", r.X)
	}
	tt.Move(160, 100)
	tt.Release(160, 100)
	if st.Collapsed || st.Size != 160 {
		t.Fatalf("dragged back: size %v collapsed=%v", st.Size, st.Collapsed)
	}
	// Keyboard: Left at the minimum collapses, Right opens, Enter toggles.
	tt.Key(0, ui.KeyHome)
	tt.Key(0, ui.KeyLeft)
	if !st.Collapsed {
		t.Fatal("Left at the minimum did not collapse")
	}
	tt.Key(0, ui.KeyRight)
	if st.Collapsed || st.Size != 100 {
		t.Fatalf("Right did not reopen at 100: size %v collapsed=%v", st.Size, st.Collapsed)
	}
	tt.Key(0, ui.KeyEnter)
	if !st.Collapsed {
		t.Fatal("Enter did not collapse")
	}
	tt.Key(0, ui.KeyEnter)
	if st.Collapsed || !tt.HasText("Files") {
		t.Fatal("Enter did not reopen")
	}
	if p.toggled != 6 {
		t.Fatalf("toggled %d times, want 6", p.toggled)
	}
}

func TestSplitPaneCollapseSecondVertical(t *testing.T) {
	st := &SplitPaneState{Size: 100}
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		SplitPane(c, st, SplitPaneOptions{Vertical: true, MinSecond: 60, Collapsible: SplitPaneSecond},
			func() { ui.Text(c, "Editor") }, func() { ui.Text(c, "Console") }).Element.Fill()
	}, 300, 201)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyDown)
	if st.Size != 116 {
		t.Fatalf("Down moved to %v, want 116", st.Size)
	}
	tt.Key(0, ui.KeyEnd)
	tt.Key(0, ui.KeyDown)
	if !st.Collapsed || tt.HasText("Console") {
		t.Fatal("Down past the second minimum did not collapse it")
	}
	tt.Key(0, ui.KeyUp)
	if st.Collapsed || st.Size != 140 {
		t.Fatalf("Up did not reopen at 140: size %v collapsed=%v", st.Size, st.Collapsed)
	}
}

func TestSplitPaneDoubleClickToggles(t *testing.T) {
	st := &SplitPaneState{Size: 200}
	tt := splitPaneTester(st, SplitPaneOptions{Collapsible: SplitPaneFirst}, &splitPaneProbe{})
	tt.ClickAt(200, 100)
	tt.ClickAt(200, 100)
	if !st.Collapsed || st.Size != 200 {
		t.Fatalf("double click: collapsed=%v size %v, want collapsed keeping 200", st.Collapsed, st.Size)
	}
}

func TestSplitPaneRejectsCollapsedWithoutSide(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("collapsed state without a Collapsible pane did not panic")
		}
	}()
	splitPaneTester(&SplitPaneState{Collapsed: true}, SplitPaneOptions{}, &splitPaneProbe{})
}

func TestSplitPaneSnapshot(t *testing.T) {
	snapshot(t, "split_pane", 420, 160, func(c *ui.Context) {
		k := core.Tokens(c)
		st := ui.Local(c.Root(), "snap", func() SplitPaneState { return SplitPaneState{Size: 140} })
		pane := func(s string, bg ui.Color) func() {
			return func() {
				ui.Box(c).Fill().Padding(12).Background(bg).Children(func() { ui.Text(c, s) })
			}
		}
		ui.Column(c).Padding(16).Children(func() {
			SplitPane(c, st, SplitPaneOptions{}, pane("Files", k.Surface), pane("Editor", k.Background)).
				Element.Height(128).Border(1, k.Border).Radius(4).Clip()
		})
	})
}

func TestSplitPaneSizesInsideBorder(t *testing.T) {
	st := &SplitPaneState{Size: 150}
	var second ui.Rect
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		SplitPane(c, st, SplitPaneOptions{MinSecond: 120}, func() { ui.Text(c, "Files") }, func() {
			second = ui.Box(c).Fill().Absolute().Bounds()
		}).Element.Fill().Border(1, ui.RGB(0, 0, 0))
	}, 401, 200)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnd)
	// 399 DIPs inside the border: 278 + 1 divider + 120.
	if st.Size != 278 || second.W != 120 {
		t.Fatalf("End gave size %v and second pane %v, want 278 and 120", st.Size, second.W)
	}
	tt.Press(1+278, 100)
	tt.Move(1+200, 100)
	tt.Release(1+200, 100)
	if st.Size != 200 {
		t.Fatalf("drag to x=201 gave %v, want 200", st.Size)
	}
}

func TestSplitPaneClampsOnResizeAlone(t *testing.T) {
	st := &SplitPaneState{Size: 299}
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		SplitPane(c, st, SplitPaneOptions{}, func() { ui.Text(c, "Files") }, func() { ui.Text(c, "Editor") }).Element.Fill()
	}, 400, 200)
	tt.SetSize(250, 200)
	// 250 less the divider and the 40 DIP second minimum.
	if r, _ := tt.Find("Editor"); st.Size != 209 || r.X != 210 {
		t.Fatalf("after resize: size %v, second pane at %v; want 209 and 210", st.Size, r.X)
	}
}
