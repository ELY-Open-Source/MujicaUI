package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/egoist/mygo/ui"

	mujica "mujicaui"
)

func dataCardDemo(c *ui.Context, log func(string)) {
	type state struct {
		decorated, footer bool
		title             string
	}
	s := ui.Local(c.Root(), "069", func() state { return state{footer: true, title: "Archive"} })
	dataSection(c, "Card", func() {
		dataControls(c, func() {
			ui.TextInput(c, &s.title).Width(180).Placeholder("Title")
			ui.Checkbox(c, &s.decorated, "Decorated")
			ui.Checkbox(c, &s.footer, "Footer")
		})
		card := func(decorated bool) {
			opts := mujica.CardOptions{Title: s.title, Decorated: decorated}
			if s.footer {
				opts.Footer = func() {
					if mujica.Button(c, "Open", mujica.ButtonOptions{Variant: mujica.Secondary}).Clicked() {
						log(fmt.Sprintf("open clicked (decorated=%v)", decorated))
					}
				}
			}
			mujica.Card(c, opts, func() {
				ui.Text(c, "Letters of the court, sealed in wine-red wax and kept in the north vault.")
			}).Width(280)
		}
		dataStates(c, func(state func(string, func())) {
			state("Options above", func() { card(s.decorated) })
			state("Plain", func() { card(false) })
			state("Decorated", func() { card(true) })
		})
	})
}

func dataAccordionDemo(c *ui.Context, log func(string)) {
	type state struct {
		multiple, keepOne bool
		open              []string
	}
	s := ui.Local(c.Root(), "074", func() state { return state{open: []string{"choir"}} })
	dataSection(c, "Accordion", func() {
		dataControls(c, func() {
			ui.Checkbox(c, &s.multiple, "Multiple")
			ui.Checkbox(c, &s.keepOne, "Keep one open")
			if mujica.Button(c, "Close all", mujica.ButtonOptions{Variant: mujica.Ghost}).Clicked() {
				s.open = nil
				log("closed all from outside")
			}
		})
		if mujica.Accordion(c, &s.open, []mujica.AccordionItem{
			{ID: "choir", Title: "Choir", Content: func() { ui.Text(c, "Twelve voices, two at the altar.") }},
			{ID: "nave", Title: "Nave", Content: func() { ui.Text(c, "Pews of black oak, polished for the feast.") }},
			{ID: "crypt", Title: "Crypt (disabled)", Disabled: true, Content: func() {}},
			{ID: "bells", Title: "Bells", Content: func() {
				if mujica.Button(c, "Ring", mujica.ButtonOptions{}).Clicked() {
					log("ring the bells")
				}
			}},
		}, mujica.AccordionOptions{Multiple: s.multiple, KeepOne: s.keepOne}).Width(420).Changed() {
			log(fmt.Sprintf("open %v", s.open))
		}
	})
}

func dataCollapsibleDemo(c *ui.Context, log func(string)) {
	type state struct{ open, disabled bool }
	s := ui.Local(c.Root(), "075", func() state { return state{} })
	dataSection(c, "Collapsible", func() {
		dataControls(c, func() {
			if mujica.Button(c, "Open from outside", mujica.ButtonOptions{Variant: mujica.Secondary}).Clicked() {
				s.open = true
				log("opened from outside")
			}
			ui.Checkbox(c, &s.disabled, "Disabled")
		})
		if mujica.Collapsible(c, &s.open, mujica.CollapsibleOptions{Title: "Rites of the evening", Disabled: s.disabled}, func() {
			ui.Text(c, "Candles are lit at the seventh bell; the choir enters from the north door.")
			if mujica.Button(c, "Light candles", mujica.ButtonOptions{Variant: mujica.Outline}).Clicked() {
				log("light candles")
			}
		}).Changed() {
			log(fmt.Sprintf("open=%v", s.open))
		}
	})
}

func dataTimelineDemo(c *ui.Context, log func(string)) {
	type state struct{ items []mujica.TimelineItem }
	s := ui.Local(c.Root(), "076", func() state {
		return state{items: []mujica.TimelineItem{
			{Time: "06:00", Title: "Matins sung", Content: "The choir gathered before dawn.", Status: mujica.TimelineDone},
			{Time: "09:00", Title: "Relics moved", Status: mujica.TimelineActive,
				Content: "The reliquaries travel from the crypt to the high altar, carried in procession by four acolytes in black and wine."},
			{Time: "12:00", Title: "Bell cracked", Status: mujica.TimelineError},
			{Time: "15:00", Title: "Storm warning", Status: mujica.TimelineWarning},
			{Time: "18:00", Title: "Vespers"},
		}}
	})
	dataSection(c, "Timeline", func() {
		dataControls(c, func() {
			if mujica.Button(c, "Add event", mujica.ButtonOptions{Variant: mujica.Secondary}).Clicked() {
				s.items = append(s.items, mujica.TimelineItem{Time: fmt.Sprintf("%02d:00", 19+len(s.items)%5), Title: "Compline"})
				log(fmt.Sprintf("%d events", len(s.items)))
			}
		})
		ui.Column(c).Width(420).Children(func() { mujica.Timeline(c, s.items) })
	})
}

