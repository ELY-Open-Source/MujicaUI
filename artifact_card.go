package mujica

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/internal"
	"mujicaui/theme"
)

// ArtifactCardOptions configure an ArtifactCard. Title is required.
type ArtifactCardOptions struct {
	Kind  ArtifactKind
	Title string
	// Detail is a short fact, such as "Go, 120 lines".
	Detail string
	// Streaming shows the artifact is still being written.
	Streaming bool
	Disabled  bool
}

var artifactKinds = [3]struct{ icon, name string }{
	{"file", "artifact_card.code"}, {"file-text", "artifact_card.document"}, {"layout-grid", "artifact_card.web"},
}

// ArtifactCard stands for an artifact; activating asks to open it.
func ArtifactCard(c *ui.Context, opts ArtifactCardOptions) ChoiceResult {
	if opts.Title == "" {
		panic("mujica: ArtifactCard needs a Title")
	}
	if opts.Kind < ArtifactCode || opts.Kind > ArtifactWeb {
		panic("mujica: ArtifactCard Kind out of range")
	}
	k := Tokens(c)
	d := Density(c)
	kind := artifactKinds[opts.Kind]
	var b *ui.Element
	ui.Row(c).Shrink(0).Disabled(opts.Disabled).Children(func() {
		b = ui.ButtonBase(c).MinWidth(260).MaxWidth(420).MinHeight(56).Padding(d.Unit()*2, d.Unit()*3).Gap(d.Unit()*3).
			Justify(ui.Start).Radius(theme.CardRadius).Background(k.Surface).Border(theme.BorderWidth, k.ControlBorder).
			Label(msg(c, kind.name) + ": " + opts.Title)
	})
	internal.Ring(b, k.Focus, theme.CardRadius)
	if b.Hovered() {
		b.Background(k.SurfaceHover)
	}
	b.Transition(colorFade)
	b.Children(func() {
		ui.Box(c).Size(36, 36).Shrink(0).Center().Radius(theme.ControlRadius).Background(k.Selection).Children(func() {
			ui.Icon(c, icons.Must(kind.icon)).Size(18, 18).TextColor(k.AccentText)
		})
		ui.Column(c).Grow(1).MinWidth(0).Gap(2).Children(func() {
			ui.Text(c, opts.Title).SingleLine().Bold().TextColor(k.Text)
			ui.Row(c).Gap(d.Unit()).AlignItems(ui.Center).Children(func() {
				detail := msg(c, kind.name)
				if opts.Detail != "" {
					detail += " · " + opts.Detail
				}
				ui.Text(c, detail).SingleLine().FontSize(FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted)
				if opts.Streaming {
					Spinner(c, SpinnerOptions{Size: 12, Label: msg(c, "artifact_card.writing")})
				}
			})
		})
		ui.Text(c, msg(c, "artifact_card.open")).FontSize(FontSize(c, theme.CaptionSize)).TextColor(k.AccentText).Shrink(0)
		ui.Icon(c, icons.Must("chevron-right")).Size(14, 14).TextColor(k.AccentText)
	})
	open := b.Clicked() && !b.IsDisabled()
	if open {
		slog.Debug("mujica: artifact_card open", "kind", opts.Kind)
	}
	return ChoiceResult{Element: b, changed: open}
}
