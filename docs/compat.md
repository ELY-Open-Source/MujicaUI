# 兼容与平台依赖

## 版本

| 项 | 值 |
| --- | --- |
| Go | 1.27.1 及以上（MyGo v0.5.5 的 go.mod 要求） |
| MyGo | v0.5.5，commit 4d8673c，固定在 go.mod 与 go.sum |
| 运行时依赖 | macOS：系统框架。Windows：无需 WebView2。Linux：GTK 3（MyGo 通过 GtkGLArea 绘制） |

## 已实测环境

| 系统 | 架构 | 会话 | 结果 |
| --- | --- | --- | --- |
| macOS（Darwin 27.2） | arm64 | Aqua | 无窗口测试、快照、Gallery 构建 |

以下项目尚未实测，不宣称支持：Windows（x64、arm64）、Linux X11、原生 Wayland、XWayland；VoiceOver、Narrator、NVDA、Orca；跨缩放显示器移动窗口；100% 以外的真机显示缩放。无窗口测试已覆盖 1x 与 2x 渲染，不替代真机。

## 快照基准

快照按系统分目录：`testdata/snapshots/<GOOS>/`。目前只有 `darwin` 基准。其他平台须在该平台运行 `go test -update .` 生成并人工检查后提交；在此之前用 `go test -skip Snapshot ./...` 运行其余测试。

## 离线

库与 Gallery 不发起网络请求，无遥测、账号或远程字体。图标与角饰随代码嵌入。浏览器镜像按 `-web` 传入的地址监听：`127.0.0.1:8080` 仅本机可达，`:8080` 对所有网络接口开放。
