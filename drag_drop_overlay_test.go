package mujica

import (
	"slices"
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestDragDropOverlayIdleShowsContentOnly(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		DragDropOverlay(c, DragDropOverlayOptions{}, func() { ui.Text(c, "Chat") })
	}, 300, 200)
	if !tt.HasText("Chat") || tt.HasText("Drop files to attach") {
		t.Fatalf("idle overlay: %q", tt.Texts())
	}
}

func TestDragDropOverlaySplitsByAccept(t *testing.T) {
	images := func(p string) bool { return strings.HasSuffix(p, ".png") }
	ok, no := dragDropSplit([]string{"/a.png", "/b.exe", "/c.png"}, images)
	if !slices.Equal(ok, []string{"/a.png", "/c.png"}) || !slices.Equal(no, []string{"/b.exe"}) {
		t.Fatalf("ok %v no %v", ok, no)
	}
	ok, no = dragDropSplit([]string{"/b.exe"}, nil)
	if len(ok) != 1 || no != nil {
		t.Fatalf("nil accept: ok %v no %v", ok, no)
	}
}

func TestDragDropOverlayNeedsContent(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic without content")
		}
	}()
	ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		DragDropOverlay(c, DragDropOverlayOptions{}, nil)
	}, 100, 100)
}

// TestDragDropOverlaySnapshot draws the mask shown while files are over the area.
func TestDragDropOverlaySnapshot(t *testing.T) {
	snapshot(t, "drag_drop_overlay", 420, 220, func(c *ui.Context) {
		ui.Box(c).Fill().Padding(12).Children(func() {
			ui.Text(c, "A conversation sits under the mask.")
			dragDropMask(c, msg(c, "drag_drop_overlay.hint"))
		})
	})
}
