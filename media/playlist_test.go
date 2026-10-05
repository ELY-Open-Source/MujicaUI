package media

import (
	"fmt"
	"slices"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

func playlistItems(n int) []PlaylistItem {
	items := make([]PlaylistItem, n)
	for i := range items {
		items[i] = PlaylistItem{ID: fmt.Sprint("t", i), Title: fmt.Sprint("Track ", i), Duration: time.Duration(180+i) * time.Second}
	}
	return items
}

func playlistIDs(items []PlaylistItem) []string {
	ids := make([]string, len(items))
	for i, it := range items {
		ids[i] = it.ID
	}
	return ids
}

func playlistTester(items *[]PlaylistItem, disabled bool, got *[]PlaylistResult) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Padding(20).Width(420).Children(func() {
			r := Playlist(c, items, PlaylistOptions{Current: "t1", Disabled: disabled})
			if r.Played != "" || r.Removed != "" || r.Moved {
				*got = append(*got, r)
			}
		})
	}, 500, 420)
}

func TestPlaylistKeyboard(t *testing.T) {
	items := playlistItems(4)
	var got []PlaylistResult
	tt := playlistTester(&items, false, &got)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	if len(got) != 1 || got[0].Played != "t0" {
		t.Fatalf("Enter played %+v", got)
	}
	tt.Key(ui.Alt, ui.KeyDown)
	if !slices.Equal(playlistIDs(items), []string{"t1", "t0", "t2", "t3"}) {
		t.Fatalf("Alt+Down order %v", playlistIDs(items))
	}
	tt.Key(0, ui.KeyDelete)
	if !slices.Equal(playlistIDs(items), []string{"t1", "t2", "t3"}) || got[len(got)-1].Removed != "t0" {
		t.Fatalf("Delete left %v, %+v", playlistIDs(items), got)
	}
}

func TestPlaylistRemoveButtonAndDrag(t *testing.T) {
	items := playlistItems(4)
	var got []PlaylistResult
	tt := playlistTester(&items, false, &got)
	if err := tt.Click("Remove Track 3"); err != nil {
		t.Fatal(err)
	}
	if len(items) != 3 || got[0].Removed != "t3" {
		t.Fatalf("remove button: %v", playlistIDs(items))
	}
	a, _ := tt.Find("Track 0")
	b, _ := tt.Find("Track 2")
	tt.Press(a.X+4, a.Y+a.H/2)
	tt.Move(a.X+4, a.Y+a.H/2+10)
	tt.Move(b.X+4, b.Y+b.H-2)
	tt.Release(b.X+4, b.Y+b.H-2)
	if !slices.Equal(playlistIDs(items), []string{"t1", "t2", "t0"}) || !got[len(got)-1].Moved {
		t.Fatalf("drag order %v", playlistIDs(items))
	}
}

func TestPlaylistBuildsOnlyVisible(t *testing.T) {
	items := playlistItems(100000)
	var got []PlaylistResult
	tt := playlistTester(&items, false, &got)
	tt.Frame()
	built := 0
	for _, s := range tt.Texts() {
		if len(s) > 6 && s[:6] == "Track " {
			built++
		}
	}
	if built == 0 || built > 40 {
		t.Fatalf("built %d rows of 100000", built)
	}
}

func TestPlaylistDisabledAndBad(t *testing.T) {
	items := playlistItems(3)
	var got []PlaylistResult
	tt := playlistTester(&items, true, &got)
	_ = tt.Click("Remove Track 0")
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyDelete)
	tt.Key(0, ui.KeyEnter)
	if len(items) != 3 || len(got) != 0 {
		t.Fatalf("disabled playlist changed: %v %+v", playlistIDs(items), got)
	}
	dup := []PlaylistItem{{ID: "a"}, {ID: "a"}}
	mediaPanics(t, "duplicate", func(c *ui.Context) { Playlist(c, &dup, PlaylistOptions{}) })
}

func TestPlaylistSnapshot(t *testing.T) {
	items := []PlaylistItem{
		{ID: "a", Title: "Ave Mujica", Subtitle: "Overture", Duration: 245 * time.Second},
		{ID: "b", Title: "Imprisoned XII", Subtitle: "Act I", Duration: 302 * time.Second},
		{ID: "c", Title: "KiLLKiSS", Subtitle: "Act II", Duration: 268 * time.Second},
	}
	snapshot(t, "playlist", 440, 240, func(c *ui.Context) {
		ui.Column(c).Padding(20).Width(400).Children(func() {
			Playlist(c, &items, PlaylistOptions{Current: "b", Height: 200})
		})
	})
}
