package agent

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/icons"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// ArtifactKind is what an artifact holds; zero is code.
type ArtifactKind int

const (
	ArtifactCode ArtifactKind = iota
	ArtifactDocument
	ArtifactWeb
)

// Artifact is something an agent made; the caller owns it.
type Artifact struct {
	Title string
	Kind  ArtifactKind
	// Language highlights code, as "go" or "json".
	Language string
	Content  string
}

// ArtifactPanelOptions configure an ArtifactPanel.
type ArtifactPanelOptions struct {
	// Preview renders documents and web pages; the caller supplies it.
	Preview func()
	// Height of the body; 320 DIPs when zero.
	Height   float32
	Disabled bool
}

// ArtifactPanelResult reports this frame's events.
type ArtifactPanelResult struct {
	Element *ui.Element
	// Copied reports the content went to the clipboard.
	Copied bool
	// Exported asks the caller to export the artifact.
	Exported bool
}

type artifactPanelState struct {
	source   bool
	copiedAt time.Time
}

// ArtifactPanel shows an artifact, with copy and export.
func ArtifactPanel(c *ui.Context, art Artifact, opts ArtifactPanelOptions) ArtifactPanelResult {
	if art.Title == "" {
		panic("mujica: ArtifactPanel needs a Title")
	}
	if art.Kind < ArtifactCode || art.Kind > ArtifactWeb {
		panic(fmt.Sprintf("mujica: unknown ArtifactKind %d", int(art.Kind)))
	}
	if art.Kind != ArtifactCode && opts.Preview == nil {
		panic("mujica: ArtifactPanel needs a Preview for documents and web pages")
	}
	if opts.Height < 0 {
		panic("mujica: ArtifactPanel Height is negative")
	}
	h := opts.Height
	if h == 0 {
		h = 320
	}
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	var r ArtifactPanelResult
	r.Element = AgentCard(c).Disabled(opts.Disabled).Label(art.Title)
	st := ui.Local(r.Element, "mujica.artifact", func() artifactPanelState { return artifactPanelState{source: art.Kind == ArtifactCode} })
	locked := r.Element.IsDisabled()
	r.Element.Children(func() {
		ui.Row(c).Gap(u * 2).AlignItems(ui.Center).Children(func() {
			ui.Text(c, art.Title).Font(theme.SerifFont).FontSize(core.FontSize(c, theme.H3Size)).FontWeight(600).SingleLine().Grow(1).Shrink(1)
			if opts.Preview != nil {
				artifactViewSwitch(c, &st.source, locked)
			}
			copied := !st.copiedAt.IsZero() && c.Now().Sub(st.copiedAt) < 2*time.Second
			label, icon := core.Msg(c, "artifact_panel.copy"), "copy"
			if copied {
				label, icon = core.Msg(c, "artifact_panel.copied"), "check"
				c.After(2*time.Second - c.Now().Sub(st.copiedAt))
			}
			if AgentToggle(c, label, icon, locked) {
				c.WriteClipboard(art.Content)
				st.copiedAt, r.Copied = c.Now(), true
			}
			if AgentToggle(c, core.Msg(c, "artifact_panel.export"), "upload", locked) {
				r.Exported = true
			}
		})
		ui.Column(c).Height(h).Shrink(0).Children(func() {
			if st.source || opts.Preview == nil {
				AgentCode(c, art.Language, art.Content, h).Grow(1)
				return
			}
			ui.Scroll(c).Grow(1).Padding(u*3).Radius(theme.ControlRadius).Background(k.Background).
				Border(theme.BorderWidth, k.Border).Children(opts.Preview)
		})
	})
	if r.Copied || r.Exported {
		slog.Debug("mujica: artifact_panel event", "title", art.Title, "copied", r.Copied, "exported", r.Exported)
	}
	return r
}

// artifactViewSwitch is the Preview / Source segmented switch.
func artifactViewSwitch(c *ui.Context, source *bool, locked bool) {
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	ui.Row(c).Shrink(0).FocusGroup(ui.Horizontal).Border(theme.BorderWidth, k.ControlBorder).Radius(theme.ControlRadius).
		Role(ui.RoleRadioGroup).Children(func() {
		for _, src := range []bool{false, true} {
			label, icon := core.Msg(c, "artifact_panel.preview"), "eye"
			if src {
				label, icon = core.Msg(c, "artifact_panel.source"), "file-text"
			}
			b := ui.ButtonBase(c).MinHeight(theme.MinHitSize).Padding(0, u*2).Gap(u).Role(ui.RoleRadio).Label(label).
				Disabled(locked)
			internal.Ring(b, k.Focus, theme.ControlRadius)
			on := *source == src
			if on {
				b.Background(k.Selection)
			} else if b.Hovered() {
				b.Background(k.SurfaceHover)
			}
			b.Children(func() {
				if on {
					ui.Icon(c, icons.Must("check")).Size(12, 12).TextColor(k.AccentText)
				}
				ui.Icon(c, icons.Must(icon)).Size(14, 14)
				ui.Text(c, label).FontSize(core.FontSize(c, theme.CaptionSize+1))
			})
			if !b.IsDisabled() && b.Clicked() && !on {
				*source = src
				slog.Debug("mujica: artifact_panel view", "source", src)
			}
		}
	})
}
