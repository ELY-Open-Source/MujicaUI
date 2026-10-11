package chat

import (
	"log/slog"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/icons"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// TableBlockOptions configure a TableBlock.
type TableBlockOptions struct {
	// Align sets each column's alignment: ui.Start, ui.Center or ui.End; missing columns start.
	Align []ui.Align
}

// TableBlock shows a small table inside a message; it scrolls sideways and copies as TSV.
func TableBlock(c *ui.Context, header []string, rows [][]string, opts TableBlockOptions) CopyResult {
	n := len(header)
	if n == 0 {
		panic("mujica: TableBlock needs a header")
	}
	if len(opts.Align) > n {
		panic("mujica: TableBlock has more alignments than columns")
	}
	for _, row := range rows {
		if len(row) != n {
			panic("mujica: TableBlock row width differs from the header")
		}
	}
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	var r CopyResult
	r.Element = ui.Column(c).Radius(theme.CardRadius).Border(theme.BorderWidth, k.Border).Background(k.Surface).Clip().MinWidth(0)
	copiedAt := ui.Local(r.Element, "copied", func() time.Time { return time.Time{} })
	r.Element.Children(func() {
		ui.Row(c).Padding(0, u, 0, u*3).AlignItems(ui.Center).Background(k.SurfaceHover).Children(func() {
			ui.Text(c, core.Msg(c, "table_block.title")).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted).Grow(1)
			icon, name := "copy", core.Msg(c, "table_block.copy")
			if left := codeBlockCopiedFor - c.Now().Sub(*copiedAt); left > 0 {
				icon, name = "check", core.Msg(c, "code_block.copied")
				c.After(left)
			}
			if core.IconAction(c, icons.Must(icon), name).Clicked() && !r.Element.IsDisabled() {
				c.WriteClipboard(tableBlockTSV(header, rows))
				*copiedAt, r.copied = c.Now(), true
				c.Announce(core.Msg(c, "code_block.copied"))
				slog.Debug("mujica: table_block copy", "rows", len(rows))
			}
		})
		tracks := make([]ui.Track, n)
		for i := range tracks {
			tracks[i] = ui.FitContent()
		}
		sc := ui.ScrollHorizontal(c).AlignItems(ui.Start)
		// Rows span the block: the grid is at least last frame's scroll width.
		w := sc.Bounds().W
		internal.Remeasure(c, sc, w, sc.Bounds().H)
		sc.Children(func() {
			ui.Grid(c).ColumnTracks(tracks...).MinWidth(w).Children(func() {
				for i, h := range header {
					tableBlockCell(c, h, tableBlockAlign(opts.Align, i), true, false)
				}
				for ri, row := range rows {
					for i, cell := range row {
						tableBlockCell(c, cell, tableBlockAlign(opts.Align, i), false, ri == len(rows)-1)
					}
				}
			})
		})
	})
	return r
}

func tableBlockAlign(a []ui.Align, i int) ui.Align {
	if i < len(a) {
		return a[i]
	}
	return ui.Start
}

func tableBlockCell(c *ui.Context, s string, align ui.Align, head, last bool) {
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	cell := ui.Box(c).Padding(u*2, u*3).Role(ui.RoleCell)
	if !last {
		cell.BorderWidth(0, 0, theme.BorderWidth, 0).BorderColor(k.Border)
	}
	if head {
		cell.Role(ui.RoleColumnHeader)
	}
	cell.Children(func() {
		t := ui.Text(c, tableBlockClean.Replace(s)).NoWrap().TextAlign(align)
		if head {
			t.FontWeight(600)
		}
	})
}

// tableBlockClean turns a cell's tabs and newlines into spaces.
var tableBlockClean = strings.NewReplacer("\t", " ", "\r\n", " ", "\n", " ")

// tableBlockTSV joins cells with tabs and rows with newlines.
func tableBlockTSV(header []string, rows [][]string) string {
	clean := tableBlockClean
	var b strings.Builder
	line := func(cells []string) {
		for i, s := range cells {
			if i > 0 {
				b.WriteByte('\t')
			}
			b.WriteString(clean.Replace(s))
		}
		b.WriteByte('\n')
	}
	line(header)
	for _, row := range rows {
		line(row)
	}
	return b.String()
}
