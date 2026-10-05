package main

const filesCategory = "Files & Collaboration"

var filesPages = []page{
	{ID: "381", Name: "FileExplorer", ZhName: "文件浏览器", Category: filesCategory, Demo: fileExplorerDemo, Code: filesFileExplorerCode, Keys: "Tab 切换；方向键选择；Enter 打开；macOS Command+[ / ] / ↑，其他平台 Alt+← / → / ↑。"},
	{ID: "382", Name: "FileGrid", ZhName: "文件网格", Category: filesCategory, Demo: fileGridDemo, Code: filesFileGridCode, Keys: "方向键选择；Command/Ctrl+点击多选；Shift 扩展；Enter 打开。"},
	{ID: "383", Name: "FileIcon", ZhName: "文件图标", Category: filesCategory, Demo: fileIconDemo, Code: filesFileIconCode, Keys: "展示组件，沿用父级焦点。"},
	{ID: "384", Name: "FilePreview", ZhName: "文件预览", Category: filesCategory, Demo: filePreviewDemo, Code: filesFilePreviewCode, Keys: "Escape 关闭；Tab 到操作；文本支持选择与复制。"},
	{ID: "385", Name: "PathBar", ZhName: "路径栏", Category: filesCategory, Demo: pathBarDemo, Code: filesPathBarCode, Keys: "Tab 进入；左右移动；Enter 激活；Escape 关闭折叠菜单。"},
	{ID: "386", Name: "RenameInline", ZhName: "行内重命名", Category: filesCategory, Demo: renameInlineDemo, Code: filesRenameInlineCode, Keys: "Enter 确认；Escape 取消；双击开始编辑；失焦确认有效名称。"},
	{ID: "387", Name: "FileOperationProgress", ZhName: "文件操作进度", Category: filesCategory, Demo: fileOperationProgressDemo, Code: filesFileOperationProgressCode, Keys: "Tab 到取消、全选规则和冲突选项；Enter / Space 激活。"},
	{ID: "388", Name: "TransferQueue", ZhName: "传输队列", Category: filesCategory, Demo: transferQueueDemo, Code: filesTransferQueueCode, Keys: "Tab 在可见行操作间移动；Enter / Space 激活。"},
	{ID: "389", Name: "StorageUsage", ZhName: "存储占用", Category: filesCategory, Demo: storageUsageDemo, Code: filesStorageUsageCode, Keys: "Tab 聚焦分段条；左右循环读取分类与剩余空间。"},
	{ID: "390", Name: "RecentFiles", ZhName: "最近文件", Category: filesCategory, Demo: recentFilesDemo, Code: filesRecentFilesCode, Keys: "上下选择文件；Enter 打开；Tab 到清除和移除按钮。"},
	{ID: "391", Name: "ArchiveViewer", ZhName: "压缩包浏览", Category: filesCategory, Demo: archiveViewerDemo, Code: filesArchiveViewerCode, Keys: "上下选择；左右展开折叠；Enter 打开；Tab 到列操作。"},
	{ID: "392", Name: "PresenceAvatars", ZhName: "协作头像", Category: filesCategory, Demo: presenceAvatarsDemo, Code: filesPresenceAvatarsCode, Keys: "Tab 进入；左右移动；Enter 激活；Escape 关闭溢出名单。"},
	{ID: "394", Name: "Reactions", ZhName: "表情回应", Category: filesCategory, Demo: reactionsDemo, Code: filesReactionsCode, Keys: "Tab 到回应与添加；Enter / Space 激活；选择器左右导航；Escape 关闭。"},
	{ID: "395", Name: "ShareDialog", ZhName: "分享对话框", Category: filesCategory, Demo: shareDialogDemo, Code: filesShareDialogCode, Keys: "Tab 在对话框内移动；Enter 邀请或激活操作；Escape 关闭并恢复焦点。"},
	{ID: "396", Name: "PermissionSelect", ZhName: "权限选择", Category: filesCategory, Demo: permissionSelectDemo, Code: filesPermissionSelectCode, Keys: "Enter 打开；上下移动；Enter 确认；Escape 关闭。"},
	{ID: "399", Name: "LiveIndicator", ZhName: "同步状态", Category: filesCategory, Demo: liveIndicatorDemo, Code: filesLiveIndicatorCode, Keys: "Tab 到重试；Enter / Space 激活。"},
	{ID: "400", Name: "RemoteCursor", ZhName: "远程光标", Category: filesCategory, Demo: remoteCursorDemo, Code: filesRemoteCursorCode, Keys: "输入穿透至下方编辑器。"},
}
