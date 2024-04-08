FROM golang:1.21-alpine3.18 AS install
RUN apk add --no-cache git make bash
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

FROM install AS build
WORKDIR /src
COPY cmd ./cmd
COPY internal ./internal
RUN --mount=type=cache,target=/root/.cache/go-build \
    go build -o /bin/api cmd/api/main.go

FROM ubuntu AS api
WORKDIR /
COPY --from=build /bin/api /bin/api
CMD ["/bin/api"]
