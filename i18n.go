package mujica

import "github.com/egoist/mygo/ui"

// Locale picks the built-in copy. The zero value is English.
type Locale string

const (
	EnUS Locale = "en-US"
	ZhCN Locale = "zh-CN"
)

// messages maps a key to its English and Chinese copy.
var messages = map[string][2]string{}

func addMessages(m map[string][2]string) {
	for k, v := range m {
		if _, dup := messages[k]; dup {
			panic("mujica: duplicate message " + k)
		}
		messages[k] = v
	}
}

// msg returns the window's copy for key.
func msg(c *ui.Context, key string) string {
	w := current(c)
	if s, ok := w.messages[key]; ok {
		return s
	}
	m, ok := messages[key]
	if !ok {
		panic("mujica: unknown message " + key)
	}
	if w.locale == ZhCN {
		return m[1]
	}
	return m[0]
}

// CurrentLocale returns the window's locale.
func CurrentLocale(c *ui.Context) Locale {
	if l := current(c).locale; l != "" {
		return l
	}
	return EnUS
}
