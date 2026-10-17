# 317 DiffStat 差异统计

## 用途

显示 +新增 −删除 行数，以及按比例填充的方块条。

## 最小示例

```go
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	git.DiffStat(c, 120, 8, git.DiffStatOptions{})
}
ui.Render(view, 400, 300, 1)
```

签名：`func DiffStat(c *ui.Context, added, deleted int, opts DiffStatOptions) ui.Element`

## 参数

`DiffStatOptions`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| Blocks | int | 0（即 5） | 比例条方块数；为负 panic。 |

added、deleted 为负时 panic。

## 状态

- 数字带 + 与 − 符号，不只靠颜色；非零一侧至少占一格。
- 无变化时方块为 Border 色。
- 朗读为“新增 n 行，删除 m 行”。

## 事件

- 无。

## 键盘操作

- 无：纯展示，不获取焦点。

## 限制

- 无。
