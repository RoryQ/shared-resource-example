ARG BASE
FROM golang:1.17 as builder

WORKDIR app

COPY go.mod app.go ./
# `skaffold debug` sets SKAFFOLD_GO_GCFLAGS to disable compiler optimizations
ARG SKAFFOLD_GO_GCFLAGS
RUN go build -gcflags="${SKAFFOLD_GO_GCFLAGS}" -o /app/protected

FROM gcr.io/distroless/base-debian11
COPY --from=builder /app/protected .

ENTRYPOINT ["/protected"]