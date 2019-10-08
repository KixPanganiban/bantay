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
		}
	}
	return config, nil
}
