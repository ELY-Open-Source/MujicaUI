# MujicaUI PRD 扩展二：281—400

版本：1.1  
日期：2026-10-05  
规则同扩展一。动画组件服从 PRD 3.4：状态过渡 120 毫秒、弹层 160 毫秒，不用弹跳与装饰性循环，减少动画模式直接切换或显示静态说明。

## 4.21 Agent 运维与动效 · 20 个

| 编号 | 组件 | 首版功能 |
| --- | --- | --- |
| 281 | InterruptButton | 中断运行中的 Agent；确认可配置；中断中状态。 |
| 282 | TokenUsageChart | 按会话或步骤的 token 用量柱线图；输入输出分色。 |
| 283 | CostBreakdown | 按模型、工具、步骤的费用分解表与合计。 |
| 284 | TraceViewer | 调用链树与时间轴；跨度耗时；选中查看属性。 |
| 285 | EvalResultTable | 评测用例、得分、通过与失败；差异查看回调。 |
| 286 | ToolRegistryPanel | 可用工具列表；参数模式预览；启停。 |
| 287 | DotsLoader | 三点加载；尺寸与说明；减少动画时静态。 |
| 288 | BarsLoader | 条形加载；同上。 |
| 289 | PulseLoader | 单点呼吸式加载；强度受限，不发光；减少动画时静态。 |
| 290 | OrbitLoader | 菱形环绕加载；与 Spinner 同一视觉语言。 |
| 291 | WaveLoader | 波形加载；同上。 |
| 292 | Shimmer | 占位扫光；仅用于加载中；减少动画时静态填充。 |
| 293 | CountUp | 数字从旧值过渡到新值；时长与精度；减少动画时直接显示。 |
| 294 | Typewriter | 逐字显示文本；速度；完成回调；跳过；减少动画时直接显示。 |
| 295 | ShimmerText | 处理中文字扫光（如“思考中”）；减少动画时静态。 |
| 296 | Stagger | 子项依次进场；间隔可配；只在出现时执行一次。 |
| 297 | ReorderTransition | 列表重排时元素移动过渡；稳定键。 |
| 298 | LayoutTransition | 布局尺寸与位置变化过渡（FLIP）。 |
| 299 | Presence | 元素进出场过渡（淡入淡出、滑动、缩放）；退出完成后移除。 |
| 300 | BlinkHighlight | 数值变化单次高亮；方向色；减少动画时改为静态标记。 |

## 4.22 版本控制 · 20 个

| 编号 | 组件 | 首版功能 |
| --- | --- | --- |
| 301 | DiffViewer | 并排与统一视图；行号；增删高亮；折叠未变更区；大文件只构建可见行。 |
| 302 | InlineDiff | 行内词级差异高亮。 |
| 303 | ThreeWayMerge | 本地、基准、远端三栏与结果栏；逐块采纳。 |
| 304 | ConflictResolver | 冲突块列表；采纳当前、传入、两者；已解决计数。 |
| 305 | ChangesList | 暂存与未暂存分组；逐文件与全部暂存；状态标记。 |
| 306 | CommitInput | 提交标题与正文；标题长度提示；提交与修订；防重入。 |
| 307 | CommitGraph | 提交图谱泳道；分支与合并线；虚拟化。 |
| 308 | CommitList | 提交列表；作者、时间、哈希、消息；选择查看。 |
| 309 | BranchSelector | 本地与远端分支搜索选择；新建分支入口。 |
| 310 | BranchList | 分支列表；当前分支、上游、领先落后数；操作回调。 |
| 311 | TagList | 标签列表；注释；跳转提交。 |
| 312 | StashList | 贮藏列表；应用、弹出、删除回调。 |
| 313 | BlameView | 逐行作者与提交；同提交行合并；悬停详情。 |
| 314 | FileHistory | 单文件提交历史与差异入口。 |
| 315 | PullRequestCard | 标题、作者、状态、检查、评审人、差异统计。 |
| 316 | ReviewComment | 行级评审评论；回复；解决与重开。 |
| 317 | DiffStat | +增 −删 统计与比例条。 |
| 318 | GitStatusBadge | M、A、D、U、R、C 状态标记；不只靠颜色。 |
| 319 | CommandBlock | 命令分块（输入、输出、退出码、耗时）；折叠；复制。 |
| 320 | AnsiText | ANSI 16 色、256 色、真彩与粗斜下划线渲染；未知序列原样保留并记录。 |

