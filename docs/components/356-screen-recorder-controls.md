# 356 ScreenRecorderControls 录屏控制

## 用途

开始、暂停、继续、停止录制并显示计时；录制本身由调用方实现。

## 最小示例

```go
st := media.RecorderState{}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	media.ScreenRecorderControls(c, st, media.ScreenRecorderControlsOptions{
		OnStart:  func() { st.Status = media.RecorderRecording },
		OnPause:  func() { st.Status = media.RecorderPaused },
		OnResume: func() { st.Status = media.RecorderRecording },
		OnStop:   func() { st = media.RecorderState{} },
	})
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func ScreenRecorderControls(c *ui.Context, st RecorderState, opts ScreenRecorderControlsOptions) ui.Element`

`RecorderState`：`Status`（`RecorderIdle`、`RecorderRecording`、`RecorderPaused`）、`Elapsed`。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `OnStart / OnPause / OnResume / OnStop` | `func()` | 必填 | 缺少任一回调时 panic。 |
| `Disabled` | `bool` | false | 禁用。 |

## 状态

就绪、录制中（红点）、已暂停；状态同时用文字表示，不只靠圆点颜色。录制中的圆点不闪烁。

## 事件

全部通过回调报告。

## 键盘操作

| 按键 | 作用 |
| --- | --- |
| Tab | 在按钮间移动 |
| Enter / Space | 按下 |

## 限制

- 计时由调用方更新 `Elapsed`。
