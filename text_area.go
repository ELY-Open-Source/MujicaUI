package mujica

import (
	"fmt"
	"unicode/utf8"

	"github.com/egoist/mygo/ui"

	"mujicaui/theme"
)

// TextAreaOptions configure a TextArea.
type TextAreaOptions struct {
	Placeholder string
	// Label names the field for assistive technology.
	Label string
	// MinHeight is the least height in DIPs; zero fits three lines.
	MinHeight float32
	// MaxHeight caps the height in DIPs, then the text scrolls; zero grows freely.
	MaxHeight float32
	// MaxLength shows a character count against it; zero hides the count.
	MaxLength int
	// SubmitKey with SubmitMods submits; a zero key means Cmd+Enter.
	SubmitKey  ui.Key
	SubmitMods ui.Modifiers
	Disabled   bool
	ReadOnly   bool
	// Error marks the value invalid and shows the message below.
	Error string
}

// TextArea is a multi-line field editing *value.
func TextArea(c *ui.Context, value *string, opts TextAreaOptions) InputResult {
	if value == nil {
		panic("mujica: TextArea needs a value")
	}
	if opts.MinHeight < 0 || opts.MaxHeight < 0 || opts.MaxHeight > 0 && opts.MaxHeight < opts.MinHeight {
		panic("mujica: TextArea needs 0 <= MinHeight <= MaxHeight")
	}
	if opts.MaxLength < 0 {
		panic("mujica: TextArea MaxLength is negative")
	}
	key, mods := opts.SubmitKey, opts.SubmitMods
	if key == 0 {
		key, mods = ui.KeyEnter, ui.Cmd
	}
	d := Density(c)
	n := utf8.RuneCountInString(*value)
	over := opts.MaxLength > 0 && n > opts.MaxLength
	errMsg := opts.Error
	if errMsg == "" && over {
		errMsg = fmt.Sprintf(msg(c, "text_area.over"), n-opts.MaxLength)
	}
	count := textAreaCount(c, n, opts.MaxLength)
	var r InputResult
	look := &textInputLook{disabled: opts.Disabled, readOnly: opts.ReadOnly, err: errMsg, area: true}
	r.Element, r.Input = textInputFrame(c, look, func(frame *ui.Element) *ui.Element {
		in := textInputEditor(c, frame, value, opts.ReadOnly || opts.Disabled, ui.TextAreaBase).Grow(1).Padding(d.Unit()*2, d.Unit()*3).
			Placeholder(opts.Placeholder).Label(opts.Label)
		if opts.MinHeight > 0 {
			in.MinHeight(opts.MinHeight)
		}
		if opts.MaxHeight > 0 {
			in.MaxHeight(opts.MaxHeight)
		}
		r.changed = in.Changed() && !opts.ReadOnly && !opts.Disabled
		// Taken before the editor, which would make a plain Enter a newline.
		submit := ui.Local(frame, "submit", func() bool { return false })
		in.HandleInput(func(ev ui.InputEvent) bool {
			if ev.Kind == ui.InputKeyDown && ev.Key == key && ev.Mods == mods {
				*submit = true
				return true
			}
			return false
		})
		r.submitted = *submit && !opts.Disabled
		*submit = false
		return in
	}, count)
	return r
}

// textAreaCount builds the character count against max, or nil without one.
func textAreaCount(c *ui.Context, n, max int) func() {
	if max == 0 {
		return nil
	}
	return func() {
		k := Tokens(c)
		color := k.TextMuted
		over := n > max
		if over {
			color = k.Danger
		}
		ui.Text(c, fmt.Sprintf("%d / %d", n, max)).TextColor(color).FontSize(FontSize(c, theme.CaptionSize)).
			FontFeatures("tnum").Shrink(0).Label(fmt.Sprintf(msg(c, "text_area.count"), n, max))
	}
}
