package main

import (
	"fmt"
	"math"
	"strings"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/icons"
	"github.com/ZacharyZhang-NY/MujicaUI/input"
)

const inputCategory = "Inputs"

// inputVariants lists the button variants for variant pickers.
var inputVariants = []input.RadioItem[input.ButtonVariant]{
	{Value: input.Primary, Label: "Primary"}, {Value: input.Secondary, Label: "Secondary"},
	{Value: input.Outline, Label: "Outline"}, {Value: input.Ghost, Label: "Ghost"}, {Value: input.Danger, Label: "Danger"},
}

// inputSection titles one part of a demo.
func inputSection(c *ui.Context, title string, fn func()) {
	ui.Column(c).Gap(8).Children(func() {
		ui.Text(c, title).FontWeight(600).TextColor(core.Tokens(c).TextMuted)
		fn()
	})
}

// inputFlags builds one checkbox per named flag, logging each toggle.
func inputFlags(c *ui.Context, log func(string), flags map[string]*bool, order ...string) {
	ui.Row(c).Gap(16).Wrap().Children(func() {
		for _, name := range order {
			if input.Checkbox(c, flags[name], name, input.CheckboxOptions{}).Changed() {
				log(fmt.Sprintf("option %s = %v", name, *flags[name]))
			}
		}
	})
}

type inputButtonDemo struct {
	variant                 input.ButtonVariant
	icon, loading, disabled bool
}

type inputTextDemo struct {
	value, err                              string
	clearable, readOnly, disabled, sl, fail bool
}

type inputAreaDemo struct {
	value              string
	max                float64
	readOnly, disabled bool
}

type inputNumberDemo struct {
	value, step                    float64
	empty, stepper, disabled, fail bool
}

type inputSearchDemo struct {
	query, find string
	current     int
	disabled    bool
}

type inputCheckDemo struct {
	single, mixed, disabled bool
	items                   [4]bool
}

type inputGroupDemo struct {
	vertical, disabled bool
}

type inputRadioDemo struct {
	value                   string
	horizontal, disabledAll bool
}

