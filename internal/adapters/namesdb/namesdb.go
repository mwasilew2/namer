package namesdb

import (
	"context"
	"embed"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/mwasilew2/go-service-template/internal/domain/models"
)

//go:embed names_transformed.csv
var namesEmbedded embed.FS
var fileWithTransformedNames = "names_transformed.csv"

type NamesDB struct {
	database map[int64]string
}

func NewNamesDB() (*NamesDB, error) {
	db := make(map[int64]string)

	// open embedded file
	fs, err := namesEmbedded.Open(fileWithTransformedNames)
	if err != nil {
		return nil, fmt.Errorf("failed to open embedded %s: %w", fileWithTransformedNames, err)
	}
	r := csv.NewReader(fs)

FILE_READING_LOOP:
	for {
		record, err := r.Read()
		if err != nil {
			switch {
			case err == io.EOF:
				break FILE_READING_LOOP
			default:
				return nil, fmt.Errorf("failed to read embedded names.csv: %w", err)

			}
		}
		id, err := strconv.ParseInt(record[0], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse id %s: %w", record[0], err)
		}
		name := record[1]
		db[id] = name
	}

	return &NamesDB{
		database: db,
	}, nil
}

var ErrNameNotFound = errors.New("name not found")

func (n NamesDB) GetName(ctx context.Context, id int64) (*models.Name, error) {
	name, ok := n.database[id]
	if !ok {
		return nil, ErrNameNotFound
	}
	return &models.Name{
		Id:    id,
		Value: name,
	}, nil
}
