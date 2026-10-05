package main

const devCategory = "Dashboards & Dev Tools"

var devtoolsPages = []page{
	{ID: "442", Name: "StatCard", ZhName: "统计卡片", Category: devCategory, Demo: devStatCardDemo, Code: devStatCardCode, Keys: "Tab 聚焦趋势或重试；方向键读取数据；Enter 激活重试。"},
	{ID: "443", Name: "AlertList", ZhName: "告警列表", Category: devCategory, Demo: devAlertListDemo, Code: devAlertListCode, Keys: "Tab 聚焦按钮；Enter 或 Space 激活。"},
	{ID: "444", Name: "IncidentCard", ZhName: "事件卡片", Category: devCategory, Demo: devIncidentCardDemo, Code: devIncidentCardCode, Keys: "展示内容随调用方更新。"},
	{ID: "445", Name: "UptimeBar", ZhName: "可用性条", Category: devCategory, Demo: devUptimeBarDemo, Code: devUptimeBarCode, Keys: "Tab 聚焦图；方向键逐日移动；Home、End 到边界；Escape 清除读取。"},
	{ID: "446", Name: "ServiceStatus", ZhName: "服务状态", Category: devCategory, Demo: devServiceStatusDemo, Code: devServiceStatusCode, Keys: "滚动列表浏览服务。"},
	{ID: "447", Name: "LogStream", ZhName: "日志流", Category: devCategory, Demo: devLogStreamDemo, Code: devLogStreamCode, Keys: "Tab 聚焦暂停和级别；方向键切换级别；Space 激活暂停。"},
	{ID: "448", Name: "TraceWaterfall", ZhName: "调用链瀑布图", Category: devCategory, Demo: devTraceWaterfallDemo, Code: devTraceWaterfallCode, Keys: "Tab 聚焦跨度；上下方向键移动；Enter 读取详情。"},
	{ID: "449", Name: "ResourceGauge", ZhName: "资源仪表", Category: devCategory, Demo: devResourceGaugeDemo, Code: devResourceGaugeCode, Keys: "Tab 聚焦仪表；方向键读取数值与阈值文字。"},
	{ID: "450", Name: "SystemMonitor", ZhName: "系统监视器", Category: devCategory, Demo: devSystemMonitorDemo, Code: devSystemMonitorCode, Keys: "Tab 到图例及图；方向键读取采样点；Home、End 到边界。"},
	{ID: "451", Name: "RefreshIntervalSelector", ZhName: "刷新间隔选择", Category: devCategory, Demo: devRefreshIntervalSelectorDemo, Code: devRefreshIntervalSelectorCode, Keys: "Tab 到选择器与按钮；方向键选择间隔；Enter 激活。"},
	{ID: "452", Name: "DashboardFilterBar", ZhName: "仪表盘过滤栏", Category: devCategory, Demo: devDashboardFilterBarDemo, Code: devDashboardFilterBarCode, Keys: "Tab 遍历控件；方向键选择；Space 切换标签；Escape 关闭弹层。"},
	{ID: "453", Name: "JsonViewer", ZhName: "JSON 查看器", Category: devCategory, Demo: devJsonViewerDemo, Code: devJsonViewerCode, Keys: "Tab 到搜索、复制和树；上下移动；左右展开折叠；Enter 切换节点。"},
	{ID: "454", Name: "JsonEditor", ZhName: "JSON 编辑器", Category: devCategory, Demo: devJsonEditorDemo, Code: devJsonEditorCode, Keys: "Tab 到格式化和文本框；标准文本编辑与剪贴板快捷键。"},
	{ID: "455", Name: "SchemaTree", ZhName: "数据库结构树", Category: devCategory, Demo: devSchemaTreeDemo, Code: devSchemaTreeCode, Keys: "Tab 进入树；上下选择；左右展开折叠。"},
	{ID: "456", Name: "ApiRequestBuilder", ZhName: "API 请求构建器", Category: devCategory, Demo: devApiRequestBuilderDemo, Code: devApiRequestBuilderCode, Keys: "Tab 遍历控件；URL 框 Enter 发送；左右切换标签页。"},
	{ID: "457", Name: "ResponseViewer", ZhName: "响应查看器", Category: devCategory, Demo: devResponseViewerDemo, Code: devResponseViewerCode, Keys: "Tab 到标签页；左右切换正文与响应头；JSON 树沿用查看器键盘操作。"},
	{ID: "458", Name: "KeyValueInput", ZhName: "键值对输入", Category: devCategory, Demo: devKeyValueInputDemo, Code: devKeyValueInputCode, Keys: "Tab 到各字段及按钮；Space 切换启用；Enter 激活增删。"},
}
