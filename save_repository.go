package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

type SaveRepository struct {
	Directory string
}

func (r SaveRepository) Save(save SaveGame) error {
	if save.Meta.PlayerName == "" {
		return os.ErrInvalid
	}

	if err := os.MkdirAll(r.Directory, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(save, "", "  ")
	if err != nil {
		return err
	}

	file := filepath.Join(r.Directory, save.Meta.PlayerName+".json")
	return os.WriteFile(file, data, 0644)
}

func (r SaveRepository) Load(playerName string) (SaveGame, error) {
	var save SaveGame

	file := filepath.Join(r.Directory, playerName+".json")
	data, err := os.ReadFile(file)
	if err != nil {
		return save, err
	}

	if err := json.Unmarshal(data, &save); err != nil {
		return save, err
	}

	return save, nil
}

func (r SaveRepository) Delete(playerName string) error {
	file := filepath.Join(r.Directory, playerName+".json")
	return os.Remove(file)
}

func (r SaveRepository) List() ([]SaveMeta, error) {
	files, err := os.ReadDir(r.Directory)
	if err != nil {
		return nil, err
	}

	var saveMetadata []SaveMeta
	for i, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".json") {
			continue
		}

		save, err := r.Load(strings.TrimSuffix(file.Name(), ".json"))
		if err != nil {
			continue
		}
		if save.Meta.PlayerName == "" {
			save.Meta.PlayerName = save.Player.Name
		}
		if save.Meta.PlayerName == "" {
			continue
		}

		saveMetadata = append(saveMetadata, SaveMeta{
			PlayerName: save.Meta.PlayerName,
			Lvl:        save.Player.Level,
			CreatedAt:  save.Meta.CreatedAt,
			LastPlayed: save.Meta.LastPlayed,
			Index:      i,
		})
	}

	return saveMetadata, nil
}
