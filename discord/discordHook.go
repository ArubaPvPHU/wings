package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus" // Assumed logging package
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

type WebhookMessage struct {
	Content  string `json:"content"`
	Username string `json:"username"`
}

const (
	apiURL = "https://admin.npterra.hu/api/application/union/"
)

var (
	bearerToken = os.Getenv("BEARER_TOKEN")
	webhookURL  = os.Getenv("WEBHOOK_URL")
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
		log.Println("SendStartingState: Error occured; " + err.Error())
	}
}

func SendWingsStarting() {
	err := SendDiscordMessage(":gear:  A **Wings** indítása megkezdődött... ", "Wings")

	if err != nil {
		log.Println("SendWingsStarting: Error occured; " + err.Error())
	}
}

func SendWingsStarted() {
	err := SendDiscordMessage(":green_circle: A **Wings** sikeresen elindult!", "Wings")

	if err != nil {
		log.Println("SendWingsStarted: Error occured; " + err.Error())
	}
}

func SendStoppingState(uuid string) {
	name := GetServerByUUID(uuid)

	err := SendDiscordMessage(":gear:  A(z) **"+name+"** szerver leállítása megkezdődött... ", name)

	if err != nil {
		log.Println("SendStoppingState: Error occured; " + err.Error())
	}
}

func SendStartedState(uuid string) {
	name := GetServerByUUID(uuid)

	err := SendDiscordMessage(":green_circle: A(z) **"+name+"** szerver sikeresen elindult!", name)

	if err != nil {
		log.Println("SendStartedState: Error occured; " + err.Error())
	}
}

func SendStoppedState(uuid string) {
	name := GetServerByUUID(uuid)

	err := SendDiscordMessage(":octagonal_sign: A(z) **"+name+"** szerver sikeresen leállt!", name)

	if err != nil {
		log.Println("SendStoppedState: Error occured; " + err.Error())
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

func CheckEnvVars() {
	log.Println("WEBHOOK_URL in discord package: %s\n", os.Getenv("WEBHOOK_URL"))
}
