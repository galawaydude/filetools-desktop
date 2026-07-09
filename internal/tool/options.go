package tool

import (
	"strconv"
	"strings"
)

// Options holds user-selected option values keyed by Option.Key. Values are
// stored as strings (as produced by the UI) and parsed on demand by helpers.
type Options map[string]string

// String returns the raw value for key, or "" if absent.
func (o Options) String(key string) string {
	return o[key]
}

// StringOr returns the value for key, or def if absent/empty.
func (o Options) StringOr(key, def string) string {
	if v, ok := o[key]; ok && v != "" {
		return v
	}
	return def
}

// Int parses the value for key, falling back to def on absence or parse error.
func (o Options) Int(key string, def int) int {
	if v, ok := o[key]; ok && v != "" {
		if n, err := strconv.Atoi(strings.TrimSpace(v)); err == nil {
			return n
		}
	}
	return def
}

// Float parses the value for key, falling back to def on absence or parse error.
func (o Options) Float(key string, def float64) float64 {
	if v, ok := o[key]; ok && v != "" {
		if f, err := strconv.ParseFloat(strings.TrimSpace(v), 64); err == nil {
			return f
		}
	}
	return def
}

// Bool interprets the value for key as a boolean ("true"/"1"/"yes"/"on").
func (o Options) Bool(key string) bool {
	switch strings.ToLower(strings.TrimSpace(o[key])) {
	case "true", "1", "yes", "on":
		return true
	}
	return false
}

// Defaults builds an Options map pre-filled with each option's Default value.
func Defaults(opts []Option) Options {
	m := make(Options, len(opts))
	for _, opt := range opts {
		m[opt.Key] = opt.Default
	}
	return m
}
