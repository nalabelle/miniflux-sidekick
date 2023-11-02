VERSION --pass-args 0.7
FROM golang:1.21-alpine
WORKDIR /build
ARG --global REGISTRY="ghcr.io"

deps:
  COPY go.mod go.sum ./
  RUN go mod download
  SAVE ARTIFACT go.mod AS LOCAL go.mod
  SAVE ARTIFACT go.sum AS LOCAL go.sum

build:
  FROM +deps
  COPY --dir cmd filter rules ./
  RUN go build -o build/miniflux-sidekick ./cmd/api
  SAVE ARTIFACT build/miniflux-sidekick /miniflux-sidekick

test:
  FROM +deps
  COPY --dir cmd filter rules ./
  RUN CGO_ENABLED=0 go test ./...

docker:
  ARG EARTHLY_GIT_PROJECT_NAME
  FROM alpine:latest
  RUN apk --update upgrade && \
      apk add curl ca-certificates && \
      update-ca-certificates && \
      rm -rf /var/cache/apk/*
  COPY +build/miniflux-sidekick /usr/local/bin/miniflux-sidekick
  RUN chmod +x /usr/local/bin/miniflux-sidekick

  # Run the image as a non-root user
  RUN adduser --uid 2000 -D user
  USER user

  CMD miniflux-sidekick

  ARG EARTHLY_GIT_ORIGIN_URL
  ARG EARTHLY_GIT_HASH
  ARG EARTHLY_GIT_COMMIT_TIMESTAMP
  ARG EARTHLY_TARGET_TAG_DOCKER
  LABEL "git.origin"="${EARTHLY_GIT_ORIGIN_URL}"
  LABEL "git.hash"="${EARTHLY_GIT_HASH}"
  LABEL "build.timestamp"="${EARTHLY_GIT_COMMIT_TIMESTAMP}"
  ARG VERSION="0.1.${EARTHLY_GIT_COMMIT_TIMESTAMP}"
  IF [ "$EARTHLY_TARGET_TAG_DOCKER" == "main" ]
    SAVE IMAGE --push ${REGISTRY}/${EARTHLY_GIT_PROJECT_NAME}:${VERSION}
    SAVE IMAGE --push ${REGISTRY}/${EARTHLY_GIT_PROJECT_NAME}:latest
  END

all:
  BUILD +build
  BUILD +test
  BUILD +docker
