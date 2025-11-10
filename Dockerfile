FROM golang:1.24-alpine AS builder

WORKDIR /src
ENV CGO_ENABLED=0
ENV GOOS=linux

RUN apk add --no-cache make git
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN make build

# runtime
FROM alpine:latest
RUN apk add --no-cache ca-certificates

RUN addgroup -S appgroup && adduser -S appuser -G appgroup
WORKDIR /home/appuser
COPY --from=builder /src/build/file-sentinel .

RUN chown -R appuser:appgroup .
USER appuser

ENTRYPOINT ["./file-sentinel"]

