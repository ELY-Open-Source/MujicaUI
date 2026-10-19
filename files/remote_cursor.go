package files

import (
	"fmt"
	"hash/fnv"
	"strings"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/code"
	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/internal"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

// RemoteCaret describes a collaborator's caret and selection offsets.
type RemoteCaret struct {
	ID, Name string
	// Offset and Anchor bound the selection, equal for a caret.
	Offset, Anchor int
}

// RemoteCursorOptions configure RemoteCursor.
type RemoteCursorOptions struct {
	// Caret maps text offsets into overlay-relative rectangles; required.
	Caret func(offset int) ui.Rect
	// Selection maps offset ranges into overlay-relative rectangles; optional.
	Selection func(from, to int) []ui.Rect
}

// RemoteCursor overlays collaborator carets, selections and names; input passes through.
// Place it filling the editor's box, after the editor.
func RemoteCursor(c *ui.Context, carets []RemoteCaret, opts RemoteCursorOptions) ui.Element {
	if opts.Caret == nil {
		panic("mujica: RemoteCursor needs Caret")
	}
	k := core.Tokens(c)
	names := make([]string, len(carets))
	for i, rc := range carets {
		if rc.ID == "" || rc.Offset < 0 || rc.Anchor < 0 {
			panic("mujica: RemoteCaret needs an ID and offsets not negative")
		}
		names[i] = rc.Name
	}
	size := core.FontSize(c, theme.CaptionSize)
	e := ui.Box(c).Absolute().Top(0).Left(0).Right(0).Bottom(0).PassThrough().Clip().Role(ui.RoleNone)
	if len(carets) > 0 {
		e.Role(ui.RoleStatus).Label(fmt.Sprintf(core.Msg(c, "remote_cursor.label"), strings.Join(names, ", ")))
	}
	return e.Draw(func(p *ui.Painter, r ui.Rect) {
		for _, rc := range carets {
			col := remoteColor(k, rc.ID)
			if opts.Selection != nil && rc.Anchor != rc.Offset {
				for _, b := range opts.Selection(min(rc.Anchor, rc.Offset), max(rc.Anchor, rc.Offset)) {
					p.Fill(remoteRect(p, ui.Rect{X: r.X + b.X, Y: r.Y + b.Y, W: b.W, H: b.H}), col.Alpha(0.22), 2)
				}
			}
			cb := opts.Caret(rc.Offset)
			x, y := internal.Snap(p, r.X+cb.X), internal.Snap(p, r.Y+cb.Y)
			p.Fill(ui.Rect{X: x - 1, Y: y, W: 2, H: cb.H}, col, 0)
			// Keep the name flag inside the overlay.
			tw, th := p.MeasureText(0, ui.Span{Text: rc.Name, Size: size})
			flag := ui.Rect{X: x - 1, Y: y - th - 2, W: min(tw+8, r.W), H: th + 2}
			if flag.Y < r.Y {
				flag.Y = y + cb.H
			}
			flag.X = max(r.X, min(flag.X, r.X+r.W-flag.W))
			flag.Y = max(r.Y, min(flag.Y, r.Y+r.H-flag.H))
			flag = remoteRect(p, flag)
			p.Fill(flag, col, 2)
			p.Clip(flag, 2, func() { p.Text(flag.X+4, flag.Y+1, rc.Name, size, k.Surface) })
		}
	})
}

func remoteRect(p *ui.Painter, r ui.Rect) ui.Rect {
	x, y := internal.Snap(p, r.X), internal.Snap(p, r.Y)
	return ui.Rect{X: x, Y: y, W: internal.Snap(p, r.X+r.W) - x, H: internal.Snap(p, r.Y+r.H) - y}
}

// remoteColor assigns a palette color by collaborator ID.
func remoteColor(k theme.Tokens, id string) ui.Color {
	h := fnv.New32a()
	h.Write([]byte(id))
	pal := code.FlamegraphPalette(k)
	return pal[h.Sum32()%uint32(len(pal))]
}
