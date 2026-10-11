package feedback

import (
	"slices"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
)

func toastIDs(q *toastQueue, visible bool) []string {
	var ids []string
	for _, it := range q.items {
		if it.visible == visible {
			ids = append(ids, it.msg.ID)
		}
	}
	return ids
}

func TestToastQueueLimitsVisible(t *testing.T) {
	var q toastQueue
	t0 := time.Unix(0, 0)
	for _, id := range []string{"a", "b", "c", "d", "e"} {
		q.push(ToastMessage{ID: id, Message: id}, t0)
	}
	if v, w := toastIDs(&q, true), toastIDs(&q, false); !slices.Equal(v, []string{"a", "b", "c"}) || !slices.Equal(w, []string{"d", "e"}) {
		t.Fatalf("visible %v waiting %v", v, w)
	}
	// d starts its 5 s only when it shows.
	q.expire(t0.Add(ToastDuration))
	if next := q.settle(t0.Add(ToastDuration), nil); next != ToastDuration {
		t.Fatalf("next wake %v", next)
	}
	if v := toastIDs(&q, true); !slices.Equal(v, []string{"d", "e"}) {
		t.Fatalf("after expiry visible %v", v)
	}
	q.expire(t0.Add(2*ToastDuration - time.Millisecond))
	if len(q.items) != 2 {
		t.Fatalf("expired early: %v", toastIDs(&q, true))
	}
	q.expire(t0.Add(2 * ToastDuration))
	if len(q.items) != 0 {
		t.Fatalf("left %v", toastIDs(&q, true))
	}
}

func TestToastQueueDedupesByID(t *testing.T) {
	var q toastQueue
	t0 := time.Unix(0, 0)
	q.push(ToastMessage{ID: "sync", Message: "Syncing 1 file"}, t0)
	q.push(ToastMessage{Message: "Saved"}, t0)
	q.push(ToastMessage{ID: "sync", Message: "Syncing 2 files"}, t0.Add(4*time.Second))
	q.push(ToastMessage{Message: "Saved"}, t0.Add(4*time.Second))
	if len(q.items) != 2 || q.items[0].msg.Message != "Syncing 2 files" || q.items[1].msg.ID != "Saved" {
		t.Fatalf("items %+v", q.items)
	}
	// The replaced toast restarted its time.
	q.expire(t0.Add(6 * time.Second))
	if len(q.items) != 2 {
		t.Fatalf("replacement kept the old time: %d left", len(q.items))
	}
}

func TestToastQueueErrorsStayAndPauseHolds(t *testing.T) {
	var q toastQueue
	t0 := time.Unix(0, 0)
	q.push(ToastMessage{ID: "e", Severity: core.SeverityError, Message: "Failed"}, t0)
	q.push(ToastMessage{ID: "w", Severity: core.SeverityWarning, Message: "Slow", Action: "Details"}, t0)
	q.settle(t0.Add(4*time.Second), map[string]bool{"w": true})
	q.expire(t0.Add(time.Hour))
	if next := q.settle(t0.Add(time.Hour), map[string]bool{"w": true}); next != 0 {
		t.Fatalf("paused toast still counts down: %v", next)
	}
	if !slices.Equal(toastIDs(&q, true), []string{"e", "w"}) {
		t.Fatalf("visible %v", toastIDs(&q, true))
	}
	// Released: the wake-up is computed after resuming, in the same call.
	if next := q.settle(t0.Add(time.Hour), nil); next != time.Second {
		t.Fatalf("resumed toast left %v, want 1s", next)
	}
	q.expire(t0.Add(time.Hour + time.Second))
	if !slices.Equal(toastIDs(&q, true), []string{"e"}) {
		t.Fatalf("visible %v", toastIDs(&q, true))
	}
}

func TestToastQueuePanics(t *testing.T) {
	for _, m := range []ToastMessage{{}, {Message: "x", Duration: -1}, {Message: "x", Severity: 9}} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("%+v: no panic", m)
				}
			}()
			var q toastQueue
			q.push(m, time.Now())
		}()
	}
}

