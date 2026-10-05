package mujica

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
	"mujicaui/theme"
)

// NewsItem is a headline with its source, time and related symbols.
type NewsItem struct {
	ID, Title, Source string
	Time              time.Time
	Symbols           []string
}

// NewsFeedState holds which items were read, by ID, and the list's place.
type NewsFeedState struct {
	Read   map[string]bool
	rows   ui.ListState
	cursor int
	count  struct {
		key    string
		unread int
	}
}

// newsFeedUnread counts unread items, recounting only when items or reads change.
func newsFeedUnread(s *NewsFeedState, items []NewsItem) int {
	key := fmt.Sprint(chartID(items), len(s.Read))
	if key != s.count.key {
		s.count.key, s.count.unread = key, 0
		for _, it := range items {
			if !s.Read[it.ID] {
				s.count.unread++
			}
		}
	}
	return s.count.unread
}

// NewsFeedOptions configure a NewsFeed.
type NewsFeedOptions struct {
	// Today decides the time stamps; zero uses the window clock.
	Today    Date
	Status   DataStatus
	Error    string
	Disabled bool
}

// NewsFeedView is a built NewsFeed and this frame's events.
type NewsFeedView struct {
	Element *ui.Element
	// Opened is the item opened, which is now read; empty for none.
	Opened string
	// Symbol is the related code clicked; empty for none.
	Symbol string
	Retry  bool
}

// NewsFeed lists headlines newest first as given, marking unread ones; only rows in view are built.
func NewsFeed(c *ui.Context, s *NewsFeedState, items []NewsItem, opts NewsFeedOptions) NewsFeedView {
	if s == nil {
		panic("mujica: NewsFeed needs a state")
	}
	if s.Read == nil {
		s.Read = map[string]bool{}
	}
	today := opts.Today
	if today.IsZero() {
		today = DateOf(c.Now())
	}
	k := Tokens(c)
	d := Density(c)
	var v NewsFeedView
	v.Element = ui.Column(c).Gap(d.Unit() * 2).MinHeight(0).Disabled(opts.Disabled)
	v.Element.Children(func() {
		unread := newsFeedUnread(s, items)
		ui.Row(c).Gap(d.Unit() * 2).AlignItems(ui.Center).Children(func() {
			ui.Text(c, fmt.Sprintf(msg(c, "news_feed.unread"), unread)).TextColor(k.TextMuted).Grow(1).Role(ui.RoleStatus)
			if Button(c, msg(c, "news_feed.mark_all"), ButtonOptions{Variant: Ghost, Disabled: opts.Disabled || unread == 0}).Clicked() {
				for _, it := range items {
					s.Read[it.ID] = true
				}
				slog.Debug("mujica: news_feed mark_all", "count", unread)
			}
		})
		if len(items) == 0 {
			v.Retry = dataPlaceholder(c, opts.Status, opts.Error, "", msg(c, "news_feed.empty"))
			return
		}
		cursor := &s.cursor
		if opts.Disabled {
			cp := s.cursor
			cursor = &cp
		}
		s.rows.Key = func(i int) any { return items[i].ID }
		s.rows.Label = func(i int) string { return items[i].Title }
		s.rows.Selected = cursor
		restore := dataSelectTheme(c, 0)
		list := ui.List(c, &s.rows, len(items), func(i int) {
			it := items[i]
			read := s.Read[it.ID]
			row := ui.Row(c).Padding(d.Unit()*2, d.Unit()*2).Gap(d.Unit() * 2).AlignItems(ui.Start).MinHeight(d.ControlHeight())
			if row.Clicked() && !row.IsDisabled() {
				*cursor, v.Opened = i, it.ID
			}
			row.Children(func() {
				ui.Box(c).Size(8, 18).Shrink(0).Draw(func(p *ui.Painter, r ui.Rect) {
					if !read {
						p.FillPath(internal.Diamond(r.X+4, r.Y+9, 8), k.AccentText)
					}
				})
				ui.Column(c).Grow(1).MinWidth(0).Gap(2).Children(func() {
					t := ui.Text(c, it.Title).MaxLines(2)
					if !read {
						t.FontWeight(700)
					} else {
						t.TextColor(k.TextMuted)
					}
					ui.Row(c).Gap(d.Unit() * 2).AlignItems(ui.Center).Wrap().Children(func() {
						ui.Text(c, it.Source+" · "+dateSeparatorWhen(c, it.Time, today)).FontSize(FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted)
						for _, sym := range it.Symbols {
							ui.Row(c).Key(sym).Children(func() {
								if SymbolBadge(c, MarketSymbol{Code: sym}, SymbolBadgeOptions{Clickable: true, Disabled: opts.Disabled}).Changed() {
									v.Symbol = sym
								}
							})
						}
					})
				})
			})
		}).Grow(1).MinHeight(0).Padding(d.Unit()).Background(k.Surface).Border(theme.BorderWidth, k.Border).Radius(theme.ControlRadius)
		restore()
		internal.Ring(list, k.Focus, theme.ControlRadius)
		if !opts.Disabled && list.Submitted() && s.cursor >= 0 && s.cursor < len(items) {
			v.Opened = items[s.cursor].ID
		}
		if v.Opened != "" && !opts.Disabled {
			s.Read[v.Opened] = true
		}
	})
	if v.Element.IsDisabled() {
		v.Opened, v.Symbol, v.Retry = "", "", false
	}
	if v.Opened != "" || v.Symbol != "" {
		slog.Debug("mujica: news_feed open", "item", v.Opened, "symbol", v.Symbol)
	}
	return v
}
