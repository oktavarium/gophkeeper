package storage

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/oktavarium/gophkeeper/internal/shared/models"
)

func (s *Storage) UpsertCard(id string, name string, number string, cvv uint32, validUntil time.Time) error {
	if !s.isInited() {
		return fmt.Errorf("storage is not inited")
	}

	newCardID := id
	if len(newCardID) == 0 {
		newCardID = uuid.New().String()
	}

	record := simpleCardData{
		Common: commonData{
			Type:     Card,
			Modified: time.Now().UTC(),
		},
		Data: simpleCardRecord{
			Name:       name,
			Number:     number,
			CVV:        cvv,
			ValidUntil: validUntil,
		},
	}

	_ = s.store.Write(func(data *storageModel) error {
		data.SimpleCardData[newCardID] = record
		return nil
	})

	return nil
}

func (s *Storage) getCards() (map[string]models.SimpleCardData, error) {
	if !s.isInited() {
		return nil, fmt.Errorf("storage is not inited")
	}

	cards := make(map[string]models.SimpleCardData)
	s.store.Read(func(data *storageModel) {
		for k, v := range data.SimpleCardData {
			if v.Common.Deleted {
				continue
			}
			cards[k] = models.SimpleCardData{
				Common: models.CommonData{
					Type:     models.Card,
					Deleted:  v.Common.Deleted,
					Modified: v.Common.Modified,
				},
				Data: models.SimpleCardRecord{
					Name:       v.Data.Name,
					Number:     v.Data.Number,
					ValidUntil: v.Data.ValidUntil,
					CVV:        v.Data.CVV,
				},
			}
		}
	})

	return cards, nil
}

func (s *Storage) getCardsEncrypted() (map[string]models.SimpleDataEncrypted, error) {
	if !s.isInited() {
		return nil, fmt.Errorf("storage is not inited")
	}

	cards := make(map[string]models.SimpleCardData)
	s.store.Read(func(data *storageModel) {
		for k, v := range data.SimpleCardData {
			cards[k] = models.SimpleCardData{
				Common: models.CommonData{
					Type:     models.Card,
					Deleted:  v.Common.Deleted,
					Modified: v.Common.Modified,
				},
				Data: models.SimpleCardRecord{
					Name:       v.Data.Name,
					Number:     v.Data.Number,
					ValidUntil: v.Data.ValidUntil,
					CVV:        v.Data.CVV,
				},
			}
		}
	})

	encryptedCards := make(map[string]models.SimpleDataEncrypted, len(cards))
	for k, v := range cards {
		binaryData, err := json.Marshal(
			&models.SimpleCardRecord{
				Name:       v.Data.Name,
				Number:     v.Data.Number,
				ValidUntil: v.Data.ValidUntil,
				CVV:        v.Data.CVV,
			})
		if err != nil {
			return nil, fmt.Errorf("error on marshaling data: %w", err)
		}

		encryptedData, err := s.crypto.EncryptData(binaryData)
		if err != nil {
			return nil, fmt.Errorf("error on encrypting data: %w", err)
		}

		encryptedCards[k] = models.SimpleDataEncrypted{
			Common: models.CommonData{
				Type:     v.Common.Type,
				Deleted:  v.Common.Deleted,
				Modified: v.Common.Modified,
			},
			Data: encryptedData,
		}
	}

	return encryptedCards, nil
}
