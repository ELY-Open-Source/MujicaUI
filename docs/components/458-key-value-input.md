# 458 KeyValueInput · 键值对输入

## 用途

编辑键值对、切换启用状态并增删行。

## 最小示例

```go
rows := []mujica.KeyValue{{Key: "Accept", Value: "application/json", Enabled: true}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	if mujica.KeyValueInput(c, &rows, mujica.KeyValueInputOptions{}).Changed() {
		fmt.Println("rows", len(rows))
	}
}
ui.Render(view, 720, 480, 1)
```

## 参数

KeyPlaceholder、ValuePlaceholder 默认内置文案；Label 为控件名称；MaxHeight 默认 320 DIP；Disabled 禁用；rows 必填。

## 状态

空集合、启用行、停用行、禁用。

## 事件

Changed() 报告编辑、启用切换或增删。

## 键盘操作

Tab 到各字段及按钮；Space 切换启用；Enter 激活增删。

## 限制

行顺序保持输入顺序；大量行采用虚拟化。
