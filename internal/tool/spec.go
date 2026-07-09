package tool

import "context"

// RunFunc is the body of a tool — the actual conversion work.
type RunFunc func(ctx context.Context, req Request, p Progress) (Result, error)

// Config declares a tool. Using Define(Config) instead of hand-writing the
// eight-method Tool interface for every operation keeps each feature to a
// single, readable block and removes boilerplate.
type Config struct {
	ID          string
	Name        string
	Description string
	Category    Category
	Input       InputKind
	Extensions  []string
	Options     []Option
	Run         RunFunc
}

// Define builds a Tool from a Config.
func Define(c Config) Tool { return &spec{c} }

// spec adapts a Config to the Tool interface.
type spec struct{ c Config }

func (s *spec) ID() string           { return s.c.ID }
func (s *spec) Name() string         { return s.c.Name }
func (s *spec) Description() string  { return s.c.Description }
func (s *spec) Category() Category   { return s.c.Category }
func (s *spec) InputKind() InputKind { return s.c.Input }
func (s *spec) Extensions() []string { return s.c.Extensions }
func (s *spec) Options() []Option    { return s.c.Options }

func (s *spec) Run(ctx context.Context, req Request, p Progress) (Result, error) {
	if p == nil {
		p = NopProgress
	}
	return s.c.Run(ctx, req, p)
}
