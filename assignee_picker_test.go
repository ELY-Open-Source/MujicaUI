package mujica

import (
	"fmt"
	"slices"
	"testing"

	"github.com/egoist/mygo/ui"
	"mujicaui/theme"
)

func assigneeTester(ids *[]string, opts *AssigneePickerOptions, n *int) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		ui.Column(c).Padding(16).Children(func() {
			if AssigneePicker(c, ids, projMembers, *opts).Changed() {
				*n++
			}
		})
	}, 520, 520)
}

func TestAssigneePickerSearchAndAvatars(t *testing.T) {
	ids, n := []string{"saki"}, 0
	tt := assigneeTester(&ids, &AssigneePickerOptions{}, &n)
	if !tt.HasText("Sakiko Togawa") {
		t.Fatalf("avatar label missing: %v", tt.Texts())
	}
	_ = tt.Click("Assignees")
	tt.Type("mutsu")
	if tt.HasText("Uika Misumi") {
		t.Fatalf("search should filter: %v", tt.Texts())
	}
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	if !slices.Contains(ids, "mutsu") || n == 0 {
		t.Fatalf("ids %v", ids)
	}
}

func TestAssigneePickerMaxDisabledInvalid(t *testing.T) {
	ids, n := []string{"saki"}, 0
	tt := assigneeTester(&ids, &AssigneePickerOptions{Disabled: true}, &n)
	_ = tt.Click("Assignees")
	tt.Type("mutsu")
	tt.Key(0, ui.KeyEnter)
	if len(ids) != 1 || n != 0 {
		t.Fatal("disabled picker changed")
	}
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic for an unknown member")
		}
	}()
	bad := []string{"ghost"}
	assigneeTester(&bad, &AssigneePickerOptions{}, &n)
}

func TestAssigneePickerSnapshot(t *testing.T) {
	ids := []string{"saki", "mutsu", "uika", "nyamu"}
	snapshot(t, "assignee_picker", 560, 80, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() { AssigneePicker(c, &ids, projMembers, AssigneePickerOptions{}) })
	})
}

func TestAssigneePickerAncestorDisableCloses(t *testing.T) {
	ids, disabled := []string{}, false
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		ui.Column(c).Disabled(disabled).Children(func() { AssigneePicker(c, &ids, projMembers, AssigneePickerOptions{}) })
	}, 500, 520)
	if err := tt.Click("Assignees"); err != nil {
		t.Fatal(err)
	}
	disabled = true
	tt.Frame()
	if tt.HasText("Select all") {
		t.Fatal("ancestor-disabled popup remains visible")
	}
	if err := tt.Click("Sakiko Togawa"); err == nil && len(ids) != 0 {
		t.Fatal("disabled popup selected a member")
	}
}

func TestAssigneePickerLargeData(t *testing.T) {
	members := make([]Member, 100000)
	for i := range members {
		members[i] = Member{ID: fmt.Sprint(i), Name: fmt.Sprintf("Member %d", i)}
	}
	var ids []string
	tt := ui.NewTester(func(c *ui.Context) { Use(c, Settings{}); AssigneePicker(c, &ids, members, AssigneePickerOptions{}) }, 500, 520)
	if err := tt.Click("Assignees"); err != nil {
		t.Fatal(err)
	}
	if built := len(tt.Texts()); built > 200 {
		t.Fatalf("built %d labels", built)
	}
	tt.Type("Member 99999")
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	if !slices.Equal(ids, []string{"99999"}) {
		t.Fatalf("selected %v", ids)
	}
}

func TestAssigneePickerFullKeyboardSkipsDisabled(t *testing.T) {
	ids, changes := []string{"mutsu"}, 0
	tt := assigneeTester(&ids, &AssigneePickerOptions{Max: 1}, &changes)
	if err := tt.Click("Assignees"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	if len(ids) != 0 || changes != 1 {
		t.Fatalf("selected = %v, changes = %d", ids, changes)
	}
}

func TestAssigneePickerVersionRefreshAndBulkClear(t *testing.T) {
	ids, changes := []string{"saki"}, 0
	opts := &AssigneePickerOptions{}
	tt := assigneeTester(&ids, opts, &changes)
	if err := tt.Click("Assignees"); err != nil {
		t.Fatal(err)
	}
	ids[0], opts.Version = "mutsu", 1
	tt.Frame()
	tt.Type("Mutsumi")
	if err := tt.Click("Select all"); err != nil {
		t.Fatal(err)
	}
	if len(ids) != 0 || changes != 1 {
		t.Fatalf("selected = %v, changes = %d", ids, changes)
	}
}

func TestAssigneePickerCompactHitAreas(t *testing.T) {
	var ids []string
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{Density: theme.Compact})
		ui.Column(c).Padding(16).Children(func() { AssigneePicker(c, &ids, []Member{{ID: "only", Name: "Only Member"}}, AssigneePickerOptions{}) })
	}, 500, 520)
	if err := tt.Click("Assignees"); err != nil {
		t.Fatal(err)
	}
	for _, label := range []string{"Select all", "Only Member"} {
		r, ok := tt.Find(label)
		if !ok || r.H < 28 {
			t.Fatalf("%s hit height = %g", label, r.H)
		}
	}
}
