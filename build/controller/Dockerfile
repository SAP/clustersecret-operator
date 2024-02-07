### build go executable
FROM --platform=$BUILDPLATFORM golang:1.22.0 as build

WORKDIR /go/src

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download -x

# Caveat: when doing changes here, double-maintain them in .github/workflows
COPY Makefile Makefile
COPY hack/ hack/
COPY pkg/ pkg/
COPY cmd/ cmd/
COPY internal/ internal/
COPY test/ test/

RUN make test-controller
ARG TARGETOS TARGETARCH
RUN echo ${TARGETOS}/${TARGETARCH} && make build-controller

### final image
FROM scratch

ENTRYPOINT ["/app/bin/controller"]

COPY --from=build /go/src/bin/controller /app/bin/controller
