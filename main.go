package main

import (
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()
	app.Use(cors.New())

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		fmt.Println("MISSING_TELEGRAM_BOT_TOKEN")
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	messages := make(chan string)
	username := make(chan string)

	go func() {
		for update := range updates {
			if update.Message != nil { // If we got a message
				log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "hallo juga")
				msg.ReplyToMessageID = update.Message.MessageID

				bot.Send(msg)

				messages <- update.Message.Text
				username <- update.Message.From.UserName
			}
		}
	}()

	app.Get("/", func(c *fiber.Ctx) error {
		select {
		case message := <-messages:
			c.Set("Content-Type", "application/json")
			return c.JSON(fiber.Map{
				"message":  message,
				"username": <-username,
			})
		default:
			return c.SendFile("index.html")
		}
	})
	app.Listen(os.Getenv("0.0.0.0:$PORT"))

}
