# 006 Grid 网格

## 用途

按行列排布子元素：列宽可为固定尺寸、按内容或按比例分配剩余空间；子项可指定位置与跨行跨列。基于 MyGo 网格布局。

## 最小示例

```go
layout.Grid(c, layout.GridOptions{Columns: []ui.Track{ui.Fixed(96), ui.Fr(1)}}, func() {
	layout.GridCell(c, layout.GridCellOptions{ColumnSpan: -1}, func() { ui.Text(c, "Header") })
	ui.Text(c, "Name")
	ui.Text(c, "Ada")
})
```

签名：

- `func Grid(c *ui.Context, opts GridOptions, children func()) ui.Element`
- `func GridCell(c *ui.Context, opts GridCellOptions, build func()) ui.Element`

## 参数

`GridOptions`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Columns` | `[]ui.Track` | 一列 `Fr(1)` | 各列尺寸：`ui.Fixed(dip)`、`ui.Fr(n)`（分配剩余空间）、`ui.FitContent()`。 |
| `Rows` | `[]ui.Track` | nil | 前几行尺寸；之后自动添加的行按内容高度。 |
| `Gap` | `Space` | `SpaceM`（3 单位） | 列间距，也是默认行间距。 |
| `RowGap` | `Space` | 同 `Gap` | 行间距。 |

`GridCellOptions`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Column`, `Row` | `int` | 0（自动放置） | 起始列、行，从 1 计；负数 panic。 |
| `ColumnSpan`, `RowSpan` | `int` | 1 | 跨越轨道数；-1 跨到最后一条；小于 -1 panic。 |

`GridCell` 内的单个子元素会拉伸铺满整个单元格，与网格直接子项一致。

## 状态

布局组件，无交互状态。

## 事件

无。

## 键盘操作

不获取焦点；Tab 顺序为构建顺序。

## 限制

- 网格会构建全部子项；大量条目请使用虚拟化列表。
- 不校验跨度是否超出列数，超出时由 MyGo 按其规则扩展。
