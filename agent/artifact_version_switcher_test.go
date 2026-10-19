package agent

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
)

var artifactVersionsSample = []ArtifactVersion{
	{Label: "v1", Note: "First draft"},
	{Label: "v2", Note: "Added the quartet"},
	{Label: "v3", Note: "Moved the wine table"},
}

func versionTester(cur *int, opts ArtifactVersionSwitcherOptions, got *[]ArtifactVersionResult) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := ArtifactVersionSwitcher(c, cur, artifactVersionsSample, opts)
		if r.Changed || r.Compare {
			*got = append(*got, r)
		}
	}, 600, 120)
}

func TestArtifactVersionSwitcherStepsAndCompare(t *testing.T) {
	cur := 2
	var got []ArtifactVersionResult
	tt := versionTester(&cur, ArtifactVersionSwitcherOptions{}, &got)
	if !tt.HasText("3 of 3") || !tt.HasText("Moved the wine table") || !tt.HasText("Compare with v2") {
		t.Fatalf("texts %q", tt.Texts())
	}
	_ = tt.Click("Compare with v2")
	if len(got) != 1 || !got[0].Compare || cur != 2 {
		t.Fatalf("compare: %+v", got)
	}
	_ = tt.Click("Older version")
	_ = tt.Click("Older version")
	if cur != 0 || tt.HasText("Compare with v1") || !tt.HasText("First draft") {
		t.Fatalf("older: cur %d texts %q", cur, tt.Texts())
	}
	_ = tt.Click("Older version") // disabled at the first version
	if cur != 0 || len(got) != 3 {
		t.Fatalf("past the first: cur %d events %d", cur, len(got))
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	if cur != 1 {
		t.Fatalf("Right: cur %d", cur)
	}
}

func TestArtifactVersionSwitcherDisabledAndRange(t *testing.T) {
	cur := 1
	var got []ArtifactVersionResult
	tt := versionTester(&cur, ArtifactVersionSwitcherOptions{Disabled: true}, &got)
	_ = tt.Click("Older version")
	_ = tt.Click("Compare with v1")
	if cur != 1 || len(got) != 0 {
		t.Fatalf("disabled moved: cur %d %+v", cur, got)
	}
	defer func() {
		if recover() == nil {
			t.Fatal("out-of-range current did not panic")
		}
	}()
	cur = 5
	versionTester(&cur, ArtifactVersionSwitcherOptions{}, &got)
}

func TestArtifactVersionSwitcherSnapshot(t *testing.T) {
	snapshot(t, "artifact_version_switcher", 560, 90, func(c *ui.Context) {
		cur := 1
		ui.Column(c).Padding(16).Children(func() {
			ArtifactVersionSwitcher(c, &cur, artifactVersionsSample, ArtifactVersionSwitcherOptions{})
		})
	})
}
