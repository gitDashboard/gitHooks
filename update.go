package main

import (
	"errors"
	"fmt"
	gdClient "github.com/gitDashboard/client/v1"
	"github.com/gitDashboard/gitHooks/core"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var gdConfig *core.GDConf

func checkAuthorization(user, refName, oldRev, newRev, gitDir, operation string) error {
	cl := &gdClient.GDClient{Url: gdConfig.Url}
	granted, err := cl.CheckAuthorization(user, gitDir, refName, "commit")
	if err != nil {
		return err
	}
	if !granted {
		return errors.New("Not authorized to " + operation + " on  " + refName)
	}
	return nil
}

func main() {
	refName := os.Args[1]
	oldRev := os.Args[2]
	newRev := os.Args[3]
	gitDir := os.Getenv("GIT_DIR")
	remoteUser := os.Getenv("REMOTE_USER")
	var err error
	if len(gitDir) == 0 {
		fmt.Println("Error GIT_DIR env not found")
		os.Exit(1)
	}
	gitDir, _ = filepath.Abs(gitDir)
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
	fmt.Printf("refname: %v, oldrev:%v, newrev:%v type:%v gitDir:%v\n", refName, oldRev, newRev, operation, gitDir)
	gdConfig, err = core.ReadGDConf(gitDir + "/" + "gitDashboard.json")
	if err != nil {
		goto fatal
	}

	err = checkAuthorization(remoteUser, refName, oldRev, newRev, gitDir, operation)
	if err != nil {
		goto fatal
	} else {
		os.Exit(0)
	}
fatal:
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}
	os.Exit(1)
}
