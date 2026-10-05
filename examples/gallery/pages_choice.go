package main

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"

	mujica "mujicaui"
	"mujicaui/icons"
	"mujicaui/theme"
)

const choiceCategory = "Choices"

var choicePages = []page{
	{ID: "031", Name: "Switch", ZhName: "开关", Category: choiceCategory, Demo: switchDemo, Code: `wifi := true
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	if mujica.Switch(c, &wifi, "Wi-Fi", mujica.SwitchOptions{}).Changed() {
		log.Println("Wi-Fi:", wifi)
	}
}`, Keys: "Tab focuses the switch; Space turns it over."},
	{ID: "032", Name: "Toggle", ZhName: "切换按钮", Category: choiceCategory, Demo: toggleDemo, Code: `bold := false
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.Toggle(c, &bold, "", mujica.ToggleOptions{Icon: icons.Must("bold"), Label: "Bold"})
}`, Keys: "Tab focuses the toggle; Space or Enter presses or releases it."},
	{ID: "033", Name: "ToggleGroup", ZhName: "切换按钮组", Category: choiceCategory, Demo: toggleGroupDemo, Code: `align := []string{"left"}
items := []mujica.ToggleItem[string]{
	{Value: "left", Label: "Left"},
	{Value: "center", Label: "Center"},
	{Value: "right", Label: "Right"},
}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.ToggleGroup(c, &align, items, mujica.ToggleGroupOptions{Label: "Align"})
}`, Keys: "Tab enters the group at the pressed item; Left/Right (Up/Down when vertical), Home and End move; Space presses."},
	{ID: "034", Name: "Slider", ZhName: "滑块", Category: choiceCategory, Demo: sliderDemo, Code: `volume := 40.0
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.Slider(c, &volume, mujica.SliderOptions{Step: 5, ShowValue: true, Label: "Volume"})
}`, Keys: "Arrows move one step; Home and End go to the ends."},
	{ID: "035", Name: "RangeSlider", ZhName: "范围滑块", Category: choiceCategory, Demo: rangeSliderDemo, Code: `low, high := 20.0, 80.0
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.RangeSlider(c, &low, &high, mujica.RangeSliderOptions{Step: 10, Label: "Price"})
}`, Keys: "Tab moves between the two thumbs; arrows move one step; Home and End go as far as the other thumb allows."},
	{ID: "036", Name: "Select", ZhName: "选择器", Category: choiceCategory, Demo: selectDemo, Code: `plan := ""
options := []mujica.SelectOption[string]{
	{Value: "free", Label: "Free"},
	{Value: "pro", Label: "Pro"},
	{Value: "team", Label: "Team", Disabled: true},
}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.Select(c, &plan, options, mujica.SelectOptions{Clearable: true, Label: "Plan"})
}`, Keys: "Enter, Space, Up or Down open; arrows, Home, End move; Enter chooses; Escape closes; typing letters on the closed select jumps to a match; Backspace/Delete clear when Clearable."},
	{ID: "037", Name: "MultiSelect", ZhName: "多选选择器", Category: choiceCategory, Demo: multiSelectDemo, Code: `tags := []int{1}
options := []mujica.SelectOption[int]{
	{Value: 1, Label: "Rose"},
	{Value: 2, Label: "Lily"},
	{Value: 3, Label: "Iris"},
}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.MultiSelect(c, &tags, options, mujica.MultiSelectOptions{Max: 2, Label: "Flowers"})
}`, Keys: "Enter, Space or Down open with the search focused; Down enters the list; Up/Down move; Space checks; Escape closes; Backspace on the field removes the last tag."},
	{ID: "038", Name: "Combobox", ZhName: "组合框", Category: choiceCategory, Demo: comboboxDemo, Code: `font := ""
fonts := []string{"Garamond", "Georgia", "Didot"}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	r := mujica.Combobox(c, &font, fonts, mujica.ComboboxOptions{AllowCustom: true, Label: "Font"})
	if r.QueryChanged() {
		log.Println("search:", r.Query())
	}
}`, Keys: "Typing filters; Down opens or moves; Up moves; Enter chooses; Escape closes; Tab leaves (keeping typed text with AllowCustom)."},
	{ID: "039", Name: "Cascader", ZhName: "级联选择", Category: choiceCategory, Demo: cascaderDemo, Code: `var path []string
level := func(p []string) mujica.CascaderLevel {
	if len(p) == 0 {
		return mujica.CascaderLevel{Nodes: []mujica.CascaderNode{{Value: "fr", Label: "France"}}}
	}
	return mujica.CascaderLevel{Nodes: []mujica.CascaderNode{{Value: "paris", Label: "Paris", Leaf: true}}}
}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.Cascader(c, &path, mujica.CascaderOptions{Level: level, Label: "City"})
}`, Keys: "Enter, Space or Down open; Up/Down move in a column; Right opens a branch; Left goes back; Enter chooses; Escape closes."},
	{ID: "040", Name: "TreeSelect", ZhName: "树选择", Category: choiceCategory, Demo: treeSelectDemo, Code: `var roles []string
nodes := []mujica.TreeSelectNode{
	{ID: "court", Label: "Court", Children: []mujica.TreeSelectNode{
		{ID: "king", Label: "King"},
		{ID: "queen", Label: "Queen"},
	}},
}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.TreeSelect(c, &roles, nodes, mujica.TreeSelectOptions{Multiple: true, Label: "Roles"})
}`, Keys: "Enter, Space or Down open with the search focused; Down enters the tree; Up/Down move; Right expands or enters; Left collapses or goes to the parent; Space checks; Escape closes."},
}

