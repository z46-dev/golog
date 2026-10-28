package golog

import (
	"fmt"
	"time"
)

func newLoader(message string, loaderType LoaderType, tps int, logger *Logger, host *Logger) *Loader {
	pattern, ok := loaders[loaderType]
	if !ok {
		pattern = LoaderPattern{Width: 20, Fill: '=', Arrow: '>', Empty: ' '}
	}
	if pattern.Width <= 0 {
		pattern.Width = 20
	}

	return &Loader{
		message:  message,
		progress: 0,
		tps:      max(1, tps),
		pattern:  pattern,
		stop:     nil,
		done:     nil,
		logger:   logger,
		host:     host,
	}
}

func (l *Loader) Start() *Loader {
	l.mu.Lock()
	if l.running {
		l.mu.Unlock()
		return l
	}

	l.stop = make(chan struct{})
	l.done = make(chan struct{})
	l.running = true
	l.paused = false
	interval := time.Second / time.Duration(l.tps)
	if l.host != nil {
		l.host.loader = l
		l.host.spinner = nil
	}
	l.mu.Unlock()

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-l.stop:
				clearSpinnerLine()
				close(l.done)
				l.mu.Lock()
				l.running = false
				l.paused = false
				if l.host != nil {
					l.host.loader = nil
				}
				l.mu.Unlock()
				return
			case <-ticker.C:
				l.render()
			}
		}
	}()
	return l
}

func (l *Loader) Stop() *Loader {
	l.mu.Lock()
	if !l.running {
		l.mu.Unlock()
		return l
	}
	close(l.stop)
	done := l.done
	l.mu.Unlock()

	<-done
	clearSpinnerLine()
	return l
}

func (l *Loader) SetProgress(progress float64) *Loader {
	l.mu.Lock()
	if progress < 0 {
		progress = 0
	}
	if progress > 1 {
		progress = 1
	}
	l.progress = progress
	running := l.running && !l.paused
	l.mu.Unlock()

	if running {
		l.render()
	}
	return l
}

// SetMessage updates the text shown beside the loader.
func (l *Loader) SetMessage(message string) *Loader {
	l.mu.Lock()
	l.message = message
	running := l.running && !l.paused
	l.mu.Unlock()

	if running {
		l.render()
	}
	return l
}

func (l *Loader) render() {
	l.mu.Lock()
	defer l.mu.Unlock()
	if !l.running || l.paused {
		return
	}

	progress := l.progress
	pattern := l.pattern
	logger := l.logger

	clearSpinnerLine()
	bar := buildLoaderBar(progress, pattern)
	if logger != nil {
		fmt.Printf("\r%s[%s] %s", logger.spinnerPrefix(), bar, l.message)
	} else {
		fmt.Printf("\r[%s] %s", bar, l.message)
	}
}

func (l *Loader) isRunning() bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.running
}
