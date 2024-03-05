package storage

import (
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"

	"github.com/oktavarium/gophkeeper/internal/shared/models"
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

func (s *Storage) getBinary() (map[string]models.SimpleBinaryData, error) {
	if !s.isInited() {
		return nil, fmt.Errorf("storage is not inited")
	}

	records := make(map[string]models.SimpleBinaryData)
	s.store.Read(func(data *storageModel) {
		for k, v := range data.SimpleBinaryData {
			if v.Common.Deleted {
				continue
			}
			records[k] = models.SimpleBinaryData{
				Common: models.CommonData{
					Type:     models.Binary,
					Deleted:  v.Common.Deleted,
					Modified: v.Common.Modified,
				},
				Data: models.SimpleBinaryRecord{
					Name: v.Data.Name,
				},
			}
		}
	})

	return records, nil
}
