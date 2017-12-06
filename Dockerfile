FROM golang:1.9.2

RUN mkdir -p /go/src/github.com/realbot/faas-dcos/

WORKDIR /go/src/github.com/realbot/faas-dcos

COPY vendor     vendor
COPY handlers	handlers
#COPY types      types
COPY read_config.go  .
COPY server.go  .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o faas-dcos .

RUN /bin/bash -c 'source $HOME/.bashrc; echo $HOME'

FROM alpine:3.6
RUN apk --no-cache add ca-certificates
WORKDIR /root/

EXPOSE 8080
ENV http_proxy      ""
ENV https_proxy     ""

COPY --from=0 /go/src/github.com/realbot/faas-dcos/faas-dcos  .

CMD ["./faas-dcos"]