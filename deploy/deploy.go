package deploy

import (
	"html/template"
	"bytes"
	"fmt"
	"time"
	"os"
	"runtime"
	"os/exec"
	"strings"
)

type Result struct {
	Code    int    `json:"code"`
	Success bool   `json:"isSuccess"`
	Message string `json:"message"`
}

const TEMPLATE = "https://jitpack.io/com/github/{{.User}}/{{.Repo}}/{{.Version}}/{{.Repo}}-{{.Version}}.jar"
const DEPLOY_SYMBOLIC = "{{.Repo}}.jar"
const DEPLOY_FILE = "{{.Repo}}-{{.Version}}-{{.Time}}.jar"

var DEPLOY_ROOT string

func Deploy(user, repo, version, deployRoot string) (*Result, error) {
	if _, err := os.Stat(deployRoot); os.IsNotExist(err) {
		return &Result{Code: 1, Success: false, Message: err.Error()}, err
	}
	if strings.HasSuffix(deployRoot, "/") {
		DEPLOY_ROOT = deployRoot
	} else {
		DEPLOY_ROOT = deployRoot + "/"
	}
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
	deploySymbolic, err := formatDeploySymbolic(repo)
	if err != nil {
		return &Result{Code: 1, Success: false, Message: err.Error()}, err
	}
	deployFile, err := formatDeployFile(repo, version, formatTimeStamp(time.Now()))
	if err != nil {
		return &Result{Code: 1, Success: false, Message: err.Error()}, err
	}
	doDeploy(deploySymbolic, tempPath, deployFile, repo)
	if err != nil {
		return &Result{Code: 1, Success: false, Message: err.Error()}, err
	}
	err = execScript()
	if err != nil {
		return &Result{Code: 1, Success: false, Message: err.Error()}, err
	} else {
		return &Result{Code: 0, Success: true, Message: fullAddr + " 部署成功！"}, nil
	}
	return &Result{Code: 1, Success: false, Message: "Internal failure"}, nil
}

func execScript() error {
	cmd := exec.Command("/bin/bash", "-c", DEPLOY_ROOT+"stop_health_diet.sh")
	fmt.Println("executing" + strings.Join(cmd.Args, " "))
	cmd.Stdout = os.Stdout
	err := cmd.Run()

	if err != nil {
		return err
	}
	cmd = exec.Command("/bin/bash", "-c", DEPLOY_ROOT+"start_health_diet.sh")
	fmt.Println("executing" + strings.Join(cmd.Args, " "))
	cmd.Stdout = os.Stdout
	err = cmd.Run()

	if err != nil {
		return err
	}
	return nil
}

func doDeploy(deploySymbolic, tempPath, deployFile, repo string) error {
	fileInfo, err := os.Lstat(deploySymbolic)
	if err != nil {
		return err
	}
	if (fileInfo.Mode() & os.ModeSymlink) != 0 {
		//os.Readlink(deployRoot) to get the real path
		if runtime.GOOS == "windows" {
			os.Rename(tempPath, "D:"+DEPLOY_ROOT+deployFile)
			os.MkdirAll("D:"+DEPLOY_ROOT, 0755)
		} else {
			os.Rename(tempPath, DEPLOY_ROOT+deployFile)
			os.MkdirAll(DEPLOY_ROOT, 0755)
			realPath, _ := os.Readlink(deploySymbolic)
			os.Remove(deploySymbolic + ".previous")
			os.Symlink(realPath, deploySymbolic+".previous")
		}
		os.Remove(deploySymbolic)
		os.Symlink(deployFile, deploySymbolic)
	} else {
		os.Rename(deploySymbolic, deploySymbolic+"-"+formatTimeStamp(time.Now())+".bak")
		if runtime.GOOS == "windows" {
			os.MkdirAll("D:"+DEPLOY_ROOT, 0755)
			os.Rename(tempPath, "D:"+DEPLOY_ROOT+deployFile)
		} else {
			os.MkdirAll(DEPLOY_ROOT, 0755)
			os.Rename(tempPath, DEPLOY_ROOT+deployFile)
		}
		os.Symlink(deployFile, deploySymbolic)
	}
	return nil
}

func formatTimeStamp(t time.Time) string {
	return fmt.Sprintf("%d%02d%02d-%02d%02d%02d",
		t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
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

func formatDeploySymbolic(repo string) (string, error) {
	var tpl bytes.Buffer
	urlTemplate := template.New("DeployRoot template")
	if runtime.GOOS == "windows" {
		urlTemplate.Parse("D:" + DEPLOY_ROOT + DEPLOY_SYMBOLIC)
	} else {
		urlTemplate.Parse(DEPLOY_ROOT + DEPLOY_SYMBOLIC)
	}
	err := urlTemplate.Execute(&tpl, struct {
		Repo string
	}{repo})
	return tpl.String(), err
}

func formatDeployFile(repo, version, time string) (string, error) {
	var tpl bytes.Buffer
	urlTemplate := template.New("DeployRootFile template")
	urlTemplate.Parse(DEPLOY_FILE)
	err := urlTemplate.Execute(&tpl, struct {
		Repo    string
		Version string
		Time    string
	}{repo, version, time})
	return tpl.String(), err
}
