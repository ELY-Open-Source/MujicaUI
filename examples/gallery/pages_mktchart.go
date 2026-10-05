package main

import "github.com/egoist/mygo/ui"

var mktchartPages = []page{
	{ID: "141", Name: "CandlestickChart", ZhName: "K 线", Category: "Market Charts", Demo: func(c *ui.Context, log func(string)) { mktDemo(c, log, "141") }, Code: mktCode141, Keys: mktKeys},
	{ID: "142", Name: "OHLCChart", ZhName: "美国线", Category: "Market Charts", Demo: func(c *ui.Context, log func(string)) { mktDemo(c, log, "142") }, Code: mktCode142, Keys: mktKeys},
	{ID: "143", Name: "HeikinAshiChart", ZhName: "平均 K 线", Category: "Market Charts", Demo: func(c *ui.Context, log func(string)) { mktDemo(c, log, "143") }, Code: mktCode143, Keys: mktKeys},
	{ID: "144", Name: "LineQuoteChart", ZhName: "分时图", Category: "Market Charts", Demo: func(c *ui.Context, log func(string)) { mktDemo(c, log, "144") }, Code: mktCode144, Keys: mktKeys},
	{ID: "145", Name: "VolumeChart", ZhName: "成交量", Category: "Market Charts", Demo: func(c *ui.Context, log func(string)) { mktDemo(c, log, "145") }, Code: mktCode145, Keys: mktKeys},
	{ID: "146", Name: "DepthChart", ZhName: "深度图", Category: "Market Charts", Demo: func(c *ui.Context, log func(string)) { mktDemo(c, log, "146") }, Code: mktCode146, Keys: mktKeys},
	{ID: "147", Name: "VolumeProfile", ZhName: "成交量分布", Category: "Market Charts", Demo: func(c *ui.Context, log func(string)) { mktDemo(c, log, "147") }, Code: mktCode147, Keys: mktKeys},
	{ID: "148", Name: "MultiChartLayout", ZhName: "多图布局", Category: "Market Charts", Demo: func(c *ui.Context, log func(string)) { mktDemo(c, log, "148") }, Code: mktCode148, Keys: mktKeys},
	{ID: "149", Name: "ChartSync", ZhName: "图表同步", Category: "Market Charts", Demo: func(c *ui.Context, log func(string)) { mktDemo(c, log, "149") }, Code: mktCode149, Keys: mktKeys},
	{ID: "150", Name: "TimeRangeSelector", ZhName: "时间范围", Category: "Market Charts", Demo: func(c *ui.Context, log func(string)) { mktDemo(c, log, "150") }, Code: mktCode150, Keys: mktKeys},
	{ID: "151", Name: "IntervalSelector", ZhName: "周期选择", Category: "Market Charts", Demo: func(c *ui.Context, log func(string)) { mktDemo(c, log, "151") }, Code: mktCode151, Keys: mktKeys},
	{ID: "152", Name: "ChartTypeSwitcher", ZhName: "图表类型", Category: "Market Charts", Demo: func(c *ui.Context, log func(string)) { mktDemo(c, log, "152") }, Code: mktCode152, Keys: mktKeys},
	{ID: "153", Name: "IndicatorOverlay", ZhName: "主图指标", Category: "Market Charts", Demo: func(c *ui.Context, log func(string)) { mktDemo(c, log, "153") }, Code: mktCode153, Keys: mktKeys},
	{ID: "154", Name: "IndicatorPane", ZhName: "指标副图", Category: "Market Charts", Demo: func(c *ui.Context, log func(string)) { mktDemo(c, log, "154") }, Code: mktCode154, Keys: mktKeys},
	{ID: "155", Name: "IndicatorSelector", ZhName: "指标选择", Category: "Market Charts", Demo: func(c *ui.Context, log func(string)) { mktDemo(c, log, "155") }, Code: mktCode155, Keys: mktKeys},
	{ID: "156", Name: "DrawingTools", ZhName: "画线", Category: "Market Charts", Demo: func(c *ui.Context, log func(string)) { mktDemo(c, log, "156") }, Code: mktCode156, Keys: mktKeys},
	{ID: "157", Name: "DrawingToolbar", ZhName: "画线工具栏", Category: "Market Charts", Demo: func(c *ui.Context, log func(string)) { mktDemo(c, log, "157") }, Code: mktCode157, Keys: mktKeys},
	{ID: "158", Name: "CompareSymbol", ZhName: "品种对比", Category: "Market Charts", Demo: func(c *ui.Context, log func(string)) { mktDemo(c, log, "158") }, Code: mktCode158, Keys: mktKeys},
	{ID: "159", Name: "PriceAlertLine", ZhName: "价格预警", Category: "Market Charts", Demo: func(c *ui.Context, log func(string)) { mktDemo(c, log, "159") }, Code: mktCode159, Keys: mktKeys},
	{ID: "160", Name: "CandleCountdown", ZhName: "K 线倒计时", Category: "Market Charts", Demo: func(c *ui.Context, log func(string)) { mktDemo(c, log, "160") }, Code: mktCode160, Keys: mktKeys},
}

const mktKeys = "Tab focuses controls and charts; arrows move focus or chart readouts; Space selects; Home/End jump; +/- zoom; Shift+arrows pan; Alt+1–6 pick drawing tools; Delete removes; Ctrl/Command+Z undoes; Escape closes panels."
