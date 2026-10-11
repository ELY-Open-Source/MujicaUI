package main

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/chart"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
)

var c1Abbey = chart.ChartNode{Name: "Abbey", Children: []chart.ChartNode{
	{Name: "Library", Children: []chart.ChartNode{{Name: "Psalters", Value: 30}, {Name: "Maps", Value: 10}, {Name: "Letters", Value: 8}}},
	{Name: "Kitchen", Value: 25},
	{Name: "Garden", Children: []chart.ChartNode{{Name: "Herbs", Value: 12}, {Name: "Vines", Value: 8}}},
	{Name: "Chapel", Children: []chart.ChartNode{{Name: "Nave", Value: 18}, {Name: "Crypt", Value: 6}}},
}}

func c1TreemapDemo(c *ui.Context, log func(string)) {
	type state struct {
		status string
		path   []string
	}
	s := ui.Local(c.Root(), "116", func() state { return state{status: "Ready"} })
	c1Section(c, "Treemap", func() {}, func() {
		frame, show := c1Frame(c, "Rooms", &s.status)
		root := c1Abbey
		if !show {
			root = chart.ChartNode{Name: "Abbey"}
			s.path = nil
		}
		r := chart.Treemap(c, &s.path, root, chart.TreemapOptions{Frame: frame})
		c1Log(log, "Treemap", r)
		if r.ValueChanged() {
			log("Treemap: path /" + strings.Join(s.path, "/"))
		}
	})
}

func c1SunburstDemo(c *ui.Context, log func(string)) {
	type state struct {
		status string
		path   []string
		depth  float64
	}
	s := ui.Local(c.Root(), "117", func() state { return state{status: "Ready", depth: 2} })
	c1Section(c, "Sunburst", func() {
		ui.Text(c, fmt.Sprintf("%.0f rings", s.depth))
		ui.StepSlider(c, &s.depth, 1, 3, 1).Width(120)
	}, func() {
		frame, show := c1Frame(c, "Rooms", &s.status)
		frame.Height = 300
		root := c1Abbey
		if !show {
			root = chart.ChartNode{Name: "Abbey"}
			s.path = nil
		}
		r := chart.Sunburst(c, &s.path, root, chart.SunburstOptions{Frame: frame, Depth: int(s.depth)})
		c1Log(log, "Sunburst", r)
		if r.ValueChanged() {
			log("Sunburst: path /" + strings.Join(s.path, "/"))
		}
	})
}

func c1SankeyDemo(c *ui.Context, log func(string)) {
	type state struct{ status string }
	s := ui.Local(c.Root(), "118", func() state { return state{status: "Ready"} })
	c1Section(c, "SankeyChart", func() {}, func() {
		frame, show := c1Frame(c, "Money", &s.status)
		nodes := []chart.SankeyNode{{ID: "tithe", Label: "Tithes"}, {ID: "alms", Label: "Alms"}, {ID: "fund", Label: "Fund"}, {ID: "roof", Label: "Roof"}, {ID: "bread", Label: "Bread"}, {ID: "bells", Label: "Bells"}}
		links := []chart.SankeyLink{{From: "tithe", To: "fund", Value: 60}, {From: "alms", To: "fund", Value: 30}, {From: "alms", To: "bread", Value: 10}, {From: "fund", To: "roof", Value: 50}, {From: "fund", To: "bells", Value: 40}}
		if !show {
			links = nil
		}
		c1Log(log, "SankeyChart", chart.SankeyChart(c, nodes, links, chart.SankeyChartOptions{Frame: frame}))
	})
}

func c1NetworkDemo(c *ui.Context, log func(string)) {
	type state struct {
		status string
		force  bool
		graph  chart.NetworkGraphState
	}
	s := ui.Local(c.Root(), "119", func() state { return state{status: "Ready", force: true} })
	c1Section(c, "NetworkGraph", func() {
		ui.Checkbox(c, &s.force, "Force layout")
		if input.Button(c, "Reset view", input.ButtonOptions{Variant: input.Secondary}).Clicked() {
			s.graph.ResetView()
			log("NetworkGraph: reset view")
		}
	}, func() {
		frame, show := c1Frame(c, "House", &s.status)
		nodes := []chart.GraphNode{{ID: "a", Label: "Abbot", X: 0, Y: 0}, {ID: "b", Label: "Prior", X: 10, Y: 0}, {ID: "c", Label: "Cellarer", X: 10, Y: 10}, {ID: "d", Label: "Novice", X: 0, Y: 10}, {ID: "e", Label: "Cantor", X: 5, Y: 5}}
		edges := []chart.GraphEdge{{From: "a", To: "b"}, {From: "a", To: "c"}, {From: "c", To: "d"}, {From: "b", To: "e"}, {From: "e", To: "d"}}
		if !show {
			nodes, edges = nil, nil
		}
		c1Log(log, "NetworkGraph", chart.NetworkGraph(c, &s.graph, nodes, edges, chart.NetworkGraphOptions{Frame: frame, Force: s.force}))
	})
}

func c1RealtimeDemo(c *ui.Context, log func(string)) {
	type state struct {
		status string
		feed   bool
		buf    *chart.RealtimeBuffer
		next   time.Time
		tick   int
		view   chart.RealtimeChartState
	}
	s := ui.Local(c.Root(), "120", func() state {
		return state{status: "Ready", feed: true, buf: chart.NewRealtimeBuffer(600, "Bell", "Organ")}
	})
	// The demo appends a sample every 500 ms while feeding.
	if s.feed && !c.Now().Before(s.next) {
		t := time.Date(2026, 10, 5, 9, 0, 0, 0, time.UTC).Add(time.Duration(s.tick) * 500 * time.Millisecond)
		x := float64(s.tick)
		s.buf.Append(t, 5+3*math.Sin(x/6), 4+2*math.Sin(x/4+1))
		s.tick++
		s.next = c.Now().Add(500 * time.Millisecond)
	}
	if s.feed {
		c.After(s.next.Sub(c.Now()))
	}
	c1Section(c, "RealtimeChart", func() { ui.Checkbox(c, &s.feed, "Feed data") }, func() {
		frame, _ := c1Frame(c, "Sound", &s.status)
		r := chart.RealtimeChart(c, &s.view, s.buf, chart.RealtimeChartOptions{Frame: frame, Window: 20 * time.Second})
		c1Log(log, "RealtimeChart", r)
		if r.ValueChanged() {
			log(fmt.Sprintf("RealtimeChart: paused=%v", s.view.Paused))
		}
	})
}
