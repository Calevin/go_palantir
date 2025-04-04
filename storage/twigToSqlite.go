package storage

import (
	"database/sql"
	"fmt"
	"github.com/Calevin/go_palantir/parser"
	_ "github.com/mattn/go-sqlite3"
)

// SQLiteControllerStorage gestiona el almacenamiento de registros de twig en SQLite.
type SQLiteTwigStorage struct {
	db *sql.DB
}

// NewSQLiteTwigStorage abre (o crea) la base de datos SQLite y se asegura de que la tabla exista.
func NewSQLiteTwigStorage(dbPath string) (*SQLiteTwigStorage, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Crea la tabla si no existe.
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS twig_paths (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		file TEXT,
		line INTEGER,
		path_param TEXT
	);`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		return nil, fmt.Errorf("error al crear la tabla: %v", err)
	}

	return &SQLiteTwigStorage{db: db}, nil
}

// SaveTwigRoutes inserta una lista de registros en la tabla.
func (s *SQLiteTwigStorage) SaveTwigRoutes(twigPaths []parser.TwigPathInfo) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO twig_paths(file, line, path_param) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, twigPath := range twigPaths {
		_, err = stmt.Exec(twigPath.File, twigPath.Line, twigPath.PathParam)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

// Close cierra la conexión con la base de datos.
func (s *SQLiteTwigStorage) Close() error {
	return s.db.Close()
}
