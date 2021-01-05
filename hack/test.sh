#!/bin/bash

GOBIN="$(go env GOBIN)"
gobin="${GOBIN:-$(go env GOPATH)/bin}"

echo "${gobin}"

${gobin}/deepcopy-gen \
    --input-dirs ../pkg/apis/kensho.ai/v1 \
    -O zz_generated.deepcopy \
    --bounding-dirs kensho.ai/v1 \
    --go-header-file ../hack/boilerplate.go.txt

${gobin}/client-gen \
    --clientset-name "${CLIENTSET_NAME_VERSIONED:-versioned}" \
    --input-base "" \
    --input "$(codegen::join , "${FQ_APIS[@]}")" \
    --output-package "${OUTPUT_PKG}/${CLIENTSET_PKG_NAME:-clientset}" \
    --go-header-file ../hack/boilerplate.go.txt