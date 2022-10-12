FROM golang:1.17-alpine

WORKDIR /server-src

COPY go.mod ./
COPY go.sum ./
# this has to go
COPY .env ./
RUN go mod download

COPY *.go ./

RUN go build -o /gamba-api

EXPOSE 8080

CMD [ "/gamba-api" ]
