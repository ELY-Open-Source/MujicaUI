package media

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/icons"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// StripFrame is one thumbnail of a VideoThumbnailStrip.
type StripFrame struct {
	// Image is a bitmap or SVG; nil leaves it empty.
	Image ui.ImageSource
	At    time.Duration
}

// VideoThumbnailStripOptions configure a VideoThumbnailStrip.
type VideoThumbnailStripOptions struct {
	// FrameWidth is each cell's width; zero means 96.
	FrameWidth float32
	Disabled   bool
}

const stripGap = 4

type stripKey struct{}

// VideoThumbnailStrip shows the frames fitting around *current.
func VideoThumbnailStrip(c *ui.Context, current *int, frames []StripFrame, opts VideoThumbnailStripOptions) core.ChoiceResult {
	if len(frames) == 0 || *current < 0 || *current >= len(frames) {
		panic("mujica: VideoThumbnailStrip current out of range")
	}
	if opts.FrameWidth < 0 {
		panic("mujica: VideoThumbnailStrip FrameWidth must not be negative")
	}
	fw := opts.FrameWidth
	if fw == 0 {
		fw = 96
	}
	fh := fw * 9 / 16
	k := core.Tokens(c)
	old := *current
	row := ui.Row(c).Gap(core.Density(c).Unit()).AlignItems(ui.Center).Disabled(opts.Disabled)
	row.Children(func() {
		prev := core.IconAction(c, icons.Must("chevron-left"), core.Msg(c, "video_thumbnail_strip.previous")).Disabled(*current == 0)
		view := ui.Row(c).Grow(1).MinWidth(fw).Height(fh + 4).Gap(stripGap).Clip().Focusable().
			Role(ui.RoleGroup).Label(core.Msg(c, "video_thumbnail_strip.label")).
			Description(fmt.Sprintf(core.Msg(c, "video_thumbnail_strip.position"), *current+1, len(frames), core.MediaClock(frames[*current].At)))
		internal.Ring(view, k.Focus, theme.ControlRadius)
		next := core.IconAction(c, icons.Must("chevron-right"), core.Msg(c, "video_thumbnail_strip.next")).Disabled(*current == len(frames)-1)
		live := !row.IsDisabled()
		switch {
		case live && (prev.Clicked() || view.Shortcut(0, ui.KeyLeft)):
			*current = max(0, *current-1)
		case live && (next.Clicked() || view.Shortcut(0, ui.KeyRight)):
			*current = min(len(frames)-1, *current+1)
		case live && view.Shortcut(0, ui.KeyHome):
			*current = 0
		case live && view.Shortcut(0, ui.KeyEnd):
			*current = len(frames) - 1
		}
		// Only frames that fit are built, around *current.
		w := view.Bounds().W
		if w <= 0 {
			// The count needs the width: lay out, then rebuild.
			c.AnimationFrame()
		}
		n := max(1, int((w+stripGap)/(fw+stripGap)))
		start := ui.Local(view, stripKey{}, func() int { return 0 })
		*start = max(0, min(*start, *current, len(frames)-n))
		if *current >= *start+n {
			*start = *current - n + 1
		}
		view.Children(func() {
			for i := *start; i < min(len(frames), *start+n); i++ {
				if stripCell(c, k, frames[i], i, i == *current, fw, fh) && live {
					*current = i
				}
			}
		})
	})
	changed := *current != old
	if changed {
		slog.Debug("mujica: video_strip jump", "index", *current)
	}
	return core.NewChoiceResult(row, changed)
}

// stripCell draws one frame and reports a click on it.
func stripCell(c *ui.Context, k theme.Tokens, f StripFrame, i int, on bool, fw, fh float32) bool {
	cell := ui.Box(c.Key(i)).Size(fw, fh).Shrink(0).Radius(theme.ControlRadius).Clip().Background(VideoBlack).
		Cursor(ui.CursorPointer).AlignSelf(ui.Center)
	edge, width := k.Border, theme.BorderWidth
	if on {
		edge, width = k.AccentText, 2
	}
	cell.Border(width, edge)
	cell.Children(func() {
		if f.Image != nil {
			ui.Image(c, f.Image).Fit(ui.Cover).Size(fw, fh)
		}
		ui.Text(c, core.MediaClock(f.At)).Absolute().Left(4).Bottom(3).FontSize(core.FontSize(c, theme.CaptionSize)).
			TextColor(ui.Hex("#F2EBDD")).Background(ui.RGBA(0, 0, 0, 0.6)).FontFeatures("tnum")
		if on {
			ui.Box(c).Absolute().Top(3).Right(3).Size(10, 10).Draw(func(p *ui.Painter, r ui.Rect) {
				p.FillPath(internal.Diamond(internal.Snap(p, r.X+5), internal.Snap(p, r.Y+5), 10), k.AccentText)
			})
		}
	})
	return cell.Clicked()
}
