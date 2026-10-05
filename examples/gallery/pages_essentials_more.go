package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"

	mujica "mujicaui"
)

func mentionInputDemo(c *ui.Context, log func(string)) {
	type state struct {
		text     string
		mentions []mujica.Mention
		disabled bool
	}
	people := []mujica.MentionCandidate{{ID: "u1", Name: "Alice"}, {ID: "u2", Name: "Albert"}, {ID: "u3", Name: "Bob"}, {ID: "u4", Name: "Sakiko"}, {ID: "u5", Name: "Mutsumi"}}
	s := ui.Local(c.Root(), "488", func() state { return state{} })
	ui.Column(c).Gap(16).FillWidth().Children(func() {
		ui.Checkbox(c, &s.disabled, "Disabled")
		feedbackSection(c, "Type @ to mention someone", func() {
			r := mujica.MentionInput(c, &s.text, &s.mentions, mujica.MentionInputOptions{
				Placeholder: "Write a comment…", Label: "Comment", Disabled: s.disabled,
				Candidates: func(q string) []mujica.MentionCandidate {
					var out []mujica.MentionCandidate
					for _, p := range people {
						if strings.HasPrefix(strings.ToLower(p.Name), strings.ToLower(q)) {
							out = append(out, p)
						}
					}
					return out
				},
			})
			r.Element.Width(420)
			if r.Changed() {
				log(fmt.Sprintf("mention: %d mentions", len(s.mentions)))
			}
			if r.Submitted() {
				log("mention: submitted")
			}
			for _, m := range s.mentions {
				ui.Text(c, fmt.Sprintf("@%s → %s", m.Name, m.ID)).Key(m.ID)
			}
		})
	})
}

func colorPaletteDemo(c *ui.Context, log func(string)) {
	type state struct {
		color            ui.Color
		custom, disabled bool
	}
	s := ui.Local(c.Root(), "492", func() state { return state{color: ui.Hex("#7D2034"), custom: true} })
	swatches := []mujica.ColorSwatch{
		{Color: ui.Hex("#7D2034"), Name: "Wine"}, {Color: ui.Hex("#A12635"), Name: "Crimson"}, {Color: ui.Hex("#AD9164"), Name: "Old gold"},
		{Color: ui.Hex("#805714"), Name: "Bronze"}, {Color: ui.Hex("#326247"), Name: "Ivy"}, {Color: ui.Hex("#355E80"), Name: "Slate blue"},
		{Color: ui.Hex("#4B3A5A"), Name: "Plum"}, {Color: ui.Hex("#211C1D"), Name: "Ink"}, {Color: ui.Hex("#6A5B55"), Name: "Taupe"},
		{Color: ui.Hex("#D2C3AF"), Name: "Parchment"}, {Color: ui.Hex("#F2EBDD"), Name: "Ivory"}, {Color: ui.Hex("#FFFFFF"), Name: "White"},
	}
	ui.Column(c).Gap(16).FillWidth().Children(func() {
		ui.Row(c).Gap(16).Children(func() {
			ui.Checkbox(c, &s.custom, "Custom entry")
			ui.Checkbox(c, &s.disabled, "Disabled")
		})
		if mujica.ColorPalette(c, &s.color, swatches, mujica.ColorPaletteOptions{Columns: 6, Custom: s.custom, Disabled: s.disabled, Label: "Tag color"}).Changed() {
			log(fmt.Sprintf("palette: %v", s.color))
		}
	})
}

func signaturePadDemo(c *ui.Context, log func(string)) {
	type state struct {
		strokes  [][]mujica.SignaturePoint
		disabled bool
	}
	s := ui.Local(c.Root(), "493", func() state { return state{} })
	ui.Column(c).Gap(16).Width(480).Children(func() {
		ui.Checkbox(c, &s.disabled, "Disabled")
		if mujica.SignaturePad(c, &s.strokes, mujica.SignaturePadOptions{Disabled: s.disabled}).Changed() {
			log(fmt.Sprintf("signature: %d strokes, %d bytes of path", len(s.strokes), len(mujica.SignaturePath(s.strokes))))
		}
	})
}

