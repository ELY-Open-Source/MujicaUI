# MujicaUI

宫廷哥特风格的 Go 桌面组件库，建立在 [MyGo](https://mygo.egoist.dev/docs/ui) 原生 UI 之上。象牙白与炭黑底色，酒红主色，旧金装饰线，衬线标题。100 个组件，浅深两套主题，三种密度。

界面由 MyGo 在 GPU 上自绘，不使用 HTML、JavaScript 或 WebView。

## 安装

```sh
go get github.com/ELY-Open-Source/MujicaUI
```

要求 Go 1.27.1 与 MyGo v0.2.10（见 [兼容表](docs/compat.md)）。

## 最小示例

即 `examples/minimal/main.go`，CI 校验两者一致。

```go
package main

import (
	"log"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/input"
)

type app struct{ name string }

func (a *app) view(c *ui.Context) {
	core.Use(c, core.Settings{Locale: core.ZhCN})
	ui.Column(c).Fill().Center().Gap(12).Children(func() {
		input.TextInput(c, &a.name, input.TextInputOptions{Placeholder: "姓名"})
		if input.Button(c, "觐见", input.ButtonOptions{}).Clicked() {
			log.Println("hello", a.name)
		}
	})
}

func main() {
	a := &app{}
	mygo.App.WhenReady(func() {
		mygo.NewWindow(mygo.WindowOptions{Title: "Mujica", Width: 480, Height: 320, Content: ui.View(a.view)})
	})
	if err := mygo.App.Run(); err != nil {
		log.Fatal(err)
	}
}
```

## 主题

每帧开头调用 `core.Use(c, settings)`。`Settings` 决定模式（`System`、`Light`、`Dark`）、密度（`theme.Compact`、`theme.Default`、`theme.Comfortable`）、语言与文案覆盖。每个窗口各自设置，切换不重置业务状态。

```go
core.Use(c, core.Settings{
	Mode:    core.Dark,
	Density: theme.Compact,
	Locale:  core.ZhCN,
})
k := core.Tokens(c) // 当前窗口的颜色
```

覆盖顺序：组件参数，窗口 `Settings`，库默认值。系统高对比度、文本缩放、减少动画始终生效。

## Gallery

```sh
go run ./examples/gallery            # 原生窗口
go run ./examples/gallery -web 127.0.0.1:8080 # 浏览器镜像，仅本机
```

`-web` 按传入的地址监听；写 `:8080` 会对所有网络接口开放。浏览器镜像不在浏览器里运行组件。Go 进程用 `ui.NewTester` 驱动同一个 Gallery 视图，把帧推给页面，再把指针、键盘、输入法和剪贴板事件送回 Go。

## 文档

- [组件文档](docs/components/)：每个组件一页。
- [兼容与平台依赖](docs/compat.md)
- [上游归并](docs/upstream.md)

## 许可

MIT。第三方声明见 [THIRD_PARTY_NOTICES.md](THIRD_PARTY_NOTICES.md)。
