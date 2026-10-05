package mujica

import (
	"fmt"
	"log/slog"
	"slices"
	"strings"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/internal"
	"mujicaui/theme"
)

// TagsInputOptions configure a TagsInput.
type TagsInputOptions struct {
	// Max caps the number of tags; zero allows any number.
	Max int
	// Separators split typed and pasted text into tags; empty uses ",;，；" and line breaks.
	Separators string
	// Validate checks a tag before it is added and returns an error message, or "".
	Validate    func(tag string) string
	Placeholder string
	// Label names the input when it is not inside a FormField.
	Label    string
	Disabled bool
	// Invalid draws the border in the danger color.
	Invalid bool
}

// TagsInputParts are the parts of a TagsInput.
type TagsInputParts struct {
	Element *ui.Element
	// Input takes the focus.
	Input   *ui.Element
	changed bool
}

// Changed reports whether tags were added or removed this frame.
func (p TagsInputParts) Changed() bool { return p.changed }

type tagsInputState struct {
	query     string
	err       string
	backspace bool
}

// TagsInput edits a list of unique tags shown as chips before a text input.
func TagsInput(c *ui.Context, tags *[]string, opts TagsInputOptions) TagsInputParts {
	if opts.Max < 0 {
		panic("mujica: TagsInput Max must not be negative")
	}
	k := Tokens(c)
	d := Density(c)
	seps := opts.Separators
	if seps == "" {
		seps = ",;，；\n\r\t"
	}
	parts := TagsInputParts{}
	outer := ui.Column(c).Gap(d.Unit())
	parts.Element = outer
	st := ui.Local(outer, "tags", func() tagsInputState { return tagsInputState{} })
	add := func(s string) bool {
		s = strings.TrimSpace(s)
		switch {
		case s == "":
			return true
		case slices.Contains(*tags, s):
			st.err = fmt.Sprintf(msg(c, "tags_input.duplicate"), s)
			return false
		case opts.Max > 0 && len(*tags) >= opts.Max:
			st.err = fmt.Sprintf(msg(c, "tags_input.max"), opts.Max)
			return false
		}
		if opts.Validate != nil {
			if e := opts.Validate(s); e != "" {
				st.err = e
				return false
			}
		}
		*tags = append(*tags, s)
		parts.changed = true
		slog.Debug("mujica: tags input added", "count", len(*tags))
		return true
	}
	remove := func(i int) {
		*tags = slices.Delete(slices.Clone(*tags), i, i+1)
		parts.changed, st.err = true, ""
		slog.Debug("mujica: tags input removed", "count", len(*tags))
	}
	if st.backspace {
		st.backspace = false
		if len(*tags) > 0 && st.query == "" {
			remove(len(*tags) - 1)
		}
	}
	var box *ui.Element
	outer.Children(func() {
		box = ui.Row(c).Wrap().AlignItems(ui.Center).Gap(d.Unit()*1.5).MinHeight(d.ControlHeight()).
			Padding(d.Unit(), d.Unit()*2).Radius(theme.ControlRadius).Background(k.Surface).Disabled(opts.Disabled)
		box.Children(func() {
			for i, tag := range *tags {
				chip := ui.Row(c).Key(tag).AlignItems(ui.Center).Padding(0, 0, 0, d.Unit()*2).Radius(theme.ControlRadius).Background(k.Selection).Shrink(0)
				chip.Children(func() {
					ui.Text(c, tag).SingleLine()
					x := ui.ButtonBase(c).Size(theme.MinHitSize, theme.MinHitSize).Radius(theme.ControlRadius).TextColor(k.TextMuted).
						Label(fmt.Sprintf(msg(c, "tags_input.remove"), tag))
					internal.Ring(x, k.Focus, theme.ControlRadius)
					if x.Hovered() {
						x.Background(k.SurfaceHover).TextColor(k.Text)
					}
					x.Children(func() { ui.Icon(c, icons.Must("x")).Size(14, 14) })
					if x.Clicked() {
						remove(i)
					}
				})
			}
			// The input keeps its identity whatever chips come before it.
			ui.Box(c).Key("input").Grow(1).MinWidth(d.Unit() * 24).Children(func() {
				in := ui.TextInputBase(c, &st.query).FocusRing(false).FillWidth().Placeholder(opts.Placeholder)
				parts.Input = in
				if opts.Label != "" {
					in.Label(opts.Label)
				}
				in.HandleInput(func(ev ui.InputEvent) bool {
					if ev.Kind == ui.InputKeyDown && ev.Mods == 0 && ev.Key == ui.KeyBackspace && st.query == "" {
						st.backspace = true
						return true
					}
					return false
				})
			})
		})
		if st.err != "" {
			ui.Text(c, st.err).TextColor(k.Danger).FontSize(FontSize(c, theme.CaptionSize)).Role(ui.RoleStatus)
		}
	})
	in := parts.Input
	if in.Changed() {
		st.err = ""
	}
	if strings.ContainsAny(st.query, seps) {
		pieces := strings.FieldsFunc(st.query+"\x00", func(r rune) bool { return strings.ContainsRune(seps, r) })
		last := strings.TrimSuffix(pieces[len(pieces)-1], "\x00")
		var kept []string
		for _, p := range pieces[:len(pieces)-1] {
			if !add(p) {
				kept = append(kept, strings.TrimSpace(p))
			}
		}
		st.query = strings.Join(append(kept, last), string([]rune(seps)[0]))
	}
	if in.Submitted() && add(st.query) {
		st.query = ""
	}
	if box.Clicked() {
		in.Focus()
	}
	border := k.ControlBorder
	switch {
	case opts.Invalid || st.err != "":
		border = k.Danger
	case in.Focused():
		border = k.Focus
	case box.Hovered() && !opts.Disabled:
		border = k.ControlBorder.Mix(k.Text, 0.25)
	}
	box.Border(theme.BorderWidth, border)
	box.DrawOver(func(p *ui.Painter, r ui.Rect) {
		if in.Focused() {
			internal.FocusRing(p, r, theme.ControlRadius, k.Focus)
		}
	})
	return parts
}
