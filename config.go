package main

type Config struct {
	Config   string
	Commit   string
	Deployer string
}

type BuildConfig struct {
	Buildpack       string    `toml:"buildpack"`
	BuildpackVendor string    `toml:"buildpack_vendor"`
	BuildImage      *string   `toml:"build_image"`
	EnvFile         *string   `toml:"env_file"`
	Env             *[]string `toml:"env"`
	Volumes         *[]string `toml:"volumes"`
}

type DeployConfig struct {
	SlugLocation  *string   `toml:"slug_location"`
	UploadServers *[]string `toml:"upload_servers"`
	Services      *[]string `toml:"services"`
}

type SlackNotifyConfig struct {
	Team       string `toml:"team"`
	Channel    string `toml:"channel"`
	Username   string `toml:"username"`
	WebhookUrl string `toml:"webhook_url"`
	OnStarted  bool   `toml:"on_started"`
	OnFailure  bool   `toml:"on_failure"`
	OnSuccess  bool   `toml:"on_success"`
}

type OpbeatNotifyConfig struct {
	OrgId     string `toml:"org"`
	AppId     string `toml:"app"`
	AppSecret string `toml:"token"`
	OnStarted bool   `toml:"on_started"`
	OnFailure bool   `toml:"on_failure"`
	OnSuccess bool   `toml:"on_success"`
}

type NotifyConfig struct {
	Slack  *SlackNotifyConfig
	Opbeat *OpbeatNotifyConfig
}

type AppConfig struct {
	App    string
	Stage  string
	Build  *BuildConfig
	Deploy *DeployConfig
	Notify *NotifyConfig
}
