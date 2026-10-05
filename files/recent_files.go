package files

import (
	"cmp"
	"fmt"
	"log/slog"
	"path"
	"slices"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/code"
	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/icons"
	"github.com/ZacharyZhang-NY/MujicaUI/input"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
)

// RecentFile is a file the user opened.
type RecentFile struct {
	Path   string
	Opened time.Time
}

// RecentFilesOptions configure a RecentFiles.
type RecentFilesOptions struct {
	// Revision changes when files are edited in place.
	Revision uint64
	// Now dates the groups; zero uses the frame's time.
	Now      time.Time
	Disabled bool
}

// RecentFilesState is a RecentFiles' choice and place.
type RecentFilesState struct {
	rows     ui.ListState
	selected int
	flat     []recentRow
	data     *RecentFile
	n        int
	day      time.Time
	started  bool
	revision uint64
}

// recentRow holds a group header or file index.
type recentRow struct{ group, file int }

// RecentFilesResult is the list and what happened this frame.
type RecentFilesResult struct {
	Element *ui.Element
	Opened  *RecentFile
	// Removed identifies one file; Cleared requests clearing all files.
	Removed string
	Cleared bool
}

// RecentFiles groups recent files by time and supports clearing.
func RecentFiles(c *ui.Context, files []RecentFile, s *RecentFilesState, opts RecentFilesOptions) RecentFilesResult {
	if s == nil {
		panic("mujica: RecentFiles needs a state")
	}
	if !s.started {
		s.started, s.selected = true, -1
	}
	now := opts.Now
	if now.IsZero() {
		now = c.Now()
	}
	if s.revision != opts.Revision {
		s.day = time.Time{}
		s.revision = opts.Revision
	}
	recentGroup(files, s, now)
	k := core.Tokens(c)
	d := core.Density(c)
	var res RecentFilesResult
	res.Element = code.CodePanel(c).Disabled(opts.Disabled)
	disabled := res.Element.IsDisabled()
	if disabled {
		copy := *s
		s = &copy
	}
	res.Element.Children(func() {
		code.CodeHeader(c, func() {
			ui.Text(c, core.Msg(c, "recent_files.title")).Font(theme.SerifFont).FontWeight(600).Grow(1)
			if len(files) > 0 && input.Button(c, core.Msg(c, "recent_files.clear"), input.ButtonOptions{Variant: input.Ghost}).MinHeight(theme.MinHitSize).Clicked() && !disabled {
				res.Cleared = true
				slog.Debug("mujica: recent_files clear")
			}
		})
		if len(files) == 0 {
			code.CodeCaption(c, core.Msg(c, "recent_files.empty")).Padding(d.Unit() * 4)
			return
		}
		s.rows.Selected = &s.selected
		s.rows.Key = func(i int) any {
			r := s.flat[i]
			if r.file < 0 {
				return r.group
			}
			return files[r.file].Path
		}
		s.rows.Header = func(i int) bool { return s.flat[i].file < 0 }
		s.rows.Label = func(i int) string {
			if r := s.flat[i]; r.file >= 0 {
				return path.Base(files[r.file].Path)
			}
			return recentGroupName(c, s.flat[i].group)
		}
		list := code.CodeRows(c, &s.rows, len(s.flat), func(i int) {
			r := s.flat[i]
			if r.file < 0 {
				ui.Text(c, recentGroupName(c, r.group)).Font(theme.SerifFont).FontSize(core.FontSize(c, theme.CaptionSize)).
					TextColor(k.Ornament.Mix(k.Text, 0.2)).Grow(1)
				return
			}
			f := files[r.file]
			FileIcon(c, f.Path, FileIconOptions{})
			ui.Column(c).Grow(1).MinWidth(0).Children(func() {
				ui.Text(c, path.Base(f.Path)).SingleLine()
				code.CodeCaption(c, path.Dir(f.Path))
			})
			when := f.Opened.Format("15:04")
			if r.group >= 2 {
				when = f.Opened.Format("Jan 2")
			}
			code.CodeCaption(c, when).Shrink(0)
			if core.IconAction(c, icons.Must("x"), fmt.Sprintf(core.Msg(c, "recent_files.remove"), path.Base(f.Path))).Clicked() && !disabled {
				res.Removed = f.Path
				slog.Debug("mujica: recent_files remove")
			}
		})
		list.Background(k.Surface)
		if list.Submitted() && !disabled && s.selected >= 0 && s.selected < len(s.flat) && s.flat[s.selected].file >= 0 {
			f := files[s.flat[s.selected].file]
			res.Opened = &f
			slog.Debug("mujica: recent_files open")
		}
	})
	return res
}

func recentGroupName(c *ui.Context, g int) string {
	return core.Msg(c, [...]string{"recent_files.today", "recent_files.yesterday", "recent_files.week", "recent_files.earlier"}[g])
}

// recentGroup caches files sorted into relative-day groups.
func recentGroup(files []RecentFile, s *RecentFilesState, now time.Time) {
	var first *RecentFile
	if len(files) > 0 {
		first = &files[0]
	}
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	if s.data == first && s.n == len(files) && s.day.Equal(today) {
		return
	}
	s.data, s.n, s.day = first, len(files), today
	idx := make([]int, len(files))
	for i := range idx {
		idx[i] = i
	}
	slices.SortStableFunc(idx, func(a, b int) int { return cmp.Compare(files[b].Opened.UnixNano(), files[a].Opened.UnixNano()) })
	group := func(t time.Time) int {
		switch {
		case !t.Before(today):
			return 0
		case !t.Before(today.AddDate(0, 0, -1)):
			return 1
		case !t.Before(today.AddDate(0, 0, -6)):
			return 2
		}
		return 3
	}
	s.flat = s.flat[:0]
	last := -1
	for _, i := range idx {
		g := group(files[i].Opened)
		if g != last {
			s.flat = append(s.flat, recentRow{g, -1})
			last = g
		}
		s.flat = append(s.flat, recentRow{g, i})
	}
}
