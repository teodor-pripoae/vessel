package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	cfg "github.com/teodor-pripoae/vessel/client/config"
)

// OnStart is called when a deploy starts
func OnStart(config cfg.Config, app cfg.AppConfig) {
	nc := app.Notify

	if nc == nil {
		return
	}

	if nc.Slack != nil && nc.Slack.OnStarted {
		notifySlack(config, app, "start")
	}

	if nc.Opbeat != nil && nc.Opbeat.OnStarted {
		notifyOpbeat(config, app, "start")
	}
}

// OnFailure is called when a deploy fails
func OnFailure(config cfg.Config, app cfg.AppConfig) {
	nc := app.Notify

	if nc == nil {
		return
	}

	if nc.Slack != nil && nc.Slack.OnFailure {
		notifySlack(config, app, "fail")
	}

	if nc.Opbeat != nil && nc.Opbeat.OnFailure {
		notifyOpbeat(config, app, "fail")
	}
}

// OnSuccess is called when a deploy is finished
func OnSuccess(config cfg.Config, app cfg.AppConfig) {
	nc := app.Notify

	if nc == nil {
		return
	}

	if nc.Slack != nil && nc.Slack.OnSuccess {
		notifySlack(config, app, "finish")
	}

	if nc.Opbeat != nil && nc.Opbeat.OnSuccess {
		notifyOpbeat(config, app, "finish")
	}
}

func notifySlack(config cfg.Config, app cfg.AppConfig, event string) {
	slack := app.Notify.Slack

	message := fmt.Sprintf("Deploy for app %s on %s %sed", app.App, app.Stage, event)

	fmt.Println("slack", slack.Channel, config.Deployer, "-", message)

	payload := map[string]string{
		"channel":    slack.Channel,
		"username":   config.Deployer,
		"text":       message,
		"icon_emoji": ":rocket:",
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Fatalf("Error while generating payload for slack %v", err)
	}

	_, err = http.PostForm(slack.WebhookURL, url.Values{"payload": {string(jsonPayload)}})
	if err != nil {
		log.Fatalf("Failed announcing to Slack %v", err)
	}
}

func notifyOpbeat(config cfg.Config, app cfg.AppConfig, event string) {
	opConfig := app.Notify.Opbeat
	webhookURL := fmt.Sprintf("https://opbeat.com/api/v1/organizations/%s/apps/%s/releases/", opConfig.OrgID, opConfig.AppID)

	payload := map[string]string{
		"rev":    config.Commit,
		"branch": "default",
		"status": "completed",
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Fatalf("Error while generating payload for slack %v", err)
	}

	req, err := http.NewRequest("POST", webhookURL, bytes.NewBufferString(string(jsonPayload)))
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", opConfig.AppSecret))
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed announcing to Opbeat %v", err)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)

	fmt.Println(buf.String())
}
