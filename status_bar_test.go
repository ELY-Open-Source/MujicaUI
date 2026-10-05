package mujica

import (
	"fmt"
	"testing"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
)

func statusBarSample() StatusBarOptions {
	return StatusBarOptions{
		Left: []StatusItem{
			{ID: "branch", Text: "main", Icon: icons.Must("folder")},
			{ID: "sync", Text: "Syncing", ShowProgress: true, Progress: 0.4},
			{ID: "lint", Text: "2 warnings", Tone: StatusItemWarning, Secondary: true},
		},
		Right: []StatusItem{
			{ID: "enc", Text: "UTF-8", Secondary: true},
			{ID: "pos", Text: "Ln 12, Col 4", Secondary: true},
			{ID: "bell", Label: "Notifications", Icon: icons.Must("bell"), Action: true},
		},
	}
}

func TestStatusBarActionClick(t *testing.T) {
	var clicked []string
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		if r := StatusBar(c, statusBarSample()); r.Clicked != "" {
			clicked = append(clicked, r.Clicked)
		}
	}, 800, 100)
	if err := tt.Click("Notifications"); err != nil {
		t.Fatal(err)
	}
	_ = tt.Click("main")
	if len(clicked) != 1 || clicked[0] != "bell" {
		t.Fatalf("clicked %v, want [bell]", clicked)
	}
}

func TestStatusBarActionKeyboard(t *testing.T) {
	n := 0
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		if StatusBar(c, statusBarSample()).Clicked == "bell" {
			n++
		}
	}, 800, 100)
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Notifications") {
		t.Fatal("Tab did not reach the action")
	}
	tt.Key(0, ui.KeyEnter)
	tt.Key(0, ui.KeySpace)
	if n != 2 {
		t.Fatalf("keyboard activations %d, want 2", n)
	}
}

func TestStatusBarCollapsesSecondary(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		StatusBar(c, statusBarSample())
	}, 800, 100)
	for _, s := range []string{"main", "Syncing", "2 warnings", "UTF-8", "Ln 12, Col 4"} {
		if !tt.HasText(s) {
			t.Fatalf("wide bar lost %q", s)
		}
	}
	tt.SetSize(330, 100)
	tt.Frame()
	for _, s := range []string{"Ln 12, Col 4", "UTF-8", "2 warnings"} {
		if tt.HasText(s) {
			t.Fatalf("narrow bar still shows secondary %q; texts %q", s, tt.Texts())
		}
	}
	if !tt.HasText("main") || !tt.HasText("Syncing") {
		t.Fatalf("narrow bar hid primary items: %q", tt.Texts())
	}
	if !tt.HasText("3 more items") {
		t.Fatalf("no hidden count named: %q", tt.Texts())
	}
	// Medium width drops only the last secondary items.
	tt.SetSize(400, 100)
	tt.Frame()
	if !tt.HasText("2 warnings") || tt.HasText("Ln 12, Col 4") {
		t.Fatalf("medium bar hid the wrong items: %q", tt.Texts())
	}
	// At every width the visible items stay inside the bar.
	for w := 320; w <= 800; w += 10 {
		tt.SetSize(w, 100)
		tt.Frame()
		bell, ok := tt.Find("Notifications")
		if !ok || bell.X+bell.W > float32(w)-8+0.5 || bell.W < 28 {
			t.Fatalf("width %d: action at %+v", w, bell)
		}
	}
	tt.SetSize(800, 100)
	tt.Frame()
	if !tt.HasText("UTF-8") {
		t.Fatal("secondary items did not come back when wide")
	}
}

func TestStatusBarManyActionsKeepMinimumWidth(t *testing.T) {
	var o StatusBarOptions
	for i := range 16 {
		o.Right = append(o.Right, StatusItem{ID: fmt.Sprint(i), Label: fmt.Sprintf("Action %d", i),
			Icon: icons.Must("bell"), Action: true, Secondary: i >= 8})
	}
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		StatusBar(c, o)
	}, 560, 100)
	shown := 0
	for i := range 16 {
		r, ok := tt.Find(fmt.Sprintf("Action %d", i))
		if !ok {
			continue
		}
		shown++
		if r.W < 28 || r.X+r.W > 560-8+0.5 {
			t.Fatalf("Action %d at %+v: narrower than 28 or past the bar", i, r)
		}
	}
	if shown == 16 {
		t.Fatal("no secondary action folded")
	}
	if _, ok := tt.Find(fmt.Sprintf("%d more items", 16-shown)); !ok {
		t.Fatalf("no button for the %d hidden actions", 16-shown)
	}
}

