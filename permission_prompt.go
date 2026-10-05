package mujica

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/theme"
)

// PermissionChoice answers a PermissionPrompt; zero is no answer.
type PermissionChoice int

const (
	PermissionNone PermissionChoice = iota
	PermissionAllowOnce
	PermissionAllowAlways
	PermissionDeny
)

// PermissionPromptOptions configure a PermissionPrompt. Title and Scope are required.
type PermissionPromptOptions struct {
	Title       string
	Description string
	// Scope says what "always" covers, as "Write access to ~/court".
	Scope string
}

// PermissionResult reports the choice made in this frame.
type PermissionResult struct {
	Element *ui.Element
	Choice  PermissionChoice
}

// PermissionPrompt asks allow once, always, or deny; Deny has focus.
func PermissionPrompt(c *ui.Context, open *bool, opts PermissionPromptOptions) PermissionResult {
	if opts.Title == "" || opts.Scope == "" {
		panic("mujica: PermissionPrompt needs a Title and a Scope")
	}
	var r PermissionResult
	was := *open
	k := Tokens(c)
	u := Density(c).Unit()
	r.Element = Dialog(c, open, DialogOptions{
		Title: opts.Title, Description: opts.Description, KeepOnBackdrop: true, HideClose: true,
		Actions: func() {
			if Button(c, msg(c, "permission_prompt.deny"), ButtonOptions{Variant: Secondary}).AutoFocus().Clicked() {
				r.Choice = PermissionDeny
			}
			if Button(c, msg(c, "permission_prompt.always"), ButtonOptions{Variant: Outline}).Clicked() {
				r.Choice = PermissionAllowAlways
			}
			if Button(c, msg(c, "permission_prompt.once"), ButtonOptions{}).Clicked() {
				r.Choice = PermissionAllowOnce
			}
		},
	}, func() {
		ui.Row(c).Gap(u*2).AlignItems(ui.Start).Padding(u*2, u*3).Background(k.Background).Radius(theme.ControlRadius).
			Border(1, k.Border).Children(func() {
			ui.Icon(c, icons.Must("lock")).Size(16, 16).TextColor(k.Info).Shrink(0).Margin(2, 0, 0, 0)
			ui.Column(c).Gap(u / 2).Shrink(1).Children(func() {
				agentCaption(c, msg(c, "permission_prompt.scope"))
				ui.Text(c, opts.Scope)
			})
		})
	})
	if was && !*open && r.Choice == PermissionNone {
		r.Choice = PermissionDeny // Escape
	}
	if r.Choice != PermissionNone {
		slog.Debug("mujica: permission_prompt choose", "choice", int(r.Choice))
		*open = false
	}
	return r
}
