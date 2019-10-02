package lib

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/KixPanganiban/bantay/log"

	"github.com/imroc/req"
)

// Check is a set of parameters matched against to see if a service is up
type Check struct {
	Name        string `yaml:"name"`
	URL         string `yaml:"url"`
	ValidStatus int    `yaml:"valid_status"`
	BodyMatch   string `yaml:"body_match"`
}

// RunCheck performs the HTTP request necessary to verify if the given Check is up
func RunCheck(c *Check) error {
	r := req.New()
	res, err := r.Get(c.URL)
	if err != nil {
		log.Warnf("[%s] Request failed: %s", c.Name, err.Error())
		return err
	}
	// Perform check on StatusCode
	response := res.Response()
	responseStatus := response.StatusCode
	if responseStatus != c.ValidStatus {
		errMsg := fmt.Sprintf("[%s] Status mismatch. Expected %d, got %d.", c.Name, c.ValidStatus, responseStatus)
		log.Warn(errMsg)
		return errors.New(errMsg)
	}
	// Perform check on Body
	if len(c.BodyMatch) > 0 {
		responseBuffer := new(bytes.Buffer)
		responseBuffer.ReadFrom(response.Body)
		responseText := responseBuffer.String()
		if !strings.Contains(responseText, c.BodyMatch) {
			errMsg := fmt.Sprintf("[%s] String '%s' not found in body.", c.Name, c.BodyMatch)
			log.Warn(errMsg)
			return errors.New(errMsg)
		}
	}
	return nil
}
