package provider

import (
	"context"
	"fmt"
)

// SourceFactory creates a Source provider from opaque options.
type SourceFactory func(ctx context.Context, opts map[string]any) (Source, error)

// DestinationFactory creates a Destination provider from opaque options.
type DestinationFactory func(ctx context.Context, opts map[string]any) (Destination, error)

// Registration holds factory functions for a provider type.
type Registration struct {
	Name            Type
	SourceFactory   SourceFactory
	DestFactory     DestinationFactory
}

// registry holds all registered providers.
//nolint:gochecknoglobals
var registry = map[Type]Registration{}

// Register adds a provider to the registry.
func Register(r Registration) {
	registry[r.Name] = r
}

// Get returns the registration for the given provider type.
func Get(t Type) (Registration, error) {
	r, ok := registry[t]
	if !ok {
		return Registration{}, fmt.Errorf("unknown provider type: %s", t)
	}
	return r, nil
}

// CreateSource creates a Source for the given provider type.
func CreateSource(ctx context.Context, t Type, opts map[string]any) (Source, error) {
	r, err := Get(t)
	if err != nil {
		return nil, err
	}
	if r.SourceFactory == nil {
		return nil, fmt.Errorf("provider %s does not implement Source", t)
	}
	return r.SourceFactory(ctx, opts)
}

// CreateDestination creates a Destination for the given provider type.
func CreateDestination(ctx context.Context, t Type, opts map[string]any) (Destination, error) {
	r, err := Get(t)
	if err != nil {
		return nil, err
	}
	if r.DestFactory == nil {
		return nil, fmt.Errorf("provider %s does not implement Destination", t)
	}
	return r.DestFactory(ctx, opts)
}

// RegisteredTypes returns all registered provider types.
func RegisteredTypes() []Type {
	types := make([]Type, 0, len(registry))
	for t := range registry {
		types = append(types, t)
	}
	return types
}
