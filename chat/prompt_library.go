package chat

import (
	"log/slog"
	"regexp"
	"slices"
	"strings"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/icons"
	"github.com/ZacharyZhang-NY/MujicaUI/input"
	"github.com/ZacharyZhang-NY/MujicaUI/internal"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
)

// PromptTemplate is a reusable prompt; {{name}} marks a variable.
type PromptTemplate struct {
	ID, Category, Title, Body string
}

// PromptLibraryOptions configure a PromptLibrary.
type PromptLibraryOptions struct {
	Label    string
	Disabled bool
}

// PromptLibraryState holds category, query, chosen template and variable values.
type PromptLibraryState struct {
	// Category filters templates; empty shows all.
	Category string
	Query    string
	Selected string
	Values   map[string]string
	shown    promptLibraryShown
}

// promptLibraryShown caches filtered templates until data, query or category change.
type promptLibraryShown struct {
	src             []PromptTemplate
	query, category string
	ids             []int
	cats            []string
	built           bool
}

// PromptLibraryView is a built PromptLibrary and inserted text.
type PromptLibraryView struct {
	Element *ui.Element
	// Inserted is the filled template, or empty.
	Inserted string
}

var promptLibraryVar = regexp.MustCompile(`\{\{\s*([A-Za-z0-9_.-]+)\s*\}\}`)

// PromptLibraryVariables lists a body's {{name}} variables in order, once each.
func PromptLibraryVariables(body string) []string {
	var out []string
	for _, m := range promptLibraryVar.FindAllStringSubmatch(body, -1) {
		if !slices.Contains(out, m[1]) {
			out = append(out, m[1])
		}
	}
	return out
}

// PromptLibraryFill replaces each {{name}} in body with values[name].
func PromptLibraryFill(body string, values map[string]string) string {
	return promptLibraryVar.ReplaceAllStringFunc(body, func(m string) string {
		return values[promptLibraryVar.FindStringSubmatch(m)[1]]
	})
}

// PromptLibrary browses, previews, fills and inserts prompt templates.
func PromptLibrary(c *ui.Context, s *PromptLibraryState, templates []PromptTemplate, opts PromptLibraryOptions) PromptLibraryView {
	if s == nil {
		panic("mujica: PromptLibrary needs a state")
	}
	if s.Values == nil {
		s.Values = map[string]string{}
	}
	k := core.Tokens(c)
	d := core.Density(c)
	promptLibraryFilter(&s.shown, templates, s.Query, s.Category)
	row := ui.Row(c).Gap(d.Unit() * 4).AlignItems(ui.Stretch).Role(ui.RoleGroup).Label(opts.Label).Disabled(opts.Disabled)
	v := PromptLibraryView{Element: row}
	row.Children(func() {
		ui.Column(c).Width(150).Shrink(0).Gap(2).FocusGroup(ui.Vertical).Children(func() {
			promptLibraryCategory(c, k, s, "", core.Msg(c, "prompt_library.all"))
			for _, cat := range s.shown.cats {
				promptLibraryCategory(c, k, s, cat, cat)
			}
		})
		ui.Column(c).Width(240).Shrink(0).Gap(d.Unit() * 2).Children(func() {
			input.SearchInput(c, &s.Query, input.SearchInputOptions{Placeholder: core.Msg(c, "prompt_library.search"), Disabled: opts.Disabled})
			ui.Scroll(c).MaxHeight(320).Children(func() {
				ui.Column(c).Gap(2).FocusGroup(ui.Vertical).Children(func() {
					if len(s.shown.ids) == 0 {
						ui.Text(c, core.Msg(c, "prompt_library.no_results")).TextColor(k.TextMuted).Padding(d.Unit() * 2)
					}
					for _, i := range s.shown.ids {
						t := templates[i]
						if selectableRow(c, k, t.ID, t.Title, s.Selected == t.ID) && s.Selected != t.ID {
							s.Selected = t.ID
							slog.Debug("mujica: prompt_library select", "id", t.ID)
						}
					}
				})
			})
		})
		ui.Column(c).Grow(1).MinWidth(200).Gap(d.Unit() * 2).Children(func() {
			i := slices.IndexFunc(templates, func(t PromptTemplate) bool { return t.ID == s.Selected })
			if i < 0 {
				ui.Text(c, core.Msg(c, "prompt_library.pick")).TextColor(k.TextMuted)
				return
			}
			v.Inserted = promptLibraryPreview(c, k, s, templates[i], opts.Disabled)
		})
	})
	if v.Inserted != "" {
		slog.Debug("mujica: prompt_library insert", "id", s.Selected)
	}
	return v
}

