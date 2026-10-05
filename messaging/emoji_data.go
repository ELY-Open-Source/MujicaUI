package messaging

type emoji struct {
	char, name string
	toned      bool
}

// emojiCategory is a tab of the picker.
type emojiCategory struct {
	key    string
	icon   string
	emojis []emoji
}

// emojiTones are the Fitzpatrick modifiers, after the default yellow.
var emojiTones = []string{"", "\U0001F3FB", "\U0001F3FC", "\U0001F3FD", "\U0001F3FE", "\U0001F3FF"}

var emojiCategories = []emojiCategory{
	{"emoji_picker.smileys", "😀", []emoji{
		{"😀", "grinning face", false}, {"😄", "smiling face with open mouth", false}, {"😂", "face with tears of joy", false},
		{"🙂", "slightly smiling face", false}, {"😉", "winking face", false}, {"😊", "smiling face with smiling eyes", false},
		{"😍", "smiling face with heart eyes", false}, {"😘", "face blowing a kiss", false}, {"🤔", "thinking face", false},
		{"😐", "neutral face", false}, {"🙄", "face with rolling eyes", false}, {"😴", "sleeping face", false},
		{"😎", "smiling face with sunglasses", false}, {"🤩", "star struck", false}, {"😢", "crying face", false},
		{"😡", "pouting face", false}, {"😱", "face screaming in fear", false}, {"🥳", "partying face", false},
	}},
	{"emoji_picker.people", "👋", []emoji{
		{"👋", "waving hand", true}, {"👍", "thumbs up", true}, {"👎", "thumbs down", true}, {"👏", "clapping hands", true},
		{"🙌", "raising hands", true}, {"🙏", "folded hands", true}, {"👌", "ok hand", true}, {"✌️", "victory hand", false},
		{"🤝", "handshake", false}, {"💪", "flexed biceps", true}, {"✍️", "writing hand", false}, {"👀", "eyes", false},
		{"🧑", "person", true}, {"👑", "crown", false}, {"🎭", "performing arts masks", false}, {"💃", "woman dancing", true},
	}},
	{"emoji_picker.nature", "🌹", []emoji{
		{"🌹", "rose", false}, {"🥀", "wilted flower", false}, {"🌸", "cherry blossom", false}, {"🌙", "crescent moon", false},
		{"⭐", "star", false}, {"🔥", "fire", false}, {"🌧️", "cloud with rain", false}, {"❄️", "snowflake", false},
		{"🐈‍⬛", "black cat", false}, {"🦇", "bat", false}, {"🐦‍⬛", "black bird", false}, {"🕊️", "dove", false},
		{"🍂", "fallen leaf", false}, {"🌲", "evergreen tree", false}, {"🌊", "water wave", false}, {"🌈", "rainbow", false},
	}},
	{"emoji_picker.food", "🍷", []emoji{
		{"🍷", "wine glass", false}, {"🥂", "clinking glasses", false}, {"☕", "hot beverage", false}, {"🍵", "teacup", false},
		{"🍰", "shortcake", false}, {"🍫", "chocolate bar", false}, {"🍓", "strawberry", false}, {"🍇", "grapes", false},
		{"🍎", "red apple", false}, {"🧀", "cheese wedge", false}, {"🥖", "baguette", false}, {"🍞", "bread", false},
	}},
	{"emoji_picker.activities", "🎻", []emoji{
		{"🎻", "violin", false}, {"🎹", "musical keyboard", false}, {"🎸", "guitar", false}, {"🥁", "drum", false},
		{"🎤", "microphone", false}, {"🎵", "musical note", false}, {"🎶", "musical notes", false}, {"🎉", "party popper", false},
		{"🏆", "trophy", false}, {"🎨", "artist palette", false}, {"♟️", "chess pawn", false}, {"🎬", "clapper board", false},
	}},
	{"emoji_picker.travel", "🏰", []emoji{
		{"🏰", "castle", false}, {"⛪", "church", false}, {"🗼", "tower", false}, {"🌃", "night with stars", false},
		{"🚂", "locomotive", false}, {"✈️", "airplane", false}, {"🚗", "automobile", false}, {"⛵", "sailboat", false},
		{"🗺️", "world map", false}, {"🏔️", "snow capped mountain", false}, {"🌅", "sunrise", false}, {"🎡", "ferris wheel", false},
	}},
	{"emoji_picker.objects", "🕯️", []emoji{
		{"🕯️", "candle", false}, {"🗝️", "old key", false}, {"📜", "scroll", false}, {"✉️", "envelope", false},
		{"📎", "paperclip", false}, {"📌", "pushpin", false}, {"📅", "calendar", false}, {"💡", "light bulb", false},
		{"🔔", "bell", false}, {"🔒", "locked", false}, {"💎", "gem stone", false}, {"⏳", "hourglass", false},
	}},
	{"emoji_picker.symbols", "❤️", []emoji{
		{"❤️", "red heart", false}, {"🖤", "black heart", false}, {"💔", "broken heart", false}, {"✨", "sparkles", false},
		{"✅", "check mark button", false}, {"❌", "cross mark", false}, {"❗", "exclamation mark", false}, {"❓", "question mark", false},
		{"💯", "hundred points", false}, {"♦️", "diamond suit", false}, {"⚜️", "fleur de lis", false}, {"🔺", "red triangle", false},
	}},
}

func emojiToned(e emoji, tone int) string {
	if !e.toned || tone == 0 {
		return e.char
	}
	return e.char + emojiTones[tone]
}
