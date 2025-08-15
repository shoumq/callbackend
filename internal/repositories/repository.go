package repositories

import (
	"database/sql"
	"fmt"
	"vasek/internal/dto"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type RequestRepository struct {
	db *sql.DB
}

func NewRequestRepository(db *sql.DB) *RequestRepository {
	return &RequestRepository{
		db: db,
	}
}

func (r *RequestRepository) InsertRequest(name, text, phone string) (int, error) {
	var id int
	err := r.db.QueryRow(
		"INSERT INTO requests (name, text, phone) VALUES ($1, $2, $3) RETURNING id",
		name, text, phone).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("could not insert request: %w", err)
	}
	return id, nil
}

func (r *RequestRepository) SelectRequests() ([]dto.Request, error) {
	rows, err := r.db.Query("SELECT id, name, text, phone FROM requests")
	if err != nil {
		return nil, fmt.Errorf("could not query requests: %w", err)
	}
	defer rows.Close()

	var requests []dto.Request
	for rows.Next() {
		var request dto.Request
		if err := rows.Scan(&request.ID, &request.Name, &request.Text, &request.Phone); err != nil {
			return nil, fmt.Errorf("could not scan request: %w", err)
		}
		requests = append(requests, request)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return requests, nil
}
