# 212 LinkPreviewCard 链接预览卡片

## 用途

显示链接的站点、标题、描述与缩略图；元数据由调用方获取。整张卡片是一个按钮，悬停显示完整地址。

## 最小示例

```go
p := chat.LinkPreview{URL: "https://example.com", Title: "The Palace", Site: "example.com"}
if chat.LinkPreviewCard(c, p).Clicked() {
	fmt.Println("open", p.URL)
}
```

签名：`func LinkPreviewCard(c *ui.Context, p LinkPreview) ui.Element`

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `URL` | `string` | 必填 | 链接地址，作为悬停提示。 |
| `Title` | `string` | 必填 | 标题，最多两行，AccentText 色。 |
| `Description` | `string` | "" | 描述，最多两行。 |
| `Site` | `string` | URL | 站点名；为空时显示 URL。 |
| `Thumbnail` | `ui.ImageSource` | nil | 缩略图，64×64 裁切显示。 |

## 状态

悬停时底色加深；聚焦时显示焦点环。

## 事件

`Clicked()`：点击或 Enter、Space；打开链接由调用方完成。

## 键盘操作

Tab 聚焦，Enter 或 Space 触发。

## 限制

- 组件不联网：标题、描述与缩略图都由调用方提供。
