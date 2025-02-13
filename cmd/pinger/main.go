package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Hanagasumiiii/docker-track/internal/models"
	"github.com/go-ping/ping"
	"log"
	"net/http"
	"time"
)

func getContainersFromAPI() ([]models.Container, error) {
	url := "http://container_api:8081/containers/get"

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get containers from backend API: %v", err)
	}
	defer resp.Body.Close()

	var containers []models.Container
	if err = json.NewDecoder(resp.Body).Decode(&containers); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return containers, nil
}

func pingContainer(ip string) (string, error) {
	pinger, err := ping.NewPinger(ip)
	if err != nil {
		return "", fmt.Errorf("failed to create pinger: %v", err)
	}

	pinger.Count = 1
	pinger.Timeout = 2 * time.Second

	err = pinger.Run()
	if err != nil {
		return "", fmt.Errorf("ping failed: %v", err)
	}

	return pinger.Statistics().AvgRtt.String(), nil
}

func updateContainerStatus(container models.Container) error {
	url := "http://container_api:8081/containers/update"

	data, err := json.Marshal(container)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed with status: %s", resp.Status)
	}

	return nil
}

func main() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			containers, err := getContainersFromAPI()
			if err != nil {
				log.Printf("Error getting containers from API: %v", err)
				continue
			}

			for _, container := range containers {
				_, err = pingContainer(container.Ip)
				if err != nil {
					log.Printf("Failed to ping container %s: %v", container.Ip, err)
					container.Status = "inactive"
				} else {
					container.Status = "active"
				}

				err = updateContainerStatus(container)
				if err != nil {
					log.Printf("Failed to update data for container %s: %v", container.Ip, err)
				} else {
					log.Printf("Successfully updated status for container %s", container.Ip)
				}
			}
		}
	}
}
