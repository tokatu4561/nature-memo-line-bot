
# nature-remo-line-bot
nature-remo を lineから操作する

## nature remo
https://home.nature.global/
api
https://swagger.nature.global/

## LINE
linebot
https://developers.line.biz/ja/docs/messaging-api/building-bot/

## start
```
env.jsonに
API_TOKEN (nature_remoのトークン)
LINE_BOT_CHANNEL_SECRET
LINE_BOT_CHANNEL_TOKEN
をセット
```

```bash
make build
make start api
ngrok 3000
```
