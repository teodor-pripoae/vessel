package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

// Build builds an app using config and returns slug_path
func Build(config Config, app AppConfig) string {
	checkStdin()

	// this calls slug builder and get container id
	containerID := buildSlug(config, app)
	// wait for build to finish, showing build output
	waitContainer(containerID)
	// copy slug from container to a temporary location
	slugPath := copySlug(containerID, config, app)
	// delete build container
	deleteContainer(containerID)

	return slugPath
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

func buildSlug(config Config, app AppConfig) string {
	args := slugBuilderCmd(config, app)

	cmd := exec.Command("docker")
	cmd.Args = args
	cmd.Stdin = os.Stdin
	containerID, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Println(string(containerID))
		log.Fatalf("Error running slugbuilder %v", err)
	}

	cid := strings.TrimSpace(string(containerID))

	fmt.Println("Container id:", cid)

	return cid
}

func waitContainer(containerID string) {
	cmd := exec.Command("docker", "logs", "-f", containerID)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

// Copy slug.tgz from container to a temporary location
func copySlug(containerID string, config Config, app AppConfig) string {
	tmpDir, err := ioutil.TempDir("", "")

	if err != nil {
		log.Fatalf("Error creating TempDir")
	}

	source := fmt.Sprintf("%s:/tmp/slug.tgz", containerID)
	output, err := exec.Command("docker", "cp", source, tmpDir).CombinedOutput()

	if err != nil {
		log.Fatalf(string(output))
	}

	return path.Join(tmpDir, "slug.tgz")
}

func deleteContainer(containerID string) {
	cmd := exec.Command("docker", "rm", containerID)
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

func slugBuilderAttachEnv(cmd []string, config Config, app AppConfig) []string {
	// -e BUILDPACK_URL= -e BUILDPACK_VENDOR_URL -e COMMIT_HASH= --env-file=...
	if app.Build.Env != nil {
		for _, e := range *app.Build.Env {
			cmd = append(cmd, "-e")
			cmd = append(cmd, e)
		}
	}

	if app.Build.Buildpack != "" {
		cmd = append(cmd, "-e")
		cmd = append(cmd, fmt.Sprintf("BUILDPACK_URL=%s", app.Build.Buildpack))
	}

	if app.Build.BuildpackVendor != "" {
		cmd = append(cmd, "-e")
		cmd = append(cmd, fmt.Sprintf("BUILDPACK_VENDOR_URL=file://%s", app.Build.BuildpackVendor))
	}

	cmd = append(cmd, "-e")
	cmd = append(cmd, fmt.Sprintf("COMMIT_HASH=%s", config.Commit))

	if app.Build.EnvFile != nil {
		cmd = append(cmd, fmt.Sprintf("--env-file=%s", *app.Build.EnvFile))
	}

	return cmd
}

func slugBuilderAttachVolumes(cmd []string, config Config, app AppConfig) []string {
	// -v /etc/buildpacks:/buildpacks -v /var/docker/build/test:/tmp/cache
	if app.Build.Volumes != nil {
		for _, v := range *app.Build.Volumes {
			cmd = append(cmd, "-v")
			cmd = append(cmd, v)
		}
	}
	return cmd
}

func slugBuilderAttachImage(cmd []string, config Config, app AppConfig) []string {
	if app.Build.BuildImage != nil {
		return append(cmd, *app.Build.BuildImage)
	}

	return append(cmd, "flynn/slugbuilder")
}
