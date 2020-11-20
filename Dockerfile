# Build stage.
FROM golang:1.15.2 as builder

WORKDIR /workspace
COPY . .

RUN CGO_ENABLED=0 go build -o server cmd/main.go

# Runtime.
FROM gcr.io/distroless/static:nonroot

WORKDIR /
COPY --from=builder /workspace/server .
COPY static/ static/

USER nonroot:nonroot
ENTRYPOINT ["/server"]
