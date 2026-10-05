package mujica

import (
	"errors"
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestAttachmentChipStatesAndActions(t *testing.T) {
	a := Attachment{Name: "guests.csv", Size: 4096, Uploading: true, Progress: 0.4}
	removed, retried := 0, 0
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		r := AttachmentChip(c, a)
		if r.Removed() {
			removed++
		}
		if r.Retried() {
			retried++
		}
	}, 400, 120)
	if !tt.HasText("Uploading 40%") {
		t.Fatalf("progress: %q", tt.Texts())
	}
	if _, ok := tt.Find("Retry upload of guests.csv"); ok {
		t.Fatal("retry shown while uploading")
	}
	a.Uploading, a.Err = false, errors.New("413")
	tt.Frame()
	_ = tt.Click("Retry upload of guests.csv")
	_ = tt.Click("Remove guests.csv")
	if retried != 1 || removed != 1 || !tt.HasText("Upload failed") {
		t.Fatalf("retried %d removed %d texts %q", retried, removed, tt.Texts())
	}
	a.Err = nil
	tt.Frame()
	if !tt.HasText("4.0 KB") {
		t.Fatalf("size: %q", tt.Texts())
	}
}

func TestAttachmentChipBadProgress(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic for progress 2")
		}
	}()
	ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		AttachmentChip(c, Attachment{Name: "a", Uploading: true, Progress: 2})
	}, 100, 100)
}

func TestAttachmentChipSnapshot(t *testing.T) {
	snapshot(t, "attachment_chip", 520, 80, func(c *ui.Context) {
		ui.Row(c).Padding(12).Gap(8).Children(func() {
			AttachmentChip(c, Attachment{Name: "guests.csv", Size: 4096})
			AttachmentChip(c, Attachment{Name: "portrait.png", Size: 812000, Uploading: true, Progress: 0.62})
		})
	})
}

func TestAttachmentChipKeyboardAndDisabled(t *testing.T) {
	off, removed, retried := false, 0, 0
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		ui.Row(c).Disabled(off).Children(func() {
			r := AttachmentChip(c, Attachment{Name: "guests.csv", Size: 0, Err: errors.New("upload failed")})
			if r.Removed() {
				removed++
			}
			if r.Retried() {
				retried++
			}
		})
	}, 400, 160)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeySpace)
	if retried != 1 || removed != 1 {
		t.Fatalf("retry %d remove %d", retried, removed)
	}
	off = true
	tt.Frame()
	tt.Key(0, ui.KeyEnter)
	_ = tt.Click("Retry upload of guests.csv")
	if retried != 1 || removed != 1 {
		t.Fatal("disabled attachment handled input")
	}
}
