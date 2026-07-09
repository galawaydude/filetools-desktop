package tool

import (
	"fmt"
	"sort"
	"sync"
)

// Registry is a thread-safe collection of Tools, indexed for fast lookup and
// grouped by category. It is the single source of truth the UI renders from.
type Registry struct {
	mu    sync.RWMutex
	order []string        // registration order, for stable listing
	byID  map[string]Tool // id -> tool
}

// NewRegistry returns an empty Registry.
func NewRegistry() *Registry {
	return &Registry{byID: make(map[string]Tool)}
}

// Register adds a tool. It panics on a duplicate id, because that is a
// programming error that should surface immediately at startup.
func (r *Registry) Register(t Tool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, dup := r.byID[t.ID()]; dup {
		panic(fmt.Sprintf("tool: duplicate registration for id %q", t.ID()))
	}
	r.byID[t.ID()] = t
	r.order = append(r.order, t.ID())
}

// Get returns the tool with the given id.
func (r *Registry) Get(id string) (Tool, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.byID[id]
	return t, ok
}

// All returns every registered tool in registration order.
func (r *Registry) All() []Tool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]Tool, 0, len(r.order))
	for _, id := range r.order {
		out = append(out, r.byID[id])
	}
	return out
}

// ByCategory returns the tools in a category, in registration order.
func (r *Registry) ByCategory(c Category) []Tool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []Tool
	for _, id := range r.order {
		if t := r.byID[id]; t.Category() == c {
			out = append(out, t)
		}
	}
	return out
}

// Counts returns the number of tools per category.
func (r *Registry) Counts() map[Category]int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	m := make(map[Category]int)
	for _, id := range r.order {
		m[r.byID[id].Category()]++
	}
	return m
}

// SortedCategories returns categories present in the registry, in display order.
func (r *Registry) SortedCategories() []Category {
	present := r.Counts()
	var out []Category
	for _, c := range Categories {
		if present[c] > 0 {
			out = append(out, c)
		}
	}
	// Any categories not in the canonical list get appended alphabetically.
	seen := map[Category]bool{}
	for _, c := range Categories {
		seen[c] = true
	}
	var extra []Category
	for c := range present {
		if !seen[c] {
			extra = append(extra, c)
		}
	}
	sort.Slice(extra, func(i, j int) bool { return extra[i] < extra[j] })
	return append(out, extra...)
}

// Default is the process-wide registry. Engine packages self-register into it
// from their init() functions (the database/sql driver pattern), so the only
// thing main needs to do is blank-import those packages.
var Default = NewRegistry()

// Register adds a tool to the Default registry.
func Register(t Tool) { Default.Register(t) }
