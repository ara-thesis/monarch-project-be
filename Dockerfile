FROM golang:1.18.1-alpine3.15

RUN mkdir /app

COPY . /app

WORKDIR /app

# server configuration
ENV PORT=8000

# postgres configuration
ENV PG_HOST=192.168.18.240
ENV PG_PORT=5432
ENV PG_USER=postgres
ENV PG_PASS=Raflis2001
ENV PG_DB=monarch-thesis

RUN go mod verify

RUN go mod tidy

RUN go build -o bin/server src/main.go

CMD [ "./bin/server" ]