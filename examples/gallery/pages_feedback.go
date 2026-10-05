package main

var feedbackPages = []page{
	{
		ID: "091", Name: "Alert", ZhName: "提示条 警告 通知横幅", Category: "Feedback",
		Demo: alertDemo,
		Code: `open := true
r := mujica.Alert(c, &open, "The archive could not be read.", mujica.AlertOptions{
	Severity: mujica.SeverityError,
	Title:    "Import failed",
	Action:   "Retry",
})
if r.Action {
	// retry the import
}`,
		Keys: "Tab moves to the action and the dismiss button; Enter or Space presses them.",
	},
	{
		ID: "092", Name: "Toast", ZhName: "轻提示 消息队列 吐司", Category: "Feedback",
		Demo: toastDemo,
		Code: `if mujica.Button(c, "Delete note", mujica.ButtonOptions{Variant: mujica.Danger}).Clicked() {
	mujica.ShowToast(c, mujica.ToastMessage{ID: "note-1", Message: "Note deleted", Action: "Undo"})
}
// Last in the view:
if r := mujica.Toaster(c); r.Action == "note-1" {
	// restore the note
}`,
		Keys: "Toasts never take the focus. Tab reaches their buttons after the page; Escape closes the focused toast. Hover or focus pauses its timer.",
	},
	{
		ID: "093", Name: "Progress", ZhName: "进度条 环形进度 计量条", Category: "Feedback",
		Demo: progressDemo,
		Code: `mujica.Progress(c, 0.62, mujica.ProgressOptions{Label: "Uploading", ShowValue: true})
mujica.Progress(c, 0, mujica.ProgressOptions{Shape: mujica.ProgressRing, Indeterminate: true, Label: "Connecting"})
mujica.Meter(c, 93, mujica.MeterOptions{
	Max: 100, Label: "Disk", ShowValue: true,
	Levels: &ui.MeterLevels{Warning: 70, Critical: 90},
})`,
		Keys: "Not focusable; assistive technology reads its name and value.",
	},
	{
		ID: "094", Name: "Spinner", ZhName: "加载标记 加载中 转圈", Category: "Feedback",
		Demo: spinnerDemo,
		Code: `mujica.Spinner(c, mujica.SpinnerOptions{Size: 20, Label: "Syncing", ShowLabel: true})`,
		Keys: "Not focusable; read as an indeterminate progress indicator named by Label.",
	},
	{
		ID: "095", Name: "Skeleton", ZhName: "骨架屏 占位", Category: "Feedback",
		Demo: skeletonDemo,
		Code: `ui.Row(c).Gap(12).Children(func() {
	mujica.Skeleton(c, mujica.SkeletonOptions{Shape: mujica.SkeletonCircle})
	ui.Column(c).Grow(1).Children(func() {
		mujica.Skeleton(c, mujica.SkeletonOptions{Lines: 2})
	})
})`,
		Keys: "Decorative; not focusable and hidden from assistive technology.",
	},
	{
		ID: "096", Name: "EmptyState", ZhName: "空状态 无数据 无搜索结果", Category: "Feedback",
		Demo: emptyStateDemo,
		Code: `r := mujica.EmptyState(c, mujica.EmptyStateOptions{
	Kind:   mujica.EmptyNoResults,
	Query:  query,
	Action: "Clear search",
})
if r.Action {
	query = ""
}`,
		Keys: "Tab reaches the action button; Enter or Space presses it.",
	},
	{
		ID: "097", Name: "Result", ZhName: "结果页 成功 失败 无权限", Category: "Feedback",
		Demo: resultDemo,
		Code: `r := mujica.Result(c, mujica.ResultOptions{
	Kind:        mujica.ResultForbidden,
	Description: "You need the Editor role to change billing.",
	Actions:     []string{"Request access", "Go back"},
})
switch r.Action {
case 0: // send the request
case 1: // navigate back
}`,
		Keys: "Tab moves through the actions, primary first; Enter or Space presses one.",
	},
	{
		ID: "098", Name: "StatusIndicator", ZhName: "状态标记 在线 离线", Category: "Feedback",
		Demo: statusDemo,
		Code: `mujica.StatusIndicator(c, mujica.StatusSuccess, mujica.StatusIndicatorOptions{Label: "API online"})
mujica.StatusIndicator(c, mujica.StatusOffline, mujica.StatusIndicatorOptions{})`,
		Keys: "Not focusable; read as a status with its word.",
	},
	{
		ID: "099", Name: "NotificationCenter", ZhName: "通知中心 消息列表", Category: "Feedback",
		Demo: notificationDemo,
		Code: `r := mujica.NotificationCenter(c, &items, mujica.NotificationCenterOptions{})
if r.Action == "a" {
	// open the build log
}`,
		Keys: "The list is one Tab stop. Arrows move among notifications and their buttons; Enter marks the focused one read; Delete or Backspace removes it.",
	},
	{
		ID: "100", Name: "LoadingOverlay", ZhName: "加载遮罩 区域加载", Category: "Feedback",
		Demo: loadingOverlayDemo,
		Code: `r := mujica.LoadingOverlay(c, loading, mujica.LoadingOverlayOptions{Cancel: "Cancel"}, func() {
	ui.Text(c, "Quarterly report")
	mujica.Button(c, "Export", mujica.ButtonOptions{})
})
if r.Canceled {
	loading = false
}`,
		Keys: "While loading, the covered region takes no input; focus inside moves to Cancel, and Escape there cancels. The rest of the window stays usable.",
	},
}
