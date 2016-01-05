package main

import (
	"fmt"
	gdClient "github.com/gitDashboard/client/v1"
	/*	"log"*/
	"os"
	"os/exec"
	"strings"
)

func main() {

	/*logF, _ := os.OpenFile("/tmp/gd_backend.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	lg := log.New(logF, "gd_backend", log.Ldate|log.Ltime)
	lg.Println("Environ:", os.Environ())
	lg.Println("Args:", os.Args)*/

	gdUrl := os.Getenv("GIT_DASHBOARD_URL")
	repoBaseDir := os.Getenv("GIT_PROJECT_ROOT")
	username := os.Getenv("REMOTE_USER")
	pathInfo := os.Getenv("PATH_INFO")
	//checking permission
	infoUrlPos := strings.LastIndex(pathInfo, "/info")
	var repoPath string
	if infoUrlPos != -1 {
		repoPath = repoBaseDir + pathInfo[0:infoUrlPos]
	} else {
		uploadPackUrlPos := strings.LastIndex(pathInfo, "/git-upload-pack")
		if uploadPackUrlPos != -1 {
			repoPath = repoBaseDir + pathInfo[0:uploadPackUrlPos]
		} else {
			receivePackUrlPos := strings.LastIndex(pathInfo, "/git-receive-pack")
			if receivePackUrlPos != -1 {
				repoPath = repoBaseDir + pathInfo[0:receivePackUrlPos]
			}
		}
	}
	/*lg.Println("Git repository path:", repoPath)*/
	gdCl := &gdClient.GDClient{Url: gdUrl}
	authorized, _ := gdCl.CheckAuthorization(username, repoPath, "/", "read")
	if !authorized {
		fmt.Println("Status:403\n")
	} else {
		//exec of git-http-backend
		gitHttpBackendPath := os.Getenv("GIT_BACKEND")
		gitBackendCmd := exec.Command(gitHttpBackendPath)
		gitBackendCmd.Stdin = os.Stdin
		gitBackendCmd.Stderr = os.Stderr
		gitBackendCmd.Stdout = os.Stdout
		gitBackendCmd.Run()
	}
	/*logF.Close()*/
}
