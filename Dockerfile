FROM --platform=$BUILDPLATFORM golang:1.25.6-alpine AS builder
WORKDIR /app

RUN apk add --no-cache git ca-certificates tzdata

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG TARGETOS
ARG TARGETARCH

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build \
    -a -installsuffix cgo \
    -ldflags="-w -s" \
    -trimpath \
    -o pm ./cmd/pm/

FROM gcr.io/distroless/static-debian13
WORKDIR /data
COPY --from=builder /app/pm /pm
EXPOSE 8080
ENTRYPOINT ["/pm"]