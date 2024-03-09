package storage

import (
	"encoding/json"
	"fmt"
	"maps"

	"github.com/oktavarium/gophkeeper/internal/shared/models"
)

func (s *Storage) GetDataEncrypted() (map[string]models.SimpleDataEncrypted, error) {
	result := make(map[string]models.SimpleDataEncrypted)
	enryptedCards, err := s.getCardsEncrypted()
	if err != nil {
		return nil, fmt.Errorf("error on getting encrypted cards: %w", err)
	}

	enryptedSimple, err := s.getSimpleEncrypted()
	if err != nil {
		return nil, fmt.Errorf("error on getting encrypted simple data: %w", err)
	}

	enryptedBinary, err := s.getBinaryEncrypted()
	if err != nil {
		return nil, fmt.Errorf("error on getting encrypted simple data: %w", err)
	}

	maps.Copy(result, enryptedCards)
	maps.Copy(result, enryptedSimple)
	maps.Copy(result, enryptedBinary)

	return result, nil
}

func (s *Storage) UpdateDataEncrypted(cards map[string]models.SimpleDataEncrypted) error {
	if !s.isInited() {
		return fmt.Errorf("storage is not inited")
	}

	if err := s.store.Write(func(data *storageModel) error {
		for k, v := range cards {
			decryptedData, err := s.crypto.DecryptData(v.Data)
			if err != nil {
				return fmt.Errorf("error on decrypting data: %w", err)
			}

			switch t := v.Common.Type; t {
			case models.Simple:
				record := &models.SimpleRecord{}
				if err := json.Unmarshal(decryptedData, record); err != nil {
					return fmt.Errorf("error on unmarshaling data: %w", err)
				}

				data.SimpleData[k] = simpleData{
					Common: commonData{
						Type:     DataType(t),
						Modified: v.Common.Modified,
						Deleted:  v.Common.Deleted,
					},
					Data: simpleRecord{
						Name: record.Name,
						Data: record.Data,
					},
				}
			case models.Card:
				cardRecord := &models.SimpleCardRecord{}
				if err := json.Unmarshal(decryptedData, cardRecord); err != nil {
					return fmt.Errorf("error on unmarshaling data: %w", err)
				}

				data.SimpleCardData[k] = simpleCardData{
					Common: commonData{
						Type:     DataType(t),
						Modified: v.Common.Modified,
						Deleted:  v.Common.Deleted,
					},
					Data: simpleCardRecord{
						Name:       cardRecord.Name,
						Number:     cardRecord.Number,
						CVV:        cardRecord.CVV,
						ValidUntil: cardRecord.ValidUntil,
					},
				}
			case models.Binary:
				record := &models.SimpleBinaryRecord{}
				if err := json.Unmarshal(decryptedData, record); err != nil {
					return fmt.Errorf("error on unmarshaling data: %w", err)
				}

				data.SimpleBinaryData[k] = simpleBinaryData{
					Common: commonData{
						Type:     DataType(t),
						Modified: v.Common.Modified,
						Deleted:  v.Common.Deleted,
					},
					Data: simpleBinaryRecord{
						Name: record.Name,
						Data: record.Data,
					},
				}
			}

		}

		return nil
	}); err != nil {
		return fmt.Errorf("error on saving data: %w", err)
	}

	return nil
}
