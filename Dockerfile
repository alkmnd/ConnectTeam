FROM golang:1.22

RUN go version
ENV GOPATH=/

COPY ./ ./

# install psql
RUN apt-get update
RUN apt-get -y install postgresql-client

# make wait-for-postgres.sh executable
#RUN chmod +x wait-for-postgres.sh

# build go app
RUN go mod download
RUN go build -o connectteam ./cmd/main.go

CMD ["./connectteam"]

#FROM alpine:latest
#
#RUN apk --no-cache add ca-certificates
#WORKDIR /root/
#
#COPY ./.bin/connectteam .
#COPY configs ./config/