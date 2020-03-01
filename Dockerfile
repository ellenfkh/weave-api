FROM golang:1.14

#WORKDIR /Users/ellen/Desktop/scratch/weave-api
#COPY . .
RUN mkdir /build 
ADD . /build/
WORKDIR /build 

RUN go build

ENTRYPOINT [ "./weave-api" ]