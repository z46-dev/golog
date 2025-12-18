package golog

import (
	"fmt"
	"time"
)

func newSpinner(message string, spinnerType SpinnerType, tps int, logger *Logger) (spinner *Spinner) {
	spinner = &Spinner{
		message:     message,
		tick:        0,
		tps:         max(1, tps),
		spinnerType: spinnerType,
		frames:      spinners[spinnerType],
		stop:        nil,
		done:        nil,
		running:     false,
		logger:      logger,
	}

	if len(spinner.frames) == 0 {
		spinner.frames = []rune{'-', '\\', '|', '/'}
	}

	return
}

func (s *Spinner) Start() {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return
	}

	s.stop = make(chan struct{})
	s.done = make(chan struct{})
	s.tick = 0
	s.running = true

	if s.logger != nil {
		s.logger.spinner = s
		s.logger.loader = nil
	}

	s.mu.Unlock()

	go func() {
		ticker := time.NewTicker(time.Second / time.Duration(s.tps))
		defer ticker.Stop()
		for {
			select {
			case <-s.stop:
				clearSpinnerLine()
				close(s.done)
				s.mu.Lock()
				s.running = false
				s.paused = false
				if s.logger != nil {
					s.logger.spinner = nil
				}
				s.mu.Unlock()
				return
			case <-ticker.C:
				s.render()
			}
		}
	}()
}

func (s *Spinner) Stop() {
	s.mu.Lock()
	if !s.running {
		s.mu.Unlock()
		return
	}

	close(s.stop)
	var done chan struct{} = s.done
	s.mu.Unlock()

	<-done
	clearSpinnerLine()
}

func (s *Spinner) render() {
	s.mu.Lock()
	if !s.running || s.paused {
		s.mu.Unlock()
		return
	}

	var frame rune = s.frames[s.tick%len(s.frames)]
	s.tick++
	var logger *Logger = s.logger
	s.mu.Unlock()

	clearSpinnerLine()
	if logger != nil {
		fmt.Printf("\r%s%s %s", logger.spinnerPrefix(), string(frame), s.message)
	} else {
		fmt.Printf("\r%s %s", string(frame), s.message)
	}
}

func (s *Spinner) isRunning() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.running
}

func clearSpinnerLine() {
	fmt.Print("\r\033[2K")
}

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}
