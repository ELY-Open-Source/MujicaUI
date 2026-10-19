package git

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/data"
	"github.com/ELY-Open-Source/MujicaUI/internal"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

// CommitGraphState is a CommitGraph's selection and its lane layout.
type CommitGraphState struct {
	List data.ListState[string]
	rows []commitGraphRow
	// lanes is the widest the graph gets.
	lanes int
	src   *Commit
	n     int
}

// CommitGraphOptions configure a CommitGraph.
type CommitGraphOptions struct {
	// TimeFormat lays out commit times; empty means "2006-01-02 15:04".
	TimeFormat string
	// MaxLanes caps the graph width in lanes; zero means 8.
	MaxLanes int
	Status   data.DataStatus
	Error    string
}

// commitGraphLink joins lane from to lane to within a row.
type commitGraphLink struct{ from, to, color int }

type commitGraphRow struct {
	node, color int
	merge       bool
	// top links enter from above; bottom links leave the node.
	top, bottom []commitGraphLink
}

type commitGraphLane struct {
	hash  string
	color int
}

// CommitGraph shows commits with branch lanes; give it a height.
func CommitGraph(c *ui.Context, s *CommitGraphState, commits []Commit, opts CommitGraphOptions) GitListResult {
	maxLanes := opts.MaxLanes
	if maxLanes == 0 {
		maxLanes = 8
	}
	if maxLanes < 0 {
		panic("mujica: CommitGraph MaxLanes is negative")
	}
	var first *Commit
	if len(commits) > 0 {
		first = &commits[0]
	}
	if first != s.src || len(commits) != s.n {
		s.rows, s.lanes = commitGraphLayout(commits)
		s.src, s.n = first, len(commits)
		slog.Debug("mujica: commit graph layout", "commits", len(commits), "lanes", s.lanes)
	}
	layout := gitTimeLayout(opts.TimeFormat)
	lane := core.Density(c).Unit() * 4
	width := float32(min(max(s.lanes, 1), maxLanes)) * lane
	h := core.Density(c).ControlHeight()
	v := data.List(c, &s.List, len(commits), data.ListOptions[string]{
		Key:    func(i int) string { return gitCommitHash(commits[i]) },
		Label:  func(i int) string { return commits[i].Subject },
		Status: opts.Status, Error: opts.Error, Empty: core.Msg(c, "commit_list.empty"),
	}, func(i int) {
		row := s.rows[i]
		ui.Box(c).Size(width, h).Shrink(0).Clip().Role(ui.RoleNone).Draw(func(p *ui.Painter, r ui.Rect) {
			commitGraphDraw(p, core.Tokens(c), r, row, lane)
		})
		commitListRow(c, commits[i], layout)
	})
	res := GitListResult{Element: v.Element, changed: v.Changed(), submitted: v.Submitted()}
	if res.changed || res.submitted {
		key, _ := s.List.Selected()
		slog.Debug("mujica: commit graph choose", "hash", key, "open", res.submitted)
	}
	return res
}

// commitGraphLayout gives each commit a lane and parent links.
func commitGraphLayout(commits []Commit) ([]commitGraphRow, int) {
	rows := make([]commitGraphRow, len(commits))
	var lanes []commitGraphLane
	colors, widest := 0, 0
	free := func() int {
		for j, l := range lanes {
			if l.hash == "" {
				return j
			}
		}
		lanes = append(lanes, commitGraphLane{})
		return len(lanes) - 1
	}
	for i, cm := range commits {
		h := gitCommitHash(cm)
		var row commitGraphRow
		row.node = -1
		for j, l := range lanes {
			if l.hash == h && row.node < 0 {
				row.node = j
			}
		}
		for j, l := range lanes {
			switch {
			case l.hash == "":
			case l.hash == h:
				row.top = append(row.top, commitGraphLink{j, row.node, l.color})
			default:
				row.top = append(row.top, commitGraphLink{j, j, l.color})
			}
		}
		if row.node < 0 {
			row.node = free()
			lanes[row.node] = commitGraphLane{h, colors % 6}
			colors++
		}
		row.color = lanes[row.node].color
		for j := range lanes {
			if lanes[j].hash == h {
				lanes[j].hash = ""
			}
		}
		row.merge = len(cm.Parents) > 1
		for pi, p := range cm.Parents {
			if pi == 0 {
				lanes[row.node] = commitGraphLane{p, row.color}
				row.bottom = append(row.bottom, commitGraphLink{row.node, row.node, row.color})
				continue
			}
			to := -1
			for j, l := range lanes {
				if l.hash == p {
					to = j
					break
				}
			}
			if to < 0 {
				to = free()
				lanes[to] = commitGraphLane{p, colors % 6}
				colors++
			}
			row.bottom = append(row.bottom, commitGraphLink{row.node, to, lanes[to].color})
		}
		for len(lanes) > 0 && lanes[len(lanes)-1].hash == "" {
			lanes = lanes[:len(lanes)-1]
		}
		widest = max(widest, len(lanes), row.node+1)
		rows[i] = row
	}
	return rows, widest
}

// commitGraphColors are the lane colors; place also tells lanes apart.
func commitGraphColors(k theme.Tokens) [6]ui.Color {
	return [6]ui.Color{k.AccentText, k.Info, k.Success, k.Warning, k.Danger, k.TextMuted}
}

func commitGraphDraw(p *ui.Painter, k theme.Tokens, r ui.Rect, row commitGraphRow, lane float32) {
	cols := commitGraphColors(k)
	x := func(j int) float32 { return internal.Snap(p, r.X+lane/2+float32(j)*lane) }
	top, mid, bot := internal.Snap(p, r.Y), internal.Snap(p, r.Y+r.H/2), internal.Snap(p, r.Y+r.H)
	seg := func(l commitGraphLink, y0, y1 float32) {
		var path ui.Path
		xa, xb := x(l.from), x(l.to)
		path.MoveTo(xa, y0)
		if xa == xb {
			path.LineTo(xb, y1)
		} else {
			m := (y0 + y1) / 2
			path.CubeTo(xa, m, xb, m, xb, y1)
		}
		p.StrokePath(&path, 2, cols[l.color])
	}
	for _, l := range row.top {
		if l.from == l.to && l.to != row.node {
			seg(l, top, bot)
			continue
		}
		seg(l, top, mid)
	}
	for _, l := range row.bottom {
		seg(l, mid, bot)
	}
	cx, col := x(row.node), cols[row.color]
	if row.merge {
		p.FillPath(internal.Diamond(cx, mid, 11), k.Surface)
		p.StrokePath(internal.Diamond(cx, mid, 9), 2, col)
		return
	}
	var dot ui.Path
	dot.Circle(cx, mid, 5.5)
	p.FillPath(&dot, k.Surface)
	var ring ui.Path
	ring.Circle(cx, mid, 4)
	p.FillPath(&ring, col)
}
