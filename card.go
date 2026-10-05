package mujica

import (
	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
	"mujicaui/theme"
)

// CardOptions configure a Card.
type CardOptions struct {
	// Title heads the card in the serif face; empty shows no header.
	Title string
	// Footer builds the bottom slot, below a divider; nil shows none.
	Footer func()
	// Decorated draws a double border and corner ornaments in Ornament.
	Decorated bool
}

// cardInset is how far the inner line of a decorated card sits inside.
const cardInset float32 = 4

// Card groups content on a surface with an outline. It adds no click behavior.
func Card(c *ui.Context, opts CardOptions, body func()) *ui.Element {
	k := Tokens(c)
	u := Density(c).Unit()
	pad := u * 4
	if opts.Decorated {
		pad += cardInset + 2
	}
	card := ui.Column(c).Background(k.Surface).Radius(theme.CardRadius).Gap(u * 3).Padding(pad).Role(ui.RoleGroup)
	if opts.Title != "" {
		card.Label(opts.Title)
	}
	if opts.Decorated {
		card.Border(theme.BorderWidth, k.Ornament).DrawOver(func(p *ui.Painter, r ui.Rect) {
			cardOrnament(p, r, k.Ornament)
		})
	} else {
		card.Border(theme.BorderWidth, k.Border)
	}
	card.Children(func() {
		if opts.Title != "" {
			ui.Text(c, opts.Title).Font(theme.SerifFont).FontSize(theme.H3Size).FontWeight(600).LineHeight(1.3)
			if opts.Decorated {
				ui.Box(c).Height(7).FillWidth().Shrink(0).Draw(func(p *ui.Painter, r ui.Rect) {
					cy := r.Y + r.H/2
					p.Line(r.X, cy, r.X+r.W/2-6, cy, 1, k.Ornament)
					p.Line(r.X+r.W/2+6, cy, r.X+r.W, cy, 1, k.Ornament)
					p.FillPath(internal.Diamond(internal.Snap(p, r.X+r.W/2), internal.Snap(p, cy), 7), k.Ornament)
				})
			}
		}
		ui.Column(c).Gap(u * 2).Children(body)
		if opts.Footer != nil {
			ui.Box(c).Height(1).FillWidth().Shrink(0).Background(k.Border)
			ui.Row(c).Gap(u * 2).Justify(ui.End).AlignItems(ui.Center).Wrap().Children(opts.Footer)
		}
	})
	return card
}

// cardOrnament draws the inner line and the four corner diamonds.
func cardOrnament(p *ui.Painter, r ui.Rect, col ui.Color) {
	in := ui.Rect{X: r.X + cardInset, Y: r.Y + cardInset, W: r.W - 2*cardInset, H: r.H - 2*cardInset}
	p.Stroke(in, col, 1, 1)
	for _, pt := range [4][2]float32{{in.X, in.Y}, {in.X + in.W, in.Y}, {in.X, in.Y + in.H}, {in.X + in.W, in.Y + in.H}} {
		p.FillPath(internal.Diamond(internal.Snap(p, pt[0]), internal.Snap(p, pt[1]), 9), col)
	}
}
