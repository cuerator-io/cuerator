package list

import (
	"github.com/dogmatiq/dyad"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// Of is a list of T.
type Of[T any] struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []T `json:"items"`
}

// DeepCopyObject returns a deep clone of l.
func (l *Of[T]) DeepCopyObject() runtime.Object {
	return dyad.Clone(l)
}
