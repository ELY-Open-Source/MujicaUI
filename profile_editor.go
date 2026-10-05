package mujica

import (
	"log/slog"
	"strings"
	"unicode/utf8"

	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
	"mujicaui/theme"
)

// Profile is what a ProfileEditor edits.
type Profile struct {
	Name, Bio string
}

// ProfileEditorOptions configure a ProfileEditor.
type ProfileEditorOptions struct {
	// Saved is the baseline for detecting unsaved edits.
	Saved Profile
	// Avatar is the current picture; nil shows initials.
	Avatar *ui.Bitmap
	// BioLimit bounds the saved biography; zero means 160.
	BioLimit int
	// Saving shows the save in progress and blocks edits.
	Saving   bool
	Disabled bool
}

// ProfileEditorResult reports save, avatar and discard requests.
type ProfileEditorResult struct {
	Element             *ui.Element
	save, avatar, reset bool
}

// Save reports a request to store the edited profile.
func (r ProfileEditorResult) Save() bool { return r.save }

// ChangeAvatar reports a request to pick a new picture.
func (r ProfileEditorResult) ChangeAvatar() bool { return r.avatar }

// Discarded reports that the edits went back to Saved.
func (r ProfileEditorResult) Discarded() bool { return r.reset }

// ProfileEditor edits name and biography, marking unsaved changes.
func ProfileEditor(c *ui.Context, p *Profile, opts ProfileEditorOptions) ProfileEditorResult {
	if opts.BioLimit < 0 {
		panic("mujica: ProfileEditor BioLimit cannot be negative")
	}
	limit := opts.BioLimit
	if limit == 0 {
		limit = 160
	}
	k := Tokens(c)
	d := Density(c)
	var r ProfileEditorResult
	r.Element = ui.Column(c).Gap(d.Unit() * 4).Width(420).Disabled(opts.Disabled)
	locked := r.Element.IsDisabled() || opts.Saving
	dirty := *p != opts.Saved
	nameErr := ""
	if strings.TrimSpace(p.Name) == "" {
		nameErr = msg(c, "form.required")
	}
	r.Element.Children(func() {
		ui.Row(c).Gap(d.Unit() * 4).AlignItems(ui.Center).Children(func() {
			Avatar(c, p.Name, AvatarOptions{Image: opts.Avatar, Size: 64})
			if Button(c, msg(c, "profile_editor.avatar"), ButtonOptions{Variant: Secondary, Disabled: locked}).Clicked() && !locked {
				r.avatar = true
			}
		})
		FormField(c, msg(c, "profile_editor.name"), FormFieldOptions{Required: true}, func() *ui.Element {
			return TextInput(c, &p.Name, TextInputOptions{Disabled: locked, Error: nameErr}).Input
		})
		FormField(c, msg(c, "profile_editor.bio"), FormFieldOptions{}, func() *ui.Element {
			return TextArea(c, &p.Bio, TextAreaOptions{MaxLength: limit, Disabled: locked}).Input
		})
		ui.Row(c).Gap(d.Unit() * 3).AlignItems(ui.Center).Children(func() {
			if dirty {
				ui.Row(c).Gap(d.Unit()).AlignItems(ui.Center).Grow(1).Role(ui.RoleStatus).Children(func() {
					ui.Box(c).Size(10, 10).Draw(func(pt *ui.Painter, rect ui.Rect) {
						pt.FillPath(internal.Diamond(internal.Snap(pt, rect.X+5), internal.Snap(pt, rect.Y+5), 10), k.Warning)
					})
					ui.Text(c, msg(c, "profile_editor.unsaved")).TextColor(k.Warning).FontSize(FontSize(c, theme.CaptionSize)).FontWeight(600)
				})
			} else {
				ui.Box(c).Grow(1)
			}
			if Button(c, msg(c, "profile_editor.discard"), ButtonOptions{Variant: Ghost, Disabled: locked || !dirty}).Clicked() && !locked {
				*p, r.reset = opts.Saved, true
			}
			valid := strings.TrimSpace(p.Name) != "" && utf8.RuneCountInString(p.Bio) <= limit
			if Button(c, msg(c, "profile_editor.save"), ButtonOptions{Loading: opts.Saving, Disabled: r.Element.IsDisabled() || !dirty || !valid}).Clicked() && !locked && valid {
				r.save = true
			}
		})
	})
	if r.save || r.reset || r.avatar {
		slog.Debug("mujica: profile_editor action", "save", r.save, "discard", r.reset, "avatar", r.avatar)
	}
	return r
}
