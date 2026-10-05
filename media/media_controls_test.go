package media

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

type controlsCalls struct {
	prev, play, pause, next int
	loops                   []LoopMode
	shuffles                []bool
}

func (m *controlsCalls) opts() MediaControlsOptions {
	return MediaControlsOptions{OnPrevious: func() { m.prev++ }, OnPlay: func() { m.play++ }, OnPause: func() { m.pause++ },
		OnNext: func() { m.next++ }, OnLoop: func(l LoopMode) { m.loops = append(m.loops, l) },
		OnShuffle: func(on bool) { m.shuffles = append(m.shuffles, on) }}
}

func controlsTester(st *MediaControlsState, o MediaControlsOptions) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Padding(20).AlignItems(ui.Start).Children(func() { MediaControls(c, *st, o) })
	}, 400, 120)
}

func TestMediaControlsClicks(t *testing.T) {
	st := MediaControlsState{Loop: LoopAll}
	var m controlsCalls
	tt := controlsTester(&st, m.opts())
	for _, name := range []string{"Previous", "Play", "Next", "Repeat all", "Shuffle"} {
		if err := tt.Click(name); err != nil {
			t.Fatal(err)
		}
	}
	if m.prev != 1 || m.play != 1 || m.next != 1 || len(m.loops) != 1 || m.loops[0] != LoopOne || len(m.shuffles) != 1 || !m.shuffles[0] {
		t.Fatalf("calls %+v", m)
	}
}

func TestMediaControlsArrowsMoveFocus(t *testing.T) {
	st := MediaControlsState{}
	var m controlsCalls
	tt := controlsTester(&st, m.opts())
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Shuffle") {
		t.Fatal("Tab did not enter at Shuffle")
	}
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyRight)
	if !tt.Focused("Play") {
		t.Fatal("arrows did not reach Play")
	}
	tt.Key(0, ui.KeyEnter)
	if m.play != 1 {
		t.Fatalf("Enter: play %d", m.play)
	}
}

func TestMediaControlsEndsAndDisabled(t *testing.T) {
	st := MediaControlsState{NoPrevious: true, NoNext: true}
	var m controlsCalls
	tt := controlsTester(&st, m.opts())
	_ = tt.Click("Previous")
	_ = tt.Click("Next")
	o := m.opts()
	o.Disabled = true
	off := controlsTester(&MediaControlsState{}, o)
	_ = off.Click("Play")
	_ = off.Click("Shuffle")
	if m.prev != 0 || m.next != 0 || m.play != 0 || len(m.shuffles) != 0 {
		t.Fatalf("blocked buttons reported %+v", m)
	}
}

func TestMediaControlsBadOptions(t *testing.T) {
	mediaPanics(t, "no callbacks", func(c *ui.Context) { MediaControls(c, MediaControlsState{}, MediaControlsOptions{}) })
	var m controlsCalls
	mediaPanics(t, "bad loop", func(c *ui.Context) { MediaControls(c, MediaControlsState{Loop: 7}, m.opts()) })
}

func TestMediaControlsSnapshot(t *testing.T) {
	var m controlsCalls
	snapshot(t, "media_controls", 300, 120, func(c *ui.Context) {
		ui.Column(c).Padding(20).Gap(12).AlignItems(ui.Start).Children(func() {
			MediaControls(c, MediaControlsState{Playing: true, Loop: LoopOne, Shuffle: true}, m.opts())
			MediaControls(c, MediaControlsState{NoPrevious: true}, m.opts())
		})
	})
}
