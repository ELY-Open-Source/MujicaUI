package diagram

import (
	"context"

	pele "github.com/Broderick-Westrope/pele-go"
	_ "github.com/Broderick-Westrope/pele-go/all"
	"github.com/egoist/mygo/ui"
)

// measurer measures labels in the fonts that draw them.
type measurer struct{ font string }

func (m measurer) Width(text string, size float64, st pele.Style) float64 {
	f := ui.Font{Family: m.font, Size: float32(size), Italic: st&pele.Italic != 0}
	if st&pele.Mono != 0 {
		f.Family = "monospace"
	}
	if st&(pele.Sub|pele.Sup) != 0 {
		f.Size /= 1.2
	}
	if st&pele.Bold != 0 {
		f.Weight = bold
	}
	w, _ := ui.ShapeText(text, f, 0).Size()
	return float64(w)
}

// render lays Mermaid source out and reads the drawing.
func render(source string, pal palette, size float32, font string) (*scene, error) {
	off := false
	res, err := pele.Render(context.Background(), source, pele.Options{
		Measurer: measurer{font}, FontSize: pele.Ptr(float64(size)), Links: &off, Images: &off})
	if err != nil {
		return nil, err
	}
	return parse(res.SVG, pal, size, font)
}
