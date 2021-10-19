FROM golang:1.14

WORKDIR /agripa-api
COPY go.mod . 
COPY go.sum .

RUN apt-get update
RUN apt-get install vim -y

RUN go mod download
RUN go get -u github.com/swaggo/swag/cmd/swag
RUN go get -u github.com/cosmtrek/air