package telegram

import "github.com/go-telegram/bot"

type BotMiddleware = func(bot.HandlerFunc) bot.HandlerFunc
