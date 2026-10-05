package mujica

import (
	"log/slog"
	"strings"

	"github.com/egoist/mygo/ui"
)

// MessageEditorOptions configure a MessageEditor.
type MessageEditorOptions struct {
	// Original is the sent text; saving it unchanged is not allowed.
	Original string
}

// MessageEditorResult is the editor's element and its outcome this frame.
type MessageEditorResult struct {
	Element         *ui.Element
	saved, canceled bool
}

// Saved reports a save; the caller makes a new branch from the draft.
func (r MessageEditorResult) Saved() bool { return r.saved }

// Canceled reports Cancel or Escape.
func (r MessageEditorResult) Canceled() bool { return r.canceled }

// MessageEditor edits a sent message in *draft; Cmd+Enter saves, Escape cancels.
func MessageEditor(c *ui.Context, draft *string, opts MessageEditorOptions) MessageEditorResult {
	if draft == nil {
		panic("mujica: MessageEditor needs a draft")
	}
	u := Density(c).Unit()
	var r MessageEditorResult
	r.Element = ui.Column(c).Gap(u * 2).MinWidth(0)
	disabled := r.Element.IsDisabled()
	r.Element.Children(func() {
		area := TextArea(c, draft, TextAreaOptions{Label: msg(c, "message_editor.label"), MaxHeight: 240})
		ok := strings.TrimSpace(*draft) != "" && *draft != opts.Original
		escape := area.Input.Shortcut(0, ui.KeyEscape)
		ui.Row(c).Gap(u * 2).Justify(ui.End).Children(func() {
			if Button(c, msg(c, "message_editor.cancel"), ButtonOptions{Variant: Secondary}).Clicked() || escape {
				r.canceled = !disabled
			}
			if Button(c, msg(c, "message_editor.save"), ButtonOptions{Disabled: !ok}).Clicked() || area.Submitted() && ok {
				r.saved = !disabled
			}
		})
	})
	if r.saved {
		slog.Debug("mujica: message_editor save", "chars", len(*draft))
	}
	if r.canceled {
		slog.Debug("mujica: message_editor cancel")
	}
	return r
}
