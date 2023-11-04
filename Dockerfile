FROM golang:1.21.2-alpine3.17 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN GOOS=linux CGO_ENABLED=0 go build -o /app/stock ./cmd

FROM alpine:latest as final
WORKDIR /app
COPY --from=builder /app/stock .
COPY --from=builder /app/static /static
COPY --from=builder /app/seed.sql /
EXPOSE 80
CMD ["./stock"]