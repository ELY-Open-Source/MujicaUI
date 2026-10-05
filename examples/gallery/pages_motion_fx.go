package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/egoist/mygo/ui"

	mujica "mujicaui"
)

// motionLoaderDemo shows a loader while loading.
func motionLoaderDemo(id string, loader func(c *ui.Context, o mujica.LoaderOptions) *ui.Element) func(c *ui.Context, log func(string)) {
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
				loader(c, mujica.LoaderOptions{Size: float32(s.size), Label: "Thinking", ShowLabel: s.label})
				loader(c, mujica.LoaderOptions{Size: float32(s.size), Label: "Indexing", Color: mujica.Tokens(c).Info})
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
		mujica.Shimmer(c, s.loading, mujica.ShimmerOptions{Height: 56, Label: "Loading report"}, func() {
			ui.Text(c, "Quarterly report: 12 chapters, 340 pages.")
		})
		mujica.Shimmer(c, s.loading, mujica.ShimmerOptions{Width: 240, Height: 20}, func() {
			ui.Text(c, "Updated today").TextColor(mujica.Tokens(c).TextMuted)
		})
	})
}

func motionCountUpDemo(c *ui.Context, log func(string)) {
	type state struct{ value, duration, decimals float64 }
	s := ui.Local(c.Root(), "293", func() state { return state{value: 1280, duration: 120} })
	navigationDemo(c, func() {
		navigationSlider(c, &s.duration, "Duration (ms)", 120, 1200, log)
		navigationSlider(c, &s.decimals, "Decimals", 0, 4, log)
		if mujica.Button(c, "Add 500", mujica.ButtonOptions{Variant: mujica.Secondary}).Clicked() {
			s.value += 500
			log(fmt.Sprintf("Value %.0f", s.value))
		}
		if mujica.Button(c, "Random", mujica.ButtonOptions{Variant: mujica.Secondary}).Clicked() {
			s.value = float64(rand.Intn(100000))
			log(fmt.Sprintf("Value %.0f", s.value))
		}
	}, func() {
		mujica.CountUp(c, s.value, mujica.CountUpOptions{Duration: time.Duration(s.duration) * time.Millisecond, Decimals: int(s.decimals)}).FontSize(mujica.FontSize(c, 32))
		mujica.CountUp(c, s.value/1000, mujica.CountUpOptions{Decimals: 2, Duration: time.Second})
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
		if mujica.Button(c, "Append", mujica.ButtonOptions{Variant: mujica.Secondary}).Clicked() {
			s.text += " The music begins."
			log("Appended")
		}
		if mujica.Button(c, "Restart", mujica.ButtonOptions{Variant: mujica.Secondary}).Clicked() {
			s.text = "A new letter arrives from the duchess."
			s.round++
			log("Restarted")
		}
	}, func() {
		ui.Column(c).Key(s.round).Children(func() {
			if mujica.Typewriter(c, s.text, mujica.TypewriterOptions{Speed: s.speed, Skippable: true}).Changed() {
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
		mujica.ShimmerText(c, "Thinking about the program…", mujica.ShimmerTextOptions{Active: s.active})
		mujica.ShimmerText(c, "思考中…", mujica.ShimmerTextOptions{Active: s.active})
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
		if mujica.Button(c, "Replay", mujica.ButtonOptions{Variant: mujica.Secondary}).Clicked() {
			s.round++
			log(fmt.Sprintf("Replay %d", s.round))
		}
	}, func() {
		ui.Column(c).Key(s.round).Children(func() {
			mujica.Stagger(c, len(guests), mujica.StaggerOptions{Interval: time.Duration(s.interval) * time.Millisecond, Horizontal: s.horizontal}, func(i int) {
				ui.Text(c, guests[i]+" arrives").Padding(6, 10).Background(mujica.Tokens(c).Surface).Radius(4).Margin(0, 0, 4, 0)
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
		if mujica.Button(c, "Shuffle", mujica.ButtonOptions{Variant: mujica.Secondary}).Clicked() {
			next := append([]string(nil), s.order...)
			rand.Shuffle(len(next), func(i, j int) { next[i], next[j] = next[j], next[i] })
			s.order = next
			log(fmt.Sprintf("Order %v", s.order))
		}
		if mujica.Button(c, "Reverse", mujica.ButtonOptions{Variant: mujica.Secondary}).Clicked() {
			next := make([]string, len(s.order))
			for i, v := range s.order {
				next[len(next)-1-i] = v
			}
			s.order = next
			log(fmt.Sprintf("Order %v", s.order))
		}
	}, func() {
		ui.Column(c).Width(280).Children(func() {
			mujica.ReorderTransition(c, s.order, mujica.ReorderTransitionOptions{Gap: float32(s.gap)}, func(name string) {
				ui.Text(c, name).Padding(8, 12).Background(mujica.Tokens(c).Surface).Radius(4)
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
			mujica.LayoutTransition(c, mujica.LayoutTransitionOptions{Clip: s.clip}, func() {
				ui.Column(c).Padding(12).Gap(6).Background(mujica.Tokens(c).Surface).Radius(4).Children(func() {
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
		if mujica.SegmentedControl(c, &s.kind, []mujica.Segment{{Label: "Fade"}, {Label: "Slide"}, {Label: "Collapse"}}, mujica.SegmentedControlOptions{Label: "Kind"}).Changed() {
			log(fmt.Sprintf("Kind %d", s.kind))
		}
	}, func() {
		ui.Column(c).Width(320).Children(func() {
			mujica.Presence(c, s.show, mujica.PresenceOptions{Kind: mujica.PresenceKind(s.kind)}, func() {
				ui.Text(c, "Your changes were saved.").Padding(10, 12).Background(mujica.Tokens(c).Surface).Radius(4)
			})
			ui.Text(c, "Content below moves with it.").TextColor(mujica.Tokens(c).TextMuted).Padding(8, 0)
		})
	})
}

func motionBlinkDemo(c *ui.Context, log func(string)) {
	type state struct{ price float64 }
	s := ui.Local(c.Root(), "300", func() state { return state{price: 101.5} })
	navigationDemo(c, func() {
		if mujica.Button(c, "Tick", mujica.ButtonOptions{Variant: mujica.Secondary}).Clicked() {
			s.price += float64(rand.Intn(200)-100) / 50
			log(fmt.Sprintf("Price %.2f", s.price))
		}
	}, func() {
		mujica.BlinkHighlight(c, s.price, mujica.BlinkHighlightOptions{})
	})
}
