package mujica

import (
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
	tt.SetSize(460, 100)
	tt.Frame()
	if !tt.HasText("2 warnings") || tt.HasText("Ln 12, Col 4") {
		t.Fatalf("medium bar hid the wrong items: %q", tt.Texts())
	}
	if _, ok := tt.Find("Notifications"); !ok {
		t.Fatal("primary action hidden")
	}
	tt.SetSize(800, 100)
	tt.Frame()
	if !tt.HasText("UTF-8") {
		t.Fatal("secondary items did not come back when wide")
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
