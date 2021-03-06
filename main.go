package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/nlopes/slack"
	"log"
	"os"
	"strconv"
)

func main() {

	var token string
	var debug bool
	var botID string
	var channelID string
	var verificationToken string
	var port int
	flag.StringVar(&token, "token", "", "'--token' bot user token / env BOT_TOKEN.")
	flag.StringVar(&botID, "botID", "", "'--botID' bot id/ env BOT_ID. ")
	flag.StringVar(&channelID, "channelID", "", "'--channelID' slack channel /env BOT_CHANNEL.")
	flag.StringVar(&verificationToken, "checkToken", "", "'--checkToken' slack channel /env CHECK_TOKEN.")
	flag.IntVar(&port, "port", 3000, "'--port' listening ports default 3000.")
	flag.BoolVar(&debug, "debug", false, "'--debug' if true, debug enabled.")
	flag.Parse()

	token = validateParam(token, "BOT_TOKEN", "bot token expected. env BOT_TOKEN or arg --token")
	botID = validateParam(botID, "BOT_ID", "botId expected. env BOT_ID or arg --botId")
	channelID = validateParam(channelID, "BOT_CHANNEL", "bot channel expected. env BOT_CHANNEL or arg --channelId")
	verificationToken = validateParam(verificationToken, "CHECK_TOKEN", "missing checkToken")

	client := slack.New(token)
	client.SetDebug(debug)

	slackListener := &SlackListener{
		client:    client,
		botID:     botID,
		channelID: channelID,
	}

	rtm := client.NewRTM()
	go rtm.ManageConnection()
	go HTTPServer(strconv.Itoa(port), verificationToken)

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

func validateParam(flag string, env string, err string) string {
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
