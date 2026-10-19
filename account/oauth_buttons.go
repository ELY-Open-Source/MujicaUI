package account

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/input"
)

// OAuthProvider is a sign-in provider; icons are plain, not trademarks.
type OAuthProvider struct {
	ID, Name string
	Icon     *ui.SVG
}

// OAuthButtonsOptions configure OAuthButtons.
type OAuthButtonsOptions struct {
	// Pending is the provider signing in; the others are blocked.
	Pending  string
	Disabled bool
}

// OAuthButtonsResult carries the selected provider.
type OAuthButtonsResult struct {
	Element ui.Element
	chosen  string
}

// Chosen returns the ID of the provider the user picked.
func (r OAuthButtonsResult) Chosen() (string, bool) { return r.chosen, r.chosen != "" }

// OAuthButtons offers a button per provider; the caller signs in.
func OAuthButtons(c *ui.Context, providers []OAuthProvider, opts OAuthButtonsOptions) OAuthButtonsResult {
	if len(providers) == 0 {
		panic("mujica: OAuthButtons needs providers")
	}
	d := core.Density(c)
	var r OAuthButtonsResult
	r.Element = ui.Column(c).Gap(d.Unit() * 2).Width(340).Disabled(opts.Disabled)
	disabled := r.Element.IsDisabled()
	r.Element.Children(func() {
		for _, p := range providers {
			if p.ID == "" || p.Name == "" {
				panic("mujica: OAuthProvider needs an ID and a Name")
			}
			label := fmt.Sprintf(core.Msg(c, "oauth_buttons.continue"), p.Name)
			ui.Column(c.Key(p.ID)).AlignItems(ui.Stretch).Children(func() {
				b := input.Button(c, label, input.ButtonOptions{Variant: input.Secondary, Icon: p.Icon, Loading: opts.Pending == p.ID,
					Disabled: opts.Pending != "" && opts.Pending != p.ID})
				if b.Clicked() && !disabled && opts.Pending == "" {
					r.chosen = p.ID
					slog.Debug("mujica: oauth_buttons chosen", "provider", p.ID)
				}
			})
		}
	})
	return r
}
