package code

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/data"
	"github.com/ZacharyZhang-NY/MujicaUI/input"
	"github.com/ZacharyZhang-NY/MujicaUI/internal"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
)

// Variable is a debugged value; Changed marks new values.
type Variable struct {
	Name, Value, Type string
	Changed           bool
	Children          []Variable
}

// VariablesPanelState is the open rows, choice and edit.
type VariablesPanelState struct {
	outline ui.OutlineState[string]
	cursor  int
	byPath  map[string]*Variable
	kids    map[string][]string
	roots   []string
	data    *Variable
	n       int
	editing string
	draft   string
	started bool
}

// VariablesPanelOptions configure a VariablesPanel.
type VariablesPanelOptions struct {
	Disabled bool
}

// VariableEdit is a typed value for a dotted Path.
type VariableEdit struct{ Path, Value string }

// VariablesPanelResult carries a committed edit.
type VariablesPanelResult struct {
	Element *ui.Element
	Edited  *VariableEdit
}

// VariablesPanel shows an editable variable tree.
func VariablesPanel(c *ui.Context, vars []Variable, s *VariablesPanelState, opts VariablesPanelOptions) VariablesPanelResult {
	if s == nil {
		panic("mujica: VariablesPanel needs a state")
	}
	if !s.started {
		s.started, s.cursor = true, -1
	}
	variablesIndex(vars, s)
	k := core.Tokens(c)
	d := core.Density(c)
	var res VariablesPanelResult
	res.Element = CodePanel(c).Disabled(opts.Disabled)
	disabled := res.Element.IsDisabled()
	if disabled {
		s.editing = ""
	}
	res.Element.Children(func() {
		if len(vars) == 0 {
			CodeCaption(c, core.Msg(c, "variables_panel.empty")).Padding(d.Unit() * 4)
			return
		}
		s.outline.List.Selected = &s.cursor
		s.outline.List.Label = func(i int) string { return s.byPath[s.outline.Item(i)].Name }
		children := func(p string) []string { return s.kids[p] }
		restore := data.DataSelectTheme(c, 0)
		e := ui.Outline(c, &s.outline, s.roots, children, func(p string) {
			v := s.byPath[p]
			ui.Row(c).MinHeight(d.ControlHeight() - 2*c.Theme().Space(1)).Gap(d.Unit() * 2).AlignItems(ui.Center).Grow(1).MinWidth(0).Children(func() {
				CodeMono(c, v.Name).TextColor(k.Info).Shrink(0)
				if s.editing == p {
					in := input.TextInput(c, &s.draft, input.TextInputOptions{Label: core.Msg(c, "variables_panel.value")})
					in.Element.Grow(1).MinWidth(0)
					in.Input.AutoFocus()
					switch {
					case in.Input.Shortcut(0, ui.KeyEscape):
						s.editing = ""
					case in.Submitted() && !disabled:
						res.Edited = &VariableEdit{Path: p, Value: s.draft}
						s.editing = ""
						slog.Debug("mujica: variables_panel edit", "path", p)
					}
					return
				}
				variablesValue(c, v)
			})
		})
		restore()
		e.Grow(1).MinHeight(d.ControlHeight() * 3)
		internal.Ring(e, k.Focus, theme.ControlRadius)
		if disabled || s.editing != "" || s.cursor < 0 || s.cursor >= s.outline.Rows() {
			return
		}
		p := s.outline.Item(s.cursor)
		if len(s.kids[p]) == 0 && s.byPath[p].Children == nil && (e.Submitted() || e.Shortcut(0, ui.KeyF2)) {
			s.editing, s.draft = p, s.byPath[p].Value
			slog.Debug("mujica: variables_panel editing", "path", p)
		}
	})
	return res
}

// variablesValue shows a value, marking changed ones.
func variablesValue(c *ui.Context, v *Variable) {
	k := core.Tokens(c)
	if v.Changed {
		ui.Box(c).Size(7, 7).Shrink(0).Role(ui.RoleNone).Draw(func(p *ui.Painter, r ui.Rect) {
			p.FillPath(internal.Diamond(r.X+r.W/2, r.Y+r.H/2, r.W), k.AccentText)
		})
	}
	val := CodeMono(c, v.Value).MinWidth(0)
	if v.Changed {
		val.FontWeight(700).TextColor(k.AccentText).Description(core.Msg(c, "variables_panel.changed"))
	}
	ui.Spacer(c)
	if v.Type != "" {
		CodeCaption(c, v.Type).Shrink(0)
	}
}

// variablesIndex keys variables by dotted path.
func variablesIndex(vars []Variable, s *VariablesPanelState) {
	var first *Variable
	if len(vars) > 0 {
		first = &vars[0]
	}
	if s.data == first && s.n == len(vars) && s.byPath != nil {
		return
	}
	s.data, s.n = first, len(vars)
	s.byPath, s.kids = map[string]*Variable{}, map[string][]string{}
	var walk func(list []Variable, prefix string) []string
	walk = func(list []Variable, prefix string) []string {
		paths := []string{}
		for i := range list {
			p := list[i].Name
			if prefix != "" {
				p = prefix + "." + p
			}
			if _, dup := s.byPath[p]; dup {
				panic("mujica: duplicate variable path " + p)
			}
			s.byPath[p] = &list[i]
			if list[i].Children != nil {
				s.kids[p] = walk(list[i].Children, p)
			}
			paths = append(paths, p)
		}
		return paths
	}
	s.roots = walk(vars, "")
}
