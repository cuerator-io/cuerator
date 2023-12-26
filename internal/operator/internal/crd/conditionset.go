package crd

import (
	"slices"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ConditionSet is a set of [metav1.Condition] objects within an object's status.
type ConditionSet []metav1.Condition

// Get returns the condition with the given type.
func (s *ConditionSet) Get(t string) metav1.Condition {
	index, ok := slices.BinarySearchFunc(
		*s,
		t,
		func(x metav1.Condition, t string) int {
			return strings.Compare(x.Type, t)
		},
	)

	if ok {
		return (*s)[index]
	}

	return metav1.Condition{
		Type:   t,
		Status: metav1.ConditionUnknown,
	}
}

// Set sets the given conditions.
func (s *ConditionSet) Set(conditions ...metav1.Condition) {
	for _, c := range conditions {
		index, ok := slices.BinarySearchFunc(
			*s,
			c.Type,
			func(x metav1.Condition, t string) int {
				return strings.Compare(x.Type, t)
			},
		)

		if !ok {
			*s = slices.Insert(*s, index, c)
		}

		x := &(*s)[index]

		if c.Status == x.Status {
			c.LastTransitionTime = x.LastTransitionTime
		} else {
			c.LastTransitionTime = metav1.Now()
		}

		*x = c
	}
}
