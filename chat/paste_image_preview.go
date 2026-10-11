package chat

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/icons"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// PastedImage contains a caller-decoded clipboard or dropped image.
type PastedImage struct {
	// ID is the image's stable identity; required.
	ID     string
	Source ui.ImageSource
	// Name describes it for assistive technology and the remove button.
	Name string
}

// pasteImageSize is the thumbnail's edge in DIPs.
const pasteImageSize float32 = 72

// PasteImagePreview displays removable image thumbnails.
func PasteImagePreview(c *ui.Context, images []PastedImage) RemoveResult {
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	var r RemoveResult
	r.Element = ui.Row(c).Gap(u * 2).Wrap().Role(ui.RoleList).Label(core.Msg(c, "paste_image_preview.label"))
	disabled := r.Element.IsDisabled()
	r.Element.Children(func() {
		for _, im := range images {
			if im.ID == "" || im.Source == nil {
				panic("mujica: PastedImage needs ID and Source")
			}
			ui.Box(c).Key(im.ID).Size(pasteImageSize, pasteImageSize).Role(ui.RoleListItem).Children(func() {
				ui.Image(c, im.Source).Fit(ui.Cover).Size(pasteImageSize, pasteImageSize).Radius(theme.ControlRadius).
					Border(theme.BorderWidth, k.Border).Label(im.Name)
				x := core.IconAction(c, icons.Must("x"), fmt.Sprintf(core.Msg(c, "paste_image_preview.remove"), im.Name)).
					Absolute().Top(-u * 3).Right(-u * 3).Radius(theme.MinHitSize / 2)
				// A 20 DIP disc inside the 28 DIP hit area.
				x.Draw(func(p *ui.Painter, r ui.Rect) {
					d := ui.Rect{X: r.X + 4, Y: r.Y + 4, W: r.W - 8, H: r.H - 8}
					p.Fill(d, k.Surface, d.W/2)
					p.Stroke(d, k.ControlBorder, d.W/2, theme.BorderWidth)
				})
				if x.Clicked() && !disabled {
					r.removed = im.ID
					slog.Debug("mujica: paste_image_preview remove", "id", im.ID)
				}
			})
		}
	})
	return r
}
