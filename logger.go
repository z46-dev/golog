package golog

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

const (
	LogFlushAlways LogFlushMode = iota
	LogFlushBuffered
)

const (
	defaultLogBufferSize    = 64 * 1024
	defaultLogFlushInterval = 1 * time.Second
)

func New() (logger *Logger) {
	logger = &Logger{
		hasPrefix:        false,
		includeTimestamp: false,
		prefix:           "",
		color:            White,
		levelWithSymbol:  false,
		levelWithColor:   true,
	}

	return
}

func (l *Logger) Builder() (builder *Builder) {
	builder = newBuilder()
	return
}

func (l *Logger) Spinner(message string, spinnerType SpinnerType, tps int) (spinner *Spinner) {
	spinner = newSpinner(message, spinnerType, tps, l)
	return
}

func (l *Logger) Loader(message string, loaderType LoaderType, tps int) (loader *Loader) {
	loader = newLoader(message, loaderType, tps, l)
	return
}

func (l *Logger) Prefix(prefix string, color ColorCode) (self *Logger) {
	l.prefix = prefix
	l.color = color
	l.hasPrefix = true
	self = l
	return
}

func (l *Logger) ClearPrefix() (self *Logger) {
	l.prefix = ""
	l.color = White
	l.hasPrefix = false
	self = l
	return
}

func (l *Logger) Timestamp() (self *Logger) {
	l.includeTimestamp = true
	self = l
	return
}

func (l *Logger) NoTimestamp() (self *Logger) {
	l.includeTimestamp = false
	self = l
	return
}

func (l *Logger) Representation(useSymbol bool, colored bool) (self *Logger) {
	l.levelWithSymbol = useSymbol
	l.levelWithColor = colored
	self = l
	return
}

func (l *Logger) LogFile(path string, mode LogFlushMode) (self *Logger) {
	l.logMu.Lock()
	defer l.logMu.Unlock()

	l.stopLogFlusherLocked()
	if l.logWriter != nil {
		_ = l.logWriter.Flush()
		l.logWriter = nil
	}
	if l.logFile != nil {
		_ = l.logFile.Close()
		l.logFile = nil
	}

	if path == "" {
		self = l
		return
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		self = l
		return
	}

	l.logFile = file
	l.logWriter = bufio.NewWriterSize(file, defaultLogBufferSize)
	l.logFlushMode = mode

	if mode == LogFlushBuffered {
		l.startLogFlusherLocked()
	}

	self = l
	return
}

func (l *Logger) timestamp() (timestamp string) {
	timestamp = time.Now().Format("01/02 15:04")
	return
}

func (l *Logger) build(level LogLevel, format string, args ...any) (output string) {
	format = "%s[%s]%s " + format

	var (
		levelText  string
		levelColor ColorCode = White
	)

	if l.levelWithSymbol {
		levelText = levelSymbols[level]
	} else {
		levelText = levelNames[level]
	}

	if l.levelWithColor {
		levelColor = levelColors[level]
	}

	args = append([]any{levelColor, levelText, Reset}, args...)

	if l.hasPrefix {
		format = "%s%s%s " + format
		args = append([]any{l.color, l.prefix, Reset}, args...)
	}

	if l.includeTimestamp {
		format = "[%s%s%s] " + format
		args = append([]any{White, l.timestamp(), Reset}, args...)
	}

	// fmt.Printf(format, args...)
	output = fmt.Sprintf(format, args...)
	return
}

// For each level, create a Level() and a Levelf() method. Level() should terminate with a \n, while Levelf() should not.

func (l *Logger) Debug(message string) {
	l.printWithSpinner(LevelDebug, l.build(LevelDebug, message), true)
}

func (l *Logger) Debugf(format string, args ...any) {
	l.printWithSpinner(LevelDebug, l.build(LevelDebug, format, args...), false)
}

func (l *Logger) Info(message string) {
	l.printWithSpinner(LevelInfo, l.build(LevelInfo, message), true)
}

func (l *Logger) Infof(format string, args ...any) {
	l.printWithSpinner(LevelInfo, l.build(LevelInfo, format, args...), false)
}

func (l *Logger) Warning(message string) {
	l.printWithSpinner(LevelWarning, l.build(LevelWarning, message), true)
}

func (l *Logger) Warningf(format string, args ...any) {
	l.printWithSpinner(LevelWarning, l.build(LevelWarning, format, args...), false)
}

func (l *Logger) Error(message string) {
	l.printWithSpinner(LevelError, l.build(LevelError, message), true)
}

func (l *Logger) Errorf(format string, args ...any) {
	l.printWithSpinner(LevelError, l.build(LevelError, format, args...), false)
}

func (l *Logger) Fatal(message string) {
	l.printWithSpinner(LevelFatal, l.build(LevelFatal, message), true)
}

