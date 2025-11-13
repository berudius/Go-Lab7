package repositories

import (
	"encoding/json"
	"os"
	"go.mod/models"
)



func LoadFromFile[T models.Saveable](path string) ([]T, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {

		if err := SaveToFile([]models.Saveable{}, path); err != nil {
			return nil, err
		}
		return []T{}, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var instances []T
	if len(data) == 0 {
		return []T{}, nil
	}

	err = json.Unmarshal(data, &instances)
	if err != nil {
		return nil, err
	}
	return instances, nil
}

func SaveToFile(instances []models.Saveable, path string) error {
	data, err := json.MarshalIndent(instances, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func AppendToFile(instancesToAdd []models.Saveable, path string) error {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, instance := range instancesToAdd {
		data, err := json.Marshal(instance)
		if err != nil {
			return err
		}

		data = append(data, '\n')

		if _, err := file.Write(data); err != nil {
			return err
		}
	}

	return nil
}
