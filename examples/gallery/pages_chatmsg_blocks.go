package main

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"

	mujica "mujicaui"
)

const chatmsgMarkdown = "## The Decree\n\nBy order of **the Queen**, *all* guests use `the east door`. See [the protocol](https://example.com/protocol).\n\n" +
	"1. Arrive at nine\n2. Bow to the throne\n   - dukes first\n\n> Silence in the Hall of Mirrors.\n\n" +
	"```python\ndef bow(guest):\n    return f\"{guest} bows\"  # deeply\n```\n\n| Guest | Seat |\n|---|--:|\n| Marie | 1 |\n| Axel | 14 |\n\n---\n\nVive la reine."

// chatmsgPicture is a flat colored picture for image demos.
func chatmsgPicture(color string) ui.ImageSource {
	return ui.MustParseSVG([]byte(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="240" height="180"><rect width="240" height="180" fill="%s"/><circle cx="120" cy="90" r="40" fill="#FBF7EF" opacity="0.5"/></svg>`, color)))
}

var chatmsgColors = []string{"#7D2034", "#AD9164", "#355E80", "#326247", "#805714", "#6A5B55", "#9C3C53", "#D2C3AF", "#211C1D"}

var chatmsgBlockPages = []page{
	{ID: "206", Name: "MarkdownView", ZhName: "Markdown 视图", Category: chatmsgCategory, Keys: "Tab 到链接，Enter 激活；代码块与表格按钮按 Tab 顺序访问。",
		Code: `if url, ok := mujica.MarkdownView(c, "# Decree\n\nSee **the guide** at [the court](https://example.com).").Link(); ok {
	fmt.Println("open", url)
}`,
		Demo: func(c *ui.Context, log func(string)) {
			type state struct {
				src    string
				stream bool
				at     time.Time
			}
			s := ui.Local(c.Root(), "206", func() state { return state{src: chatmsgMarkdown} })
			inputSection(c, "Parameters", func() {
				ui.Row(c).Children(func() {
					if mujica.Button(c, "Stream it again", mujica.ButtonOptions{Variant: mujica.Secondary}).Clicked() {
						s.src, s.stream = "", true
						log("streaming started")
					}
				})
				mujica.TextArea(c, &s.src, mujica.TextAreaOptions{Label: "Markdown source", MaxHeight: 160})
			})
			if s.stream && len(s.src) < len(chatmsgMarkdown) && c.Now().Sub(s.at) > 40*time.Millisecond {
				s.src = chatmsgMarkdown[:min(len(chatmsgMarkdown), len(s.src)+6)]
				s.at = c.Now()
				if len(s.src) == len(chatmsgMarkdown) {
					s.stream = false
					log("streaming finished; finished blocks were never parsed again")
				}
			}
			if s.stream {
				c.After(40 * time.Millisecond)
			}
			inputSection(c, "Live", func() {
				if url, ok := mujica.MarkdownView(c, s.src).Link(); ok {
					log("link " + url)
				}
			})
		}},
	{ID: "207", Name: "CodeBlock", ZhName: "代码块", Category: chatmsgCategory, Keys: "Tab 到换行、复制、折叠按钮；Space 或 Enter 触发。",
		Code: `mujica.CodeBlock(c, "func bow() {}", mujica.CodeBlockOptions{Language: "go", LineNumbers: true})`,
		Demo: func(c *ui.Context, log func(string)) {
			type state struct {
				lang          string
				numbers, wrap bool
			}
			s := ui.Local(c.Root(), "207", func() state { return state{lang: "go", numbers: true} })
			samples := map[string]string{
				"go":     "// greet the court\nfunc greet(name string) int {\n\treturn len(\"Your Majesty, \" + name + \", the ball awaits you in the Hall of Mirrors\") + 42\n}",
				"python": "def bow(guest):\n    # deeply, always\n    return f\"{guest} bows\" * 3",
				"sql":    "SELECT name, seat FROM guests WHERE rank > 2 ORDER BY seat LIMIT 10; -- the envoys",
			}
			inputSection(c, "Parameters", func() {
				langs := []mujica.RadioItem[string]{{Value: "go", Label: "Go"}, {Value: "python", Label: "Python"}, {Value: "sql", Label: "SQL"}}
				if mujica.RadioGroup(c, &s.lang, langs, mujica.RadioGroupOptions{Horizontal: true, Label: "Language"}).Changed() {
					log("language " + s.lang)
				}
				inputFlags(c, log, map[string]*bool{"LineNumbers": &s.numbers, "Wrap at first": &s.wrap}, "LineNumbers", "Wrap at first")
			})
			inputSection(c, "Live", func() {
				ui.Column(c).Key(s.lang + fmt.Sprint(s.wrap)).Children(func() {
					if mujica.CodeBlock(c, samples[s.lang], mujica.CodeBlockOptions{Language: s.lang, LineNumbers: s.numbers, Wrap: s.wrap}).Copied() {
						log("code copied")
					}
				})
			})
			inputSection(c, "States", func() {
				mujica.CodeBlock(c, samples["go"], mujica.CodeBlockOptions{Language: "go", Collapsed: true})
			})
		}},
	{ID: "208", Name: "TableBlock", ZhName: "表格块", Category: chatmsgCategory, Keys: "Tab 到复制按钮，Enter 或 Space 复制为 TSV。",
		Code: `if mujica.TableBlock(c, []string{"Guest", "Seat"}, [][]string{{"Marie", "1"}}, mujica.TableBlockOptions{}).Copied() {
	fmt.Println("copied as TSV")
}`,
		Demo: func(c *ui.Context, log func(string)) {
			wide := ui.Local(c.Root(), "208", func() bool { return false })
			inputSection(c, "Parameters", func() {
				inputFlags(c, log, map[string]*bool{"Wide table": wide}, "Wide table")
			})
			header := []string{"Guest", "Title", "Seat"}
			rows := [][]string{{"Marie", "Queen of France", "1"}, {"Axel", "Count of Fersen", "14"}}
			if *wide {
				header = append(header, "Arrival", "Carriage", "Retinue", "Notes")
				for i := range rows {
					rows[i] = append(rows[i], "21:00", "Gilded coach", "Twelve", "Seated near the orchestra, away from the draughts")
				}
			}
			inputSection(c, "Live", func() {
				ui.Column(c).MaxWidth(520).Children(func() {
					if mujica.TableBlock(c, header, rows, mujica.TableBlockOptions{Align: []ui.Align{ui.Start, ui.Start, ui.End}}).Copied() {
						log("table copied as TSV")
					}
				})
			})
		}},
	{ID: "209", Name: "ImageGrid", ZhName: "图片网格", Category: chatmsgCategory, Keys: "Tab 聚焦图片，Enter 打开查看器；查看器内 Left/Right 切换，Escape 关闭。",
		Code: `if i, ok := mujica.ImageGrid(c, []mujica.GridImage{{Source: pic, Alt: "Portrait"}}).Retried(); ok {
	fmt.Println("reload", i)
}`,
		Demo: func(c *ui.Context, log func(string)) {
			type state struct {
				n             float64
				loading, fail bool
			}
			s := ui.Local(c.Root(), "209", func() state { return state{n: 5} })
			inputSection(c, "Parameters", func() {
				if mujica.NumberInput(c, &s.n, mujica.NumberInputOptions{Min: 1, Max: 9, Label: "Images", StepperOnly: true}).Changed() {
					log(fmt.Sprintf("%v images", s.n))
				}
				inputFlags(c, log, map[string]*bool{"Last loading": &s.loading, "Second failed": &s.fail}, "Last loading", "Second failed")
			})
			imgs := make([]mujica.GridImage, int(s.n))
			for i := range imgs {
				imgs[i] = mujica.GridImage{Source: chatmsgPicture(chatmsgColors[i]), Alt: fmt.Sprintf("Portrait %d", i+1)}
			}
			if s.loading {
				imgs[len(imgs)-1] = mujica.GridImage{Alt: "Loading portrait", Loading: true}
			}
			if s.fail && len(imgs) > 1 {
				imgs[1] = mujica.GridImage{Alt: "Lost portrait", Err: errors.New("404")}
			}
			inputSection(c, "Live", func() {
				if i, ok := mujica.ImageGrid(c, imgs).Retried(); ok {
					s.fail = false
					log(fmt.Sprintf("retry image %d", i))
				}
			})
		}},
	{ID: "210", Name: "FileMessage", ZhName: "文件消息", Category: chatmsgCategory, Keys: "Tab 到文件与下载按钮，Enter 或 Space 触发。",
		Code: `r := mujica.FileMessage(c, mujica.FileMessageOptions{Name: "treaty.pdf", Size: 2621440})
if r.Downloaded() {
	fmt.Println("save treaty.pdf")
}`,
		Demo: func(c *ui.Context, log func(string)) {
			disabled := ui.Local(c.Root(), "210", func() bool { return false })
			inputSection(c, "Parameters", func() {
				inputFlags(c, log, map[string]*bool{"Disabled": disabled}, "Disabled")
			})
			inputSection(c, "Live", func() {
				ui.Column(c).Gap(8).Disabled(*disabled).Children(func() {
					for _, f := range []mujica.FileMessageOptions{{Name: "treaty.pdf", Size: 2621440}, {Name: "portrait.png", Size: 812000},
						{Name: "ledger.zip", Size: 73400320}, {Name: "decree.go", Size: 1840}, {Name: "notes", Size: 12}} {
						r := mujica.FileMessage(c, f)
						if r.Opened() {
							log("open " + f.Name)
						}
						if r.Downloaded() {
							log("download " + f.Name)
						}
					}
				})
			})
		}},
}

// chatmsgWave is a fixed waveform for the voice demo.
var chatmsgWave = func() []float32 {
	w := make([]float32, 48)
	for i := range w {
		w[i] = float32(0.25 + 0.7*float64((i*37)%11)/10)
		if w[i] > 1 {
			w[i] = 1
		}
	}
	return w
}()

// chatmsgSnippet trims s to n runes for logs.
func chatmsgSnippet(s string, n int) string {
	r := []rune(strings.TrimSpace(s))
	if len(r) > n {
		return string(r[:n]) + "…"
	}
	return string(r)
}
