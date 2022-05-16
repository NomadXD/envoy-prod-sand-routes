FROM --platform=linux/amd64 golang:1.17

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY *.go ./

RUN go build -o ./backend-server

EXPOSE 8001
EXPOSE 8002

CMD [ "./backend-server" ]