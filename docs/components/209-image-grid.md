# 209 ImageGrid 图片网格

## 用途

在消息中排列 1 至 9 张图片：1 张大图，2 与 4 张两列，其余三列。点击已加载的图片打开 ImageViewer；加载中与失败状态复用 Image 组件，失败时可重试。

## 最小示例

```go
if i, ok := chat.ImageGrid(c, []chat.GridImage{{Source: pic, Alt: "Portrait"}}).Retried(); ok {
	fmt.Println("reload", i)
}
```

签名：`func ImageGrid(c *ui.Context, images []GridImage) ImageGridResult`

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `images` | `[]GridImage` | 必填 | 1 至 9 张，否则 panic。 |
| `Source` | `ui.ImageSource` | nil | 已加载的位图或 SVG。 |
| `Alt` | `string` | "" | 替代说明，也是按钮名称与查看器标题。 |
| `Loading / Err` | `bool / error` | false / nil | 调用方的加载中与失败状态。 |

## 状态

已加载、加载中（菱形加载标记）、失败（重试按钮）；只有已加载的图片可点击。

## 事件

`Retried()` 返回本帧点击了重试的图片下标。查看器在组件内部打开与关闭。

## 键盘操作

Tab 聚焦图片，Enter 或 Space 打开查看器；查看器内的键盘操作见 ImageViewer。

## 限制

- 查看器只包含已加载的图片，序号按已加载图片计算。
