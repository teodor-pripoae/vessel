#vessel
======

Simple Heroku clone on your own servers. Based on [gitreceive](https://github.com/progrium/gitreceive), [slugbuilder](https://github.com/flynn/flynn/tree/master/slugbuilder) and [slugrunner](https://github.com/flynn/flynn/tree/master/slugrunner).

This script takes a payload from gitreceive, runs the buildpack based on your configuration, uploads the slug to your servers and then restarts your services.


### Compiling

```bash
$ go get github.com/BurntSushi/toml
$ CGO_ENABLED=0 GOOS=linux go build -a -tags netgo -ldflags '-w' .
$ # Upload vessel executable to your server and put it in /usr/local/bin
```

### Installation

In your gitreceive `receiver` script:

```bash
#!/bin/bash

set -e

config=/etc/vessel/conf/$1.toml

cat | vessel $config $2 $3
```

Next, put your application configuration in `/etc/vessel/conf/myapp.toml`

```toml
app = "Awesome App"
stage = "production"

[build]
buildpack = "https://github.com/teodor-pripoae/heroku-buildpack-ruby.git"
buildpack_vendor = "/buildpacks/heroku-buildpack-ruby"
env_file = "/etc/vessel/env/myapp"
env = [
  "GEOIP_DIRECTORY=/buildpacks/geoip"
]
volumes = [
  "/etc/vessel/buildpacks:/buildpacks",
  "/var/vessel/build/myapp:/tmp/cache:rw"
]

[deploy]
slug_location = "/etc/vessel/myapp/slug_current.tgz"
upload = [
  "ssh://myuser@web1.myapp.com",
  "ssh://myuser@web2.myapp.com",
  "ssh://myuser@worker1.myapp.com",
  "ssh://myuser@worker2.myapp.com",
  "webdav://myserver:8000/myapp/slugs/current.tgz"
]
services = [
  "ssh://myuser@web1.myapp.com/rails",
  "ssh://myuser@web2.myapp.com/rails",
  "ssh://myuser@worker1.myapp.com/sidekiq",
  "ssh://myuser@worker2.myapp.com/sidekiq",
  "mrt://myserver:8080/rails",
  "mrt://myserver:8080/sidekiq"
]

[notify]
  [notify.slack]
    channel = "#notifications"
    webhook_url = "https://hooks.slack.com/services/token/anothertoken/yetanothertoken"
    on_started = true
    on_failure = true
    on_success = true
  [notify.opbeat]
    org = "myorg_id"
    app = "myapp_id"
    token = "mytoken"
    on_started = false
    on_failure = false
    on_success = true
```

On your `web` or `worker` servers you need to have `rails` or `sidekiq` upstart services present. This services need to expose your compiled slug to `slugrunner` container.

Example command to run slugrunner using slug:

```bash
# Start Web
docker run -d --name=myapp -v /etc/vessel/myapp/:/slug -e SLUG_URL=file:///slug/slug_current.tgz --env-file=/etc/vessel/env/myapp.conf -a stdout -a stderr flynn/slugrunner start web

# Start Worker
docker run -d --name=myapp -v /etc/vessel/myapp/:/slug -e SLUG_URL=file:///slug/slug_current.tgz --env-file=/etc/vessel/env/myapp.conf -a stdout -a stderr flynn/slugrunner start worker
```

Example [Procfile](https://devcenter.heroku.com/articles/procfile):

```yaml
web: bin/puma -C config/puma.rb
worker: bin/sidekiq
```

### Todo

- [X] support plugin system for uploaders (SSH, WebDav)
- [X] support plugin system for services (SSH, Marathon)
- [ ] support plugin for S3
- [ ] support plugin for Kubernetes / Systemd, etc
- [ ] fix notifying on failure (currently it does not work properly)
- [ ] fix cleaning up the mess (tmp slug) on failure
- [ ] more notifiers
