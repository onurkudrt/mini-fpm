package helper

import (
	"os"
	"path/filepath"
	"strings"
	"errors"
)

func FindRootDir() (string, error) {
	currentDir, err := os.Getwd()

	if err != nil {
		panic("Failed to get current working directory")
	}

	currentIter, iterMax := 0, 20

	for ; currentIter < iterMax; currentIter++ { // while (true)

		if currentDir == "/" || strings.HasSuffix(currentDir, ":\\\\") || strings.HasSuffix(currentDir, "://") {
			break
		}

		if _, err := os.Stat(currentDir + "/go.mod"); err == nil {
			return currentDir, nil
		}

		parent := filepath.Dir(currentDir)

		if parent == currentDir {
			break
		}

		currentDir = parent
	}

	return "", errors.New("project root not found: could not locate go.mod")
}

func GetConfigFilePath() (string, error) {
	rootDir, err := FindRootDir()
	if err != nil {
		return rootDir, err
	}
	
	return filepath.Join(rootDir, "resources", "pm.config"), nil
}