// choiceDemo lays out a demo: parameter controls, then the samples.
func choiceDemo(c *ui.Context, controls, samples func()) {
	k := mujica.Tokens(c)
	ui.Column(c).Gap(20).AlignItems(ui.Start).Children(func() {
		ui.Text(c, "Parameters").Font(theme.SerifFont).FontSize(theme.H3Size).TextColor(k.Text)
		ui.Row(c).Gap(20).Wrap().AlignItems(ui.Center).Children(controls)
		ui.Box(c).Height(1).FillWidth().Background(k.Ornament)
		ui.Text(c, "States").Font(theme.SerifFont).FontSize(theme.H3Size).TextColor(k.Text)
		ui.Column(c).Gap(16).AlignItems(ui.Start).Children(samples)
	})
}

// choiceFlag is a parameter switch that logs its change.
func choiceFlag(c *ui.Context, on *bool, label string, log func(string)) {
	if mujica.Switch(c, on, label, mujica.SwitchOptions{}).Changed() {
		log(fmt.Sprintf("%s = %v", label, *on))
	}
}

// choiceLabeled shows a caption above a sample.
func choiceLabeled(c *ui.Context, caption string, fn func()) {
	ui.Column(c).Gap(6).Children(func() {
		ui.Text(c, caption).FontSize(theme.CaptionSize).TextColor(mujica.Tokens(c).TextMuted)
		fn()
	})
}

func switchDemo(c *ui.Context, log func(string)) {
	s := *ui.Local(c.Root(), "031", func() *struct{ wifi, disabled, loading, pending bool } {
		return &struct{ wifi, disabled, loading, pending bool }{wifi: true}
	})
	choiceDemo(c, func() {
		choiceFlag(c, &s.disabled, "Disabled", log)
		choiceFlag(c, &s.loading, "Loading", log)
	}, func() {
		if mujica.Switch(c, &s.wifi, "Wi-Fi", mujica.SwitchOptions{Disabled: s.disabled, Loading: s.loading}).Changed() {
			log(fmt.Sprintf("Wi-Fi changed: %v", s.wifi))
		}
		// A caller applying the change: loading for a moment, then done.
		at := ui.Local(c.Root(), "031-at", func() time.Time { return time.Time{} })
		if mujica.Switch(c, &s.pending, "Backup (saves for 1 s)", mujica.SwitchOptions{Loading: !at.IsZero()}).Changed() {
			*at = c.Now()
			log(fmt.Sprintf("Backup requested: %v", s.pending))
		}
		if !at.IsZero() {
			if c.Now().Sub(*at) >= time.Second {
				*at = time.Time{}
				log("Backup saved")
			} else {
				c.After(time.Second - c.Now().Sub(*at))
			}
		}
		off, on := false, true
		ui.Row(c).Gap(24).Children(func() {
			mujica.Switch(c, &off, "Disabled off", mujica.SwitchOptions{Disabled: true})
			mujica.Switch(c, &on, "Disabled on", mujica.SwitchOptions{Disabled: true})
		})
	})
}

