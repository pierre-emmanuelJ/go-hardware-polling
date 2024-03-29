FROM golang:1.9.2
MAINTAINER Pierre-Emmanuel Jacquier <pierre-emmanuel.jacquier@epitech.eu>

WORKDIR /go/src/github.com/pierre-emmanuelJ/go-exercises
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o exercises .

FROM alpine:3.6
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0  /go/src/github.com/pierre-emmanuelJ/go-exercises/exercises .
CMD ["./exercises", "-i", "1", "-p", "tmpfs", "-n", "eth0"]
