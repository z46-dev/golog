package golog

const (
	LevelDebug = iota
	LevelInfo
	LevelWarning
	LevelError
	LevelFatal
)

var levelNames map[LogLevel]string = map[LogLevel]string{
	LevelDebug:   "Debug",
	LevelInfo:    "Info",
	LevelWarning: "Warning",
	LevelError:   "Error",
	LevelFatal:   "Fatal",
}

var levelSymbols map[LogLevel]string = map[LogLevel]string{
	LevelDebug:   "?",
	LevelInfo:    ":",
	LevelWarning: "!",
	LevelError:   ")",
	LevelFatal:   ">",
}

var levelColors map[LogLevel]ColorCode = map[LogLevel]ColorCode{
	LevelDebug:   BoldCyan,
	LevelInfo:    BoldWhite,
	LevelWarning: BoldYellow,
	LevelError:   BoldRed,
	LevelFatal:   RedBg,
}
