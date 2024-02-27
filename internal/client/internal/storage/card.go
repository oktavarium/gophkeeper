package storage

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/oktavarium/gophkeeper/internal/shared/dto"
)

func (s *JsonStorage) SaveNewCard(name string, number string, cvv uint32, validUntil time.Time) error {
	if !s.isInited() {
		return fmt.Errorf("storage is not inited")
	}

	record := simpleCardRecord{
		Common: commonData{
			Name:     name,
			Modified: time.Now(),
		},
		Number:     number,
		CVV:        cvv,
		ValidUntil: validUntil,
	}

	_ = s.store.Write(func(data *storageModel) error {
		if data.SimpleCardData == nil {
			data.SimpleCardData = make(map[string]simpleCardRecord)
		}
		data.SimpleCardData[uuid.New().String()] = record
		return nil
	})

	return nil
}

func (s *JsonStorage) GetCards() (map[string]dto.SimpleCardRecord, error) {
	if !s.isInited() {
		return nil, fmt.Errorf("storage is not inited")
	}

	cards := make(map[string]dto.SimpleCardRecord)
	s.store.Read(func(data *storageModel) {
		for k, v := range data.SimpleCardData {
			cards[k] = dto.SimpleCardRecord{
				Common: dto.CommonData{
					Name:      v.Common.Name,
					IsDeleted: v.Common.IsDeleted,
					Modified:  v.Common.Modified,
				},
				Number:     v.Number,
				ValidUntil: v.ValidUntil,
				CVV:        v.CVV,
			}
		}
	})

	return cards, nil
}
