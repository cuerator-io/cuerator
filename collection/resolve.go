package collection

import (
	"cuelang.org/go/cue"
	"cuelang.org/go/cue/build"
	"cuelang.org/go/cue/cuecontext"
)

// Resolve passes the given inputs to each of the given CUE instances and
// returns the resulting outputs.
func Resolve(
	inputs map[string]any,
	instances []*build.Instance,
) (map[string]map[string]any, error) {
	cc := cuecontext.New()

	base := cc.
		CompileBytes(baseCUE).
		FillPath(
			cue.ParsePath("#cuerator.inputs"),
			inputs,
		)

	// for i, inst := range instances {
	// 	if inst.PkgName == "cuerator" {
	// 		v := cc.BuildInstance(inst)
	// 		if v.Err() != nil {
	// 			return nil, v.Err()
	// 		}

	// 		base = base.Unify(v)
	// 		instances = slices.Delete(instances, i, i+1)
	// 		break
	// 	}
	// }

	outputs := map[string]map[string]any{}

	for _, inst := range instances {
		if len(inst.Files) == 0 {
			continue
		}

		v := cc.BuildInstance(inst, cue.Scope(base))
		if v.Err() != nil {
			return nil, v.Err()
		}

		v = v.Unify(base)

		var x map[string]any
		if err := v.Decode(&x); err != nil {
			return nil, err
		}

		outputs[inst.PkgName] = x
	}

	return outputs, nil
}
