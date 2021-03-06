package lib

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/KixPanganiban/bantay/log"
	"github.com/KixPanganiban/bantay/version"

	"github.com/imroc/req"
)

// Check is a set of parameters matched against to see if a service is up
type Check struct {
	Name        string `yaml:"name"`
	URL         string `yaml:"url"`
	ValidStatus int    `yaml:"valid_status"`
	BodyMatch   string `yaml:"body_match"`
}

//CheckResult contains a fail/success flag and a message
type CheckResult struct {
	Name    string
	Success bool
	Message string
	Latency time.Duration
}

// RunCheck performs the HTTP request necessary to verify if the given Check is up
func RunCheck(c Check, resChan chan<- CheckResult) {
	r := req.New()
	r.SetFlags(req.Lcost)
	res, err := r.Get(
		c.URL,
		req.Header{"User-Agent": fmt.Sprintf("Bantay %s", version.Version)},
	)
	if err != nil {
		log.Warnf("Request failed: %s", err.Error())
		resChan <- CheckResult{Name: c.Name, Success: false, Message: err.Error(), Latency: 0 * time.Second}
		return
	}
	// Perform check on StatusCode
	response := res.Response()
	responseStatus := response.StatusCode
	if responseStatus != c.ValidStatus {
		errMsg := fmt.Sprintf("Status mismatch. Expected %d, got %d.", c.ValidStatus, responseStatus)
		resChan <- CheckResult{Name: c.Name, Success: false, Message: errMsg, Latency: res.Cost()}
		return
	}
	// Perform check on Body
	if len(c.BodyMatch) > 0 {
		responseBuffer := new(bytes.Buffer)
		responseBuffer.ReadFrom(response.Body)
		responseText := responseBuffer.String()
		if !strings.Contains(responseText, c.BodyMatch) {
			errMsg := fmt.Sprintf("String '%s' not found in body.", c.BodyMatch)
			resChan <- CheckResult{Name: c.Name, Success: false, Message: errMsg, Latency: res.Cost()}
			return
		}
	}
	resChan <- CheckResult{Name: c.Name, Success: true, Latency: res.Cost()}
	return
}

// RunChecks calls RunCheck for every Check provided in slice cs and returns counts for failed, successful, total
func RunChecks(cs *[]Check, r *[]Reporter, downCounter map[string]int) (int, int, int) {
	var (
		failed     int
		successful int
		total      int
	)
	failed, successful = 0, 0
	total = len(*cs)

	resChan := make(chan CheckResult, total)
	for _, c := range *cs {
		go RunCheck(c, resChan)
	}
	func() {
		for i := 0; i < total; i++ {
			res := <-resChan
			for _, reporter := range *r {
				reporter.Report(res, &downCounter)
			}
			if res.Success == true {
				successful++
				downCounter[res.Name] = 0
			} else {
				failed++
				downCounter[res.Name]++
			}
		}
	}()
	return failed, successful, total
}
