package cascader

import "fmt"

// Layer represents a named environment layer in the cascade chain.
type Layer struct {
	Name string
	Env  map[string]string
}

// Options controls cascade behaviour.
type Options struct {
	// Overwrite allows later layers to override earlier ones.
	// When false, the first definition of a key wins.
	Overwrite bool

	// Strict causes Cascade to return an error if any layer is nil.
	Strict bool
}

// DefaultOptions returns sensible defaults: first-wins, non-strict.
func DefaultOptions() Options {
	return Options{
		Overwrite: false,
		Strict:    false,
	}
}

// Result holds the merged environment and metadata about each key's origin.
type Result struct {
	// Env is the final merged map.
	Env map[string]string

	// Origins maps each key to the name of the layer it came from.
	Origins map[string]string
}

// Cascade merges layers in order according to opts.
// When Overwrite is false, the first layer that defines a key wins.
// When Overwrite is true, each subsequent layer can override previous values.
func Cascade(layers []Layer, opts Options) (*Result, error) {
	result := &Result{
		Env:     make(map[string]string),
		Origins: make(map[string]string),
	}

	for _, layer := range layers {
		if layer.Env == nil {
			if opts.Strict {
				return nil, fmt.Errorf("cascader: layer %q has nil env map", layer.Name)
			}
			continue
		}

		for k, v := range layer.Env {
			_, exists := result.Env[k]
			if !exists || opts.Overwrite {
				result.Env[k] = v
				result.Origins[k] = layer.Name
			}
		}
	}

	return result, nil
}
