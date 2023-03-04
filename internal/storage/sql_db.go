package storage

import (
	"database/sql"

	"github.com/AlexEkdahl/govelocity/internal/process"
	_ "github.com/mattn/go-sqlite3"
)

type SqliteDatabase struct {
	Db       *sql.DB
	Filename string
}

func (s *SqliteDatabase) Open() error {
	db, err := sql.Open("sqlite3", s.Filename)
	if err != nil {
		return err
	}
	s.Db = db
	return nil
}

func OpenDatabase(filename string) (*SqliteDatabase, error) {
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil, err
	}

	storer := &SqliteDatabase{
		Db:       db,
		Filename: filename,
	}

	if err := storer.CreateTable(); err != nil {
		return nil, err
	}

	return storer, nil
}

func (s *SqliteDatabase) RemoveProcess(pid int) error {
	query := `
         DELETE FROM processes
         WHERE pid = ?
      `
	_, err := s.Db.Exec(query, pid)
	if err != nil {
		return err
	}
	return nil
}

func (s *SqliteDatabase) Close() error {
	return s.Db.Close()
}

func (s *SqliteDatabase) CreateTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS processes (
			id INTEGER PRIMARY KEY,
			name TEXT,
			pid INTEGER,
			path TEXT
		)
	`
	_, err := s.Db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (s *SqliteDatabase) InsertProcess(p *process.Process) error {
	query := `
		INSERT INTO processes (name, pid, path)
		VALUES (?, ?, ?)
	`
	_, err := s.Db.Exec(query, p.Name, p.PID, p.Path)
	if err != nil {
		return err
	}

	return nil
}

func (s *SqliteDatabase) GetProcesses() ([]*process.Process, error) {
	query := `
		SELECT id, name, pid, path
		FROM processes
	`
	rows, err := s.Db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var processes []*process.Process
	for rows.Next() {
		p := &process.Process{}
		if err := rows.Scan(&p.PID, &p.Name, &p.PID, &p.Path); err != nil {
			return nil, err
		}
		processes = append(processes, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return processes, nil
}
