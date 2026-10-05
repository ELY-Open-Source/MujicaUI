package datetime

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

var calAttendees = []Attendee{
	{Name: "Sakiko Togawa", Email: "sakiko@court.example", Status: AttendeeAccepted, Organizer: true},
	{Name: "Mutsumi Wakaba", Status: AttendeeTentative},
	{Name: "Uika Misumi", Email: "uika@court.example", Status: AttendeeDeclined, Optional: true},
	{Name: "Nyamu Yuutenji", Status: AttendeePending},
}

func attendeeTester(list *[]Attendee, opts AttendeeListOptions) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if i, ok := AttendeeList(c, *list, opts).Removed(); ok {
			*list = append((*list)[:i:i], (*list)[i+1:]...)
		}
	}, 480, 400)
}

func TestAttendeeListStatusesAndSummary(t *testing.T) {
	list := append([]Attendee(nil), calAttendees...)
	tt := attendeeTester(&list, AttendeeListOptions{})
	for _, s := range []string{"1 accepted · 1 maybe · 1 declined · 1 awaiting", "Sakiko Togawa · Organizer",
		"Uika Misumi · Optional", "Accepted", "Maybe", "Declined", "Awaiting reply", "uika@court.example"} {
		if !tt.HasText(s) {
			t.Fatalf("missing %q in %v", s, tt.Texts())
		}
	}
	if tt.HasText("Remove Mutsumi Wakaba") {
		t.Fatal("remove buttons without Removable")
	}
}

func TestAttendeeListRemoveKeyboardAndOrganizer(t *testing.T) {
	list := append([]Attendee(nil), calAttendees...)
	tt := attendeeTester(&list, AttendeeListOptions{Removable: true})
	if tt.HasText("Remove Sakiko Togawa") {
		t.Fatal("the organizer cannot be removed")
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	if len(list) != 3 || list[0].Name != "Sakiko Togawa" || list[1].Name != "Uika Misumi" {
		t.Fatalf("keyboard remove left %v", list)
	}
	_ = tt.Click("Remove Uika Misumi")
	if len(list) != 2 || !tt.HasText("1 accepted · 0 maybe · 0 declined · 1 awaiting") {
		t.Fatalf("click remove left %v", list)
	}
}

func TestAttendeeListDisabledEmptyInvalid(t *testing.T) {
	list := append([]Attendee(nil), calAttendees...)
	tt := attendeeTester(&list, AttendeeListOptions{Removable: true, Disabled: true})
	_ = tt.Click("Remove Mutsumi Wakaba")
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	if len(list) != 4 {
		t.Fatal("disabled list removed someone")
	}
	empty := []Attendee{}
	if tt := attendeeTester(&empty, AttendeeListOptions{}); !tt.HasText("No attendees yet.") {
		t.Fatal("no empty state")
	}
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic on bad status")
		}
	}()
	bad := []Attendee{{Name: "X", Status: 7}}
	attendeeTester(&bad, AttendeeListOptions{})
}

func TestAttendeeListSnapshot(t *testing.T) {
	snapshot(t, "attendee_list", 420, 260, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() { AttendeeList(c, calAttendees, AttendeeListOptions{Removable: true}) })
	})
}