func bannerDemo(c *ui.Context, log func(string)) {
	type state struct {
		sev                  int
		action, closable, op bool
	}
	s := ui.Local(c.Root(), "494", func() state { return state{sev: 2, action: true, closable: true, op: true} })
	ui.Column(c).Gap(16).FillWidth().Children(func() {
		feedbackSection(c, "Parameters", func() {
			ui.Segmented(c, &s.sev, "Info", "Success", "Warning", "Error")
			ui.Row(c).Gap(16).Children(func() {
				ui.Checkbox(c, &s.action, "Action")
				ui.Checkbox(c, &s.closable, "Closable")
			})
		})
		opts := mujica.BannerOptions{Severity: feedbackSeverities[s.sev]}
		if s.action {
			opts.Action = "Details"
		}
		var open *bool
		if s.closable {
			open = &s.op
		}
		r := mujica.Banner(c, open, "Scheduled maintenance tonight from 22:00 to 23:00.", opts)
		if r.Action {
			log("banner: action")
		}
		if r.Closed {
			log("banner: closed")
		}
		if !s.op && mujica.Button(c, "Show again", mujica.ButtonOptions{Variant: mujica.Secondary}).Clicked() {
			s.op = true
			log("banner: reopened")
		}
	})
}

func undoToastDemo(c *ui.Context, log func(string)) {
	type state struct {
		notes   []string
		removed string
		open    bool
		seconds float64
	}
	s := ui.Local(c.Root(), "495", func() state {
		return state{notes: []string{"Call the florist", "Book the hall", "Order candles"}, seconds: 5}
	})
	ui.Column(c).Gap(16).FillWidth().Children(func() {
		ui.Row(c).Gap(12).AlignItems(ui.Center).Children(func() {
			ui.Text(c, fmt.Sprintf("Countdown %.0f s", s.seconds))
			ui.Slider(c, &s.seconds, 2, 10).Width(200)
		})
		for i, n := range s.notes {
			ui.Row(c).Key(n).Gap(12).AlignItems(ui.Center).Children(func() {
				ui.Text(c, n).Grow(1)
				if mujica.Button(c, "Delete", mujica.ButtonOptions{Variant: mujica.Ghost}).Clicked() {
					s.removed, s.open = n, true
					s.notes = append(s.notes[:i:i], s.notes[i+1:]...)
					log("undo toast: deleted " + n)
				}
			})
		}
	})
	if !s.open {
		return
	}
	removed := s.removed
	r := mujica.UndoToast(c, &s.open, "Deleted “"+removed+"”", mujica.UndoToastOptions{
		Duration: time.Duration(s.seconds * float64(time.Second)),
		OnUndo:   func() { s.notes = append(s.notes, removed) },
	})
	if r.Undone {
		log("undo toast: undone " + removed)
	}
	if r.Expired {
		log("undo toast: expired")
	}
}

func masonryDemo(c *ui.Context, log func(string)) {
	type state struct {
		count, width float64
	}
	s := ui.Local(c.Root(), "496", func() state { return state{count: 40, width: 180} })
	k := mujica.Tokens(c)
	ui.Column(c).Gap(16).FillWidth().Children(func() {
		ui.Row(c).Gap(12).AlignItems(ui.Center).Children(func() {
			ui.Text(c, fmt.Sprintf("Items %.0f", s.count))
			if ui.Slider(c, &s.count, 1, 10000).Width(200).Changed() {
				log(fmt.Sprintf("masonry: %d items", int(s.count)))
			}
			ui.Text(c, fmt.Sprintf("Column %.0f", s.width))
			ui.Slider(c, &s.width, 120, 320).Width(160)
		})
		mujica.Masonry(c, int(s.count), mujica.MasonryOptions{Key: strconv.Itoa, ColumnWidth: float32(s.width)}, func(i int) {
			card := ui.Column(c).Height(float32(60+i*37%90)).Padding(10).Background(k.Surface).Border(1, k.Border).Radius(4).Cursor(ui.CursorPointer)
			card.Children(func() { ui.Text(c, "Card "+strconv.Itoa(i+1)) })
			if card.Clicked() {
				log("masonry: clicked card " + strconv.Itoa(i+1))
			}
		}).FillWidth().Height(420)
	})
}

func resizablePanelGroupDemo(c *ui.Context, log func(string)) {
	type state struct {
		sizes              []float32
		vertical, disabled bool
	}
	s := ui.Local(c.Root(), "497", func() state { return state{} })
	k := mujica.Tokens(c)
	pane := func(name string) func() {
		return func() {
			ui.Column(c).Fill().Padding(12).Background(k.Surface).Children(func() { ui.Text(c, name) })
		}
	}
	ui.Column(c).Gap(16).FillWidth().Children(func() {
		ui.Row(c).Gap(16).Children(func() {
			if ui.Checkbox(c, &s.vertical, "Vertical").Changed() {
				s.sizes = nil
			}
			ui.Checkbox(c, &s.disabled, "Disabled")
		})
		r := mujica.ResizablePanelGroup(c, &s.sizes, []mujica.ResizablePanel{
			{Label: "Files", Collapsible: true, Build: pane("Files")},
			{Label: "Editor", Min: 160, Build: pane("Editor")},
			{Label: "Outline", Collapsible: true, Build: pane("Outline")},
		}, mujica.ResizablePanelGroupOptions{Vertical: s.vertical, Disabled: s.disabled})
		r.Element.FillWidth().Height(320)
		if r.Committed {
			log(fmt.Sprintf("panels: save %.2f", s.sizes))
		}
	})
}

