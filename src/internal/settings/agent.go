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
	data, err := json.Marshal(s)
	if err != nil {
		return err
	}
	p, err := getAgentPath()
	if err != nil {
		return err
	}
	err = os.WriteFile(p, data, 0755)
	if err != nil {
		return err
	}
	return nil
}

func getAgentPath() (string, error) {
	path, err := getSettingsBasePath()
	if err != nil {
		return "", err
	}
	err = os.MkdirAll(path, 0644)
	if err != nil {
		return "", err
	}
	sPath := filepath.Join(path, "/agent.json")
	return sPath, nil
}

func LoadAgentSettings() (*AgentSettings, error) {
	path, err := getAgentPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
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
	err = json.Unmarshal(data, &s)
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
