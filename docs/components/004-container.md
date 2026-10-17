# 004 Container 容器

## 用途

为一组内容提供内边距、最大宽度、水平居中、填充父级和裁剪，适合正文、设置页等需要可读行宽的区域。

## 最小示例

```go
layout.Container(c, layout.ContainerOptions{MaxWidth: 320, Center: true}, func() {
	ui.Text(c, "A centered column of prose.")
})
```

签名：`func Container(c *ui.Context, opts ContainerOptions, children func()) ui.Element`

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Padding` | `Space` | `SpaceL`（4 单位，默认密度 16 DIP） | 四周内边距，随密度缩放。 |
| `MaxWidth` | `float32` | 0（不限） | 最大宽度 DIP；负数 panic。 |
| `Center` | `bool` | false | 占满父级宽度后受 `MaxWidth` 限制并水平居中。 |
| `Fill` | `bool` | false | 占满父级内容区（宽和高）。 |
| `Clip` | `bool` | false | 裁剪超出的内容。 |
| `children` | `func()` | 必填 | 子内容，纵向排列；nil 时 panic。 |

`Space` 取值：`SpaceDefault`（组件默认）、`SpaceNone`、`SpaceXS`（1 单位）、`SpaceS`（2）、`SpaceM`（3）、`SpaceL`（4）、`SpaceXL`（6）、`SpaceXXL`（8）。单位为 Compact 3、Default 4、Comfortable 5 DIP。非法值 panic。

## 状态

布局组件，无交互状态。不绘制背景，可在返回的元素上调用 `Background`、`Border`。

## 事件

无。

## 键盘操作

不获取焦点。

## 限制

- 只做水平居中；垂直居中请用 `Stack` 的对齐或 MyGo `Center()`。
- 子内容固定纵向排列，需要横向时在内部放 `Stack`。
