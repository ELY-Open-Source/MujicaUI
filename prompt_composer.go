package mujica

import (
	"log/slog"
	"strings"

	"github.com/egoist/mygo/ui"
)

// PromptComposerOptions configure a PromptComposer.
type PromptComposerOptions struct {
	Placeholder string
	// Label names the field; empty uses the built-in "Message".
	Label string
	// CmdEnterSends sends with Cmd+Enter (Ctrl+Enter) and makes Enter a newline.
	CmdEnterSends bool
	// Composing blocks submission while the caller's input method composes.
	Composing bool
	// MaxHeight caps input height; zero means 200 DIPs.
	MaxHeight float32
	// Attachments builds the row above the field, as AttachmentChips.
	Attachments func()
	// Tools and Actions build opposite ends of the toolbar.
	Tools, Actions func()
	Disabled       bool
}

// promptHooks hold menu handlers for the current frame.
type promptHooks struct{ list []func(ui.InputEvent) bool }

type promptHooksKey struct{}

// PromptComposer edits prompts with configurable send shortcuts.
// Submitted never fires for blank text.
func PromptComposer(c *ui.Context, value *string, opts PromptComposerOptions) InputResult {
	if value == nil {
		panic("mujica: PromptComposer needs a value")
	}
	if opts.MaxHeight < 0 {
		panic("mujica: PromptComposer MaxHeight must not be negative")
	}
	u := Density(c).Unit()
	maxH := opts.MaxHeight
	if maxH == 0 {
		maxH = 200
	}
	label := opts.Label
	if label == "" {
		label = msg(c, "prompt_composer.label")
	}
	var r InputResult
	r.Element = ui.Column(c).Gap(u * 2).MinWidth(0).Disabled(opts.Disabled)
	hooks := ui.Local(r.Element, promptHooksKey{}, func() *promptHooks { return &promptHooks{} })
	h := *hooks
	h.list = h.list[:0]
	r.Element.Children(func() {
		if opts.Attachments != nil {
			ui.Row(c).Gap(u * 2).Wrap().Children(opts.Attachments)
		}
		ta := TextAreaOptions{Label: label, Placeholder: opts.Placeholder, MaxHeight: maxH,
			keys: func(ev ui.InputEvent) bool {
				if r.Element.IsDisabled() {
					return true
				}
				if opts.Composing && ev.Kind == ui.InputKeyDown {
					return ev.Key == ui.KeyEnter
				}
				for _, f := range h.list {
					if f(ev) {
						return true
					}
				}
				return false
			}}
		if !opts.CmdEnterSends {
			ta.SubmitKey = ui.KeyEnter
		}
		area := TextArea(c, value, ta)
		r.Input, r.changed = area.Input, area.Changed()
		r.submitted = area.Submitted() && !opts.Composing && strings.TrimSpace(*value) != ""
		if opts.Tools != nil || opts.Actions != nil {
			ui.Row(c).Gap(u * 2).AlignItems(ui.Center).Children(func() {
				ui.Row(c).Gap(u).AlignItems(ui.Center).Grow(1).Children(func() {
					if opts.Tools != nil {
						opts.Tools()
					}
				})
				if opts.Actions != nil {
					opts.Actions()
				}
			})
		}
	})
	if r.submitted {
		slog.Debug("mujica: prompt_composer submit", "chars", len(*value))
	}
	return r
}

// promptComposerHook registers a handler before editor input.
func promptComposerHook(composer InputResult, f func(ui.InputEvent) bool) {
	if composer.Element == nil {
		panic("mujica: menus need a built PromptComposer")
	}
	h := *ui.Local(composer.Element, promptHooksKey{}, func() *promptHooks { return &promptHooks{} })
	h.list = append(h.list, f)
}
