package installation

import (
	"slices"
	"sync"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
)

type references struct {
	m          sync.RWMutex
	byReferrer map[types.NamespacedName]map[referent]struct{}
	byReferent map[referent]map[types.NamespacedName]struct{}
}

type referent struct {
	schema.GroupVersionKind
	types.NamespacedName
}

func (r *references) Referrers(to referent) []types.NamespacedName {
	r.m.RLock()
	defer r.m.RUnlock()

	names := make([]types.NamespacedName, 0, len(r.byReferent[to]))

	for n := range r.byReferent[to] {
		names = append(names, n)
	}

	return names
}

func (r *references) Update(from types.NamespacedName, to []referent) {
	r.m.Lock()
	defer r.m.Unlock()

	if r.byReferrer == nil {
		r.byReferrer = map[types.NamespacedName]map[referent]struct{}{}
		r.byReferent = map[referent]map[types.NamespacedName]struct{}{}
	}

	referents := r.byReferrer[from]

	for _, ref := range to {
		if _, ok := referents[ref]; ok {
			continue
		}

		if referents == nil {
			referents = map[referent]struct{}{}
			r.byReferrer[from] = referents
		}

		referrers := r.byReferent[ref]
		if referrers == nil {
			referrers = map[types.NamespacedName]struct{}{}
			r.byReferent[ref] = referrers
		}

		referents[ref] = struct{}{}
		referrers[from] = struct{}{}
	}

	for ref := range referents {
		if slices.Contains(to, ref) {
			continue
		}

		referrers := r.byReferent[ref]

		delete(referents, ref)
		if len(referents) == 0 {
			delete(r.byReferrer, from)
		}

		delete(referrers, from)
		if len(referrers) == 0 {
			delete(r.byReferent, ref)
		}
	}
}
