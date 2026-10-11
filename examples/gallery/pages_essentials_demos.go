package main

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/overlay"
)

func splitButtonDemo(c *ui.Context, log func(string)) {
	type state struct{ variant, disabled bool }
	s := ui.Local(c.Root(), "481", func() state { return state{} })
	ui.Column(c).Gap(16).FillWidth().Children(func() {
		feedbackSection(c, "Parameters", func() {
			ui.Row(c).Gap(16).Children(func() {
				ui.Checkbox(c, &s.variant, "Secondary")
				ui.Checkbox(c, &s.disabled, "Disabled")
			})
		})
		feedbackSection(c, "Live", func() {
			v := input.Primary
			if s.variant {
				v = input.Secondary
			}
			ui.Row(c).Children(func() {
				r := overlay.SplitButton(c, "Save", overlay.SplitButtonOptions{Variant: v, Disabled: s.disabled}, func(m *overlay.PopupMenu) {
					for _, it := range []string{"Save as…", "Save a copy", "Save all"} {
						if m.Item(it, overlay.PopupMenuItemOptions{}) {
							log("split button: menu " + it)
						}
					}
				})
				if r.Changed() {
					log("split button: main action")
				}
			})
		})
	})
}

func copyButtonDemo(c *ui.Context, log func(string)) {
	type state struct {
		text     string
		disabled bool
	}
	s := ui.Local(c.Root(), "482", func() state { return state{text: "https://example.com/invite/42"} })
	ui.Column(c).Gap(16).FillWidth().Children(func() {
		feedbackSection(c, "Parameters", func() {
			ui.TextInput(c, &s.text).Width(320)
			ui.Checkbox(c, &s.disabled, "Disabled")
		})
		feedbackSection(c, "Live", func() {
			ui.Row(c).Gap(12).Children(func() {
				if input.CopyButton(c, s.text, input.CopyButtonOptions{Disabled: s.disabled}).Changed() {
					log(fmt.Sprintf("copy button: copied %d characters", len([]rune(s.text))))
				}
				if input.CopyButton(c, s.text, input.CopyButtonOptions{Label: "Copy link", Disabled: s.disabled}).Changed() {
					log("copy button: copied with label")
				}
			})
		})
	})
}

func holdToConfirmDemo(c *ui.Context, log func(string)) {
	type state struct {
		seconds  float64
		disabled bool
		deleted  int
	}
	s := ui.Local(c.Root(), "483", func() state { return state{seconds: 1.5} })
	ui.Column(c).Gap(16).FillWidth().Children(func() {
		feedbackSection(c, "Parameters", func() {
			ui.Row(c).Gap(12).AlignItems(ui.Center).Children(func() {
				ui.Text(c, fmt.Sprintf("Hold %.1f s", s.seconds))
				ui.Slider(c, &s.seconds, 0.5, 4).Width(200)
			})
			ui.Checkbox(c, &s.disabled, "Disabled")
		})
		feedbackSection(c, "Live (turn on reduced motion to see the countdown text)", func() {
			ui.Row(c).Children(func() {
				d := time.Duration(s.seconds * float64(time.Second))
				if input.HoldToConfirm(c, "Delete project", input.HoldToConfirmOptions{Duration: d, Disabled: s.disabled}).Changed() {
					s.deleted++
					log("hold to confirm: confirmed")
				}
			})
			ui.Text(c, fmt.Sprintf("Confirmed %d times", s.deleted))
		})
	})
}

func bulkActionBarDemo(c *ui.Context, log func(string)) {
	type state struct {
		picked   [6]bool
		disabled bool
	}
	s := ui.Local(c.Root(), "484", func() state { return state{picked: [6]bool{true, false, true}} })
	ui.Column(c).Gap(16).FillWidth().Children(func() {
		feedbackSection(c, "Rows", func() {
			for i := range s.picked {
				if ui.Checkbox(c, &s.picked[i], fmt.Sprintf("Invoice %d", 1001+i)).Changed() {
					log(fmt.Sprintf("bulk bar: row %d %v", i, s.picked[i]))
				}
			}
			ui.Checkbox(c, &s.disabled, "Disable the bar")
		})
		n := 0
		for _, p := range s.picked {
			if p {
				n++
			}
		}
		r := input.BulkActionBar(c, n, input.BulkActionBarOptions{Disabled: s.disabled}, func() {
			if input.Button(c, "Archive", input.ButtonOptions{Variant: input.Secondary}).Clicked() {
				log(fmt.Sprintf("bulk bar: archive %d", n))
			}
			if input.Button(c, "Delete", input.ButtonOptions{Variant: input.Danger}).Clicked() {
				log(fmt.Sprintf("bulk bar: delete %d", n))
			}
		})
		if r.Changed() {
			s.picked = [6]bool{}
			log("bulk bar: cleared")
		}
	})
}

func scrubInputDemo(c *ui.Context, log func(string)) {
	type state struct {
		opacity, size float64
		disabled      bool
	}
	s := ui.Local(c.Root(), "485", func() state { return state{opacity: 0.6, size: 14} })
	ui.Column(c).Gap(16).FillWidth().Children(func() {
		ui.Checkbox(c, &s.disabled, "Disabled")
		feedbackSection(c, "Drag a label sideways, or type", func() {
			if input.ScrubInput(c, "Opacity", &s.opacity, input.ScrubInputOptions{Max: 1, Step: 0.05, Disabled: s.disabled}).Changed() {
				log(fmt.Sprintf("scrub: opacity %.2f", s.opacity))
			}
			if input.ScrubInput(c, "Size", &s.size, input.ScrubInputOptions{Min: 8, Max: 72, Disabled: s.disabled}).Changed() {
				log(fmt.Sprintf("scrub: size %.0f", s.size))
			}
		})
	})
}

