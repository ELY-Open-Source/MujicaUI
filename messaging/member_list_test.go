package messaging

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

func TestMemberListGroupsAndChoose(t *testing.T) {
	var s MemberListState
	members := memberTestList(8)
	opened := ""
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if r := MemberList(c, &s, members, MemberListOptions{}); r.Opened != "" {
			opened = r.Opened
		}
	}, 400, 500)
	if !tt.HasText("Online — 6") || !tt.HasText("Offline — 2") || !tt.HasText("Owner") {
		t.Fatalf("texts %v", tt.Texts())
	}
	on, _ := tt.Find("Member 1")
	off, _ := tt.Find("Member 0")
	if off.Y <= on.Y {
		t.Fatal("offline member above the online section")
	}
	if err := tt.Click("Member 1"); err != nil {
		t.Fatal(err)
	}
	if s.Selected != "m1" {
		t.Fatalf("selected %q", s.Selected)
	}
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	if s.Selected != "m2" || opened != "m2" {
		t.Fatalf("keyboard chose %q, opened %q", s.Selected, opened)
	}
}

func TestMemberListBuildsOnlyVisibleAndGroupsOnce(t *testing.T) {
	var s MemberListState
	members := memberTestList(100000)
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		MemberList(c, &s, members, MemberListOptions{})
	}, 400, 500)
	gen := s.cache.gen
	tt.Frame()
	tt.Frame()
	built := 0
	for _, x := range tt.Texts() {
		if len(x) > 7 && x[:7] == "Member " {
			built++
		}
	}
	if built == 0 || built > 60 || s.cache.gen != gen {
		t.Fatalf("built %d rows, regrouped %d times", built, s.cache.gen-gen)
	}
}

func TestMemberListSnapshot(t *testing.T) {
	members := []Member{{ID: "1", Name: "Sakiko Togawa", Presence: PresenceOnline, Role: "Leader"}, {ID: "2", Name: "Uika Misumi", Presence: PresenceBusy},
		{ID: "3", Name: "Mutsumi Wakaba", Presence: PresenceAway}, {ID: "4", Name: "Umiri Yahata", Presence: PresenceOffline}}
	snapshot(t, "member_list", 340, 300, func(c *ui.Context) {
		s := ui.Local(c.Root(), "members", func() *MemberListState { return &MemberListState{} })
		ui.Column(c).Padding(20).Width(300).Children(func() { MemberList(c, *s, members, MemberListOptions{Height: 260}) })
	})
}

func TestMemberListEmptyDisabledAndReplacement(t *testing.T) {
	members := []Member{{ID: "1", Name: "One", Presence: PresenceOnline}}
	var s MemberListState
	disabled, events := true, 0
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Disabled(disabled).Children(func() {
			r := MemberList(c, &s, members, MemberListOptions{})
			if r.Changed() || r.Opened != "" {
				events++
			}
		})
	}, 400, 400)
	_ = tt.Click("One")
	tt.Key(0, ui.KeyEnter)
	if s.Selected != "" || events != 0 {
		t.Fatal("disabled member changed")
	}
	disabled = false
	members = []Member{{ID: "1", Name: "One", Presence: PresenceOffline}}
	tt.Frame()
	if !tt.HasText("Offline — 1") || tt.HasText("Online — 1") {
		t.Fatal("replacement kept stale groups")
	}
	members = nil
	tt.Frame()
	if !tt.HasText("No members") {
		t.Fatal("empty state missing")
	}
}
