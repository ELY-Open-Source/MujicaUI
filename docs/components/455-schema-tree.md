# 455 SchemaTree · 数据库结构树

## 用途

显示数据库、表、视图、列和索引树。

## 最小示例

```go
roots := []devtools.SchemaNode{{Name: "app", Children: []devtools.SchemaNode{{Name: "users", Kind: devtools.SchemaTable, Children: []devtools.SchemaNode{{Name: "id", Kind: devtools.SchemaColumn, Detail: "bigint"}, {Name: "users_pk", Kind: devtools.SchemaIndex, Detail: "(id)"}}}}}}
var state devtools.SchemaTreeState
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if devtools.SchemaTree(c, &state, roots, devtools.SchemaTreeOptions{}).Changed() {
		path, _ := state.Selected()
		fmt.Println("selected", path)
	}
}
ui.Render(view, 720, 480, 1)
```

## 参数

state 必填；Height 默认 360 DIP；Name 必填，兄弟名称唯一；Kind 使用枚举，Detail 为类型说明。

## 状态

展开、折叠、选择、空结构、父级禁用。

## 事件

Changed() 报告选择变化；state.Selected() 返回路径。

## 键盘操作

Tab 进入树；上下选择；左右展开折叠。

## 限制

路径以斜线分隔；节点名称须保持路径唯一；更新结构提交新切片。
