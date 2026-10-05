# 181 OptionChain

## 用途

期权链。

## 最小示例

```go
var state finance.OptionChainState
rows := []finance.OptionStrike{{Strike: 100, Call: finance.OptionQuote{Bid: 3, Ask: 3.2, Last: 3.1}, Put: finance.OptionQuote{Bid: 2, Ask: 2.2, Last: 2.1}}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := finance.OptionChain(c, &state, []string{"Oct 17", "Nov 21"}, rows, finance.OptionChainOptions{Spot: 100.4})
	if r.Changed || r.Submitted || r.ExpiryChanged {
		fmt.Println(state.Expiry, state.Strike, state.Side)
	}
}
ui.Render(view, 960, 440, 1)
```

## 参数

`s` 保存到期日和合约；`expiries` 至少一个；`strikes` 升序。`Spot` 默认 0 隐藏现价；`Decimals` 默认 2；`Disabled` 默认 false。

## 状态

默认；交互区域支持聚焦与选中，Disabled 参数控制禁用。

## 事件

if r.Changed || r.Submitted || r.ExpiryChanged { 调用方回调(state.Expiry, state.Strike, state.Side) }

## 键盘操作

Tab 切换到期日或列表；上下、Home/End 选择行；左右切换认购认沽；Enter 确认。

## 限制

到期日切换后，由调用方更换对应报价；排序与稳定合约标识由调用方保证。
