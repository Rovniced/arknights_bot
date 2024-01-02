package utils

import (
	"arknights_bot/config"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
	"strings"
)

func ParseBilibiliDynamic() (string, []string) {
	var text string
	var pics []string
	url := config.GetString("api.bilibili_dynamic")
	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Add("user-agent", "Mozilla/5.0")
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
			if RedisSetIsExists("tg_azurlane", link) {
				return "", nil
			}
			RedisAddSet("tg_azurlane", link)
			break
		}
	}
	return text, pics
}
