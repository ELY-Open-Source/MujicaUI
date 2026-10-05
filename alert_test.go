package mujica

import (
	"slices"
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestAlertActionAndClose(t *testing.T) {
	open := true
	actions, closes := 0, 0
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		r := Alert(c, &open, "Disk almost full", AlertOptions{Severity: SeverityWarning, Title: "Storage", Action: "Free space"})
		if r.Action {
			actions++
		}
		if r.Closed {
			closes++
		}
	}, 500, 300)
	if err := tt.Click("Free space"); err != nil {
		t.Fatal(err)
	}
	if actions != 1 || !open {
		t.Fatalf("actions %d open %v", actions, open)
	}
	if err := tt.Click("Dismiss"); err != nil {
		t.Fatal(err)
	}
	if open || closes != 1 {
		t.Fatalf("after dismiss open %v closes %d", open, closes)
	}
	if tt.HasText("Disk almost full") {
		t.Fatal("closed alert still shows")
	}
}

func TestAlertKeyboard(t *testing.T) {
	open := true
	actions := 0
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		if Alert(c, &open, "Sync failed", AlertOptions{Severity: SeverityError, Action: "Retry"}).Action {
			actions++
		}
	}, 500, 300)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	if actions != 1 {
		t.Fatalf("Enter on action: %d", actions)
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeySpace)
	if open {
		t.Fatal("Space on dismiss left it open")
	}
}

func TestAlertWithoutOpenHasNoDismiss(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		Alert(c, nil, "Read only", AlertOptions{})
	}, 500, 300)
	if _, ok := tt.Find("Dismiss"); ok {
		t.Fatal("dismiss button without open")
	}
	if _, ok := tt.Find("Information"); !ok {
		t.Fatal("alert lacks its severity name")
	}
}

func TestAlertAnnouncesWarningsOnce(t *testing.T) {
	sev := SeverityInfo
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{Locale: ZhCN})
		Alert(c, nil, "磁盘已满", AlertOptions{Severity: sev})
	}, 500, 300)
	if a := tt.Announcements(); len(a) != 0 {
		t.Fatalf("info announced %q", a)
	}
	sev = SeverityError
	tt.Frame()
	tt.Frame()
	if a := tt.Announcements(); !slices.Equal(a, []string{"错误: 磁盘已满"}) {
		t.Fatalf("error announcements %q", a)
	}
}

func TestAlertPanicsWithoutMessage(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic")
		}
	}()
	ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		Alert(c, nil, "", AlertOptions{})
	}, 300, 200)
}

func TestAlertSnapshot(t *testing.T) {
	open := true
	snapshot(t, "alert", 560, 400, func(c *ui.Context) {
		ui.Column(c).Padding(20).Gap(12).Children(func() {
			Alert(c, nil, "A new version is ready to install.", AlertOptions{})
			Alert(c, &open, "Your changes were saved.", AlertOptions{Severity: SeveritySuccess})
			Alert(c, nil, "The certificate expires in 3 days.", AlertOptions{Severity: SeverityWarning, Title: "Certificate"})
			Alert(c, &open, "The archive could not be read.", AlertOptions{Severity: SeverityError, Title: "Import failed", Action: "Retry"})
		})
	})
}
