package main

import (
	"errors"
	"fmt"
	gdClient "github.com/gitDashboard/client/v1"
	"github.com/gitDashboard/client/v1/misc"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var cl *gdClient.GDClient

func checkAuthorization(user, refName, oldRev, newRev, gitDir, operation string) error {
	granted, locked, err := cl.CheckAuthorization(user, gitDir, refName, "commit")
	if err != nil {
		return err
	}
	if !granted {
		return errors.New("Not authorized to " + operation + " on  " + refName)
	}
	if locked {
		return errors.New("Repository locked, please contact the git administrator")
	}
	return nil
}

func main() {
	os.Exit(mainExec())
}

func mainExec() int {
	refName := os.Args[1]
	oldRev := os.Args[2]
	newRev := os.Args[3]
	gitDir := os.Getenv("GIT_DIR")
	remoteUser := os.Getenv("REMOTE_USER")
	var err error
	var eventId uint
	var gdUrl string
	if len(gitDir) == 0 {
		fmt.Println("Error GIT_DIR env not found")
		return 1
	}
	gitDir, _ = filepath.Abs(gitDir)

	gdUrl = os.Getenv("GIT_DASHBOARD_URL")
	if gdUrl == "" {
		fmt.Println("Error GIT_DASHBOARD_URL env not found")
		return 1
	}
	//understand operation
	zero := "0000000000000000000000000000000000000000"
	var operation string

	if newRev == zero {
		operation = "delete"
	} else {
		operationOut, err := exec.Command("git", "cat-file", "-t", newRev).Output()
		if err != nil {
			goto fatal
		}
		operation = strings.TrimSuffix(string(operationOut), "\n")
	}

	cl = &gdClient.GDClient{Url: gdUrl}
	eventId, err = cl.StartEvent(gitDir, operation, remoteUser, refName, "", misc.EventLevel_INFO)
	if err != nil {
		goto fatal
	}
	defer cl.FinishEvent(eventId)

	err = checkAuthorization(remoteUser, refName, oldRev, newRev, gitDir, operation)
	if err != nil {
		goto fatal
	} else {
		return 0
	}
fatal:
	if err != nil {
		cl.AddEvent(gitDir, operation, remoteUser, refName, "Error:"+err.Error(), misc.EventLevel_ERROR)
		fmt.Println("Error:", err.Error())
		return 1
	}
	return 1
}
