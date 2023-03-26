package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/XiaoMengXinX/OneTextAPI-Go"
	"github.com/XiaoMengXinX/Telegram-OneText-bot/font"
	"github.com/XiaoMengXinX/Telegram-OneText-bot/utils"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kolesa-team/go-webp/decoder"
	"github.com/kolesa-team/go-webp/webp"
)

var onetextJSON []byte

type ShortURLResp struct {
	Token string `json:"token"`
	Error string `json:"error"`
}

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

	update, err := bot.HandleUpdate(r)
	if err != nil {
		return
	}

	if update.Message == nil {
		return
	}
	if update.Message.Command() == "onetext" {
		o := onetext.New()
		o.ReadBytes(onetextJSON)
		var sentence onetext.Sentence
		if update.Message.CommandArguments() != "" {
			index, _ := strconv.Atoi(update.Message.CommandArguments())
			sentence = o.Get(index)
		}
		if sentence.Text == "" {
			sentence = o.Random()
		}
		if err := sendOnetextImg(bot, utils.OnetextData{Sentence: sentence}, update.Message.Chat.ID, update.Message.MessageID); err != nil {
			log.Println(err)
			return
		}
	}
	if update.Message.Command() == "quote" {
		if update.Message.ReplyToMessage == nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Please reply to a message")
			msg.ParseMode = tgbotapi.ModeMarkdownV2
			msg.ReplyToMessageID = update.Message.MessageID
			if _, err := bot.Send(msg); err != nil {
				log.Println(err)
			}
			return
		}
		s := utils.OnetextData{
			Sentence: onetext.Sentence{
				Text: update.Message.ReplyToMessage.Text,
				By:   update.Message.ReplyToMessage.From.FirstName + " " + update.Message.ReplyToMessage.From.LastName,
			}}
		if update.Message.ReplyToMessage.ForwardSenderName != "" {
			s.By = update.Message.ReplyToMessage.ForwardSenderName
		}
		if len(update.Message.ReplyToMessage.Photo) != 0 {
			photoURL, _ := bot.GetFileDirectURL(update.Message.ReplyToMessage.Photo[len(update.Message.ReplyToMessage.Photo)-1].FileID)
			photo, err := getFile(photoURL)
			if err != nil {
				s.Text = fmt.Sprintf("[图片]\n%v", err)
			}
			s.Image, _, err = image.Decode(bytes.NewReader(photo))
			if err != nil {
				s.Text = fmt.Sprintf("[图片]\n%v", err)
			}
			if update.Message.ReplyToMessage.Caption != "" {
				s.Text += "\n" + update.Message.ReplyToMessage.Caption
			}
		}
		if update.Message.ReplyToMessage.Sticker != nil {
			stickerURL, _ := bot.GetFileDirectURL(update.Message.ReplyToMessage.Sticker.FileID)
			sticker, err := getFile(stickerURL)
			if err != nil {
				s.Text = fmt.Sprintf("[贴纸]\n%v", err)
			}
			s.Image, err = webp.Decode(bytes.NewReader(sticker), &decoder.Options{})
			if err != nil {
				s.Text = fmt.Sprintf("[贴纸]\n%v", err)
			}
		}
		if err := sendOnetextImg(bot, s, update.Message.Chat.ID, update.Message.MessageID); err != nil {
			log.Println(err)
			return
		}
	}
	if update.Message.Command() == "custom" {
		if update.Message.CommandArguments() == "" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Please input your custom text\\. Spilt arguments by newline\\.\nFor example:\n```\n/custom Some\\\\nrandom\\\\ntext\nAuthor\nSource\nhttps://example.com```")
			msg.ParseMode = tgbotapi.ModeMarkdownV2
			msg.ReplyToMessageID = update.Message.MessageID
			if _, err := bot.Send(msg); err != nil {
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
			if i == 3 {
				short, err := shortURL(arg)
				if err != nil {
					log.Println(err)
				} else {
					s.Uri = short
				}
			}
		}
		if err := sendOnetextImg(bot, utils.OnetextData{Sentence: s}, update.Message.Chat.ID, update.Message.MessageID); err != nil {
			log.Println(err)
			return
		}
	}
}

func sendOnetextImg(bot *tgbotapi.BotAPI, s utils.OnetextData, chatID int64, messageID int) (err error) {
	img, err := utils.CreateOnetextImage(s, font.BuiltinFont)
	if err != nil {
		return err
	}
	msg := tgbotapi.NewPhoto(chatID, tgbotapi.FileBytes{
		Name:  "onetext.png",
		Bytes: img,
	})
	msg.ReplyToMessageID = messageID
	_, err = bot.Send(msg)
	return err
}

func shortURL(u string) (short string, err error) {
	resp, _ := http.PostForm("https://xve.me/api/url", url.Values{"url": {u}})
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var respData ShortURLResp
	if err = json.Unmarshal(body, &respData); err != nil {
		return
	}
	if respData.Error != "" {
		return "", errors.New(respData.Error)
	}
	short = fmt.Sprintf("https://xve.me/%s", respData.Token)
	return
}

func getFile(url string) (file []byte, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	file, err = io.ReadAll(resp.Body)
	return
}
