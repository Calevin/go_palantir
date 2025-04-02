package storage

import (
	"database/sql"
	"fmt"
	"github.com/Calevin/go_palantir/parser"
	_ "github.com/mattn/go-sqlite3"
)

// SQLiteControllerStorage gestiona el almacenamiento de registros de controller en SQLite.
type SQLiteControllerStorage struct {
	db *sql.DB
}

// NewSQLiteControllerStorage abre (o crea) la base de datos SQLite y se asegura de que la tabla exista.
func NewSQLiteControllerStorage(dbPath string) (*SQLiteControllerStorage, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Crea la tabla si no existe.
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS controller_routes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		file TEXT,
		line INTEGER,
		url TEXT,
		name_url TEXT,
		method TEXT
	);`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		return nil, fmt.Errorf("error al crear la tabla: %v", err)
	}

	return &SQLiteControllerStorage{db: db}, nil
}

// SaveControllerRoutes inserta una lista de registros en la tabla.
func (s *SQLiteControllerStorage) SaveControllerRoutes(routes []parser.RouteInfo) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO controller_routes(file, line, url, name_url, method) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, route := range routes {
		_, err = stmt.Exec(route.File, route.Line, route.URL, route.NameURL, route.Method)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

// Close cierra la conexión con la base de datos.
func (s *SQLiteControllerStorage) Close() error {
	return s.db.Close()
}
