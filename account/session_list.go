package account

import (
	"fmt"
	"log/slog"
	"slices"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/display"
	"github.com/ELY-Open-Source/MujicaUI/input"
	"github.com/ELY-Open-Source/MujicaUI/overlay"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

// Session describes a device without carrying its authentication token.
type Session struct {
	ID, Device, Location string
	LastActive           time.Time
	// Current marks the session this window runs in.
	Current bool
}

// SessionListOptions configure a SessionList.
type SessionListOptions struct {
	Disabled bool
}

// SessionListResult reports individual and bulk sign-out requests.
type SessionListResult struct {
	Element ui.Element
	revoke  string
	others  bool
}

// Revoked returns the ID of a session to sign out.
func (r SessionListResult) Revoked() (string, bool) { return r.revoke, r.revoke != "" }

// OthersRevoked reports confirmed bulk sign-out.
func (r SessionListResult) OthersRevoked() bool { return r.others }

type sessionListState struct {
	confirm bool
	source  []Session
	others  int
	ready   bool
}

// SessionList displays devices and requests sign-out through events.
func SessionList(c *ui.Context, sessions []Session, opts SessionListOptions) SessionListResult {
	k := core.Tokens(c)
	d := core.Density(c)
	var r SessionListResult
	r.Element = ui.Column(c).Gap(d.Unit() * 3).MinWidth(420).Disabled(opts.Disabled)
	st := ui.Local(r.Element, "sessions", func() sessionListState { return sessionListState{} })
	disabled := r.Element.IsDisabled()
	if !st.ready || !slices.Equal(st.source, sessions) {
		st.source, st.others, st.ready = slices.Clone(sessions), 0, true
		ids := map[string]bool{}
		for _, s := range sessions {
			if s.ID == "" || s.Device == "" || ids[s.ID] {
				panic("mujica: Session needs a unique ID and a Device")
			}
			ids[s.ID] = true
			if !s.Current {
				st.others++
			}
		}
	}
	others := st.others
	r.Element.Children(func() {
		ui.Column(c).Role(ui.RoleList).Label(core.Msg(c, "session_list.label")).Children(func() {
			accountRows(c, len(sessions), func(i int) any { return sessions[i].ID }, func(i int) {
				s := sessions[i]
				ui.Row(c.Key(s.ID)).MinHeight(d.ControlHeight()+d.Unit()*4).Gap(d.Unit()*3).AlignItems(ui.Center).
					BorderWidth(0, 0, theme.BorderWidth, 0).BorderColor(k.Border).Role(ui.RoleListItem).Children(func() {
					ui.Column(c).Grow(1).MinWidth(0).Gap(2).Children(func() {
						ui.Row(c).Gap(d.Unit() * 2).AlignItems(ui.Center).Children(func() {
							ui.Text(c, s.Device).FontWeight(600).SingleLine()
							if s.Current {
								display.Badge(c, display.BadgeOptions{Text: core.Msg(c, "session_list.current")}, nil)
							}
						})
						ui.Text(c, sessionListDetail(c, s)).TextColor(k.TextMuted).FontSize(core.FontSize(c, theme.CaptionSize)).SingleLine()
					})
					if !s.Current {
						b := input.Button(c, core.Msg(c, "session_list.sign_out"), input.ButtonOptions{Variant: input.Secondary, Disabled: disabled}).
							Label(fmt.Sprintf(core.Msg(c, "session_list.sign_out_of"), s.Device))
						if b.Clicked() && !disabled {
							r.revoke = s.ID
						}
					}
				})
			})
		})
		ui.Row(c).Children(func() {
			if input.Button(c, core.Msg(c, "session_list.sign_out_others"), input.ButtonOptions{Variant: input.Danger, Disabled: disabled || others == 0}).Clicked() && !disabled {
				st.confirm = true
			}
		})
	})
	if disabled || others == 0 {
		st.confirm = false
	}
	res := overlay.AlertDialog(c, &st.confirm, overlay.AlertDialogOptions{Title: core.Msg(c, "session_list.confirm_title"),
		Description: fmt.Sprintf(core.Msg(c, "session_list.confirm_text"), others), ConfirmLabel: core.Msg(c, "session_list.sign_out_others")})
	r.others = res.Confirmed && !disabled
	if r.revoke != "" || r.others {
		slog.Debug("mujica: session_list sign out", "others", r.others)
	}
	return r
}

// sessionListDetail formats the location and last active time.
func sessionListDetail(c *ui.Context, s Session) string {
	when := core.Msg(c, "session_list.now")
	if ago := c.Now().Sub(s.LastActive); s.LastActive.IsZero() {
		when = core.Msg(c, "session_list.unknown")
	} else if ago >= time.Minute {
		when = s.LastActive.Format("2006-01-02 15:04")
	}
	if s.Location == "" {
		return when
	}
	return s.Location + " · " + when
}
