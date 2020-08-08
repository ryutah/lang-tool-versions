package main

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/pkg/errors"
)

type versionLoader func() (versions []string, err error)

var versionLoaders = map[string]versionLoader{
	"python":                 loadPythonVersions,
	// "golang":                 loadGoVersions,
	// "node":                   loadNodeVersions,
	// "docker_compose":         loadDockerComposeVersions,
	// "go_task_task":           loadGotaskTask,
	// "golangci_golangci_lint": loadGolangciLint,
	// "google_cloud_sdl":       loadGoogleCloudSDK,
}

func main() {
	results := make(nameVersions)

	for name, loader := range versionLoaders {
		versions, err := loader()
		if err != nil {
			log.Fatalf("%+v", err)
		}
		results[name] = versions
	}

	f, err := os.OpenFile("versions.json", os.O_WRONLY, 0)
	if os.IsNotExist(err) {
		f, err = os.Create("versions.json")
	}
	if err != nil {
		log.Fatalf("%+v", err)
	}
	defer f.Close()

	if err := writeVersions(f, results); err != nil {
		log.Fatalf("%+v", err)
	}
}

type nameVersions map[string][]string

func writeVersions(w io.Writer, versions nameVersions) (err error) {
	defer func() {
		if err != nil {
			err = errors.WithStack(err)
		}
	}()

	body, e := json.MarshalIndent(versions, "", "  ")
	if err = e; err != nil {
		return
	}
	bufw := bufio.NewWriter(w)
	_, err = bufw.Write(body)
	if err != nil {
		return
	}
	err = bufw.Flush()
	return
}
