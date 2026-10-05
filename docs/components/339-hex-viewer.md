# 339 HexViewer · 十六进制查看器

## 用途

以偏移、十六进制与 ASCII 三列对照显示字节；两列共享同一选择与光标，只构建可见行。

## 最小示例

```go
data := []byte("MUJICA court ledger\x00\x01")
var state mujica.HexViewerState
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.HexViewer(c, data, &state, mujica.HexViewerOptions{}).Element.Height(160)
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func HexViewer(c *ui.Context, data []byte, s *HexViewerState, opts HexViewerOptions) ChoiceResult`。每行 16 字节；`s.Selection()` 返回选中范围；光标或锚点越界 panic。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `HexViewerState.Cursor`, `Anchor` | `int` | 0 | 光标字节与 Shift 选择的起点。 |
| `Disabled` | `bool` | false | 不能移动与选择。 |

## 状态

选中字节在两列都显示 Selection 背景，光标字节在两列都加 AccentText 框；不可打印字节显示为点；聚焦时外侧焦点环；空数据；禁用。

## 事件

`Changed()`：光标或选择变化。

## 键盘操作

- Tab 聚焦；←/→ 移动一字节，↑/↓ 移动一行，Home/End 到行首行尾；加 Shift 扩展选择。
- Cmd+C（macOS）/ Ctrl+C：以空格分隔的十六进制复制选中字节。
- 点击任一列选择字节，Shift+点击扩展。

## 限制

- 只读，不支持编辑与搜索。
