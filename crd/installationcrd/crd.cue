import "github.com/cuerator-io/cuerator/crd/internal/manifest"

manifest.#CustomResourceDefinition

spec: {
	names: {
		kind:   "Installation"
		plural: "installations"
		shortNames: ["install", "inst"]
	}
	versions: [
		{
			name: "v1alpha1"
			schema: openAPIV3Schema: {
				properties:
				{
					spec: {
						required: [
							"image",
							"versionConstraint",
						]
						properties: {
							image: {
								description: "The name of the image containing the Cuerator Collection."
								type:        "string"
							}
							versionConstraint: {
								description: "A semantic version constraint that determines which version of the image to use."
								type:        "string"
							}
							inputs: {
								description: "A list of references to Kubernetes resources that are used as inputs to the Cuerator Collection."
								type:        "array"
								items: {
									type: "object"
									required: [
										"name",
										"kind",
									]
									properties: {
										kind: {
											description: "The kind of resource from which input values are read."
											type:        "string"
											enum: ["Secret", "ConfigMap"]
										}
										name: {
											description: "The name of the input resource."
											type:        "string"
										}
									}
								}
							}
						}
					}
					status: {
						properties: {
							tag: {
								description: "The image tag that the version constraint resolves to."
								type:        "object"
								required: [
									"name",
									"version",
								]
								properties: {
									name: {
										description: "The tag name."
										type:        "string"
									}
									version: {
										description: "The normalized semantic version."
										type:        "string"
									}
									digest: {
										description: "The image digest for the tag."
										type:        "string"
									}
								}
							}
						}
					}
				}
			}
			additionalPrinterColumns: [
				{
					name:        "Image"
					description: "The name of the image containing the Cuerator Collection."
					type:        "string"
					jsonPath:    ".spec.image"
				},
				{
					name:        "Constraint"
					description: "The semantic version constraint that determines which version of the image to use."
					type:        "string"
					jsonPath:    ".spec.versionConstraint"
				},
				{
					name:        "Tag"
					description: "The name of the image tag the the version constraint resolves to."
					type:        "string"
					jsonPath:    ".status.tag.name"
				},
				{
					name:        "Age"
					description: "The time at which the installation was created."
					type:        "date"
					jsonPath:    ".metadata.creationTimestamp"
				},
			]
		},
	]
}
