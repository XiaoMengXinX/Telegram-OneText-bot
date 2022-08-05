# Telegram OneText Bot

A telegram bot to create onetext image

Embed font: [LxgwWenKai](https://github.com/lxgw/LxgwWenKai)

## Usage

On telegram group or private chat, send a special command to the bot:

`/onetext` - Get an onetext image

`/custom` - Create an onetext image with your own text

For example, send the following message to the bot:

```
/custom Some random text
Author
Source
```

Then you will get an onetext image like this:

<img src="https://user-images.githubusercontent.com/19994286/183112581-f8e46e4e-53ad-4765-bdd7-d843a22af9fb.png" width="80%" height="80%">

## Deploy

[![Deploy with Vercel](https://vercel.com/button)](https://vercel.com/new/clone?repository-url=https%3A%2F%2Fgithub.com%2FXiaoMengXinX%2FTelegram-OneText-bot)

After your deployment on vercel, set the webhook url of the bot by requesting bot api like this:

```
https://api.telegram.org/bot<YOUR_BOT_TOKEN>/setWebhook?url=https://<YOUR_DEPLOYMENT_URL>/<YOUR_BOT_TOKEN>
```

Then you can use your bot on telegram.
