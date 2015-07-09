package main

// Config keeps settings related to CLI (commit, config file, deployer)
type Config struct {
	Config   string
	Commit   string
	Deployer string
}

// BuildConfig keeps settings related to build from config file
type BuildConfig struct {
	Buildpack       string    `toml:"buildpack"`
	BuildpackVendor string    `toml:"buildpack_vendor"`
	BuildImage      *string   `toml:"build_image"`
	EnvFile         *string   `toml:"env_file"`
	Env             *[]string `toml:"env"`
	Volumes         *[]string `toml:"volumes"`
}

// DeployConfig keeps settings related to deploy from config file
type DeployConfig struct {
	SlugLocation  *string   `toml:"slug_location"`
	UploadServers *[]string `toml:"upload_servers"`
	Services      *[]string `toml:"services"`
}

// SlackNotifyConfig keeps settings for notifying slack
type SlackNotifyConfig struct {
	Team       string `toml:"team"`
	Channel    string `toml:"channel"`
	Username   string `toml:"username"`
	WebhookURL string `toml:"webhook_url"`
	OnStarted  bool   `toml:"on_started"`
	OnFailure  bool   `toml:"on_failure"`
	OnSuccess  bool   `toml:"on_success"`
}

// OpbeatNotifyConfig keeps settings for notifying opbeat
type OpbeatNotifyConfig struct {
	OrgID     string `toml:"org"`
	AppID     string `toml:"app"`
	AppSecret string `toml:"token"`
	OnStarted bool   `toml:"on_started"`
	OnFailure bool   `toml:"on_failure"`
	OnSuccess bool   `toml:"on_success"`
}

// NotifyConfig keeps a list of available notifiers
type NotifyConfig struct {
	Slack  *SlackNotifyConfig
	Opbeat *OpbeatNotifyConfig
}

// AppConfig keeps config file parts toghether(build, deploy, notify)
type AppConfig struct {
	App    string
	Stage  string
	Build  *BuildConfig
	Deploy *DeployConfig
	Notify *NotifyConfig
}
