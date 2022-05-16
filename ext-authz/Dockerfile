FROM --platform=linux/amd64 golang:1.17

WORKDIR /app

COPY go.mod ./

COPY go.sum ./

RUN go mod download

COPY *.go ./

RUN go build -o ./ext-authz-server

EXPOSE 8080
EXPOSE 50051

CMD [ "./ext-authz-server" ]