package diagram

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/theme"
)

func TestTokensFollowTheTheme(t *testing.T) {
	for name, k := range map[string]theme.Tokens{"light": theme.Light(), "dark": theme.Dark()} {
		pal := tokens(k)
		for token, want := range map[string]ui.Color{"--_bg": k.Background, "--_fg": k.Text, "--_m": k.TextMuted, "--_l": k.Text,
			"--_s": k.Surface, "--_a": k.SurfaceHover, "--_b": k.ControlBorder, "--_c": k.AccentText, "--pele-critical": k.Danger} {
			if got, ok, err := pal.color("var("+token+")", k.Text); !ok || err != nil || got != want {
				t.Errorf("%s %s = %v %v %v, want %v", name, token, got, ok, err, want)
			}
		}
		seen := map[ui.Color]bool{}
		for i, c := range pal.series {
			if seen[c] {
				t.Errorf("%s series %d repeats %v", name, i+1, c)
			}
			seen[c] = true
			if r := theme.Contrast(c, k.Surface); r < 3 {
				t.Errorf("%s series %d contrast %.2f against the surface, want 3", name, i+1, r)
			}
		}
		if pal.radius != theme.CardRadius {
			t.Errorf("%s radius %v", name, pal.radius)
		}
	}
	sc, err := render("flowchart LR\n  Score[Write score] --> Stage[Take the stage]", tokens(theme.Dark()), theme.BodySize, "")
	if err != nil {
		t.Fatal(err)
	}
	paths, texts := marks(sc)
	if k := theme.Dark(); paths[2].fill != k.Surface || paths[2].stroke != k.ControlBorder || texts[0].fill != k.Text {
		t.Errorf("box %v edged %v, label %v: want the dark tokens", paths[2].fill, paths[2].stroke, texts[0].fill)
	}
}
