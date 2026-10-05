package devtools

import (
	"fmt"
	"log/slog"
	"slices"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/icons"
	"github.com/ZacharyZhang-NY/MujicaUI/input"
)

// DashboardFilters holds filters; zero uses caller defaults.
type DashboardFilters struct {
	// Range is the look-back; zero uses the caller's default.
	Range time.Duration
	Env   string
	Tags  []string
}

// DashboardFilterBarOptions configure a DashboardFilterBar.
type DashboardFilterBarOptions struct {
	// Ranges defaults to 15m, 1h, 24h and 7d.
	Ranges []time.Duration
	// Envs and Tags are the offered environments and tags.
	Envs, Tags []string
	Disabled   bool
}

var dashboardRanges = []time.Duration{15 * time.Minute, time.Hour, 24 * time.Hour, 7 * 24 * time.Hour}

// DashboardFilterBar edits time, environment and tags with clear.
func DashboardFilterBar(c *ui.Context, f *DashboardFilters, opts DashboardFilterBarOptions) core.ChoiceResult {
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
	for _, tag := range f.Tags {
		if !slices.Contains(opts.Tags, tag) {
			panic("mujica: DashboardFilterBar tag is not offered")
		}
	}
	u := core.Density(c).Unit()
	var r core.ChoiceResult
	r.Element = ui.Row(c).Shrink(0).Wrap().Gap(u * 3).AlignItems(ui.Center).Disabled(opts.Disabled).
		Label(core.Msg(c, "dashboard_filter_bar.label"))
	locked := r.Element.IsDisabled()
	r.Element.Children(func() {
		rs := []input.SelectOption[time.Duration]{{Value: 0, Label: core.Msg(c, "dashboard_filter_bar.default_range")}}
		for _, d := range ranges {
			rs = append(rs, input.SelectOption[time.Duration]{Value: d, Label: fmt.Sprintf(core.Msg(c, "dashboard_filter_bar.last"), dashboardRangeLabel(c, d))})
		}
		if f.Range != 0 && !slices.Contains(ranges, f.Range) {
			panic("mujica: DashboardFilterBar Range is not offered")
		}
		r = core.NewChoiceResult(r.Element, input.Select(c, &f.Range, rs, input.SelectOptions{Label: core.Msg(c, "dashboard_filter_bar.range"), Disabled: locked}).Changed())
		envs := []input.SelectOption[string]{{Value: "", Label: core.Msg(c, "dashboard_filter_bar.all_envs")}}
		for _, e := range opts.Envs {
			envs = append(envs, input.SelectOption[string]{Value: e, Label: e})
		}
		r = core.NewChoiceResult(r.Element, input.Select(c, &f.Env, envs, input.SelectOptions{Label: core.Msg(c, "dashboard_filter_bar.env"), Disabled: locked}).Changed() || r.Changed())
		if len(opts.Tags) > 0 {
			tags := make([]input.SelectOption[string], len(opts.Tags))
			for i, t := range opts.Tags {
				tags[i] = input.SelectOption[string]{Value: t, Label: t}
			}
			r = core.NewChoiceResult(r.Element, input.MultiSelect(c, &f.Tags, tags, input.MultiSelectOptions{Placeholder: core.Msg(c, "dashboard_filter_bar.tags"),
				Label: core.Msg(c, "dashboard_filter_bar.tags"), Disabled: locked}).Changed() || r.Changed())
		}
		empty := f.Range == 0 && f.Env == "" && len(f.Tags) == 0
		if input.Button(c, core.Msg(c, "dashboard_filter_bar.clear"), input.ButtonOptions{Variant: input.Ghost, Icon: icons.Must("x"), Disabled: locked || empty}).Clicked() && !locked {
			*f, r = DashboardFilters{}, core.NewChoiceResult(r.Element, true)
		}
	})
	if r.Changed() {
		slog.Debug("mujica: dashboard_filter_bar change", "range", f.Range, "env", f.Env, "tags", len(f.Tags))
	}
	return r
}

// dashboardRangeLabel names a look-back, as "15 min" or "7 d".
func dashboardRangeLabel(c *ui.Context, d time.Duration) string {
	switch {
	case d <= 0:
	case d%(24*time.Hour) == 0:
		return fmt.Sprintf(core.Msg(c, "dashboard_filter_bar.days"), int(d/(24*time.Hour)))
	case d%time.Hour == 0:
		return fmt.Sprintf(core.Msg(c, "dashboard_filter_bar.hours"), int(d/time.Hour))
	case d%time.Minute == 0:
		return fmt.Sprintf(core.Msg(c, "dashboard_filter_bar.minutes"), int(d/time.Minute))
	}
	panic("mujica: DashboardFilterBar ranges must be positive whole minutes")
}
