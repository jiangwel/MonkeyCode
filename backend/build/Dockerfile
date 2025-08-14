FROM --platform=$BUILDPLATFORM golang:1.24-alpine AS builder

RUN apk add --no-cache git

WORKDIR /src
ENV CGO_ENABLED=0

# Configure git to use the GitHub token for private repositories
ARG GITHUB_TOKEN
RUN if [ -n "$GITHUB_TOKEN" ]; then \
        git config --global url."https://$GITHUB_TOKEN:x-oauth-basic@github.com/".insteadOf "https://github.com/"; \
    fi

COPY go.* .
ARG GOMODCACHE
RUN --mount=type=cache,target=${GOMODCACHE} \
    go mod download

ARG TARGETOS TARGETARCH GOCACHE
ARG VERSION
ARG BUILD_TIME
ARG GIT_COMMIT
RUN --mount=type=bind,target=. \
--mount=type=cache,target=${GOMODCACHE} \
--mount=type=cache,target=${GOCACHE} \
GOOS=$TARGETOS GOARCH=$TARGETARCH \
go build \
-ldflags "-w -s -X 'github.com/chaitin/MonkeyCode/backend/pkg/version.Version=${VERSION}' -X 'github.com/chaitin/MonkeyCode/backend/pkg/version.BuildTime=${BUILD_TIME}' -X 'github.com/chaitin/MonkeyCode/backend/pkg/version.GitCommit=${GIT_COMMIT}'" \
-o /out/main \
cmd/server/main.go cmd/server/wire_gen.go 

FROM alpine:3.22.1 AS binary

WORKDIR /app

ADD migration ./migration
ADD assets/vsix ./assets/vsix

COPY --from=builder /out/main /app/main

CMD [ "./main" ]