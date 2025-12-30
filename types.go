package golog

import (
	"bufio"
	"os"
	"sync"
)

type (
	LogLevel uint8 // LogLevel represents the level of the log message.

	ColorCode string // colorCode represents the color code for the logger output.

	LogFlushMode uint8 // LogFlushMode controls how log file writes are flushed.

	TimestampPrecision uint8 // TimestampPrecision represents the precision for timestamps.

	SpinnerType uint8 // SpinnerType represents the type of spinner to use.
	LoaderType  uint8 // LoaderType represents the type of loader to use.

	Builder struct {
		message string
	} // Builder is a struct that represents a log message builder. Can be used to create complex color messages.

	Logger struct {
		hasPrefix, includeTimestamp     bool
		prefix                          string
		color                           ColorCode
		levelWithSymbol, levelWithColor bool
		timestampPrecision              TimestampPrecision
		spinner                         *Spinner
		loader                          *Loader
		logFile                         *os.File
		logWriter                       *bufio.Writer
		logFlushMode                    LogFlushMode
		logFlushStop                    chan struct{}
		logMu                           sync.Mutex
	} // Logger is a struct that represents a logger with various configurations.

	Spinner struct {
		message     string
		tick, tps   int
		spinnerType SpinnerType
		frames      []rune
		stop        chan struct{}
		done        chan struct{}
		running     bool
		paused      bool
		mu          sync.Mutex
		logger      *Logger
	} // Spinner is a struct that represents a spinner with a message, tick interval, and ticks per second (tps).

	LoaderPattern struct {
		Width int
		Fill  rune
		Arrow rune
		Empty rune
	} // LoaderPattern is a struct that represents the pattern of a loader.

	Loader struct {
		message  string
		progress float64
		tps      int
		pattern  LoaderPattern
		stop     chan struct{}
		done     chan struct{}
		running  bool
		paused   bool
		mu       sync.Mutex
		logger   *Logger
	} // Loader is a struct that represents a loader with a message, progress, ticks per second (tps), and pattern.
)
