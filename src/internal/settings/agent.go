package settings

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type AgentSettings struct {
	AgentID int `json:"agent_id"`
}

func LoadAgentSettings() (*AgentSettings, error) {
	path, err := getSettingsBasePath()
	if err != nil {
		return nil, err
	}
	sPath := filepath.Join(path, "/agent.json")

	c, err := os.ReadFile(sPath)
	if err != nil {
		if os.IsNotExist(err) {
			return createDefaultAgentSettings(), nil
		}
		return nil, err
	}

	var s AgentSettings
	err = json.Unmarshal(c, &s)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func createDefaultAgentSettings() *AgentSettings {
	return &AgentSettings{}
}
