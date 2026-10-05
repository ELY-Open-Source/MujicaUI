package mujica

import (
	"fmt"
	"log/slog"
	"slices"
	"strings"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/internal"
	"mujicaui/theme"
)

// EmojiPickerOptions configure an EmojiPicker.
type EmojiPickerOptions struct {
	// Recent stores up to 24 emoji, newest first.
	Recent *[]string
	// Tone stores the selected skin tone, from 0 through 5.
	Tone     *int
	Disabled bool
}

// EmojiResult exposes the element and this frame’s chosen emoji.
type EmojiResult struct {
	Element *ui.Element
	// Picked contains the chosen emoji, or an empty string.
	Picked string
}

// emojiRecentMax is how many recent emoji are kept.
const emojiRecentMax = 24

type emojiState struct {
	tab    int
	query  string
	recent bool
}

type emojiKey struct{}

// EmojiPicker supports categories, search, history and skin tones.
func EmojiPicker(c *ui.Context, opts EmojiPickerOptions) EmojiResult {
	if opts.Tone != nil && (*opts.Tone < 0 || *opts.Tone >= len(emojiTones)) {
		panic("mujica: EmojiPicker Tone outside 0..5")
	}
	k := Tokens(c)
	u := Density(c).Unit()
	tone := 0
	if opts.Tone != nil {
		tone = *opts.Tone
	}
	var r EmojiResult
	col := ui.Column(c).Width(8*36+u*4).Gap(u*2).Padding(u*2).Radius(theme.PopupRadius).Background(k.Surface).
		Border(theme.BorderWidth, k.Border).Role(ui.RoleGroup).Label(msg(c, "emoji_picker.label")).Disabled(opts.Disabled)
	r.Element = col
	recent := opts.Recent != nil && len(*opts.Recent) > 0
	st := ui.Local(col, emojiKey{}, func() emojiState { return emojiState{recent: recent} })
	if recent != st.recent {
		if recent {
			st.tab++
		} else {
			st.tab = max(0, st.tab-1)
		}
		st.recent = recent
	}
	col.Children(func() {
		TextInput(c, &st.query, TextInputOptions{Placeholder: msg(c, "emoji_picker.search"), Label: msg(c, "emoji_picker.search"),
			Leading: func() { ui.Icon(c, icons.Must("search")).Size(14, 14).TextColor(k.TextMuted) }, Clearable: true})
		if opts.Tone != nil {
			emojiToneRow(c, k, opts.Tone)
		}
		if strings.TrimSpace(st.query) == "" {
			emojiTabs(c, k, &st.tab, recent)
		}
		ui.Scroll(c).Height(36*5 + u*2).Children(func() {
			grid := ui.Row(c).Wrap().FocusGroup(ui.Horizontal | ui.Vertical).Role(ui.RoleList)
			grid.Children(func() {
				for _, e := range emojiShown(c, st, opts.Recent, recent) {
					if s := emojiCell(c, k, e, tone); s != "" && !col.IsDisabled() {
						r.Picked = s
					}
				}
			})
		})
	})
	if r.Picked != "" {
		if opts.Recent != nil {
			rec := slices.DeleteFunc(slices.Clone(*opts.Recent), func(s string) bool { return s == r.Picked })
			*opts.Recent = slices.Insert(rec, 0, r.Picked)[:min(emojiRecentMax, len(rec)+1)]
		}
		slog.Debug("mujica: emoji_picker pick", "emoji", r.Picked)
	}
	return r
}

func emojiShown(c *ui.Context, st *emojiState, recentList *[]string, recent bool) []emoji {
	if q := strings.ToLower(strings.TrimSpace(st.query)); q != "" {
		var out []emoji
		for _, cat := range emojiCategories {
			for _, e := range cat.emojis {
				if strings.Contains(e.name, q) || strings.Contains(emojiLabel(c, e), q) {
					out = append(out, e)
				}
			}
		}
		if len(out) == 0 {
			ui.Text(c, msg(c, "emoji_picker.none")).TextColor(Tokens(c).TextMuted).Padding(Density(c).Unit() * 2)
		}
		return out
	}
	tab := st.tab
	if recent {
		if tab == 0 {
			out := make([]emoji, len(*recentList))
			for i, s := range *recentList {
				out[i] = emoji{char: s, name: s}
			}
			return out
		}
		tab--
	}
	return emojiCategories[min(tab, len(emojiCategories)-1)].emojis
}

