ARG BASE
FROM golang:1.17 as builder

WORKDIR app

COPY go.mod go.sum ./
COPY vendor vendor
COPY app.go ./

# `skaffold debug` sets SKAFFOLD_GO_GCFLAGS to disable compiler optimizations
ARG SKAFFOLD_GO_GCFLAGS
RUN go build -gcflags="${SKAFFOLD_GO_GCFLAGS}" -o /app/takelock

FROM gcr.io/distroless/base-debian11
COPY --from=builder /app/takelock .

ENTRYPOINT ["/takelock"]
