FROM golang:1.24.0-alpine AS build
WORKDIR /mtg-price-scrapper
# Copy dependencies list
COPY api ./api
WORKDIR /mtg-price-scrapper/api
RUN go mod download
# Build
RUN env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o main cmd/main.go
# Copy artifacts to a clean image
FROM alpine
COPY --from=build /mtg-price-scrapper/api/main /main
ENTRYPOINT [ "/main" ]
