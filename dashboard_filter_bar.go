package mujica

import (
	"fmt"
	"log/slog"
	"slices"
	"time"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
)

// DashboardFilters is what a dashboard is filtered by; the zero value filters nothing.
type DashboardFilters struct {
	// Range is how far back to look; zero means the caller's default.
	Range time.Duration
	Env   string
	Tags  []string
}

// DashboardFilterBarOptions configure a DashboardFilterBar.
type DashboardFilterBarOptions struct {
	// Ranges are the offered look-backs; nil uses 15 min, 1 h, 24 h and 7 d.
	Ranges []time.Duration
	// Envs and Tags are the offered environments and tags.
	Envs, Tags []string
	Disabled   bool
}

var dashboardRanges = []time.Duration{15 * time.Minute, time.Hour, 24 * time.Hour, 7 * 24 * time.Hour}

// DashboardFilterBar filters by time range, environment and tags, with a clear button.
func DashboardFilterBar(c *ui.Context, f *DashboardFilters, opts DashboardFilterBarOptions) ChoiceResult {
	if f == nil {
		panic("mujica: DashboardFilterBar needs filters")
	}
	ranges := opts.Ranges
	if ranges == nil {
		ranges = dashboardRanges
	}
	if f.Env != "" && !slices.Contains(opts.Envs, f.Env) {
		panic("mujica: DashboardFilterBar Env is not offered")
	}
	u := Density(c).Unit()
	var r ChoiceResult
	r.Element = ui.Row(c).Shrink(0).Wrap().Gap(u * 3).AlignItems(ui.Center).Disabled(opts.Disabled).
		Label(msg(c, "dashboard_filter_bar.label"))
	locked := r.Element.IsDisabled()
	r.Element.Children(func() {
		rs := []SelectOption[time.Duration]{{Value: 0, Label: msg(c, "dashboard_filter_bar.default_range")}}
		for _, d := range ranges {
			rs = append(rs, SelectOption[time.Duration]{Value: d, Label: fmt.Sprintf(msg(c, "dashboard_filter_bar.last"), dashboardRangeLabel(c, d))})
		}
		if f.Range != 0 && !slices.Contains(ranges, f.Range) {
			panic("mujica: DashboardFilterBar Range is not offered")
		}
		r.changed = Select(c, &f.Range, rs, SelectOptions{Label: msg(c, "dashboard_filter_bar.range"), Disabled: locked}).Changed()
		envs := []SelectOption[string]{{Value: "", Label: msg(c, "dashboard_filter_bar.all_envs")}}
		for _, e := range opts.Envs {
			envs = append(envs, SelectOption[string]{Value: e, Label: e})
		}
		r.changed = Select(c, &f.Env, envs, SelectOptions{Label: msg(c, "dashboard_filter_bar.env"), Disabled: locked}).Changed() || r.changed
		if len(opts.Tags) > 0 {
			tags := make([]SelectOption[string], len(opts.Tags))
			for i, t := range opts.Tags {
				tags[i] = SelectOption[string]{Value: t, Label: t}
			}
			r.changed = MultiSelect(c, &f.Tags, tags, MultiSelectOptions{Placeholder: msg(c, "dashboard_filter_bar.tags"),
				Label: msg(c, "dashboard_filter_bar.tags"), Disabled: locked}).Changed() || r.changed
		}
		empty := f.Range == 0 && f.Env == "" && len(f.Tags) == 0
		if Button(c, msg(c, "dashboard_filter_bar.clear"), ButtonOptions{Variant: Ghost, Icon: icons.Must("x"), Disabled: locked || empty}).Clicked() && !locked {
			*f, r.changed = DashboardFilters{}, true
		}
	})
	if r.changed {
		slog.Debug("mujica: dashboard_filter_bar change", "range", f.Range, "env", f.Env, "tags", len(f.Tags))
	}
	return r
}

// dashboardRangeLabel names a look-back, as "15 min" or "7 d".
func dashboardRangeLabel(c *ui.Context, d time.Duration) string {
	switch {
	case d <= 0:
	case d%(24*time.Hour) == 0:
		return fmt.Sprintf(msg(c, "dashboard_filter_bar.days"), int(d/(24*time.Hour)))
	case d%time.Hour == 0:
		return fmt.Sprintf(msg(c, "dashboard_filter_bar.hours"), int(d/time.Hour))
	case d%time.Minute == 0:
		return fmt.Sprintf(msg(c, "dashboard_filter_bar.minutes"), int(d/time.Minute))
	}
	panic("mujica: DashboardFilterBar ranges must be positive whole minutes")
}