func dataDescriptionDemo(c *ui.Context, log func(string)) {
	type state struct {
		columns string
		width   float64
	}
	s := ui.Local(c.Root(), "077", func() state { return state{columns: "2", width: 560} })
	dataSection(c, "DescriptionList", func() {
		dataControls(c, func() {
			ui.Select(c, &s.columns, []string{"1", "2", "3"})
			ui.Text(c, "Width")
			ui.Slider(c, &s.width, 240, 720).Width(200)
		})
		cols, _ := strconv.Atoi(s.columns)
		ui.Column(c).Width(float32(s.width)).Children(func() {
			if i, ok := mujica.DescriptionList(c, []mujica.DescriptionItem{
				{Label: "Archivist", Value: "Sister Agnès", Copyable: true},
				{Label: "Shelf", Value: "North vault, row 7"},
				{Label: "Catalogue", Value: "MS-1347-B", Copyable: true},
				{Label: "Seal", Value: "Wine-red wax, cracked at the rim", Span: -1},
			}, mujica.DescriptionListOptions{Columns: cols}).Copied(); ok {
				log(fmt.Sprintf("copied item %d", i))
			}
		})
	})
}

func dataStatisticDemo(c *ui.Context, log func(string)) {
	type state struct {
		value, change float64
		precision     string
	}
	s := ui.Local(c.Root(), "078", func() state { return state{value: 12840, change: 12.5, precision: "0"} })
	dataSection(c, "Statistic", func() {
		dataControls(c, func() {
			ui.Text(c, "Value")
			if ui.Slider(c, &s.value, 0, 100000).Width(160).Changed() {
				log(fmt.Sprintf("value %.0f", s.value))
			}
			ui.Text(c, "Change")
			ui.Slider(c, &s.change, -50, 50).Width(160)
			ui.Select(c, &s.precision, []string{"0", "1", "2"})
		})
		p, _ := strconv.Atoi(s.precision)
		change := s.change
		dataStates(c, func(state func(string, func())) {
			state("With change", func() {
				mujica.Statistic(c, s.value, mujica.StatisticOptions{Title: "Pilgrims", Unit: "souls", Precision: p,
					Change: &change, ChangeUnit: "%", ChangeNote: "vs last feast"})
			})
			state("Prefix", func() {
				mujica.Statistic(c, s.value/1000, mujica.StatisticOptions{Title: "Tithe", Prefix: "£", Suffix: "k", Precision: p})
			})
		})
	})
}

var dataPictures = func() []mujica.ViewerImage {
	art := func(bg, fg string, shape string) *ui.SVG {
		return ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" width="480" height="320"><rect width="480" height="320" fill="` + bg + `"/>` +
			`<path d="` + shape + `" fill="` + fg + `"/></svg>`))
	}
	return []mujica.ViewerImage{
		{Source: art("#AD9164", "#7D2034", "M240 40 L380 160 L240 280 L100 160 Z"), Alt: "Rose window"},
		{Source: art("#211C1D", "#BDA172", "M80 280 L80 120 L240 40 L400 120 L400 280 Z"), Alt: "Chapel arch"},
		{Source: art("#F2EBDD", "#651729", "M60 60 H420 V100 H60 Z M60 140 H420 V180 H60 Z M60 220 H420 V260 H60 Z"), Alt: "Choir stalls"},
	}
}()

func dataViewerDemo(c *ui.Context, log func(string)) {
	type state struct {
		open  bool
		index int
	}
	s := ui.Local(c.Root(), "079", func() state { return state{} })
	dataSection(c, "ImageViewer", func() {
		k := mujica.Tokens(c)
		ui.Text(c, "Choose a picture to open the viewer.").TextColor(k.TextMuted)
		ui.Row(c).Gap(12).Children(func() {
			for i, p := range dataPictures {
				thumb := ui.ButtonBase(c).Label("Open "+p.Alt).Radius(4).Border(1, k.Border).Clip()
				thumb.Children(func() { ui.Image(c, p.Source).Size(144, 96) })
				if thumb.Clicked() {
					s.open, s.index = true, i
					log("open " + p.Alt)
				}
			}
		})
		wasOpen := s.open
		if mujica.ImageViewer(c, &s.open, &s.index, dataPictures).Changed() {
			log(fmt.Sprintf("image %d", s.index+1))
		}
		if wasOpen && !s.open {
			log("closed")
		}
	})
}

func dataCarouselDemo(c *ui.Context, log func(string)) {
	type state struct {
		index          int
		autoplay, loop bool
		interval       string
	}
	s := ui.Local(c.Root(), "080", func() state { return state{interval: "3s"} })
	dataSection(c, "Carousel", func() {
		dataControls(c, func() {
			ui.Checkbox(c, &s.autoplay, "Autoplay")
			ui.Checkbox(c, &s.loop, "Loop")
			ui.Select(c, &s.interval, []string{"2s", "3s", "5s"})
		})
		d, err := time.ParseDuration(s.interval)
		if err != nil {
			panic(err)
		}
		v := mujica.Carousel(c, &s.index, len(dataPictures), mujica.CarouselOptions{Label: "Windows", Autoplay: s.autoplay, Loop: s.loop, Interval: d}, func(i int) {
			ui.Image(c, dataPictures[i].Source).Grow(1).Label(dataPictures[i].Alt)
		})
		v.Size(520, 300)
		if v.Changed() {
			log(fmt.Sprintf("slide %d", s.index+1))
		}
	})
}
