package mujica

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
	"mujicaui/theme"
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
func CheckpointList(c *ui.Context, current *string, points []Checkpoint, opts CheckpointListOptions) ChoiceResult {
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
	k := Tokens(c)
	var r ChoiceResult
	r.Element = ui.Column(c).Shrink(0).Disabled(opts.Disabled).Label(msg(c, "checkpoint_list.label"))
	st := ui.Local(r.Element, "mujica.checkpoints", func() checkpointListState { return checkpointListState{} })
	locked := r.Element.IsDisabled()
	if locked {
		st.asking = false // a disabled list closes its confirmation
	}
	r.Element.Children(func() {
		if len(points) == 0 {
			agentCaption(c, msg(c, "checkpoint_list.empty"))
			return
		}
		agentRows(c, len(points), func(i int) {
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
		res := AlertDialog(c, &st.asking, AlertDialogOptions{Title: fmt.Sprintf(msg(c, "checkpoint_list.confirm"), label),
			Description: msg(c, "checkpoint_list.warning"), ConfirmLabel: msg(c, "checkpoint_list.rollback")})
		if res.Confirmed && !locked && seen[st.target] {
			*current, r.changed = st.target, true
			slog.Debug("mujica: checkpoint_list rollback", "to", st.target)
		}
	})
	return r
}

// checkpointRow shows one checkpoint and reports a rollback request.
func checkpointRow(c *ui.Context, p Checkpoint, current, locked bool) bool {
	k := Tokens(c)
	d := Density(c)
	u := d.Unit()
	asked := false
	row := ui.Row(c).Key(p.ID).Gap(u*3).AlignItems(ui.Center).MinHeight(d.ControlHeight()+u*2).Padding(u*2, u*3)
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
				agentCaption(c, p.Note)
			}
		})
		if current {
			agentCaption(c, msg(c, "checkpoint_list.current"))
			return
		}
		asked = agentToggle(c, msg(c, "checkpoint_list.rollback"), "refresh-cw", locked)
	})
	row.Label(p.Label)
	return asked
}