func toastTester(m ToastMessage, res *ToasterResult) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if input.Button(c, "Notify", input.ButtonOptions{}).Clicked() {
			ShowToast(c, m)
		}
		if r := Toaster(c); r != (ToasterResult{}) {
			*res = r
		}
	}, 600, 400)
}

func TestToastActionAndClose(t *testing.T) {
	var res ToasterResult
	tt := toastTester(ToastMessage{ID: "del", Message: "Note deleted", Action: "Undo"}, &res)
	if err := tt.Click("Notify"); err != nil {
		t.Fatal(err)
	}
	if !tt.HasText("Note deleted") {
		t.Fatal("toast not shown")
	}
	if err := tt.Click("Undo"); err != nil {
		t.Fatal(err)
	}
	if res.Action != "del" || tt.HasText("Note deleted") {
		t.Fatalf("after undo %+v, still shown %v", res, tt.HasText("Note deleted"))
	}
	res = ToasterResult{}
	_ = tt.Click("Notify")
	if err := tt.Click("Dismiss notification"); err != nil {
		t.Fatal(err)
	}
	if res.Closed != "del" || tt.HasText("Note deleted") {
		t.Fatalf("after close %+v", res)
	}
}

func TestToastKeepsFocusAndEscapeCloses(t *testing.T) {
	var res ToasterResult
	tt := toastTester(ToastMessage{ID: "s", Message: "Saved"}, &res)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	if !tt.HasText("Saved") || !tt.Focused("Notify") {
		t.Fatalf("shown %v, focus kept %v", tt.HasText("Saved"), tt.Focused("Notify"))
	}
	if a := tt.Announcements(); !slices.Contains(a, "Saved") {
		t.Fatalf("announcements %q", a)
	}
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Dismiss notification") {
		t.Fatal("Tab does not reach the toast")
	}
	tt.Key(0, ui.KeyEscape)
	if res.Closed != "s" || tt.HasText("Saved") {
		t.Fatalf("Escape: %+v", res)
	}
}

func TestToastExpiresAndHoverPauses(t *testing.T) {
	var res ToasterResult
	const life = 40 * time.Millisecond
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if input.Button(c, "Notify", input.ButtonOptions{}).Clicked() {
			ShowToast(c, ToastMessage{Message: "Copied", Duration: life})
		}
		res = Toaster(c)
	}, 600, 400)
	_ = tt.Click("Notify")
	if res.wake.IsZero() {
		t.Fatal("no wake-up scheduled for a running toast")
	}
	r, _ := tt.Find("Copied")
	tt.Move(r.X+r.W/2, r.Y+r.H/2)
	if !res.wake.IsZero() {
		t.Fatal("hovered toast still schedules its expiry")
	}
	time.Sleep(2 * life)
	released := time.Now()
	tt.Move(5, 395)
	done := time.Now()
	// The frame that ends the hover must schedule the expiry itself.
	if res.wake.IsZero() || res.wake.Before(released) || res.wake.After(done.Add(life)) {
		t.Fatalf("release scheduled wake %v, released %v", res.wake, released)
	}
	if !tt.HasText("Copied") {
		t.Fatal("hovered toast expired")
	}
	// Stands in for the frame MyGo's timer runs at res.wake.
	time.Sleep(time.Until(res.wake) + 5*time.Millisecond)
	tt.Frame()
	if tt.HasText("Copied") {
		t.Fatal("toast stayed after its wake-up")
	}
}

func TestToastSnapshot(t *testing.T) {
	snapshot(t, "toast", 480, 300, func(c *ui.Context) {
		if len(toasts(c).items) == 0 {
			ShowToast(c, ToastMessage{ID: "1", Message: "Settings saved", Severity: core.SeveritySuccess})
			ShowToast(c, ToastMessage{ID: "2", Message: "Note moved to trash", Action: "Undo"})
			ShowToast(c, ToastMessage{ID: "3", Message: "Upload failed: network lost", Severity: core.SeverityError})
			ShowToast(c, ToastMessage{ID: "4", Message: "Queued and not shown"})
		}
		Toaster(c)
	})
}
