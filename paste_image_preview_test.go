package mujica

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func pasteImages() []PastedImage {
	return []PastedImage{
		{ID: "a", Source: imageGridPicture("#7D2034"), Name: "screenshot 1"},
		{ID: "b", Source: imageGridPicture("#355E80"), Name: "screenshot 2"},
	}
}

func TestPasteImagePreviewRemoves(t *testing.T) {
	imgs, off := pasteImages(), false
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		ui.Column(c).Padding(16).Disabled(off).Children(func() {
			if id, ok := PasteImagePreview(c, imgs).Removed(); ok {
				for i, im := range imgs {
					if im.ID == id {
						imgs = append(imgs[:i], imgs[i+1:]...)
						break
					}
				}
			}
		})
	}, 400, 150)
	if _, ok := tt.Find("screenshot 2"); !ok {
		t.Fatal("thumbnail not named")
	}
	_ = tt.Click("Remove image screenshot 1")
	if len(imgs) != 1 || imgs[0].ID != "b" {
		t.Fatalf("left %v", imgs)
	}
	off = true
	_ = tt.Click("Remove image screenshot 2")
	if len(imgs) != 1 {
		t.Fatal("disabled preview removed an image")
	}
}

func TestPasteImagePreviewNeedsSource(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic without a source")
		}
	}()
	ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		PasteImagePreview(c, []PastedImage{{ID: "x"}})
	}, 100, 100)
}

func TestPasteImagePreviewSnapshot(t *testing.T) {
	snapshot(t, "paste_image_preview", 240, 110, func(c *ui.Context) {
		ui.Column(c).Padding(20).Children(func() { PasteImagePreview(c, pasteImages()) })
	})
}
