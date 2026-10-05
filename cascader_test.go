package mujica

import (
	"errors"
	"slices"
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"
)

// cascaderFake is a caller's store: roots are known, children load on request.
type cascaderFake struct {
	levels map[string]CascaderLevel
	loads  []string
}

func newCascaderFake() *cascaderFake {
	return &cascaderFake{levels: map[string]CascaderLevel{
		"": {Nodes: []CascaderNode{{Value: "eu", Label: "Europe"}, {Value: "as", Label: "Asia"}, {Value: "aq", Label: "Antarctica", Disabled: true}}},
	}}
}

func (f *cascaderFake) level(path []string) CascaderLevel {
	if lv, ok := f.levels[strings.Join(path, "/")]; ok {
		return lv
	}
	return CascaderLevel{Status: CascaderIdle}
}

func (f *cascaderFake) load(path []string) {
	key := strings.Join(path, "/")
	f.loads = append(f.loads, key)
	f.levels[key] = CascaderLevel{Status: CascaderLoading}
}

func (f *cascaderFake) finish(key string, lv CascaderLevel) { f.levels[key] = lv }

func cascaderTester(path *[]string, f *cascaderFake, anyLevel bool, changes *int) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		ui.Column(c).Padding(10).Children(func() {
			if Cascader(c, path, CascaderOptions{Level: f.level, Load: f.load, AnyLevel: anyLevel, Label: "Place"}).Changed() {
				*changes++
			}
		})
	}, 700, 500)
}

func TestCascaderLoadsAndChoosesLeaf(t *testing.T) {
	var path []string
	changes := 0
	f := newCascaderFake()
	tt := cascaderTester(&path, f, false, &changes)
	_ = tt.Click("Select…")
	_ = tt.Click("Europe")
	if !slices.Equal(f.loads, []string{"eu"}) || !tt.HasText("Loading…") {
		t.Fatalf("loads %v, texts %v", f.loads, tt.Texts())
	}
	if len(path) != 0 {
		t.Fatalf("branch chosen without AnyLevel: %v", path)
	}
	f.finish("eu", CascaderLevel{Nodes: []CascaderNode{{Value: "fr", Label: "France"}}})
	tt.Frame()
	_ = tt.Click("France")
	f.finish("eu/fr", CascaderLevel{Nodes: []CascaderNode{{Value: "paris", Label: "Paris", Leaf: true}}})
	tt.Frame()
	_ = tt.Click("Paris")
	if !slices.Equal(path, []string{"eu", "fr", "paris"}) || changes != 1 {
		t.Fatalf("path %v changes %d", path, changes)
	}
	if !tt.HasText("Europe / France / Paris") || tt.HasText("Asia") {
		t.Fatalf("field shows %v", tt.Texts())
	}
}

func TestCascaderRetryAfterFailure(t *testing.T) {
	var path []string
	changes := 0
	f := newCascaderFake()
	tt := cascaderTester(&path, f, false, &changes)
	_ = tt.Click("Select…")
	_ = tt.Click("Asia")
	f.finish("as", CascaderLevel{Status: CascaderFailed, Err: errors.New("timeout")})
	tt.Frame()
	if !tt.HasText("Couldn't load: timeout") {
		t.Fatalf("texts %v", tt.Texts())
	}
	if err := tt.Click("Retry"); err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(f.loads, []string{"as", "as"}) || !tt.HasText("Loading…") {
		t.Fatalf("loads %v", f.loads)
	}
}

func TestCascaderAnyLevelAndDisabled(t *testing.T) {
	var path []string
	changes := 0
	f := newCascaderFake()
	tt := cascaderTester(&path, f, true, &changes)
	_ = tt.Click("Select…")
	_ = tt.Click("Antarctica")
	if len(path) != 0 {
		t.Fatalf("disabled node chosen: %v", path)
	}
	_ = tt.Click("Asia")
	if !slices.Equal(path, []string{"as"}) {
		t.Fatalf("AnyLevel path %v", path)
	}
}

func TestCascaderKeyboard(t *testing.T) {
	var path []string
	changes := 0
	f := newCascaderFake()
	f.finish("eu", CascaderLevel{Nodes: []CascaderNode{{Value: "fr", Label: "France", Leaf: true}, {Value: "it", Label: "Italy", Leaf: true}}})
	tt := cascaderTester(&path, f, false, &changes)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyDown)
	if !tt.Focused("Europe") {
		t.Fatal("first row not focused on open")
	}
	tt.Key(0, ui.KeyRight)
	if !tt.Focused("France") {
		t.Fatal("Right did not move into the children")
	}
	tt.Key(0, ui.KeyLeft)
	if !tt.Focused("Europe") {
		t.Fatal("Left did not go back")
	}
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	if !slices.Equal(path, []string{"eu", "it"}) {
		t.Fatalf("path %v", path)
	}
	if !tt.Focused("Place") {
		t.Fatal("focus did not return to the field")
	}
}

func TestCascaderNeedsLevel(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic")
		}
	}()
	var path []string
	ui.Render(func(c *ui.Context) {
		Use(c, Settings{})
		Cascader(c, &path, CascaderOptions{})
	}, 100, 100, 1)
}

func TestCascaderSnapshot(t *testing.T) {
	f := newCascaderFake()
	f.finish("eu", CascaderLevel{Nodes: []CascaderNode{{Value: "fr", Label: "France"}}})
	f.finish("eu/fr", CascaderLevel{Nodes: []CascaderNode{{Value: "paris", Label: "Paris", Leaf: true}}})
	snapshot(t, "cascader", 460, 110, func(c *ui.Context) {
		none, some := []string{}, []string{"eu", "fr", "paris"}
		ui.Row(c).Padding(20).Gap(16).Children(func() {
			Cascader(c, &none, CascaderOptions{Level: f.level, Load: f.load})
			Cascader(c, &some, CascaderOptions{Level: f.level, Load: f.load})
		})
	})
}