func TestStatusBarHiddenActionsByKeyboard(t *testing.T) {
	var clicked []string
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		if id := StatusBar(c, StatusBarOptions{
			Left:  []StatusItem{{ID: "main", Text: "A long primary status text that fills the bar"}},
			Right: []StatusItem{{ID: "fmt", Text: "Format", Action: true, Secondary: true}},
		}).Clicked; id != "" {
			clicked = append(clicked, id)
		}
	}, 300, 300)
	if tt.HasText("Format") {
		t.Fatal("secondary action not folded")
	}
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("1 more items") {
		t.Fatal("Tab did not reach the hidden-items button")
	}
	tt.Key(0, ui.KeyEnter)
	if !tt.HasText("Format") {
		t.Fatal("Enter did not open the hidden items")
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	if len(clicked) != 1 || clicked[0] != "fmt" {
		t.Fatalf("clicked %v, want [fmt]", clicked)
	}
	if tt.HasText("Format") {
		t.Fatal("running a hidden action did not close the list")
	}
	if !tt.Focused("1 more items") {
		t.Fatal("focus did not return to the hidden-items button")
	}
	tt.Key(0, ui.KeyEnter)
	tt.Key(0, ui.KeyEscape)
	if tt.HasText("Format") || !tt.Focused("1 more items") {
		t.Fatal("Escape did not close the list back to its button")
	}
}

func TestStatusBarIconSize(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		StatusBar(c, StatusBarOptions{Left: []StatusItem{{ID: "i", Label: "Online", Icon: icons.Must("wifi-off")}}})
	}, 400, 100)
	// Padding of one unit each side around a 16 DIP icon.
	if r, _ := tt.Find("Online"); r.W != 24 {
		t.Fatalf("icon item width %v, want 24", r.W)
	}
}

func TestStatusBarTextScale(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		StatusBar(c, StatusBarOptions{Left: []StatusItem{{ID: "r", Text: "Ready"}}})
	}, 400, 200)
	plain, _ := tt.Find("Ready")
	tt.SetPreferences(ui.Preferences{TextScale: 2})
	tt.Frame()
	big, _ := tt.Find("Ready")
	bar, _ := tt.Find("Status bar")
	if big.H < 1.8*plain.H {
		t.Fatalf("text height %v at 200%%, was %v", big.H, plain.H)
	}
	if bar.H < big.H {
		t.Fatalf("bar %v shorter than its text %v", bar.H, big.H)
	}
}

func TestStatusBarRejectsBadItems(t *testing.T) {
	for _, o := range []StatusBarOptions{
		{Left: []StatusItem{{Text: "no id"}}},
		{Left: []StatusItem{{ID: "a", Text: "a"}}, Right: []StatusItem{{ID: "a", Text: "b"}}},
		{Left: []StatusItem{{ID: "p", Text: "p", ShowProgress: true, Progress: 1.5}}},
		{Left: []StatusItem{{ID: "n", Icon: icons.Must("bell")}}},
	} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("%+v did not panic", o)
				}
			}()
			ui.NewTester(func(c *ui.Context) {
				Use(c, Settings{})
				StatusBar(c, o)
			}, 400, 100)
		}()
	}
}

func TestStatusBarSnapshot(t *testing.T) {
	snapshot(t, "status_bar", 560, 60, func(c *ui.Context) {
		ui.Column(c).Fill().Justify(ui.End).Children(func() {
			o := statusBarSample()
			o.Left = append(o.Left, StatusItem{ID: "ok", Text: "Saved", Tone: StatusItemSuccess, Icon: icons.Must("check")})
			StatusBar(c, o)
		})
	})
}
