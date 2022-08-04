# Telegram OneText Bot

A telegram bot to create onetext image

## Usage

On telegram group or private chat, send a special command to the bot:

`/onetext` - Get a onetext image

`/custom` - Create a onetext image with your own text

For example, send the following message to the bot:

```
/custom Some random text
Author
Source
```

Then you will get a onetext image like this:

<img src="https://user-images.githubusercontent.com/19994286/182917358-8c3c9efe-509a-4f47-b2ff-67f55afd432b.jpg" width="60%" height="60%">

## Deploy

[![Deploy with Vercel](https://vercel.com/button)](https://vercel.com/new/clone?repository-url=https%3A%2F%2Fgithub.com%2FXiaoMengXinX%2FTelegram-OneText-bot)

After your deployment on vercel, set the webhook url of the bot by requesting bot api like this:

```
https://api.telegram.org/bot<YOUR_BOT_TOKEN>/setWebhook?url=https://<YOUR_DEPLOYMENT_URL>/<YOUR_BOT_TOKEN>
```

Then you can use your bot on telegram.
