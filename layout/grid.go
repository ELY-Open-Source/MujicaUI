package layout

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

// GridOptions configure a Grid.
type GridOptions struct {
	// Columns size each column with ui.Fixed, ui.Fr or ui.FitContent; nil is one Fr(1) column.
	Columns []ui.Track
	// Rows size the first rows; rows added after fit their content.
	Rows []ui.Track
	// Gap defaults to SpaceM.
	Gap Space
	// RowGap defaults to Gap.
	RowGap Space
}

// Grid lays children out in columns and rows; GridCell places and spans them.
func Grid(c *ui.Context, opts GridOptions, children func()) ui.Element {
	if children == nil {
		panic("mujica: Grid needs children")
	}
	d := core.Density(c)
	gap := containerSpace(d, opts.Gap, SpaceM)
	rowGap := gap
	if opts.RowGap != SpaceDefault {
		rowGap = containerSpace(d, opts.RowGap, SpaceM)
	}
	cols := opts.Columns
	if len(cols) == 0 {
		cols = []ui.Track{ui.Fr(1)}
	}
	e := ui.Grid(c).ColumnTracks(cols...).GapX(gap).GapY(rowGap)
	if len(opts.Rows) > 0 {
		e.RowTracks(opts.Rows...)
	}
	return e.Children(children)
}

// GridCellOptions place one cell. Lines count from 1; zero places the cell automatically.
type GridCellOptions struct {
	Column, Row int
	// ColumnSpan and RowSpan default to 1; -1 spans to the last track.
	ColumnSpan, RowSpan int
}

// GridCell is one child of a Grid, placed and spanned by opts; build makes one element, stretched over the cell.
func GridCell(c *ui.Context, opts GridCellOptions, build func()) ui.Element {
	if build == nil {
		panic("mujica: GridCell needs build")
	}
	if opts.Column < 0 || opts.Row < 0 {
		panic("mujica: GridCell Column and Row count from 1")
	}
	if opts.ColumnSpan < -1 || opts.RowSpan < -1 {
		panic("mujica: GridCell spans must be positive, or -1 for the rest")
	}
	// A one-track grid stretches its child over the whole cell, as the grid does its own children.
	e := ui.Grid(c).Columns(1).GridRows(1)
	if opts.Column > 0 {
		e.ColumnStart(opts.Column)
	}
	if opts.Row > 0 {
		e.RowStart(opts.Row)
	}
	if opts.ColumnSpan != 0 {
		e.ColumnSpan(opts.ColumnSpan)
	}
	if opts.RowSpan != 0 {
		e.RowSpan(opts.RowSpan)
	}
	return e.Children(build)
}
