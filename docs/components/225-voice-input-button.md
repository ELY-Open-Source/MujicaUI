# 225 VoiceInputButton 语音输入按钮

## 用途

录音开始与停止；录音时长；权限失败提示；录音由调用方实现。

## 最小示例

```go
var rec mujica.VoiceInputOptions
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	if mujica.VoiceInputButton(c, rec).Changed() {
		rec = mujica.VoiceInputOptions{Recording: !rec.Recording, Started: time.Now()}
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

`Recording=false`；录音时 `Started` 必填；`Denied=false`，表示权限失败；`Disabled=false`。

## 状态

默认、悬停、键盘焦点与祖先禁用；业务状态由调用方提供。

## 事件

ChoiceResult.Changed() 表示开始、停止或重试录音请求。

## 键盘操作

Tab 到按钮；Enter、Space 开始或停止。

## 限制

录音与权限请求由调用方实现。
