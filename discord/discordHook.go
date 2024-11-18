package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus" // Assumed logging package
	"io/ioutil"
	"net/http"
	"strconv"
)

type WebhookMessage struct {
	Content  string `json:"content"`
	Username string `json:"username"`
}

const (
	apiURL      = "https://admin.npterra.hu/api/application/union/"
	bearerToken = "ptla_SdxLtQjF1j6losVJw7C1QDfS6LueOSfULEUZufeGU9B"
	webhookURL  = "https://discord.com/api/webhooks/1263780637106638871/ytU7ObFrwINIBIvIgDcIqMjNTvO6eSMPEFJar9Wlau_DsSDlh3VWHmrt7iZMiBqmaD8j"
)

type ServerResponse struct {
	Object     string `json:"object"`
	Attributes struct {
		Name string `json:"name"`
	} `json:"attributes"`
}

func SendDiscordMessage(message string, username string) error {
	msg := WebhookMessage{Content: message, Username: username}
	jsonData, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("could not marshal JSON: %w", err)
	}

	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("could not create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("could not send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("received non-204 response: %d", resp.StatusCode)
	}

	return nil
}

func SendStartingState(uuid string) {
	name := GetServerByUUID(uuid)

	err := SendDiscordMessage(":gear:  A(z) **"+name+"** szerver indítása megkezdődött... ", name)

	if err != nil {
	}
}

func SendWingsStarting() {
	err := SendDiscordMessage(":gear:  A **Wings** indítása megkezdődött... ", "Wings")

	if err != nil {
	}
}

func SendWingsStarted() {
	err := SendDiscordMessage(":green_circle: A **Wings** sikeresen elindult!", "Wings")

	if err != nil {
	}
}

func SendStoppingState(uuid string) {
	name := GetServerByUUID(uuid)

	err := SendDiscordMessage(":gear:  A(z) **"+name+"** szerver leállítása megkezdődött... ", name)

	if err != nil {
	}
}

func SendStartedState(uuid string) {
	name := GetServerByUUID(uuid)

	err := SendDiscordMessage(":green_circle: A(z) **"+name+"** szerver sikeresen elindult!", name)

	if err != nil {
	}
}

func SendStoppedState(uuid string) {
	name := GetServerByUUID(uuid)

	err := SendDiscordMessage(":octagonal_sign: A(z) **"+name+"** szerver sikeresen leállt!", name)

	if err != nil {
	}
}

// GetServerByUUID retrieves server information by UUID
func GetServerByUUID(uuid string) string {
	apiEndpoint := apiURL + uuid
	logrus.WithField("path", apiEndpoint).Debug("Sending GET request")

	req, err := http.NewRequest("GET", apiEndpoint, nil)
	if err != nil {
		logrus.WithField("error", err.Error()).Error("Failed to create request")
		return ""
	}

	req.Header.Add("Authorization", "Bearer "+bearerToken)
	req.Header.Add("Accept", "application/json")
	logrus.WithField("Authorization", "Bearer "+bearerToken).WithField("Content-Type", "application/json").Debug("Headers set")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logrus.WithField("error", err.Error()).Error("Failed to send request")
		return ""
	}
	defer resp.Body.Close()
	logrus.WithField("status", resp.StatusCode).Debug("Request sent")

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.WithField("error", err.Error()).Error("Failed to read response body")
		return ""
	}
	logrus.WithField("body", strconv.QuoteToASCII(string(body))).Debug("Response body read")

	// Check for HTML content
	if resp.Header.Get("Content-Type") != "application/json" {
		logrus.WithField("contentType", resp.Header.Get("Content-Type")).Error("Unexpected content type")
		return ""
	}

	var response ServerResponse
	if err := json.Unmarshal(body, &response); err != nil {
		logrus.WithField("body", strconv.QuoteToASCII(string(body))).WithField("error", err.Error()).Error("Failed to parse JSON")
		return ""
	}
	logrus.WithField("serverName", response.Attributes.Name).Debug("JSON parsed successfully")

	return response.Attributes.Name
}
