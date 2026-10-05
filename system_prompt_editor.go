package mujica

import (
	"fmt"
	"log/slog"
	"runtime"
	"unicode/utf8"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/internal"
	"mujicaui/theme"
)

// SystemPromptEditorOptions configure a SystemPromptEditor.
type SystemPromptEditorOptions struct {
	// Default is what Restore default puts back.
	Default string
	// Saved is the text last saved; a different value shows as unsaved.
	Saved string
	// MaxLength caps the characters; zero is no cap.
	MaxLength int
	// Saving shows the caller's save in progress and blocks editing.
	Saving   bool
	Disabled bool
}

// SystemPromptEditorView is a built editor and what the user asked for.
type SystemPromptEditorView struct {
	Element *ui.Element
	// Save asks the caller to save *value; Restored reports the default was put back.
	Save, Restored bool
}

// SystemPromptEditor edits a system prompt with a character count, restore default, an unsaved mark and Mod+S to save.
func SystemPromptEditor(c *ui.Context, value *string, opts SystemPromptEditorOptions) SystemPromptEditorView {
	if value == nil {
		panic("mujica: SystemPromptEditor needs a value")
	}
	if opts.MaxLength < 0 {
		panic("mujica: SystemPromptEditor MaxLength is negative")
	}
	k := Tokens(c)
	d := Density(c)
	blocked := opts.Disabled || opts.Saving
	col := ui.Column(c).Gap(d.Unit() * 2).Disabled(blocked)
	v := SystemPromptEditorView{Element: col}
	dirty := *value != opts.Saved
	col.Children(func() {
		ui.Row(c).Gap(d.Unit() * 2).AlignItems(ui.Center).Children(func() {
			ui.Text(c, msg(c, "system_prompt_editor.title")).Font(theme.SerifFont).FontSize(FontSize(c, theme.H3Size)).TextColor(k.Text).Grow(1)
			if dirty {
				ui.Row(c).Gap(d.Unit()).AlignItems(ui.Center).Role(ui.RoleStatus).Children(func() {
					ui.Box(c).Size(8, 8).Draw(func(p *ui.Painter, r ui.Rect) {
						p.FillPath(internal.Diamond(r.X+4, r.Y+4, 8), k.Warning)
					})
					ui.Text(c, msg(c, "system_prompt_editor.unsaved")).FontSize(FontSize(c, theme.CaptionSize)).TextColor(k.Warning)
				})
			}
		})
		TextArea(c, value, TextAreaOptions{Label: msg(c, "system_prompt_editor.title"), MinHeight: 120, MaxHeight: 320,
			MaxLength: opts.MaxLength, Disabled: blocked})
		ui.Row(c).Gap(d.Unit() * 2).AlignItems(ui.Center).Children(func() {
			n := utf8.RuneCountInString(*value)
			count := fmt.Sprintf(msg(c, "system_prompt_editor.count"), n)
			if opts.MaxLength > 0 {
				count = fmt.Sprintf(msg(c, "system_prompt_editor.count_of"), n, opts.MaxLength)
			}
			ui.Text(c, count).FontSize(FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted).Grow(1)
			if Button(c, msg(c, "system_prompt_editor.restore"), ButtonOptions{Variant: Ghost, Icon: icons.Must("refresh-cw"),
				Disabled: blocked || *value == opts.Default}).Clicked() && !col.IsDisabled() {
				*value, v.Restored = opts.Default, true
			}
			v.Save = Button(c, msg(c, "system_prompt_editor.save"), ButtonOptions{Loading: opts.Saving, Disabled: opts.Disabled || !dirty}).Clicked()
		})
	})
	mods := ui.Ctrl
	if runtime.GOOS == "darwin" {
		mods = ui.Cmd
	}
	if col.Shortcut(mods, ui.KeyS) && dirty {
		v.Save = true
	}
	if col.IsDisabled() {
		v.Save = false
	}
	switch {
	case v.Save:
		slog.Debug("mujica: system_prompt_editor save", "len", len(*value))
	case v.Restored:
		slog.Debug("mujica: system_prompt_editor restore")
	}
	return v
}
