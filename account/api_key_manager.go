package account

import (
	"fmt"
	"log/slog"
	"slices"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/feedback"
	"github.com/ZacharyZhang-NY/MujicaUI/icons"
	"github.com/ZacharyZhang-NY/MujicaUI/input"
	"github.com/ZacharyZhang-NY/MujicaUI/overlay"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
)

// APIKey contains public metadata and the visible prefix.
type APIKey struct {
	ID, Name, Prefix  string
	Created, LastUsed time.Time
}

// ApiKeyManagerState holds creation, confirmation and one-time reveal state.
type ApiKeyManagerState struct {
	creating, pending bool
	name, secret      string
	err               string
	revoking          string
	keys              []APIKey
	index             map[string]APIKey
}

// Reveal presents the generated secret until dismissal.
func (s *ApiKeyManagerState) Reveal(secret string) {
	if !s.pending || secret == "" {
		panic("mujica: ApiKeyManagerState.Reveal needs a create in progress and a secret")
	}
	s.pending, s.secret, s.err = false, secret, ""
	slog.Debug("mujica: api_key_manager revealed")
}

// Fail displays a caller-provided creation error.
func (s *ApiKeyManagerState) Fail(err string) {
	if !s.pending {
		panic("mujica: ApiKeyManagerState.Fail needs a create in progress")
	}
	s.creating, s.pending, s.err = true, false, err
}

// ApiKeyManagerOptions configure an ApiKeyManager.
type ApiKeyManagerOptions struct {
	Disabled bool
}

// ApiKeyManagerResult reports creation and revocation requests.
type ApiKeyManagerResult struct {
	Element        *ui.Element
	create, revoke string
}

// Create returns the requested name; finish with Reveal or Fail.
func (r ApiKeyManagerResult) Create() (string, bool) { return r.create, r.create != "" }

// Revoked returns the confirmed key ID.
func (r ApiKeyManagerResult) Revoked() (string, bool) { return r.revoke, r.revoke != "" }

// ApiKeyManager masks stored keys and confirms revocation.
func ApiKeyManager(c *ui.Context, s *ApiKeyManagerState, keys []APIKey, opts ApiKeyManagerOptions) ApiKeyManagerResult {
	k := core.Tokens(c)
	d := core.Density(c)
	var r ApiKeyManagerResult
	r.Element = ui.Column(c).Gap(d.Unit() * 3).MinWidth(460).Disabled(opts.Disabled)
	disabled := r.Element.IsDisabled()
	if disabled {
		s.creating = s.creating && (s.pending || s.err != "")
		s.revoking, s.secret = "", ""
	}
	if s.index == nil || !slices.Equal(s.keys, keys) {
		s.keys, s.index = slices.Clone(keys), map[string]APIKey{}
		for _, key := range keys {
			if key.ID == "" || key.Name == "" || key.Prefix == "" {
				panic("mujica: APIKey needs an ID, a Name and a Prefix")
			}
			if _, ok := s.index[key.ID]; ok {
				panic("mujica: APIKey IDs must be unique")
			}
			s.index[key.ID] = key
		}
	}
	byID := s.index
	r.Element.Children(func() {
		ui.Row(c).Children(func() {
			if input.Button(c, core.Msg(c, "api_key_manager.create"), input.ButtonOptions{Icon: icons.Must("plus"), Disabled: disabled}).Clicked() && !disabled {
				s.creating, s.name, s.err = true, "", ""
			}
		})
		if len(keys) == 0 {
			feedback.EmptyState(c, feedback.EmptyStateOptions{Title: core.Msg(c, "api_key_manager.empty")})
		}
		accountRows(c, len(keys), func(i int) any { return keys[i].ID }, func(i int) {
			key := keys[i]
			ui.Row(c).Key(key.ID).MinHeight(d.ControlHeight()+d.Unit()*4).Gap(d.Unit()*3).AlignItems(ui.Center).
				BorderWidth(0, 0, theme.BorderWidth, 0).BorderColor(k.Border).Children(func() {
				ui.Column(c).Grow(1).MinWidth(0).Gap(2).Children(func() {
					ui.Text(c, key.Name).FontWeight(600).SingleLine()
					ui.Text(c, apiKeyMasked(key.Prefix)).Font(theme.MonoFont).TextColor(k.TextMuted).SingleLine()
					ui.Text(c, apiKeyDates(c, key)).TextColor(k.TextMuted).FontSize(core.FontSize(c, theme.CaptionSize))
				})
				b := input.Button(c, core.Msg(c, "api_key_manager.revoke"), input.ButtonOptions{Variant: input.Secondary, Disabled: disabled}).
					Label(fmt.Sprintf(core.Msg(c, "api_key_manager.revoke_named"), key.Name))
				if b.Clicked() && !disabled {
					s.revoking = key.ID
				}
			})
		})
	})
	if name, ok := apiKeyCreateDialog(c, s, disabled); ok && !disabled {
		r.create = name
		slog.Debug("mujica: api_key_manager create")
	}
	victim, exists := byID[s.revoking]
	open := s.revoking != "" && exists
	res := overlay.AlertDialog(c, &open, overlay.AlertDialogOptions{Title: fmt.Sprintf(core.Msg(c, "api_key_manager.revoke_title"), victim.Name),
		Description: core.Msg(c, "api_key_manager.revoke_text"), ConfirmLabel: core.Msg(c, "api_key_manager.revoke")})
	if res.Confirmed && !disabled {
		r.revoke = victim.ID
		slog.Debug("mujica: api_key_manager revoke")
	}
	if !open {
		s.revoking = ""
	}
	return r
}

