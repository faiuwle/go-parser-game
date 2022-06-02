package rage

import (
	"encoding/json"
	"io"
)

func ReadConfig(reader io.Reader) (*GameData, error) {
	configData, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	var gameData GameData
	err = json.Unmarshal(configData, &gameData)
	if err != nil {
		return nil, err
	}
	return &gameData, nil
}
