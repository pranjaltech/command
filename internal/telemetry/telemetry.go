package telemetry

import (
	"context"

	lf "github.com/henomis/langfuse-go"
	lfmodel "github.com/henomis/langfuse-go/model"
)

// Tracker records LLM usage.
// Methods are no-op when telemetry is disabled.
type Tracker interface {
	Generation(prompt, model string, commands []string)
}

// NewFromEnv creates a Langfuse tracker using environment variables.
func NewFromEnv(ctx context.Context) Tracker {
	return &langfuseTracker{lf: lf.New(ctx)}
}

type langfuseTracker struct{ lf *lf.Langfuse }

func (l *langfuseTracker) Generation(prompt, model string, commands []string) {
	if l.lf == nil {
		return
	}
	_, _ = l.lf.Generation(&lfmodel.Generation{
		Name:   "cmd-generation",
		Model:  model,
		Input:  map[string]any{"prompt": prompt},
		Output: commands,
	}, nil)
	l.lf.Flush(context.Background())
}

type noop struct{}

func (noop) Generation(string, string, []string) {}

// Disabled returns a tracker that does nothing.
func Disabled() Tracker { return noop{} }