func dockPanelDemo(c *ui.Context, log func(string)) {
	type state struct {
		layout   mujica.DockLayout
		disabled bool
	}
	s := ui.Local(c.Root(), "498", func() state {
		return state{layout: mujica.DockLayout{Zones: [4][]string{{"editor"}, {"files"}, {"outline"}, {"terminal", "problems"}}}}
	})
	var items []mujica.DockItem
	for _, id := range []string{"files", "editor", "outline", "terminal", "problems"} {
		items = append(items, mujica.DockItem{ID: id, Title: strings.ToUpper(id[:1]) + id[1:], Build: func() { ui.Text(c, "Contents of "+id) }})
	}
	ui.Column(c).Gap(16).FillWidth().Children(func() {
		ui.Checkbox(c, &s.disabled, "Disabled")
		r := mujica.DockPanel(c, &s.layout, items, mujica.DockPanelOptions{Disabled: s.disabled})
		r.Element.FillWidth().Height(440)
		if r.Changed() {
			log(fmt.Sprintf("dock: save %v", s.layout.Zones))
		}
	})
}

func qrCodeDemo(c *ui.Context, log func(string)) {
	type state struct {
		text  string
		level int
		size  float64
	}
	s := ui.Local(c.Root(), "499", func() state { return state{text: "https://example.com/mujica", level: 1, size: 200} })
	levels := []mujica.QRLevel{mujica.QRLow, mujica.QRMedium, mujica.QRQuartile, mujica.QRHigh}
	ui.Column(c).Gap(16).FillWidth().Children(func() {
		feedbackSection(c, "Parameters", func() {
			if ui.TextInput(c, &s.text).Width(360).Changed() {
				log(fmt.Sprintf("qr: %d bytes", len(s.text)))
			}
			if ui.Segmented(c, &s.level, "L", "M", "Q", "H").Changed() {
				log("qr: level " + levels[s.level].String())
			}
			ui.Row(c).Gap(12).AlignItems(ui.Center).Children(func() {
				ui.Text(c, fmt.Sprintf("Size %.0f", s.size))
				ui.Slider(c, &s.size, 100, 320).Width(200)
			})
		})
		mujica.QRCode(c, s.text, mujica.QRCodeOptions{Level: levels[s.level], Size: float32(s.size)})
		if m, err := mujica.EncodeQR([]byte(s.text), levels[s.level]); err == nil {
			ui.Text(c, fmt.Sprintf("Version %d, mask %d, %d × %d modules", m.Version, m.Mask, m.Size(), m.Size()))
		}
	})
}

func activityFeedDemo(c *ui.Context, log func(string)) {
	type state struct{ events []mujica.Activity }
	s := ui.Local(c.Root(), "500", func() state {
		now := time.Now()
		var ev []mujica.Activity
		actors := []string{"Sakiko", "Mutsumi", "Uika", "Umiri", "Nyamu"}
		verbs := []string{"commented on", "edited", "approved", "closed", "opened"}
		for i := range 400 {
			ev = append(ev, mujica.Activity{ID: strconv.Itoa(i), Actor: actors[i%5], Action: verbs[i*3%5], Object: fmt.Sprintf("Score #%d", 400-i), Time: now.Add(-time.Duration(i) * 97 * time.Minute)})
		}
		return state{events: ev}
	})
	ui.Column(c).Gap(16).FillWidth().Children(func() {
		if mujica.Button(c, "Add event", mujica.ButtonOptions{Variant: mujica.Secondary}).AlignSelf(ui.Start).Clicked() {
			e := mujica.Activity{ID: fmt.Sprint("n", len(s.events)), Actor: "You", Action: "added", Object: "a new note", Time: time.Now()}
			s.events = append([]mujica.Activity{e}, s.events...)
			log("activity: added " + e.ID)
		}
		mujica.ActivityFeed(c, s.events, mujica.ActivityFeedOptions{}).FillWidth().Height(420)
	})
}
