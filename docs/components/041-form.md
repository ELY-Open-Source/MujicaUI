# 041 Form

表单。

## 用途

组合多个字段：注册字段、同步与调用方提供的异步校验、提交与重置、提交中防重入、提交失败时定位第一个错误字段。Form 不保存业务模型，不调用远程接口，也不做持久化；值始终由应用持有。

## 最小示例

```go
var state input.FormState
name := ""
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	res := input.Form(c, &state, input.FormOptions{}, func(f *input.FormScope) {
		input.FormItem(c, f, "name", &name, input.FormItemOptions[string]{Label: "Name", Required: true}, func() ui.Element {
			return input.InputGroup(c, &name, input.InputGroupOptions{Invalid: state.Error("name") != ""}).Input
		})
	})
	if res.Submitted() {
		state.EndSubmit("") // after saving name
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `FormOptions.Horizontal` | `bool` | false | 标签放在控件左侧一列，右对齐。 |
| `FormOptions.SubmitLabel` | `string` | 内置“Submit/提交” | 提交按钮文字。 |
| `FormOptions.ResetLabel` | `string` | 内置“Reset/重置” | 重置按钮文字。 |
| `FormItemOptions[T].Label` | `string` | 必填 | 字段标签；空值 panic。 |
| `FormItemOptions[T].Required` | `bool` | false | 标记必填；零值、空字符串、空切片或空映射视为缺失。 |
| `FormItemOptions[T].Description` | `string` | "" | 说明文字；异步校验进行中显示“Checking…”。 |
| `FormItemOptions[T].Initial` | `T` | 零值 | 重置时写回的初始值。 |
| `FormItemOptions[T].Validate` | `func(T) string` | nil | 同步校验，返回错误文字或 ""；跨字段规则在此读取其他值。 |
| `FormItemOptions[T].Async` | `func(v T, version uint64)` | nil | 启动异步校验；完成后在界面线程调用 `FormState.Resolve(name, version, msg)`。 |

## 状态

默认；字段错误（触碰后或首次提交后才显示）；异步校验中（说明显示 Checking…）；提交中（提交按钮显示加载并阻止重复提交，重置禁用）；提交失败（表单级错误显示在按钮上方，保留全部值）。

## 事件

`FormParts.Submitted()`：全部校验通过的一次提交，处理完毕后调用 `FormState.EndSubmit(err)`。`FormParts.Reset()`：已恢复初始值并清除校验与提交状态。`FormState.Submit()` 可在 Enter 等处请求提交。`FormState.Error(name)` 返回当前应显示的错误，用于给控件设置 Invalid。

## 键盘操作

Tab / Shift+Tab 在字段间移动，离开字段即校验。提交失败后焦点移到第一个错误字段并滚动到可见处，同时向屏幕阅读器播报。点击标签聚焦控件。

## 限制

- 控件由调用方构建，`FormItem` 的 control 回调必须返回可获取焦点的元素，用于错误定位。
- 校验在失焦与提交时触发；已触碰字段在值变化时重新同步校验，异步校验只在下一次失焦或提交时重跑。
- 异步结果以版本号对齐：值变化后旧版本的 `Resolve` 返回 false 并被丢弃。`FormState` 的方法只能在界面线程调用，后台任务需经 `Window.Update` 回到界面线程。
- 必填标记“*”只作视觉显示，辅助技术通过说明“必填”获知。
