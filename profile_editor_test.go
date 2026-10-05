package mujica

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestProfileEditorUnsavedSaveAndDiscard(t *testing.T) {
	saved := Profile{Name: "Sakiko", Bio: "Keyboard."}
	p := saved
	var saves, avatars, discards int
	opts := ProfileEditorOptions{Saved: saved}
	tt := accountTester(480, 500, func(c *ui.Context) {
		r := ProfileEditor(c, &p, opts)
		if r.Save() {
			saves++
		}
		if r.ChangeAvatar() {
			avatars++
		}
		if r.Discarded() {
			discards++
		}
	})
	if tt.HasText("Unsaved changes") {
		t.Fatal("clean profile marked unsaved")
	}
	_ = tt.Click("Save")
	if saves != 0 {
		t.Fatal("saved without changes")
	}
	_ = tt.Click("Name")
	tt.Type(" Togawa")
	if !tt.HasText("Unsaved changes") || p.Name != "Sakiko Togawa" {
		t.Fatalf("edit: %q", p.Name)
	}
	_ = tt.Click("Save")
	_ = tt.Click("Change picture")
	if saves != 1 || avatars != 1 {
		t.Fatalf("saves %d avatars %d", saves, avatars)
	}
	_ = tt.Click("Discard")
	if p != saved || discards != 1 || tt.HasText("Unsaved changes") {
		t.Fatalf("discard: %+v", p)
	}
	_ = tt.Click("Name")
	tt.Key(ui.Cmd, ui.KeyA)
	tt.Key(0, ui.KeyBackspace)
	_ = tt.Click("Save")
	if saves != 1 || !tt.HasText("Required") {
		t.Fatal("saved an empty name")
	}
}

func TestProfileEditorSavingBlocksEdits(t *testing.T) {
	p := Profile{Name: "Uika"}
	tt := accountTester(480, 500, func(c *ui.Context) {
		ProfileEditor(c, &p, ProfileEditorOptions{Saved: Profile{Name: "Misumi"}, Saving: true})
	})
	_ = tt.Click("Name")
	tt.Type("x")
	if p.Name != "Uika" {
		t.Fatal("edited while saving")
	}
}

func TestProfileEditorSnapshot(t *testing.T) {
	snapshot(t, "profile_editor", 460, 400, func(c *ui.Context) {
		p := ui.Local(c.Root(), "p", func() Profile { return Profile{Name: "Sakiko Togawa", Bio: "Keyboard, Ave Mujica."} })
		ui.Column(c).Padding(16).Children(func() { ProfileEditor(c, p, ProfileEditorOptions{Saved: Profile{Name: "Sakiko"}}) })
	})
}

func TestProfileEditorDisabledAndBioLimit(t *testing.T) {
	p := Profile{Name: "Uika"}
	disabled := false
	saves := 0
	tt := accountTester(600, 500, func(c *ui.Context) {
		if ProfileEditor(c, &p, ProfileEditorOptions{Saved: Profile{Name: "Uika"}, BioLimit: 3, Disabled: disabled}).Save() {
			saves++
		}
	})
	if err := tt.Click("Bio"); err != nil {
		t.Fatal(err)
	}
	tt.Type("初音未来")
	if err := tt.Click("Save"); err != nil {
		t.Fatal(err)
	}
	if p.Bio != "初音未来" || saves != 0 {
		t.Fatal("over-limit bio is saved or modified")
	}
	disabled = true
	tt.Frame()
	tt.Type("edit")
	if err := tt.Click("Discard"); err != nil {
		t.Fatal(err)
	}
	if p.Bio != "初音未来" {
		t.Fatal("disabled profile changes")
	}
	accountPanics(t, func(c *ui.Context) { ProfileEditor(c, &p, ProfileEditorOptions{BioLimit: -1}) })
}
