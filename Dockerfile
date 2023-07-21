FROM golang:alpine

RUN apk update && apk add --no-cache git

WORKDIR /app

COPY . .

RUN go mod tidy

EXPOSE 8080

RUN go build -o binary

ENTRYPOINT ["/app/binary"]

