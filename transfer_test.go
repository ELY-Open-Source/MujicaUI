package mujica

import (
	"slices"
	"testing"

	"github.com/egoist/mygo/ui"
)

var transferCast = []TransferItem{
	{ID: "uika", Label: "Uika"}, {ID: "mutsumi", Label: "Mutsumi"}, {ID: "umiri", Label: "Umiri"},
	{ID: "nyamu", Label: "Nyamu", Disabled: true}, {ID: "sakiko", Label: "Sakiko"},
}

func transferTester(target *[]string, items *[]TransferItem, opts TransferOptions, changes *int) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		if Transfer(c, target, *items, opts).Changed() {
			*changes++
		}
	}, 640, 420)
}

func TestTransferMoveSelected(t *testing.T) {
	target, items, n := []string{"sakiko"}, slices.Clone(transferCast), 0
	tt := transferTester(&target, &items, TransferOptions{}, &n)
	_ = tt.Click("Mutsumi")
	_ = tt.Click("Uika")
	if err := tt.Click("Move to Chosen"); err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(target, []string{"sakiko", "uika", "mutsumi"}) || n != 1 {
		t.Fatalf("target %q changes %d", target, n)
	}
	_ = tt.Click("Sakiko")
	_ = tt.Click("Move to Available")
	if !slices.Equal(target, []string{"uika", "mutsumi"}) || n != 2 {
		t.Fatalf("back: %q", target)
	}
}

func TestTransferSelectAllSkipsDisabledAndFollowsSearch(t *testing.T) {
	target, items, n := []string{}, slices.Clone(transferCast), 0
	tt := transferTester(&target, &items, TransferOptions{}, &n)
	_ = tt.Click("Search Available")
	tt.Type("u")
	if tt.HasText("Sakiko") {
		t.Fatal("search did not filter")
	}
	_ = tt.Click("Select all in Available")
	_ = tt.Click("Move to Chosen")
	if !slices.Equal(target, []string{"uika", "mutsumi", "umiri"}) {
		t.Fatalf("select all moved %q", target)
	}
	_ = tt.Click("Search Available")
	tt.Type("zzz")
	if !tt.HasText("No matches") {
		t.Fatal("no empty search state")
	}
}

func TestTransferKeyboardAndStableIDs(t *testing.T) {
	target, items, n := []string{}, slices.Clone(transferCast), 0
	tt := transferTester(&target, &items, TransferOptions{}, &n)
	_ = tt.Click("Umiri")
	// Items reorder and relabel; the check stays on the same ID.
	items = []TransferItem{{ID: "umiri", Label: "Yahata Umiri"}, {ID: "uika", Label: "Uika"}}
	tt.Frame()
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeySpace)
	_ = tt.Click("Move to Chosen")
	if !slices.Equal(target, []string{"umiri", "uika"}) {
		t.Fatalf("target %q", target)
	}
}

func TestTransferDisabledAndBadIDs(t *testing.T) {
	target, items, n := []string{}, slices.Clone(transferCast), 0
	tt := transferTester(&target, &items, TransferOptions{Disabled: true}, &n)
	_ = tt.Click("Uika")
	_ = tt.Click("Move to Chosen")
	if len(target) != 0 {
		t.Fatal("disabled transfer moved")
	}
	defer func() {
		if recover() == nil {
			t.Fatal("unknown target ID accepted")
		}
	}()
	bad := []string{"ghost"}
	transferTester(&bad, &items, TransferOptions{}, &n)
}

func TestTransferSnapshot(t *testing.T) {
	target := []string{"sakiko", "mutsumi"}
	snapshot(t, "transfer", 620, 380, func(c *ui.Context) {
		ui.Column(c).Padding(20).Children(func() {
			p := Transfer(c, &target, transferCast, TransferOptions{Height: 200})
			st := ui.Local(p.Element, "transfer", func() transferState { return transferState{} })
			if !st.checked["uika"] {
				st.checked["uika"] = true
				c.Invalidate()
			}
		})
	})
}

// MyGo's tester cannot read descriptions, so this checks the text the select-all box gets.
func TestTransferPartialDescription(t *testing.T) {
	var got [3]string
	ui.Render(func(c *ui.Context) {
		Use(c, Settings{})
		got = [3]string{transferPartial(c, 2, 3), transferPartial(c, 3, 3), transferPartial(c, 0, 3)}
	}, 100, 100, 1)
	if got != [3]string{"2 of 3 selected", "", ""} {
		t.Fatalf("descriptions %q", got)
	}
}
