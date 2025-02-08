package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Hanagasumiiii/docker-track/internal/models"
	_ "github.com/lib/pq"
)

type Storage struct {
	*sql.DB
}

func Connect(path string) (*Storage, error) {
	const op = "storage.Connect"

	db, err := sql.Open("postgres", path)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS containers (
    	ID SERIAL PRIMARY KEY,
    	IP VARCHAR(255) NOT NULL,
    	STATUS VARCHAR(255) NOT NULL);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db}, nil
}

func (s *Storage) SaveContainer(container models.Container) error {
	const op = "storage.SaveContainer"

	// TODO check on repeat

	stmt, err := s.DB.Prepare(`
	INSERT INTO containers (ip, status) 
	VALUES ($1, $2)
	`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(container.Ip, container.Status)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetContainers() ([]models.Container, error) {
	const op = "storage.GetContainers"

	stmt, err := s.DB.Prepare(`
	SELECT * FROM containers
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := stmt.Query()
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("%s: %w", op, errors.New("no containers found"))
	}
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var containers []models.Container

	for rows.Next() {
		var c models.Container
		err = rows.Scan(&c.Id, &c.Ip, &c.Status)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		containers = append(containers, c)
	}

	return containers, nil
}

func (s *Storage) DeleteContainer(ip string) error {
	const op = "storage.DeleteContainer"

	stmt, err := s.DB.Prepare(`
	DELETE from containers WHERE ip = $1
	`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(ip)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
