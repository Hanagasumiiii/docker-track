package main

import (
	"bytes"
	"docker-track/internal/models"
	"encoding/json"
	"fmt"
	"github.com/go-ping/ping"
	"log"
	"net/http"
	"time"
)

func getContainersFromAPI() ([]models.Container, error) {
	url := "http://localhost:8080/containers/get"

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

func sendDataToFrontend(containers models.Container) error {
	url := "http://frontend-api/containers" // API frontend для получения данных

	data, err := json.Marshal(containers)

	fmt.Println(data)

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
	ticker := time.NewTicker(1 * time.Minute)
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
				// TODO after PingTime down
				_, err := pingContainer(container.Ip)
				if err != nil {
					log.Printf("Failed to ping container %s: %v", container.Ip, err)
					continue
				}

				//container.PingTime = pingTime

				err = sendDataToFrontend(container)
				if err != nil {
					log.Printf("Failed to send data for container %s: %v", container.Ip, err)
				} else {
					log.Printf("Successfully sent data for container %s", container.Ip)
				}
			}
		}
	}
}
