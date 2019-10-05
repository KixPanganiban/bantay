package lib

import (
	"fmt"

	"github.com/KixPanganiban/bantay/log"
	"github.com/nlopes/slack"
)

// Reporter consumes a CheckResult to flush into some predefined sink
type Reporter interface {
	Report(CheckResult) error
}

// LogReporter implements Reporter by writing to log
type LogReporter struct{}

// Report writes to log
func (lr LogReporter) Report(c CheckResult) error {
	switch c.Success {
	case true:
		{
			log.Infof("[%s] Check successful.", c.Name)
		}
	case false:
		{
			log.Debugf("[%s] Check failed. Reason: %s", c.Name, c.Message)
		}
	}
	return nil
}

// SlackReporter reports check outputs to Slack
type SlackReporter struct {
	SlackToken   string `yaml:"slack_token"`
	SlackChannel string `yaml:"slack_channel"`
	FailedOnly   bool   `yaml:"failed_only"`
}

// Report sends an update to Slack
func (sr SlackReporter) Report(c CheckResult) error {
	client := slack.New(sr.SlackToken)
	switch c.Success {
	case true:
		{
			if sr.FailedOnly == false {
				attachment := slack.Attachment{
					Color:  "#36a64f",
					Footer: "bantay uptime check",
					Text:   fmt.Sprintf("%s check succeeded.", c.Name),
				}
				_, _, err := client.PostMessage(
					sr.SlackChannel,
					slack.MsgOptionAsUser(false),
					slack.MsgOptionUsername("bantay"),
					slack.MsgOptionAttachments(attachment),
				)
				if err != nil {
					return err
				}
			}
		}
	case false:
		{
			attachment := slack.Attachment{
				Color: "#bd2f2f",
				Fields: []slack.AttachmentField{
					slack.AttachmentField{
						Title: "Reason",
						Value: c.Message,
					},
				},
				Footer: "bantay uptime check",
				Text:   fmt.Sprintf("%s check failed.", c.Name),
			}
			_, _, err := client.PostMessage(
				sr.SlackChannel,
				slack.MsgOptionAsUser(false),
				slack.MsgOptionUsername("bantay"),
				slack.MsgOptionAttachments(attachment),
			)
			if err != nil {
				return err
			}
		}

	}
	return nil
}
