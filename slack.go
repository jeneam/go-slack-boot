package main

import (
	"fmt"
	"github.com/nlopes/slack"
	"strings"
)

// SlackListener ...
type SlackListener struct {
	client    *slack.Client
	botID     string
	channelID string
}

const (
	// action is used for slack attament action.
	actionSelect = "select"
	actionStart  = "start"
	actionCancel = "cancel"
)

// handleMesageEvent handles message events.
func (s *SlackListener) handleMessageEvent(ev *slack.MessageEvent) error {

	// Only response in specific channel. Ignore else.
	if ev.Channel != s.channelID {
		Log(fmt.Sprintf("%s %s", ev.Channel, ev.Msg.Text))
		return nil
	}

	// Only response mention to bot. Ignore else.
	if !strings.HasPrefix(ev.Msg.Text, fmt.Sprintf("<@%s> ", s.botID)) {
		return nil
	}

	// Parse message
	m := strings.Split(strings.TrimSpace(ev.Msg.Text), " ")[1:]
	if len(m) == 0 || m[0] != "hey" {
		return fmt.Errorf("invalid message")
	}

	// value is passed to message handler when request is approved.
	attachment := slack.Attachment{
		Text:       "Which beer do you want? :beer:",
		Color:      "#f9a41b",
		CallbackID: "beer",
		Actions: []slack.AttachmentAction{
			{
				Name: actionSelect,
				Type: "select",
				Options: []slack.AttachmentActionOption{
					{
						Text:  "Asahi Super Dry",
						Value: "Asahi Super Dry",
					},
					{
						Text:  "Kirin Lager Beer",
						Value: "Kirin Lager Beer",
					},
					{
						Text:  "Sapporo Black Label",
						Value: "Sapporo Black Label",
					},
					{
						Text:  "Suntory Malts",
						Value: "Suntory Malts",
					},
					{
						Text:  "Yona Yona Ale",
						Value: "Yona Yona Ale",
					},
				},
			},

			{
				Name:  actionCancel,
				Text:  "Cancel",
				Type:  "button",
				Style: "danger",
			},
		},
	}

	params := slack.PostMessageParameters{
		Attachments: []slack.Attachment{
			attachment,
		},
	}

	if _, _, err := s.client.PostMessage(ev.Channel, "", params); err != nil {
		return fmt.Errorf("failed to post message: %s", err)
	}

	return nil
}
