package storage

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/oktavarium/gophkeeper/internal/shared/dto"
)

func (s *JsonStorage) SaveNewCard(name string, number string, cvv uint32, validUntil time.Time) error {
	if !s.isInited() {
		return fmt.Errorf("storage is not inited")
	}

	record := simpleCardData{
		Common: commonData{

			Modified: time.Now(),
		},
		Data: simpleCardRecord{
			Name:       name,
			Number:     number,
			CVV:        cvv,
			ValidUntil: validUntil,
		},
	}

	_ = s.store.Write(func(data *storageModel) error {
		if data.SimpleCardData == nil {
			data.SimpleCardData = make(map[string]simpleCardData)
		}
		data.SimpleCardData[uuid.New().String()] = record
		return nil
	})

	return nil
}

func (s *JsonStorage) GetCards() (map[string]dto.SimpleCardData, error) {
	if !s.isInited() {
		return nil, fmt.Errorf("storage is not inited")
	}

	cards := make(map[string]dto.SimpleCardData)
	s.store.Read(func(data *storageModel) {
		for k, v := range data.SimpleCardData {
			cards[k] = dto.SimpleCardData{
				Common: dto.CommonData{
					IsDeleted: v.Common.IsDeleted,
					Modified:  v.Common.Modified,
				},
				Data: dto.SimpleCardRecord{
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

func (s *JsonStorage) GetCardsEncrypted() (map[string]dto.SimpleCardDataEncrypted, error) {
	if !s.isInited() {
		return nil, fmt.Errorf("storage is not inited")
	}

	cards := make(map[string]dto.SimpleCardData)
	s.store.Read(func(data *storageModel) {
		for k, v := range data.SimpleCardData {
			cards[k] = dto.SimpleCardData{
				Common: dto.CommonData{
					IsDeleted: v.Common.IsDeleted,
					Modified:  v.Common.Modified,
				},
				Data: dto.SimpleCardRecord{
					Name:       v.Data.Name,
					Number:     v.Data.Number,
					ValidUntil: v.Data.ValidUntil,
					CVV:        v.Data.CVV,
				},
			}
		}
	})

	encryptedCards := make(map[string]dto.SimpleCardDataEncrypted, len(cards))
	for k, v := range cards {
		binaryData, err := json.Marshal(
			&dto.SimpleCardRecord{
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

		encryptedCards[k] = dto.SimpleCardDataEncrypted{
			Common: dto.CommonData{
				IsDeleted: v.Common.IsDeleted,
				Modified:  v.Common.Modified,
			},
			Data: encryptedData,
		}
	}

	return encryptedCards, nil
}

func (s *JsonStorage) UpdateCardsEncrypted(cards map[string]dto.SimpleCardDataEncrypted) error {
	if !s.isInited() {
		return fmt.Errorf("storage is not inited")
	}

	if err := s.store.Write(func(data *storageModel) error {
		for k, v := range cards {
			decryptedData, err := s.crypto.DecryptData(v.Data)
			if err != nil {
				return fmt.Errorf("error on decrypting data: %w", err)
			}

			cardRecord := &dto.SimpleCardRecord{}
			if err := json.Unmarshal(decryptedData, cardRecord); err != nil {
				return fmt.Errorf("error on unmarshaling data: %w", err)
			}

			data.SimpleCardData[k] = simpleCardData{
				Common: commonData{
					Modified:  v.Common.Modified,
					IsDeleted: v.Common.IsDeleted,
				},
				Data: simpleCardRecord{
					Name:       cardRecord.Name,
					Number:     cardRecord.Number,
					CVV:        cardRecord.CVV,
					ValidUntil: cardRecord.ValidUntil,
				},
			}
		}

		return nil
	}); err != nil {
		return fmt.Errorf("error on saving data: %w", err)
	}

	return nil
}
