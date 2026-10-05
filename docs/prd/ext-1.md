# MujicaUI PRD 扩展一：101—280

版本：1.1  
日期：2026-10-05  
来源：ELY GPUI 组件清单（`Ely-GPUI-Components/gpui-components.md`）中与本库路线相符的条目。

本扩展沿用 PRD 第 2、3、5、6、8、9、10 节的全部规则。每行一个公开组件；子部件、变体、尺寸不单独计数。图表、行情与媒体组件以 MyGo Painter 绘制，不建图表引擎、音视频解码或富文本编辑引擎：数据、帧与解码由调用方提供。

## 4.12 通用图表一 · 20 个

| 编号 | 组件 | 首版功能 |
| --- | --- | --- |
| 101 | LineChart | 多序列折线；数值与时间轴；缺值断开；悬停读数；序列显隐。 |
| 102 | AreaChart | 面积图；单序列与堆叠；基线可配置；悬停读数。 |
| 103 | BarChart | 纵横向柱状；分组与堆叠；负值；柱标签；键盘逐柱读数。 |
| 104 | PieChart | 饼图与环形；扇区标签；最小角合并为“其他”；中心总计。 |
| 105 | ScatterChart | 散点；点大小映射（气泡）；分组颜色；框选。 |
| 106 | RadarChart | 多维雷达；多序列；刻度环；悬停读数。 |
| 107 | HeatmapChart | 二维矩阵热力；色阶图例；单元格读数；空值样式。 |
| 108 | CalendarHeatmap | 按日贡献图；年份切换；周起始可配置；色阶与读数。 |
| 109 | Histogram | 分箱计数；箱宽与箱数可配置；累计模式。 |
| 110 | BoxPlot | 四分位、须线、离群点；多组对比。 |
| 111 | WaterfallChart | 增减累计；小计与总计柱；正负语义色。 |
| 112 | FunnelChart | 漏斗阶段；转化率；阶段读数。 |
| 113 | BulletChart | 实际值、目标值、区间带；纵横向。 |
| 114 | Sparkline | 行内迷你折线、柱形与面积；末值标记；无坐标轴。 |
| 115 | Gauge | 半圆或整圆仪表；阈值分段；数值与单位；不能只靠颜色表达区间。 |
| 116 | Treemap | 层级矩形；点击下钻与返回；标签按面积省略。 |
| 117 | Sunburst | 层级环形；下钻；路径面包屑。 |
| 118 | SankeyChart | 节点与流量带；悬停高亮路径；节点排序稳定。 |
| 119 | NetworkGraph | 节点与边；调用方给定或简单力导布局；拖动节点；缩放平移。 |
| 120 | RealtimeChart | 实时追加数据；固定窗口滚动；暂停；隐藏时停止重绘。 |

## 4.13 通用图表二与图表部件 · 20 个

| 编号 | 组件 | 首版功能 |
| --- | --- | --- |
| 121 | ChartLegend | 序列图例；点击显隐；键盘切换；横纵布局与溢出折叠。 |
| 122 | ChartTooltip | 跟随指针或焦点的读数浮层；多序列同点读数；边缘避让。 |
| 123 | ChartCrosshair | 十字光标；轴上读数标签；吸附最近数据点。 |
| 124 | ChartAxis | 数值、时间、类别轴；刻度自动与指定；格式化函数；对数轴。 |
| 125 | ChartGrid | 网格线；主次刻度；按设备像素对齐。 |
| 126 | ChartBrush | 区间选择与缩放；拖动、调整两端；键盘移动；回调区间。 |
| 127 | ChartAnnotation | 水平线、垂直线、区间带与文字标注；拖动可选。 |
| 128 | ChartExport | 以当前主题渲染图表为 PNG 字节；尺寸与缩放可配置；不写文件。 |
| 129 | ViolinPlot | 核密度轮廓；中位数与四分位；多组。 |
| 130 | ParallelCoordinates | 多维平行轴；轴上刷选；线高亮。 |
| 131 | ChordDiagram | 关系矩阵弦图；悬停高亮；组弧标签。 |
| 132 | SparkBar | 行内迷你柱；正负值；高亮末项。 |
| 133 | TrendIndicator | 趋势箭头与变化值；正负语义；不只靠颜色。 |
| 134 | MetricCard | 指标名、值、变化、迷你图插槽；加载与错误状态。 |
| 135 | ProgressChart | 多目标完成度环或条；目标线；百分比读数。 |
| 136 | StackedPercentBar | 百分比堆叠条；分段标签；悬停读数。 |
| 137 | DualAxisChart | 左右双轴组合（柱与线）；各自刻度与格式。 |
| 138 | DensityPlot | 一维核密度曲线；带宽可配置；多组叠加。 |
| 139 | QuadrantChart | 四象限散点；象限标题与中线可配置。 |
| 140 | ChartEmptyState | 图表的无数据、加载、错误与重试占位；保持图表尺寸。 |

