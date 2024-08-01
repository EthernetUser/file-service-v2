package postgres

import (
	"database/sql"
	"file-service/m/internal/config"
	"file-service/m/internal/database"
	"fmt"
	"sync"

	_ "github.com/lib/pq"
)

type Postgres struct {
	db   *sql.DB
	once sync.Once
}

func New(cfg config.DatabaseConfig) (*Postgres, error) {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)

	db, err := sql.Open("postgres", connString)

	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %v", err)
	}

	err = db.Ping()

	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS files (
		id SERIAL PRIMARY KEY,
		original_name TEXT NOT NULL,
		name TEXT NOT NULL UNIQUE,
		path TEXT NOT NULL,
		size BIGINT NOT NULL,
		storage_type TEXT NOT NULL,
		timestamp TIMESTAMP NOT NULL DEFAULT NOW(),
		is_deleted BOOLEAN NOT NULL DEFAULT FALSE
	);`)

	if err != nil {
		return nil, fmt.Errorf("failed to create table: %v", err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %v", err)
	}

	stmt, err = db.Prepare(`CREATE INDEX IF NOT EXISTS files_name_idx ON files (name);`)

	if err != nil {
		return nil, fmt.Errorf("failed to create index: %v", err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("failed to create index: %v", err)
	}

	return &Postgres{
		db: db,
	}, nil
}

func (p *Postgres) Close() error {
	var errorOnClose error

	p.once.Do(func() {
		err := p.db.Close()
		if err != nil {
			errorOnClose = fmt.Errorf("failed to close database: %v", err)
		}
	})

	return errorOnClose
}

func (p *Postgres) SaveFile(file database.FileToSave) (int64, error) {
	const op = "postgres.InsertFile"

	query := `INSERT INTO files (name, original_name, path, size, storage_type) VALUES ($1, $2, $3, $4, $5)`

	tx, err := p.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	defer tx.Rollback()

	stmt, err := tx.Prepare(query)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(file.Name, file.OriginalName, file.Path, file.Size, file.StorageType)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	var id int64
	stmt, err = tx.Prepare("SELECT id FROM files WHERE name = $1")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	err = stmt.QueryRow(file.Name).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	err = tx.Commit()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (p *Postgres) GetFile(id int64, isDeleted bool) (*database.File, error) {
	const op = "postgres.GetFile"

	query := `SELECT * FROM files WHERE id = $1 and is_deleted = $2`

	tx, err := p.db.Begin()
	if err != nil {
		return &database.File{}, fmt.Errorf("%s: %w", op, err)
	}

	defer tx.Rollback()

	stmt, err := tx.Prepare(query)
	if err != nil {
		return &database.File{}, fmt.Errorf("%s: %w", op, err)
	}

	var file database.File
	err = stmt.
		QueryRow(id, isDeleted).
		Scan(
			&file.Id,
			&file.OriginalName,
			&file.Name,
			&file.Path,
			&file.Size,
			&file.StrorageType,
			&file.Timestamp,
			&file.IsDeleted,
		)

	if err != nil {
		return &database.File{}, fmt.Errorf("%s: %w", op, err)
	}

	err = tx.Commit()
	if err != nil {
		return &database.File{}, fmt.Errorf("%s: %w", op, err)
	}

	return &file, nil
}

func (p *Postgres) SetFileIsDeleted(id int64) (int64, error) {

	const op = "postgres.DeleteFile"

	query := `UPDATE files SET is_deleted = true WHERE id = $1 and is_deleted = false`

	tx, err := p.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	defer tx.Rollback()

	stmt, err := tx.Prepare(query)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	r, err := stmt.Exec(id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	resultRowsAffected, err := r.RowsAffected()

	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	tx.Commit()

	if resultRowsAffected == 0 {
		return 0, fmt.Errorf("%s: %w", op, database.ErrorNotFound)
	}

	return resultRowsAffected, nil
}

func (p *Postgres) DeleteFile(id int64) (int64, error) {
	const op = "postgres.DeleteFile"

	query := `DELETE FROM files WHERE id = $1 and is_deleted = true`

	tx, err := p.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	defer tx.Rollback()

	stmt, err := p.db.Prepare(query)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	r, err := stmt.Exec(id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	resultRowsAffected, err := r.RowsAffected()

	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	tx.Commit()

	if resultRowsAffected == 0 {
		return 0, fmt.Errorf("%s: %w", op, database.ErrorNotFound)
	}

	return resultRowsAffected, nil
}
