package main

import (
	"./core"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var gdConfig *core.GDConfig

func commitOnBranch(refName, oldRev, newRev, gitDir string) error {

	return nil
}

func main() {
	refName := os.Args[1]
	oldRev := os.Args[2]
	newRev := os.Args[3]
	gitDir := os.Getenv("GIT_DIR")

	if len(gitDir) == 0 {
		fmt.Println("Error GIT_DIR env not found")
		os.Exit(1)
	}
	gitDir, _ = filepath.Abs(gitDir)
	//understand newRevType
	zero := "0000000000000000000000000000000000000000"
	var newRevType string

	if newRev == zero {
		newRevType = "delete"
	} else {
		newRevTypeOut, err := exec.Command("git", "cat-file", "-t", newRev).Output()
		if err != nil {
			goto fatal
		}
		newRevType = strings.TrimSuffix(string(newRevTypeOut), "\n")
	}
	fmt.Printf("refname: %v, oldrev:%v, newrev:%v type:%v gitDir:%v\n", refName, oldRev, newRev, newRevType, gitDir)
	gdConfig, err := core.ReadGDConf(gitDir + "/" + "gitDashboard.json")
	if err != nil {
		goto fatal
	}
	switch {
	case strings.HasPrefix(refName, "refs/heads/") && newRevType == "commit":
		err := commitOnBranch(refName, oldRev, newRev, gitDir)
		if err != nil {
			goto fatal
		}
		break
	}
fatal:
	if err != nil {
		fmt.Println("Fatal:", err.Error())
		os.Exit(1)
	}
	os.Exit(1)
}
