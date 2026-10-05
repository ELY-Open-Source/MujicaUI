package main

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	mujica "mujicaui"
)

func c2ViolinDemo(c *ui.Context, log func(string))  { c2DistributionDemo(c, log, "129", "ViolinPlot") }
func c2DensityDemo(c *ui.Context, log func(string)) { c2DistributionDemo(c, log, "138", "DensityPlot") }

func c2DistributionDemo(c *ui.Context, log func(string), id, name string) {
	type state struct {
		status    string
		bandwidth float64
		disabled  bool
		groups    []mujica.BoxGroup
	}
	s := ui.Local(c.Root(), id, func() state {
		return state{status: "Ready", groups: []mujica.BoxGroup{
			{Name: "Matins", Values: []float64{1, 2, 2, 3, 4, 5, 8, 9}},
			{Name: "Vespers", Values: []float64{3, 4, 5, 6, 6, 7, 8, 11}},
		}}
	})
	c1Section(c, name, func() {
		c2Number(c, &s.bandwidth, "Bandwidth", 0, 5, 0.5, log)
		c2Check(c, &s.disabled, "Disabled", log)
	}, func() {
		f, show := c2Frame(c, name, &s.status, log)
		groups := s.groups
		if !show {
			groups = nil
		}
		ui.Column(c).Disabled(s.disabled).Children(func() {
			if id == "129" {
				c1Log(log, name, mujica.ViolinPlot(c, groups, mujica.ViolinPlotOptions{Frame: f, Bandwidth: s.bandwidth}))
			} else {
				c1Log(log, name, mujica.DensityPlot(c, groups, mujica.DensityPlotOptions{Frame: f, Bandwidth: s.bandwidth}))
			}
		})
	})
}

func c2ParallelDemo(c *ui.Context, log func(string)) {
	type state struct {
		status   string
		disabled bool
		filters  map[string][2]float64
		rows     [][]float64
	}
	s := ui.Local(c.Root(), "130", func() state {
		return state{status: "Ready", filters: map[string][2]float64{}, rows: [][]float64{{20, 3, 9}, {35, 8, 2}, {50, 5, 5}, {65, 1, 4}}}
	})
	c1Section(c, "ParallelCoordinates", func() {
		c2Check(c, &s.disabled, "Disabled", log)
		if mujica.Button(c, "Clear filters", mujica.ButtonOptions{Disabled: s.disabled}).Clicked() {
			s.filters = map[string][2]float64{}
			log("ParallelCoordinates: clear filters")
		}
	}, func() {
		f, show := c2Frame(c, "ParallelCoordinates", &s.status, log)
		rows := s.rows
		names := []string{"Agnes", "Clare", "Hild", "Bede"}
		if !show {
			rows = nil
			names = nil
		}
		ui.Column(c).Disabled(s.disabled).Children(func() {
			c1Log(log, "ParallelCoordinates", mujica.ParallelCoordinates(c, &s.filters, []string{"Age", "Hymns", "Vigils"}, rows, mujica.ParallelCoordinatesOptions{Frame: f, Names: names}))
		})
		ui.Text(c, fmt.Sprintf("Filters: %v", s.filters))
	})
}

func c2ChordDemo(c *ui.Context, log func(string)) {
	type state struct {
		status   string
		disabled bool
		amount   float64
		matrix   [][]float64
	}
	s := ui.Local(c.Root(), "131", func() state {
		return state{status: "Ready", amount: 5, matrix: [][]float64{{0, 5, 2}, {4, 0, 1}, {1, 2, 0}}}
	})
	c1Section(c, "ChordDiagram", func() {
		old := s.amount
		c2Number(c, &s.amount, "Nave to Choir", 0, 20, 1, log)
		if old != s.amount {
			s.matrix = [][]float64{{0, s.amount, 2}, {4, 0, 1}, {1, 2, 0}}
		}
		c2Check(c, &s.disabled, "Disabled", log)
	}, func() {
		f, show := c2Frame(c, "ChordDiagram", &s.status, log)
		names, matrix := []string{"Nave", "Choir", "Crypt"}, s.matrix
		if !show {
			names = nil
			matrix = nil
		}
		ui.Column(c).Disabled(s.disabled).Children(func() {
			c1Log(log, "ChordDiagram", mujica.ChordDiagram(c, names, matrix, mujica.ChordDiagramOptions{Frame: f}))
		})
	})
}
