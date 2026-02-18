FROM golang:1.25.6-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo \
    -ldflags="-w -s" -trimpath -o survey ./cmd/survey/

FROM gcr.io/distroless/static-debian13
COPY --from=builder /app/survey /survey
EXPOSE 8080
ENTRYPOINT ["/survey"]