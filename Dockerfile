FROM golang:1.20 as builder

WORKDIR /app


COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

WORKDIR /app/cmd
RUN go build -o /app/main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/


COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]