func toggleDemo(c *ui.Context, log func(string)) {
	s := *ui.Local(c.Root(), "032", func() *struct{ bold, italic, underline, disabled bool } {
		return &struct{ bold, italic, underline, disabled bool }{bold: true}
	})
	choiceDemo(c, func() {
		choiceFlag(c, &s.disabled, "Disabled", log)
	}, func() {
		ui.Row(c).Gap(8).Children(func() {
			for _, t := range []struct {
				on         *bool
				label, ico string
			}{{&s.bold, "Bold", "bold"}, {&s.italic, "Italic", "italic"}, {&s.underline, "Underline", "underline"}} {
				if mujica.Toggle(c, t.on, "", mujica.ToggleOptions{Icon: icons.Must(t.ico), Label: t.label, Disabled: s.disabled}).Changed() {
					log(fmt.Sprintf("%s: %v", t.label, *t.on))
				}
			}
		})
		if mujica.Toggle(c, &s.bold, "Bold text", mujica.ToggleOptions{Disabled: s.disabled}).Changed() {
			log(fmt.Sprintf("Bold: %v", s.bold))
		}
	})
}

func toggleGroupDemo(c *ui.Context, log func(string)) {
	type state struct {
		value                                    []string
		multiple, allowEmpty, vertical, disabled bool
	}
	s := *ui.Local(c.Root(), "033", func() *state { return &state{value: []string{"left"}} })
	items := []mujica.ToggleItem[string]{
		{Value: "left", Label: "Left"}, {Value: "center", Label: "Center"},
		{Value: "right", Label: "Right"}, {Value: "justify", Label: "Justify", Disabled: true},
	}
	choiceDemo(c, func() {
		if mujica.Switch(c, &s.multiple, "Multiple", mujica.SwitchOptions{}).Changed() {
			s.value = s.value[:min(len(s.value), 1)]
			log(fmt.Sprintf("Multiple = %v", s.multiple))
		}
		choiceFlag(c, &s.allowEmpty, "AllowEmpty", log)
		choiceFlag(c, &s.vertical, "Vertical", log)
		choiceFlag(c, &s.disabled, "Disabled", log)
	}, func() {
		opts := mujica.ToggleGroupOptions{Multiple: s.multiple, AllowEmpty: s.allowEmpty, Vertical: s.vertical, Disabled: s.disabled, Label: "Align"}
		if mujica.ToggleGroup(c, &s.value, items, opts).Changed() {
			log(fmt.Sprintf("value: %v", s.value))
		}
		ui.Text(c, fmt.Sprintf("Value: %v", s.value)).TextColor(mujica.Tokens(c).TextMuted)
	})
}

