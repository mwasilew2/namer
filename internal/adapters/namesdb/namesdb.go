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
	maxId    int64
}

func NewNamesDB() (*NamesDB, error) {
	db := make(map[int64]string)

	// open embedded file
	fs, err := namesEmbedded.Open(fileWithTransformedNames)
	if err != nil {
		return nil, fmt.Errorf("failed to open embedded %s: %w", fileWithTransformedNames, err)
	}
	r := csv.NewReader(fs)

	var maxId int64
FILE_READING_LOOP:
	for {
		record, err := r.Read()
		if err != nil {
			switch {
			case err == io.EOF:
				break FILE_READING_LOOP
			default:
				return nil, fmt.Errorf("failed to read embedded %s: %w\nretrieved record: %v", fileWithTransformedNames, err, record)

			}
		}
		id, err := strconv.ParseInt(record[0], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse id %s: %w", record[0], err)
		}
		name := record[1]
		db[id] = name
		if id > maxId {
			maxId = id
		}
	}

	return &NamesDB{
		database: db,
		maxId:    maxId,
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

func (n NamesDB) GetPage(ctx context.Context, year int64, page int64, limit int64) ([]*models.Name, error) {
	start := page * limit
	end := start + limit
	if end > n.maxId {
		end = n.maxId
	}

	// generate response
	var names []*models.Name
	for i := start; i < end; i++ {
		name, ok := n.database[i]
		if !ok {
			return nil, fmt.Errorf("name with id %d not found", i)
		}
		names = append(names, &models.Name{
			Id:    i,
			Value: name,
		})
	}

	return names, nil
}

func (n NamesDB) GetNoOfEntries(ctx context.Context, year int64) (int64, error) {
	return int64(len(n.database)), nil
}
