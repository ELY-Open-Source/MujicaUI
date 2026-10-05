package chat

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/agent"
	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

func TestArtifactCardOpens(t *testing.T) {
	opens := 0
	tt := chatTester(480, 160, func(c *ui.Context) {
		if ArtifactCard(c, ArtifactCardOptions{Kind: agent.ArtifactDocument, Title: "Decree", Detail: "2 pages"}).Changed() {
			opens++
		}
	})
	if !tt.HasText("Document · 2 pages") {
		t.Fatalf("texts %q", tt.Texts())
	}
	_ = tt.Click("Decree")
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	if opens != 2 {
		t.Fatalf("opens %d", opens)
	}
}

func TestArtifactCardStreamingDisabledAndChecks(t *testing.T) {
	opens := 0
	tt := chatTester(480, 160, func(c *ui.Context) {
		if ArtifactCard(c, ArtifactCardOptions{Kind: agent.ArtifactWeb, Title: "Menu", Streaming: true, Disabled: true}).Changed() {
			opens++
		}
	})
	_ = tt.Click("Menu")
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	if opens != 0 || !tt.HasText("Web page") {
		t.Fatalf("opens %d texts %q", opens, tt.Texts())
	}
	for _, opts := range []ArtifactCardOptions{{}, {Title: "x", Kind: 7}} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("%+v: no panic", opts)
				}
			}()
			ui.Render(func(c *ui.Context) {
				core.Use(c, core.Settings{})
				ArtifactCard(c, opts)
			}, 100, 100, 1)
		}()
	}
}

func TestArtifactCardSnapshot(t *testing.T) {
	snapshot(t, "artifact_card", 460, 210, func(c *ui.Context) {
		ui.Column(c).Padding(16).Gap(8).Children(func() {
			ArtifactCard(c, ArtifactCardOptions{Kind: agent.ArtifactCode, Title: "seating.go", Detail: "Go · 120 lines"})
			ArtifactCard(c, ArtifactCardOptions{Kind: agent.ArtifactDocument, Title: "Royal decree", Streaming: true})
			ArtifactCard(c, ArtifactCardOptions{Kind: agent.ArtifactWeb, Title: "Ball invitation", Disabled: true})
		})
	})
}
