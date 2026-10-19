package chat

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
)

func TestVoiceInputButtonStartStopAndTime(t *testing.T) {
	var opts VoiceInputOptions
	presses := 0
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if VoiceInputButton(c, opts).Changed() {
			presses++
			if opts.Recording {
				opts = VoiceInputOptions{}
			} else {
				opts = VoiceInputOptions{Recording: true, Started: c.Now().Add(-65 * time.Second)}
			}
		}
	}, 300, 100)
	_ = tt.Click("Start voice input")
	tt.Frame()
	if !opts.Recording || !tt.HasText("1:05") {
		t.Fatalf("recording %v texts %q", opts.Recording, tt.Texts())
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeySpace)
	if opts.Recording || presses != 2 {
		t.Fatalf("stop: recording %v presses %d", opts.Recording, presses)
	}
}

func TestVoiceInputButtonDeniedAndDisabled(t *testing.T) {
	opts := VoiceInputOptions{Denied: true}
	presses := 0
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if VoiceInputButton(c, opts).Changed() {
			presses++
		}
	}, 400, 100)
	if !tt.HasText("Microphone access denied") {
		t.Fatalf("texts %q", tt.Texts())
	}
	_ = tt.Click("Start voice input")
	opts.Disabled = true
	_ = tt.Click("Start voice input")
	if presses != 1 {
		t.Fatalf("presses %d, want 1 (retry once, then disabled)", presses)
	}
}

func TestVoiceInputButtonRecordingNeedsStart(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic for recording without Started")
		}
	}()
	ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		VoiceInputButton(c, VoiceInputOptions{Recording: true})
	}, 100, 100)
}

func TestVoiceInputButtonSnapshot(t *testing.T) {
	snapshot(t, "voice_input_button", 360, 170, func(c *ui.Context) {
		ui.Column(c).Padding(12).Gap(8).Children(func() {
			VoiceInputButton(c, VoiceInputOptions{})
			VoiceInputButton(c, VoiceInputOptions{Recording: true, Started: c.Now().Add(-12 * time.Second)})
			VoiceInputButton(c, VoiceInputOptions{Denied: true})
		})
	})
}
