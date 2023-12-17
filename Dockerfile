FROM golang:1.21.3-alpine3.18 as build
WORKDIR /mtg-price-scrapper
# Copy dependencies list
COPY go.mod go.sum ./
COPY scrapper ./scrapper
COPY main.go .
RUN go mod download
# Build
RUN env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o main main.go
# Copy artifacts to a clean image
FROM alpine:3.18
COPY --from=build /mtg-price-scrapper/main /main
ENTRYPOINT [ "/main" ]
