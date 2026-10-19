package agent

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/internal"
	"github.com/ELY-Open-Source/MujicaUI/overlay"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

// Checkpoint is a saved point to roll back to.
type Checkpoint struct {
	ID    string
	Label string
	Note  string
}

// CheckpointListOptions configure a CheckpointList.
type CheckpointListOptions struct {
	// MaxHeight caps the scrolling list; zero means 360.
	MaxHeight float32
	Disabled  bool
}

type checkpointListState struct {
	asking bool
	target string
}

// CheckpointList marks the current point and confirms rollbacks.
func CheckpointList(c *ui.Context, current *string, points []Checkpoint, opts CheckpointListOptions) core.ChoiceResult {
	if current == nil {
		panic("mujica: CheckpointList needs current")
	}
	if opts.MaxHeight < 0 {
		panic("mujica: CheckpointList MaxHeight is negative")
	}
	seen := map[string]bool{}
	for i, p := range points {
		if p.ID == "" || p.Label == "" || seen[p.ID] {
			panic(fmt.Sprintf("mujica: CheckpointList point %d needs a unique ID and a Label", i))
		}
		seen[p.ID] = true
	}
	maxH := opts.MaxHeight
	if maxH == 0 {
		maxH = 360
	}
	k := core.Tokens(c)
	var r core.ChoiceResult
	r.Element = ui.Column(c).Shrink(0).Disabled(opts.Disabled).Label(core.Msg(c, "checkpoint_list.label"))
	st := ui.Local(r.Element, "mujica.checkpoints", func() checkpointListState { return checkpointListState{} })
	locked := r.Element.IsDisabled()
	if locked {
		st.asking = false // a disabled list closes its confirmation
	}
	r.Element.Children(func() {
		if len(points) == 0 {
			AgentCaption(c, core.Msg(c, "checkpoint_list.empty"))
			return
		}
		AgentRows(c, len(points), func(i int) {
			if checkpointRow(c, points[i], points[i].ID == *current, locked) {
				st.asking, st.target = true, points[i].ID
			}
		}).Shrink(1).MaxHeight(maxH).Background(k.Surface).Border(theme.BorderWidth, k.Border).
			Radius(theme.CardRadius).Dividers(theme.BorderWidth, k.Border)
		label := st.target
		for _, p := range points {
			if p.ID == st.target {
				label = p.Label
			}
		}
		res := overlay.AlertDialog(c, &st.asking, overlay.AlertDialogOptions{Title: fmt.Sprintf(core.Msg(c, "checkpoint_list.confirm"), label),
			Description: core.Msg(c, "checkpoint_list.warning"), ConfirmLabel: core.Msg(c, "checkpoint_list.rollback")})
		if res.Confirmed && !locked && seen[st.target] {
			*current, r = st.target, core.NewChoiceResult(r.Element, true)
			slog.Debug("mujica: checkpoint_list rollback", "to", st.target)
		}
	})
	return r
}

// checkpointRow shows one checkpoint and reports a rollback request.
func checkpointRow(c *ui.Context, p Checkpoint, current, locked bool) bool {
	k := core.Tokens(c)
	d := core.Density(c)
	u := d.Unit()
	asked := false
	row := ui.Row(c.Key(p.ID)).Gap(u*3).AlignItems(ui.Center).MinHeight(d.ControlHeight()+u*2).Padding(u*2, u*3)
	if current {
		row.Background(k.Selection)
	}
	row.Children(func() {
		ui.Box(c).Size(14, 14).Shrink(0).Draw(func(pt *ui.Painter, r ui.Rect) {
			cx, cy := internal.Snap(pt, r.X+r.W/2), internal.Snap(pt, r.Y+r.H/2)
			if current {
				pt.FillPath(internal.Diamond(cx, cy, 12), k.AccentText)
				return
			}
			pt.StrokePath(internal.Diamond(cx, cy, 10), 1.5, k.ControlBorder)
		})
		ui.Column(c).Grow(1).Shrink(1).Gap(u / 2).Children(func() {
			ui.Text(c, p.Label).FontWeight(600).SingleLine()
			if p.Note != "" {
				AgentCaption(c, p.Note)
			}
		})
		if current {
			AgentCaption(c, core.Msg(c, "checkpoint_list.current"))
			return
		}
		asked = AgentToggle(c, core.Msg(c, "checkpoint_list.rollback"), "refresh-cw", locked)
	})
	row.Label(p.Label)
	return asked
}
