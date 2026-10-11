package files

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

func permissionSelectTester(v *Permission, disabled bool, changed *int) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Padding(12).Children(func() {
			if PermissionSelect(c, v, PermissionSelectOptions{Label: "Ada", Disabled: disabled}).Changed() {
				*changed++
			}
		})
	}, 360, 260)
}

func TestPermissionSelectChoose(t *testing.T) {
	v, n := PermissionView, 0
	tt := permissionSelectTester(&v, false, &n)
	if err := tt.Click("Can view"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("Can edit"); err != nil {
		t.Fatal(err)
	}
	if v != PermissionEdit || n != 1 {
		t.Fatalf("value %d changed %d", v, n)
	}
}

func TestPermissionSelectKeyboard(t *testing.T) {
	v, n := PermissionEdit, 0
	tt := permissionSelectTester(&v, false, &n)
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Ada") {
		t.Fatal("Tab did not reach the select")
	}
	tt.Key(0, ui.KeyEnter)
	tt.Key(0, ui.KeyUp)
	tt.Key(0, ui.KeyEnter)
	if v != PermissionComment {
		t.Fatalf("keys chose %d", v)
	}
}

func TestPermissionSelectDisabledAndInvalid(t *testing.T) {
	v, n := PermissionComment, 0
	tt := permissionSelectTester(&v, true, &n)
	_ = tt.Click("Can comment")
	_ = tt.Click("Can edit")
	if v != PermissionComment || n != 0 {
		t.Fatal("disabled select changed")
	}
	defer func() {
		if recover() == nil {
			t.Fatal("bad permission did not panic")
		}
	}()
	bad := Permission(7)
	permissionSelectTester(&bad, false, &n)
}

func TestPermissionSelectSnapshot(t *testing.T) {
	snapshot(t, "permission_select", 260, 70, func(c *ui.Context) {
		v := PermissionComment
		ui.Column(c).Padding(12).Children(func() { PermissionSelect(c, &v, PermissionSelectOptions{}) })
	})
}

func TestPermissionSelectDisabledClosesAndStaysClosed(t *testing.T) {
	v, disabled := PermissionView, false
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Disabled(disabled).Children(func() { PermissionSelect(c, &v, PermissionSelectOptions{Label: "Permission"}) })
	}, 400, 300)
	if err := tt.Click("Can view"); err != nil {
		t.Fatal(err)
	}
	disabled = true
	tt.Frame()
	if tt.HasText("Can edit") {
		t.Fatal("disabled permission popup stayed open")
	}
	disabled = false
	tt.Frame()
	if tt.HasText("Can edit") {
		t.Fatal("reenabling permission reopened popup")
	}
}
