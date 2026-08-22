package settings

import (
	"os"
	"path/filepath"
)

func getSettingsBasePath() (string, error) {
	confDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	dir := filepath.Join(confDir, "/deeploy")
	return dir, nil
}
