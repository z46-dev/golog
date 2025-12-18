package golog

import (
	"fmt"
	"time"
)

func newLoader(message string, loaderType LoaderType, tps int, logger *Logger) *Loader {
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
	}
}

func (l *Loader) Start() {
	l.mu.Lock()
	if l.running {
		l.mu.Unlock()
		return
	}

	l.stop = make(chan struct{})
	l.done = make(chan struct{})
	l.running = true
	l.paused = false
	interval := time.Second / time.Duration(l.tps)
	if l.logger != nil {
		l.logger.loader = l
		l.logger.spinner = nil
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
				if l.logger != nil {
					l.logger.loader = nil
				}
				l.mu.Unlock()
				return
			case <-ticker.C:
				l.render()
			}
		}
	}()
}

func (l *Loader) Stop() {
	l.mu.Lock()
	if !l.running {
		l.mu.Unlock()
		return
	}
	close(l.stop)
	done := l.done
	l.mu.Unlock()

	<-done
	clearSpinnerLine()
}

func (l *Loader) SetProgress(progress float64) {
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
}

func (l *Loader) render() {
	l.mu.Lock()
	if !l.running || l.paused {
		l.mu.Unlock()
		return
	}

	progress := l.progress
	pattern := l.pattern
	logger := l.logger
	l.mu.Unlock()

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
