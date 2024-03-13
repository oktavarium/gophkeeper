package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"

	"github.com/oktavarium/gophkeeper/internal/shared/models"
)

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

func (s *Storage) getBinaryEncrypted() (map[string]models.SimpleDataEncrypted, error) {
	if !s.isInited() {
		return nil, fmt.Errorf("storage is not inited")
	}

	records := make(map[string]models.SimpleBinaryData)
	s.store.Read(func(data *storageModel) {
		for k, v := range data.SimpleBinaryData {
			records[k] = models.SimpleBinaryData{
				Common: models.CommonData{
					Type:     models.Binary,
					Deleted:  v.Common.Deleted,
					Modified: v.Common.Modified,
				},
				Data: models.SimpleBinaryRecord{
					Name: v.Data.Name,
					Data: v.Data.Data,
				},
			}
		}
	})

	encryptedRecords := make(map[string]models.SimpleDataEncrypted, len(records))
	for k, v := range records {
		binaryData, err := json.Marshal(
			&models.SimpleBinaryRecord{
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
				Type:     v.Common.Type,
				Deleted:  v.Common.Deleted,
				Modified: v.Common.Modified,
			},
			Data: encryptedData,
		}
	}

	return encryptedRecords, nil
}

func (s *Storage) SaveFile(id string) error {
	if !s.isInited() {
		return fmt.Errorf("storage is not inited")
	}

	var name string
	var binary []byte
	s.store.Read(func(data *storageModel) {
		name = data.SimpleBinaryData[id].Data.Name
		binary = data.SimpleBinaryData[id].Data.Data
	})

	if err := os.WriteFile(name, binary, 0600); err != nil {
		return fmt.Errorf("error on saving file: %w", err)
	}

	return nil
}
