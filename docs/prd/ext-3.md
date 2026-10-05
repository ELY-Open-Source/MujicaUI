# MujicaUI PRD 扩展三：401—500

版本：1.1  
日期：2026-10-05  
规则同扩展一。账户类组件只负责界面与输入，不实现认证、不保存凭据；凭据只经调用方回调传出，不进日志、无障碍值或测试报告。

## 4.27 即时通讯与邮件 · 20 个

| 编号 | 组件 | 首版功能 |
| --- | --- | --- |
| 401 | ChannelList | 频道与私信；未读数、静音、置顶；分组折叠。 |
| 402 | ChatMessage | 连续消息合并、时间戳、编辑标记、回复与回应插槽。 |
| 403 | UnreadDivider | 未读分隔线与计数。 |
| 404 | ReadReceipt | 已发送、已送达、已读；多人已读头像。 |
| 405 | MemberList | 成员分组（在线、离线）；角色。 |
| 406 | UserProfileCard | 头像、名称、状态、角色、操作插槽。 |
| 407 | OnlineStatus | 在线、离开、忙碌、离线；形状加文字。 |
| 408 | StatusSetter | 自定义状态文字与表情；清除时间。 |
| 409 | CallControls | 静音、摄像头、共享、挂断；状态与快捷键。 |
| 410 | VideoCallGrid | 参与者宫格；发言者高亮；帧由调用方提供。 |
| 411 | IncomingCallDialog | 来电接听与拒绝；默认聚焦拒绝可配置。 |
| 412 | EmojiPicker | 分类、搜索、最近使用、肤色。 |
| 413 | PinnedMessages | 置顶消息列表；跳转；取消置顶。 |
| 414 | ThreadPanel | 话题侧栏；原消息与回复。 |
| 415 | MailList | 未读粗体、星标、附件图标；多选；虚拟化。 |
| 416 | MailReader | 邮件头、正文（纯文本与安全子集）、附件。 |
| 417 | MailComposer | 收件人、抄送、密送、主题、正文、附件；草稿保存回调。 |
| 418 | RecipientInput | 收件人标签输入；地址校验；联系人建议。 |
| 419 | QuotedText | 引用文字折叠与展开。 |
| 420 | SnoozePicker | 稍后提醒快捷项与自定义时间。 |

## 4.28 日历与时间 · 20 个

| 编号 | 组件 | 首版功能 |
| --- | --- | --- |
| 421 | CalendarMonthView | 月视图事件格；溢出“+N”；拖动改期。 |
| 422 | CalendarWeekView | 周视图时间网格；全天行；拖动创建与调整。 |
| 423 | CalendarDayView | 日视图；重叠事件并列。 |
| 424 | AgendaView | 日程列表；按日分组。 |
| 425 | YearView | 年视图；事件密度标记。 |
| 426 | EventChip | 事件块；颜色、时间、标题；截断。 |
| 427 | EventEditor | 标题、时间、地点、参与者、提醒、重复；校验。 |
| 428 | EventPopover | 事件摘要浮层；编辑、删除入口。 |
| 429 | CurrentTimeIndicator | 当前时间线；隐藏时停止更新。 |
| 430 | RecurrenceEditor | 每日、每周、每月、每年与自定义；结束条件。 |
| 431 | AvailabilityPicker | 空闲时段选择网格；多人重叠。 |
| 432 | TimezoneSelect | 时区搜索；偏移显示；常用时区。 |
| 433 | ReminderPicker | 提前提醒选项。 |
| 434 | AttendeeList | 参与者与回复状态。 |
| 435 | WeekPicker | 周选择；周起始可配置。 |
| 436 | DurationPicker | 时长输入（时分秒）；范围。 |
| 437 | CronEditor | Cron 表达式编辑与可视化；校验；下次运行时间预览。 |
| 438 | RelativeTime | 相对时间文本；自动更新频率随时间远近；隐藏时停止。 |
| 439 | Countdown | 倒计时；结束回调；格式。 |
| 440 | Stopwatch | 秒表；计次；暂停与复位。 |

## 4.29 仪表盘与开发工具 · 20 个

