package settings

import (
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"
)

const AgentSettingsVersion = "v1.0"

type AgentSettings struct {
	path             string
	Version          string `json:"version"`
	AgentID          int    `json:"agent_id"`
	ControllerHostWS string `json:"controller_host_ws"`
}

func createDefaultAgentSettings(path string) *AgentSettings {
	return &AgentSettings{
		path:             path,
		Version:          AgentSettingsVersion,
		AgentID:          0,
		ControllerHostWS: "ws://localhost:42066/ws",
	}
}

func (s *AgentSettings) Save() error {
	data, err := json.MarshalIndent(s, "", " ")
	if err != nil {
		return err
	}
	err = os.WriteFile(s.path, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func GetAgentPath() (string, error) {
	path := "./"
	sPath := filepath.Join(path, "/agent.json")
	return sPath, nil
}

func LoadAgentSettings(path string) (*AgentSettings, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			slog.Info("No agent config found. Creating default one.", "path", path)
			s := createDefaultAgentSettings(path)
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

	s.path = path

	return &s, nil
}
