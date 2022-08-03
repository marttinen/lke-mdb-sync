package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func cleanKubeconfigPath(path string) (string, error) {
	if strings.HasPrefix(path, "~/") {
		homedir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to resolve home dir: %s", err.Error())
		}
		return homedir + path[1:], nil
	}

	return filepath.Abs(path)
}

func equalAllowLists(mdbAllowList, targetAllowList []string) bool {
	if len(mdbAllowList) != len(targetAllowList) {
		return false
	}

	for _, mIP := range mdbAllowList {
		for _, tIP := range targetAllowList {
			if mIP == tIP {
				goto nextIP
			}
		}
		return false
	nextIP:
	}

	return true
}