## 4.14 行情图表与技术分析 · 20 个

| 编号 | 组件 | 首版功能 |
| --- | --- | --- |
| 141 | CandlestickChart | K 线；涨跌色可配置（红涨或绿涨）；缩放平移；万根级只绘可见区。 |
| 142 | OHLCChart | 美国线；与 K 线共用数据接口。 |
| 143 | HeikinAshiChart | 平均 K 线；由调用方 OHLC 计算；与原 K 线可切换。 |
| 144 | LineQuoteChart | 分时图；昨收基线；均价线；盘前盘后分段。 |
| 145 | VolumeChart | 成交量柱；涨跌着色；与主图共享时间轴。 |
| 146 | DepthChart | 买卖深度累计面积；中间价；悬停读数。 |
| 147 | VolumeProfile | 价格分布成交量；价值区与控制点。 |
| 148 | MultiChartLayout | 1、2、4、6 宫格多图；独立品种与周期；布局切换不丢状态。 |
| 149 | ChartSync | 多图十字光标与时间范围同步；可逐图解除。 |
| 150 | TimeRangeSelector | 1D、5D、1M、3M、YTD、1Y、5Y、自定义；键盘导航。 |
| 151 | IntervalSelector | 1m 至 1M 周期；分组；收藏周期。 |
| 152 | ChartTypeSwitcher | K 线、美国线、折线、面积、平均 K 线切换。 |
| 153 | IndicatorOverlay | 主图叠加 MA、EMA、BOLL、VWAP；参数可配；计算由调用方或内置辅助函数。 |
| 154 | IndicatorPane | 副图 MACD、RSI、KDJ、OBV；多副图堆叠与高度调整。 |
| 155 | IndicatorSelector | 指标搜索、添加、移除、参数编辑。 |
| 156 | DrawingTools | 趋势线、水平线、斐波那契、矩形、文字标注；选中、移动、删除、撤销。 |
| 157 | DrawingToolbar | 画线工具栏；当前工具；锁定与清除；键盘快捷键提示。 |
| 158 | CompareSymbol | 叠加对比品种；百分比归一；序列颜色与移除。 |
| 159 | PriceAlertLine | 价格预警线；拖动调整；触发回调；不自行通知。 |
| 160 | CandleCountdown | 当前 K 线收盘倒计时；随周期变化；隐藏时停止计时重绘。 |

## 4.15 报价、盘口与交易 · 20 个

