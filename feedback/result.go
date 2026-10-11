package feedback

import (
	"cmp"
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/icons"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// ResultKind is the outcome a Result page reports. The zero value is ResultSuccess.
type ResultKind int

const (
	ResultSuccess ResultKind = iota
	ResultFailure
	ResultForbidden
)

// ResultOptions configure a Result.
type ResultOptions struct {
	Kind ResultKind
	// Title uses the kind's built-in copy when empty.
	Title       string
	Description string
	// Actions label the next steps; the first is the primary button.
	Actions []string
}

// ResultResult is what a Result built and reported this frame.
type ResultResult struct {
	Element *ui.Element
	// Action is the index of the clicked action, -1 for none.
	Action int
}

// Result shows a task outcome and reports the chosen next step; the caller navigates.
func Result(c *ui.Context, opts ResultOptions) ResultResult {
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	var tone ui.Color
	var icon *ui.SVG
	var key string
	switch opts.Kind {
	case ResultSuccess:
		tone, icon, key = k.Success, icons.Must("circle-check"), "result.success.title"
	case ResultFailure:
		tone, icon, key = k.Danger, icons.Must("circle-x"), "result.failure.title"
	case ResultForbidden:
		tone, icon, key = k.Warning, icons.Must("shield-off"), "result.forbidden.title"
	default:
		panic(fmt.Sprintf("mujica: unknown ResultKind %d", int(opts.Kind)))
	}
	title := cmp.Or(opts.Title, core.Msg(c, key))
	res := ResultResult{Action: -1}
	res.Element = ui.Column(c).FillWidth().AlignItems(ui.Center).Gap(u*3).Padding(u*10, u*6).Role(ui.RoleGroup).Label(title)
	res.Element.Children(func() {
		ui.Icon(c, icon).Size(48, 48).TextColor(tone)
		ui.Text(c, title).Font(theme.SerifFont).FontSize(core.FontSize(c, theme.H2Size)).TextColor(k.Text).TextAlign(ui.Center)
		if opts.Description != "" {
			ui.Text(c, opts.Description).TextColor(k.TextMuted).TextAlign(ui.Center).MaxWidth(420).LineHeight(theme.LineHeight)
		}
		resultOrnament(c, k)
		if len(opts.Actions) == 0 {
			return
		}
		ui.Row(c).Gap(u * 3).Wrap().Justify(ui.Center).Children(func() {
			for i, a := range opts.Actions {
				v := input.Secondary
				if i == 0 {
					v = input.Primary
				}
				if input.Button(c, a, input.ButtonOptions{Variant: v}).Key(i).Clicked() {
					slog.Debug("mujica: result action", "kind", opts.Kind, "index", i)
					res.Action = i
				}
			}
		})
	})
	return res
}

// resultOrnament draws two hairlines meeting a diamond, in ornament gold.
func resultOrnament(c *ui.Context, k theme.Tokens) {
	ui.Box(c).Size(160, 12).Role(ui.RoleNone).Draw(func(p *ui.Painter, r ui.Rect) {
		internal.RuleDiamond(p, r, 9, 8, k.Ornament)
	})
}
