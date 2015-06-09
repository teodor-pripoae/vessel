package main

import (
  "fmt"
  "log"
  "os"
  "os/exec"
  "strings"
)

func Deploy(slug_path string, config Config, app AppConfig) {
  for _, server := range *app.Deploy.UploadServers {
    copyDeploySlug(slug_path, server, config, app)
  }

  for _, entry := range *app.Deploy.Services {
    split := strings.Split(entry, "/")

    if len(split) < 2 {
      log.Fatalf("Please enter service in the format: myuser@myapp.com/my_service")
    }
    restartService(split[0], split[1], config, app)
  }
}

func copyDeploySlug(slug_path string, server string, config Config, app AppConfig) {
  if app.Deploy.SlugLocation == nil {
    log.Fatalf("Deploy slug location not defined")
  }
  destination := fmt.Sprintf("%s:%s", server, *app.Deploy.SlugLocation)

  cmd := exec.Command("scp", slug_path, destination)
  cmd.Stdout = os.Stdout
  cmd.Stderr = os.Stderr
  cmd.Run()
}

func restartService(server string, service string, config Config, app AppConfig) {
  restart_cmd := fmt.Sprintf("sudo service %s restart", service)
  cmd := exec.Command("ssh", server, restart_cmd)
  cmd.Stdout = os.Stdout
  cmd.Stderr = os.Stderr
  cmd.Run() 
}