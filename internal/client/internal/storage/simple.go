package storage

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/oktavarium/gophkeeper/internal/shared/models"
)

func (s *Storage) getSimple() (map[string]models.SimpleData, error) {
	if !s.isInited() {
		return nil, fmt.Errorf("storage is not inited")
	}

	records := make(map[string]models.SimpleData)
	s.store.Read(func(data *storageModel) {
		for k, v := range data.SimpleData {
			if v.Common.Deleted {
				continue
			}
			records[k] = models.SimpleData{
				Common: models.CommonData{
					Type:     models.Simple,
					Deleted:  v.Common.Deleted,
					Modified: v.Common.Modified,
				},
				Data: models.SimpleRecord{
					Name: v.Data.Name,
					Data: v.Data.Data,
				},
			}
		}
	})

	return records, nil
}

func (s *Storage) UpsertSimple(id string, name string, data string) error {
	if !s.isInited() {
		return fmt.Errorf("storage is not inited")
	}

	newID := id
	if len(newID) == 0 {
		newID = uuid.New().String()
	}

	record := simpleData{
		Common: commonData{
			Type:     Simple,
			Modified: time.Now().UTC(),
		},
		Data: simpleRecord{
			Name: name,
			Data: data,
		},
	}

	_ = s.store.Write(func(data *storageModel) error {
		data.SimpleData[newID] = record
		return nil
	})

	return nil
}

func (s *Storage) DeleteSimple(id string) error {
	if !s.isInited() {
		return fmt.Errorf("storage is not inited")
	}

	if err := s.store.Write(func(data *storageModel) error {
		data.SimpleData[id] = simpleData{
			Common: commonData{
				Modified: time.Now().UTC(),
				Deleted:  true,
				Type:     Card,
			},
		}
		return nil
	}); err != nil {
		return fmt.Errorf("error on deleting data: %w", err)
	}

	return nil
}

func (s *Storage) getSimpleEncrypted() (map[string]models.SimpleDataEncrypted, error) {
	if !s.isInited() {
		return nil, fmt.Errorf("storage is not inited")
	}

	records := make(map[string]models.SimpleData)
	s.store.Read(func(data *storageModel) {
		for k, v := range data.SimpleData {
			records[k] = models.SimpleData{
				Common: models.CommonData{
					Type:     models.Simple,
					Deleted:  v.Common.Deleted,
					Modified: v.Common.Modified,
				},
				Data: models.SimpleRecord{
					Name: v.Data.Name,
					Data: v.Data.Data,
				},
			}
		}
	})

	encryptedRecords := make(map[string]models.SimpleDataEncrypted, len(records))
	for k, v := range records {
		binaryData, err := json.Marshal(
			&models.SimpleRecord{
				Name: v.Data.Name,
				Data: v.Data.Data,
			})
		if err != nil {
			return nil, fmt.Errorf("error on marshaling data: %w", err)
		}

		encryptedData, err := s.crypto.EncryptData(binaryData)
		if err != nil {
			return nil, fmt.Errorf("error on encrypting data: %w", err)
		}

		encryptedRecords[k] = models.SimpleDataEncrypted{
			Common: models.CommonData{
				Deleted:  v.Common.Deleted,
				Modified: v.Common.Modified,
			},
			Data: encryptedData,
		}
	}

	return encryptedRecords, nil
}
