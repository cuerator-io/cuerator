package collection

import (
	"cuelang.org/go/cue/build"
	"cuelang.org/go/cue/load"
)

// Load returns the CUE instances in the given directory.
func Load(dir string) ([]*build.Instance, error) {
	cfg := &load.Config{
		Package:    "*",
		Dir:        dir,
		ModuleRoot: dir,
	}

	instances := load.Instances(nil, cfg)

	for _, inst := range instances {
		if inst.Err != nil {
			return nil, inst.Err
		}
	}

	return instances, nil
}
