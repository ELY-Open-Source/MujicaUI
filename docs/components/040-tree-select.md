# 040 TreeSelect · 树选择

## 用途

在树形数据中选择一个或多个节点，支持搜索、展开折叠、半选显示，父子联动规则可配置。

## 最小示例

```go
var roles []string
nodes := []input.TreeSelectNode{
	{ID: "court", Label: "Court", Children: []input.TreeSelectNode{
		{ID: "king", Label: "King"},
		{ID: "queen", Label: "Queen"},
	}},
}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	input.TreeSelect(c, &roles, nodes, input.TreeSelectOptions{Multiple: true, Label: "Roles"})
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func TreeSelect(c *ui.Context, selected *[]string, nodes []TreeSelectNode, opts TreeSelectOptions) ChoiceResult`。单选模式下 `*selected` 最多一个 ID，否则 panic。节点 ID 在整棵树中必须唯一。

`TreeSelectNode`：`ID`、`Label`、`Children []TreeSelectNode`、`Disabled`。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Multiple` | `bool` | false | 多选（复选框）；单选时点击节点即选择并关闭。 |
| `Rule` | `TreeSelectRule` | `TreeSelectCascade` | 多选父子规则：`TreeSelectCascade` 勾选父节点即勾选其全部可用叶节点，值只含叶节点 ID，父节点按子节点显示全选/半选；`TreeSelectIndependent` 每个节点独立勾选，值可含任意 ID。 |
| `Placeholder` | `string` | 内置“请选择…” | 无值时提示。 |
| `MaxTags` | `int` | 0 → 3 | 多选时字段中显示的标签数，其余折叠为“+N”。 |
| `Disabled` | `bool` | false | 禁用。 |
| `Label` | `string` | "" | 无障碍名称。 |

## 状态

字段：占位、单值文字、多值标签与折叠、展开、聚焦、禁用。弹层：搜索框、展开/折叠箭头、未选/选中/半选方框（单选为玫瑰底高亮）、禁用节点（联动时被跳过）、无结果。搜索时显示匹配节点及其祖先并自动展开。

## 事件

`Changed()` 在选择变化的那一帧返回 true。

## 键盘操作

- Enter / Space / ↓：打开弹层，焦点在搜索框；搜索框中 ↓ 进入树。
- ↑/↓、Home、End：在可见节点间移动。
- →：展开；已展开时进入第一个子节点。←：折叠；已折叠或叶节点时回到父节点。
- Space：勾选（多选）；Enter / Space：选择（单选）。
- Escape：关闭，焦点回到字段。

## 限制

- 树行基于 `ui.CheckboxBase` / `ui.ButtonBase` 自绘。MyGo 的 `TreeItem` 只能显示文字标签，无法放入复选框（ui/tree.go:44），故未复用。
- MyGo 不对自定义元素公开 expanded 与 mixed 状态的设置（ui/tree.go:98、ui/element.go:282 仅内部设置）。组件改用 `Description` 读出“已展开/已收起”“部分选中”，展开箭头带“展开/收起 + 节点名”的名称；屏幕阅读器读不到真正的 expanded/mixed 状态属性。单选行通过 `Selected` 报告选中，并显示勾号。
- `CollapsibleBase` 的触发器能报告 expanded，但它自身要占一个焦点位，会让方向键在箭头与节点间交替停留，故未采用。
- 列表未虚拟化，适合数千节点以内。
