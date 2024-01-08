package main

// import (
// 	"log"
// 	"os"

// 	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
// 	"gopkg.in/yaml.v2"
// )

// var bot *tgbotapi.BotAPI
// var QuizData []Quiz

// //go:embed quizdata.yaml
// var yamlFile []byte

// func init() {
// 	var err error
// 	bot, err = tgbotapi.NewBotAPI(os.Getenv("TOKEN"))
// 	if err != nil {
// 		log.Panic(err)
// 	}

// 	if err := yaml.Unmarshal(yamlFile, &QuizData); err != nil {
// 		log.Fatalf("Unmarshal error: %v", err)
// 	}

// 	log.Printf("Authorized on account %s", bot.Self.UserName)
// }

// func main() {
// 	u := tgbotapi.NewUpdate(0)
// 	u.Timeout = 60

// 	updates := bot.GetUpdatesChan(u)

// 	for update := range updates {
// 		processUpdate(update)
// 	}
// }
