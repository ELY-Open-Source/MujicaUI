package mujica

import (
	"fmt"
	"log/slog"
	"math"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/theme"
)

// Attachment is a file attached to a prompt; the caller uploads it.
type Attachment struct {
	// Name is required.
	Name string
	// Type is the MIME type; empty picks the icon from the name.
	Type string
	Size int64
	// Uploading shows Progress, 0 to 1.
	Uploading bool
	Progress  float64
	// Err shows a failed upload with retry.
	Err error
}

// AttachmentChipResult is a chip's element and the user's request this frame.
type AttachmentChipResult struct {
	Element          *ui.Element
	removed, retried bool
}

// Removed reports a press of the remove button.
func (r AttachmentChipResult) Removed() bool { return r.removed }

// Retried reports a press of retry on a failed upload.
func (r AttachmentChipResult) Retried() bool { return r.retried }

// AttachmentChip shows an attachment's type, name, size, upload progress or failure, and remove.
func AttachmentChip(c *ui.Context, a Attachment) AttachmentChipResult {
	if a.Name == "" || a.Size < 0 {
		panic("mujica: Attachment needs a Name and a non-negative Size")
	}
	if a.Uploading && !(a.Progress >= 0 && a.Progress <= 1) {
		panic("mujica: Attachment Progress must be from 0 to 1")
	}
	k := Tokens(c)
	u := Density(c).Unit()
	var r AttachmentChipResult
	border := k.Border
	if a.Err != nil {
		border = k.Danger
	}
	r.Element = ui.Row(c).Gap(u*2).AlignItems(ui.Center).Padding(u, u, u, u*2).Width(220).
		Radius(theme.ControlRadius).Background(k.Surface).Border(theme.BorderWidth, border).Role(ui.RoleGroup).Label(a.Name)
	disabled := r.Element.IsDisabled()
	r.Element.Children(func() {
		icon, color := fileMessageIcon(a.Name, a.Type), k.AccentText
		if a.Err != nil {
			icon, color = "circle-alert", k.Danger
		}
		ui.Icon(c, icons.Must(icon)).Size(20, 20).Shrink(0).TextColor(color)
		ui.Column(c).Grow(1).MinWidth(0).Gap(u / 2).Children(func() {
			ui.Text(c, a.Name).SingleLine().FontSize(FontSize(c, theme.CaptionSize)).FontWeight(600)
			small := FontSize(c, theme.CaptionSize)
			switch {
			case a.Err != nil:
				ui.Text(c, msg(c, "attachment_chip.failed")).SingleLine().FontSize(small).TextColor(k.Danger)
			case a.Uploading:
				pct := int(math.Round(a.Progress * 100))
				ui.Text(c, fmt.Sprintf(msg(c, "attachment_chip.uploading"), pct)).SingleLine().FontSize(small).TextColor(k.TextMuted).FontFeatures("tnum")
				frac := float32(a.Progress)
				ui.Box(c).Height(3).Radius(1.5).Background(k.Border).Role(ui.RoleProgress).Label(fmt.Sprintf(msg(c, "attachment_chip.uploading"), pct)).
					Draw(func(p *ui.Painter, b ui.Rect) {
						p.Fill(ui.Rect{X: b.X, Y: b.Y, W: b.W * frac, H: b.H}, k.AccentText, 1.5)
					})
			default:
				ui.Text(c, fileMessageSize(c, a.Size)).SingleLine().FontSize(small).TextColor(k.TextMuted)
			}
		})
		if a.Err != nil && iconAction(c, icons.Must("refresh-cw"), fmt.Sprintf(msg(c, "attachment_chip.retry"), a.Name)).Clicked() && !disabled {
			r.retried = true
			slog.Debug("mujica: attachment_chip retry")
		}
		if iconAction(c, icons.Must("x"), fmt.Sprintf(msg(c, "attachment_chip.remove"), a.Name)).Clicked() && !disabled {
			r.removed = true
			slog.Debug("mujica: attachment_chip remove")
		}
	})
	return r
}