| 编号 | 组件 | 首版功能 |
| --- | --- | --- |
| 161 | TickerTape | 滚动行情条；悬停与聚焦暂停；减少动画时静态分页。 |
| 162 | QuoteCard | 代码、名称、现价、涨跌、迷你图插槽。 |
| 163 | PriceText | 价格精度与千分位；涨跌色；变化闪烁一次（减少动画时无闪烁）。 |
| 164 | PriceChangeBadge | 涨跌额与幅度；箭头与符号；不只靠颜色。 |
| 165 | OrderBook | 买卖档位；深度条；价格聚合档位；点击价格回传。 |
| 166 | Level2Quotes | 多档买卖与委托量；逐档变化高亮。 |
| 167 | TimeAndSales | 逐笔成交；大单标记；虚拟化；暂停跟随。 |
| 168 | DepthLadder | 价格阶梯（DOM）；买卖量列；当前价居中；点击下单回调。 |
| 169 | SpreadIndicator | 买卖价差与百分比；阈值提示。 |
| 170 | BidAskBar | 买卖力量对比条；数值标签。 |
| 171 | OrderEntry | 限价、市价、止损、条件单；数量、价格、有效期；校验与预估金额；提交防重入。 |
| 172 | QuickTradeButtons | 快速买卖；数量预设；危险操作确认可配置。 |
| 173 | PositionTable | 持仓；成本、市值、浮动盈亏；平仓操作回调；稳定行 ID。 |
| 174 | OrderTable | 委托；状态筛选；撤单回调；部分成交进度。 |
| 175 | TradeHistoryTable | 成交记录；时间范围筛选；导出回调。 |
| 176 | PnLDisplay | 已实现、未实现、当日盈亏；正负语义；隐私遮蔽切换。 |
| 177 | LeverageSlider | 杠杆档位；风险提示区间；确认后生效。 |
| 178 | MarginIndicator | 保证金率；警戒与强平线；不只靠颜色。 |
| 179 | RiskMeter | 风险等级刻度与说明文字。 |
| 180 | OrderConfirmDialog | 下单摘要与风险提示；默认聚焦取消；二次确认可配置。 |

## 4.16 市场与资产 · 20 个

| 编号 | 组件 | 首版功能 |
| --- | --- | --- |
| 181 | OptionChain | 行权价居中的认购认沽表；到期日切换；价内高亮。 |
| 182 | GreeksTable | Delta、Gamma、Theta、Vega、Rho；精度与说明。 |
| 183 | PayoffDiagram | 期权组合到期损益；盈亏平衡点；当前价标记。 |
| 184 | Watchlist | 自选分组；拖动排序；实时刷新不丢选择；列可配置。 |
| 185 | SymbolSearch | 代码与名称搜索；市场筛选；最近搜索；异步结果由调用方提供。 |
| 186 | SymbolBadge | 代码、市场、类型标签；可点击。 |
| 187 | MarketHeatmap | 板块与个股热力；面积映射市值；涨跌色阶；下钻。 |
| 188 | Screener | 条件构建（字段、比较、值）；结果表；保存条件回调。 |
| 189 | MarketOverview | 指数卡片组与涨跌家数；时间戳。 |
| 190 | IndexCard | 指数现价、涨跌、迷你图。 |
| 191 | MarketStatus | 开盘、收盘、盘前、盘后、休市；下一状态时间。 |
| 192 | TradingSessionClock | 多市场交易时段时钟；时区由调用方指定。 |
| 193 | EconomicCalendar | 经济事件；重要性；预期、前值、实际；国家筛选。 |
| 194 | EarningsCalendar | 财报日程；盘前盘后；预期与实际。 |
| 195 | NewsFeed | 新闻流；来源、时间、相关代码；已读状态；虚拟化。 |
| 196 | SentimentGauge | 情绪刻度（恐惧到贪婪）；数值与文字。 |
| 197 | PortfolioSummary | 总资产、当日盈亏、累计收益；分账户。 |
| 198 | AssetAllocationChart | 资产配置环形与列表；目标与偏离。 |
| 199 | PerformanceChart | 收益曲线、基准对比、最大回撤区间。 |
| 200 | CurrencyConverter | 币种选择与互换；汇率由调用方提供；精度规则。 |

## 4.17 AI 对话：消息 · 20 个

