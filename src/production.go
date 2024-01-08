package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"net/http"
	"os"

	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"gopkg.in/yaml.v2"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var bot *tgbotapi.BotAPI
var QuizData []Quiz

//go:embed quizdata.yaml
var yamlFile []byte

func init() {
	var err error
	bot, err = tgbotapi.NewBotAPI(os.Getenv("TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	if err := yaml.Unmarshal(yamlFile, &QuizData); err != nil {
		log.Fatalf("Unmarshal error: %v", err)
	}

	log.Printf("Initialized cold start of bot: %s", bot.Self.UserName)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var update tgbotapi.Update

	if err := json.Unmarshal([]byte(req.Body), &update); err != nil {
		log.Panic(err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest}, nil
	}

	if err := processUpdate(update); err != nil {
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, os.Getenv("ERROR_MESSAGE")))
		log.Panic(err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest}, nil
	}

	return events.APIGatewayProxyResponse{StatusCode: http.StatusOK}, nil
}

func main() {
	lambda.Start(Handler)
}
