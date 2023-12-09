#cuerator: close({
	// Inputs is the result of unifying the values from theConfigMaps and Secrets
	// defined as inputs in the Installation resource.
	inputs: *{} | _
})
