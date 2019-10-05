package lib

import "github.com/KixPanganiban/bantay/log"

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
