package store

import (
	common_db "common/db"
	"context"
	"database/sql"
	"fmt"
)

type PingPongStore struct {
	dbService *common_db.DBService
}

func NewPingPongStore(db *common_db.DBService) *PingPongStore {
	return &PingPongStore{
		dbService: db,
	}
}

func (ps *PingPongStore) GetCurr() (int, error) {
	var count int
	query := `
	SELECT count
	FROM pingpong_counter
	WHERE id = 1
	`

	err := ps.dbService.DB.QueryRow(query).Scan(&count)
	if err == sql.ErrNoRows {
		return -1, fmt.Errorf("No rows in pingpong_count.count for row with id 1")
	}

	if err != nil {
		return -1, err
	}

	return count, nil
}

func (ps *PingPongStore) Update() (int, error) {
	query := `
	INSERT INTO pingpong_counter (id, count) VALUES(1, 1)
	ON CONFLICT (id) DO UPDATE
	SET count = pingpong_counter.count + 1
	RETURNING count
	`

	var newCount int
	err := ps.dbService.DB.QueryRowContext(context.Background(), query).Scan(&newCount)
	if err != nil {
		return -1, err
	}

	return newCount, nil
}