var inputPages = []page{
	{ID: "021", Name: "Button", ZhName: "按钮", Category: inputCategory, Keys: "Tab 聚焦；Enter 或 Space 触发。",
		Code: `if input.Button(c, "Save", input.ButtonOptions{Icon: icons.Must("check")}).Clicked() {
	fmt.Println("saved")
}`,
		Demo: func(c *ui.Context, log func(string)) {
			s := ui.Local(c.Root(), "021", func() inputButtonDemo { return inputButtonDemo{} })
			inputSection(c, "Parameters", func() {
				input.RadioGroup(c, &s.variant, inputVariants, input.RadioGroupOptions{Horizontal: true, Label: "Variant"})
				inputFlags(c, log, map[string]*bool{"Icon": &s.icon, "Loading": &s.loading, "Disabled": &s.disabled}, "Icon", "Loading", "Disabled")
			})
			inputSection(c, "Live", func() {
				opts := input.ButtonOptions{Variant: s.variant, Loading: s.loading, Disabled: s.disabled}
				if s.icon {
					opts.Icon = icons.Must("check")
				}
				ui.Row(c).Children(func() {
					if input.Button(c, "Save changes", opts).Clicked() {
						log("Button clicked")
					}
				})
			})
			inputSection(c, "Variants", func() {
				ui.Row(c).Gap(12).Wrap().Children(func() {
					for _, v := range inputVariants {
						if input.Button(c, v.Label, input.ButtonOptions{Variant: v.Value}).Clicked() {
							log(v.Label + " clicked")
						}
					}
					input.Button(c, "Loading", input.ButtonOptions{Loading: true})
					input.Button(c, "Disabled", input.ButtonOptions{Disabled: true})
				})
			})
		}},
	{ID: "022", Name: "IconButton", ZhName: "图标按钮", Category: inputCategory, Keys: "Tab 聚焦；Enter 或 Space 触发；悬停显示提示。",
		Code: `if input.IconButton(c, icons.Must("settings"), "Settings", input.IconButtonOptions{Variant: input.Ghost}).Clicked() {
	fmt.Println("open settings")
}`,
		Demo: func(c *ui.Context, log func(string)) {
			s := ui.Local(c.Root(), "022", func() inputButtonDemo { return inputButtonDemo{variant: input.Ghost} })
			inputSection(c, "Parameters", func() {
				input.RadioGroup(c, &s.variant, inputVariants, input.RadioGroupOptions{Horizontal: true, Label: "Variant"})
				inputFlags(c, log, map[string]*bool{"Loading": &s.loading, "Disabled": &s.disabled}, "Loading", "Disabled")
			})
			inputSection(c, "Live", func() {
				ui.Row(c).Gap(8).Children(func() {
					for _, n := range []string{"bold", "italic", "underline", "copy", "trash-2"} {
						opts := input.IconButtonOptions{Variant: s.variant, Loading: s.loading, Disabled: s.disabled}
						if input.IconButton(c, icons.Must(n), n, opts).Clicked() {
							log("IconButton " + n + " clicked")
						}
					}
				})
			})
		}},
	{ID: "023", Name: "ButtonGroup", ZhName: "按钮组", Category: inputCategory, Keys: "Tab 进入组；方向键在组内移动并跳过禁用项；Enter 或 Space 触发。",
		Code: `items := []input.ButtonGroupItem{{Label: "Day"}, {Label: "Week"}, {Label: "Month"}}
if i := input.ButtonGroup(c, items, input.ButtonGroupOptions{Label: "Range"}).ClickedItem(); i >= 0 {
	fmt.Println("chose", items[i].Label)
}`,
		Demo: func(c *ui.Context, log func(string)) {
			s := ui.Local(c.Root(), "023", func() inputGroupDemo { return inputGroupDemo{} })
			inputSection(c, "Parameters", func() {
				inputFlags(c, log, map[string]*bool{"Vertical": &s.vertical, "Disable group": &s.disabled}, "Vertical", "Disable group")
			})
			inputSection(c, "Live", func() {
				items := []input.ButtonGroupItem{{Label: "Day"}, {Label: "Week", Disabled: true}, {Label: "Month"}, {Icon: icons.Must("calendar"), Name: "Pick a date"}}
				r := input.ButtonGroup(c, items, input.ButtonGroupOptions{Vertical: s.vertical, Disabled: s.disabled, Label: "Range"})
				if i := r.ClickedItem(); i >= 0 {
					log(fmt.Sprintf("ButtonGroup item %d clicked", i))
				}
			})
		}},
	{ID: "024", Name: "TextInput", ZhName: "文本输入框", Category: inputCategory, Keys: "Tab 聚焦；Enter 提交；Escape 清空（可清空时）；Cmd/Ctrl+Z 撤销、+C/X/V 剪贴板；支持输入法组合。",
		Code: `r := input.TextInput(c, &name, input.TextInputOptions{Label: "Name", Placeholder: "Your name", Clearable: true})
if r.Submitted() {
	fmt.Println("hello", name)
}`,
		Demo: func(c *ui.Context, log func(string)) {
			s := ui.Local(c.Root(), "024", func() inputTextDemo { return inputTextDemo{clearable: true, sl: true} })
			inputSection(c, "Parameters", func() {
				inputFlags(c, log, map[string]*bool{"Clearable": &s.clearable, "Slots": &s.sl, "ReadOnly": &s.readOnly, "Disabled": &s.disabled}, "Clearable", "Slots", "ReadOnly", "Disabled")
				input.TextInput(c, &s.err, input.TextInputOptions{Label: "Error message", Placeholder: "Error message (empty for none)"})
			})
			inputSection(c, "Live", func() {
				opts := input.TextInputOptions{Label: "Name", Placeholder: "Your name", Clearable: s.clearable, ReadOnly: s.readOnly, Disabled: s.disabled, Error: s.err}
				if s.sl {
					opts.Leading = func() { ui.Icon(c, icons.Must("user")).Size(16, 16).TextColor(core.Tokens(c).TextMuted) }
					opts.Trailing = func() { ui.Text(c, fmt.Sprintf("%d", len([]rune(s.value)))).TextColor(core.Tokens(c).TextMuted) }
				}
				r := input.TextInput(c, &s.value, opts)
				if r.Changed() {
					log(fmt.Sprintf("TextInput changed: %q", s.value))
				}
				if r.Submitted() {
					log(fmt.Sprintf("TextInput submitted: %q", s.value))
				}
			})
			inputSection(c, "States", func() {
				ui.Row(c).Gap(12).Children(func() {
					e, ro, off := "1789-13-01", "Read only", "Disabled"
					input.TextInput(c, &e, input.TextInputOptions{Label: "Error state", Error: "No such date"}).Element.Grow(1)
					input.TextInput(c, &ro, input.TextInputOptions{Label: "Read-only state", ReadOnly: true}).Element.Grow(1)
					input.TextInput(c, &off, input.TextInputOptions{Label: "Disabled state", Disabled: true}).Element.Grow(1)
				})
			})
		}},
	{ID: "025", Name: "TextArea", ZhName: "多行文本框", Category: inputCategory, Keys: "Enter 换行；Cmd/Ctrl+Enter 提交（可配置）；滚轮与方向键滚动；撤销、剪贴板与输入法同 TextInput。",
		Code: `r := input.TextArea(c, &notes, input.TextAreaOptions{Label: "Notes", MaxHeight: 160, MaxLength: 200})
if r.Submitted() {
	fmt.Println("send", notes)
}`,
		Demo: func(c *ui.Context, log func(string)) {
			s := ui.Local(c.Root(), "025", func() inputAreaDemo { return inputAreaDemo{max: 120} })
			inputSection(c, "Parameters", func() {
				inputFlags(c, log, map[string]*bool{"ReadOnly": &s.readOnly, "Disabled": &s.disabled}, "ReadOnly", "Disabled")
				if input.NumberInput(c, &s.max, input.NumberInputOptions{Min: 0, Max: 1000, Step: 10, Label: "MaxLength"}).Changed() {
					log(fmt.Sprintf("option MaxLength = %v", s.max))
				}
			})
			inputSection(c, "Live", func() {
				r := input.TextArea(c, &s.value, input.TextAreaOptions{Label: "Notes", Placeholder: "Write a dedication…", MinHeight: 80, MaxHeight: 160,
					MaxLength: int(s.max), ReadOnly: s.readOnly, Disabled: s.disabled})
				if r.Changed() {
					log(fmt.Sprintf("TextArea changed: %d characters", len([]rune(s.value))))
				}
				if r.Submitted() {
					log("TextArea submitted")
				}
			})
		}},
	{ID: "026", Name: "NumberInput", ZhName: "数字输入框", Category: inputCategory, Keys: "Up/Down 按步长增减；步进模式下 Up/Down/Left/Right、Home、End。",
		Code: `if input.NumberInput(c, &copies, input.NumberInputOptions{Min: 1, Max: 99, Label: "Copies"}).Changed() {
	fmt.Println("copies", copies)
}`,
		Demo: func(c *ui.Context, log func(string)) {
			s := ui.Local(c.Root(), "026", func() inputNumberDemo { return inputNumberDemo{value: 3, step: 1} })
			inputSection(c, "Parameters", func() {
				inputFlags(c, log, map[string]*bool{"AllowEmpty": &s.empty, "StepperOnly": &s.stepper, "Disabled": &s.disabled, "Error": &s.fail}, "AllowEmpty", "StepperOnly", "Disabled", "Error")
				steps := []input.RadioItem[float64]{{Value: 1, Label: "Step 1"}, {Value: 0.5, Label: "Step 0.5"}, {Value: 0.01, Label: "Step 0.01"}}
				if input.RadioGroup(c, &s.step, steps, input.RadioGroupOptions{Horizontal: true, Label: "Step"}).Changed() {
					log(fmt.Sprintf("option Step = %v", s.step))
				}
			})
			if s.stepper {
				s.empty = false
			}
			if !s.empty && math.IsNaN(s.value) {
				s.value = 0
			}
			inputSection(c, "Live", func() {
				opts := input.NumberInputOptions{Min: 0, Max: 10, Step: s.step, AllowEmpty: s.empty, StepperOnly: s.stepper, Disabled: s.disabled, Label: "Quantity", Placeholder: "Empty"}
				if s.fail {
					opts.Error = "Not in stock"
				}
				ui.Box(c).Width(240).Children(func() {
					if input.NumberInput(c, &s.value, opts).Changed() {
						log(fmt.Sprintf("NumberInput changed: %v", s.value))
					}
				})
			})
			inputSection(c, "States", func() {
				ui.Row(c).Gap(12).AlignItems(ui.Start).Children(func() {
					a, b, e := 4.0, 2.0, 7.5
					input.NumberInput(c, &a, input.NumberInputOptions{Min: 1, Max: 8, StepperOnly: true, Label: "Stepper state"})
					input.NumberInput(c, &b, input.NumberInputOptions{Label: "Disabled state", Disabled: true}).Element.Width(160)
					input.NumberInput(c, &e, input.NumberInputOptions{Step: 0.5, Label: "Error state", Error: "Must be whole"}).Element.Width(160)
				})
			})
		}},
	{ID: "027", Name: "PasswordInput", ZhName: "密码输入框", Category: inputCategory, Keys: "Tab 依次聚焦输入框与显隐按钮；Space 或 Enter 切换显隐；Enter 在输入框中提交；粘贴可用，遮蔽时不可复制。",
		Code: `if input.PasswordInput(c, &password, input.PasswordInputOptions{Label: "Password"}).Submitted() {
	fmt.Println("signing in")
}`,
		Demo: func(c *ui.Context, log func(string)) {
			s := ui.Local(c.Root(), "027", func() inputTextDemo { return inputTextDemo{} })
			inputSection(c, "Parameters", func() {
				inputFlags(c, log, map[string]*bool{"Disabled": &s.disabled, "Error": &s.fail}, "Disabled", "Error")
			})
			inputSection(c, "Live", func() {
				opts := input.PasswordInputOptions{Label: "Password", Placeholder: "Password", Disabled: s.disabled}
				if s.fail {
					opts.Error = "At least 8 characters"
				}
				r := input.PasswordInput(c, &s.value, opts)
				if r.Changed() {
					// The contents never reach the log.
					log(fmt.Sprintf("PasswordInput changed: %d characters", len([]rune(s.value))))
				}
				if r.Submitted() {
					log("PasswordInput submitted")
				}
			})
		}},
	{ID: "028", Name: "SearchInput", ZhName: "搜索框", Category: inputCategory, Keys: "Enter 提交（搜索模式）或下一项（查找栏）；Shift+Enter 上一项；Cmd+G / F3 下一项，Shift 加上为上一项；Escape 清空。",
		Code: `matches := 0
if query != "" {
	matches = strings.Count(text, query)
}
r := input.SearchInput(c, &query, input.SearchInputOptions{Find: true, Matches: matches, Current: &current})
if r.Moved() {
	fmt.Println("show match", current)
}`,
		Demo: func(c *ui.Context, log func(string)) {
			s := ui.Local(c.Root(), "028", func() inputSearchDemo { return inputSearchDemo{} })
			const text = "The rose of the court, the rose of the chapel, and the rose of the garden all bloom in May."
			inputSection(c, "Parameters", func() {
				inputFlags(c, log, map[string]*bool{"Disabled": &s.disabled}, "Disabled")
			})
			inputSection(c, "Search mode", func() {
				r := input.SearchInput(c, &s.query, input.SearchInputOptions{Disabled: s.disabled})
				if r.Changed() {
					log(fmt.Sprintf("SearchInput changed: %q", s.query))
				}
				if r.Submitted() {
					log(fmt.Sprintf("SearchInput submitted: %q", s.query))
				}
			})
			inputSection(c, "Find bar mode", func() {
				n := 0
				if s.find != "" {
					n = strings.Count(text, s.find)
				}
				r := input.SearchInput(c, &s.find, input.SearchInputOptions{Find: true, Matches: n, Current: &s.current, Disabled: s.disabled})
				if r.Changed() {
					log(fmt.Sprintf("find query: %q, %d matches", s.find, n))
				}
				if r.Moved() {
					log(fmt.Sprintf("find moved to match %d", s.current+1))
				}
				ui.Text(c, text).TextColor(core.Tokens(c).TextMuted)
			})
		}},
	{ID: "029", Name: "Checkbox", ZhName: "复选框", Category: inputCategory, Keys: "Tab 聚焦；Space 切换；点击标签同样切换；全选框在组内勾选全部或清空。",
		Code: `input.CheckboxGroup(c, "Notifications", []input.CheckboxItem{
	{Label: "Mail", Value: &mail},
	{Label: "Calendar", Value: &calendar},
}, input.CheckboxGroupOptions{})`,
		Demo: func(c *ui.Context, log func(string)) {
			s := ui.Local(c.Root(), "029", func() inputCheckDemo { return inputCheckDemo{items: [4]bool{true, false, false, true}} })
			inputSection(c, "Parameters", func() {
				inputFlags(c, log, map[string]*bool{"Indeterminate": &s.mixed, "Disabled": &s.disabled}, "Indeterminate", "Disabled")
			})
			inputSection(c, "Live", func() {
				if input.Checkbox(c, &s.single, "Accept the terms", input.CheckboxOptions{Indeterminate: s.mixed, Disabled: s.disabled}).Changed() {
					s.mixed = false
					log(fmt.Sprintf("Checkbox changed: %v", s.single))
				}
			})
			inputSection(c, "Group with select all", func() {
				items := []input.CheckboxItem{
					{Label: "Mail", Value: &s.items[0]}, {Label: "Calendar", Value: &s.items[1]},
					{Label: "Messages", Value: &s.items[2]}, {Label: "Security alerts (required)", Value: &s.items[3], Disabled: true},
				}
				if input.CheckboxGroup(c, "Notifications", items, input.CheckboxGroupOptions{Disabled: s.disabled}).Changed() {
					log(fmt.Sprintf("CheckboxGroup changed: %v", s.items))
				}
			})
		}},
	{ID: "030", Name: "RadioGroup", ZhName: "单选组", Category: inputCategory, Keys: "Tab 进入组并落在已选项；方向键选择并跳过禁用项；Home、End 选择首末项。",
		Code: `items := []input.RadioItem[string]{{Value: "s", Label: "Small"}, {Value: "m", Label: "Medium"}, {Value: "l", Label: "Large"}}
if input.RadioGroup(c, &size, items, input.RadioGroupOptions{Label: "Size"}).Changed() {
	fmt.Println("size", size)
}`,
		Demo: func(c *ui.Context, log func(string)) {
			s := ui.Local(c.Root(), "030", func() inputRadioDemo { return inputRadioDemo{value: "m"} })
			inputSection(c, "Parameters", func() {
				inputFlags(c, log, map[string]*bool{"Horizontal": &s.horizontal, "Disabled": &s.disabledAll}, "Horizontal", "Disabled")
			})
			inputSection(c, "Live", func() {
				items := []input.RadioItem[string]{{Value: "s", Label: "Small"}, {Value: "m", Label: "Medium"}, {Value: "l", Label: "Large"}, {Value: "xl", Label: "Sold out", Disabled: true}}
				if input.RadioGroup(c, &s.value, items, input.RadioGroupOptions{Horizontal: s.horizontal, Disabled: s.disabledAll, Label: "Size"}).Changed() {
					log("RadioGroup chose " + s.value)
				}
			})
		}},
}
