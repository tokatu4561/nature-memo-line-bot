package line

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/line/line-bot-sdk-go/linebot"
)


type Line struct {
	ChannelSecret string
	ChannelToken  string
	Client        *linebot.Client
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