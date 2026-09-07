# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /cv-server ./cmd/server

# Runtime stage
FROM alpine:3.20

RUN apk add --no-cache chromium ca-certificates tzdata

ENV CHROME_PATH=/usr/bin/chromium-browser

WORKDIR /app

COPY --from=builder /cv-server /app/cv-server

RUN mkdir -p /app/data

EXPOSE 8080

CMD ["/app/cv-server"]