| 编号 | 组件 | 首版功能 |
| --- | --- | --- |
| 201 | ChatContainer | 消息区、输入区、侧栏插槽；新消息自动滚底，用户上翻时停止跟随。 |
| 202 | MessageList | 虚拟化消息列表；稳定消息 ID；新消息提示；按日期分隔。 |
| 203 | MessageBubble | 用户、助手、系统、错误四种；头像与名称插槽；宽度限制。 |
| 204 | MessageHeader | 模型名、时间、状态。 |
| 205 | StreamingText | 流式追加文本；平滑出字可关闭；流式光标；减少动画时直接显示。 |
| 206 | MarkdownView | Markdown 子集渲染：标题、段落、强调、列表、引用、链接、代码、表格、分隔线；流式增量渲染不重排已完成块。 |
| 207 | CodeBlock | 语言标签、行号、基础语法着色、复制、折叠、换行切换。 |
| 208 | TableBlock | 消息内表格；横向滚动；复制为 TSV。 |
| 209 | ImageGrid | 1 至 9 张图片网格；点击打开 ImageViewer；加载与失败。 |
| 210 | FileMessage | 文件名、类型图标、大小、打开与下载回调。 |
| 211 | AudioMessage | 语音条；波形；播放进度；时长；播放由调用方驱动。 |
| 212 | LinkPreviewCard | 标题、描述、站点、缩略图；元数据由调用方提供。 |
| 213 | QuoteReply | 引用原消息摘要；点击跳转原消息。 |
| 214 | MessageActions | 复制、重新生成、编辑、赞、踩、朗读、分享；按角色显示。 |
| 215 | MessageEditor | 编辑已发送消息；保存与取消；保存后产生新分支回调。 |
| 216 | BranchNavigator | 分支切换（< 2/3 >）；键盘操作。 |
| 217 | ErrorMessage | 失败原因、重试；限流倒计时（RateLimit）；不吞掉原错误文本。 |
| 218 | ThinkingBlock | 可折叠思考过程；思考中动画与耗时；默认折叠可配置。 |
| 219 | ScrollToBottomButton | 离开底部时出现；未读计数；点击回到底部。 |
| 220 | TypingIndicator | 对方正在输入三点动画；减少动画时显示静态文字。 |

## 4.18 AI 对话：输入与引用 · 20 个

| 编号 | 组件 | 首版功能 |
| --- | --- | --- |
| 221 | PromptComposer | 多行输入；Enter 发送、Shift+Enter 换行可配置；输入法组合期间不发送；附件与工具插槽。 |
| 222 | AttachmentChip | 附件名、类型、大小、上传进度、移除；失败重试。 |
| 223 | PasteImagePreview | 粘贴或拖入图片的预览缩略图；移除。 |
| 224 | DragDropOverlay | 拖入文件时的全区域遮罩与提示；类型检查回调。 |
| 225 | VoiceInputButton | 录音开始与停止；录音时长；权限失败提示；录音由调用方实现。 |
| 226 | VoiceWaveform | 输入电平波形；静音提示。 |
| 227 | SlashCommandMenu | 输入“/”弹出命令；筛选；方向键选择；Tab 补全。 |
| 228 | ContextMentionMenu | 输入“@”选择文件、文档、网页等上下文；分组；异步结果。 |
| 229 | ContextChips | 已附加上下文列表；移除；溢出折叠。 |
| 230 | ModelSelector | 模型列表；能力标签与上下文长度；不可用项说明。 |
| 231 | ModeSelector | 对话、搜索、深度研究、Agent 模式切换。 |
| 232 | ToolToggleMenu | 联网、代码执行等工具开关；需授权项标记。 |
| 233 | ParameterPanel | 温度、Top-p、最大长度等参数；范围与默认值；重置。 |
| 234 | TokenCounter | 当前输入 token 数与上限；估算值由调用方提供。 |
| 235 | ContextWindowMeter | 上下文使用率条；分段（系统、历史、附件）；阈值警告。 |
| 236 | SendButton | 发送与停止两态；禁用条件；快捷键提示。 |
| 237 | CitationBadge | 行内引用编号；悬停预览来源；点击定位来源。 |
| 238 | SourceCard | 来源标题、站点、摘要、时间；打开回调。 |
| 239 | SourcesPanel | 来源列表；按引用顺序；与正文引用联动高亮。 |
| 240 | SearchProgress | 搜索中关键词与已读来源数；阶段列表；完成后折叠。 |

## 4.19 AI 对话：会话与引导 · 20 个

