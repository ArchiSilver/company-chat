FROM golang:1.24-alpine AS builder
RUN apk add --no-cache git ca-certificates
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . ./

# Build server binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o /build/company-chat ./cmd/app
# Build migrate binary so container can apply migrations before start
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o /build/migrate ./cmd/migrate

FROM alpine:3.18
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /build/company-chat /app/company-chat
COPY --from=builder /build/migrate /app/migrate
RUN mkdir -p /app/uploads
EXPOSE 8080
# Run migrations first then start the server
ENTRYPOINT ["/bin/sh", "-c", "/app/migrate && /app/company-chat"]
