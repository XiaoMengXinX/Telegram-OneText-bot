package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

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

	fmt.Println(update.Message.Text)

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
		msg := tgbotapi.NewChatPhoto(update.Message.Chat.ID, tgbotapi.FileBytes{
			Name:  "onetext.png",
			Bytes: img,
		})
		_, err = bot.Send(msg)
		if err != nil {
			log.Println(err)
			return
		}
	}
}
