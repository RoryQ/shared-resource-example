FROM golang:1.17-bullseye as build

ENV GH_USER=roryq
ENV GH_REPO=wrench
ENV GH_VERSION=1.5.1

RUN wget "https://github.com/$GH_USER/$GH_REPO/releases/download/v${GH_VERSION}/${GH_REPO}_${GH_VERSION}_Linux_arm64.tar.gz" \
    -O "$GH_REPO-$GH_BRANCH.tar.gz" && \
    tar -xzvf ./"$GH_REPO-$GH_BRANCH.tar.gz" && \
    rm ./"$GH_REPO-$GH_BRANCH.tar.gz"

FROM gcr.io/distroless/base-debian11:debug
COPY --from=build go/wrench /bin/
COPY migrations migrations

ENTRYPOINT ["/busybox/sh", "-c", "wrench migrate up && wrench migrate history"]

