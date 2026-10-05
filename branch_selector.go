package mujica

import (
	"log/slog"
	"strings"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/internal"
	"mujicaui/theme"
)

// BranchSelectorOptions configure a BranchSelector.
type BranchSelectorOptions struct {
	// AllowCreate offers to create a branch named by the search.
	AllowCreate bool
	Disabled    bool
}

// BranchSelectorResult is a built BranchSelector and its events.
type BranchSelectorResult struct {
	Element *ui.Element
	changed bool
	create  string
}

// Changed reports a new branch in *selected this frame.
func (r BranchSelectorResult) Changed() bool { return r.changed }

// Created returns a branch name to create this frame.
func (r BranchSelectorResult) Created() (string, bool) { return r.create, r.create != "" }

type branchSelectorState struct {
	open, wasOpen bool
	query         string
	hl            int
	reveal        bool
	// found caches the filter of src for query fq.
	found []Branch
	fq    string
	src   *Branch
	n     int
}

// BranchSelector chooses a branch from a searchable popup.
func BranchSelector(c *ui.Context, selected *string, branches []Branch, opts BranchSelectorOptions) BranchSelectorResult {
	for _, b := range branches {
		branchKey(b)
	}
	k := Tokens(c)
	d := Density(c)
	u := d.Unit()
	tr := ui.ButtonBase(c).MinHeight(d.ControlHeight()).Padding(0, u*3).Gap(u*2).Radius(theme.ControlRadius).
		Background(k.Surface).Border(theme.BorderWidth, k.ControlBorder).Disabled(opts.Disabled).
		Label(msg(c, "branch_selector.label") + ": " + *selected)
	internal.Ring(tr, k.Focus, theme.ControlRadius)
	if tr.Hovered() {
		tr.Background(k.SurfaceHover)
	}
	tr.Transition(colorFade)
	st := ui.Local(tr, "state", func() *branchSelectorState { return &branchSelectorState{} })
	s := *st
	res := BranchSelectorResult{Element: tr}
	if tr.Clicked() && !tr.IsDisabled() {
		s.open = !s.open
	}
	if tr.IsDisabled() {
		s.open = false
	}
	if s.open && !s.wasOpen {
		s.query, s.hl = "", 0
		slog.Debug("mujica: branch selector open")
	}
	s.wasOpen = s.open
	tr.Children(func() {
		ui.Icon(c, icons.Must("git-branch")).Size(16, 16).TextColor(k.TextMuted)
		ui.Text(c, *selected).SingleLine().MaxWidth(220)
		ui.Icon(c, icons.Must("chevrons-up-down")).Size(14, 14).TextColor(k.TextMuted)
	})
	b := tr.Bounds()
	ui.PopoverBase(c, tr, &s.open, func(panel *ui.Element) {
		overlayStyle(c, panel, overlayAnchored(c, panel, b))
		panel.Margin(6, 0, 0, 0).Padding(u * 2).Gap(u * 2).Width(300).MaxHeight(popoverRoom(c, b)).AlignItems(ui.Stretch)
		var first *Branch
		if len(branches) > 0 {
			first = &branches[0]
		}
		if s.found == nil || s.fq != s.query || s.src != first || s.n != len(branches) {
			s.found, s.fq, s.src, s.n = branchSelectorFilter(branches, s.query), s.query, first, len(branches)
		}
		found := s.found
		create := opts.AllowCreate && strings.TrimSpace(s.query) != "" && !branchSelectorHas(branches, strings.TrimSpace(s.query))
		n := len(found)
		if create {
			n++
		}
		in := TextInput(c, &s.query, TextInputOptions{Placeholder: msg(c, "branch_selector.search"), Label: msg(c, "branch_selector.search"),
			Leading: func() { ui.Icon(c, icons.Must("search")).Size(14, 14).TextColor(k.TextMuted) }})
		in.Input.AutoFocus()
		if in.Changed() {
			s.hl, s.reveal = 0, true
		}
		switch {
		case n == 0:
		case in.Input.Shortcut(0, ui.KeyDown):
			s.hl, s.reveal = (s.hl+1)%n, true
		case in.Input.Shortcut(0, ui.KeyUp):
			s.hl, s.reveal = (s.hl-1+n)%n, true
		}
		s.hl = max(0, min(s.hl, n-1))
		pick := -1
		if in.Submitted() && n > 0 {
			pick = s.hl
		}
		ui.Scroll(c).Shrink(1).Gap(2).Role(ui.RoleList).Children(func() {
			if n == 0 {
				ui.Text(c, msg(c, "branch_selector.empty")).TextColor(k.TextMuted).Padding(u * 3)
			}
			group := -1
			for i, br := range found {
				if g := branchSelectorGroup(br); g != group {
					group = g
					ui.Text(c, msg(c, []string{"branch_selector.local", "branch_selector.remote"}[g])).Font(theme.SerifFont).
						FontSize(FontSize(c, theme.CaptionSize)).TextColor(k.Ornament.Mix(k.Text, 0.35)).Padding(u*2, u*2, u, u*2)
				}
				if branchSelectorRow(c, s, i, br.Name, br.Name == *selected, "git-branch") {
					pick = i
				}
			}
			if create && branchSelectorRow(c, s, len(found), msg(c, "branch_selector.create")+" “"+strings.TrimSpace(s.query)+"”", false, "plus") {
				pick = len(found)
			}
		})
		switch {
		case pick < 0:
		case pick == len(found):
			res.create = strings.TrimSpace(s.query)
			s.open = false
			slog.Debug("mujica: branch selector create", "name", res.create)
		default:
			s.open = false
			if found[pick].Name != *selected {
				*selected, res.changed = found[pick].Name, true
				slog.Debug("mujica: branch selector choose", "branch", *selected)
			}
		}
	})
	return res
}

