FROM golang:1.18.1-alpine3.15

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN go mod verify

RUN go mod tidy

RUN go build -o bin/server src/main.go

CMD [ "./bin/server" ]