func (l *Logger) Fatalf(format string, args ...any) {
	l.printWithSpinner(LevelFatal, l.build(LevelFatal, format, args...), false)
}

func (l *Logger) Panic(message string) {
	panic(l.build(LevelFatal, message))
}

func (l *Logger) Panicf(format string, args ...any) {
	panic(l.build(LevelFatal, format, args...))
}

func (l *Logger) printWithSpinner(level LogLevel, output string, newline bool) {
	l.writeLogFile(level, output, newline)

	if ld := l.loader; ld != nil && ld.isRunning() {
		ld.mu.Lock()
		ld.paused = true

		clearSpinnerLine()
		if newline {
			fmt.Println(output)
		} else {
			fmt.Print(output)
		}

		progress := ld.progress
		pattern := ld.pattern
		msg := ld.message
		prefix := ""
		if ld.logger != nil {
			prefix = ld.logger.spinnerPrefix()
		}

		bar := buildLoaderBar(progress, pattern)
		fmt.Printf("\r%s[%s] %s", prefix, bar, msg)

		ld.paused = false
		ld.mu.Unlock()
		return
	}

	if s := l.spinner; s != nil && s.isRunning() {
		s.mu.Lock()
		s.paused = true

		clearSpinnerLine()
		if newline {
			fmt.Println(output)
		} else {
			fmt.Print(output)
		}

		frame := s.frames[s.tick%len(s.frames)]
		s.tick++
		prefix := ""
		if s.logger != nil {
			prefix = s.logger.spinnerPrefix()
		}
		fmt.Printf("\r%s%s %s", prefix, string(frame), s.message)

		s.paused = false
		s.mu.Unlock()
		return
	}

	if newline {
		fmt.Println(output)
	} else {
		fmt.Print(output)
	}
}

// spinnerPrefix builds the prefix (timestamp, custom prefix) without level.
func (l *Logger) spinnerPrefix() string {
	format := ""
	var args []any

	if l.hasPrefix {
		format = "%s%s%s " + format
		args = append([]any{l.color, l.prefix, Reset}, args...)
	}

	if l.includeTimestamp {
		format = "[%s%s%s] " + format
		args = append([]any{White, l.timestamp(), Reset}, args...)
	}

	if format == "" {
		return ""
	}

	// Ensure we always end with a reset so spinner colors don't leak.
	format += "%s"
	args = append(args, Reset)

	return fmt.Sprintf(format, args...)
}

func buildLoaderBar(progress float64, pattern LoaderPattern) string {
	if pattern.Width <= 0 {
		pattern.Width = 20
	}

	bar := make([]rune, pattern.Width)
	for i := range bar {
		bar[i] = pattern.Empty
	}

	filled := int(progress * float64(pattern.Width))
	if filled > pattern.Width {
		filled = pattern.Width
	}
	for i := 0; i < filled; i++ {
		bar[i] = pattern.Fill
	}

	arrowPos := filled
	if arrowPos >= pattern.Width {
		arrowPos = pattern.Width - 1
	}
	if arrowPos >= 0 {
		bar[arrowPos] = pattern.Arrow
	}

	return string(bar)
}

func (l *Logger) writeLogFile(level LogLevel, output string, newline bool) {
	l.logMu.Lock()
	defer l.logMu.Unlock()

	if l.logWriter == nil {
		return
	}

	output = stripAnsi(output)
	if newline {
		output += "\n"
	}

	_, _ = l.logWriter.WriteString(output)

	if l.logFlushMode == LogFlushAlways || level >= LevelWarning {
		_ = l.logWriter.Flush()
	}
}

func stripAnsi(input string) string {
	// Remove ANSI escape sequences (e.g., color codes) from log output.
	out := make([]byte, 0, len(input))
	for i := 0; i < len(input); i++ {
		if input[i] != 0x1b || i+1 >= len(input) || input[i+1] != '[' {
			out = append(out, input[i])
			continue
		}

		j := i + 2
		for ; j < len(input); j++ {
			b := input[j]
			if b >= 'A' && b <= 'Z' || b >= 'a' && b <= 'z' {
				break
			}
		}

		if j >= len(input) {
			break
		}

		i = j
	}

	return string(out)
}

func (l *Logger) startLogFlusherLocked() {
	if l.logFlushStop != nil {
		return
	}

	l.logFlushStop = make(chan struct{})
	stop := l.logFlushStop

	go func() {
		ticker := time.NewTicker(defaultLogFlushInterval)
		defer ticker.Stop()
		for {
			select {
			case <-stop:
				return
			case <-ticker.C:
				l.logMu.Lock()
				if l.logWriter != nil {
					_ = l.logWriter.Flush()
				}
				l.logMu.Unlock()
			}
		}
	}()
}

func (l *Logger) stopLogFlusherLocked() {
	if l.logFlushStop == nil {
		return
	}

	close(l.logFlushStop)
	l.logFlushStop = nil
}
