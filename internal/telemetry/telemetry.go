package telemetry

import (
	"context"
	"os"
	"time"

	lf "github.com/henomis/langfuse-go"
	lfmodel "github.com/henomis/langfuse-go/model"

	"command/internal/log"
)

// Tracker records LLM usage.
// Methods are no-op when telemetry is disabled.
type Tracker interface {
	Generation(prompt, model string, commands []string)
}

// NewFromEnv creates a Langfuse tracker using environment variables.
func NewFromEnv(ctx context.Context, debug bool) Tracker {
	if os.Getenv("LANGFUSE_PUBLIC_KEY") == "" || os.Getenv("LANGFUSE_SECRET_KEY") == "" {
		return noop{}
	}
	return &langfuseTracker{lf: lf.New(ctx), debug: debug}
}

type langfuseTracker struct {
	lf    *lf.Langfuse
	debug bool
}

func (l *langfuseTracker) Generation(prompt, model string, commands []string) {
	if l.lf == nil {
		return
	}
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		if !l.debug {
			devNull, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
			if err == nil {
				orig := os.Stdout
				os.Stdout = devNull
				defer func() {
					os.Stdout = orig
					_ = devNull.Close()
				}()
			}
		}

		if _, err := l.lf.Generation(&lfmodel.Generation{
			Name:   "cmd-generation",
			Model:  model,
			Input:  map[string]any{"prompt": prompt},
			Output: commands,
		}, nil); err == nil {
			l.lf.Flush(ctx)
		} else if l.debug {
			log.Debugf("telemetry error: %v", err)
		}
	}()
}

type noop struct{}

func (noop) Generation(string, string, []string) {}

// Disabled returns a tracker that does nothing.
func Disabled() Tracker { return noop{} }
