package main;

import (
  "bytes"
  "fmt"
  "log"
  "net/http"
  "net/url"
  "encoding/json"
)

func NotifyOnStart(config Config, app AppConfig) {
  nc := app.Notify

  if nc == nil {
    return;
  }

  if nc.Slack != nil && nc.Slack.OnStarted {
    NotifySlack(config, app, "start")
  }

  if nc.Opbeat != nil && nc.Opbeat.OnStarted {
    NotifyOpbeat(config, app, "start")
  }
}

func NotifyOnFailure(config Config, app AppConfig) {
  nc := app.Notify

  if nc == nil {
    return;
  }

  if nc.Slack != nil && nc.Slack.OnFailure {
    NotifySlack(config, app, "fail")
  }

  if nc.Opbeat != nil && nc.Opbeat.OnFailure {
    NotifyOpbeat(config, app, "fail")
  }
}

func NotifyOnSuccess(config Config, app AppConfig) {
  nc := app.Notify

  if nc == nil {
    return;
  }

  if nc.Slack != nil && nc.Slack.OnSuccess {
    NotifySlack(config, app, "finish")
  }

  if nc.Opbeat != nil && nc.Opbeat.OnSuccess {
    NotifyOpbeat(config, app, "finish")
  }
}

func NotifySlack(config Config, app AppConfig, event string) {
  slack := app.Notify.Slack

  message := fmt.Sprintf("Deploy for app %s on %s %sed", app.App, app.Stage, event)

  fmt.Println("slack", slack.Channel, config.Deployer, "-", message)

  payload := map[string]string{
    "channel": slack.Channel, 
    "username": config.Deployer, 
    "text": message, 
    "icon_emoji": ":rocket:",
  }

  json_payload, err := json.Marshal(payload)
  if err != nil {
    log.Fatalf("Error while generating payload for slack", err)
  }

  _, err = http.PostForm(slack.WebhookUrl, url.Values{"payload": {string(json_payload)}})
  if err != nil {
    log.Fatalf("Failed announcing to Slack", err)
  }
}

func NotifyOpbeat(config Config, app AppConfig, event string) {
  op_config := app.Notify.Opbeat
  webhook_url := fmt.Sprintf("https://opbeat.com/api/v1/organizations/%s/apps/%s/releases/", op_config.OrgId, op_config.AppId)
  
  payload  := map[string] string{
    "rev": config.Commit, 
    "branch": "default", 
    "status": "completed",
  }

  json_payload, err := json.Marshal(payload)
  if err != nil {
    log.Fatalf("Error while generating payload for slack", err)
  }
  
  req, err := http.NewRequest("POST", webhook_url, bytes.NewBufferString(string(json_payload)))
  req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", op_config.AppSecret))
  req.Header.Add("Content-Type", "application/json")

  client := &http.Client{}
  resp, err := client.Do(req)
  if err != nil {
    log.Fatalf("Failed announcing to Opbeat", err)
  }

  buf := new(bytes.Buffer)
  buf.ReadFrom(resp.Body)

  fmt.Println(buf.String())
}