## 4.23 终端、代码与调试 · 20 个

| 编号 | 组件 | 首版功能 |
| --- | --- | --- |
| 321 | Terminal | 终端网格：ANSI 光标移动与清屏、滚动缓冲、选择复制、输入回调；PTY 由调用方提供。 |
| 322 | TerminalTabs | 多终端标签；重命名；关闭确认。 |
| 323 | TerminalSearch | 滚动缓冲内搜索；匹配数与上下跳转。 |
| 324 | CommandHistory | 命令历史列表；搜索；重放回调。 |
| 325 | LogViewer | 日志级别过滤、跟随尾部、关键字高亮、虚拟化十万行。 |
| 326 | ProcessList | 进程表；CPU、内存排序；结束进程回调与确认。 |
| 327 | CodeViewer | 只读代码视图：行号、基础语法着色、折叠、当前行、选择复制；大文件虚拟化。 |
| 328 | Minimap | 代码缩略图；可见区框；拖动跳转。 |
| 329 | FindWidget | 查找与替换；大小写、整词、正则；匹配数。 |
| 330 | SearchPanel | 多文件搜索结果；按文件分组；匹配预览；包含与排除模式。 |
| 331 | ProblemsPanel | 错误、警告、信息列表；按文件分组；跳转回调。 |
| 332 | OutputPanel | 多通道输出；清空；跟随尾部。 |
| 333 | SymbolOutline | 符号大纲树；类型图标；筛选；跟随光标。 |
| 334 | CompletionMenu | 补全候选；类型图标与详情；方向键与 Tab 接受。 |
| 335 | DebugToolbar | 继续、单步、步入、步出、重启、停止；状态禁用规则。 |
| 336 | BreakpointList | 断点列表；启用开关；条件断点编辑。 |
| 337 | CallStack | 调用栈帧；当前帧；线程分组。 |
| 338 | VariablesPanel | 变量树形展开；值修改回调；变化高亮。 |
| 339 | HexViewer | 十六进制与 ASCII 对照；偏移列；选择同步；虚拟化。 |
| 340 | Flamegraph | 火焰图；点击放大子树；搜索高亮；耗时读数。 |

## 4.24 媒体 · 20 个

| 编号 | 组件 | 首版功能 |
| --- | --- | --- |
| 341 | ImageCropper | 裁剪框；比例锁定；旋转 90 度；输出裁剪矩形。 |
| 342 | ImageAnnotator | 矩形、箭头、文字标注；选择与删除；导出标注数据。 |
| 343 | ImageThumbnail | 缩略图；比例裁切；选中与加载失败状态。 |
| 344 | Lightbox | 全窗图片浏览；左右切换；说明文字；Escape 关闭与焦点恢复。 |
| 345 | VideoPlayer | 帧由调用方提供的播放器外壳：播放、进度、倍速、字幕、全屏、画中画回调。 |
| 346 | VideoScrubber | 时间轴拖动；缓冲段；悬停时间与缩略图插槽。 |
| 347 | VideoThumbnailStrip | 缩略帧条；当前帧；点击跳转。 |
| 348 | AudioPlayer | 播放、暂停、进度、时长；播放由调用方驱动。 |
| 349 | AudioWaveform | 波形绘制；播放进度覆盖；区间选择。 |
| 350 | AudioSpectrum | 频谱柱；数据由调用方提供；隐藏时停止重绘。 |
| 351 | VolumeControl | 音量滑块与静音；方向键。 |
| 352 | PlaybackSpeedControl | 倍速选择 0.5×—2×。 |
| 353 | MediaControls | 上一首、播放、下一首、循环、随机；键盘。 |
| 354 | Playlist | 播放列表；当前项；拖动排序；移除。 |
| 355 | SubtitleEditor | 字幕条目时间与文本编辑；时间重叠校验。 |
| 356 | ScreenRecorderControls | 开始、暂停、停止、计时；录制由调用方实现。 |
| 357 | CameraPreview | 摄像头帧预览外壳；无信号状态；镜像切换。 |
| 358 | MicLevelMeter | 麦克风电平条；峰值保持。 |
| 359 | DeviceSelector | 麦克风、摄像头、扬声器选择；设备列表由调用方提供。 |
| 360 | ImageCompare | 前后对比滑块；键盘调整。 |

