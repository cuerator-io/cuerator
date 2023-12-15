import "github.com/cuerator-io/cuerator/internal/operator/internal/crd"

crd.#CustomResourceDefinition

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
				properties: {
					spec: {
						required: [
							"image",
							"versionConstraint",
						]
						properties: {
							image: {
								type:        "string"
								description: "The name of the image containing the Cuerator Collection."
							}
							versionConstraint: {
								type:        "string"
								description: "A semantic version constraint that determines which version of the image to use."
							}
							inputs: {
								type: "array"
								items: {
									type: "object"
									required: ["path"]
									properties: {
										path: {
											type:        "string"
											description: "The path within the Cuerator inputs struct at which the value is merged."
										}
										value: {
											description:                            "The value, expressed as a YAML value (not a CUE expression)."
											"x-kubernetes-preserve-unknown-fields": true
										}
										valueFrom: {
											type:        "object"
											description: "Obtain the input value from another Kubernetes resource."
											properties: {
												configMapKeyRef: {
													type:        "object"
													description: "Obtain the input value from a specific key of a ConfigMap."
													required: ["name", "key"]
													properties: {
														name: {
															type:        "string"
															description: "The name of the ConfigMap."
														}
														key: {
															type:        "string"
															description: "The key within the ConfigMap."
														}
														optional: {
															type:        "boolean"
															description: "Indicates whether to ignore this reference if the ConfigMap or key does not exist; otherwise, an error occurs."
														}
													}
												}
												secretKeyRef: {
													type:        "object"
													description: "Obtain the input value from a specific key of a Secret."
													required: ["name", "key"]
													properties: {
														name: {
															type:        "string"
															description: "The name of the Secret."
														}
														key: {
															type:        "string"
															description: "The key within the Secret."
														}
														optional: {
															type:        "boolean"
															description: "Indicates whether to ignore this reference if the Secret or key does not exist; otherwise, an error occurs."
														}
													}
												}
											}
											value: {
												description: "The value of the input."
											}
											valueFrom: {
												description: "Obtain the input value from another Kubernetes resource."
												type:        "object"
												properties: {
													configMapKeyRef: {
														description: "Obtain the input value from a specific key of a ConfigMap."
														type:        "object"
														required: ["name", "key"]
														properties: {
															name: {
																description: "The name of the ConfigMap."
																type:        "string"
															}
															key: {
																description: "The key of the ConfigMap entry."
																type:        "string"
															}
															optional: {
																description: "If true, the input value is optional. If the ConfigMap or key does not exist, the input value is not set."
																type:        "boolean"
															}
														}
													}
													secretKeyRef: {
														description: "Obtain the input value from a specific key of a ConfigMap."
														type:        "object"
														required: ["name", "key"]
														properties: {
															name: {
																description: "The name of the ConfigMap."
																type:        "string"
															}
															key: {
																description: "The key of the ConfigMap entry."
																type:        "string"
															}
															optional: {
																description: "If true, the input value is optional. If the ConfigMap or key does not exist, the input value is not set."
																type:        "boolean"
															}
														}
													}
												}
											}
										}
									}
								}
							}
							inputsFrom: {
								description: "A list of references to Kubernetes resources that are used as inputs to the Cuerator Collection."
								type:        "array"
								items: {
									type: "object"
									properties: {
										path: {
											type:        "string"
											description: "The path at which the values are placed within the Cuerator inputs, defaults to the root of the inputs struct."
										}
										configMapRef: {
											type:        "object"
											description: "Obtain input values from a ConfigMap."
											required: ["name"]
											properties: {
												name: {
													type:        "string"
													description: "The name of the ConfigMap."
												}
												optional: {
													type:        "boolean"
													description: "Indicates whether to ignore this reference if the ConfigMap does not exist; otherwise, an error occurs."
												}
											}
										}
										secretRef: {
											type:        "object"
											description: "Obtain input values from a Secret."
											required: ["name"]
											properties: {
												name: {
													type:        "string"
													description: "The name of the Secret."
												}
												optional: {
													type:        "boolean"
													description: "Indicates whether to ignore this reference if the Secret does not exist; otherwise, an error occurs."
												}
											}
										}
									}
								}
							}
						}
					}
					status: {
						type: "object"
						properties: {
							tag: {
								description: "The image tag that the version constraint resolves to."
								type:        "object"
								required: [
									"image",
									"name",
									"digest",
									"normalizedVersion",
								]
								properties: {
									image: {
										description: "The name of the image that was used to resolve the tag."
										type:        "string"
									}
									name: {
										description: "The tag name."
										type:        "string"
									}
									digest: {
										description: "The image digest for the tag."
										type:        "string"
									}
									normalizedVersion: {
										description: "The normalized semantic version."
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
