package mujica

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/internal"
	"mujicaui/theme"
)

// DragDropOverlayOptions configure a DragDropOverlay.
type DragDropOverlayOptions struct {
	// Hint replaces the built-in "Drop files to attach".
	Hint string
	// Accept checks each dropped path; nil accepts all.
	Accept func(path string) bool
}

// DropResult carries the overlay and dropped files.
type DropResult struct {
	Element  *ui.Element
	accepted []string
	rejected []string
}

// Accepted returns accepted paths; Rejected returns refused paths.
func (r DropResult) Accepted() []string { return r.accepted }
func (r DropResult) Rejected() []string { return r.rejected }

// DragDropOverlay covers content during external file dragging.
func DragDropOverlay(c *ui.Context, opts DragDropOverlayOptions, content func()) DropResult {
	if content == nil {
		panic("mujica: DragDropOverlay needs content")
	}
	k := Tokens(c)
	u := Density(c).Unit()
	var r DropResult
	r.Element = ui.Box(c).MinHeight(0).MinWidth(0)
	disabled := r.Element.IsDisabled()
	rejected := ui.Local(r.Element, "rejected", func() []string { return nil })
	over := !disabled && r.Element.FileDragOver()
	if files := r.Element.DroppedFiles(); len(files) > 0 && !disabled {
		r.accepted, r.rejected = dragDropSplit(files, opts.Accept)
		*rejected = r.rejected
		slog.Debug("mujica: drag_drop_overlay drop", "accepted", len(r.accepted), "rejected", len(r.rejected))
	}
	if over {
		*rejected = nil
	}
	hint := opts.Hint
	if hint == "" {
		hint = msg(c, "drag_drop_overlay.hint")
	}
	r.Element.Children(func() {
		content()
		if len(*rejected) > 0 && !over {
			names := make([]string, len(*rejected))
			for i, p := range *rejected {
				names[i] = filepath.Base(p)
			}
			ui.Text(c, fmt.Sprintf(msg(c, "drag_drop_overlay.rejected"), strings.Join(names, ", "))).
				TextColor(k.Danger).FontSize(FontSize(c, theme.CaptionSize)).Role(ui.RoleStatus).Padding(u, u*2)
		}
		if over {
			dragDropMask(c, hint)
		}
	})
	return r
}

// dragDropMask covers its parent with the drop hint.
func dragDropMask(c *ui.Context, hint string) {
	k := Tokens(c)
	u := Density(c).Unit()
	mask := ui.Column(c).Absolute().Left(0).Top(0).Right(0).Bottom(0).Center().Gap(u * 2).
		Background(k.Background.Alpha(0.92)).Radius(theme.CardRadius).PassThrough().Role(ui.RoleStatus).Label(hint)
	mask.Opacity(overlayEnter(c, mask))
	mask.Draw(func(p *ui.Painter, b ui.Rect) { internal.Frame(p, b, u*2, 5, k.Ornament) })
	mask.Children(func() {
		ui.Icon(c, icons.Must("upload")).Size(24, 24).TextColor(k.AccentText)
		ui.Text(c, hint).Font(theme.SerifFont).FontSize(FontSize(c, theme.H3Size)).TextColor(k.Text)
	})
}

// dragDropSplit sorts dropped paths by accept, nil accepting all.
func dragDropSplit(files []string, accept func(string) bool) (ok, no []string) {
	for _, f := range files {
		if accept == nil || accept(f) {
			ok = append(ok, f)
		} else {
			no = append(no, f)
		}
	}
	return ok, no
}
