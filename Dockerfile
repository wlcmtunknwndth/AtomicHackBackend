ARG GO_VERSION=1.22
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION} AS build

WORKDIR /src
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x

ARG TARGETARCH

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    CGO_ENABLED=0 GOARCH=$TARGETARCH go build -o /bin/server ./cmd

################################################################################

FROM alpine:latest AS final

RUN --mount=type=cache,target=/var/cache/apk \
    apk --update add \
        ca-certificates \
        tzdata \
        && \
        update-ca-certificates

RUN mkdir /bin/server

ARG UID=10001
USER root

COPY --from=build /bin/server /bin/server/
COPY  ./config/config.yaml /bin/server/

EXPOSE 63342

ENTRYPOINT [ "/bin/server/server" ]
