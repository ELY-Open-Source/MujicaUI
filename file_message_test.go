package mujica

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestFileMessageOpenDownloadAndSize(t *testing.T) {
	opened, downloaded, off := 0, 0, false
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		ui.Column(c).Disabled(off).Children(func() {
			r := FileMessage(c, FileMessageOptions{Name: "treaty.pdf", Size: 2_621_440})
			if r.Opened() {
				opened++
			}
			if r.Downloaded() {
				downloaded++
			}
		})
	}, 400, 200)
	if !tt.HasText("2.5 MB") || !tt.HasText("treaty.pdf") {
		t.Fatalf("texts %q", tt.Texts())
	}
	_ = tt.Click("Open treaty.pdf, 2.5 MB")
	_ = tt.Click("Download treaty.pdf")
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	if opened != 2 || downloaded != 1 {
		t.Fatalf("opened %d downloaded %d", opened, downloaded)
	}
	off = true
	_ = tt.Click("Open treaty.pdf, 2.5 MB")
	_ = tt.Click("Download treaty.pdf")
	if opened != 2 || downloaded != 1 {
		t.Fatalf("disabled file acted: opened %d downloaded %d", opened, downloaded)
	}
}

func TestFileMessageIconsAndSizes(t *testing.T) {
	for name, want := range map[string]string{"a.png": "image", "b.zip": "file-archive", "c.go": "file-code", "d.md": "file-text", "e.bin": "file", "f.c": "file-code", "g.s": "file"} {
		if got := fileMessageIcon(name, ""); got != want {
			t.Errorf("%s: icon %s, want %s", name, got, want)
		}
	}
	if fileMessageIcon("x", "image/heic") != "image" {
		t.Error("MIME type ignored")
	}
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		for _, n := range []int64{0, 1023, 1024, 5 << 30} {
			ui.Text(c, fileMessageSize(c, n))
		}
	}, 200, 200)
	for _, s := range []string{"0 B", "1023 B", "1.0 KB", "5.0 GB"} {
		if !tt.HasText(s) {
			t.Errorf("missing size %q in %q", s, tt.Texts())
		}
	}
}

func TestFileMessageNeedsName(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic without a name")
		}
	}()
	ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		FileMessage(c, FileMessageOptions{})
	}, 100, 100)
}

func TestFileMessageSnapshot(t *testing.T) {
	snapshot(t, "file_message", 400, 160, func(c *ui.Context) {
		ui.Column(c).Padding(16).Gap(8).Children(func() {
			FileMessage(c, FileMessageOptions{Name: "seating-plan.xlsx", Size: 48_213})
			FileMessage(c, FileMessageOptions{Name: "portrait.png", Size: 3_400_000})
		})
	})
}
