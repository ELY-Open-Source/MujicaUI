# 462 SettingsSearch

设置搜索。

## 用途

按标签、说明与关键词搜索设置项，高亮匹配片段；选中后切到对应分区并把该设置行滚入视野、聚焦其控件。

## 最小示例

```go
query, page := "", "general"
entries := []account.SettingEntry{{ID: "sound", Page: "general", Label: "Sound", Keywords: []string{"audio"}}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if e, ok := account.SettingsSearch(c, &query, entries, account.SettingsSearchOptions{}).Chosen(); ok {
		page = e.Page
	}
	fmt.Println("page", page)
}
ui.Render(view, 400, 300, 1)
```

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `query` | `*string` | 必填 | 搜索文字。 |
| `entries` | `[]SettingEntry` | 必填 | ID（对应 SettingsRow ID）、Page、Label、Description、Keywords。 |
| `Limit` | `int` | 8 | 最多显示条数。 |

## 状态

无输入、有结果（当前项底色）、无结果。

## 事件

`SettingsSearchResult.Chosen()` 返回选中的条目；调用方据此设置分区。

## 键盘操作

在搜索框中 ↑/↓ 移动当前项，Enter 跳转；也可点击结果。

## 限制

大集合使用虚拟列表，构建当前可见行；搜索与分组在数据或查询变化时更新。

- 匹配不区分大小写，按子串查找。
