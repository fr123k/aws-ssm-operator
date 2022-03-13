# Accept the Go version for the image to be set as a build argument.
ARG GO_VERSION=1.17.6

# Second stage: build the executable
FROM --platform=linux/amd64 golang:${GO_VERSION}-alpine AS builder

# Create the user and group files that will be used in the running container to
# run the process an unprivileged user.
RUN mkdir /user && \
    echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
    echo 'nobody:x:65534:' > /user/group

# Import the Certificate-Authority certificates for the app to be able to send
# requests to HTTPS endpoints.
RUN apk add --no-cache ca-certificates

# Accept the version of the app that will be injected into the compiled
# executable.
ARG APP_VERSION=undefined

# Set the environment variables for the build command.
ENV CGO_ENABLED=0 GOFLAGS=-mod=vendor

# Set the working directory outside $GOPATH to enable the support for modules.
WORKDIR /app

COPY go.mod ./
COPY go.sum ./
COPY vendor ./vendor

# Import the code from the first stage.
COPY cmd ./cmd
COPY pkg ./pkg

# inject the version as a global variable.RUN
RUN --mount=type=cache,id=gomod,target=/go/pkg/mod \
    --mount=type=cache,id=gobuild,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build \
    -installsuffix 'static' \
    -ldflags "-X main.Version=${APP_VERSION}" \
    -o /usr/local/bin/aws-ssm-operator -v \
    ./cmd/manager/main.go
# RUN go test -v -timeout 60s ./...

# Final stage: the running container
FROM --platform=linux/amd64 bash AS final

# This variable should be pass in the ci build pipeline
ARG APP_VERSION=undefined
ARG GIT_COMMIT=undefined
ARG BUILD_DATE=undefined

LABEL org.opencontainers.image.created="$BUILD_DATE"
LABEL org.opencontainers.image.description="aws-ssm-operator"
LABEL org.opencontainers.image.source="https://github.com/fr123k/aws-ssm-operator"
LABEL org.opencontainers.image.revision="$GIT_COMMIT"
LABEL org.opencontainers.image.version="$APP_VERSION"
LABEL go-version="${GO_VERSION}"

# Declare the port on which the application will be run.
EXPOSE 8080

# Import the user and group files.
COPY --from=builder /user/group /user/passwd /etc/

# Import the Certificate-Authority certificates for enabling HTTPS.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Import the compiled executable from the second stage.
COPY --from=builder /usr/local/bin/aws-ssm-operator /usr/local/bin/aws-ssm-operator

# Run the container as an unprivileged user.
USER nobody:nobody
