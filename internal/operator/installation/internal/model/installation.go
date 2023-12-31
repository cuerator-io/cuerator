package model

import (
	"github.com/cuerator-io/cuerator/internal/operator/internal/crd"
	"github.com/dogmatiq/dyad"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// APIVersion is the version of the CRD's API.
const APIVersion = "v1alpha1"

// Installation is the root of an "installation" resource.
type Installation struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              Spec   `json:"spec,omitempty"`
	Status            Status `json:"status,omitempty"`
}

// DeepCopyObject returns a deep clone of i.
func (s *Installation) DeepCopyObject() runtime.Object {
	return dyad.Clone(s)
}

// InstallationList is a list of [Installation] objects.
type InstallationList struct {
	crd.List[Installation] `json:",inline"`
}
