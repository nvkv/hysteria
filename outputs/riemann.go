package outputs

import (
	"github.com/amir/raidman"
	. "github.com/semka/hysteria/models"
	"os"
	"time"
)

type Riemann struct {
	URL         string
	Client      *raidman.Client
	ServiceName string
}

func (r *Riemann) Connect() error {
	c, err := raidman.Dial("tcp", r.URL)
	r.Client = c
	return err
}

func (r *Riemann) SendTestResult(res *TestResult) error {

	var state string
	if res.IsPassed {
		state = "ok"
	} else {
		state = "critical"
	}

	hostname, err := os.Hostname()
	if err != nil {
		return err
	}

	event := &raidman.Event{
		State:   state,
		Host:    hostname,
		Service: r.ServiceName,
		Metric:  1,
		Attributes: map[string]string{
			"description": res.LogLine(),
			"test":        res.File.Name(),
		},
	}
	r.Client.Send(event)
	return nil
}

func (r *Riemann) SendBulk(results []TestResult) error {
	for _, res := range results {
		if err := r.SendTestResult(&res); err == nil {
			// A bit hackish way to pass events to riemann.
			// Riemann can process events from the same service with ~1s precision,
			// that means if we push all events simultaneously only one of them
			// will be processed. That's why this stupid delay between push
			// events is necessary.
			time.Sleep(time.Second)
		} else {
			return err
		}
	}
	return nil
}
