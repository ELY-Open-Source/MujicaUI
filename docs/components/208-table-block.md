# 208 TableBlock 表格块

## 用途

在消息中显示小型表格：表头加粗，列按内容宽度排列并铺满消息宽度，过宽时横向滚动；可复制为 TSV。

## 最小示例

```go
if chat.TableBlock(c, []string{"Guest", "Seat"}, [][]string{{"Marie", "1"}}, chat.TableBlockOptions{}).Copied() {
	fmt.Println("copied as TSV")
}
```

签名：`func TableBlock(c *ui.Context, header []string, rows [][]string, opts TableBlockOptions) CopyResult`

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `header` | `[]string` | 必填 | 表头；为空 panic。 |
| `rows` | `[][]string` | nil | 数据行；每行列数必须与表头相同，否则 panic。 |
| `Align` | `[]ui.Align` | 全部 `ui.Start` | 各列对齐：`ui.Start`、`ui.Center`、`ui.End`；多于列数 panic。 |

## 状态

复制后按钮 2 秒内显示对勾。

## 事件

`Copied()`：本帧写入了 TSV（制表符分列、换行分行；单元格内的制表符与换行替换为空格）。

## 键盘操作

Tab 到复制按钮，Enter 或 Space 复制。

## 限制

- 不排序、不选择；需要这些功能请使用 DataTable（069 起）。
