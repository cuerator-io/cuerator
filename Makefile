DOCKER_REPO = ghcr.io/cuerator-io/cuerator
DOCKER_PLATFORMS += linux/amd64
DOCKER_PLATFORMS += linux/arm64

-include .makefiles/Makefile
-include .makefiles/pkg/go/v1/Makefile
-include .makefiles/pkg/go/v1/with-ferrite.mk
-include .makefiles/pkg/docker/v1/Makefile

CRD_CUE_FILES += $(shell find . -name 'crd.cue')
CRD_YAML_FILES += $(foreach f,$(CRD_CUE_FILES:.cue=.yaml),$(if $(findstring /_,/$f),,$f))

GENERATED_FILES += $(CRD_YAML_FILES)

.PHONY: run
run: $(GO_DEBUG_DIR)/cuerator
	$< $(args)

.PHONY: run-operator
run-operator: $(GO_DEBUG_DIR)/cuerator
	kubectl apply $(foreach f,$(CRD_YAML_FILES),-f $f)
	DEBUG=true $< run $(args)

%/crd.yaml: %/crd.cue
	cue fmt $<
	cue def --out yaml --force --outfile $@ $<

.makefiles/%:
	@curl -sfL https://makefiles.dev/v1 | bash /dev/stdin "$@"
