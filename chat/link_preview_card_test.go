package chat

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

var linkPreviewSample = LinkPreview{
	URL: "https://example.com/versailles", Title: "The Palace of Versailles",
	Description: "A royal château in the Île-de-France region, seat of the court from 1682.", Site: "example.com",
}

func TestLinkPreviewCardShowsMetadataAndOpens(t *testing.T) {
	opened := 0
	p := linkPreviewSample
	p.Thumbnail = imageGridPicture("#7D2034")
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if LinkPreviewCard(c, p).Clicked() {
			opened++
		}
	}, 500, 200)
	for _, s := range []string{"example.com", "The Palace of Versailles", "A royal château"} {
		if !tt.HasText(s) {
			t.Errorf("missing %q in %q", s, tt.Texts())
		}
	}
	_ = tt.Click("The Palace of Versailles")
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	if opened != 2 {
		t.Fatalf("opened %d", opened)
	}
}

func TestLinkPreviewCardSiteFallsBackToURL(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		LinkPreviewCard(c, LinkPreview{URL: "https://a.io/x", Title: "X"})
	}, 400, 200)
	if !tt.HasText("https://a.io/x") {
		t.Fatalf("texts %q", tt.Texts())
	}
	defer func() {
		if recover() == nil {
			t.Fatal("no panic without a title")
		}
	}()
	ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		LinkPreviewCard(c, LinkPreview{URL: "https://a.io"})
	}, 100, 100)
}

func TestLinkPreviewCardSnapshot(t *testing.T) {
	p := linkPreviewSample
	p.Thumbnail = imageGridPicture("#AD9164")
	snapshot(t, "link_preview_card", 460, 130, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() { LinkPreviewCard(c, p) })
	})
}
