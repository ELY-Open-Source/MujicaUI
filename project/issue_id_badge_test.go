package project

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

func TestIssueIdBadgeCopy(t *testing.T) {
	n, disabled := 0, false
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Children(func() {
			if IssueIdBadge(c, "ENG-123", IssueIdBadgeOptions{Copyable: true, Disabled: disabled}).Changed() {
				n++
			}
			IssueIdBadge(c, "OPS-7", IssueIdBadgeOptions{})
		})
	}, 300, 120)
	_ = tt.Click("Copy ENG-123")
	if tt.Clipboard() != "ENG-123" || n != 1 || !tt.HasText("Copied") {
		t.Fatalf("copy %q %d %v", tt.Clipboard(), n, tt.Texts())
	}
	if a := tt.Announcements(); len(a) != 1 || a[0] != "ENG-123 copied" {
		t.Fatalf("announce %q", a)
	}
	tt.SetClipboard("")
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	if tt.Clipboard() != "ENG-123" || n != 2 {
		t.Fatal("Enter should copy")
	}
	if !tt.HasText("OPS-7") || tt.HasText("Copy OPS-7") {
		t.Fatal("plain badge should not be a button")
	}
	disabled = true
	tt.SetClipboard("")
	tt.Frame()
	_ = tt.Click("Copy ENG-123")
	if tt.Clipboard() != "" || n != 2 {
		t.Fatal("disabled badge copied")
	}
}

func TestIssueIdBadgeCopiedFades(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		IssueIdBadge(c, "ENG-1", IssueIdBadgeOptions{Copyable: true})
	}, 300, 120)
	_ = tt.Click("Copy ENG-1")
	if !tt.HasText("Copied") {
		t.Fatal("no Copied note")
	}
	time.Sleep(issueCopiedFor + 50*time.Millisecond)
	tt.Frame()
	if tt.HasText("Copied") {
		t.Fatal("Copied should go away")
	}
}

func TestIssueIdBadgeInvalid(t *testing.T) {
	for _, id := range []string{"eng-1", "ENG-0", "ENG123", ""} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("%q: expected panic", id)
				}
			}()
			ui.NewTester(func(c *ui.Context) { core.Use(c, core.Settings{}); IssueIdBadge(c, id, IssueIdBadgeOptions{}) }, 10, 10)
		}()
	}
}

func TestIssueIdBadgeSnapshot(t *testing.T) {
	snapshot(t, "issue_id_badge", 300, 60, func(c *ui.Context) {
		ui.Row(c).Padding(16).Gap(12).AlignItems(ui.Center).Children(func() {
			IssueIdBadge(c, "ENG-123", IssueIdBadgeOptions{Copyable: true})
			IssueIdBadge(c, "OPS-7", IssueIdBadgeOptions{})
		})
	})
}
