FROM golang:alpine

RUN apk update && apk add --no-cache git && apk add --no-cach bash && apk add build-base

RUN mkdir /app
WORKDIR /app

COPY go.mod go.sum ./
COPY scrapper ./scrapper
COPY *.go ./

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /build

EXPOSE 8080

CMD ["/build"]