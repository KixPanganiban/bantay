package lib

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/KixPanganiban/bantay/log"
	"github.com/hako/durafmt"
	"github.com/mailgun/mailgun-go"
	"github.com/nlopes/slack"
)

// Reporter consumes a CheckResult to flush into some predefined sink
type Reporter interface {
	Report(CheckResult, *map[string]int) error
}

// LogReporter implements Reporter by writing to log
type LogReporter struct {
	ServerConfig ParsedServer
}

// Report writes to log
func (lr LogReporter) Report(c CheckResult, dc *map[string]int) error {
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
	ServerConfig ParsedServer
	SlackToken   string
	SlackChannel string
	FailedOnly   bool
}

// Report sends an update to Slack
func (sr SlackReporter) Report(c CheckResult, dc *map[string]int) error {
	client := slack.New(sr.SlackToken)
	switch c.Success {
	case true:
		{
			if sr.FailedOnly == false && (*dc)[c.Name] == 0 {
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
			} else if (*dc)[c.Name] != 0 {
				attachment := slack.Attachment{
					Color:  "#36a64f",
					Footer: "bantay uptime check",
					Text:   fmt.Sprintf("%s is back up.", c.Name),
					Fields: []slack.AttachmentField{
						slack.AttachmentField{
							Title: "Failed Check Count",
							Value: strconv.Itoa((*dc)[c.Name]),
						},
						slack.AttachmentField{
							Title: "Total Downtime",
							Value: durafmt.Parse(time.Duration(math.Ceil((float64((*dc)[c.Name]) * float64(sr.ServerConfig.PollInterval)))) * time.Second).String(),
						},
					},
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
			var attachment slack.Attachment
			if (*dc)[c.Name] == 0 {
				attachment = slack.Attachment{
					Color: "#bd2f2f",
					Fields: []slack.AttachmentField{
						slack.AttachmentField{
							Title: "Reason",
							Value: c.Message,
						},
					},
					Footer: "bantay uptime check",
					Text:   fmt.Sprintf("%s went down.", c.Name),
				}
			} else {
				attachment = slack.Attachment{
					Color: "#bd2f2f",
					Fields: []slack.AttachmentField{
						slack.AttachmentField{
							Title: "Reason",
							Value: c.Message,
						},
						slack.AttachmentField{
							Title: "Failed Check Count",
							Value: strconv.Itoa((*dc)[c.Name] + 1),
						},
					},
					Footer: "bantay uptime check",
					Text:   fmt.Sprintf("%s is still down.", c.Name),
				}
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

// MailgunReporter reports up and down events via email using Mailgun
type MailgunReporter struct {
	ServerConfig      ParsedServer
	MailgunDomain     string
	MailgunPrivateKey string
	MailgunRecipients []string
	MailgunSender     string
	MailgunExclude    []string
}

// Report sends an email via Mailgun
func (mr MailgunReporter) Report(c CheckResult, dc *map[string]int) error {
	mg := mailgun.NewMailgun(mr.MailgunDomain, mr.MailgunPrivateKey)
	var (
		body    string
		subject string
	)
	for _, e := range mr.MailgunExclude {
		if c.Name == e {
			return nil
		}
	}
	if c.Success == true && (*dc)[c.Name] != 0 {
		body = fmt.Sprintf(
			"%s is back up. Estimated total downtime: %s.",
			c.Name,
			durafmt.Parse(time.Duration(math.Ceil((float64((*dc)[c.Name])*float64(mr.ServerConfig.PollInterval))))*time.Second).String(),
		)
		subject = fmt.Sprintf(
			"[%s] %s is back up",
			time.Now().Format("01/02/06 15:04:05 MST"),
			c.Name,
		)
	} else if c.Success == false && (*dc)[c.Name] == 0 {
		body = fmt.Sprintf(
			"%s went down. Reason: %s",
			c.Name,
			c.Message,
		)
		subject = fmt.Sprintf(
			"[%s] %s went down",
			time.Now().Format("01/02/06 15:04:05 MST"),
			c.Name,
		)
	}
	message := mg.NewMessage(mr.MailgunSender, subject, body, mr.MailgunRecipients...)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, _, err := mg.Send(ctx, message)

	if err != nil {
		return err
	}

	return nil
}
