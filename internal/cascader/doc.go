// Package cascader merges multiple named environment layers into a single
// resolved map, tracking the origin layer for every key.
//
// Layers are processed in the order they are supplied. By default the first
// layer to define a key wins (first-wins semantics). Setting Options.Overwrite
// to true switches to last-wins semantics, where each subsequent layer may
// replace values set by earlier ones.
//
// Example usage:
//
//	layers := []cascader.Layer{
//		{Name: "defaults", Env: defaultEnv},
//		{Name: "staging",  Env: stagingEnv},
//		{Name: "local",    Env: localEnv},
//	}
//	result, err := cascader.Cascade(layers, cascader.DefaultOptions())
package cascader
