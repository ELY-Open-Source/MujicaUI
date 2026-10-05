package input

import (
	"slices"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

var chipItems = []ToggleItem[string]{{Value: "open", Label: "Open"}, {Value: "draft", Label: "Draft"}, {Value: "closed", Label: "Closed", Disabled: true}, {Value: "merged", Label: "Merged"}}

func chipTester(v *[]string, opts *ChoiceChipsOptions, n *int) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if ChoiceChips(c, v, chipItems, *opts).Changed() {
			*n++
		}
	}, 600, 200)
}

func TestChoiceChipsSingleAndMultiple(t *testing.T) {
	var v []string
	opts := ChoiceChipsOptions{}
	n := 0
	tt := chipTester(&v, &opts, &n)
	_ = tt.Click("Open")
	_ = tt.Click("Draft")
	if !slices.Equal(v, []string{"draft"}) || n != 2 {
		t.Fatalf("single: %v %d", v, n)
	}
	_ = tt.Click("Draft")
	if len(v) != 0 {
		t.Fatalf("single chip did not turn off: %v", v)
	}
	_ = tt.Click("Closed")
	if len(v) != 0 || n != 3 {
		t.Fatalf("disabled chip: %v %d", v, n)
	}
	opts.Multiple = true
	_ = tt.Click("Open")
	_ = tt.Click("Merged")
	if !slices.Equal(v, []string{"open", "merged"}) {
		t.Fatalf("multiple: %v", v)
	}
}

func TestChoiceChipsKeyboard(t *testing.T) {
	var v []string
	opts := ChoiceChipsOptions{Multiple: true}
	n := 0
	tt := chipTester(&v, &opts, &n)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeySpace)
	tt.Key(0, ui.KeyRight) // skips the disabled chip
	tt.Key(0, ui.KeySpace)
	if !slices.Equal(v, []string{"draft", "merged"}) {
		t.Fatalf("keyboard: %v", v)
	}
	opts.Disabled = true
	tt.Frame()
	_ = tt.Click("Open")
	if len(v) != 2 {
		t.Fatalf("disabled set changed: %v", v)
	}
}

func TestChoiceChipsRejectsTwoInSingle(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic")
		}
	}()
	v := []string{"open", "draft"}
	ui.NewTester(func(c *ui.Context) { core.Use(c, core.Settings{}); ChoiceChips(c, &v, chipItems, ChoiceChipsOptions{}) }, 300, 100)
}

func TestChoiceChipsSnapshot(t *testing.T) {
	v := []string{"open", "merged"}
	snapshot(t, "choice_chips", 480, 80, func(c *ui.Context) {
		ui.Row(c).Padding(16).Children(func() {
			ChoiceChips(c, &v, chipItems, ChoiceChipsOptions{Multiple: true})
		})
	})
}
