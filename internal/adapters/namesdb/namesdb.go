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

var ErrYearNotFound = errors.New("year not found")
var ErrNameNotFound = errors.New("name not found")

type Entries map[int64]string
type YearDB struct {
	Entries
	maxId int64
}

type NamesDB struct {
	database map[int64]*YearDB
	years    map[int64]struct{}
}

func NewNamesDB() (*NamesDB, error) {
	namesDB := &NamesDB{
		database: make(map[int64]*YearDB),
		years:    make(map[int64]struct{}),
	}

	// open embedded file
	fs, err := namesEmbedded.Open(fileWithTransformedNames)
	if err != nil {
		return nil, fmt.Errorf("failed to open embedded %s: %w", fileWithTransformedNames, err)
	}
	r := csv.NewReader(fs)

	// read the embedded file line by line
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
		// year
		year := record[0]
		yearInt, err := strconv.ParseInt(year, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse year %s: %w", year, err)
		}
		yearDB, exists := namesDB.database[yearInt]
		if !exists {
			yearDB = &YearDB{
				Entries: make(map[int64]string),
				maxId:   0,
			}
			namesDB.years[yearInt] = struct{}{}
			namesDB.database[yearInt] = yearDB
		}
		// id
		id, err := strconv.ParseInt(record[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse id %s: %w", record[0], err)
		}
		// name
		name := record[2]

		// write entry
		yearDB.Entries[id] = name
		if id > yearDB.maxId {
			yearDB.maxId = id
		}
	}

	return namesDB, nil
}

func (n NamesDB) GetName(ctx context.Context, year int64, id int64) (*models.Name, error) {
	yearDB, ok := n.database[year]
	if !ok {
		return nil, ErrYearNotFound
	}
	name, ok := yearDB.Entries[id]
	if !ok {
		return nil, ErrNameNotFound
	}
	return &models.Name{
		Id:    id,
		Value: name,
	}, nil
}

func (n NamesDB) GetPage(ctx context.Context, year int64, page int64, limit int64) ([]*models.Name, error) {
	yearDB, ok := n.database[year]
	if !ok {
		return nil, ErrYearNotFound
	}
	start := page * limit
	end := start + limit
	if end > yearDB.maxId {
		end = yearDB.maxId
	}

	// generate response
	var names []*models.Name
	for i := start; i < end; i++ {
		name, ok := yearDB.Entries[i]
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

func (n NamesDB) GetYearsAvailable(ctx context.Context) (map[int64]struct{}, error) {
	return n.years, nil
}

func (n NamesDB) GetNoOfEntries(ctx context.Context, year int64) (int64, error) {
	return int64(len(n.database[year].Entries)), nil
}
