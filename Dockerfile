FROM golang:1.18-alpine AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY ./pkg ./pkg
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app pkg/main.go

FROM scratch

WORKDIR /app

COPY --from=build /app/main ./main
EXPOSE 8080
ENTRYPOINT ["./main"]