package golog

const (
	Reset      ColorCode = "\033[0m"
	Bold       ColorCode = "\033[1m"
	Red        ColorCode = "\033[31m"
	Green      ColorCode = "\033[32m"
	Yellow     ColorCode = "\033[33m"
	Blue       ColorCode = "\033[34m"
	Magenta    ColorCode = "\033[35m"
	Cyan       ColorCode = "\033[36m"
	White      ColorCode = "\033[37m"
	BoldRed    ColorCode = "\033[1;31m"
	BoldGreen  ColorCode = "\033[1;32m"
	BoldYellow ColorCode = "\033[1;33m"
	BoldBlue   ColorCode = "\033[1;34m"
	BoldPurple ColorCode = "\033[1;35m"
	BoldCyan   ColorCode = "\033[1;36m"
	BoldWhite  ColorCode = "\033[1;37m"
	BlackBg    ColorCode = "\033[40m"
	RedBg      ColorCode = "\033[41m"
	GreenBg    ColorCode = "\033[42m"
	YellowBg   ColorCode = "\033[43m"
	BlueBg     ColorCode = "\033[44m"
	MagentaBg  ColorCode = "\033[45m"
	CyanBg     ColorCode = "\033[46m"
	WhiteBg    ColorCode = "\033[47m"
)

var (
	RainbowTheme           = []ColorCode{Red, Yellow, Green, Blue, Magenta, Cyan, White}
	BoldRainbowTheme       = []ColorCode{BoldRed, BoldYellow, BoldGreen, BoldBlue, BoldPurple, BoldCyan, BoldWhite}
	RainbowBackgroundTheme = []ColorCode{RedBg, YellowBg, GreenBg, BlueBg, MagentaBg, CyanBg, WhiteBg}
)