func currencyInputDemo(c *ui.Context, log func(string)) {
	type state struct {
		cents, yen int64
		disabled   bool
	}
	s := ui.Local(c.Root(), "486", func() state { return state{cents: 123450, yen: 98000} })
	ui.Column(c).Gap(16).FillWidth().Children(func() {
		ui.Checkbox(c, &s.disabled, "Disabled")
		feedbackSection(c, "Live", func() {
			ui.Row(c).Gap(12).Children(func() {
				if input.CurrencyInput(c, &s.cents, input.CurrencyInputOptions{Symbol: "$", Label: "Price", Disabled: s.disabled}).Changed() {
					log(fmt.Sprintf("currency: %d cents", s.cents))
				}
				if input.CurrencyInput(c, &s.yen, input.CurrencyInputOptions{Symbol: "¥", NoFraction: true, Label: "Yen", Disabled: s.disabled}).Changed() {
					log(fmt.Sprintf("currency: %d yen", s.yen))
				}
			})
			ui.Text(c, "Stored as "+input.CurrencyFormat(s.cents, 2, true)+" and "+input.CurrencyFormat(s.yen, 0, true))
		})
	})
}

func maskedInputDemo(c *ui.Context, log func(string)) {
	type state struct {
		mask     int
		value    string
		disabled bool
	}
	masks := []string{"(999) 999-9999", "9999 9999 9999 9999", "99/99/9999", "aa-9999"}
	s := ui.Local(c.Root(), "487", func() state { return state{} })
	ui.Column(c).Gap(16).FillWidth().Children(func() {
		feedbackSection(c, "Parameters", func() {
			if ui.Segmented(c, &s.mask, masks...).Changed() {
				s.value = ""
				log("masked: mask " + masks[s.mask])
			}
			ui.Checkbox(c, &s.disabled, "Disabled")
		})
		feedbackSection(c, "Live (try an input method: nothing commits until it does)", func() {
			r := input.MaskedInput(c, &s.value, input.MaskedInputOptions{Mask: masks[s.mask], Label: "Masked field", Disabled: s.disabled})
			r.Element.Width(260)
			if r.Changed() {
				log(fmt.Sprintf("masked: %d characters filled", len([]rune(s.value))))
			}
			if r.Submitted() {
				log("masked: submitted")
			}
		})
	})
}

func choiceChipsDemo(c *ui.Context, log func(string)) {
	type state struct {
		value              []string
		multiple, disabled bool
	}
	s := ui.Local(c.Root(), "489", func() state { return state{multiple: true, value: []string{"open"}} })
	items := []input.ToggleItem[string]{{Value: "open", Label: "Open"}, {Value: "draft", Label: "Draft"}, {Value: "review", Label: "In review"}, {Value: "closed", Label: "Closed", Disabled: true}, {Value: "merged", Label: "Merged"}}
	ui.Column(c).Gap(16).FillWidth().Children(func() {
		feedbackSection(c, "Parameters", func() {
			ui.Row(c).Gap(16).Children(func() {
				if ui.Checkbox(c, &s.multiple, "Multiple").Changed() && !s.multiple && len(s.value) > 1 {
					s.value = s.value[:1]
				}
				ui.Checkbox(c, &s.disabled, "Disabled")
			})
		})
		feedbackSection(c, "Live", func() {
			if input.ChoiceChips(c, &s.value, items, input.ChoiceChipsOptions{Multiple: s.multiple, Disabled: s.disabled, Label: "Status"}).Changed() {
				log(fmt.Sprintf("chips: %v", s.value))
			}
		})
	})
}

func knobDemo(c *ui.Context, log func(string)) {
	type state struct {
		gain, mix float64
		disabled  bool
	}
	s := ui.Local(c.Root(), "490", func() state { return state{gain: 35, mix: 0.5} })
	ui.Column(c).Gap(16).FillWidth().Children(func() {
		ui.Checkbox(c, &s.disabled, "Disabled")
		ui.Row(c).Gap(32).Children(func() {
			if input.Knob(c, &s.gain, input.KnobOptions{Label: "Gain", Ticks: 11, Disabled: s.disabled}).Changed() {
				log(fmt.Sprintf("knob: gain %.1f", s.gain))
			}
			if input.Knob(c, &s.mix, input.KnobOptions{Max: 1, Step: 0.1, Size: 72, Label: "Mix", Disabled: s.disabled}).Changed() {
				log(fmt.Sprintf("knob: mix %.1f", s.mix))
			}
		})
	})
}

func verticalSliderDemo(c *ui.Context, log func(string)) {
	type state struct {
		levels             [4]float64
		showValue, disable bool
	}
	s := ui.Local(c.Root(), "491", func() state { return state{levels: [4]float64{40, 70, 20, 90}, showValue: true} })
	ui.Column(c).Gap(16).FillWidth().Children(func() {
		ui.Row(c).Gap(16).Children(func() {
			ui.Checkbox(c, &s.showValue, "Show value")
			ui.Checkbox(c, &s.disable, "Disabled")
		})
		ui.Row(c).Gap(16).Children(func() {
			for i := range s.levels {
				opts := input.VerticalSliderOptions{Step: 10, ShowValue: s.showValue, Disabled: s.disable, Label: fmt.Sprintf("Band %d", i+1)}
				if i == 0 {
					opts.Marks = []input.SliderMark{{Value: 0, Label: "0"}, {Value: 50, Label: "50"}, {Value: 100, Label: "100"}}
				}
				if input.VerticalSlider(c, &s.levels[i], opts).Changed() {
					log(fmt.Sprintf("vertical slider: band %d = %.0f", i+1, s.levels[i]))
				}
			}
		})
	})
}