## 4.25 项目与任务管理 · 20 个

| 编号 | 组件 | 首版功能 |
| --- | --- | --- |
| 361 | KanbanBoard | 列与卡片；拖动跨列移动与排序；列 WIP 上限；键盘移动卡片。 |
| 362 | TaskList | 任务列表；分组；完成过滤；虚拟化。 |
| 363 | TaskItem | 复选、标题、优先级、截止日期、负责人；行内编辑。 |
| 364 | TaskDetailPanel | 任务详情侧栏；字段编辑；活动记录插槽。 |
| 365 | SubtaskList | 子任务；进度；增删与排序。 |
| 366 | GanttChart | 任务条、依赖线、里程碑；拖动调整日期；今日线；时间缩放。 |
| 367 | Roadmap | 季度或月度路线图；泳道；条目拖动。 |
| 368 | SprintBoard | 迭代信息、剩余天数、看板。 |
| 369 | BurndownChart | 理想线与实际线；剩余工作量。 |
| 370 | PriorityIndicator | 紧急、高、中、低、无；图形与文字。 |
| 371 | StatusSelect | 工作流状态选择；状态分组。 |
| 372 | AssigneePicker | 成员搜索与多选；头像。 |
| 373 | DueDatePicker | 截止日期快捷项（今天、明天、下周）与日历；逾期标记。 |
| 374 | LabelManager | 标签创建、编辑颜色、删除、选择。 |
| 375 | IssueCard | 编号、标题、状态、优先级、负责人、标签。 |
| 376 | IssueIdBadge | ENG-123 样式编号；复制。 |
| 377 | MilestoneProgress | 里程碑完成度与日期。 |
| 378 | WorkloadView | 成员负载条；超载提示。 |
| 379 | TimeTracker | 计时开始停止；累计时长；手动修正。 |
| 380 | PomodoroTimer | 专注与休息周期；暂停；完成回调；隐藏时停止重绘。 |

## 4.26 文件与协作 · 20 个

| 编号 | 组件 | 首版功能 |
| --- | --- | --- |
| 381 | FileExplorer | 文件夹树与文件列表组合；导航历史；选择。 |
| 382 | FileGrid | 文件网格视图；缩略图插槽；多选。 |
| 383 | FileIcon | 按扩展名与类型的图标；文件夹开合。 |
| 384 | FilePreview | 快速预览外壳：图片、文本、代码、不支持类型说明。 |
| 385 | PathBar | 路径分段；点击跳转；长路径中间折叠。 |
| 386 | RenameInline | 行内重命名；默认选中扩展名前部分；校验回调。 |
| 387 | FileOperationProgress | 复制、移动、删除进度；取消；冲突询问。 |
| 388 | TransferQueue | 上传下载队列；暂停、继续、重试、移除。 |
| 389 | StorageUsage | 存储占用分段条与明细。 |
| 390 | RecentFiles | 最近文件；时间分组；清除。 |
| 391 | ArchiveViewer | 压缩包条目树；大小与压缩率；条目由调用方提供。 |
| 392 | PresenceAvatars | 在线协作者头像；溢出计数；跟随入口。 |
| 393 | CommentThread | 评论与回复；编辑、删除、解决；提及插槽。 |
| 394 | Reactions | 表情回应与计数；添加回应选择器。 |
| 395 | ShareDialog | 分享对象、权限、链接复制回调。 |
| 396 | PermissionSelect | 可查看、可评论、可编辑。 |
| 397 | ActivityLog | 操作记录；按日分组；过滤。 |
| 398 | VersionHistory | 版本列表；预览与恢复回调；命名版本。 |
| 399 | LiveIndicator | 实时同步状态；断线提示。 |
| 400 | RemoteCursor | 远程光标与选区标记；姓名标签。 |
