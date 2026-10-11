package files

import (
	"fmt"
	"log/slog"
	"sort"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/code"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// StorageSegment is one category of used space.
type StorageSegment struct {
	Label string
	Bytes int64
}

// StorageUsageOptions configure a StorageUsage.
type StorageUsageOptions struct {
	// Capacity is the total space; required.
	Capacity int64
	// Label names the volume.
	Label string
	// Revision changes when segments are edited in place.
	Revision uint64
}

type storageUsageState struct {
	data           *StorageSegment
	n              int
	revision       uint64
	capacity, used int64
	ends           []int64
	rows           ui.ListState
}

// StorageUsage displays segmented capacity with keyboard readouts and details.
func StorageUsage(c *ui.Context, segments []StorageSegment, opts StorageUsageOptions) *ui.Element {
	if opts.Capacity <= 0 {
		panic("mujica: StorageUsage needs a positive Capacity")
	}
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	e := ui.Column(c).Gap(u * 3).Role(ui.RoleGroup).Label(opts.Label)
	s := ui.Local(e, "storage_usage.data", func() storageUsageState { return storageUsageState{} })
	var first *StorageSegment
	if len(segments) > 0 {
		first = &segments[0]
	}
	if s.data != first || s.n != len(segments) || s.capacity != opts.Capacity || s.revision != opts.Revision {
		s.used, s.ends = 0, s.ends[:0]
		for _, seg := range segments {
			if seg.Bytes < 0 || seg.Bytes > opts.Capacity-s.used {
				panic("mujica: StorageUsage segment size out of range")
			}
			s.used += seg.Bytes
			s.ends = append(s.ends, s.used)
		}
		s.data, s.n, s.capacity, s.revision = first, len(segments), opts.Capacity, opts.Revision
	}
	used := s.used
	pal := code.FlamegraphPalette(k)
	pct := func(b int64) float64 { return 100 * float64(b) / float64(opts.Capacity) }
	bytes := func(b int64) string { return code.ProcessListBytes(uint64(b)) }
	readout := func(i int) string {
		if i < 0 || i >= len(segments) {
			return fmt.Sprintf(core.Msg(c, "storage_usage.free"), bytes(opts.Capacity-used), pct(opts.Capacity-used))
		}
		return fmt.Sprintf("%s: %s (%.1f%%)", segments[i].Label, bytes(segments[i].Bytes), pct(segments[i].Bytes))
	}
	summary := fmt.Sprintf(core.Msg(c, "storage_usage.used"), bytes(used), bytes(opts.Capacity))
	focus := ui.Local(e, "storage_usage.focus", func() int { return 0 })
	*focus = min(*focus, len(segments))
	e.Children(func() {
		ui.Row(c).Gap(u * 2).Children(func() {
			if opts.Label != "" {
				ui.Text(c, opts.Label).Font(theme.SerifFont).FontWeight(600).Grow(1)
			}
			ui.Text(c, summary).TextColor(k.TextMuted)
		})
		bar := ui.Box(c).Height(14).Focusable().Role(ui.RoleMeter).Label(summary).Description(readout(*focus))
		n := len(segments) + 1
		moved := true
		switch {
		case bar.Shortcut(0, ui.KeyRight) && !bar.IsDisabled():
			*focus = (*focus + 1) % n
		case bar.Shortcut(0, ui.KeyLeft) && !bar.IsDisabled():
			*focus = (*focus + n - 1) % n
		default:
			moved = false
		}
		if moved {
			c.Announce(readout(*focus))
			slog.Debug("mujica: storage_usage focus", "segment", *focus)
		}
		internal.Ring(bar, k.Focus, theme.ControlRadius)
		focused, at := bar.Focused(), *focus
		bar.Draw(func(p *ui.Painter, r ui.Rect) {
			p.Fill(r, k.Border.Mix(k.Surface, 0.4), 3)
			if len(segments) > int(r.W) {
				for px := float32(0); px < r.W; px++ {
					b := int64(float64(px) / float64(r.W) * float64(opts.Capacity))
					i := sort.Search(len(s.ends), func(i int) bool { return s.ends[i] > b })
					if i == len(segments) {
						break
					}
					p.Fill(ui.Rect{X: internal.Snap(p, r.X+px), Y: r.Y, W: min(1, r.W-px), H: r.H}, pal[i%len(pal)], 0)
				}
				return
			}
			x := r.X
			for i, s := range segments {
				w := internal.Snap(p, r.W*float32(s.Bytes)/float32(opts.Capacity))
				if w <= 0 {
					continue
				}
				seg := ui.Rect{X: x, Y: r.Y, W: max(w-2, 1), H: r.H}
				p.Fill(seg, pal[i%len(pal)], 2)
				if focused && i == at {
					p.Stroke(ui.Rect{X: seg.X - 1, Y: seg.Y - 1, W: seg.W + 2, H: seg.H + 2}, k.Text, 3, 1.5)
				}
				x += w
			}
		})
		if focused {
			ui.Text(c, readout(at)).FontSize(core.FontSize(c, theme.CaptionSize)).Role(ui.RoleStatus)
		}
		ui.List(c, &s.rows, len(segments)+1, func(i int) {
			if i == len(segments) {
				storageLegend(c, k.Border.Mix(k.Surface, 0.4), core.Msg(c, "storage_usage.freeLabel"), bytes(opts.Capacity-used), pct(opts.Capacity-used))
				return
			}
			seg := segments[i]
			storageLegend(c, pal[i%len(pal)], seg.Label, bytes(seg.Bytes), pct(seg.Bytes))
		}).Height(min(float32(len(segments)+1)*core.Density(c).ControlHeight(), core.Density(c).ControlHeight()*6))
	})
	return e
}

func storageLegend(c *ui.Context, col ui.Color, label, size string, pct float64) {
	k := core.Tokens(c)
	ui.Row(c).MinHeight(core.Density(c).ControlHeight()).Gap(core.Density(c).Unit() * 2).AlignItems(ui.Center).Children(func() {
		ui.Box(c).Size(10, 10).Radius(2).Shrink(0).Background(col).Role(ui.RoleNone)
		ui.Text(c, label).Grow(1).MinWidth(0).SingleLine()
		ui.Text(c, size).TextColor(k.TextMuted).Shrink(0)
		ui.Text(c, fmt.Sprintf("%.1f%%", pct)).TextColor(k.TextMuted).FontFeatures("tnum").Width(52).TextAlign(ui.End).Shrink(0)
	})
}
