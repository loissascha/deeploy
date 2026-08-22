package settings

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type AgentSettings struct {
	AgentID int `json:"agent_id"`
}

func (s *AgentSettings) Save() error {
	return nil
}

func LoadAgentSettings() (*AgentSettings, error) {
	path, err := getSettingsBasePath()
	if err != nil {
		return nil, err
	}
	err = os.MkdirAll(path, 0644)
	if err != nil {
		return nil, err
	}
	sPath := filepath.Join(path, "/agent.json")

	c, err := os.ReadFile(sPath)
	if err != nil {
		if os.IsNotExist(err) {
			s := createDefaultAgentSettings()
			err = s.Save()
			if err != nil {
				return nil, err
			}
			return s, nil
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
	return &AgentSettings{
		AgentID: 1,
	}
}
