package deploy

import (
	"bytes"
	"html/template"
	"os"
	"net/http"
	"errors"
	"io"
	"runtime"
)

const TEMP_PATH = "/tmp/{{.Repo}}-{{.Version}}.jar"

func formatTempPath(repo string, version string) (string, error) {
	var tpl2 bytes.Buffer
	tempPathTemplate := template.New("temp path template")
	if runtime.GOOS == "windows" {
		tempPathTemplate.Parse("D:" + TEMP_PATH)
	} else {
		tempPathTemplate.Parse(TEMP_PATH)
	}
	err := tempPathTemplate.Execute(&tpl2, struct {
		Repo    string
		Version string
	}{repo, version})
	return tpl2.String(), err
}

func DownloadFile(path string, url string) error {
	// Create the file
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	if resp.StatusCode == http.StatusNotFound {
		return errors.New(url + " not found, please check cleanly.")
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
