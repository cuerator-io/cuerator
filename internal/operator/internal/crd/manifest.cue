package crd

import apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"

#CustomResourceDefinition: apiextensions.#CustomResourceDefinition & {
	apiVersion: "apiextensions.k8s.io/v1"
	kind:       "CustomResourceDefinition"
	metadata: {
		name: "\(spec.names.plural).\(spec.group)"
	}
	spec: {
		group: "cuerator.io"
		scope: *"Namespaced" | "Cluster"
		names: {
			categories: ["cuerator"]
		}
		versions: [#Version, ...#Version]
	}
}

#Version: {
	served:  *true | _
	storage: *true | _
	subresources: status: {}
	schema: openAPIV3Schema: {
		type: "object"
		required: ["spec"]
		properties: {
			spec:   #Spec
			status: #Status
		}
	}
}

#Spec: {
	type: "object"
	properties: {}
}

#Status: {
	type: "object"
	properties: {
		conditions: {
			"x-kubernetes-list-type": "map"
			"x-kubernetes-list-map-keys": ["type"]
			description: "List of conditions that indicate the status of object."
			type:        "array"
			items: {
				type: "object"
				required: [
					"status",
					"type",
				]
				properties: {
					type: {
						description: "Type of the condition."
						type:        "string"
					}
					status: {
						description: "Status of the condition."
						type:        "string"
						enum: ["Unknown", "True", "False"]
					}
					reason: {
						description: "A machine-readable explanation for the condition's last transition."
						type:        "string"
					}
					message: {
						description: "A human-readable description that complements the reason."
						type:        "string"
					}
					observedGeneration: {
						description: "The generation of the DNS-SD resource that was known to the controller when this condition was set."
						type:        "integer"
						format:      "int64"
					}
					lastTransitionTime: {
						description: "The time at which this condition was last changed."
						type:        "string"
						format:      "date-time"
					}
				}
			}
		}
	}
}
