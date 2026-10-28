package golog

import (
	"io"
	"os"
	"strings"
	"testing"
)

func TestAnimationMethodsReturnSelf(t *testing.T) {
	s := newSpinner("spinner", SpinnerLine, 1, nil, nil)
	if s.Start() != s || s.Start() != s || s.SetMessage("working") != s || s.Stop() != s || s.Stop() != s {
		t.Fatal("spinner method did not return its receiver")
	}

	l := newLoader("loader", LoaderBar, 1, nil, nil)
	if l.Start() != l || l.Start() != l || l.SetProgress(0.5) != l || l.SetMessage("working") != l || l.Stop() != l || l.Stop() != l {
		t.Fatal("loader method did not return its receiver")
	}
}

func TestSetMessageRendersImmediately(t *testing.T) {
	tests := []struct {
		name string
		run  func()
	}{
		{"spinner", func() {
			s := newSpinner("before", SpinnerLine, 1, nil, nil)
			s.running = true
			s.SetMessage("after")
		}},
		{"loader", func() {
			l := newLoader("before", LoaderBlocks, 1, nil, nil)
			l.running = true
			l.SetMessage("after")
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader, writer, err := os.Pipe()
			if err != nil {
				t.Fatal(err)
			}
			original := os.Stdout
			os.Stdout = writer
			defer func() { os.Stdout = original }()

			tt.run()
			writer.Close()
			output, err := io.ReadAll(reader)
			reader.Close()
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(string(output), "after") {
				t.Fatalf("updated message was not rendered: %q", output)
			}
		})
	}
}
