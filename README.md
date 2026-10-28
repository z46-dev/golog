# golog

A simple golang logging module that adds fancy headers and timestamps

## Spinners and loaders

```go
log := golog.New()
log.Spinner("Connecting...", golog.SpinnerDots, 10).
	Start().
	SetMessage("Downloading...").
	Stop()

log.Loader("Downloading...", golog.LoaderBlocks, 10).
	Start().
	SetProgress(0.5). // 0 to 1
	SetMessage("Installing...").
	Stop()
```

Loader styles: `LoaderBar` (classic), `LoaderBlocks`, `LoaderDots`, and `LoaderASCII`.
