package layout

import (
	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
)

// TitleBarOptions configure a TitleBar.
type TitleBarOptions struct {
	// Leading holds tools before the title, such as a sidebar toggle.
	Leading func()
	// Trailing holds tools after the title.
	Trailing func()
}

// TitleBar is a window title row that drags the window and clears the system buttons.
func TitleBar(c *ui.Context, title string, opts TitleBarOptions) *ui.Element {
	k := core.Tokens(c)
	d := core.Density(c)
	bar := c.TitleBar()
	pad := d.Unit() * 3
	e := ui.Row(c).MinHeight(max(bar.Height, d.ControlHeight()+d.Unit()*2)).Shrink(0).
		Padding(0, bar.Right+pad, 0, bar.Left+pad).Gap(d.Unit()*2).
		Background(k.Surface).BorderWidth(0, 0, theme.BorderWidth, 0).BorderColor(k.Border).
		DragWindow()
	return e.Children(func() {
		if opts.Leading != nil {
			ui.Row(c).Gap(d.Unit()).Shrink(0).Children(opts.Leading)
		}
		ui.Text(c, title).Font(theme.SerifFont).FontSize(core.FontSize(c, theme.BodySize+1)).FontWeight(600).
			SingleLine().Grow(1).Shrink(1).MinWidth(0)
		if opts.Trailing != nil {
			ui.Row(c).Gap(d.Unit()).Shrink(0).Children(opts.Trailing)
		}
	})
}