func sliderDemo(c *ui.Context, log func(string)) {
	type state struct {
		value, step                float64
		marks, showValue, disabled bool
	}
	s := *ui.Local(c.Root(), "034", func() *state { return &state{value: 40, step: 5, showValue: true} })
	choiceDemo(c, func() {
		choiceFlag(c, &s.marks, "Marks", log)
		choiceFlag(c, &s.showValue, "ShowValue", log)
		choiceFlag(c, &s.disabled, "Disabled", log)
		choiceLabeled(c, fmt.Sprintf("Step %v", s.step), func() {
			if mujica.Slider(c, &s.step, mujica.SliderOptions{Max: 25, Step: 5, Label: "Step"}).Changed() {
				log(fmt.Sprintf("Step = %v", s.step))
			}
		})
	}, func() {
		var marks []mujica.SliderMark
		if s.marks {
			for v := 0.0; v <= 100; v += 25 {
				marks = append(marks, mujica.SliderMark{Value: v, Label: fmt.Sprint(v)})
			}
		}
		ui.Column(c).Width(360).Children(func() {
			opts := mujica.SliderOptions{Step: s.step, Marks: marks, ShowValue: s.showValue, Disabled: s.disabled, Label: "Volume"}
			if mujica.Slider(c, &s.value, opts).Changed() {
				log(fmt.Sprintf("Volume: %v", s.value))
			}
		})
		ui.Text(c, fmt.Sprintf("Value: %v", s.value)).TextColor(mujica.Tokens(c).TextMuted)
	})
}

func rangeSliderDemo(c *ui.Context, log func(string)) {
	type state struct {
		low, high                  float64
		showValue, disabled, marks bool
	}
	s := *ui.Local(c.Root(), "035", func() *state { return &state{low: 20, high: 80, showValue: true} })
	choiceDemo(c, func() {
		choiceFlag(c, &s.marks, "Marks", log)
		choiceFlag(c, &s.showValue, "ShowValue", log)
		choiceFlag(c, &s.disabled, "Disabled", log)
	}, func() {
		var marks []mujica.SliderMark
		if s.marks {
			marks = []mujica.SliderMark{{Value: 0, Label: "£0"}, {Value: 500, Label: "£500"}, {Value: 1000, Label: "£1000"}}
		}
		ui.Column(c).Width(360).Children(func() {
			opts := mujica.RangeSliderOptions{Max: 1000, Step: 50, Marks: marks, ShowValue: s.showValue, Disabled: s.disabled, Label: "Price",
				Format: func(v float64) string { return fmt.Sprintf("£%.0f", v) }}
			if mujica.RangeSlider(c, &s.low, &s.high, opts).Changed() {
				log(fmt.Sprintf("Price: %v–%v", s.low, s.high))
			}
		})
		ui.Text(c, fmt.Sprintf("Range: £%.0f – £%.0f", s.low, s.high)).TextColor(mujica.Tokens(c).TextMuted)
	})
}

// choiceDelay is how long the demos' pretend server takes.
const choiceDelay = 700 * time.Millisecond

// choiceDue reports whether a request started at "at" has finished, asking for a frame when it will.
func choiceDue(c *ui.Context, at time.Time) bool {
	if left := choiceDelay - c.Now().Sub(at); left > 0 {
		c.After(left)
		return false
	}
	return true
}

var choiceFruits = []mujica.SelectOption[string]{
	{Value: "apple", Label: "Apple"}, {Value: "banana", Label: "Banana"}, {Value: "blueberry", Label: "Blueberry"},
	{Value: "cherry", Label: "Cherry", Disabled: true}, {Value: "date", Label: "Date"}, {Value: "fig", Label: "Fig"},
}

func selectDemo(c *ui.Context, log func(string)) {
	type state struct {
		value                       string
		clearable, disabled, noData bool
	}
	s := *ui.Local(c.Root(), "036", func() *state { return &state{value: "banana", clearable: true} })
	choiceDemo(c, func() {
		choiceFlag(c, &s.clearable, "Clearable", log)
		choiceFlag(c, &s.disabled, "Disabled", log)
		choiceFlag(c, &s.noData, "No options", log)
	}, func() {
		options := choiceFruits
		if s.noData {
			options = nil
		}
		choiceLabeled(c, "Fruit", func() {
			if mujica.Select(c, &s.value, options, mujica.SelectOptions{Clearable: s.clearable, Disabled: s.disabled, Label: "Fruit"}).Changed() {
				log(fmt.Sprintf("Fruit: %q", s.value))
			}
		})
		empty := ""
		choiceLabeled(c, "Placeholder", func() {
			mujica.Select(c, &empty, choiceFruits, mujica.SelectOptions{Placeholder: "Pick a fruit"})
		})
	})
}

