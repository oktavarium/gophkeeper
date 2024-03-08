package storage

import (
	"fmt"
	"time"

	"github.com/oktavarium/gophkeeper/internal/shared/models"
)

func (s *Storage) GetData() (
	map[string]models.SimpleCardData,
	map[string]models.SimpleData,
	map[string]models.SimpleBinaryData,
	error,
) {

	cd, err := s.getCards()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error on getting cards: %w", err)
	}

	sd, err := s.getSimple()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error on getting simple data: %w", err)
	}

	sb, err := s.getBinary()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error on getting binary data: %w", err)
	}

	return cd, sd, sb, nil
}

func (s *Storage) DeleteData(id string) error {
	if !s.isInited() {
		return fmt.Errorf("storage is not inited")
	}

	if err := s.store.Write(func(data *storageModel) error {
		if _, ok := data.SimpleCardData[id]; ok {
			data.SimpleCardData[id] = simpleCardData{
				Common: commonData{
					Modified: time.Now().UTC(),
					Deleted:  true,
					Type:     Card,
				},
			}
		} else if _, ok := data.SimpleData[id]; ok {
			data.SimpleData[id] = simpleData{
				Common: commonData{
					Modified: time.Now().UTC(),
					Deleted:  true,
					Type:     Simple,
				},
			}
		} else if _, ok := data.SimpleBinaryData[id]; ok {
			data.SimpleBinaryData[id] = simpleBinaryData{
				Common: commonData{
					Modified: time.Now().UTC(),
					Deleted:  true,
					Type:     Binary,
				},
			}
		}

		return nil
	}); err != nil {
		return fmt.Errorf("error on deleting data: %w", err)
	}

	return nil
}
