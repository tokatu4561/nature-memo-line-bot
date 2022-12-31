package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/tokatu4561/nature-memo-line-bot/line"
)

type Appliance struct {
	Id string
	Type string
	Nickname string
}

const AWS_REGION = "ap-northeast-1"

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	line, err := line.SetUpLineClient()
	if err != nil {
		return events.APIGatewayProxyResponse{Body: "LINE接続エラー", StatusCode: 500}, err
	}

	lineEvents, err := line.ParseRequest(request)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: "LINE接続エラー", StatusCode: 500}, err
	}

	for _, event := range lineEvents {
		// イベントがメッセージの受信だった場合
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {

			case *linebot.TextMessage:
				replyMessage := message.Text

				switch replyMessage {
				case "エアコン":
					actions := []linebot.TemplateAction {
						// linebot.NewPostbackAction("On", "appliances={air}&status=on", "on", "on"),
						linebot.NewPostbackAction("On", "on", "on", ""),
						linebot.NewPostbackAction("Off","off", "off", ""),
					}
					res := linebot.NewTemplateMessage(
						"エアコンの電源を入れますか？",
						linebot.NewButtonsTemplate("", "エアコンの電源を入れますか？", "please select", actions...),
					)

					_, err = line.Client.ReplyMessage(event.ReplyToken, res).Do()

					log.Println(err)

					if err != nil {
						return events.APIGatewayProxyResponse{
							Body:       err.Error(),
							StatusCode: 500,
						}, nil
					}
				case "照明":
					// appliances, err := fetchAppliances()
					// if err != nil {
					// 	log.Println(err)
					// 	return events.APIGatewayProxyResponse{
					// 		Body:       err.Error(),
					// 		StatusCode: 500,
					// 	}, nil
					// }

					// var lightApp *Appliance
					// for _, appliance := range appliances {
					// 	if (appliance.Type == "light") {
					// 		lightApp = appliance
					// 	}
					// }
				case "テレビ":
				}
				break
			default:
			}
		} else if event.Type == linebot.EventTypePostback {
			postBackData := event.Postback.Data

			appliances, err := fetchAppliances()

			var lightApp *Appliance
			for _, appliance := range appliances {
				if (appliance.Type == "LIGHT") {
					lightApp = appliance
				}
			}

			switch postBackData {
			case "on":
				err = switchPowerAppliance(lightApp, true)
			case "off":
				err = switchPowerAppliance(lightApp, false)
			}
			if err != nil {
				return events.APIGatewayProxyResponse{
					Body:       err.Error(),
					StatusCode: 500,
				}, nil
			}
		} else {}
	}

	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("Hello, %v", string("hello")),
		StatusCode: 200,
	}, nil
}

func switchPowerAppliance(app *Appliance, on bool) error{
	var switchText string
	if (on) {
		switchText = "on"
	} else {
		switchText = "off"
	}
	
	values := url.Values{}
    values.Set("button", switchText)
	
	baseUrl := os.Getenv("API_URL")
	path := fmt.Sprintf("1/appliances/%s/%s", app.Id, strings.ToLower(app.Type))
    endpoint := fmt.Sprintf("%s/%s", baseUrl, path)
	
	log.Println(switchText)
    req, err := http.NewRequest("POST", endpoint, strings.NewReader(values.Encode()))
	if err != nil {
        return err
    }

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("API_TOKEN")))

	client := new(http.Client)
	res, err := client.Do(req)
	if err != nil {
        return err
    }
    defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
        log.Printf("http status code: %d", res.StatusCode)
		log.Println(res.Body)
    }

	return nil
}

func fetchAppliances() ([]*Appliance, error){
	endpoint := fmt.Sprintf("%s/%s", os.Getenv("API_URL"), "1/appliances")
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
        return nil, err
    }

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("API_TOKEN")))

	client := new(http.Client)
	res, _ := client.Do(req)
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
        return nil, errors.New(fmt.Sprintf("http status code %d", res.StatusCode))
    }

	var appliances []*Appliance
	decoder := json.NewDecoder(res.Body)

	err = decoder.Decode(&appliances)
    if err != nil {
        return nil, err
    }

	return appliances, nil
}

func main() {
	lambda.Start(handler)
}