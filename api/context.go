package api

import (
	"errors"

	"golang.org/x/net/context"
)

// unexported key prevents collisions
type key int

const (
	specKey key = iota
)

var (
	errNoSpecFromContext = errors.New("api: Context missing a Spec")
)

// withSpec returns a copy of ctx that stores the given Spec.
func withSpec(ctx context.Context, spec *Spec) context.Context {
	return context.WithValue(ctx, specKey, spec)
}

// specFromContext returns the Spec from the ctx.
func specFromContext(ctx context.Context) (*Spec, error) {
	spec, ok := ctx.Value(specKey).(*Spec)
	if !ok {
		return nil, errNoSpecFromContext
	}
	return spec, nil
}
