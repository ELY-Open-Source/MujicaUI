package agent

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/icons"
)

// ArtifactVersion is one saved version of an artifact.
type ArtifactVersion struct {
	Label string
	Note  string
}

// ArtifactVersionSwitcherOptions configure an ArtifactVersionSwitcher.
type ArtifactVersionSwitcherOptions struct {
	Disabled bool
}

// ArtifactVersionResult reports this frame's events.
type ArtifactVersionResult struct {
	Element ui.Element
	// Changed reports *current moved.
	Changed bool
	// Compare asks to diff current against the previous version.
	Compare bool
}

// ArtifactVersionSwitcher steps between artifact versions.
func ArtifactVersionSwitcher(c *ui.Context, current *int, versions []ArtifactVersion, opts ArtifactVersionSwitcherOptions) ArtifactVersionResult {
	if current == nil || len(versions) == 0 || *current < 0 || *current >= len(versions) {
		panic("mujica: ArtifactVersionSwitcher needs versions and a current index inside them")
	}
	for i, v := range versions {
		if v.Label == "" {
			panic(fmt.Sprintf("mujica: ArtifactVersionSwitcher version %d needs a Label", i))
		}
	}
	u := core.Density(c).Unit()
	n := len(versions)
	var r ArtifactVersionResult
	r.Element = ui.Row(c).Shrink(0).Gap(u * 2).AlignItems(ui.Center).Disabled(opts.Disabled).
		Label(core.Msg(c, "artifact_version_switcher.label"))
	locked := r.Element.IsDisabled()
	move := func(d int) {
		if to := *current + d; !locked && to >= 0 && to < n {
			*current, r.Changed = to, true
		}
	}
	r.Element.Children(func() {
		if core.IconAction(c, icons.Must("chevron-left"), core.Msg(c, "artifact_version_switcher.older")).Disabled(locked || *current == 0).Clicked() {
			move(-1)
		}
		v := versions[*current]
		ui.Column(c).Shrink(1).Role(ui.RoleStatus).Label(v.Label).Children(func() {
			ui.Row(c).Gap(u * 2).AlignItems(ui.Center).Children(func() {
				ui.Text(c, v.Label).FontWeight(600).SingleLine()
				AgentCaption(c, fmt.Sprintf(core.Msg(c, "artifact_version_switcher.count"), *current+1, n))
			})
			if v.Note != "" {
				AgentCaption(c, v.Note).SingleLine()
			}
		})
		if core.IconAction(c, icons.Must("chevron-right"), core.Msg(c, "artifact_version_switcher.newer")).Disabled(locked || *current == n-1).Clicked() {
			move(1)
		}
		if *current > 0 {
			label := fmt.Sprintf(core.Msg(c, "artifact_version_switcher.compare"), versions[*current-1].Label)
			if AgentToggle(c, label, "arrow-up-down", locked) {
				r.Compare = true
			}
		}
	})
	if !locked {
		if r.Element.Shortcut(0, ui.KeyLeft) {
			move(-1)
		}
		if r.Element.Shortcut(0, ui.KeyRight) {
			move(1)
		}
	}
	if r.Changed || r.Compare {
		slog.Debug("mujica: artifact_version_switcher event", "current", *current, "compare", r.Compare)
	}
	return r
}
