package files

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/code"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/icons"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// TransferStatus is where a transfer is.
type TransferStatus int

const (
	TransferQueued TransferStatus = iota
	TransferActive
	TransferPaused
	TransferFailed
	TransferDone
)

// QueuedTransfer is one upload or download the caller runs.
type QueuedTransfer struct {
	ID, Name    string
	Upload      bool
	Status      TransferStatus
	Done, Total int64
	// Error explains a failure.
	Error string
}

// TransferAction is what the user asked of a transfer.
type TransferAction int

const (
	TransferNone TransferAction = iota
	TransferPause
	TransferResume
	TransferRetry
	TransferRemove
)

// TransferQueueOptions configure a TransferQueue.
type TransferQueueOptions struct {
	Disabled bool
	// Revision changes when items are edited in place.
	Revision uint64
}

// TransferQueueResult exposes the queue and requested item action.
type TransferQueueResult struct {
	Element *ui.Element
	Action  TransferAction
	// ID is the transfer the action is for.
	ID string
}

// TransferQueueState is the queue's place.
type TransferQueueState struct {
	rows           ui.ListState
	data           *QueuedTransfer
	n              int
	revision       uint64
	active, failed int
	started        bool
}

// TransferQueue virtualizes uploads, downloads and their available actions.
func TransferQueue(c *ui.Context, items []QueuedTransfer, s *TransferQueueState, opts TransferQueueOptions) TransferQueueResult {
	if s == nil {
		panic("mujica: TransferQueue needs a state")
	}
	var first *QueuedTransfer
	if len(items) > 0 {
		first = &items[0]
	}
	if !s.started || s.data != first || s.n != len(items) || s.revision != opts.Revision {
		s.active, s.failed = 0, 0
		ids := make(map[string]bool, len(items))
		for _, t := range items {
			if t.ID == "" || ids[t.ID] {
				panic("mujica: Transfer needs a unique ID")
			}
			ids[t.ID] = true
			if t.Status < TransferQueued || t.Status > TransferDone || t.Done < 0 || t.Total < 0 || t.Total > 0 && t.Done > t.Total {
				panic("mujica: Transfer " + t.ID + " has an invalid status or size")
			}
			switch t.Status {
			case TransferActive, TransferQueued:
				s.active++
			case TransferFailed:
				s.failed++
			}
		}
		s.data, s.n, s.revision, s.started = first, len(items), opts.Revision, true
	}
	d := core.Density(c)
	var res TransferQueueResult
	res.Element = code.CodePanel(c).Disabled(opts.Disabled)
	disabled := res.Element.IsDisabled()
	res.Element.Children(func() {
		code.CodeHeader(c, func() {
			ui.Text(c, fmt.Sprintf(core.Msg(c, "transfer_queue.summary"), s.active, s.failed)).Role(ui.RoleStatus)
		})
		if len(items) == 0 {
			code.CodeCaption(c, core.Msg(c, "transfer_queue.empty")).Padding(d.Unit() * 4)
			return
		}
		s.rows.Key = func(i int) any { return items[i].ID }
		s.rows.Label = func(i int) string { return items[i].Name }
		ui.List(c, &s.rows, len(items), func(i int) {
			if a := transferRow(c, items[i]); a != TransferNone && !disabled {
				res.Action, res.ID = a, items[i].ID
				slog.Debug("mujica: transfer_queue action", "action", a, "id", items[i].ID)
			}
		}).Grow(1).MinHeight(d.ControlHeight() * 3)
	})
	return res
}

// transferRow builds one transfer and returns the action clicked.
func transferRow(c *ui.Context, t QueuedTransfer) TransferAction {
	k := core.Tokens(c)
	d := core.Density(c)
	u := d.Unit()
	action := TransferNone
	ui.Row(c).MinHeight(d.ControlHeight()+2*u).Padding(u, u*3).Gap(u * 3).AlignItems(ui.Center).Children(func() {
		FileIcon(c, t.Name, FileIconOptions{})
		ui.Column(c).Grow(1).MinWidth(0).Gap(u).Children(func() {
			ui.Row(c).Gap(u).AlignItems(ui.Center).Children(func() {
				dir, key := "arrow-down", "transfer_queue.download"
				if t.Upload {
					dir, key = "arrow-up", "transfer_queue.upload"
				}
				ui.Icon(c, icons.Must(dir)).Size(14, 14).Shrink(0).TextColor(k.TextMuted).Label(core.Msg(c, key)).Role(ui.RoleImage)
				ui.Text(c, t.Name).SingleLine().MinWidth(0)
			})
			frac := float32(0)
			if t.Total > 0 {
				frac = float32(t.Done) / float32(t.Total)
			}
			if t.Status == TransferDone {
				frac = 1
			}
			ui.Box(c).Height(4).Shrink(0).Role(ui.RoleNone).Draw(func(p *ui.Painter, r ui.Rect) {
				p.Fill(r, k.Border, 2)
				fill := k.AccentText
				if t.Status == TransferFailed {
					fill = k.Danger
				}
				if frac > 0 {
					p.Fill(ui.Rect{X: internal.Snap(p, r.X), Y: internal.Snap(p, r.Y), W: internal.Snap(p, max(r.W*frac, r.H)), H: r.H}, fill, 2)
				}
			})
			transferStatusText(c, t, frac)
		})
		var acts []TransferAction
		switch t.Status {
		case TransferActive, TransferQueued:
			acts = []TransferAction{TransferPause}
		case TransferPaused:
			acts = []TransferAction{TransferResume}
		case TransferFailed:
			acts = []TransferAction{TransferRetry}
		}
		acts = append(acts, TransferRemove)
		for _, a := range acts {
			icon, key := transferActionIcon(a)
			if core.IconAction(c, icon, fmt.Sprintf(core.Msg(c, key), t.Name)).Clicked() {
				action = a
			}
		}
	})
	return action
}

func transferStatusText(c *ui.Context, t QueuedTransfer, frac float32) {
	k := core.Tokens(c)
	text, col := "", k.TextMuted
	switch t.Status {
	case TransferQueued:
		text = core.Msg(c, "transfer_queue.queued")
	case TransferActive:
		text = fmt.Sprintf("%.0f%% · %s / %s", frac*100, code.ProcessListBytes(uint64(t.Done)), code.ProcessListBytes(uint64(t.Total)))
	case TransferPaused:
		text = fmt.Sprintf(core.Msg(c, "transfer_queue.paused"), frac*100)
	case TransferFailed:
		text, col = fmt.Sprintf(core.Msg(c, "transfer_queue.failed"), t.Error), k.Danger
	case TransferDone:
		text, col = core.Msg(c, "transfer_queue.done"), k.Success
	}
	ui.Text(c, text).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(col).SingleLine()
}

var transferPauseIcon = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><path d="M9 5v14M15 5v14"/></svg>`))
var transferPlayIcon = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linejoin="round"><path d="M8 5l11 7-11 7z"/></svg>`))

func transferActionIcon(a TransferAction) (*ui.SVG, string) {
	switch a {
	case TransferPause:
		return transferPauseIcon, "transfer_queue.pause"
	case TransferResume:
		return transferPlayIcon, "transfer_queue.resume"
	case TransferRetry:
		return icons.Must("refresh-cw"), "transfer_queue.retry"
	}
	return icons.Must("x"), "transfer_queue.remove"
}
