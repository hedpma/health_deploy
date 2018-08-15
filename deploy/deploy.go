package deploy

import (
	"html/template"
	"bytes"
	"fmt"
	"time"
)

type Result struct {
	Code    int    `json:"code"`
	Success bool   `json:"isSuccess"`
	Message string `json:"message"`
}

const TEMPLATE = "https://jitpack.io/com/github/{{.User}}/{{.Repo}}/{{.Version}}/{{.Repo}}-{{.Version}}.jar"
const DEPLOY_ROOT = "/opt/HealthDiet/"

func Deploy(user, repo, version string) (*Result, error) {
	fullAddr, err := formatUrl(user, repo, version)
	if err != nil {
		return &Result{Code: 1, Success: false, Message: err.Error()}, err
	}
	fmt.Println("downloading from " + fullAddr)
	tempPath, err := formatTempPath(repo, version)
	if err != nil {
		return &Result{Code: 1, Success: false, Message: err.Error()}, err
	}
	err = DownloadFile(tempPath, fullAddr)
	if err != nil {
		return &Result{Code: 1, Success: false, Message: err.Error()}, err
	}
	return &Result{Code: 1, Success: false, Message: "Internal failure"}, nil
}

func formatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func formatUrl(user string, repo string, version string) (string, error) {
	var tpl bytes.Buffer
	urlTemplate := template.New("url template")
	urlTemplate.Parse(TEMPLATE)
	err := urlTemplate.Execute(&tpl, struct {
		User    string
		Repo    string
		Version string
	}{user, repo, version})
	return tpl.String(), err
}
