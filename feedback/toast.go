package feedback

import (
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/icons"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// Toast timing and stacking per window.
const (
	ToastDuration = 5 * time.Second
	ToastVisible  = 3
)

// ToastMessage is one toast. Message is required.
type ToastMessage struct {
	// ID dedupes: showing an ID already queued replaces it. Empty uses Message.
	ID       string
	Severity core.Severity
	Message  string
	// Action labels a button; Toaster reports its click by ID.
	Action string
	// Duration is ToastDuration when zero; errors then stay until closed.
	Duration time.Duration
}

// ToasterResult reports the toast events of a frame by message ID.
type ToasterResult struct {
	Action string
	Closed string
	wake   time.Time // when this frame asked for the next one; zero for none
}

type toastItem struct {
	msg       ToastMessage
	life      time.Duration // zero stays until closed
	spent     time.Duration
	runFrom   time.Time // zero while paused or queued
	visible   bool
	announced bool
}

// toastQueue holds a window's toasts, visible ones first in order shown.
type toastQueue struct {
	items []toastItem
}

type toastKey struct{}

func toastLife(m ToastMessage) time.Duration {
	switch {
	case m.Duration < 0:
		panic("mujica: negative toast Duration")
	case m.Duration > 0:
		return m.Duration
	case m.Severity == core.SeverityError:
		return 0
	}
	return ToastDuration
}

func (q *toastQueue) push(m ToastMessage, now time.Time) {
	if m.Message == "" {
		panic("mujica: toast needs a Message")
	}
	_ = m.Severity.String() // panics on an unknown Severity
	if m.ID == "" {
		m.ID = m.Message
	}
	for i := range q.items {
		it := &q.items[i]
		if it.msg.ID != m.ID {
			continue
		}
		it.msg, it.life, it.spent, it.announced = m, toastLife(m), 0, false
		if it.visible && !it.runFrom.IsZero() {
			it.runFrom = now
		}
		slog.Debug("mujica: toast replaced", "id", m.ID)
		return
	}
	q.items = append(q.items, toastItem{msg: m, life: toastLife(m)})
	slog.Debug("mujica: toast queued", "id", m.ID)
	q.promote(now)
}

// promote shows queued toasts while fewer than ToastVisible show.
func (q *toastQueue) promote(now time.Time) {
	n := 0
	for i := range q.items {
		it := &q.items[i]
		if !it.visible && n < ToastVisible {
			it.visible, it.runFrom = true, now
			slog.Debug("mujica: toast shown", "id", it.msg.ID)
		}
		if it.visible {
			n++
		}
	}
}

func (it *toastItem) elapsed(now time.Time) time.Duration {
	if it.runFrom.IsZero() {
		return it.spent
	}
	return it.spent + now.Sub(it.runFrom)
}

// expire drops toasts whose time ran out and shows queued ones.
func (q *toastQueue) expire(now time.Time) {
	live := q.items[:0]
	for _, it := range q.items {
		if it.visible && it.life > 0 && it.elapsed(now) >= it.life {
			slog.Debug("mujica: toast expired", "id", it.msg.ID)
			continue
		}
		live = append(live, it)
	}
	q.items = live
	q.promote(now)
}

// settle pauses the toasts held, resumes the rest, and returns the wait until one expires.
func (q *toastQueue) settle(now time.Time, held map[string]bool) (next time.Duration) {
	for _, it := range q.items {
		if it.visible {
			q.pause(it.msg.ID, held[it.msg.ID], now)
		}
	}
	for _, it := range q.items {
		if it.visible && it.life > 0 && !it.runFrom.IsZero() {
			left := it.life - it.elapsed(now)
			if next == 0 || left < next {
				next = left
			}
		}
	}
	return next
}

func (q *toastQueue) pause(id string, paused bool, now time.Time) {
	for i := range q.items {
		it := &q.items[i]
		if it.msg.ID != id || !it.visible {
			continue
		}
		switch {
		case paused && !it.runFrom.IsZero():
			it.spent, it.runFrom = it.elapsed(now), time.Time{}
			slog.Debug("mujica: toast paused", "id", id)
		case !paused && it.runFrom.IsZero():
			it.runFrom = now
			slog.Debug("mujica: toast resumed", "id", id)
		}
	}
}

func (q *toastQueue) remove(id string, now time.Time) {
	for i, it := range q.items {
		if it.msg.ID == id {
			q.items = append(q.items[:i], q.items[i+1:]...)
			break
		}
	}
	q.promote(now)
}

func toasts(c *ui.Context) *toastQueue {
	return ui.Local(c.Root(), toastKey{}, func() toastQueue { return toastQueue{} })
}

// ShowToast queues m in the window. Call it on the UI thread.
func ShowToast(c *ui.Context, m ToastMessage) {
	toasts(c).push(m, c.Now())
	c.Invalidate()
}

// DismissToast closes the toast with id, shown or queued.
func DismissToast(c *ui.Context, id string) {
	toasts(c).remove(id, c.Now())
	c.Invalidate()
}

// Toaster draws the window's toasts, never taking the focus; call it last in the view each frame.
func Toaster(c *ui.Context) ToasterResult {
	q := toasts(c)
	now := c.Now()
	q.expire(now)
	held := map[string]bool{}
	var res ToasterResult
	if len(q.items) == 0 {
		return res
	}
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	w, _ := c.Size()
	ui.Overlay(c, func() {
		stack := ui.Column(c).Absolute().Right(u * 4).Bottom(u * 4).Gap(u * 2).AlignItems(ui.End).
			PassThrough().Role(ui.RoleGroup).Label(core.Msg(c, "toast.region"))
		stack.Children(func() {
			for i := range q.items {
				it := q.items[i]
				if !it.visible {
					continue
				}
				if !it.announced {
					q.items[i].announced = true
					c.Announce(it.msg.Message)
				}
				action, closed, hold := toastCard(c, k, u, min(360, w-u*8), it.msg)
				held[it.msg.ID] = hold
				switch {
				case action:
					res.Action = it.msg.ID
					slog.Debug("mujica: toast action", "id", it.msg.ID)
				case closed:
					res.Closed = it.msg.ID
					slog.Debug("mujica: toast closed", "id", it.msg.ID)
				}
			}
		})
	})
	for _, id := range []string{res.Action, res.Closed} {
		if id != "" {
			q.remove(id, now)
			c.Invalidate()
		}
	}
	// Scheduled after resuming, so a toast released this frame still expires.
	if next := q.settle(now, held); next > 0 {
		res.wake = now.Add(next)
		c.After(next)
	}
	return res
}

// toastCard builds one toast and reports its events and whether hover or focus holds it.
func toastCard(c *ui.Context, k theme.Tokens, u, width float32, m ToastMessage) (action, closed, held bool) {
	tone, icon := core.SeverityLook(k, m.Severity)
	card := ui.Row(c).Key(m.ID).Width(width).AlignItems(ui.Center).Gap(u*3).Padding(u*2, u*2, u*2, u*4).
		Radius(theme.PopupRadius).Background(k.Surface).Border(theme.BorderWidth, k.ControlBorder).
		Role(ui.RoleStatus).Label(m.Message)
	card.Shadow(0, 4, 12, 0, ui.RGBA(0, 0, 0, 0.14))
	card.Transition(ui.ElementTransition{
		Duration: theme.PopupDuration,
		Enter:    &ui.Motion{Y: 8},
		Exit:     &ui.Motion{},
	})
	card.Draw(func(p *ui.Painter, r ui.Rect) {
		p.Fill(ui.Rect{X: r.X, Y: r.Y + theme.PopupRadius, W: 3, H: r.H - 2*theme.PopupRadius}, tone, 0)
	})
	card.Children(func() {
		ui.Icon(c, icon).Size(20, 20).TextColor(tone)
		ui.Text(c, m.Message).TextColor(k.Text).Grow(1).Shrink(1).LineHeight(theme.LineHeight)
		if m.Action != "" && input.Button(c, m.Action, input.ButtonOptions{Variant: input.Ghost}).Clicked() {
			action = true
		}
		x := input.Button(c, "", input.ButtonOptions{Variant: input.Ghost, Icon: icons.Must("x"), Label: core.Msg(c, "toast.close")})
		x.TextColor(k.TextMuted).Padding(0, u)
		if x.Clicked() {
			closed = true
		}
	})
	if card.Shortcut(0, ui.KeyEscape) {
		closed = true
	}
	return action, closed, card.Hovered() || card.FocusWithin()
}
