FROM golang:1.24-alpine AS builder

WORKDIR /src

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY app ./app
COPY internal ./internal

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -trimpath \
    -ldflags="-s -w -X github.com/deanlofts/devops-roundrobin/internal/build.Version=local -X github.com/deanlofts/devops-roundrobin/internal/build.Commit=docker -X github.com/deanlofts/devops-roundrobin/internal/build.Date=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
    -o /out/app ./app

FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /
COPY --from=builder /out/app /app

EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/app"]

