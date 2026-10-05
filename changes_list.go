package mujica

import (
	"fmt"
	"log/slog"
	"path"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/theme"
)

// ChangedFile is a changed file in the tree or index.
type ChangedFile struct {
	Path   string
	Status GitStatus
	// Staged puts the file in the index group.
	Staged bool
}

// ChangesListState is a ChangesList's selection and grouping.
type ChangesListState struct {
	list   ui.ListState
	cursor int
	// rows maps rows to files; negatives are group headers.
	rows    []int
	staged  int
	src     *ChangedFile
	n       int
	started bool
}

// Selected returns the chosen file's index, or -1.
func (s *ChangesListState) Selected() int {
	if s.cursor < 0 || s.cursor >= len(s.rows) {
		return -1
	}
	return s.rows[s.cursor]
}

// ChangesListOptions configure a ChangesList.
type ChangesListOptions struct {
	Disabled bool
}

// ChangesList groups staged and unstaged files; give it a height.
func ChangesList(c *ui.Context, s *ChangesListState, files []ChangedFile, opts ChangesListOptions) GitListResult {
	if !s.started {
		s.started, s.cursor = true, -1
	}
	var first *ChangedFile
	if len(files) > 0 {
		first = &files[0]
	}
	if first != s.src || len(files) != s.n {
		s.rows, s.staged = changesListRows(files)
		s.src, s.n, s.cursor = first, len(files), -1
		slog.Debug("mujica: changes list regroup", "files", len(files), "staged", s.staged)
	}
	k := Tokens(c)
	var r GitListResult
	s.list.Selected = &s.cursor
	s.list.Header = func(i int) bool { return s.rows[i] < 0 }
	s.list.Key = func(i int) any {
		if f := s.rows[i]; f >= 0 {
			return fmt.Sprint(files[f].Staged, files[f].Path)
		}
		return s.rows[i]
	}
	s.list.Label = func(i int) string {
		if f := s.rows[i]; f >= 0 {
			return files[f].Path
		}
		return ""
	}
	prev := s.cursor
	outer := ui.Column(c).AlignItems(ui.Stretch).Disabled(opts.Disabled)
	var e *ui.Element
	outer.Children(func() {
		restore := dataSelectTheme(c, 0)
		defer restore()
		e = ui.List(c, &s.list, len(s.rows), func(i int) { changesListRow(c, s, files, i, &r) })
		e.Grow(1).Background(k.Surface).Border(theme.BorderWidth, k.Border).Radius(theme.ControlRadius)
	})
	r.Element, r.changed, r.submitted = outer, e.Changed(), e.Submitted()
	if outer.IsDisabled() {
		s.cursor, r.changed, r.submitted, r.action, r.item = prev, false, false, "", ""
	}
	if r.action != "" {
		slog.Debug("mujica: changes list action", "action", r.action, "path", r.item)
	}
	return r
}

// changesListRow builds row i: a header or a file.
func changesListRow(c *ui.Context, s *ChangesListState, files []ChangedFile, i int, r *GitListResult) {
	k := Tokens(c)
	d := Density(c)
	u := d.Unit()
	f := s.rows[i]
	row := ui.Row(c).AlignItems(ui.Center).Gap(u*2).MinHeight(d.ControlHeight()).Padding(0, u*2, 0, u*3)
	if f < 0 {
		key, all, icon, count := "changes_list.staged", "unstage_all", "minus", s.staged
		if f == -2 {
			key, all, icon, count = "changes_list.unstaged", "stage_all", "plus", len(files)-s.staged
		}
		row.Background(k.Background).Children(func() {
			ui.Text(c, fmt.Sprintf(msg(c, key), count)).Font(theme.SerifFont).FontWeight(600).
				FontSize(FontSize(c, theme.CaptionSize)).TextColor(k.Ornament.Mix(k.Text, 0.35)).Grow(1)
			b := iconAction(c, icons.Must(icon), msg(c, "changes_list."+all)).Disabled(count == 0)
			if b.Tooltip(msg(c, "changes_list."+all)).Clicked() && !b.IsDisabled() {
				r.action = all
			}
		})
		return
	}
	cf := files[f]
	if s.cursor == i {
		row.Draw(func(p *ui.Painter, rc ui.Rect) {
			p.Fill(ui.Rect{X: rc.X, Y: rc.Y + 4, W: 3, H: rc.H - 8}, k.AccentText, 1.5)
		})
	}
	row.Children(func() {
		GitStatusBadge(c, cf.Status, GitStatusBadgeOptions{})
		dir := path.Dir(cf.Path)
		ui.Text(c, path.Base(cf.Path)).SingleLine().Shrink(0).MaxWidth(240)
		if dir != "." {
			ui.Text(c, dir).FontSize(FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted).SingleLine().MinWidth(0)
		}
		ui.Spacer(c)
		id, icon := "stage", "plus"
		if cf.Staged {
			id, icon = "unstage", "minus"
		}
		b := iconAction(c, icons.Must(icon), msg(c, "changes_list."+id)+" "+cf.Path).Tooltip(msg(c, "changes_list."+id))
		if b.Clicked() && !b.IsDisabled() {
			r.action, r.item = id, cf.Path
		}
	})
}

// changesListRows orders staged files, then unstaged, each under a header.
func changesListRows(files []ChangedFile) ([]int, int) {
	rows := []int{-1}
	for i, f := range files {
		if f.Path == "" {
			panic("mujica: ChangedFile needs a Path")
		}
		gitStatusIndex(f.Status)
		if f.Staged {
			rows = append(rows, i)
		}
	}
	staged := len(rows) - 1
	rows = append(rows, -2)
	for i, f := range files {
		if !f.Staged {
			rows = append(rows, i)
		}
	}
	return rows, staged
}
