package log

import (
	"fmt"
	"io"
	"os"
)

var (
	enabled bool
	out     io.Writer = os.Stderr
)

// Enable turns on debug logging. If w is nil, os.Stderr is used.
func Enable(w io.Writer) {
	enabled = true
	if w != nil {
		out = w
	}
}

// Debugf prints a formatted debug message in dim gray when enabled.
func Debugf(format string, args ...interface{}) {
	if !enabled {
		return
	}
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintf(out, "\x1b[90m%s\x1b[0m\n", msg)
}
