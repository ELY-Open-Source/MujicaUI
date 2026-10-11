package display

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

var avatarGroupPeople = []AvatarItem{
	{Name: "Sakiko Togawa"}, {Name: "Uika Misumi"}, {Name: "Mutsumi Wakaba"},
	{Name: "Umiri Yahata"}, {Name: "Nyamu Yutenji"}, {Name: "Tomori Takamatsu"},
}

func TestAvatarGroupMaxAndRest(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		AvatarGroup(c, avatarGroupPeople, AvatarGroupOptions{Max: 3})
	}, 400, 300)
	tt.Frame()
	if !tt.HasText("+3") {
		t.Fatalf("texts %v", tt.Texts())
	}
	if _, ok := tt.Find("Umiri Yahata"); ok {
		t.Fatal("fourth avatar shown past Max")
	}
}

func TestAvatarGroupNoOverflowNoButton(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		AvatarGroup(c, avatarGroupPeople[:2], AvatarGroupOptions{})
	}, 400, 300)
	tt.Frame()
	if tt.HasText("+") {
		t.Fatalf("overflow button with 2 of 4: %v", tt.Texts())
	}
}

func TestAvatarGroupExpandByKeyboard(t *testing.T) {
	var toggles []bool
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if v := AvatarGroup(c, avatarGroupPeople, AvatarGroupOptions{Max: 2}); v.Toggled {
			toggles = append(toggles, v.Open)
		}
	}, 400, 400)
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Show all 6 people") {
		t.Fatal("overflow button not focused")
	}
	tt.Key(0, ui.KeyEnter)
	for _, p := range avatarGroupPeople {
		if !tt.HasText(p.Name) {
			t.Fatalf("list lacks %q: %v", p.Name, tt.Texts())
		}
	}
	tt.Key(0, ui.KeyEscape)
	if tt.HasText("Tomori Takamatsu") {
		t.Fatal("list still open after Escape")
	}
	if len(toggles) != 2 || !toggles[0] || toggles[1] {
		t.Fatalf("toggles %v, want [true false]", toggles)
	}
}

func TestAvatarGroupExpandByClick(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		AvatarGroup(c, avatarGroupPeople, AvatarGroupOptions{Max: 4})
	}, 400, 400)
	if err := tt.Click("+2"); err != nil {
		t.Fatal(err)
	}
	if !tt.HasText("Tomori Takamatsu") {
		t.Fatal("list not open")
	}
	tt.ClickAt(390, 390)
	if tt.HasText("Tomori Takamatsu") {
		t.Fatal("list open after outside click")
	}
}

func TestAvatarGroupSnapshot(t *testing.T) {
	snapshot(t, "avatar_group", 320, 260, func(c *ui.Context) {
		ui.Column(c).Padding(20).Gap(16).AlignItems(ui.Start).Children(func() {
			AvatarGroup(c, avatarGroupPeople[:3], AvatarGroupOptions{})
			AvatarGroup(c, avatarGroupPeople, AvatarGroupOptions{Max: 3, Size: 40})
		})
	})
}

func TestAvatarGroupSmallButtonHitArea(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Row(c).AlignItems(ui.Start).Children(func() {
			AvatarGroup(c, avatarGroupPeople, AvatarGroupOptions{Max: 2, Size: 20})
		})
	}, 400, 300)
	tt.Frame()
	r, ok := tt.Find("Show all 6 people")
	if !ok || r.W < 28 || r.H < 28 {
		t.Fatalf("hit area %v", r)
	}
}
