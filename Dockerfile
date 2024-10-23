# Container image that runs your code
FROM golang:1.23-alpine as builder
WORKDIR /workspace

# copy Go Module manifests
COPY go.mod go.mod
COPY go.sum go.sum

# cache the dependancies
RUN go mod download

COPY ./ ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 \
    go build -a -ldflags '-extldflags "-static"' \
    -o SailPoint .

# copy to a thin image
FROM gcr.io/distroless/static:latest 
WORKDIR /
COPY --from=builder /workspace/SailPoint .
ENTRYPOINT ["/SailPoint"]