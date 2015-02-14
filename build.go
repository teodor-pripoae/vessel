package main

import (
  "fmt"
  "log"
  "os"
  "os/exec"
  "strings"
  "io/ioutil"
  "path"
)

func Build(config Config, app AppConfig) string {
  checkStdin()

  // this calls slug builder and get container id
  container_id := buildSlug(config, app)
  // wait for build to finish, showing build output
  waitContainer(container_id)
  // copy slug from container to a temporary location
  slug_path := copySlug(container_id, config, app)
  // delete build container
  deleteContainer(container_id)

  return slug_path
}

func checkStdin() {
  fi, err := os.Stdin.Stat()
  if err != nil {
    panic(err)
  }

  if (fi.Mode() & os.ModeCharDevice) != 0 {
    log.Fatalf("Stdin is empty, please pipe repo content")
  }
}

func buildSlug(config Config, app AppConfig) string{
  args := slugBuilderCmd(config, app)
    
  cmd := exec.Command("docker")
  cmd.Args = args
  cmd.Stdin = os.Stdin
  container_id, err := cmd.CombinedOutput()

  if err != nil {
    fmt.Println(string(container_id))
    log.Fatalf("Error running slugbuilder", err)
  }

  cid := strings.TrimSpace(string(container_id))

  fmt.Println("Container id:", cid)

  return cid
}

func waitContainer(container_id string) {
  cmd := exec.Command("docker", "logs", "-f", container_id)
  cmd.Stdout = os.Stdout
  cmd.Stderr = os.Stderr
  cmd.Run()
}

// Copy slug.tgz from container to a temporary location
func copySlug(container_id string, config Config, app AppConfig) string {
  tmp_dir, err := ioutil.TempDir("", "")

  if err != nil {
    log.Fatalf("Error creating TempDir")
  }

  source := fmt.Sprintf("%s:/tmp/slug.tgz", container_id)
  output, err := exec.Command("docker", "cp", source, tmp_dir).CombinedOutput()
  
  if err != nil {
    log.Fatalf(string(output))
  }

  slug_path := path.Join(tmp_dir, "slug.tgz")
  return slug_path
}

func deleteContainer(container_id string) {
  cmd := exec.Command("docker", "rm", container_id)
  cmd.Stdout = os.Stdout
  cmd.Stderr = os.Stderr
  cmd.Run()
}

func slugBuilderCmd(config Config, app AppConfig) []string {
  cmd := []string{"docker", "run"}
  cmd = slugBuilderAttachNetwork(cmd)
  cmd = slugBuilderAttachStdin(cmd)
  cmd = slugBuilderAttachEnv(cmd, config, app)
  cmd = slugBuilderAttachVolumes(cmd, config, app)
  cmd = slugBuilderAttachImage(cmd, config, app)

  // fmt.Println(cmd)
  return cmd
}

func slugBuilderAttachNetwork(cmd []string) []string {
  return append(cmd, "--net='host'")
}

func slugBuilderAttachStdin(cmd []string) []string {
  // -i -a stdin
  cmd = append(cmd, "-i")
  cmd = append(cmd, "-a")
  return append(cmd, "stdin")
}

func slugBuilderAttachEnv(cmd []string, config Config, app AppConfig) [] string {
  // -e BUILDPACK_URL= -e BUILDPACK_VENDOR_URL -e COMMIT_HASH= --env-file=...
  if app.Build.Env != nil {
    for _, e := range *app.Build.Env {
      cmd = append(cmd, "-e")
      cmd = append(cmd, e)
    }
  }

  cmd = append(cmd, "-e")
  cmd = append(cmd, fmt.Sprintf("BUILDPACK_URL=%s", app.Build.Buildpack))

  cmd = append(cmd, "-e")
  cmd = append(cmd, fmt.Sprintf("BUILDPACK_VENDOR_URL=file://%s", app.Build.BuildpackVendor))
  
  cmd = append(cmd, "-e")
  cmd = append(cmd, fmt.Sprintf("COMMIT_HASH=%s", config.Commit))

  if app.Build.EnvFile != nil {
    cmd = append(cmd, fmt.Sprintf("--env-file=%s", *app.Build.EnvFile))  
  }

  return cmd
}

func slugBuilderAttachVolumes(cmd [] string, config Config, app AppConfig) [] string {
  // -v /etc/buildpacks:/buildpacks -v /var/docker/build/test:/tmp/cache  
  if app.Build.Volumes != nil {
    for _, v := range *app.Build.Volumes {
      cmd = append(cmd, "-v")
      cmd = append(cmd, v)
    }
  }
  return cmd
}

func slugBuilderAttachImage(cmd [] string, config Config, app AppConfig) [] string {
  return append(cmd, "flynn/slugbuilder")
}