func multiSelectDemo(c *ui.Context, log func(string)) {
	type state struct {
		values         []string
		limit, disable bool
	}
	s := *ui.Local(c.Root(), "037", func() *state { return &state{values: []string{"apple", "date"}} })
	choiceDemo(c, func() {
		choiceFlag(c, &s.limit, "Max 3", log)
		choiceFlag(c, &s.disable, "Disabled", log)
	}, func() {
		limit := 0
		if s.limit {
			limit = 3
		}
		choiceLabeled(c, "Fruits", func() {
			r := mujica.MultiSelect(c, &s.values, choiceFruits, mujica.MultiSelectOptions{Max: limit, MaxTags: 2, Disabled: s.disable, Label: "Fruits"})
			if r.Width(320).Changed() {
				log(fmt.Sprintf("Fruits: %v", s.values))
			}
		})
	})
}

var choiceFonts = []string{"Baskerville", "Bodoni", "Caslon", "Didot", "Garamond", "Georgia", "Jenson", "Palatino"}

func comboboxDemo(c *ui.Context, log func(string)) {
	type state struct {
		value, query, asked string
		custom, remote      bool
		at                  time.Time
		results             []string
	}
	s := *ui.Local(c.Root(), "038", func() *state { return &state{custom: true} })
	choiceDemo(c, func() {
		choiceFlag(c, &s.custom, "AllowCustom", log)
		choiceFlag(c, &s.remote, "Remote results", log)
	}, func() {
		options, loading := choiceFonts, false
		if s.remote {
			// The caller fetches results for the query and hands them back.
			if s.asked != s.query && choiceDue(c, s.at) {
				s.asked = s.query
				s.results = nil
				for _, f := range choiceFonts {
					if strings.Contains(strings.ToLower(f), strings.ToLower(s.query)) {
						s.results = append(s.results, f)
					}
				}
				log(fmt.Sprintf("results for %q: %d", s.query, len(s.results)))
			}
			options, loading = s.results, s.asked != s.query
		}
		choiceLabeled(c, "Typeface", func() {
			r := mujica.Combobox(c, &s.value, options, mujica.ComboboxOptions{AllowCustom: s.custom, ExternalFilter: s.remote, Loading: loading, Label: "Typeface"})
			if r.QueryChanged() {
				s.query, s.at = r.Query(), c.Now()
				log(fmt.Sprintf("query: %q", s.query))
			}
			if r.Changed() {
				log(fmt.Sprintf("Typeface: %q", s.value))
			}
		})
		disabled := "Didot"
		choiceLabeled(c, "Disabled", func() {
			mujica.Combobox(c, &disabled, choiceFonts, mujica.ComboboxOptions{Disabled: true})
		})
	})
}

// choiceAtlas is the cascader demo's pretend server data.
var choiceAtlas = map[string][]mujica.CascaderNode{
	"":      {{Value: "eu", Label: "Europe"}, {Value: "as", Label: "Asia"}, {Value: "aq", Label: "Antarctica", Disabled: true}},
	"eu":    {{Value: "fr", Label: "France"}, {Value: "it", Label: "Italy"}},
	"as":    {{Value: "jp", Label: "Japan"}, {Value: "cn", Label: "China"}},
	"eu/fr": {{Value: "paris", Label: "Paris", Leaf: true}, {Value: "lyon", Label: "Lyon", Leaf: true}},
	"eu/it": {{Value: "rome", Label: "Rome", Leaf: true}},
	"as/jp": {{Value: "tokyo", Label: "Tokyo", Leaf: true}, {Value: "kyoto", Label: "Kyoto", Leaf: true}},
	"as/cn": {{Value: "beijing", Label: "Beijing", Leaf: true}},
}

