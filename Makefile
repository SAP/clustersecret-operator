BUILDOS := $(shell go env GOOS)
BUILDARCH := $(shell go env GOARCH)
TESTOS ?= $(BUILDOS)
TESTARCH ?= $(BUILDARCH)
TARGETOS ?= $(TESTOS)
TARGETARCH ?= $(TESTARCH)

# build (used by Dockerfile)
.PHONY: build
build: build-controller build-webhook

.PHONY: build-controller
build-controller:
	@CGO_ENABLED=0 GOOS=$(TARGETOS) GOARCH=$(TARGETARCH) go build -o ./bin/controller ./cmd/controller

.PHONY: build-webhook
build-webhook:
	@CGO_ENABLED=0 GOOS=$(TARGETOS) GOARCH=$(TARGETARCH) go build -o ./bin/webhook ./cmd/webhook

# build and install executables
.PHONY: install
install: install-controller install-webhook

.PHONY: install-controller
install-controller:
	@CGO_ENABLED=0 GOOS=$(TARGETOS) GOARCH=$(TARGETARCH) go install ./cmd/controller

.PHONY: install-webhook
install-webhook:
	@CGO_ENABLED=0 GOOS=$(TARGETOS) GOARCH=$(TARGETARCH) go install ./cmd/webhook

# run unit tests
.PHONY: test
test: test-controller test-webhook

.PHONY: test-controller
test-controller:
	@go test -v ./internal/controller

.PHONY: test-webhook
test-webhook:
	@go test -v ./internal/admission

# generate code
.PHONY: generate
generate:
	@hack/update-codegen.sh --parallel

# format code
.PHONY: format
format:
	@go fmt ./cmd/... ./internal/... ./pkg/apis/... ./test...

# prepare local developement environment
.PHONY: local-generate
local-generate:
	@hack/generate-local.sh

.PHONY: local-setup
local-setup:
	@hack/setup-local.sh

# Set the year for SPDX header updates (default: current year)
YEAR ?= $(shell date +%Y)

.PHONY: update-header-year
update-header-year:
    # Go + TXT
	@find . -type f \( -name "*.go" -o -name "*.txt" \) -exec sed -i \
	's/^SPDX-FileCopyrightText: [0-9]\{4\}\( SAP SE or an SAP affiliate company and [^"]\+ contributors\)/SPDX-FileCopyrightText: $(YEAR)\1/' {} +

    # TOML
	@find . -type f -name "*.toml" -exec sed -i \
	's/^SPDX-FileCopyrightText = "[0-9]\{4\}\( SAP SE or an SAP affiliate company and [^"]\+ contributors\)"/SPDX-FileCopyrightText = "$(YEAR)\1"/' {} +