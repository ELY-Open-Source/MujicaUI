package main

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/egoist/mygo/ui"
	mujica "mujicaui"
)

func chatinputDemoPromptComposer(c *ui.Context, log func(string)) {
	type state struct {
		text                                               string
		disabled, cmd, composing, attached, web, streaming bool
	}
	s := ui.Local(c.Root(), "221", func() state { return state{attached: true} })
	choiceDemo(c, func() {
		choiceFlag(c, &s.disabled, "Disabled", log)
		choiceFlag(c, &s.cmd, "Cmd+Enter sends", log)
		choiceFlag(c, &s.composing, "Composition guard", log)
		choiceFlag(c, &s.attached, "Attachment", log)
	}, func() {
		ui.Column(c).Width(520).Children(func() {
			r := mujica.PromptComposer(c, &s.text, mujica.PromptComposerOptions{Label: "Prompt", Placeholder: "Ask the herald", Disabled: s.disabled, CmdEnterSends: s.cmd, Composing: s.composing,
				Attachments: func() {
					if s.attached && mujica.AttachmentChip(c, mujica.Attachment{Name: "guests.csv", Size: 4096}).Removed() {
						s.attached = false
						log("remove attachment")
					}
				},
				Tools: func() {
					if _, ok := mujica.ToolToggleMenu(c, []mujica.ToolItem{{ID: "web", Label: "Web search", On: &s.web}}).Changed(); ok {
						log("toggle web search")
					}
				},
				Actions: func() {
					r := mujica.SendButton(c, mujica.SendButtonOptions{Streaming: s.streaming, Disabled: strings.TrimSpace(s.text) == "" || s.composing, Shortcut: "Enter"})
					if r.Sent() {
						s.streaming = true
						s.text = ""
						log("send")
					}
					if r.Stopped() {
						s.streaming = false
						log("stop")
					}
				},
			})
			if r.Changed() {
				log(fmt.Sprintf("edit %d characters", len(s.text)))
			}
			if r.Submitted() {
				s.text = ""
				s.streaming = true
				log("submit")
			}
		})
	})
}

func chatinputDemoAttachmentChip(c *ui.Context, log func(string)) {
	type state struct {
		upload, fail, disabled, removed bool
		progress                        float64
	}
	s := ui.Local(c.Root(), "222", func() state { return state{upload: true, progress: 0.4} })
	choiceDemo(c, func() {
		choiceFlag(c, &s.upload, "Uploading", log)
		choiceFlag(c, &s.fail, "Failed", log)
		choiceFlag(c, &s.disabled, "Disabled", log)
		if mujica.Slider(c, &s.progress, mujica.SliderOptions{Min: 0, Max: 1, Step: 0.1, Label: "Upload progress"}).Changed() {
			log(fmt.Sprint("progress ", s.progress))
		}
		if mujica.Button(c, "Restore attachment", mujica.ButtonOptions{}).Clicked() {
			s.removed = false
			log("restore")
		}
	}, func() {
		ui.Column(c).Disabled(s.disabled).Children(func() {
			if s.removed {
				ui.Text(c, "Attachment removed")
				return
			}
			a := mujica.Attachment{Name: "guests.csv", Size: 4096, Uploading: s.upload, Progress: s.progress}
			if s.fail {
				a.Err = errors.New("upload rejected")
			}
			r := mujica.AttachmentChip(c, a)
			if r.Removed() {
				s.removed = true
				log("remove")
			}
			if r.Retried() {
				s.fail = false
				s.upload = true
				s.progress = 0
				log("retry")
			}
		})
	})
}

