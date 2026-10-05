package mujica

import (
	"fmt"
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"
)

type shareProbe struct {
	open     bool
	members  []ShareMember
	disabled bool
	invites  []ShareInvite
	changed  []string
	removed  []string
	copied   int
}

func newShareProbe() *shareProbe {
	return &shareProbe{open: true, members: []ShareMember{
		{ID: "o", Name: "Ada Lovelace", Email: "ada@court.org", Permission: PermissionEdit, Owner: true},
		{ID: "b", Name: "Blaise Pascal", Email: "blaise@court.org", Permission: PermissionView},
	}}
}

func (p *shareProbe) view(c *ui.Context) {
	Use(c, Settings{})
	r := ShareDialog(c, &p.open, &p.members, ShareDialogOptions{Title: "Ledger", Link: "https://court.org/ledger", Disabled: p.disabled})
	if r.Invited != nil {
		p.invites = append(p.invites, *r.Invited)
	}
	if r.Changed != "" {
		p.changed = append(p.changed, r.Changed)
	}
	if r.Removed != "" {
		p.removed = append(p.removed, r.Removed)
	}
	if r.Copied {
		p.copied++
	}
}

func TestShareDialogInvite(t *testing.T) {
	p := newShareProbe()
	tt := ui.NewTester(p.view, 640, 520)
	if !tt.HasText("Share “Ledger”") || !tt.HasText("Owner") {
		t.Fatalf("texts %q", tt.Texts())
	}
	for !tt.Focused("Email to invite") {
		tt.Key(0, ui.KeyTab)
	}
	tt.Type("carl@court.org")
	tt.Key(0, ui.KeyEnter)
	if len(p.invites) != 1 || p.invites[0] != (ShareInvite{"carl@court.org", PermissionView}) {
		t.Fatalf("invites %v", p.invites)
	}
}

func TestShareDialogChangeRemoveCopy(t *testing.T) {
	p := newShareProbe()
	tt := ui.NewTester(p.view, 640, 520)
	if err := tt.Click("Copy link"); err != nil {
		t.Fatal(err)
	}
	if p.copied != 1 || tt.Clipboard() != "https://court.org/ledger" {
		t.Fatalf("copied %d clipboard %q", p.copied, tt.Clipboard())
	}
	if err := tt.Click("Remove Blaise Pascal"); err != nil {
		t.Fatal(err)
	}
	if len(p.removed) != 1 || len(p.members) != 1 || tt.HasText("Blaise Pascal") {
		t.Fatalf("removed %v members %v", p.removed, p.members)
	}
	if tt.HasText("Remove Ada Lovelace") {
		t.Fatal("owner can be removed")
	}
}

func TestShareDialogPermissionChange(t *testing.T) {
	p := newShareProbe()
	tt := ui.NewTester(p.view, 640, 520)
	// Two "Can view" selects show: the invite's, then Blaise's.
	inv, _ := tt.Find("Permission for the invite")
	r, _ := tt.Find("Blaise Pascal")
	tt.ClickAt(inv.X+inv.W/2, r.Y+r.H/2)
	if err := tt.Click("Can comment"); err != nil {
		t.Fatal(err)
	}
	if p.members[1].Permission != PermissionComment || len(p.changed) != 1 || p.changed[0] != "b" {
		t.Fatalf("permission %d changed %v", p.members[1].Permission, p.changed)
	}
}

func TestShareDialogDisabled(t *testing.T) {
	p := newShareProbe()
	p.disabled = true
	tt := ui.NewTester(p.view, 640, 520)
	_ = tt.Click("Remove Blaise Pascal")
	_ = tt.Click("Copy link")
	if len(p.members) != 2 || p.copied != 0 {
		t.Fatal("disabled dialog acted")
	}
}

func TestShareDialogSnapshot(t *testing.T) {
	snapshot(t, "share_dialog", 560, 440, func(c *ui.Context) {
		p := ui.Local(c.Root(), "snap", newShareProbe)
		ShareDialog(c, &(*p).open, &(*p).members, ShareDialogOptions{Title: "Ledger", Link: "https://court.org/ledger"})
	})
}

func TestShareDialogAncestorDisabledCloses(t *testing.T) {
	open, disabled := true, false
	members := []ShareMember{{ID: "a", Name: "Ada"}}
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		ui.Column(c).Disabled(disabled).Children(func() {
			ShareDialog(c, &open, &members, ShareDialogOptions{Title: "Decree"})
		})
	}, 600, 400)
	disabled = true
	tt.Frame()
	if open || tt.HasText("Share “Decree”") {
		t.Fatal("disabled dialog stayed open")
	}
}

func TestShareDialogCopyCallback(t *testing.T) {
	open := true
	members := []ShareMember{}
	var copied []string
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		ShareDialog(c, &open, &members, ShareDialogOptions{Title: "Ledger", Link: "https://court.org/ledger", OnCopy: func(s string) { copied = append(copied, s) }})
	}, 600, 400)
	if err := tt.Click("Copy link"); err != nil {
		t.Fatal(err)
	}
	if len(copied) != 1 || copied[0] != "https://court.org/ledger" || tt.Clipboard() != "" {
		t.Fatalf("callback count %d, clipboard %q", len(copied), tt.Clipboard())
	}
}
func TestShareDialogLargeList(t *testing.T) {
	open := true
	members := make([]ShareMember, 100000)
	for i := range members {
		members[i] = ShareMember{ID: fmt.Sprint(i), Name: fmt.Sprintf("Member %d", i)}
	}
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		ShareDialog(c, &open, &members, ShareDialogOptions{Title: "Roster"})
	}, 640, 520)
	n := 0
	for _, text := range tt.Texts() {
		if strings.HasPrefix(text, "Member ") {
			n++
		}
	}
	if n == 0 || n > 60 {
		t.Fatalf("built %d member labels", n)
	}
}

func TestShareDialogSameTitleIndependent(t *testing.T) {
	openA, openB := true, false
	membersA, membersB := []ShareMember{}, []ShareMember{}
	var invited *ShareInvite
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		ui.Column(c).Children(func() {
			ShareDialog(c, &openA, &membersA, ShareDialogOptions{Title: "Ledger"})
			r := ShareDialog(c, &openB, &membersB, ShareDialogOptions{Title: "Ledger"})
			if r.Invited != nil {
				invited = r.Invited
			}
		})
	}, 640, 520)
	for i := 0; i < 10 && !tt.Focused("Email to invite"); i++ {
		tt.Key(0, ui.KeyTab)
	}
	tt.Type("ada@court.org")
	openA, openB = false, true
	tt.Frame()
	if err := tt.Click("Invite"); err != nil {
		t.Fatal(err)
	}
	if invited != nil {
		t.Fatal("second instance inherited invitation draft")
	}
}
