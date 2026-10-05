# 186 SymbolBadge

## 用途

代码徽标。

## 最小示例

```go
symbol := mujica.MarketSymbol{Code: "AAPL", Name: "Apple", Market: "NASDAQ", Type: "Stock"}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	r := mujica.SymbolBadge(c, symbol, mujica.SymbolBadgeOptions{Clickable: true})
	if r.Changed() {
		fmt.Println(symbol.Code)
	}
}
ui.Render(view, 960, 440, 1)
```

## 参数

`MarketSymbol.Code` 必填，Name、Market、Type 为调用方标签。`Clickable/Disabled` 默认 false。

## 状态

默认；交互区域支持聚焦与选中，Disabled 参数控制禁用。

## 事件

if r.Changed() { 调用方回调(symbol.Code) }

## 键盘操作

Clickable 开启时 Tab 聚焦，Enter/Space 激活。

## 限制

市场和类型属于业务标签，由调用方本地化。
