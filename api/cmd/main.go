package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"mtg-price-scrapper-sg/handler"
	"mtg-price-scrapper-sg/pkg/config"
)

func main() {
	if config.IsTestEnv {
		start := time.Now()
		log.Println(handler.Search(context.Background(), events.APIGatewayProxyRequest{}))
		log.Println(fmt.Sprintf("Took: %s", time.Since(start)))
	} else {
		lambda.Start(handler.Search)
	}
}
