package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	onetext "github.com/XiaoMengXinX/OneTextAPI-Go"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"utils"
)

var onetextJSON []byte
var bot *tgbotapi.BotAPI

func init() {
	resp, _ := http.Get("https://raw.githubusercontent.com/lz233/OneText-Library/master/OneText-Library.json")
	onetextJSON, _ = io.ReadAll(resp.Body)
	bot, _ = tgbotapi.NewBotAPIWithClient(os.Getenv("BOT_TOKEN"), tgbotapi.APIEndpoint, &http.Client{})
}

func BotHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	body, _ := io.ReadAll(r.Body)

	var update tgbotapi.Update

	err := json.Unmarshal(body, &update)
	if err != nil {
		log.Println(err)
		return
	}

	if update.Message == nil {
		return
	}
	if update.Message.Command() == "onetext" {
		o := onetext.New()
		o.ReadBytes(onetextJSON)
		s := o.Random()
		img, err := utils.CreateOnetextImage(s)
		if err != nil {
			log.Println(err)
			return
		}
		msg := tgbotapi.NewPhoto(update.Message.Chat.ID, tgbotapi.FileBytes{
			Name:  "onetext.png",
			Bytes: img,
		})
		msg.ReplyToMessageID = update.Message.MessageID
		_, err = bot.Send(msg)
		if err != nil {
			log.Println(err)
			return
		}
	}
	if update.Message.Command() == "custom" {
		/*
			text:xxxx
			from:xxxx
			by:xxxx
			record_time:xxxx
			create_time:xxxx
		*/
		var s onetext.Sentence
		msgText := update.Message.CommandArguments()
		args := strings.Split(msgText, "\n")
		for _, arg := range args {
			if strings.HasPrefix(arg, "text:") {
				s.Text = fmt.Sprintf("%s", strings.TrimPrefix(arg, "text:"))
			}
			if strings.HasPrefix(arg, "from:") {
				s.From = strings.TrimPrefix(arg, "from:")
			}
			if strings.HasPrefix(arg, "by:") {
				s.By = strings.TrimPrefix(arg, "by:")
			}
			if strings.HasPrefix(arg, "record_time:") {
				s.Time = []string{strings.TrimPrefix(arg, "record_time:")}
			}
			if strings.HasPrefix(arg, "create_time:") && len(s.Time) > 0 {
				s.Time = append(s.Time, strings.TrimPrefix(arg, "create_time:"))
			}
		}
		img, err := utils.CreateOnetextImage(s)
		if err != nil {
			log.Println(err)
			return
		}
		msg := tgbotapi.NewPhoto(update.Message.Chat.ID, tgbotapi.FileBytes{
			Name:  "onetext.png",
			Bytes: img,
		})
		msg.ReplyToMessageID = update.Message.MessageID
		_, err = bot.Send(msg)
		if err != nil {
			log.Println(err)
			return
		}
	}
}
