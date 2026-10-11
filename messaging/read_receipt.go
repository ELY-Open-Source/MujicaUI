package messaging

import (
	"fmt"
	"strings"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/display"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// ReceiptState tracks delivery; its zero value means sent.
type ReceiptState int

const (
	ReceiptSent ReceiptState = iota
	ReceiptDelivered
	ReceiptRead
)

// ReadReceiptOptions configure a ReadReceipt.
type ReadReceiptOptions struct {
	// Readers supplies avatars for the read state.
	Readers []display.AvatarItem
	// Max avatars show before "+N"; zero means 3.
	Max int
}

// ReadReceipt displays delivery marks, text and reader avatars.
func ReadReceipt(c *ui.Context, state ReceiptState, opts ReadReceiptOptions) *ui.Element {
	if state < ReceiptSent || state > ReceiptRead || opts.Max < 0 {
		panic(fmt.Sprintf("mujica: bad ReadReceipt state %d or Max %d", state, opts.Max))
	}
	limit := opts.Max
	if limit == 0 {
		limit = 3
	}
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	word := core.Msg(c, []string{"read_receipt.sent", "read_receipt.delivered", "read_receipt.read"}[state])
	label := word
	if state == ReceiptRead && len(opts.Readers) > 0 {
		names := make([]string, len(opts.Readers))
		for i, r := range opts.Readers {
			names[i] = r.Name
		}
		label = fmt.Sprintf(core.Msg(c, "read_receipt.readBy"), strings.Join(names, ", "))
	}
	ink := k.TextMuted
	if state == ReceiptRead {
		ink = k.AccentText
	}
	row := ui.Row(c).Gap(u * 1.5).AlignItems(ui.Center).Shrink(0).Role(ui.RoleStatus).Label(label).Tooltip(label)
	row.Children(func() {
		ui.Box(c).Size(18, 12).Shrink(0).Draw(func(p *ui.Painter, r ui.Rect) { receiptChecks(p, r, state != ReceiptSent, ink) })
		ui.Text(c, word).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(ink).SingleLine()
		if state != ReceiptRead || len(opts.Readers) == 0 {
			return
		}
		ui.Row(c).Children(func() {
			for i, r := range opts.Readers[:min(limit, len(opts.Readers))] {
				cell := ui.Box(c).Key(i).Padding(1).Radius(9).Background(k.Background)
				if i > 0 {
					cell.Margin(0, 0, 0, -5)
				}
				cell.Children(func() { display.Avatar(c, r.Name, display.AvatarOptions{Image: r.Image, Size: 16}) })
			}
		})
		if rest := len(opts.Readers) - limit; rest > 0 {
			ui.Textf(c, "+%d", rest).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted)
		}
	})
	return row
}

// receiptChecks draws one check, or two overlapping checks.
func receiptChecks(p *ui.Painter, r ui.Rect, double bool, col ui.Color) {
	check := func(x float32) {
		var path ui.Path
		path.MoveTo(internal.Snap(p, x), internal.Snap(p, r.Y+r.H*0.55))
		path.LineTo(internal.Snap(p, x+r.H*0.35), internal.Snap(p, r.Y+r.H*0.9))
		path.LineTo(internal.Snap(p, x+r.H*0.95), internal.Snap(p, r.Y+r.H*0.15))
		p.StrokePath(&path, 1.6, col)
	}
	check(r.X)
	if double {
		check(r.X + r.H*0.5)
	}
}
