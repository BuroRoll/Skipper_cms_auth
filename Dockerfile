FROM golang:1.18-alpine AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY ./pkg ./pkg


EXPOSE 8080

RUN go build pkg/main.go

FROM scratch
COPY --from=build /app/main /bin
ENTRYPOINT ["/bin/main"]
#CMD [ "./main" ]
