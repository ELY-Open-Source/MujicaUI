package mujica

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"mujicaui/theme"
)

func TestUseLineHeight(t *testing.T) {
	var h float32
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		h = ui.Text(c, "one\ntwo").Bounds().H
	}, 300, 200)
	tt.Frame()
	if want := 2 * theme.BodySize * theme.LineHeight; h < want-0.5 || h > want+0.5 {
		t.Fatalf("two lines are %v DIPs, want %v", h, want)
	}
}

func TestUseRejectsInvalidSettings(t *testing.T) {
	for _, s := range []Settings{{Mode: 9}, {Density: 9}, {Locale: "fr-FR"}} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("%+v accepted", s)
				}
			}()
			ui.NewTester(func(c *ui.Context) { Use(c, s) }, 100, 100)
		}()
	}
}

func TestTokensWithoutUsePanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("Tokens before Use did not panic")
		}
	}()
	ui.NewTester(func(c *ui.Context) { Tokens(c) }, 100, 100)
}
