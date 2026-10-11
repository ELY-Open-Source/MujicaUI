package files

import (
	"fmt"
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

func presenceSample() []Collaborator {
	return []Collaborator{{ID: "a", Name: "Ada Lovelace"}, {ID: "b", Name: "Blaise Pascal", Away: true},
		{ID: "c", Name: "Carl Gauss"}, {ID: "d", Name: "Dora Maar"}, {ID: "e", Name: "Emmy Noether"}}
}

type presenceProbe struct {
	people    []Collaborator
	following string
	max       int
	disabled  bool
	events    []string
}

func (p *presenceProbe) view(c *ui.Context) {
	core.Use(c, core.Settings{})
	ui.Column(c).Padding(12).Children(func() {
		r := PresenceAvatars(c, p.people, PresenceAvatarsOptions{Max: p.max, Following: p.following, Disabled: p.disabled})
		if r.Follow != "" {
			p.following = r.Follow
			p.events = append(p.events, "follow "+r.Follow)
		}
		if r.Unfollow {
			p.following = ""
			p.events = append(p.events, "unfollow")
		}
	})
}

func TestPresenceAvatarsFollow(t *testing.T) {
	p := &presenceProbe{people: presenceSample(), max: 3}
	tt := ui.NewTester(p.view, 400, 300)
	if !tt.HasText("+2") {
		t.Fatalf("overflow missing: %q", tt.Texts())
	}
	if err := tt.Click("Follow Ada Lovelace"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("Stop following Ada Lovelace"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("2 more people"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("Follow Emmy Noether"); err != nil {
		t.Fatal(err)
	}
	want := []string{"follow a", "unfollow", "follow e"}
	if len(p.events) != 3 || p.events[2] != want[2] {
		t.Fatalf("events %v", p.events)
	}
}

func TestPresenceAvatarsKeyboardAndDisabled(t *testing.T) {
	p := &presenceProbe{people: presenceSample()}
	tt := ui.NewTester(p.view, 400, 300)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyEnter)
	if p.following != "b" {
		t.Fatalf("keyboard followed %q", p.following)
	}
	p = &presenceProbe{people: presenceSample(), disabled: true}
	tt = ui.NewTester(p.view, 400, 300)
	_ = tt.Click("Follow Ada Lovelace")
	if p.following != "" {
		t.Fatal("disabled avatars followed")
	}
}

func TestPresenceAvatarsInvalid(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("missing ID did not panic")
		}
	}()
	ui.NewTester((&presenceProbe{people: []Collaborator{{Name: "x"}}}).view, 400, 300)
}

func TestPresenceAvatarsSnapshot(t *testing.T) {
	snapshot(t, "presence_avatars", 300, 70, func(c *ui.Context) {
		ui.Column(c).Padding(12).Children(func() { PresenceAvatars(c, presenceSample(), PresenceAvatarsOptions{Following: "c"}) })
	})
}

func TestPresenceAvatarsSmallHitArea(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		PresenceAvatars(c, presenceSample(), PresenceAvatarsOptions{Size: 12, Max: 1})
	}, 300, 200)
	for _, name := range []string{"Follow Ada Lovelace", "4 more people"} {
		r, ok := tt.Find(name)
		if !ok || r.W < 28 || r.H < 28 {
			t.Fatalf("%s: hit area %v", name, r)
		}
	}
}

func TestPresenceAvatarsLargeOverflow(t *testing.T) {
	people := make([]Collaborator, 100000)
	for i := range people {
		people[i] = Collaborator{ID: fmt.Sprint(i), Name: fmt.Sprintf("Member %d", i)}
	}
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		PresenceAvatars(c, people, PresenceAvatarsOptions{Max: 1})
	}, 400, 300)
	if err := tt.Click("99999 more people"); err != nil {
		t.Fatal(err)
	}
	n := 0
	for _, text := range tt.Texts() {
		if strings.HasPrefix(text, "Member ") {
			n++
		}
	}
	if n == 0 || n > 60 {
		t.Fatalf("built %d overflow member labels", n)
	}
}
