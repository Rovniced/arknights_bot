package scheduled

import (
	"arknights_bot/bot/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
)

func BilibiliNews() func() {
	bilibiliNews := func() {
		text, pics := utils.ParseBilibiliDynamic()
		if text != "" {
			groups := utils.GetJoinedGroups()
			if pics == nil {
				for _, group := range groups {
					groupNumber, _ := strconv.ParseInt(group, 10, 64)
					sendMessage := tgbotapi.NewMessage(groupNumber, text)
					utils.SendMessage(sendMessage)
				}
			} else if len(pics) == 1 {
				for _, group := range groups {
					groupNumber, _ := strconv.ParseInt(group, 10, 64)
					sendPhoto := tgbotapi.NewPhoto(groupNumber, tgbotapi.FileURL(pics[0]))
					sendPhoto.Caption = text
					utils.SendPhoto(sendPhoto)
				}
			} else {
				for _, group := range groups {
					groupNumber, _ := strconv.ParseInt(group, 10, 64)
					var mediaGroup tgbotapi.MediaGroupConfig
					var media []interface{}
					mediaGroup.ChatID = groupNumber
					for i, pic := range pics {
						var inputPhoto tgbotapi.InputMediaPhoto
						inputPhoto.Media = tgbotapi.FileURL(pic)
						inputPhoto.Type = "photo"
						if i == 0 {
							inputPhoto.Caption = text
						}
						media = append(media, inputPhoto)
					}
					mediaGroup.Media = media
					utils.SendMediaGroup(mediaGroup)
				}
			}
		}
	}
	return bilibiliNews
}