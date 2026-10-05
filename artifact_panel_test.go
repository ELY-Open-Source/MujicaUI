package mujica

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

var artifactCodeSample = Artifact{Title: "seating.go", Language: "go", Content: "package court\n\n// Seats maps guests to tables.\nfunc Seats() map[string]int {\n\treturn map[string]int{\"Sakiko\": 1}\n}\n"}

func artifactTester(art Artifact, opts ArtifactPanelOptions, got *[]ArtifactPanelResult) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		r := ArtifactPanel(c, art, opts)
		if r.Copied || r.Exported {
			*got = append(*got, r)
		}
	}, 600, 460)
}

func TestArtifactPanelCodeCopyExport(t *testing.T) {
	var got []ArtifactPanelResult
	tt := artifactTester(artifactCodeSample, ArtifactPanelOptions{}, &got)
	if !tt.HasText("// Seats maps guests to tables.") || tt.HasText("Preview") {
		t.Fatalf("texts %q", tt.Texts())
	}
	_ = tt.Click("Copy")
	if len(got) != 1 || !got[0].Copied || tt.Clipboard() != artifactCodeSample.Content || !tt.HasText("Copied") {
		t.Fatalf("copy: %+v clipboard %q", got, tt.Clipboard())
	}
	_ = tt.Click("Export")
	if len(got) != 2 || !got[1].Exported {
		t.Fatalf("export: %+v", got)
	}
}

func TestArtifactPanelPreviewSourceSwitch(t *testing.T) {
	doc := Artifact{Title: "Invitation", Kind: ArtifactDocument, Content: "# You are invited\n"}
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		ArtifactPanel(c, doc, ArtifactPanelOptions{Preview: func() { ui.Text(c, "Rendered invitation") }})
	}, 600, 460)
	if !tt.HasText("Rendered invitation") || tt.HasText("# You are invited") {
		t.Fatalf("preview first: %q", tt.Texts())
	}
	_ = tt.Click("Source")
	if !tt.HasText("# You are invited") || tt.HasText("Rendered invitation") {
		t.Fatalf("source: %q", tt.Texts())
	}
	tt.Key(0, ui.KeyLeft) // the switch is a focus group
	tt.Key(0, ui.KeyEnter)
	if !tt.HasText("Rendered invitation") {
		t.Fatal("keyboard back to preview failed")
	}
}

func TestArtifactPanelDisabledAndRules(t *testing.T) {
	var got []ArtifactPanelResult
	tt := artifactTester(artifactCodeSample, ArtifactPanelOptions{Disabled: true}, &got)
	_ = tt.Click("Copy")
	_ = tt.Click("Export")
	if len(got) != 0 || tt.Clipboard() != "" {
		t.Fatalf("disabled panel acted: %+v", got)
	}
	defer func() {
		if recover() == nil {
			t.Fatal("web artifact without Preview did not panic")
		}
	}()
	artifactTester(Artifact{Title: "x", Kind: ArtifactWeb}, ArtifactPanelOptions{}, &got)
}

func TestArtifactPanelSnapshot(t *testing.T) {
	snapshot(t, "artifact_panel", 560, 300, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() {
			ArtifactPanel(c, artifactCodeSample, ArtifactPanelOptions{Height: 170, Preview: func() {}})
		})
	})
}
