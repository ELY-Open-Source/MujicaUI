package git

import (
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/internal"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
)

// BlameLine is a source line and its last commit.
type BlameLine struct {
	Hash    string
	Author  string
	When    time.Time
	Summary string
	Text    string
}

// BlameViewState keeps scroll place and highlight state.
type BlameViewState struct {
	list   ui.ListState
	blocks []bool
	src    *BlameLine
	n      int
	lang   string
}

// BlameViewOptions configure a BlameView.
type BlameViewOptions struct {
	// Language picks syntax colors, as "go"; empty is plain.
	Language string
	// TimeFormat lays out commit dates; empty means "2006-01-02".
	TimeFormat string
}

// BlameView shows code with its commits; give it a height.
func BlameView(c *ui.Context, s *BlameViewState, lines []BlameLine, opts BlameViewOptions) *ui.Element {
	var first *BlameLine
	if len(lines) > 0 {
		first = &lines[0]
	}
	if first != s.src || len(lines) != s.n || opts.Language != s.lang {
		s.blocks = blameViewBlocks(lines, opts.Language)
		s.src, s.n, s.lang = first, len(lines), opts.Language
		slog.Debug("mujica: blame view index", "lines", len(lines))
	}
	layout := opts.TimeFormat
	if layout == "" {
		layout = "2006-01-02"
	}
	k := core.Tokens(c)
	d := core.Density(c)
	u := d.Unit()
	digits := float32(len(strconv.Itoa(len(lines))))
	s.list.Label = func(i int) string { return lines[i].Text }
	e := ui.List(c, &s.list, len(lines), func(i int) {
		ln := lines[i]
		if ln.Hash == "" {
			panic("mujica: BlameLine needs a Hash")
		}
		head := i == 0 || lines[i-1].Hash != ln.Hash
		row := ui.Row(c).AlignItems(ui.Stretch).MinHeight(theme.MinHitSize)
		if head && i > 0 {
			row.Draw(func(p *ui.Painter, r ui.Rect) { p.Fill(ui.Rect{X: r.X, Y: r.Y, W: r.W, H: 1}, k.Border, 0) })
		}
		row.Children(func() {
			gut := ui.Row(c).Width(240).Shrink(0).AlignItems(ui.Center).Gap(u*2).Padding(0, u*3).Background(k.Surface).
				Label(fmt.Sprintf(core.Msg(c, "blame_view.line"), i+1, ln.Author, gitShortHash(ln.Hash)))
			gut.Tooltip(fmt.Sprintf("%s\n%s · %s\n%s", ln.Hash, ln.Author, ln.When.Format(layout), ln.Summary))
			if gut.Hovered() {
				gut.Background(k.SurfaceHover)
			}
			gut.Transition(core.ColorFade)
			gut.Children(func() {
				if !head {
					return
				}
				small := core.FontSize(c, theme.CaptionSize)
				ui.Text(c, gitShortHash(ln.Hash)).Font(theme.MonoFont).FontSize(small).TextColor(k.AccentText).Shrink(0)
				ui.Text(c, ln.Author).FontSize(small).SingleLine().Grow(1).MinWidth(0)
				ui.Text(c, ln.When.Format(layout)).FontSize(small).FontFeatures("tnum").TextColor(k.TextMuted).Shrink(0)
			})
			ui.Text(c, strconv.Itoa(i+1)).Font(theme.MonoFont).FontSize(core.FontSize(c, theme.CaptionSize)).FontFeatures("tnum").
				TextColor(k.TextMuted).Width(digits*8+u*4).TextAlign(ui.End).Padding(u, u*2, 0, 0).Shrink(0)
			spans, _ := internal.CodeSpans(k, opts.Language, ln.Text, s.blocks[i])
			ui.RichText(c, spans...).Font(theme.MonoFont).NoWrap().Padding(u, u*2).AlignSelf(ui.Center)
		})
	})
	e.Background(k.Background).Border(theme.BorderWidth, k.Border).Radius(theme.ControlRadius)
	return e
}

// blameViewBlocks records where block comments are open.
func blameViewBlocks(lines []BlameLine, language string) []bool {
	blocks := make([]bool, len(lines))
	open := false
	for i, ln := range lines {
		blocks[i] = open
		_, open = internal.Highlight(language, ln.Text, open)
	}
	return blocks
}
