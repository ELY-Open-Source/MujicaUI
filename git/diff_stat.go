package git

import (
	"fmt"
	"math"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
)

// DiffStatOptions configure a DiffStat.
type DiffStatOptions struct {
	// Blocks is the ratio bar length; zero means 5.
	Blocks int
}

// DiffStat shows added and deleted counts with a ratio bar.
func DiffStat(c *ui.Context, added, deleted int, opts DiffStatOptions) *ui.Element {
	if added < 0 || deleted < 0 {
		panic("mujica: DiffStat counts must not be negative")
	}
	n := opts.Blocks
	if n == 0 {
		n = 5
	}
	if n < 0 {
		panic("mujica: DiffStat Blocks is negative")
	}
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	plus, minus := diffStatBlocks(added, deleted, n)
	row := ui.Row(c).AlignItems(ui.Center).Gap(u * 2).Shrink(0).
		Label(fmt.Sprintf(core.Msg(c, "diff_stat.label"), added, deleted))
	row.Children(func() {
		mono := func(s string, col ui.Color) {
			ui.Text(c, s).Font(theme.MonoFont).FontFeatures("tnum").FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(col).SingleLine()
		}
		mono(fmt.Sprintf("+%d", added), k.Success)
		mono(fmt.Sprintf("−%d", deleted), k.Danger)
		ui.Row(c).Gap(2).Shrink(0).Role(ui.RoleNone).Children(func() {
			for i := range n {
				col := k.Border
				switch {
				case i < plus:
					col = k.Success
				case i < plus+minus:
					col = k.Danger
				}
				ui.Box(c).Size(8, 8).Radius(1).Background(col)
			}
		})
	})
	return row
}

// diffStatBlocks splits n blocks, one at least per nonzero side.
func diffStatBlocks(added, deleted, n int) (plus, minus int) {
	total := added + deleted
	if total == 0 {
		return 0, 0
	}
	plus = int(math.Round(float64(added) / float64(total) * float64(n)))
	minus = n - plus
	if added > 0 && plus == 0 {
		plus, minus = 1, n-1
	}
	if deleted > 0 && minus == 0 && n > 1 {
		plus, minus = n-1, 1
	}
	return plus, minus
}
