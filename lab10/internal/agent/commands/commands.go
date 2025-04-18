package commands

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func Execute(cmd string) (string, error) {
	cmd = strings.TrimSpace(cmd)

	switch cmd {
	case "pwd":
		return Pwd()
	case "hostname":
		return Hostname()
	case "whoami":
		return Whoami()
	default:
		return "", fmt.Errorf("unknown command: %s", cmd)
	}
}

func Pwd() (string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return pwd, nil
}

func Hostname() (string, error) {
	host, err := os.Hostname()
	if err != nil {
		return "", err
	}
	return host, nil
}

func Whoami() (string, error) {
	cmd := exec.Command("whoami")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}
