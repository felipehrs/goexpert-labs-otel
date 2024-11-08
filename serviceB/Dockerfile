FROM golang:1.23 AS builder

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o cloudrun ./cmd/server/main.go

FROM gcr.io/distroless/static
WORKDIR /app
COPY --from=builder /app/cloudrun /app/cloudrun
ENTRYPOINT [ "./cloudrun" ]