func promptLibraryFilter(sh *promptLibraryShown, templates []PromptTemplate, query, category string) {
	same := sh.built && sh.query == query && sh.category == category && len(sh.src) == len(templates) &&
		(len(templates) == 0 || &sh.src[0] == &templates[0])
	if same {
		return
	}
	*sh = promptLibraryShown{src: templates, query: query, category: category, built: true}
	q := strings.ToLower(strings.TrimSpace(query))
	for i, t := range templates {
		if t.ID == "" || t.Title == "" {
			panic("mujica: PromptTemplate needs an ID and a Title")
		}
		if t.Category != "" && !slices.Contains(sh.cats, t.Category) {
			sh.cats = append(sh.cats, t.Category)
		}
		if category != "" && t.Category != category {
			continue
		}
		if q == "" || strings.Contains(strings.ToLower(t.Title+"\n"+t.Body), q) {
			sh.ids = append(sh.ids, i)
		}
	}
}

func promptLibraryCategory(c *ui.Context, k theme.Tokens, s *PromptLibraryState, cat, label string) {
	if selectableRow(c, k, "cat:"+cat, label, s.Category == cat) && s.Category != cat {
		s.Category = cat
	}
}

// selectableRow is a focusable row marked chosen with a diamond.
func selectableRow(c *ui.Context, k theme.Tokens, key, label string, chosen bool) bool {
	d := core.Density(c)
	b := ui.ButtonBase(c).Key(key).MinHeight(d.ControlHeight()).Padding(0, d.Unit()*2).Gap(d.Unit() * 2).Justify(ui.Start).
		Radius(theme.ControlRadius).Label(label).Selected(chosen)
	b.Background(ui.Transparent).TextColor(k.Text)
	internal.Ring(b, k.Focus, theme.ControlRadius)
	switch {
	case chosen:
		b.Background(k.Selection).TextColor(k.AccentText)
	case b.Hovered():
		b.Background(k.SurfaceHover)
	}
	b.Children(func() {
		ui.Box(c).Size(8, 8).Shrink(0).Draw(func(p *ui.Painter, r ui.Rect) {
			if chosen {
				p.FillPath(internal.Diamond(r.X+4, r.Y+4, 8), k.AccentText)
			}
		})
		ui.Text(c, label).SingleLine()
	})
	return b.Clicked() && !b.IsDisabled()
}

// promptLibraryPreview shows the body, variable fields and Insert.
func promptLibraryPreview(c *ui.Context, k theme.Tokens, s *PromptLibraryState, t PromptTemplate, disabled bool) string {
	d := core.Density(c)
	vars := PromptLibraryVariables(t.Body)
	ui.Text(c, t.Title).Font(theme.SerifFont).FontSize(core.FontSize(c, theme.H3Size)).TextColor(k.Text)
	missing := false
	var spans []ui.Span
	at := 0
	for _, m := range promptLibraryVar.FindAllStringSubmatchIndex(t.Body, -1) {
		spans = append(spans, ui.Span{Text: t.Body[at:m[0]]})
		name := t.Body[m[2]:m[3]]
		if val := s.Values[name]; val != "" {
			spans = append(spans, ui.Span{Text: val, Weight: 600, Color: k.AccentText})
		} else {
			missing = true
			spans = append(spans, ui.Span{Text: "[" + name + "]", Background: k.Selection, Color: k.AccentText, Underline: true})
		}
		at = m[1]
	}
	spans = append(spans, ui.Span{Text: t.Body[at:]})
	ui.Column(c).Padding(d.Unit()*3).Background(k.Background).Border(theme.BorderWidth, k.Border).Radius(theme.ControlRadius).
		Children(func() { ui.RichText(c, spans...).Selectable() })
	for _, name := range vars {
		val := s.Values[name]
		ui.Column(c).Key("var:" + name).Children(func() {
			input.FormField(c, name, input.FormFieldOptions{Required: true}, func() *ui.Element {
				r := input.TextInput(c, &val, input.TextInputOptions{Label: name, Disabled: disabled})
				if r.Changed() {
					s.Values[name] = val
				}
				return r.Element
			})
		})
	}
	insert := input.Button(c, core.Msg(c, "prompt_library.insert"), input.ButtonOptions{Icon: icons.Must("plus"), Disabled: missing || disabled}).Clicked()
	if insert && !missing {
		return PromptLibraryFill(t.Body, s.Values)
	}
	return ""
}
