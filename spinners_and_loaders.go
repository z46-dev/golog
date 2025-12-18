package golog

const (
	SpinnerDots SpinnerType = iota
	SpinnerDots2
	SpinnerDots3
	SpinnerLine
	SpinnerClock
	SpinnerEarth
	SpinnerMoon
	SpinnerRunner
	SpinnerWeather
	SpinnerMindblown
	SpinnerOrangePulse
	SpinnerBluePulse
	SpinnerOrangeBluePulse
	SpinnerTimeTravel
)

const (
	LoaderBar LoaderType = iota
)

var spinners = map[SpinnerType][]rune{
	SpinnerDots: {'⠋', '⠙', '⠹', '⠸', '⠼', '⠴', '⠦', '⠧', '⠇', '⠏'},

	SpinnerDots2: {'⣾', '⣽', '⣻', '⢿', '⡿', '⣟', '⣯', '⣷'},

	SpinnerDots3: {'⠋', '⠙', '⠚', '⠞', '⠖', '⠦', '⠴', '⠲', '⠳', '⠓'},

	SpinnerLine: {'-', '\\', '|', '/'},

	SpinnerClock: {'🕛', '🕐', '🕑', '🕒', '🕓', '🕔', '🕕', '🕖', '🕗', '🕘', '🕙', '🕚'},

	SpinnerEarth: {'🌍', '🌎', '🌏'},

	SpinnerMoon: {'🌑', '🌒', '🌓', '🌔', '🌕', '🌖', '🌗', '🌘'},

	SpinnerRunner: {'🚶', '🏃'},

	SpinnerWeather: {
		'☀', '☀', '☀', '🌤', '⛅', '🌥', '☁', '🌧', '🌨', '🌧', '🌨', '🌧',
		'🌨', '⛈', '🌨', '🌧', '🌨', '☁', '🌥', '⛅', '🌤', '☀', '☀',
	},

	SpinnerMindblown: {'😐', '😐', '😮', '😮', '😦', '😦', '😧', '😧', '🤯', '💥', '✨', ' ', ' ', ' '},

	SpinnerOrangePulse: {'🔸', '🔶', '🟠', '🟠', '🔶'},

	SpinnerBluePulse: {'🔹', '🔷', '🔵', '🔵', '🔷'},

	SpinnerOrangeBluePulse: {'🔸', '🔶', '🟠', '🟠', '🔶', '🔹', '🔷', '🔵', '🔵', '🔷'},

	SpinnerTimeTravel: {'🕛', '🕚', '🕙', '🕘', '🕗', '🕖', '🕕', '🕔', '🕓', '🕒', '🕑', '🕐'},
}

var loaders = map[LoaderType]LoaderPattern{
	LoaderBar: {
		Width: 20,
		Fill:  '=',
		Arrow: '>',
		Empty: ' ',
	},
}
