package account

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
)

var settingEntries = []SettingEntry{
	{ID: "sound", Page: "general", Label: "Sound", Description: "Chime on new messages.", Keywords: []string{"audio"}},
	{ID: "track", Page: "privacy", Label: "Usage reports", Description: "Send anonymous usage."},
}

func TestSettingsSearchFindsHighlightsAndReveals(t *testing.T) {
	page, query, sound, track, n := "privacy", "", false, false, 0
	var chosen SettingEntry
	tt := accountTester(1000, 600, func(c *ui.Context) {
		ui.Row(c).Fill().Children(func() {
			ui.Box(c).Width(300).Children(func() {
				if e, ok := SettingsSearch(c, &query, settingEntries, SettingsSearchOptions{}).Chosen(); ok {
					chosen, page = e, e.Page
				}
			})
			ui.Box(c).Grow(1).FillHeight().Children(func() { settingsView(c, &page, &sound, &track, &n) })
		})
	})
	_ = tt.Click("Search settings")
	tt.Type("AUD")
	if _, ok := tt.Find("Sound"); !ok || tt.HasText("Send anonymous usage.") {
		t.Fatalf("keyword search: %q", tt.Texts())
	}
	_ = tt.Click("Search settings")
	tt.Key(ui.Cmd, ui.KeyA)
	tt.Type("ou")
	if !tt.HasText("ou") || !tt.HasText("nd") || !tt.HasText("Usage reports") {
		t.Fatalf("match not split for marking: %q", tt.Texts())
	}
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	if chosen.ID != "sound" || page != "general" {
		t.Fatalf("chose %q page %q", chosen.ID, page)
	}
	tt.Frame()
	if !tt.Focused("Sound") || query != "" {
		t.Fatalf("row not revealed, query %q", query)
	}
	_ = tt.Click("Search settings")
	tt.Type("zzz")
	if !tt.HasText("No settings match") {
		t.Fatal("no empty result")
	}
}

func TestSettingsSearchSnapshot(t *testing.T) {
	snapshot(t, "settings_search", 360, 220, func(c *ui.Context) {
		q := ui.Local(c.Root(), "q", func() string { return "us" })
		ui.Column(c).Padding(16).Children(func() { SettingsSearch(c, q, settingEntries, SettingsSearchOptions{}) })
	})
}

func TestSettingsSearchHighlightsUnicodeCaseMapping(t *testing.T) {
	var spans []ui.Span
	tt := accountTester(400, 220, func(c *ui.Context) {
		q := "i"
		spans = settingsSearchSpans(c, "İstanbul", q, 14)
		SettingsSearch(c, &q, []SettingEntry{{ID: "language", Page: "general", Label: "İstanbul"}}, SettingsSearchOptions{})
	})
	if len(spans) != 2 || spans[0].Text != "İ" || spans[0].Weight != 700 || spans[1].Text != "stanbul" || !tt.HasText("İstanbul") {
		t.Fatal("Unicode case mapping loses match highlighting")
	}
}

func TestSettingsSearchPlainSpansHaveReadableColor(t *testing.T) {
	accountTester(400, 220, func(c *ui.Context) {
		for _, span := range settingsSearchSpans(c, "Usage reports", "us", 14) {
			if span.Color != core.Tokens(c).Text {
				t.Error("search text does not use the readable Text token")
			}
		}
	})
}

func TestSettingsSearchCachedPlanSeesInPlaceChanges(t *testing.T) {
	st := settingsSearchState{}
	entries := []SettingEntry{{ID: "theme", Page: "general", Label: "Appearance", Keywords: []string{"theme"}}}
	settingsSearchPlan(&st, entries, "theme", 8)
	if len(st.hits) != 1 {
		t.Fatal("initial keyword fails")
	}
	before := &st.hits[0]
	settingsSearchPlan(&st, entries, "theme", 8)
	if before != &st.hits[0] {
		t.Fatal("unchanged data rebuilt search plan")
	}
	entries[0].Keywords[0] = "color"
	settingsSearchPlan(&st, entries, "theme", 8)
	if len(st.hits) != 0 {
		t.Fatal("cached search misses in-place keyword change")
	}
}

func TestSettingsSearchDisabledAndLimit(t *testing.T) {
	q, chosen := "o", 0
	tt := accountTester(600, 500, func(c *ui.Context) {
		ui.Column(c).Disabled(true).Children(func() {
			if _, ok := SettingsSearch(c, &q, settingEntries, SettingsSearchOptions{Limit: 1}).Chosen(); ok {
				chosen++
			}
		})
	})
	if err := tt.Click("Sound"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	tt.Type("edited")
	if q != "o" || chosen != 0 || tt.HasText("Send anonymous usage.") {
		t.Fatal("disabled or limited search acts incorrectly")
	}
	accountPanics(t, func(c *ui.Context) { SettingsSearch(c, &q, settingEntries, SettingsSearchOptions{Limit: -1}) })
}
