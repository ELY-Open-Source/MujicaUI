package mujica

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestSourceCardOpensAndLights(t *testing.T) {
	var st CitationState
	opened, off := 0, false
	var bg []byte
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		ui.Column(c).Padding(16).Width(360).Disabled(off).Children(func() {
			if SourceCard(c, 1, sourcesSample[0], &st).Clicked() {
				opened++
			}
		})
	}, 400, 200)
	for _, s := range []string{"[1]", "example.com", "2025-05-04", "The Palace of Versailles", "Seat of the court from 1682."} {
		if !tt.HasText(s) {
			t.Errorf("missing %q in %q", s, tt.Texts())
		}
	}
	bg = encode(t, tt.Image())
	_ = tt.Click("Open The Palace of Versailles")
	tt.Key(0, ui.KeyEnter)
	if opened != 2 {
		t.Fatalf("opened %d", opened)
	}
	st.Active = 1
	tt.Move(390, 190)
	tt.Frame()
	if string(encode(t, tt.Image())) == string(bg) {
		t.Fatal("active card does not light up")
	}
	off = true
	_ = tt.Click("Open The Palace of Versailles")
	if opened != 2 {
		t.Fatal("disabled card opened")
	}
}

func TestSourceCardNeedsURL(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic without a URL")
		}
	}()
	var st CitationState
	ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		SourceCard(c, 1, Source{Title: "x"}, &st)
	}, 100, 100)
}

func TestSourceCardSnapshot(t *testing.T) {
	snapshot(t, "source_card", 400, 260, func(c *ui.Context) {
		st := ui.Local(c.Root(), "st", func() CitationState { return CitationState{Active: 2} })
		ui.Column(c).Padding(16).Gap(8).Children(func() {
			SourceCard(c, 1, sourcesSample[0], st)
			SourceCard(c, 2, sourcesSample[1], st)
		})
	})
}
