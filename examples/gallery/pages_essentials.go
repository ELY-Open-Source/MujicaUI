package main

var essentialsPages = []page{
	{
		ID: "481", Name: "SplitButton", ZhName: "拆分按钮 主操作 下拉", Category: "Essentials",
		Demo: splitButtonDemo,
		Code: `view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := overlay.SplitButton(c, "Save", overlay.SplitButtonOptions{}, func(m *overlay.PopupMenu) {
		if m.Item("Save as…", overlay.PopupMenuItemOptions{}) {
			// ask for a file name
		}
	})
	if r.Changed() {
		// save
	}
}`,
		Keys: "Tab reaches the main half, then the menu half. Enter or Space runs the action; on the menu half, Enter, Space or Down opens the menu.",
	},
	{
		ID: "482", Name: "CopyButton", ZhName: "复制按钮 剪贴板", Category: "Essentials",
		Demo: copyButtonDemo,
		Code: `view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	input.CopyButton(c, "https://example.com/invite/42", input.CopyButtonOptions{Label: "Copy link"})
}`,
		Keys: "Tab focuses; Enter or Space copies. The copied state is announced.",
	},
	{
		ID: "483", Name: "HoldToConfirm", ZhName: "按住确认 危险操作", Category: "Essentials",
		Demo: holdToConfirmDemo,
		Code: `view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if input.HoldToConfirm(c, "Delete project", input.HoldToConfirmOptions{}).Changed() {
		// delete the project
	}
}`,
		Keys: "Focus it and hold Space for the full duration; a click, Enter or a short Space does nothing.",
	},
	{
		ID: "484", Name: "BulkActionBar", ZhName: "批量操作条 多选", Category: "Essentials",
		Demo: bulkActionBarDemo,
		Code: `selected := 3
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := input.BulkActionBar(c, selected, input.BulkActionBarOptions{}, func() {
		input.Button(c, "Archive", input.ButtonOptions{Variant: input.Secondary})
	})
	if r.Changed() {
		selected = 0
	}
}`,
		Keys: "Tab moves through the bar's actions; Escape inside it clears the selection.",
	},
	{
		ID: "485", Name: "ScrubInput", ZhName: "拖动数值 拖动标签", Category: "Essentials",
		Demo: scrubInputDemo,
		Code: `opacity := 0.6
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	input.ScrubInput(c, "Opacity", &opacity, input.ScrubInputOptions{Max: 1, Step: 0.1})
}`,
		Keys: "The label is pointer-only; in the field, Up and Down step and typing sets the value.",
	},
	{
		ID: "486", Name: "CurrencyInput", ZhName: "货币输入 金额 千分位", Category: "Essentials",
		Demo: currencyInputDemo,
		Code: `cents := int64(123450)
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	input.CurrencyInput(c, &cents, input.CurrencyInputOptions{Symbol: "$", Label: "Price"})
}`,
		Keys: "Tab focuses the field, which shows plain digits while focused and groups thousands when left.",
	},
	{
		ID: "487", Name: "MaskedInput", ZhName: "掩码输入 电话 卡号", Category: "Essentials",
		Demo: maskedInputDemo,
		Code: `phone := ""
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	input.MaskedInput(c, &phone, input.MaskedInputOptions{Mask: "(999) 999-9999", Label: "Phone"})
}`,
		Keys: "Arrows, Home and End move the caret; Backspace and Delete remove; paste keeps the characters that fit; Enter submits.",
	},
	{
		ID: "488", Name: "MentionInput", ZhName: "提及输入 @ 候选", Category: "Essentials",
		Demo: mentionInputDemo,
		Code: `text, mentions := "", []input.Mention(nil)
people := []input.MentionCandidate{{ID: "u1", Name: "Alice"}, {ID: "u2", Name: "Bob"}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	input.MentionInput(c, &text, &mentions, input.MentionInputOptions{
		Label: "Comment",
		Candidates: func(q string) []input.MentionCandidate {
			var out []input.MentionCandidate
			for _, p := range people {
				if strings.HasPrefix(strings.ToLower(p.Name), strings.ToLower(q)) {
					out = append(out, p)
				}
			}
			return out
		},
	})
}`,
		Keys: "After @, Up and Down pick a candidate, Enter or Tab inserts it, Escape closes the list. Arrows step over a mention; Backspace removes it whole.",
	},
	{
		ID: "489", Name: "ChoiceChips", ZhName: "筛选芯片 单选 多选", Category: "Essentials",
		Demo: choiceChipsDemo,
		Code: `var status []string
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	input.ChoiceChips(c, &status, []input.ToggleItem[string]{
		{Value: "open", Label: "Open"}, {Value: "closed", Label: "Closed"},
	}, input.ChoiceChipsOptions{Multiple: true, Label: "Status"})
}`,
		Keys: "The chips are one Tab stop; Left and Right move, skipping disabled chips; Space or Enter toggles.",
	},
	{
		ID: "490", Name: "Knob", ZhName: "旋钮 调节", Category: "Essentials",
		Demo: knobDemo,
		Code: `gain := 35.0
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	input.Knob(c, &gain, input.KnobOptions{Label: "Gain", Ticks: 11})
}`,
		Keys: "Arrows step the value; Home and End go to the ends; drag up or down with the pointer.",
	},
	{
		ID: "491", Name: "VerticalSlider", ZhName: "纵向滑块", Category: "Essentials",
		Demo: verticalSliderDemo,
		Code: `level := 40.0
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	input.VerticalSlider(c, &level, input.VerticalSliderOptions{Step: 10, Label: "Level"})
}`,
		Keys: "Up and Right raise the value, Down and Left lower it, Home and End go to the ends.",
	},
	{
		ID: "492", Name: "ColorPalette", ZhName: "预设色板 颜色", Category: "Essentials",
		Demo: colorPaletteDemo,
		Code: `color := ui.Hex("#7D2034")
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	input.ColorPalette(c, &color, []input.ColorSwatch{
		{Color: ui.Hex("#7D2034"), Name: "Wine"},
		{Color: ui.Hex("#AD9164"), Name: "Old gold"},
		{Color: ui.Hex("#355E80"), Name: "Slate blue"},
	}, input.ColorPaletteOptions{Custom: true, Label: "Tag color"})
}`,
		Keys: "The palette is one Tab stop; Left and Right move by one swatch, Up and Down by a row; Enter or Space chooses.",
	},
	{
		ID: "493", Name: "SignaturePad", ZhName: "签名板 手写", Category: "Essentials",
		Demo: signaturePadDemo,
		Code: `var strokes [][]input.SignaturePoint
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if input.SignaturePad(c, &strokes, input.SignaturePadOptions{}).Changed() {
		_ = input.SignaturePath(strokes) // SVG path data to store
	}
}`,
		Keys: "Tab focuses the pad; Delete or Backspace clears it. Drawing needs the pointer.",
	},
	{
		ID: "494", Name: "Banner", ZhName: "横幅 页面提示", Category: "Essentials",
		Demo: bannerDemo,
		Code: `open := true
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := input.Banner(c, &open, "Maintenance tonight at 22:00.", input.BannerOptions{
		Severity: core.SeverityWarning,
		Action:   "Details",
	})
	if r.Action {
		// show the details
	}
}`,
		Keys: "Tab reaches the action and the dismiss button; Enter or Space presses them.",
	},
	{
		ID: "495", Name: "UndoToast", ZhName: "撤销提示 倒计时", Category: "Essentials",
		Demo: undoToastDemo,
		Code: `open := true
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	feedback.UndoToast(c, &open, "Note deleted", feedback.UndoToastOptions{
		OnUndo: func() {
			// restore the note
		},
	})
}`,
		Keys: "It never takes the focus. Tab reaches Undo and the close button after the page; Escape inside closes it. Hover or focus pauses the countdown.",
	},
	{
		ID: "496", Name: "Masonry", ZhName: "瀑布流 卡片", Category: "Essentials",
		Demo: masonryDemo,
		Code: `titles := []string{"Sea", "Forest", "Desert", "Glacier"}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	input.Masonry(c, len(titles), input.MasonryOptions{
		Key: func(i int) string { return titles[i] },
	}, func(i int) {
		ui.Text(c, titles[i])
	}).Size(380, 280)
}`,
		Keys: "Scroll with the wheel or the scrollbar; only cards in view are built.",
	},
	{
		ID: "497", Name: "ResizablePanelGroup", ZhName: "可调面板 分隔条", Category: "Essentials",
		Demo: resizablePanelGroupDemo,
		Code: `var sizes []float32
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := input.ResizablePanelGroup(c, &sizes, []input.ResizablePanel{
		{Label: "Files", Collapsible: true, Build: func() { ui.Text(c, "Files") }},
		{Label: "Editor", Build: func() { ui.Text(c, "Editor") }},
	}, input.ResizablePanelGroupOptions{})
	if r.Committed {
		// save sizes
	}
}`,
		Keys: "Tab focuses a divider; arrows move it, Home and End go to the limits, Enter collapses or restores a collapsible neighbor.",
	},
	{
		ID: "498", Name: "DockPanel", ZhName: "停靠面板 标签", Category: "Essentials",
		Demo: dockPanelDemo,
		Code: `var layout input.DockLayout
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := input.DockPanel(c, &layout, []input.DockItem{
		{ID: "files", Title: "Files", Build: func() { ui.Text(c, "Files") }},
		{ID: "terminal", Title: "Terminal", Build: func() { ui.Text(c, "$") }},
	}, input.DockPanelOptions{})
	if r.Changed() {
		// save layout
	}
}`,
		Keys: "Each zone's tabs are one Tab stop; Left and Right switch tabs. Shift+Alt+Left, Right or Down docks the focused tab to that edge; Shift+Alt+Up returns it to the center.",
	},
	{
		ID: "499", Name: "QRCode", ZhName: "二维码", Category: "Essentials",
		Demo: qrCodeDemo,
		Code: `view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	input.QRCode(c, "https://example.com/mujica", input.QRCodeOptions{Level: input.QRQuartile})
}`,
		Keys: "Not focusable; assistive technology reads its name.",
	},
	{
		ID: "500", Name: "ActivityFeed", ZhName: "活动流 动态", Category: "Essentials",
		Demo: activityFeedDemo,
		Code: `events := []feedback.Activity{
	{ID: "1", Actor: "Lin", Action: "commented on", Object: "Quarterly report", Time: time.Now()},
}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	feedback.ActivityFeed(c, events, feedback.ActivityFeedOptions{}).Size(380, 280)
}`,
		Keys: "Scroll with the wheel or the scrollbar; only rows in view are built.",
	},
}
