package chat

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/agent"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/icons"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// ArtifactCardOptions configure an ArtifactCard. Title is required.
type ArtifactCardOptions struct {
	Kind  agent.ArtifactKind
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
func ArtifactCard(c *ui.Context, opts ArtifactCardOptions) core.ChoiceResult {
	if opts.Title == "" {
		panic("mujica: ArtifactCard needs a Title")
	}
	if opts.Kind < agent.ArtifactCode || opts.Kind > agent.ArtifactWeb {
		panic("mujica: ArtifactCard Kind out of range")
	}
	k := core.Tokens(c)
	d := core.Density(c)
	kind := artifactKinds[opts.Kind]
	var b *ui.Element
	ui.Row(c).Shrink(0).Disabled(opts.Disabled).Children(func() {
		b = ui.ButtonBase(c).MinWidth(260).MaxWidth(420).MinHeight(56).Padding(d.Unit()*2, d.Unit()*3).Gap(d.Unit()*3).
			Justify(ui.Start).Radius(theme.CardRadius).Background(k.Surface).Border(theme.BorderWidth, k.ControlBorder).
			Label(core.Msg(c, kind.name) + ": " + opts.Title)
	})
	internal.Ring(b, k.Focus, theme.CardRadius)
	if b.Hovered() {
		b.Background(k.SurfaceHover)
	}
	b.Transition(core.ColorFade)
	b.Children(func() {
		ui.Box(c).Size(36, 36).Shrink(0).Center().Radius(theme.ControlRadius).Background(k.Selection).Children(func() {
			ui.Icon(c, icons.Must(kind.icon)).Size(18, 18).TextColor(k.AccentText)
		})
		ui.Column(c).Grow(1).MinWidth(0).Gap(2).Children(func() {
			ui.Text(c, opts.Title).SingleLine().Bold().TextColor(k.Text)
			ui.Row(c).Gap(d.Unit()).AlignItems(ui.Center).Children(func() {
				detail := core.Msg(c, kind.name)
				if opts.Detail != "" {
					detail += " · " + opts.Detail
				}
				ui.Text(c, detail).SingleLine().FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted)
				if opts.Streaming {
					core.Spinner(c, core.SpinnerOptions{Size: 12, Label: core.Msg(c, "artifact_card.writing")})
				}
			})
		})
		ui.Text(c, core.Msg(c, "artifact_card.open")).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.AccentText).Shrink(0)
		ui.Icon(c, icons.Must("chevron-right")).Size(14, 14).TextColor(k.AccentText)
	})
	open := b.Clicked() && !b.IsDisabled()
	if open {
		slog.Debug("mujica: artifact_card open", "kind", opts.Kind)
	}
	return core.NewChoiceResult(b, open)
}