// emojiTabs picks a category; Recent leads when there is history.
func emojiTabs(c *ui.Context, k theme.Tokens, tab *int, recent bool) {
	type tabItem struct{ face, key string }
	var tabs []tabItem
	if recent {
		tabs = append(tabs, tabItem{"🕘", "emoji_picker.recent"})
	}
	for _, cat := range emojiCategories {
		tabs = append(tabs, tabItem{cat.icon, cat.key})
	}
	*tab = min(*tab, len(tabs)-1)
	ui.Row(c).FocusGroup(ui.Horizontal).Role(ui.RoleTabList).BorderWidth(0, 0, 1, 0).BorderColor(k.Border).Children(func() {
		for i, t := range tabs {
			b := ui.ButtonBase(c).Key(t.key).Size(32, 32).Radius(theme.ControlRadius).Role(ui.RoleTab).Label(msg(c, t.key)).
				Selected(i == *tab).Background(ui.Transparent).Tooltip(msg(c, t.key))
			internal.Ring(b, k.Focus, theme.ControlRadius)
			switch {
			case i == *tab:
				b.Background(k.Selection).BorderWidth(0, 0, 2, 0).BorderColor(k.AccentText)
			case b.Hovered():
				b.Background(k.SurfaceHover)
			}
			b.Children(func() { ui.Text(c, t.face).FontSize(FontSize(c, 16)) })
			if b.Clicked() && !b.IsDisabled() {
				*tab = i
			}
		}
	})
}

func emojiToneRow(c *ui.Context, k theme.Tokens, tone *int) {
	ui.Row(c).Gap(2).AlignItems(ui.Center).FocusGroup(ui.Horizontal).Role(ui.RoleRadioGroup).Label(msg(c, "emoji_picker.tone")).Children(func() {
		ui.Text(c, msg(c, "emoji_picker.tone")).FontSize(FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted).Margin(0, 6, 0, 0)
		for i := range emojiTones {
			b := ui.ButtonBase(c).Key(i).Size(theme.MinHitSize, theme.MinHitSize).Radius(theme.ControlRadius).Role(ui.RoleRadio).
				Label(fmt.Sprintf(msg(c, "emoji_picker.toneN"), i)).Selected(i == *tone).Background(ui.Transparent)
			internal.Ring(b, k.Focus, theme.ControlRadius)
			if i == *tone {
				b.Background(k.Selection).Border(theme.BorderWidth, k.AccentText)
			} else if b.Hovered() {
				b.Background(k.SurfaceHover)
			}
			b.Children(func() { ui.Text(c, emojiToned(emoji{char: "👋", toned: true}, i)).FontSize(FontSize(c, 15)) })
			if b.Clicked() && !b.IsDisabled() {
				*tone = i
				slog.Debug("mujica: emoji_picker tone", "tone", i)
			}
		}
	})
}

func emojiCell(c *ui.Context, k theme.Tokens, e emoji, tone int) string {
	s := emojiToned(e, tone)
	b := ui.ButtonBase(c).Key(e.char).Size(36, 36).Radius(theme.ControlRadius).Label(emojiLabel(c, e)).Tooltip(emojiLabel(c, e)).Transition(colorFade)
	internal.Ring(b, k.Focus, theme.ControlRadius)
	if b.Hovered() || b.Focused() {
		b.Background(k.SurfaceHover)
	}
	b.Children(func() { ui.Text(c, s).FontSize(FontSize(c, 20)) })
	if b.Clicked() {
		return s
	}
	return ""
}

func emojiLabel(c *ui.Context, e emoji) string {
	if e.name == e.char {
		return e.char
	}
	return msg(c, "emoji_picker.name_"+strings.ReplaceAll(e.name, " ", "_"))
}
