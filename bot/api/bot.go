package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
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
}

func BotHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	bot := &tgbotapi.BotAPI{
		Token:  strings.ReplaceAll(r.URL.Path, "/", ""),
		Client: &http.Client{},
		Buffer: 100,
	}
	bot.SetAPIEndpoint(tgbotapi.APIEndpoint)

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
		img, err := utils.CreateOnetextImage(s, utils.FontFile, 0.9)
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
	if update.Message.Command() == "quote" {
		if update.Message.ReplyToMessage == nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Please reply to a message")
			msg.ParseMode = tgbotapi.ModeMarkdownV2
			msg.ReplyToMessageID = update.Message.MessageID
			if _, err = bot.Send(msg); err != nil {
				log.Println(err)
			}
			return
		}
		s := onetext.Sentence{
			Text: update.Message.ReplyToMessage.Text,
			By:   update.Message.ReplyToMessage.From.FirstName + " " + update.Message.ReplyToMessage.From.LastName,
		}
		if len(update.Message.ReplyToMessage.Photo) != 0 {
			s.Text = "[图片]"
			if update.Message.Caption != "" {
				s.Text += "\n" + update.Message.Caption
			}
		}
		if update.Message.ReplyToMessage.Sticker != nil {
			s.Text = "[贴纸]"
		}
		img, err := utils.CreateOnetextImage(s, utils.FontFile, 0.9)
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
		if update.Message.CommandArguments() == "" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Please input your custom text\\. Spilt arguments by newline\\.\nFor example:\n```\n/custom Some random text\nAuthor\nSource\n```")
			msg.ParseMode = tgbotapi.ModeMarkdownV2
			msg.ReplyToMessageID = update.Message.MessageID
			if _, err = bot.Send(msg); err != nil {
				log.Println(err)
			}
			return
		}
		var s onetext.Sentence
		msgText := update.Message.CommandArguments()
		args := strings.Split(msgText, "\n")
		for i, arg := range args {
			if i == 0 {
				s.Text = strings.ReplaceAll(arg, "\\n", "\n")
			}
			if i == 1 {
				s.By = strings.ReplaceAll(arg, "\\n", "\n")
			}
			if i == 2 {
				s.From = strings.ReplaceAll(arg, "\\n", "\n")
			}
		}
		img, err := utils.CreateOnetextImage(s, utils.FontFile, 0.9)
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
