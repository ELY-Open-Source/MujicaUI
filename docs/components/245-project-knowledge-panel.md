# 245 ProjectKnowledgePanel · 项目知识面板

## 用途

列出项目知识文件，提供添加、移除，并显示 token 用量与容量提示。选择与读取文件由调用方完成。

## 最小示例

```go
files := []mujica.KnowledgeFile{{ID: "1", Name: "etiquette.pdf", Tokens: 12000}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	v := mujica.ProjectKnowledgePanel(c, files, mujica.ProjectKnowledgePanelOptions{Capacity: 200000})
	if v.Add {
		log.Println("pick files")
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func ProjectKnowledgePanel(c *ui.Context, files []KnowledgeFile, opts ProjectKnowledgePanelOptions) ProjectKnowledgePanelView`。

`KnowledgeFile`：`ID`、`Name`、`Tokens`（不能为负）。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Capacity` | `int` | 必填 | 知识容量（tokens），必须为正，否则 panic。 |
| `Title` | `string` | 内置“项目知识” | 面板标题。 |
| `Disabled` | `bool` | false | 禁用添加与移除。 |

## 状态

正常用量（酒红进度）、≥ 80%（警告色 + 三角图标 + “即将用满”）、超出容量（危险色 + 圆形警示图标 + “超出容量”）、空列表、禁用。

## 事件

`Add`：请调用方打开文件选择。`Removed`：要移除的文件 ID。

## 键盘操作

- Tab：添加按钮，再进入文件列表；↑/↓ 在移除按钮间移动；Enter/Space 移除。

## 限制

- token 数由调用方计算；面板只做求和与比例显示。
