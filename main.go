package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/nlopes/slack"
	"log"
	"net/http"
	"os"
)

func main() {

	var token string
	var debug bool
	var botId string
	var channelId string

	flag.StringVar(&token, "token", "", "'--token' bot user token / env BOT_TOKEN.")
	flag.StringVar(&botId, "botId", "", "'--botId' bot id/ env BOT_ID. ")
	flag.StringVar(&channelId, "channelId", "", "'--channelId' slack channel /env BOT_CHANNEL.")
	flag.BoolVar(&debug, "debug", false, "'--debug' if true, debug enabled.")
	flag.Parse()

	token = ValidateParam(token, "BOT_TOKEN", "bot token expected. env BOT_TOKEN or arg --token")
	botId = ValidateParam(botId, "BOT_ID", "botId expected. env BOT_ID or arg --botId")
	channelId = ValidateParam(channelId, "BOT_CHANNEL", "bot channel expected. env BOT_CHANNEL or arg --channelId")

	client := slack.New(token)
	client.SetDebug(debug)

	slackListener := &SlackListener{
		client:    client,
		botID:     botId,
		channelID: channelId,
	}

	rtm := client.NewRTM()
	go rtm.ManageConnection()
	go HttpHandler()

	for {
		select {
		case msg := <-rtm.IncomingEvents:
			Log("Event Received: ")
			switch ev := msg.Data.(type) {
			case *slack.ConnectedEvent:
				Log(fmt.Sprintln("Connection counter:", ev.ConnectionCount))
			case *slack.MessageEvent:
				if err := slackListener.handleMessageEvent(ev); err != nil {
					Log(err.Error())
				}
			case *slack.RTMError:
				Log(ev.Error())

			case *slack.InvalidAuthEvent:
				LogError(errors.New("Invalid credentials"))

			default:
				//Take no action
			}
		}
	}
}
func HttpHandler() {
	http.Handle("/interaction", interactionHandler{
		verificationToken: "xoxb-336749255621-uyvA4SbsYrHmI5UZTz37dhEl",
	})

	Log(fmt.Sprintf("[INFO] Server listening on :%s", "8082"))
	if err := http.ListenAndServe(":"+"8082", nil); err != nil {
		LogError(fmt.Errorf("[ERROR] %s", err))
	}
}

func ValidateParam(flag string, env string, err string) string {
	data := flag
	if data == "" {
		data = os.Getenv(env)
	}

	if data == "" {
		LogError(errors.New(err))
	}

	return data
}

// Log logs message to console.
func Log(message string) {
	fmt.Println(message)
}

// LogError logs an error and exit if it's not nil.
func LogError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}