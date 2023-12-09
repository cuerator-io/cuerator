package installationcrd

import (
	"github.com/cuerator-io/cuerator/crd/internal/list"
	"github.com/dogmatiq/dyad"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// Version is the version of the API/CRDs.
const Version = "v1alpha1"

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
	list.Of[Installation] `json:",inline"`
}