// apiKeyMasked shows the prefix and hides the rest.
func apiKeyMasked(prefix string) string { return prefix + strings.Repeat("•", 12) }

// apiKeyDates says when a key was made and last used.
func apiKeyDates(c *ui.Context, key APIKey) string {
	used := core.Msg(c, "api_key_manager.never")
	if !key.LastUsed.IsZero() {
		used = key.LastUsed.Format("2006-01-02")
	}
	return fmt.Sprintf(core.Msg(c, "api_key_manager.dates"), key.Created.Format("2006-01-02"), used)
}

// apiKeyCreateDialog collects a name and reveals the generated secret.
func apiKeyCreateDialog(c *ui.Context, s *ApiKeyManagerState, disabled bool) (string, bool) {
	k := core.Tokens(c)
	d := core.Density(c)
	open := !disabled && (s.creating || s.pending || s.secret != "")
	was := open
	asked := false
	title := core.Msg(c, "api_key_manager.create_title")
	if s.secret != "" {
		title = core.Msg(c, "api_key_manager.secret_title")
	}
	overlay.Dialog(c, &open, overlay.DialogOptions{Title: title, KeepOnBackdrop: true, KeepOnEscape: s.pending, HideClose: s.pending, Actions: func() {
		switch {
		case s.secret != "":
			if input.Button(c, core.Msg(c, "api_key_manager.done"), input.ButtonOptions{}).Clicked() {
				open = false
			}
		default:
			if input.Button(c, core.Msg(c, "api_key_manager.cancel"), input.ButtonOptions{Variant: input.Secondary, Disabled: s.pending}).Clicked() {
				open = false
			}
			ok := strings.TrimSpace(s.name) != ""
			if input.Button(c, core.Msg(c, "api_key_manager.create"), input.ButtonOptions{Loading: s.pending, Disabled: !ok}).Clicked() && ok && !s.pending {
				asked = true
			}
		}
	}}, func() {
		if s.secret != "" {
			ui.Row(c).Gap(d.Unit() * 2).AlignItems(ui.Start).Padding(d.Unit() * 3).Radius(theme.ControlRadius).Background(k.Warning.Alpha(0.12)).Children(func() {
				ui.Icon(c, icons.Must("triangle-alert")).Size(16, 16).TextColor(k.Warning).Shrink(0)
				ui.Text(c, core.Msg(c, "api_key_manager.once")).Grow(1)
			})
			ui.Row(c).Gap(d.Unit()*2).AlignItems(ui.Center).Padding(d.Unit()*2, d.Unit()*3).Radius(theme.ControlRadius).
				Background(k.Background).Border(theme.BorderWidth, k.Border).Children(func() {
				ui.Text(c, s.secret).Font(theme.MonoFont).Selectable().Label(core.Msg(c, "api_key_manager.secret_title")).Grow(1).MinWidth(0)
				if core.IconAction(c, icons.Must("copy"), core.Msg(c, "api_key_manager.copy")).Clicked() {
					c.WriteClipboard(s.secret)
					c.Announce(core.Msg(c, "api_key_manager.copied"))
				}
			})
			return
		}
		res := input.TextInput(c, &s.name, input.TextInputOptions{Label: core.Msg(c, "api_key_manager.name"), Placeholder: core.Msg(c, "api_key_manager.name"),
			Disabled: s.pending, Error: s.err})
		res.Input.AutoFocus()
		if res.Submitted() && strings.TrimSpace(s.name) != "" && !s.pending {
			asked = true
		}
	})
	if was && !open && !s.pending {
		// Dismissing wipes the one-time secret.
		s.creating, s.secret, s.name = false, "", ""
	}
	if asked {
		s.pending, s.err = true, ""
		return strings.TrimSpace(s.name), true
	}
	return "", false
}