func chatinputPictures() []mujica.PastedImage {
	var out []mujica.PastedImage
	for i, col := range []string{"#7D2034", "#355E80"} {
		svg := fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="72" height="72"><rect width="72" height="72" fill="%s"/><circle cx="36" cy="36" r="18" fill="#F2EBDD"/></svg>`, col)
		out = append(out, mujica.PastedImage{ID: fmt.Sprint(i), Name: fmt.Sprintf("Picture %d", i+1), Source: ui.MustParseSVG([]byte(svg))})
	}
	return out
}

func chatinputDemoPasteImagePreview(c *ui.Context, log func(string)) {
	type state struct {
		images   []mujica.PastedImage
		disabled bool
	}
	s := ui.Local(c.Root(), "223", func() state { return state{images: chatinputPictures()} })
	choiceDemo(c, func() {
		choiceFlag(c, &s.disabled, "Disabled", log)
		if mujica.Button(c, "Restore images", mujica.ButtonOptions{}).Clicked() {
			s.images = chatinputPictures()
			log("restore images")
		}
	}, func() {
		ui.Column(c).Padding(12).Disabled(s.disabled).Children(func() {
			if id, ok := mujica.PasteImagePreview(c, s.images).Removed(); ok {
				for i, im := range s.images {
					if im.ID == id {
						s.images = append(s.images[:i], s.images[i+1:]...)
						break
					}
				}
				log("remove image " + id)
			}
		})
		ui.Text(c, "Decoded images arrive from the caller's paste or drop handler.")
	})
}

func chatinputDemoDragDropOverlay(c *ui.Context, log func(string)) {
	type state struct {
		imagesOnly, disabled bool
		accepted, rejected   int
	}
	s := ui.Local(c.Root(), "224", func() state { return state{imagesOnly: true} })
	choiceDemo(c, func() {
		choiceFlag(c, &s.imagesOnly, "PNG only", log)
		choiceFlag(c, &s.disabled, "Disabled", log)
	}, func() {
		ui.Column(c).Width(520).Height(180).Disabled(s.disabled).Children(func() {
			r := mujica.DragDropOverlay(c, mujica.DragDropOverlayOptions{Accept: func(p string) bool { return !s.imagesOnly || strings.EqualFold(filepath.Ext(p), ".png") }}, func() {
				ui.Box(c).Fill().Border(1, mujica.Tokens(c).ControlBorder).Center().Children(func() { ui.Text(c, "Drop files here") })
			})
			r.Element.Fill()
			if n := len(r.Accepted()); n > 0 {
				s.accepted += n
				log(fmt.Sprintf("accept %d files", n))
			}
			if n := len(r.Rejected()); n > 0 {
				s.rejected += n
				log(fmt.Sprintf("reject %d files", n))
			}
		})
		ui.Text(c, fmt.Sprintf("Accepted %d · Rejected %d", s.accepted, s.rejected))
	})
}

func chatinputDemoVoiceInputButton(c *ui.Context, log func(string)) {
	s := ui.Local(c.Root(), "225", func() mujica.VoiceInputOptions { return mujica.VoiceInputOptions{} })
	choiceDemo(c, func() {
		choiceFlag(c, &s.Disabled, "Disabled", log)
		if mujica.Button(c, "Permission failure", mujica.ButtonOptions{}).Clicked() {
			s.Recording = false
			s.Denied = true
			log("permission denied")
		}
	}, func() {
		if mujica.VoiceInputButton(c, *s).Changed() {
			s.Recording = !s.Recording
			s.Denied = false
			s.Started = c.Now()
			log(fmt.Sprintf("recording %v", s.Recording))
		}
		ui.Text(c, "The caller records audio; this demo changes recording state.")
	})
}

func chatinputDemoVoiceWaveform(c *ui.Context, log func(string)) {
	type state struct {
		silent bool
		level  float64
	}
	s := ui.Local(c.Root(), "226", func() state { return state{level: 0.6} })
	choiceDemo(c, func() {
		choiceFlag(c, &s.silent, "Silent", log)
		if mujica.Slider(c, &s.level, mujica.SliderOptions{Min: 0, Max: 1, Step: 0.1, Label: "Input level"}).Changed() {
			log(fmt.Sprint("level ", s.level))
		}
	}, func() {
		levels := make([]float32, 40)
		for i := range levels {
			if !s.silent {
				levels[i] = float32(s.level) * float32(i%7+1) / 7
			}
		}
		mujica.VoiceWaveform(c, levels, mujica.VoiceWaveformOptions{})
	})
}
