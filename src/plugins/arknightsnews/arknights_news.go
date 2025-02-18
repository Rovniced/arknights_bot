package arknightsnews

import (
	"arknights_bot/config"
	"arknights_bot/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
	"strings"
)

func BilibiliNews() func() {
	return func() {
		text, pics := ParseBilibiliDynamic()
		if len(text) == 0 {
			return
		}
		groups := utils.GetJoinedGroups()
		if pics == nil {
			for _, group := range groups {
				sendMessage := tgbotapi.NewMessage(group, text)
				config.Arknights.Send(sendMessage)
			}
			return
		}

		if len(pics) == 1 {
			for _, group := range groups {
				sendPhoto := tgbotapi.NewPhoto(group, tgbotapi.FileURL(pics[0]))
				sendPhoto.Caption = text
				config.Arknights.Send(sendPhoto)
			}
			return
		}

		for _, group := range groups {
			var mediaGroup tgbotapi.MediaGroupConfig
			var media []interface{}
			mediaGroup.ChatID = group
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
			config.Arknights.SendMediaGroup(mediaGroup)
		}
	}
}

func ParseBilibiliDynamic() (string, []string) {
	var text string
	var pics []string
	url := viper.GetString("api.bilibili_dynamic")
	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Add("user-agent", "Mozilla/5.0")
	request.Header.Add("Cookie", "buvid3=4BCEC5AA-D998-B5AA-7869-5C38389B6F3338331infoc;buvid4=891E0F8F-21A6-2816-285B-DC05A009D61538331-024011112-ojR14xwcWnZIblFyUE6/cLQdGDJg2ZsOVAhfgVvhVv8vC3DWjzJwtgXMtg1ADDmT")
	request.Header.Add("Referer", "https://www.bilibili.com")
	resp, _ := http.DefaultClient.Do(request)
	readAll, _ := io.ReadAll(resp.Body)
	result := gjson.ParseBytes(readAll)
	items := result.Get("data.items").Array()
	for _, item := range items {
		top := item.Get("modules.module_tag.text").String()
		if top != "置顶" {
			dynamicType := item.Get("type").String()
			id := item.Get("id_str").String()
			link := "https://t.bilibili.com/" + id
			//publishTime := time.Unix(item.Get("modules.module_author.pub_ts").Int(), 0).Format("2006-01-02 15:04:05")
			if dynamicType == "DYNAMIC_TYPE_DRAW" {
				for _, pic := range item.Get("modules.module_dynamic.major.opus.pics").Array() {
					pics = append(pics, pic.Get("url").String())
				}
				text = item.Get("modules.module_dynamic.major.opus.summary.text").String()
			}
			if dynamicType == "DYNAMIC_TYPE_WORD" {
				text = item.Get("modules.module_dynamic.major.opus.summary.text").String()
			}
			if dynamicType == "DYNAMIC_TYPE_AV" {
				title := item.Get("modules.module_dynamic.major.archive.title").String() + "\n\n"
				desc := item.Get("modules.module_dynamic.major.archive.desc").String()
				cover := item.Get("modules.module_dynamic.major.archive.cover").String()
				vUrl := "https:" + item.Get("modules.module_dynamic.major.archive.jump_url").String()
				text = title + desc + "\n视频链接：" + vUrl
				pics = append(pics, cover)
			}
			if dynamicType == "DYNAMIC_TYPE_FORWARD" {
				desc := item.Get("modules.module_dynamic.desc.text").String()
				for _, pic := range item.Get("orig.modules.module_dynamic.major.opus.pics").Array() {
					pics = append(pics, pic.Get("url").String())
				}
				text = desc + "\n\n" + item.Get("orig.modules.module_dynamic.major.opus.summary.text").String()
			}
			if dynamicType == "DYNAMIC_TYPE_ARTICLE" {
				summary := item.Get("modules.module_dynamic.major.opus.summary.text").String()
				for _, pic := range item.Get("modules.module_dynamic.major.opus.pics").Array() {
					pics = append(pics, pic.Get("url").String())
				}
				text = strings.ReplaceAll(summary, "[图片]", "") + "\n\n专栏地址：https:" + item.Get("modules.module_dynamic.major.opus.jump_url").String()
			}
			if utils.RedisSetIsExists("tg_azurlane", link) {
				return "", nil
			}
			utils.RedisAddSet("tg_azurlane", link)
			break
		}
	}
	return text, pics
}