package storage

import (
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
)

func (s *Storage) DeleteBinary(id string) error {
	if !s.isInited() {
		return fmt.Errorf("storage is not inited")
	}

	if err := s.store.Write(func(data *storageModel) error {
		data.SimpleBinaryData[id] = simpleBinaryData{
			Common: commonData{
				Modified: time.Now().UTC(),
				Deleted:  true,
				Type:     Binary,
			},
		}
		return nil
	}); err != nil {
		return fmt.Errorf("error on deleting data: %w", err)
	}

	return nil
}

func (s *Storage) UpsertBinary(id string, name string, path string) error {
	if !s.isInited() {
		return fmt.Errorf("storage is not inited")
	}

	b, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("error on opening file: %w", err)
	}

	newID := id
	if len(newID) == 0 {
		newID = uuid.New().String()
	}

	record := simpleBinaryData{
		Common: commonData{
			Type:     Binary,
			Modified: time.Now().UTC(),
		},
		Data: simpleBinaryRecord{
			Name: name,
			Data: b,
		},
	}

	_ = s.store.Write(func(data *storageModel) error {
		data.SimpleBinaryData[newID] = record
		return nil
	})

	return nil
}
