FROM golang:1.19-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY ./pkg ./pkg


EXPOSE 8080

RUN go build pkg/main.go

CMD [ "./main" ]
