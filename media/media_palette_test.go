package media

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// TestMediaChartInksContrast keeps every chart ink at 3:1 against Surface in both themes.
func TestMediaChartInksContrast(t *testing.T) {
	for _, k := range []theme.Tokens{theme.Light(), theme.Dark()} {
		inks := map[string]ui.Color{
			"played (AccentText)": k.AccentText, "unplayed (TextMuted)": k.TextMuted, "cursor (Text)": k.Text,
			"mic safe": micZone(k, 0), "mic loud": micZone(k, micSegments*7/10), "mic clip": micZone(k, micSegments-1),
		}
		for name, c := range inks {
			if r := theme.Contrast(c, k.Surface); r < 3 {
				t.Errorf("dark=%v %s: %.2f:1 against Surface", k.Dark, name, r)
			}
		}
	}
}
