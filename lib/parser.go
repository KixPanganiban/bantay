package lib

import (
	"errors"

	"gopkg.in/yaml.v2"
)

// ParsedReporter represents the unmarshalled reporters config
type ParsedReporter struct {
	Type    string                 `yaml:"type"`
	Options map[string]interface{} `yaml:"options"`
}

// ParsedServer represents server config
type ParsedServer struct {
	PollInterval uint32 `yaml:"poll_interval"`
}

// ParsedConfig represents the unmarshalled YAML file
type ParsedConfig struct {
	Server            ParsedServer     `yaml:"server"`
	Checks            *[]Check         `yaml:"checks"`
	Reporters         []ParsedReporter `yaml:"reporters"`
	ExportedReporters []Reporter       `yaml:"donotunmarshal"`
}

// ParseYAML parses the given YAML File and outputs a ParsedConfig struct
func ParseYAML(b []byte) (ParsedConfig, error) {
	config := ParsedConfig{}
	err := yaml.Unmarshal(b, &config)
	if err != nil {
		return ParsedConfig{}, err
	}
	config.ExportedReporters = []Reporter{}
	for _, rconfig := range config.Reporters {
		switch rconfig.Type {
		case "log":
			{
				config.ExportedReporters = append(config.ExportedReporters,
					LogReporter{
						ServerConfig: config.Server,
					})
			}
		case "slack":
			{
				slackChannel, ok := rconfig.Options["slack_channel"].(string)
				if !ok {
					return ParsedConfig{}, errors.New("can't parse required Slack config slack_channel")
				}
				slackToken, ok := rconfig.Options["slack_token"].(string)
				if !ok {
					return ParsedConfig{}, errors.New("can't parse required Slack config slack_token")
				}
				failedOnly, ok := rconfig.Options["failed_only"].(bool)
				if !ok {
					failedOnly = true
				}
				config.ExportedReporters = append(
					config.ExportedReporters,
					SlackReporter{
						ServerConfig: config.Server,
						SlackChannel: slackChannel,
						SlackToken:   slackToken,
						FailedOnly:   failedOnly,
					})
			}
		case "mailgun":
			{
				mailgunDomain, ok := rconfig.Options["mailgun_domain"].(string)
				if !ok {
					return ParsedConfig{}, errors.New("can't parse required Mailgun config mailgun_domain")
				}
				mailgunPrivateKey, ok := rconfig.Options["mailgun_private_key"].(string)
				if !ok {
					return ParsedConfig{}, errors.New("can't parse required Mailgun config mailgun_private_key")
				}
				mailgunSender, ok := rconfig.Options["mailgun_sender"].(string)
				if !ok {
					return ParsedConfig{}, errors.New("can't parse required Mailgun config mailgun_sender")
				}
				parsedMailgunRecipients, ok := rconfig.Options["mailgun_recipients"].([]interface{})
				if !ok {
					return ParsedConfig{}, errors.New("can't parse required Mailgun config mailgun_recipients")
				}
				var mailgunRecipients []string = make([]string, len(parsedMailgunRecipients))
				for i, pmr := range parsedMailgunRecipients {
					mailgunRecipients[i] = pmr.(string)
				}
				var mailgunExclude []string
				parsedMailgunExclude, ok := rconfig.Options["mailgun_exclude"].([]interface{})
				if ok {
					mailgunExclude = make([]string, len(parsedMailgunExclude))
					for i, pme := range parsedMailgunExclude {
						mailgunExclude[i] = pme.(string)
					}
				} else {
					mailgunExclude = make([]string, 0)
				}
				config.ExportedReporters = append(
					config.ExportedReporters,
					MailgunReporter{
						ServerConfig:      config.Server,
						MailgunDomain:     mailgunDomain,
						MailgunPrivateKey: mailgunPrivateKey,
						MailgunSender:     mailgunSender,
						MailgunRecipients: mailgunRecipients,
						MailgunExclude:    mailgunExclude,
					},
				)
			}
		}
	}
	return config, nil
}
