package messaging

import (
	"fmt"
	"log/slog"
	"net/mail"
	"slices"
	"strings"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/display"
	"github.com/ELY-Open-Source/MujicaUI/input"
	"github.com/ELY-Open-Source/MujicaUI/internal"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

// Contact is a person RecipientInput can suggest.
type Contact struct {
	Name, Email string
}

// RecipientInputOptions configure a RecipientInput.
type RecipientInputOptions struct {
	// Contacts supplies suggestions matched by name or address.
	Contacts    []Contact
	Placeholder string
	// Label names the field, as "To".
	Label    string
	Disabled bool
}

// recipientMax is how many suggestions show.
const recipientMax = 6

type recipientState struct {
	query     string
	matchedOn string
	matches   []Contact
	contacts  *Contact
	count     int
	added     []string
	hi        int
	dismissed string
}

type recipientKey struct{}

// RecipientInput validates addresses and offers contact suggestions.
func RecipientInput(c *ui.Context, recipients *[]string, opts RecipientInputOptions) input.TagsInputParts {
	k := core.Tokens(c)
	if opts.Label == "" {
		opts.Label = core.Msg(c, "recipient_input.label")
	}
	wrap := ui.Column(c).Disabled(opts.Disabled)
	st := ui.Local(wrap, recipientKey{}, func() *recipientState { return &recipientState{hi: -1} })
	s := *st
	// Suggestions are found again only when the typed text changes.
	var first *Contact
	if len(opts.Contacts) > 0 {
		first = &opts.Contacts[0]
	}
	if q := strings.ToLower(strings.TrimSpace(s.query)); q != s.matchedOn || first != s.contacts || len(opts.Contacts) != s.count || !slices.Equal(s.added, *recipients) {
		s.matchedOn, s.matches, s.hi = q, recipientFind(opts.Contacts, *recipients, q), -1
		s.contacts, s.count, s.added = first, len(opts.Contacts), slices.Clone(*recipients)
	}
	typed := s.query
	var parts input.TagsInputParts
	wrap.Children(func() {
		parts = input.TagsInput(c, recipients, input.TagsInputOptions{
			Placeholder: opts.Placeholder, Label: opts.Label, Disabled: opts.Disabled, Query: &s.query,
			Validate: func(tag string) string {
				if s.hi >= 0 && tag == strings.TrimSpace(typed) {
					return "" // replaced by the highlighted contact below
				}
				if _, err := mail.ParseAddress(tag); err != nil {
					return fmt.Sprintf(core.Msg(c, "recipient_input.invalid"), tag)
				}
				return ""
			},
		})
	})
	parts.Element = wrap
	in := parts.Input
	live := !wrap.IsDisabled()
	showing := live && len(s.matches) > 0 && s.dismissed != s.query && s.query != ""
	if showing && in.Focused() {
		switch {
		case in.Shortcut(0, ui.KeyDown):
			s.hi = min(len(s.matches)-1, s.hi+1)
		case in.Shortcut(0, ui.KeyUp):
			s.hi = max(-1, s.hi-1)
		case in.Shortcut(0, ui.KeyEscape):
			s.dismissed, s.hi = s.query, -1
		}
	}
	pick := func(ct Contact) {
		addr := recipientFormat(ct)
		if n := len(*recipients); n > 0 && (*recipients)[n-1] == strings.TrimSpace(typed) {
			*recipients = (*recipients)[:n-1] // the raw text Enter added
		}
		if !slices.Contains(*recipients, addr) {
			*recipients = append(slices.Clone(*recipients), addr)
		}
		s.query, s.hi = "", -1
		in.Focus()
		slog.Debug("mujica: recipient_input pick", "count", len(*recipients))
	}
	if live && in.Submitted() && s.hi >= 0 && s.hi < len(s.matches) {
		pick(s.matches[s.hi])
		parts = input.NewTagsInputParts(wrap, in, true)
	}
	if showing {
		b := wrap.Bounds()
		ui.Overlay(c, func() {
			ww, wh := c.Size()
			w := min(max(240, b.W), max(0, ww-8))
			h := min(float32(len(s.matches))*core.Density(c).ControlHeight()+16, max(0, wh-8))
			x, y := core.PlaceOverlay(c, b, w, h, 4)
			panel := ui.Column(c).Absolute().Left(x).Top(y).Width(w).MaxHeight(h).Padding(4).Role(ui.RoleList).
				Label(core.Msg(c, "recipient_input.suggestions"))
			core.OverlayStyle(c, panel, core.OverlayEnter(c, panel))
			panel.Children(func() {
				ui.Scroll(c).Shrink(1).Padding(4).Children(func() {
					for i, ct := range s.matches {
						row := ui.ButtonBase(c.Key(ct.Email)).MinHeight(core.Density(c).ControlHeight()).Justify(ui.Start).Gap(8).Padding(0, 8).
							Radius(theme.ControlRadius).Role(ui.RoleListItem).Label(recipientFormat(ct)).FocusRing(false)
						if i == s.hi || row.Hovered() {
							row.Background(k.Selection)
						}
						if i == s.hi || row.FocusVisible() {
							row.DrawOver(func(p *ui.Painter, r ui.Rect) { internal.FocusRing(p, r, theme.ControlRadius, k.Focus) })
						}
						row.Children(func() {
							display.Avatar(c, ct.Name, display.AvatarOptions{Size: 22})
							ui.Text(c, ct.Name).SingleLine()
							ui.Text(c, ct.Email).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted).SingleLine()
						})
						if row.Clicked() && !row.IsDisabled() {
							pick(ct)
							parts = input.NewTagsInputParts(parts.Element, parts.Input, true)
						}
					}
				})
			})
			if !in.Focused() && !panel.FocusWithin() {
				panel.Invisible()
			}
		})
	}
	return parts
}

func recipientFind(contacts []Contact, added []string, q string) []Contact {
	if q == "" {
		return nil
	}
	var out []Contact
	for _, ct := range contacts {
		if len(out) == recipientMax {
			break
		}
		if (strings.Contains(strings.ToLower(ct.Name), q) || strings.Contains(strings.ToLower(ct.Email), q)) &&
			!slices.Contains(added, recipientFormat(ct)) {
			out = append(out, ct)
		}
	}
	return out
}

func recipientFormat(ct Contact) string {
	if ct.Name == "" {
		return ct.Email
	}
	plain := ct.Name + " <" + ct.Email + ">"
	if _, err := mail.ParseAddress(plain); err == nil {
		return plain
	}
	address := mail.Address{Name: ct.Name, Address: ct.Email}
	return address.String()
}
