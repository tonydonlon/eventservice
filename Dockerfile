# syntax=docker/dockerfile:1
FROM golang:1.16 AS builder
WORKDIR /go/src/github.com/tonydonlon/eventservice
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -o app .

FROM alpine:latest  
WORKDIR /root/
RUN apk --no-cache add curl
EXPOSE 8080
COPY --from=builder /go/src/github.com/tonydonlon/eventservice/app .
CMD ["./app"]
