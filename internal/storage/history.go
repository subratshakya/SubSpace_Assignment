package storage

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

const dbPath = "data/history.db"

// HistoryDB manages the SQLite database for tracking actions
type HistoryDB struct {
	db *sql.DB
}

// NewHistoryDB creates a new HistoryDB instance
func NewHistoryDB() (*HistoryDB, error) {
	// Ensure data directory exists
	if err := ensureDataDir(); err != nil {
		return nil, err
	}
	
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	
	historyDB := &HistoryDB{db: db}
	
	// Initialize schema
	if err := historyDB.initSchema(); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}
	
	return historyDB, nil
}

// ensureDataDir creates the data directory if it doesn't exist
func ensureDataDir() error {
	return os.MkdirAll("data", 0755)
}

// initSchema creates the necessary tables
func (h *HistoryDB) initSchema() error {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS profiles_processed (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		profile_url TEXT UNIQUE NOT NULL,
		action_type TEXT NOT NULL,
		processed_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	
	CREATE INDEX IF NOT EXISTS idx_profile_url ON profiles_processed(profile_url);
	CREATE INDEX IF NOT EXISTS idx_processed_at ON profiles_processed(processed_at);
	`
	
	_, err := h.db.Exec(createTableSQL)
	return err
}

// IsProcessed checks if a profile URL has already been processed
func (h *HistoryDB) IsProcessed(profileURL string) (bool, error) {
	query := `SELECT COUNT(*) FROM profiles_processed WHERE profile_url = ?`
	var count int
	err := h.db.QueryRow(query, profileURL).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// MarkProcessed marks a profile URL as processed
func (h *HistoryDB) MarkProcessed(profileURL string, actionType string) error {
	query := `INSERT OR IGNORE INTO profiles_processed (profile_url, action_type) VALUES (?, ?)`
	_, err := h.db.Exec(query, profileURL, actionType)
	return err
}

// GetProcessedCount returns the number of profiles processed today
func (h *HistoryDB) GetProcessedCount(actionType string) (int, error) {
	query := `SELECT COUNT(*) FROM profiles_processed 
	          WHERE action_type = ? AND DATE(processed_at) = DATE('now')`
	var count int
	err := h.db.QueryRow(query, actionType).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// Close closes the database connection
func (h *HistoryDB) Close() error {
	return h.db.Close()
}

