package repository

import (
	"database/sql"
)

// SessionRepository handles database operations for study sessions.
type SessionRepository struct {
	DB *sql.DB
}

// NewSessionRepository creates a new repository instance.
func NewSessionRepository(db *sql.DB) *SessionRepository {
	return &SessionRepository{DB: db}
}

// InitSchema ensures the necessary tables are created upon startup.
func (r *SessionRepository) InitSchema() error {
	query := `
	CREATE TABLE IF NOT EXISTS study_sessions (
		id SERIAL PRIMARY KEY,
		student_id VARCHAR(50) NOT NULL,
		date DATE NOT NULL,
		subject VARCHAR(100) NOT NULL,
		duration_minutes INT NOT NULL,
		client_session_id VARCHAR(100) UNIQUE NOT NULL, -- Serves as the idempotency key
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`
	_, err := r.DB.Exec(query)
	return err
}

// SaveSession persists a study session into the database.
// It gracefully ignores duplicate entries by utilizing the unique client_session_id.
func (r *SessionRepository) SaveSession(studentID, date, subject string, duration int, clientSessionID string) error {
	query := `
	INSERT INTO study_sessions (student_id, date, subject, duration_minutes, client_session_id)
	VALUES ($1, $2, $3, $4, $5)
	ON CONFLICT (client_session_id) DO NOTHING;` // Ensures idempotency at the DB level

	_, err := r.DB.Exec(query, studentID, date, subject, duration, clientSessionID)
	return err
}