func cascaderDemo(c *ui.Context, log func(string)) {
	type state struct {
		path               []string
		levels             map[string]mujica.CascaderLevel
		pending            map[string]time.Time
		anyLevel, fail, ds bool
	}
	s := *ui.Local(c.Root(), "039", func() *state {
		return &state{levels: map[string]mujica.CascaderLevel{"": {Nodes: choiceAtlas[""]}}, pending: map[string]time.Time{}}
	})
	for key, at := range s.pending {
		if choiceDue(c, at) {
			delete(s.pending, key)
			if s.fail {
				s.levels[key] = mujica.CascaderLevel{Status: mujica.CascaderFailed, Err: errors.New("server busy")}
				log(fmt.Sprintf("load %q failed", key))
				continue
			}
			s.levels[key] = mujica.CascaderLevel{Nodes: choiceAtlas[key]}
			log(fmt.Sprintf("loaded %q", key))
		}
	}
	level := func(p []string) mujica.CascaderLevel {
		if lv, ok := s.levels[strings.Join(p, "/")]; ok {
			return lv
		}
		return mujica.CascaderLevel{Status: mujica.CascaderIdle}
	}
	load := func(p []string) {
		key := strings.Join(p, "/")
		s.levels[key] = mujica.CascaderLevel{Status: mujica.CascaderLoading}
		s.pending[key] = c.Now()
		log(fmt.Sprintf("load %q", key))
	}
	choiceDemo(c, func() {
		choiceFlag(c, &s.anyLevel, "AnyLevel", log)
		choiceFlag(c, &s.fail, "Fail loads", log)
		choiceFlag(c, &s.ds, "Disabled", log)
	}, func() {
		choiceLabeled(c, "City (children load on demand)", func() {
			opts := mujica.CascaderOptions{Level: level, Load: load, AnyLevel: s.anyLevel, Disabled: s.ds, Label: "City"}
			if mujica.Cascader(c, &s.path, opts).Width(320).Changed() {
				log(fmt.Sprintf("path: %v", s.path))
			}
		})
	})
}

var choiceRoles = []mujica.TreeSelectNode{
	{ID: "court", Label: "Court", Children: []mujica.TreeSelectNode{
		{ID: "king", Label: "King"}, {ID: "queen", Label: "Queen"}, {ID: "jester", Label: "Jester", Disabled: true},
		{ID: "guard", Label: "Royal guard", Children: []mujica.TreeSelectNode{{ID: "captain", Label: "Captain"}, {ID: "archer", Label: "Archer"}}},
	}},
	{ID: "church", Label: "Church", Children: []mujica.TreeSelectNode{{ID: "bishop", Label: "Bishop"}, {ID: "monk", Label: "Monk"}}},
	{ID: "merchant", Label: "Merchant"},
}

func treeSelectDemo(c *ui.Context, log func(string)) {
	type state struct {
		sel                             []string
		multiple, independent, disabled bool
	}
	s := *ui.Local(c.Root(), "040", func() *state { return &state{sel: []string{"queen", "bishop"}, multiple: true} })
	choiceDemo(c, func() {
		if mujica.Switch(c, &s.multiple, "Multiple", mujica.SwitchOptions{}).Changed() {
			s.sel = s.sel[:min(len(s.sel), 1)]
			log(fmt.Sprintf("Multiple = %v", s.multiple))
		}
		choiceFlag(c, &s.independent, "Independent rule", log)
		choiceFlag(c, &s.disabled, "Disabled", log)
	}, func() {
		rule := mujica.TreeSelectCascade
		if s.independent {
			rule = mujica.TreeSelectIndependent
		}
		choiceLabeled(c, "Roles", func() {
			opts := mujica.TreeSelectOptions{Multiple: s.multiple, Rule: rule, MaxTags: 2, Disabled: s.disabled, Label: "Roles"}
			if mujica.TreeSelect(c, &s.sel, choiceRoles, opts).Width(320).Changed() {
				log(fmt.Sprintf("Roles: %v", s.sel))
			}
		})
	})
}
