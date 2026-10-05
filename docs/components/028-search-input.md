# 028 SearchInput 搜索框

## 用途

搜索模式：放大镜、输入、清空与提交。查找栏模式（`Find`）：在此基础上显示匹配数量、上一项与下一项，对应 MyGo 的 Search field 与 Find bar。匹配由调用方计算。

## 最小示例

```go
matches := 0
if query != "" {
	matches = strings.Count(text, query)
}
r := mujica.SearchInput(c, &query, mujica.SearchInputOptions{Find: true, Matches: matches, Current: &current})
if r.Moved() {
	fmt.Println("show match", current)
}
```

签名：`func SearchInput(c *ui.Context, query *string, opts SearchInputOptions) SearchInputResult`

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `query` | `*string` | 必填 | 查询文字。 |
| `Placeholder` | `string` | “搜索”或“查找” | 占位文字。 |
| `Label` | `string` | 占位文字 | 无障碍名称。 |
| `Disabled` | `bool` | `false` | 禁用（查找快捷键也停用）。 |
| `Find` | `bool` | `false` | 开启查找栏模式。 |
| `Matches` | `int` | `0` | 匹配数量；负数 panic。 |
| `Current` | `*int` | `nil` | 当前匹配下标；`Find` 时必填，组件会把它夹在有效范围内。 |

## 状态

同 TextInput；查找栏模式下显示“第 n 项，共 m 项”或“无匹配”（状态角色，变化时可被朗读），无匹配时上一项/下一项禁用。

## 事件

- `Changed()`：查询改变（含清空）。
- `Submitted()`：搜索模式下 Enter。
- `Moved()`：查找栏模式下当前匹配改变。

## 键盘操作

搜索模式：Enter 提交，Escape 清空。查找栏模式：Enter 下一项，Shift+Enter 上一项，macOS 上 Cmd+G / Shift+Cmd+G、其他平台 F3 / Shift+F3（窗口级），循环到首尾；Escape 清空。

## 限制

- MyGo 单行编辑器把 Shift+Enter 当作普通提交且不公开修饰键（`editor.go:489`，`submitMods` 未导出），组件用 `HandleInput` 先截获 Shift+Enter。
- 没有“完成/关闭”按钮；是否显示查找栏由调用方决定。
