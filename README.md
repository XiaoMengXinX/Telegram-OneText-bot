# Telegram OneText Bot

A telegram bot to create onetext image.

Demo: [OneText Bot](https://t.me/one_text_bot)

Built-in font: [LxgwWenKai](https://github.com/lxgw/LxgwWenKai)

## Usage

On telegram group or private chat, send the listed command to the bot:

`/onetext` - Get an onetext image

`/quote` - Reply to a message and quote it

`/custom` - Create an onetext image with your own text

For example, send the following message to the bot:

```
/custom Some random text
Author
Source
```

Then you will get an onetext image like this:

<img src="https://user-images.githubusercontent.com/19994286/183112855-35acb999-6388-4736-94d5-2f672579c7bc.jpg" width="80%" height="80%">

## Deploy

**Notice: The Root Directory need to be set to "bot".**

[![Deploy with Vercel](https://vercel.com/button)](https://vercel.com/new/clone?repository-url=https%3A%2F%2Fgithub.com%2FXiaoMengXinX%2FTelegram-OneText-bot)

After your deployment on vercel, set the webhook url of the bot by requesting bot api like this:

```
https://api.telegram.org/bot<YOUR_BOT_TOKEN>/setWebhook?url=https://<YOUR_DEPLOYMENT_URL>/<YOUR_BOT_TOKEN>
```

Then you can use your bot on telegram.
