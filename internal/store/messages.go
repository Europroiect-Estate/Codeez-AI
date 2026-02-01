package store

import "time"

// Message represents a chat message.
type Message struct {
	ID        int64
	SessionID int64
	Role      string
	Content   string
	CreatedAt time.Time
}

// AppendMessage inserts a message and returns its ID.
func (s *Store) AppendMessage(sessionID int64, role, content string) (int64, error) {
	res, err := s.DB.Exec(
		`INSERT INTO messages (session_id, role, content) VALUES (?, ?, ?)`,
		sessionID, role, content,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// GetMessages returns all messages for a session in order.
func (s *Store) GetMessages(sessionID int64) ([]Message, error) {
	rows, err := s.DB.Query(
		`SELECT id, session_id, role, content, created_at FROM messages WHERE session_id = ? ORDER BY id`,
		sessionID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Message
	for rows.Next() {
		var m Message
		var createdAt string
		if err := rows.Scan(&m.ID, &m.SessionID, &m.Role, &m.Content, &createdAt); err != nil {
			return nil, err
		}
		m.CreatedAt = parseTime(createdAt)
		out = append(out, m)
	}
	return out, rows.Err()
}
