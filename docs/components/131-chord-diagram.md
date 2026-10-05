# 131 ChordDiagram 弦图

## 用途

以组弧和弦带展示关系矩阵，悬停或键盘读数组时高亮相关连线。

## 最小示例

```go
names := []string{"Nave", "Choir", "Crypt"}
matrix := [][]float64{{0, 5, 2}, {4, 0, 1}, {1, 2, 0}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.ChordDiagram(c, names, matrix, mujica.ChordDiagramOptions{Frame: mujica.ChartFrame{Label: "Paths"}})
}
ui.Render(view, 400, 300, 1)
```

## 参数

`names []string` 为唯一非空组名；`matrix [][]float64` 为同阶方阵，数值须有限且非负。`Frame` 控制名称、尺寸和加载状态。`Format` 自定义关系数值。

## 状态

默认、组高亮、弦高亮、加载、错误、空矩阵、祖先禁用；对角线值形成自身关系。

## 事件

`Point` 先枚举组，再枚举关系；`Changed()` 表示读数变化；`Retried()` 表示重试。

## 键盘操作

Tab 聚焦；方向键依次读取组的流入流出与关系值；Home/End 首尾；Escape 清除。

## 限制

高密度矩阵按绘图区半径抽样绘制弦带，当前键盘关系保留；布局仅在外层 matrix 切片身份变化时计算。更新数据需替换外层切片。组弧标签在有足够空间时展示，完整名称始终包含在键盘读数中。
