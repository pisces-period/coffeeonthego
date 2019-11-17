# https://docs.docker.com/develop/develop-images/dockerfile_best-practices/
# multi-stage go app Dockerfile

# 1) GO BUILD
FROM golang:1.13.4-alpine3.10 AS go-build

LABEL author="Yuri Neves <pisces.period@gmail.com>" \
      app="coffee-on-the-go" \
      description="Golang HTCPCP Implementation" \
      version="v1.0"

WORKDIR /go/
COPY . .

# alphanumerical sorting of multi-line arguments
RUN apk --no-cache add \
    bzr \
    git \
    mercurial && \
    go get gopkg.in/mgo.v2 && \
    go get gopkg.in/mgo.v2/bson && \ 
    go build -o ./... goapp

# 2) GO RUN
FROM alpine:3.10.3 AS go-run
COPY --from=go-build /go/ ./
CMD ["./..."]
