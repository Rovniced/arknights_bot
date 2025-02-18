package player

import (
	bot "arknights_bot/config"
	"arknights_bot/plugins/account"
	"arknights_bot/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"strings"
)

const (
	OP_STATE = "state" // 实时数据
	OP_BOX   = "box"   // 我的干员
	OP_GACHA = "gacha" // 抽卡记录
)

// PlayerData 角色数据
func PlayerData(callBack tgbotapi.Update) (bool, error) {
	callbackQuery := callBack.CallbackQuery
	data := callBack.CallbackData()
	d := strings.Split(data, ",")

	if len(d) < 6 {
		return true, nil
	}

	userId := callbackQuery.From.ID
	chatId := callbackQuery.Message.Chat.ID
	operate := d[1]
	clickUserId, _ := strconv.ParseInt(d[2], 10, 64)
	uid := d[3]
	messageId, _ := strconv.Atoi(d[4])
	param := d[5]

	if userId != clickUserId {
		answer := tgbotapi.NewCallbackWithAlert(callbackQuery.ID, "这不是你的角色！")
		bot.Arknights.Send(answer)
		return true, nil
	}

	// 获取账号信息
	var userAccount account.UserAccount
	utils.GetAccountByUserId(userId).Scan(&userAccount)

	// 判断操作类型
	switch operate {
	case OP_STATE:
		return State(uid, userAccount, chatId, messageId)
	case OP_BOX:
		return Box(uid, userAccount, chatId, messageId, param)
	case OP_GACHA:
		return Gacha(uid, userAccount, chatId, messageId)
	}

	return true, nil
}
