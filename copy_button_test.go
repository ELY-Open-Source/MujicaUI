package mujica

import (
	"slices"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestCopyButtonCopiesAndAnnounces(t *testing.T) {
	n := 0
	var w float32
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		ui.Row(c).Children(func() {
			r := CopyButton(c, "sk-live-123", CopyButtonOptions{Label: "Copy key"})
			if r.Changed() {
				n++
			}
			w = r.Element.Bounds().W
		})
	}, 400, 200)
	tt.Frame()
	before := w
	if err := tt.Click("Copy key"); err != nil {
		t.Fatal(err)
	}
	if tt.Clipboard() != "sk-live-123" || n != 1 {
		t.Fatalf("clipboard %q, copies %d", tt.Clipboard(), n)
	}
	if !tt.HasText("Copied") || !slices.Contains(tt.Announcements(), "Copied") {
		t.Fatalf("texts %q", tt.Texts())
	}
	tt.Frame()
	if w != before {
		t.Fatalf("width %v changed from %v", w, before)
	}
}

func TestCopyButtonKeyboardRevertAndDisabled(t *testing.T) {
	dis := false
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{Locale: ZhCN})
		CopyButton(c, "abc", CopyButtonOptions{Disabled: dis})
	}, 300, 200)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeySpace)
	if tt.Clipboard() != "abc" {
		t.Fatalf("Space copied %q", tt.Clipboard())
	}
	if _, ok := tt.Find("已复制"); !ok {
		t.Fatal("icon button lacks the copied name")
	}
	time.Sleep(copyButtonShown + 50*time.Millisecond)
	tt.Frame()
	if _, ok := tt.Find("复制"); !ok {
		t.Fatal("did not revert to Copy")
	}
	tt.SetClipboard("")
	dis = true
	tt.Frame()
	_ = tt.Click("复制")
	if tt.Clipboard() != "" {
		t.Fatal("disabled button copied")
	}
}

func TestCopyButtonSnapshot(t *testing.T) {
	snapshot(t, "copy_button", 360, 80, func(c *ui.Context) {
		ui.Row(c).Padding(16).Gap(12).Children(func() {
			CopyButton(c, "x", CopyButtonOptions{})
			CopyButton(c, "x", CopyButtonOptions{Label: "Copy link"})
			CopyButton(c, "x", CopyButtonOptions{Label: "Copy", Disabled: true})
		})
	})
}
