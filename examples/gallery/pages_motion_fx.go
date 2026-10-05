package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/feedback"
	"github.com/ZacharyZhang-NY/MujicaUI/input"
	"github.com/ZacharyZhang-NY/MujicaUI/navigation"
)

// motionLoaderDemo shows a loader while loading.
func motionLoaderDemo(id string, loader func(c *ui.Context, o feedback.LoaderOptions) *ui.Element) func(c *ui.Context, log func(string)) {
	return func(c *ui.Context, log func(string)) {
		type state struct {
			loading, label bool
			size           float64
		}
		s := ui.Local(c.Root(), id, func() state { return state{loading: true, label: true, size: 20} })
		navigationDemo(c, func() {
			navigationCheck(c, &s.loading, "Loading", log)
			navigationCheck(c, &s.label, "Show label", log)
			navigationSlider(c, &s.size, "Size", 12, 48, log)
		}, func() {
			ui.Row(c).Gap(24).AlignItems(ui.Center).Children(func() {
				if !s.loading {
					ui.Text(c, "Done.")
					return
				}
				loader(c, feedback.LoaderOptions{Size: float32(s.size), Label: "Thinking", ShowLabel: s.label})
				loader(c, feedback.LoaderOptions{Size: float32(s.size), Label: "Indexing", Color: core.Tokens(c).Info})
			})
		})
	}
}

func motionShimmerDemo(c *ui.Context, log func(string)) {
	type state struct{ loading bool }
	s := ui.Local(c.Root(), "292", func() state { return state{loading: true} })
	navigationDemo(c, func() {
		navigationCheck(c, &s.loading, "Loading", log)
	}, func() {
		feedback.Shimmer(c, s.loading, feedback.ShimmerOptions{Height: 56, Label: "Loading report"}, func() {
			ui.Text(c, "Quarterly report: 12 chapters, 340 pages.")
		})
		feedback.Shimmer(c, s.loading, feedback.ShimmerOptions{Width: 240, Height: 20}, func() {
			ui.Text(c, "Updated today").TextColor(core.Tokens(c).TextMuted)
		})
	})
}

func motionCountUpDemo(c *ui.Context, log func(string)) {
	type state struct{ value, duration, decimals float64 }
	s := ui.Local(c.Root(), "293", func() state { return state{value: 1280, duration: 120} })
	navigationDemo(c, func() {
		navigationSlider(c, &s.duration, "Duration (ms)", 120, 1200, log)
		navigationSlider(c, &s.decimals, "Decimals", 0, 4, log)
		if input.Button(c, "Add 500", input.ButtonOptions{Variant: input.Secondary}).Clicked() {
			s.value += 500
			log(fmt.Sprintf("Value %.0f", s.value))
		}
		if input.Button(c, "Random", input.ButtonOptions{Variant: input.Secondary}).Clicked() {
			s.value = float64(rand.Intn(100000))
			log(fmt.Sprintf("Value %.0f", s.value))
		}
	}, func() {
		feedback.CountUp(c, s.value, feedback.CountUpOptions{Duration: time.Duration(s.duration) * time.Millisecond, Decimals: int(s.decimals)}).FontSize(core.FontSize(c, 32))
		feedback.CountUp(c, s.value/1000, feedback.CountUpOptions{Decimals: 2, Duration: time.Second})
	})
}

func motionTypewriterDemo(c *ui.Context, log func(string)) {
	type state struct {
		text  string
		speed float64
		round int
	}
	s := ui.Local(c.Root(), "294", func() state { return state{text: "The court assembles at dusk.", speed: 30} })
	navigationDemo(c, func() {
		navigationSlider(c, &s.speed, "Speed", 5, 120, log)
		if input.Button(c, "Append", input.ButtonOptions{Variant: input.Secondary}).Clicked() {
			s.text += " The music begins."
			log("Appended")
		}
		if input.Button(c, "Restart", input.ButtonOptions{Variant: input.Secondary}).Clicked() {
			s.text = "A new letter arrives from the duchess."
			s.round++
			log("Restarted")
		}
	}, func() {
		ui.Column(c).Key(s.round).Children(func() {
			if feedback.Typewriter(c, s.text, feedback.TypewriterOptions{Speed: s.speed, Skippable: true}).Changed() {
				log("Finished")
			}
		})
	})
}

func motionShimmerTextDemo(c *ui.Context, log func(string)) {
	type state struct{ active bool }
	s := ui.Local(c.Root(), "295", func() state { return state{active: true} })
	navigationDemo(c, func() {
		navigationCheck(c, &s.active, "Thinking", log)
	}, func() {
		feedback.ShimmerText(c, "Thinking about the program…", feedback.ShimmerTextOptions{Active: s.active})
		feedback.ShimmerText(c, "思考中…", feedback.ShimmerTextOptions{Active: s.active})
	})
}

