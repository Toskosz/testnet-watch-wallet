package db

import (
	"database/sql"
)

type DB struct {
	*sql.DB
}

func NewDB(dbPath string) (*DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Create addresses table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS addresses (
			address TEXT PRIMARY KEY,
			alias TEXT PRIMARY KEY,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

func (db *DB) StoreAddress(address string, alias string) error {
	_, err := db.Exec("INSERT INTO addresses (address, alias) VALUES (?, ?)", address, alias)
	return err
}

func (db *DB) DeleteAddress(alias string) error {
	_, err := db.Exec("DELETE FROM addresses WHERE alias = ?", alias)
	return err
}

func (db *DB) GetAddresses() ([]string, error) {
	rows, err := db.Query("SELECT address, alias FROM addresses")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var addresses []string
	for rows.Next() {
		var addr string
		if err := rows.Scan(&addr); err != nil {
			return nil, err
		}
		addresses = append(addresses, addr)
	}
	return addresses, nil
}

// GetAddress returns the actual address for a given address or alias
func (db *DB) GetAddress(input string) (string, error) {
	var address string
	err := db.QueryRow("SELECT address FROM addresses WHERE address = ? OR alias = ?", input, input).Scan(&address)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return address, err
}