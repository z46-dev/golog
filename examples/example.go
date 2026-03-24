package main

import (
	"time"

	"github.com/z46-dev/golog"
)

func main() {
	var log *golog.Logger = golog.New().Prefix("[TEST]", golog.BoldBlue).Timestamp().Representation(false, true)
	var sLog *golog.Logger = log.SpawnChild().Prefix("[CHILD]", golog.BoldGreen)

	log.Info("Hello, world!")
	sLog.Info("Hello from the child logger!")
	log.Debug(log.Builder().C(golog.Red).A("Hello").R().A(", ").C(golog.Blue).A("world").R().A("!").B())
	log.Info(log.Builder().ThemeColors("Hello, world!", golog.RainbowTheme).B())
	log.Info(log.Builder().ThemeColors("Hello, world!", golog.BoldRainbowTheme).B())
	log.Info(log.Builder().ThemeColors("Hello, world!", golog.RainbowBackgroundTheme).B())
	log.Info(log.Builder().ThemeColors("USA", []golog.ColorCode{golog.BoldRed, golog.BoldWhite, golog.BoldBlue}).B())

	var spinner *golog.Spinner = log.Spinner("Loading...", golog.SpinnerRunner, 5)
	spinner.Start()

	log.Warning("One")
	time.Sleep(2 * time.Second)

	log.Error("Two")
	time.Sleep(2 * time.Second)

	log.Info("Three")
	time.Sleep(2 * time.Second)

	spinner.Stop()

	var loader *golog.Loader = log.Loader("Loading...", golog.LoaderBar, 10)

	loader.Start()

	log.Info("Thinking about it...")
	loader.SetProgress(0.1)

	time.Sleep(2 * time.Second)

	log.Info("Still thinking...")
	loader.SetProgress(0.5)

	time.Sleep(2 * time.Second)

	log.Info("Almost there...")
	loader.SetProgress(0.9)

	time.Sleep(2 * time.Second)

	loader.Stop()

	log = golog.New().Prefix("[TEST]", golog.BoldCyan).Timestamp().Representation(true, true).LogFile("./test.log", golog.LogFlushAlways)
	log.Info("This will be logged to a file!")
	log.Error("An error occurred!")

	log.Panic("Panic!")
}
