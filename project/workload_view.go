package project

import (
	"fmt"
	"math"
	"strconv"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/display"
	"github.com/ELY-Open-Source/MujicaUI/icons"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

// Workload is one member's assigned and available hours.
type Workload struct {
	Member   ProjectMember
	Assigned float64
	Capacity float64
}

// WorkloadViewOptions configure a WorkloadView.
type WorkloadViewOptions struct {
	// Unit labels the amounts; empty uses hours.
	Unit  string
	Label string
	// Version changes when loads change in place.
	Version int
}

// WorkloadView virtualizes member loads with capacity and overload marks.
func WorkloadView(c *ui.Context, loads []Workload, opts WorkloadViewOptions) ui.Element {
	k := core.Tokens(c)
	d := core.Density(c)
	unit := opts.Unit
	if unit == "" {
		unit = core.Msg(c, "workload_view.hours")
	}
	label := opts.Label
	if label == "" {
		label = core.Msg(c, "workload_view.label")
	}
	root := ui.Column(c).Grow(1).MinHeight(0).Role(ui.RoleList).Label(label)
	cache := ui.Local(root, "workload-summary", func() workloadCache { return workloadCache{} })
	var first *Workload
	if len(loads) > 0 {
		first = &loads[0]
	}
	if !cache.ready || cache.first != first || cache.n != len(loads) || cache.version != opts.Version {
		*cache = workloadSummary(loads)
		cache.first, cache.n, cache.version, cache.ready = first, len(loads), opts.Version, true
	}
	top := cache.top
	state := ui.Local(root, "workload-list", func() *ui.ListState { return &ui.ListState{} })
	st := *state
	st.Key = func(i int) any { return loads[i].Member.ID }
	num := func(v float64) string { return strconv.FormatFloat(v, 'f', -1, 64) }
	root.Children(func() {
		over := cache.over
		if over > 0 {
			ui.Row(c).Gap(d.Unit()).AlignItems(ui.Center).TextColor(k.Danger).Role(ui.RoleStatus).Padding(d.Unit()).Children(func() {
				ui.Icon(c, icons.Must("triangle-alert")).Size(14, 14)
				ui.Text(c, core.CalendarCount(c, "workload_view.overloaded", over)).FontSize(core.FontSize(c, theme.CaptionSize))
			})
		}
		ui.List(c, st, len(loads), func(i int) {
			l := loads[i]
			hot := l.Assigned > l.Capacity
			text := fmt.Sprintf("%s / %s %s", num(l.Assigned), num(l.Capacity), unit)
			note := ""
			if hot {
				note = fmt.Sprintf(core.Msg(c, "workload_view.over"), num(l.Assigned-l.Capacity), unit)
			}
			ui.Row(c).MinHeight(d.ControlHeight()).Gap(d.Unit() * 3).AlignItems(ui.Center).PaddingX(d.Unit()).Role(ui.RoleListItem).
				Label(l.Member.Name + ", " + text + " " + note).Children(func() {
				display.Avatar(c, l.Member.Name, display.AvatarOptions{Size: 24})
				ui.Text(c, l.Member.Name).Width(d.ControlHeight() * 4).Shrink(0).SingleLine().Ellipsis("…").TextColor(k.Text)
				ui.Box(c).Grow(1).Height(10).Draw(func(p *ui.Painter, r ui.Rect) { workloadBar(p, r, l, top, k.Border, k.AccentText, k.Danger, k.Text) })
				fg := k.TextMuted
				if hot {
					fg = k.Danger
					ui.Icon(c, icons.Must("triangle-alert")).Size(14, 14).TextColor(fg)
				}
				ui.Text(c, text).Width(d.ControlHeight() * 3).Shrink(0).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(fg).FontFeatures("tnum").TextAlign(ui.End)
			})
		}).Grow(1).MinHeight(d.ControlHeight() * 3)
	})
	return root
}

// workloadBar draws assigned work, overload hatching and capacity marks.
func workloadBar(p *ui.Painter, r ui.Rect, l Workload, top float64, track, fill, hot, tick ui.Color) {
	x := func(v float64) float32 { return r.X + r.W*float32(v/top) }
	p.Fill(r, track, 3)
	p.Fill(ui.Rect{X: r.X, Y: r.Y, W: x(min(l.Assigned, l.Capacity)) - r.X, H: r.H}, fill, 3)
	if l.Assigned > l.Capacity {
		p.Fill(ui.Rect{X: x(l.Capacity), Y: r.Y, W: x(l.Assigned) - x(l.Capacity), H: r.H}, hot, 0)
		for sx := x(l.Capacity) + 2; sx < x(l.Assigned); sx += 5 {
			p.Line(sx, r.Y+r.H, sx+3, r.Y, 1, track)
		}
	}
	cx := x(l.Capacity)
	p.Line(cx, r.Y-3, cx, r.Y+r.H+3, 2, tick)
}

type workloadCache struct {
	first            *Workload
	n, version, over int
	top              float64
	ready            bool
}

func workloadSummary(loads []Workload) workloadCache {
	out := workloadCache{}
	seen := map[string]bool{}
	for _, l := range loads {
		if l.Member.ID == "" || seen[l.Member.ID] || l.Assigned < 0 || l.Capacity <= 0 || math.IsNaN(l.Assigned) || math.IsNaN(l.Capacity) || math.IsInf(l.Assigned, 0) || math.IsInf(l.Capacity, 0) {
			panic("mujica: WorkloadView needs unique member IDs and finite non-negative loads with positive capacity")
		}
		seen[l.Member.ID] = true
		out.top = max(out.top, l.Assigned, l.Capacity)
		if l.Assigned > l.Capacity {
			out.over++
		}
	}
	return out
}