// branchSelectorRow builds row i and reports a click on it.
func branchSelectorRow(c *ui.Context, s *branchSelectorState, i int, label string, current bool, icon string) bool {
	k := Tokens(c)
	d := Density(c)
	u := d.Unit()
	row := ui.Row(c).AlignItems(ui.Center).Gap(u*2).MinHeight(d.ControlHeight()).Padding(0, u*2).Radius(theme.ControlRadius).
		Role(ui.RoleListItem).Label(label).Cursor(ui.CursorPointer)
	if i == s.hl {
		row.Background(k.Selection)
		if s.reveal {
			row.ScrollIntoView()
			s.reveal = false
		}
	} else if row.Hovered() {
		row.Background(k.SurfaceHover)
	}
	row.Transition(colorFade)
	row.Children(func() {
		ui.Icon(c, icons.Must(icon)).Size(14, 14).TextColor(k.TextMuted)
		t := ui.Text(c, label).SingleLine().Grow(1).MinWidth(0)
		if current {
			t.FontWeight(600)
			ui.Icon(c, icons.Must("check")).Size(14, 14).TextColor(k.AccentText)
		}
	})
	return row.Clicked()
}

func branchSelectorGroup(b Branch) int {
	if b.Remote {
		return 1
	}
	return 0
}

// branchSelectorFilter keeps branches whose name holds the query, locals first.
func branchSelectorFilter(branches []Branch, query string) []Branch {
	q := strings.ToLower(strings.TrimSpace(query))
	var local, remote []Branch
	for _, b := range branches {
		if !strings.Contains(strings.ToLower(b.Name), q) {
			continue
		}
		if b.Remote {
			remote = append(remote, b)
		} else {
			local = append(local, b)
		}
	}
	return append(append(make([]Branch, 0, len(local)+len(remote)), local...), remote...)
}

func branchSelectorHas(branches []Branch, name string) bool {
	for _, b := range branches {
		if b.Name == name {
			return true
		}
	}
	return false
}
