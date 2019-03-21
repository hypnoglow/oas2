# Accept the Go version for the image to be set as a build argument.
# Default to Go 1.12
ARG GO_VERSION=1.12

# First stage: build the executable.
FROM golang:${GO_VERSION}-alpine AS build

WORKDIR /src

COPY . .

RUN apk add --no-cache git

RUN go mod vendor && \
    CGO_ENABLED=0 go build -mod=vendor -o ./bin/oas-expand ./cmd/oas-expand

# Final stage: the running container.
FROM scratch AS final

COPY --from=build /src/bin/oas-expand /oas-expand
