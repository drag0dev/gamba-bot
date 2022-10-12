FROM golang:1.17-alpine

WORKDIR /bot-src

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY concurrency ./concurrency/
COPY scraping ./scraping/
COPY *.go ./

RUN go build -o /bot

EXPOSE 8080

CMD [ "/bot" ]
