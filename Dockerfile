
FROM golang:latest
MAINTAINER Daniel Rocha
WORKDIR /var/www
COPY . /var/www
RUN go get -d -v ./...
RUN go install -v ./...
EXPOSE 1337