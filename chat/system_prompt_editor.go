package chat

import (
	"fmt"
	"log/slog"
	"runtime"
	"unicode/utf8"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/icons"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// SystemPromptEditorOptions configure a SystemPromptEditor.
type SystemPromptEditorOptions struct {
	// Default is what Restore default puts back.
	Default string
	// Saved is the text last saved; differences show as unsaved.
	Saved string
	// MaxLength limits characters, blocking Save beyond it; zero is none.
	MaxLength int
	// Saving shows the caller's save in progress and blocks editing.
	Saving   bool
	Disabled bool
}

// SystemPromptEditorView is a built editor and its requests.
type SystemPromptEditorView struct {
	Element *ui.Element
	// Save asks to save *value; Restored reports the default returned.
	Save, Restored bool
}

// SystemPromptEditor edits a system prompt with count, restore and save.
func SystemPromptEditor(c *ui.Context, value *string, opts SystemPromptEditorOptions) SystemPromptEditorView {
	if value == nil {
		panic("mujica: SystemPromptEditor needs a value")
	}
	if opts.MaxLength < 0 {
		panic("mujica: SystemPromptEditor MaxLength is negative")
	}
	k := core.Tokens(c)
	d := core.Density(c)
	blocked := opts.Disabled || opts.Saving
	col := ui.Column(c).Gap(d.Unit() * 2).Disabled(blocked)
	v := SystemPromptEditorView{Element: col}
	dirty := *value != opts.Saved
	over := opts.MaxLength > 0 && utf8.RuneCountInString(*value) > opts.MaxLength
	col.Children(func() {
		ui.Row(c).Gap(d.Unit() * 2).AlignItems(ui.Center).Children(func() {
			ui.Text(c, core.Msg(c, "system_prompt_editor.title")).Font(theme.SerifFont).FontSize(core.FontSize(c, theme.H3Size)).TextColor(k.Text).Grow(1)
			if dirty {
				ui.Row(c).Gap(d.Unit()).AlignItems(ui.Center).Role(ui.RoleStatus).Children(func() {
					ui.Box(c).Size(8, 8).Draw(func(p *ui.Painter, r ui.Rect) {
						p.FillPath(internal.Diamond(r.X+4, r.Y+4, 8), k.Warning)
					})
					ui.Text(c, core.Msg(c, "system_prompt_editor.unsaved")).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.Warning)
				})
			}
		})
		input.TextArea(c, value, input.TextAreaOptions{Label: core.Msg(c, "system_prompt_editor.title"), MinHeight: 120, MaxHeight: 320,
			MaxLength: opts.MaxLength, Disabled: blocked})
		ui.Row(c).Gap(d.Unit() * 2).AlignItems(ui.Center).Children(func() {
			// With MaxLength the TextArea shows "n / max" itself.
			count := ""
			if opts.MaxLength == 0 {
				count = fmt.Sprintf(core.Msg(c, "system_prompt_editor.count"), utf8.RuneCountInString(*value))
			}
			ui.Text(c, count).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted).Grow(1)
			if input.Button(c, core.Msg(c, "system_prompt_editor.restore"), input.ButtonOptions{Variant: input.Ghost, Icon: icons.Must("refresh-cw"),
				Disabled: blocked || *value == opts.Default}).Clicked() && !col.IsDisabled() {
				*value, v.Restored = opts.Default, true
			}
			v.Save = input.Button(c, core.Msg(c, "system_prompt_editor.save"), input.ButtonOptions{Loading: opts.Saving, Disabled: opts.Disabled || !dirty || over}).Clicked()
		})
	})
	mods := ui.Ctrl
	if runtime.GOOS == "darwin" {
		mods = ui.Cmd
	}
	if col.Shortcut(mods, ui.KeyS) && dirty && !over {
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
