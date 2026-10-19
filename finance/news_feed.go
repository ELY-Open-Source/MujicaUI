package finance

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/chart"
	"github.com/ELY-Open-Source/MujicaUI/chat"
	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/data"
	"github.com/ELY-Open-Source/MujicaUI/input"
	"github.com/ELY-Open-Source/MujicaUI/internal"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

type NewsItem struct {
	ID, Title, Source string
	Time              time.Time
	Symbols           []string
}

type NewsFeedState struct {
	Read        map[string]bool
	rows        ui.ListState
	cursor      int
	readVersion int
	count       struct {
		key    string
		unread int
	}
}

func newsFeedUnread(s *NewsFeedState, items []NewsItem, version int) int {
	key := fmt.Sprint(chart.ChartID(items), version, s.readVersion)
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
	Today    core.Date
	Version  int
	Status   data.DataStatus
	Error    string
	Disabled bool
}

// NewsFeedView is a built NewsFeed and this frame's events.
type NewsFeedView struct {
	Element ui.Element
	Opened  string
	// Symbol is the related code clicked; empty for none.
	Symbol      string
	Retry       bool
	ReadChanged bool
}

func NewsFeed(c *ui.Context, s *NewsFeedState, items []NewsItem, opts NewsFeedOptions) NewsFeedView {
	if s == nil {
		panic("mujica: NewsFeed needs a state")
	}

	today := opts.Today
	if today.IsZero() {
		today = core.DateOf(c.Now())
	}
	k := core.Tokens(c)
	d := core.Density(c)
	var v NewsFeedView
	v.Element = ui.Column(c).Gap(d.Unit() * 2).Grow(1).MinHeight(0).Disabled(opts.Disabled)
	opts.Disabled = v.Element.IsDisabled()
	if opts.Disabled {
		copy := *s
		s = &copy
	}
	if s.Read == nil {
		s.Read = map[string]bool{}
	}
	v.Element.Children(func() {
		unread := newsFeedUnread(s, items, opts.Version)
		ui.Row(c).Gap(d.Unit() * 2).AlignItems(ui.Center).Children(func() {
			ui.Text(c, fmt.Sprintf(core.Msg(c, "news_feed.unread"), unread)).TextColor(k.TextMuted).Grow(1).Role(ui.RoleStatus)
			if input.Button(c, core.Msg(c, "news_feed.mark_all"), input.ButtonOptions{Variant: input.Ghost, Disabled: opts.Disabled || unread == 0}).Clicked() {
				for _, it := range items {
					s.Read[it.ID] = true
				}
				s.readVersion++
				v.ReadChanged = true
				slog.Debug("mujica: news_feed mark_all", "count", unread)
			}
		})
		if len(items) == 0 {
			v.Retry = data.DataPlaceholder(c, opts.Status, opts.Error, "", core.Msg(c, "news_feed.empty"))
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
		restore := data.DataSelectTheme(c, 0)
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
						ui.Text(c, it.Source+" · "+chat.DateSeparatorWhen(c, it.Time, today)).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted)
						for _, sym := range it.Symbols {
							ui.Row(c.Key(sym)).Children(func() {
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
		if !list.IsDisabled() && list.Submitted() && s.cursor >= 0 && s.cursor < len(items) {
			v.Opened = items[s.cursor].ID
		}
		if v.Opened != "" && !opts.Disabled {
			if !s.Read[v.Opened] {
				s.readVersion++
				v.ReadChanged = true
			}
			s.Read[v.Opened] = true
		}
	})
	if v.Element.IsDisabled() {
		v.Opened, v.Symbol, v.Retry, v.ReadChanged = "", "", false, false
	}
	if v.Opened != "" || v.Symbol != "" {
		slog.Debug("mujica: news_feed open", "item", v.Opened, "symbol", v.Symbol)
	}
	return v
}
