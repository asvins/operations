FROM golang:1.5.1

WORKDIR /go/src/app
COPY . /go/src/app

RUN go get -u .
CMD app

EXPOSE 8080
