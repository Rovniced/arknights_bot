package init

import (
	"arknights_bot/bot/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

var Kawakaze *tgbotapi.BotAPI

func Bot() {
	token := config.GetString("bot.token")
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Println(err)
		return
	}
	Kawakaze = bot
	log.Println("机器人初始化完成")
}