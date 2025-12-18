package golog

import "fmt"

func newBuilder() (builder *Builder) {
	builder = &Builder{
		message: "",
	}

	return
}

// Apply a color code to the message after this call
func (b *Builder) C(color ColorCode) (self *Builder) {
	b.message += string(color)
	self = b
	return
}

// Reset the color code after this call
func (b *Builder) R() (self *Builder) {
	b.message += string(Reset)
	self = b
	return
}

// Append a string to the message
func (b *Builder) A(str string) (self *Builder) {
	b.message += str
	self = b
	return
}

// Append a formatted string to the message
func (b *Builder) F(format string, args ...any) (self *Builder) {
	b.message += fmt.Sprintf(format, args...)
	self = b
	return
}

// Build the message and reset the builder
func (b *Builder) B() (message string) {
	message = b.message
	b.message = ""
	return
}

// Splits the text into characters and applies a color to each character
// Color is applied in a round-robin fashion (colors[idx % len(colors)])
func (b *Builder) ThemeColors(text string, colors []ColorCode) (self *Builder) {
	for i, c := range text {
		b.C(colors[i%len(colors)]).A(string(c))
	}

	b.R()

	self = b
	return
}
