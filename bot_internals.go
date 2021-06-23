package main

import (
	tb "gopkg.in/tucnak/telebot.v2"
	"time"
)

func createBot() *tb.Bot {
	// Creates a bot using the default URL setting (https://api.telegram.org), this can be overwritten if desired.
	bot, err := tb.NewBot(tb.Settings{
		Token: getToken(),
		// The 10 second timeout it's the default and recommended setting, a lower number seems to affect performance.
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	handleError(err, true)

	// Logging
	log.Infoln("Bot connected to Telegram Servers")

	// Change the default URL if you override the URL on tb.Settings.
	sendMessageToAdmin(bot, "Connection successful to "+tb.DefaultApiURL)

	return bot
}

func handleEndpoint(bot *tb.Bot, route string, message string, privateMsg bool) {
	bot.Handle(route, func(src *tb.Message) {
		chatID := tb.ChatID(src.Chat.ID)
		if privateMsg {
			_, errSend := bot.Send(src.Sender, message, "html")
			if !handleError403(bot, chatID, errSend) {
				handleError(errSend, false)
				sendMessage(bot, chatID, "Te lo he enviado por privado shur")
			}
			logEndpointUsage(src, route)
		} else {
			// This is the default way on handling endpoints.
			_, errSend := bot.Send(chatID, message, "html")
			handleError(errSend, false)
			logEndpointUsage(src, route)
		}
	})
}

func sendMessage(bot *tb.Bot, chatID tb.ChatID, message string) {
	_, err := bot.Send(chatID, message, "html")
	handleError(err, false)
}

func sendMessageToAdmin(bot *tb.Bot, message string) {
	// This is hardcoded for now, since I'm the only admin
	chatID := tb.ChatID(1099020633)
	_, err := bot.Send(chatID, message)
	handleError(err, false)
	log.Warnln("Message: " + message + ". Sent to admins.")
}