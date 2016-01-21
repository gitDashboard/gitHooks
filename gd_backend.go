package main

import (
	"fmt"
	gdClient "github.com/gitDashboard/client/v1"
	"github.com/gitDashboard/client/v1/misc"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {

	logF, _ := os.OpenFile("/tmp/gd_backend.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	defer logF.Close()
	lg := log.New(logF, "gd_backend", log.Ldate|log.Ltime)
	lg.Println("Environ:", os.Environ())
	lg.Println("Args:", os.Args)

	gdUrl := os.Getenv("GIT_DASHBOARD_URL")
	//repoBaseDir := os.Getenv("GIT_PROJECT_ROOT")
	username := os.Getenv("REMOTE_USER")
	pathInfo := os.Getenv("PATH_INFO")
	//checking permission
	lg.Println("pathInfo:", pathInfo)
	infoUrlPos := strings.LastIndex(pathInfo, "/info")
	var repoPath string
	if infoUrlPos != -1 {
		repoPath = pathInfo[1:infoUrlPos]
	} else {
		uploadPackUrlPos := strings.LastIndex(pathInfo, "/git-upload-pack")
		if uploadPackUrlPos != -1 {
			repoPath = pathInfo[1:uploadPackUrlPos]
		} else {
			receivePackUrlPos := strings.LastIndex(pathInfo, "/git-receive-pack")
			if receivePackUrlPos != -1 {
				repoPath = pathInfo[1:receivePackUrlPos]
			}
		}
	}
	lg.Println("Git repository path:", repoPath)
	gdCl := &gdClient.GDClient{Url: gdUrl}
	authorized, locked, _ := gdCl.CheckAuthorization(username, repoPath, "/", "read")
	eventId, err := gdCl.StartEvent(repoPath, "access", username, "", pathInfo, misc.EventLevel_INFO)
	//lg.Println("eventId", eventId)
	if err != nil {
		gdCl.AddEvent(repoPath, "access", username, "", err.Error(), misc.EventLevel_ERROR)
		lg.Println("Error:" + err.Error())
		fmt.Println("Status:500\n")
	}
	defer gdCl.FinishEvent(eventId)
	switch {
	case !authorized:
		gdCl.AddEvent(repoPath, "access", username, "", "Unatorized", misc.EventLevel_WARN)
		lg.Println("Status:403")
		fmt.Println("Status:403\n")
	case locked:
		gdCl.AddEvent(repoPath, "access", username, "", "Locked", misc.EventLevel_WARN)
		lg.Println("Status:503")
		fmt.Println("Status:503\n")
	default:
		//exec of git-http-backend
		lg.Println("exec of git-http-backend")
		gitHttpBackendPath := os.Getenv("GIT_BACKEND")
		gitBackendCmd := exec.Command(gitHttpBackendPath)
		gitBackendCmd.Stdin = os.Stdin
		gitBackendCmd.Stderr = os.Stderr
		gitBackendCmd.Stdout = os.Stdout
		gitBackendCmd.Run()
	}

}
