package account

import (
	"fmt"
	"math"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/data"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/icons"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// QuotaItem is one metered resource: Used of Limit, in Unit.
type QuotaItem struct {
	Label       string
	Used, Limit float64
	Unit        string
}

// UsageQuotaOptions configure a UsageQuota.
type UsageQuotaOptions struct {
	// ResetAt provides the reset timestamp; zero hides it.
	ResetAt time.Time
	// Warn sets the warning threshold; zero means 0.8.
	Warn float64
}

// UsageQuota displays usage bars and reset time.
func UsageQuota(c *ui.Context, items []QuotaItem, opts UsageQuotaOptions) ui.Element {
	warn := opts.Warn
	if warn == 0 {
		warn = 0.8
	}
	if warn < 0 || warn > 1 || math.IsNaN(warn) {
		panic("mujica: UsageQuota Warn must be between 0 and 1")
	}
	k := core.Tokens(c)
	d := core.Density(c)
	e := ui.Column(c).Gap(d.Unit() * 4).MinWidth(320)
	e.Children(func() {
		for i, it := range items {
			if it.Label == "" || !(it.Limit > 0) || it.Used < 0 || math.IsNaN(it.Used) || math.IsInf(it.Used, 0) || math.IsInf(it.Limit, 0) {
				panic(fmt.Sprintf("mujica: QuotaItem %q needs a Label, a positive Limit and Used of zero or more", it.Label))
			}
			share := it.Used / it.Limit
			col, icon, note := k.AccentText, "", ""
			switch {
			case share >= 1:
				col, icon, note = k.Danger, "circle-x", core.Msg(c, "usage_quota.over")
			case share >= warn:
				col, icon, note = k.Warning, "triangle-alert", core.Msg(c, "usage_quota.near")
			}
			amount := fmt.Sprintf(core.Msg(c, "usage_quota.amount"), data.StatisticFormat(it.Used, 0), data.StatisticFormat(it.Limit, 0), it.Unit)
			ui.Column(c.Key(i)).Gap(d.Unit()).Role(ui.RoleMeter).Label(it.Label).
				Description(fmt.Sprintf("%s, %s%% %s", amount, data.StatisticFormat(share*100, 0), note)).Children(func() {
				ui.Row(c).Gap(d.Unit() * 2).AlignItems(ui.Center).Children(func() {
					ui.Text(c, it.Label).FontWeight(500).Grow(1)
					if icon != "" {
						ui.Row(c).Gap(d.Unit() / 2).AlignItems(ui.Center).TextColor(col).Children(func() {
							ui.Icon(c, icons.Must(icon)).Size(12, 12)
							ui.Text(c, note).FontSize(core.FontSize(c, theme.CaptionSize)).FontWeight(600)
						})
					}
					ui.Text(c, amount).TextColor(k.TextMuted).FontSize(core.FontSize(c, theme.CaptionSize)).FontFeatures("tnum")
				})
				ui.Box(c).Height(6).FillWidth().Draw(func(p *ui.Painter, r ui.Rect) {
					p.Fill(r, k.SurfaceHover, theme.ControlRadius)
					p.Fill(ui.Rect{X: r.X, Y: r.Y, W: r.W * float32(math.Min(1, share)), H: r.H}, col, theme.ControlRadius)
				})
			})
		}
		if !opts.ResetAt.IsZero() {
			ui.Row(c).Gap(d.Unit()).AlignItems(ui.Center).TextColor(k.TextMuted).Children(func() {
				ui.Icon(c, icons.Must("clock")).Size(12, 12)
				ui.Text(c, usageQuotaReset(c, opts.ResetAt)).FontSize(core.FontSize(c, theme.CaptionSize))
			})
		}
	})
	return e
}

// usageQuotaReset formats hours or the reset date.
func usageQuotaReset(c *ui.Context, at time.Time) string {
	left := at.Sub(c.Now())
	switch {
	case left <= 0:
		return core.Msg(c, "usage_quota.reset_due")
	case left < 24*time.Hour:
		return fmt.Sprintf(core.Msg(c, "usage_quota.reset_hours"), int(math.Ceil(left.Hours())))
	}
	return fmt.Sprintf(core.Msg(c, "usage_quota.reset_on"), at.Format("2006-01-02"))
}
