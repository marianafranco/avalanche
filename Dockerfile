# Build the manager binary
FROM golang:1.17 as builder

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum

ENV GOPROXY direct

# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY cmd/ cmd/
COPY metrics/ metrics/
COPY pkg/ pkg/
COPY scripts/ scripts/

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o avalanche cmd/avalanche.go
# RUN GO111MODULE=on /Users/marfram/go/bin/promu build --prefix /Users/marfram/workplace/avalanche

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/avalanche .
USER 65532:65532

ENTRYPOINT ["/avalanche"]