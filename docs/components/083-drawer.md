# 083 Drawer 抽屉

## 用途

从窗口四边之一滑出的面板，用于筛选、详情和设置。可选模态或非模态。

## 最小示例

```go
open := true
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	if mujica.Button(c, "Filters", mujica.ButtonOptions{}).Clicked() {
		open = true
	}
	mujica.Drawer(c, &open, mujica.DrawerOptions{Title: "Filters", Side: mujica.DrawerRight, Size: 320}, func() {
		ui.Text(c, "Guests are placed by rank.")
	})
}
ui.Render(view, 400, 300, 1)
```

签名：`func Drawer(c *ui.Context, open *bool, opts DrawerOptions, content func()) *ui.Element`。

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Title` | `string` | 必填 | 标题与无障碍名称；为空时 panic。 |
| `Side` | `DrawerSide` | `DrawerRight` | `DrawerRight`、`DrawerLeft`、`DrawerTop`、`DrawerBottom`；越界时 panic。 |
| `Size` | `float32` | 左右 360，上下 280 | 宽度或高度；限制在 160 DIP 与窗口的 90% 之间；负数 panic。 |
| `NonModal` | `bool` | `false` | 非模态：不加遮罩、不限制焦点，窗口其余部分仍可操作。 |
| `Actions` | `func()` | `nil` | 页脚操作区。 |

## 状态

- 打开时从所在边滑入并淡入 160 ms；减少动画时直接出现。
- 面板贴边，只在朝内一侧画 1 DIP 分隔线。
- 内容超出时在内容区滚动。

## 事件

关闭按钮与 Escape 把 `*open` 设为 false；模态抽屉点击遮罩同样关闭。

## 键盘操作

| 按键 | 行为 |
| --- | --- |
| Tab / Shift+Tab | 模态：在抽屉内循环。非模态：按窗口焦点顺序，抽屉内容排在窗口内容之后。 |
| Escape | 关闭（先交给正在编辑的控件和更上层浮层）。 |

## 限制

- 非模态抽屉不在打开时移动焦点，关闭时也不自动恢复焦点：MyGo 记录与恢复打开者的 `openOverlay` 不公开（`ui/scope.go:123`），也没有读取当前焦点的公开接口。
- 非模态抽屉的 Escape 是窗口级快捷键，优先级最低；同时打开多个非模态浮层时由先构建的那个响应。
