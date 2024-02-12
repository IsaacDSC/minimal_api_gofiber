FROM golang:1.21 as builder

WORKDIR /usr/app

COPY . .

RUN GOOS=linux GOARCH=amd64 go build -o server

# FROM scratch
FROM alpine

COPY --from=builder /usr/app /app

WORKDIR /app

EXPOSE 3333

CMD ["/app/server"]