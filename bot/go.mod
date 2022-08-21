module bot

go 1.18

require (
	github.com/XiaoMengXinX/OneTextAPI-Go v0.0.0-20220121152125-864a1aabac29
	github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.5.1
)

require utils v0.0.0

require (
	github.com/fogleman/gg v1.3.0 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/skip2/go-qrcode v0.0.0-20200617195104-da1b6568686e // indirect
	golang.org/x/image v0.0.0-20220722155232-062f8c9fd539 // indirect
	golang.org/x/text v0.3.7 // indirect
)

replace utils => ../
