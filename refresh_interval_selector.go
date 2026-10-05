package mujica

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
)

// RefreshIntervalSelectorOptions configure a RefreshIntervalSelector.
type RefreshIntervalSelectorOptions struct {
	// Intervals to offer besides Off; nil uses 5 s, 10 s, 30 s, 1 m and 5 m.
	Intervals []time.Duration
	// Refreshing shows the refresh in progress and holds the timer.
	Refreshing bool
	Disabled   bool
}

// RefreshResult reports this frame's events.
type RefreshResult struct {
	Element *ui.Element
	// Refresh asks the caller to reload, by hand or because the interval elapsed.
	Refresh bool
	// Changed reports a new interval.
	Changed bool
}

type refreshState struct{ last time.Time }

var refreshDefaults = []time.Duration{5 * time.Second, 10 * time.Second, 30 * time.Second, time.Minute, 5 * time.Minute}

// RefreshIntervalSelector picks an auto-refresh interval, zero for off, and refreshes by hand.
func RefreshIntervalSelector(c *ui.Context, interval *time.Duration, opts RefreshIntervalSelectorOptions) RefreshResult {
	if interval == nil {
		panic("mujica: RefreshIntervalSelector needs an interval")
	}
	choices := opts.Intervals
	if choices == nil {
		choices = refreshDefaults
	}
	options := []SelectOption[time.Duration]{{Value: 0, Label: msg(c, "refresh_interval_selector.off")}}
	known := *interval == 0
	for _, d := range choices {
		if d < time.Second || d%time.Second != 0 {
			panic("mujica: RefreshIntervalSelector intervals must be whole seconds")
		}
		known = known || d == *interval
		label := fmt.Sprintf(msg(c, "refresh_interval_selector.seconds"), int(d/time.Second))
		if d%time.Minute == 0 {
			label = fmt.Sprintf(msg(c, "refresh_interval_selector.minutes"), int(d/time.Minute))
		}
		options = append(options, SelectOption[time.Duration]{Value: d, Label: label})
	}
	if !known {
		panic("mujica: RefreshIntervalSelector interval is not offered")
	}
	u := Density(c).Unit()
	var r RefreshResult
	r.Element = ui.Row(c).Shrink(0).Gap(u * 2).AlignItems(ui.Center).Disabled(opts.Disabled)
	st := ui.Local(r.Element, "mujica.refresh", func() refreshState { return refreshState{last: c.Now()} })
	locked := r.Element.IsDisabled()
	r.Element.Children(func() {
		label := msg(c, "refresh_interval_selector.refresh")
		if opts.Refreshing {
			Spinner(c, SpinnerOptions{Size: 14, Label: msg(c, "refresh_interval_selector.refreshing")})
		}
		if Button(c, label, ButtonOptions{Variant: Secondary, Icon: icons.Must("refresh-cw"), Disabled: locked || opts.Refreshing}).Clicked() && !locked {
			r.Refresh = true
		}
		r.Changed = Select(c, interval, options, SelectOptions{Label: msg(c, "refresh_interval_selector.label"), Disabled: locked}).Changed()
	})
	if r.Changed || opts.Refreshing {
		st.last = c.Now()
	}
	if *interval > 0 && !locked && !opts.Refreshing && !r.Refresh {
		if wait := *interval - c.Now().Sub(st.last); wait > 0 {
			c.After(wait)
		} else {
			r.Refresh = true
		}
	}
	if r.Refresh {
		st.last = c.Now()
		slog.Debug("mujica: refresh_interval_selector refresh", "interval", *interval)
	}
	return r
}
