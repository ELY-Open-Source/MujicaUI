package mujica

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestVideoCallGridTilesAndSpeaker(t *testing.T) {
	people := []CallParticipant{{ID: "1", Name: "Sakiko", Speaking: true}, {ID: "2", Name: "Uika", Muted: true, Frame: mediaTestBitmap(32, 18)},
		{ID: "3", Name: "Mutsumi", CameraOff: true}, {ID: "4", Name: "Umiri"}, {ID: "5", Name: "Nyamu"}}
	var rects []ui.Rect
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		ui.Column(c).Width(600).Children(func() { VideoCallGrid(c, people, VideoCallGridOptions{}) })
	}, 700, 600)
	for _, name := range []string{"Sakiko, speaking", "Uika, muted", "Mutsumi, camera off", "Umiri", "Nyamu"} {
		r, ok := tt.Find(name)
		if !ok {
			t.Fatalf("no tile %q in %v", name, tt.Texts())
		}
		rects = append(rects, r)
	}
	if rects[3].Y <= rects[0].Y || rects[1].Y != rects[0].Y || rects[2].Y != rects[0].Y {
		t.Fatalf("tiles %v", rects)
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	if !tt.Focused("Uika, muted") {
		t.Fatal("arrows do not move between tiles")
	}
	mediaPanics(t, "duplicate", func(c *ui.Context) {
		VideoCallGrid(c, []CallParticipant{{ID: "a", Name: "A"}, {ID: "a", Name: "B"}}, VideoCallGridOptions{})
	})
}

func TestVideoCallGridSnapshot(t *testing.T) {
	people := []CallParticipant{{ID: "1", Name: "Sakiko", Speaking: true, Frame: mediaTestBitmap(160, 90)}, {ID: "2", Name: "Uika", Muted: true},
		{ID: "3", Name: "Mutsumi", CameraOff: true}, {ID: "4", Name: "Umiri", Frame: mediaTestBitmap(90, 160)}}
	snapshot(t, "video_call_grid", 520, 340, func(c *ui.Context) {
		ui.Column(c).Padding(20).Width(480).Children(func() { VideoCallGrid(c, people, VideoCallGridOptions{}) })
	})
}