| 编号 | 组件 | 首版功能 |
| --- | --- | --- |
| 441 | DashboardGrid | 卡片网格；拖动与调整尺寸；布局保存回调。 |
| 442 | StatCard | 统计卡片；值、变化、迷你图；加载与错误。 |
| 443 | AlertList | 告警列表；级别、来源、时间；确认与静默回调。 |
| 444 | IncidentCard | 事件标题、状态、影响、时间线。 |
| 445 | UptimeBar | 按日可用性条；悬停详情。 |
| 446 | ServiceStatus | 服务状态列表；正常、降级、中断。 |
| 447 | LogStream | 实时日志流；暂停；级别过滤；虚拟化。 |
| 448 | TraceWaterfall | 请求瀑布图；跨度层级与耗时。 |
| 449 | ResourceGauge | CPU、内存、磁盘仪表；阈值。 |
| 450 | SystemMonitor | 多资源实时折线组合；采样数据由调用方提供。 |
| 451 | RefreshIntervalSelector | 刷新间隔选择；手动刷新。 |
| 452 | DashboardFilterBar | 时间范围、环境、标签过滤；清除。 |
| 453 | JsonViewer | JSON 树；展开折叠；路径复制；搜索。 |
| 454 | JsonEditor | JSON 文本编辑；格式化；解析错误定位。 |
| 455 | SchemaTree | 数据库、表、列、索引树。 |
| 456 | ApiRequestBuilder | 方法、URL、参数、请求头、正文；发送回调。 |
| 457 | ResponseViewer | 状态码、耗时、大小、正文（JSON 与文本）、响应头。 |
| 458 | KeyValueInput | 键值对行编辑；启用开关；增删。 |
| 459 | EnvironmentSelector | 环境切换；变量预览。 |
| 460 | RegexTester | 正则输入、测试文本、匹配高亮与分组；语法错误提示。 |

## 4.30 设置、账户与引导 · 20 个

| 编号 | 组件 | 首版功能 |
| --- | --- | --- |
| 461 | SettingsLayout | 设置分区导航、分节、设置行（标签、说明、控件）。 |
| 462 | SettingsSearch | 设置项搜索；匹配高亮；跳转。 |
| 463 | ThemeSelector | 系统、浅色、深色预览卡选择。 |
| 464 | AccentColorPicker | 主色候选；对比度不足的候选标记。 |
| 465 | ShortcutRecorder | 录制快捷键；冲突提示；平台修饰键显示。 |
| 466 | KeyboardShortcutsList | 快捷键表；分组；搜索。 |
| 467 | LoginForm | 账号、密码、记住我；提交中与错误；凭据仅经回调传出。 |
| 468 | OAuthButtons | 第三方登录按钮组；名称由调用方提供，不使用商标图形。 |
| 469 | TwoFactorInput | 六位验证码；粘贴；重发倒计时。 |
| 470 | ProfileEditor | 头像、名称、简介；未保存标记。 |
| 471 | AccountSwitcher | 多账户切换；添加账户入口。 |
| 472 | WorkspaceSwitcher | 工作区搜索与切换。 |
| 473 | UserMenu | 用户头像菜单；状态、设置、退出。 |
| 474 | SessionList | 登录会话；设备与位置；注销其他会话回调。 |
| 475 | ApiKeyManager | 密钥列表；创建时只显示一次；撤销确认；明文不入日志。 |
| 476 | UsageQuota | 用量与配额条；重置时间。 |
| 477 | OnboardingWizard | 多步引导；进度；跳过；完成回调。 |
| 478 | Coachmark | 功能提示锚点；步骤导航；Escape 关闭。 |
| 479 | WhatsNewDialog | 更新内容列表；不再显示。 |
| 480 | FeedbackWidget | 反馈类型、文字、截图附件插槽；提交回调。 |

## 4.31 通用补充 · 20 个

| 编号 | 组件 | 首版功能 |
| --- | --- | --- |
| 481 | SplitButton | 主操作与下拉菜单；两部分独立焦点。 |
| 482 | CopyButton | 复制并短暂显示已复制；屏幕阅读器播报。 |
| 483 | HoldToConfirm | 按住确认危险操作；进度；键盘按住空格；减少动画时显示计时文字。 |
| 484 | BulkActionBar | 多选后出现的批量操作条；选中数；清除选择。 |
| 485 | ScrubInput | 拖动标签调整数值；精度与范围；输入模式。 |
| 486 | CurrencyInput | 货币金额输入；币种符号；千分位；精度。 |
| 487 | MaskedInput | 格式掩码输入（电话、卡号等）；输入法安全。 |
| 488 | MentionInput | @ 提及输入；候选；提及标记为整体删除。 |
| 489 | ChoiceChips | 单选或多选筛选芯片；键盘导航。 |
| 490 | Knob | 旋钮；拖动与方向键；刻度与数值。 |
| 491 | VerticalSlider | 纵向滑块；同 Slider 规则。 |
| 492 | ColorPalette | 预设色板；选择与键盘导航；自定义色入口。 |
| 493 | SignaturePad | 手写签名画板；清除；导出路径数据。 |
| 494 | Banner | 页面顶部横幅；信息、警告、错误；操作与关闭。 |
| 495 | UndoToast | 可撤销操作的消息；倒计时；撤销回调。 |
| 496 | Masonry | 瀑布流布局；列数随宽度；稳定键。 |
| 497 | ResizablePanelGroup | 多面板可调尺寸；最小值；折叠；尺寸保存回调。 |
| 498 | DockPanel | 可停靠面板区；拖动标签到边缘停靠；布局状态由调用方保存。 |
| 499 | QRCode | 文本生成二维码；纠错级别；尺寸；离线生成。 |
| 500 | ActivityFeed | 活动流；参与者、动作、对象、时间；按日分组。 |
