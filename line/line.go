package line

import (
	"encoding/json"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/line/line-bot-sdk-go/linebot"
)


type Line struct {
	ChannelSecret string
	ChannelToken  string
	Client        *linebot.Client
}

func SetUpLineClient() (*Line, error) {
	line := &Line{
		ChannelSecret: os.Getenv("LINE_BOT_CHANNEL_SECRET"),
		ChannelToken:  os.Getenv("LINE_BOT_CHANNEL_TOKEN"),
	}

	bot, err := linebot.New(
		line.ChannelSecret,
		line.ChannelToken,
	)
	if err != nil {
		return nil, err
	}

	line.Client = bot

	return line, nil
}

func (l *Line) ParseRequest(r events.APIGatewayProxyRequest) ([]*linebot.Event, error) {
	req := &struct {
		Events []*linebot.Event `json:"events"`
	}{}
	if err := json.Unmarshal([]byte(r.Body), req); err != nil {
		return nil, err
	}

	return req.Events, nil
}

// NewSelectMessage select button message
func (l *Line) NewSelectMessage(displayMsg string, selectActions ...linebot.TemplateAction) *linebot.TemplateMessage {
	return linebot.NewTemplateMessage(
		"エアコンの電源を入れますか？",
		linebot.NewButtonsTemplate("", "エアコンの電源を入れますか？", "please select", selectActions...),
	)
}