func motionStaggerDemo(c *ui.Context, log func(string)) {
	type state struct {
		round      int
		interval   float64
		horizontal bool
	}
	s := ui.Local(c.Root(), "296", func() state { return state{interval: 60} })
	guests := []string{"Sakiko", "Mutsumi", "Uika", "Umiri", "Nyamu"}
	navigationDemo(c, func() {
		navigationSlider(c, &s.interval, "Interval (ms)", 20, 300, log)
		navigationCheck(c, &s.horizontal, "Horizontal", log)
		if input.Button(c, "Replay", input.ButtonOptions{Variant: input.Secondary}).Clicked() {
			s.round++
			log(fmt.Sprintf("Replay %d", s.round))
		}
	}, func() {
		ui.Column(c).Key(s.round).Children(func() {
			feedback.Stagger(c, len(guests), feedback.StaggerOptions{Interval: time.Duration(s.interval) * time.Millisecond, Horizontal: s.horizontal}, func(i int) {
				ui.Text(c, guests[i]+" arrives").Padding(6, 10).Background(core.Tokens(c).Surface).Radius(4).Margin(0, 0, 4, 0)
			})
		})
	})
}

func motionReorderDemo(c *ui.Context, log func(string)) {
	type state struct {
		order []string
		gap   float64
	}
	s := ui.Local(c.Root(), "297", func() state { return state{order: []string{"Overture", "Masquerade", "Interlude", "Finale"}, gap: 4} })
	navigationDemo(c, func() {
		navigationSlider(c, &s.gap, "Gap", 0, 24, log)
		if input.Button(c, "Shuffle", input.ButtonOptions{Variant: input.Secondary}).Clicked() {
			next := append([]string(nil), s.order...)
			rand.Shuffle(len(next), func(i, j int) { next[i], next[j] = next[j], next[i] })
			s.order = next
			log(fmt.Sprintf("Order %v", s.order))
		}
		if input.Button(c, "Reverse", input.ButtonOptions{Variant: input.Secondary}).Clicked() {
			next := make([]string, len(s.order))
			for i, v := range s.order {
				next[len(next)-1-i] = v
			}
			s.order = next
			log(fmt.Sprintf("Order %v", s.order))
		}
	}, func() {
		ui.Column(c).Width(280).Children(func() {
			feedback.ReorderTransition(c, s.order, feedback.ReorderTransitionOptions{Gap: float32(s.gap)}, func(name string) {
				ui.Text(c, name).Padding(8, 12).Background(core.Tokens(c).Surface).Radius(4)
			})
		})
	})
}

func motionLayoutDemo(c *ui.Context, log func(string)) {
	type state struct{ open, clip bool }
	s := ui.Local(c.Root(), "298", func() state { return state{clip: true} })
	navigationDemo(c, func() {
		navigationCheck(c, &s.open, "Expanded", log)
		navigationCheck(c, &s.clip, "Clip", log)
	}, func() {
		ui.Column(c).AlignItems(ui.Start).Children(func() {
			feedback.LayoutTransition(c, feedback.LayoutTransitionOptions{Clip: s.clip}, func() {
				ui.Column(c).Padding(12).Gap(6).Background(core.Tokens(c).Surface).Radius(4).Children(func() {
					ui.Text(c, "Program").FontWeight(600)
					if s.open {
						for _, p := range []string{"Overture", "Masquerade", "Interlude", "Finale"} {
							ui.Text(c, p)
						}
					}
				})
			})
		})
	})
}

func motionPresenceDemo(c *ui.Context, log func(string)) {
	type state struct {
		show bool
		kind int
	}
	s := ui.Local(c.Root(), "299", func() state { return state{show: true} })
	navigationDemo(c, func() {
		navigationCheck(c, &s.show, "Show", log)
		if navigation.SegmentedControl(c, &s.kind, []navigation.Segment{{Label: "Fade"}, {Label: "Slide"}, {Label: "Collapse"}}, navigation.SegmentedControlOptions{Label: "Kind"}).Changed() {
			log(fmt.Sprintf("Kind %d", s.kind))
		}
	}, func() {
		ui.Column(c).Width(320).Children(func() {
			feedback.Presence(c, s.show, feedback.PresenceOptions{Kind: feedback.PresenceKind(s.kind)}, func() {
				ui.Text(c, "Your changes were saved.").Padding(10, 12).Background(core.Tokens(c).Surface).Radius(4)
			})
			ui.Text(c, "Content below moves with it.").TextColor(core.Tokens(c).TextMuted).Padding(8, 0)
		})
	})
}

func motionBlinkDemo(c *ui.Context, log func(string)) {
	type state struct{ price float64 }
	s := ui.Local(c.Root(), "300", func() state { return state{price: 101.5} })
	navigationDemo(c, func() {
		if input.Button(c, "Tick", input.ButtonOptions{Variant: input.Secondary}).Clicked() {
			s.price += float64(rand.Intn(200)-100) / 50
			log(fmt.Sprintf("Price %.2f", s.price))
		}
	}, func() {
		feedback.BlinkHighlight(c, s.price, feedback.BlinkHighlightOptions{})
	})
}
