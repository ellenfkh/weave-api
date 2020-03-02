FROM golang:1.14

RUN mkdir /build 
ADD . /build/
WORKDIR /build 

RUN go build

ENTRYPOINT [ "./weave-api" ]