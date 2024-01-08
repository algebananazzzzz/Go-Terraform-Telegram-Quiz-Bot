package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var passphrase = os.Getenv("PASSPHRASE")

func defaultHandler(update tgbotapi.Update, user User) {
	switch update.Message.Command() {
	case "start":

		message := os.Getenv("START_MESSAGE")
		if passphrase == "" {
			replyKeyboard := tgbotapi.NewReplyKeyboard(
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton("Start Quiz"),
				),
			)
			message += " Start quiz?"
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, message)
			msg.ReplyMarkup = replyKeyboard
			bot.Send(msg)
		} else {
			message += " Provide passphrase to continue."
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, message)
			bot.Send(msg)
		}
	case "":
		var msg tgbotapi.MessageConfig

		if update.Message.Text == passphrase || passphrase == "" {
			message := "Please choose a Quiz."
			buttons := make([][]tgbotapi.InlineKeyboardButton, len(QuizData))

			for index, quiz := range QuizData {
				buttons[index] = tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(quiz.Name, fmt.Sprint(quiz.Id)))
			}

			inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(buttons...)
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, message)
			msg.ReplyMarkup = inlineKeyboard

			user.UserConvState = QuizState
			if err := dumpUserData(user); err != nil {
				log.Fatal(err)
			}
		} else {
			message := "Incorrect Passphrase. Try again."
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, message)
		}
		bot.Send(msg)
	default:
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, unknownCommandMessage)
		bot.Send(msg)
	}
}

func genPoll(chatId int64, question Question) tgbotapi.SendPollConfig {
	poll := tgbotapi.NewPoll(chatId, question.Title, question.Options...)
	poll.Type = "quiz"
	poll.IsAnonymous = false
	poll.CorrectOptionID = int64(question.Answer)
	return poll
}

func quizHandler(update tgbotapi.Update, user User) {
	if update.CallbackQuery != nil {
		quizId, err := strconv.Atoi(update.CallbackQuery.Data)
		if err != nil {
			log.Fatal(err)
		}
		quiz := QuizData[quizId]
		editMsg := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, fmt.Sprintf("Starting quiz: %s. There are %d questions.", quiz.Name, quiz.Length))
		bot.Send(editMsg)

		poll := genPoll(update.CallbackQuery.Message.Chat.ID, quiz.Questions[0])
		bot.Send(poll)

		userData := UserDataStruct{Score: 0, QuizId: quiz.Id, NextQnId: 1, PrevQnAnswer: quiz.Questions[0].Answer}
		user.UserData = userData

		dumpUserData(user)
	} else if update.PollAnswer != nil {
		quiz := QuizData[user.UserData.QuizId]
		if update.PollAnswer.OptionIDs[0] == int(user.UserData.PrevQnAnswer) {
			user.UserData.Score += 1
		} else {
			msg := tgbotapi.NewMessage(update.PollAnswer.User.ID, fmt.Sprintf("Wrong answer. Your current score is %d out of %d.", user.UserData.Score, quiz.Length))
			bot.Send(msg)
		}

		if user.UserData.NextQnId == quiz.Length {
			msg := tgbotapi.NewMessage(update.PollAnswer.User.ID, fmt.Sprintf("You have come to the end of quiz: %s. Your final score is %d out of %d.", quiz.Name, user.UserData.Score, quiz.Length))
			bot.Send(msg)
			user.UserData = UserDataStruct{}
			user.UserConvState = DefaultState
			dumpUserData(user)
			return
		}

		nextQuestion := quiz.Questions[user.UserData.NextQnId]
		poll := genPoll(update.PollAnswer.User.ID, nextQuestion)
		bot.Send(poll)

		user.UserData.PrevQnAnswer = nextQuestion.Answer
		user.UserData.NextQnId += 1

		dumpUserData(user)
	} else if update.Message != nil {
		switch update.Message.Command() {
		case "exit":
			quiz := QuizData[user.UserData.QuizId]
			msg := tgbotapi.NewMessage(update.Message.From.ID, fmt.Sprintf("Exiting quiz: %s. Your final score is %d out of %d.", quiz.Name, user.UserData.Score, quiz.Length))
			bot.Send(msg)
			user.UserData = UserDataStruct{}
			user.UserConvState = DefaultState
			dumpUserData(user)
		case "":
			msg := tgbotapi.NewMessage(update.Message.From.ID, "/exit to exit quiz.")
			bot.Send(msg)
		default:
			msg := tgbotapi.NewMessage(update.Message.From.ID, unknownCommandMessage)
			bot.Send(msg)
		}
	} else {
		user.UserConvState = DefaultState
		dumpUserData(user)
		log.Fatal("Not valid stage. Exiting Quiz state.")
	}

}

func processUpdate(update tgbotapi.Update) error {
	updateUser := update.SentFrom()

	var userId int64
	if updateUser != nil {
		userId = updateUser.ID
	} else if update.PollAnswer != nil {
		userId = update.PollAnswer.User.ID
	} else {
		log.Fatal("No valid userid")
	}

	user, err := getUserData(userId)
	ConvState := user.UserConvState

	if err != nil {
		log.Fatal(err)
	}

	switch ConvState {
	case DefaultState:
		defaultHandler(update, user)
	case QuizState:
		quizHandler(update, user)
	}
	return nil
}
