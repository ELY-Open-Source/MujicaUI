# 333 SymbolOutline · 符号大纲

## 用途

以树显示文件中的符号，带类型标记，可筛选，并跟随编辑器光标选中所在符号。

## 最小示例

```go
symbols := []code.Symbol{{Name: "Court", Kind: code.SymbolType, Line: 3, EndLine: 40, Children: []code.Symbol{
	{Name: "Seal", Kind: code.SymbolMethod, Line: 10, EndLine: 20},
}}}
var state code.SymbolOutlineState
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := code.SymbolOutline(c, symbols, &state, code.SymbolOutlineOptions{CursorLine: 12})
	r.Element.Height(240)
	if r.Jumped != nil {
		fmt.Println("go to line", r.Jumped.Line)
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func SymbolOutline(c *ui.Context, symbols []Symbol, s *SymbolOutlineState, opts SymbolOutlineOptions) SymbolOutlineResult`。符号树在数据变化时编号一次。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Symbol.Name`, `Detail` | `string` | — | 名称与补充说明（如签名）。 |
| `Symbol.Kind` | `SymbolKind` | `SymbolFunction` | 函数 f、方法 m、类型 T、接口 I、字段 p、变量 v、常量 c、模块 M；越界 panic。 |
| `Symbol.Line`, `EndLine` | `int` | — | 所在行范围，从 1 开始。 |
| `SymbolOutlineState.Query` | `string` | "" | 筛选文字；匹配项与其祖先保留并自动展开。 |
| `CursorLine` | `int` | 0（不跟随） | 变化时选中包含该行的最内层符号并展开祖先。 |
| `Disabled` | `bool` | false | 不能跳转。 |

## 状态

类型以字母徽记加颜色区分；选中行 Selection 与菱形；筛选高亮；无符号；禁用。

## 事件

`Jumped`：用户要跳转的符号。

## 键盘操作

- Tab 到达筛选框与树；↑/↓ 移动，→ 展开或进入子项，← 折叠或回到父项；输入字母按名称定位。
- Enter 或双击：跳转到符号。

## 限制

- 光标跟随只在 `CursorLine` 变化时生效，之后用户可自由移动选择。
