FROM golang:1.23 AS builder

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o main ./cmd/server/main.go

FROM gcr.io/distroless/static
WORKDIR /app
EXPOSE 8081
COPY --from=builder /app/main /app/main
ENTRYPOINT [ "./main" ]