| 编号 | 组件 | 首版功能 |
| --- | --- | --- |
| 241 | ConversationList | 按今天、昨天、7 天内、更早分组；置顶；虚拟化。 |
| 242 | ConversationItem | 标题、时间、重命名、删除、置顶、导出回调；行内重命名。 |
| 243 | ConversationSearch | 标题与内容搜索；匹配高亮；异步结果。 |
| 244 | ProjectList | 项目与文件夹；会话归属；拖动移入。 |
| 245 | ProjectKnowledgePanel | 项目知识文件列表；添加、移除、容量提示。 |
| 246 | WelcomeScreen | 问候、能力卡片、起始提示插槽。 |
| 247 | SuggestionChips | 起始提示与追问建议；点击填入或直接发送可配置。 |
| 248 | PromptLibrary | 提示模板分类、搜索、预览、插入；变量占位符填写。 |
| 249 | CapabilityCards | 能力说明卡片；点击示例。 |
| 250 | FeedbackForm | 赞踩后的原因选项与文字；提交与取消。 |
| 251 | DateSeparator | 消息日期分隔；今天、昨天与日期格式本地化。 |
| 252 | ThinkingIndicator | 思考中状态动画与阶段文字；减少动画时静态。 |
| 253 | StopGeneratingButton | 停止生成；生成结束后隐藏；快捷键。 |
| 254 | SystemPromptEditor | 系统提示多行编辑；字数；恢复默认；未保存标记。 |
| 255 | CostEstimator | 输入、输出 token 与费用估算；单价由调用方提供。 |
| 256 | SharedConversationView | 只读会话视图；分享者与时间；复制链接回调。 |
| 257 | ConversationExport | 导出为 Markdown、JSON、纯文本；范围选择；内容由组件生成，写文件由调用方完成。 |
| 258 | RegenerateMenu | 重新生成；换模型重试；更长、更短等变体选项。 |
| 259 | ArtifactCard | 消息内产物卡片（代码、文档、网页）；打开产物面板回调。 |
| 260 | ChatEmptyState | 新会话、无搜索结果、离线三种空状态。 |

## 4.20 AI Agent 与工具调用 · 20 个

| 编号 | 组件 | 首版功能 |
| --- | --- | --- |
| 261 | ToolCallCard | 工具名、参数、结果、耗时、状态；参数与结果折叠；JSON 高亮。 |
| 262 | ToolCallGroup | 连续调用折叠为一组；总耗时与失败数。 |
| 263 | ToolApprovalDialog | 执行前审批；参数预览；允许、拒绝；默认聚焦拒绝。 |
| 264 | PermissionPrompt | 允许一次、始终允许、拒绝；作用范围说明。 |
| 265 | AgentStepList | 执行步骤时间线；当前步；步骤展开详情；失败步重试回调。 |
| 266 | AgentPlan | 计划清单；实时勾选；新增与跳过项；进度。 |
| 267 | AgentStatus | 规划中、执行中、等待、完成、失败；不只靠颜色；状态动画可关闭。 |
| 268 | AgentProgress | 总体进度与当前动作文字；不确定模式。 |
| 269 | SubAgentTree | 子 Agent 层级；各自状态；展开查看。 |
| 270 | FileChangeCard | 文件变更预览；增删统计；接受与拒绝；打开差异。 |
| 271 | MultiFileDiffReview | 多文件审阅；逐文件接受拒绝；全部接受；进度。 |
| 272 | CommandExecutionCard | 命令、工作目录、输出（ANSI）、退出码、耗时；中止回调。 |
| 273 | ScreenshotStream | Agent 画面帧序列查看；帧由调用方提供；暂停与逐帧。 |
| 274 | ArtifactPanel | 产物面板；代码、文档、网页预览插槽；复制与导出回调。 |
| 275 | ArtifactVersionSwitcher | 产物版本切换与差异入口。 |
| 276 | SandboxStatus | 沙箱运行状态、资源与剩余时间。 |
| 277 | MCPServerList | 服务器列表；连接状态；启停与重连回调；工具数。 |
| 278 | MemoryPanel | 记忆条目列表；搜索；编辑与删除。 |
| 279 | CheckpointList | 检查点；回滚确认；当前点标记。 |
| 280 | HumanInputRequest | Agent 请求用户输入；问题、选项与自由回答；超时提示。 |
